package main

import (
	"fmt"
	"time"
	"os"
	"io"
	"log"
	"net/http"
	"gopkg.in/ini.v1"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gin-gonic/gin"
	"github.com/zsais/go-gin-prometheus"
)

func main() {

	// Configuration load
	
	cfg, err := ini.Load("dts.ini")
		
    if err != nil {
        fmt.Println("Fail to read file: %v", err)
        os.Exit(1)
    }
    
    // stdout DTS runner
    
    fmt.Printf("\n-----------------------------------Starting API Server ...----------------------------------------\n\n")
    fmt.Println("App Mode:", cfg.Section("").Key("app_mode").String())
    fmt.Println("Host ID:", cfg.Section("").Key("host_id").String())
    fmt.Printf("\n-----------------------------------Starting API Server OK!----------------------------------------\n\n")
	
	// MySQL params
	
    sqluser := cfg.Section("mysql").Key("user").String()
    sqlpassword := cfg.Section("mysql").Key("password").String()
    sqlserver := cfg.Section("mysql").Key("server").String()
    sqlport := cfg.Section("mysql").Key("port").String()
    sqldatabase := cfg.Section("mysql").Key("database").String()
    sqlconn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", sqluser, sqlpassword,sqlserver,sqlport,sqldatabase)
  
    db, err := sql.Open("mysql", sqlconn)
    defer db.Close()

    if err != nil {
        log.Fatal(err)
    }
    
    var mysql_version string
    
    // Check MySQL connection
    
    err2 := db.QueryRow("SELECT VERSION()").Scan(&mysql_version)
    fmt.Println("MySQL Connection: " + mysql_version + "\n")
	
    if err2 != nil {
        log.Fatal(err2)
    }
	
	// Read app parameters
	
	app_mode := cfg.Section("").Key("app_mode").String()
	app_listen := cfg.Section("server").Key("listen")
	app_port := cfg.Section("server").Key("http_port")
	app_param := fmt.Sprintf( "%s:%s" ,app_listen,app_port)
	app_allow_from := cfg.Section("server").Key("allow_from").String()
	
	// DTS Mode

	if  app_mode == "development" {
		gin.SetMode(gin.DebugMode)
	} else if app_mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	
	// Logging
	
	log_file := cfg.Section("log").Key("file").String()
	log_console := cfg.Section("log").Key("console").String()
	log_dir := cfg.Section("log").Key("path").String()
	log_path := cfg.Section("log").Key("path").String() + "dts-" + time.Now().Format("20060102150405") + ".log"
	
	// If log dir not exist, creating it
	
	if _, err := os.Stat(log_dir); os.IsNotExist(err) {
    	err := os.Mkdir(log_dir, 0700)
		if err != nil {
       		log.Fatal(err)
    	}
	}
	
	// console color output, file logging without colors, mixed mode, changed via dts.ini
	
	if log_file == "yes" && log_console == "no" {
		gin.DisableConsoleColor()
		f, _ := os.Create(log_path)
		gin.DefaultWriter = io.MultiWriter(f)
	}  else if log_file == "no" && log_console == "yes" {
		gin.ForceConsoleColor()
		gin.DefaultWriter = io.MultiWriter(os.Stdout)
	} else {
		gin.ForceConsoleColor()
		f, _ := os.Create(log_path)
		gin.DefaultWriter = io.MultiWriter(f,os.Stdout)
	}
	
	// Run statement
	
	router := gin.Default()
	router.SetTrustedProxies([]string{app_allow_from})
	
	// Prometheus exporter
	
	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)
	
	// Routes
	
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, "Hello world!")
	})
	
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	
	router.Run(app_param)
}
