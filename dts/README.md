# DTS - is a very light Data Transformation Service with Rest API

### Written on golang 1.18.x
<img src="https://github.com/dn-smirnov/golang-labs/blob/master/dts/dts-01.png" width=75% height=75%>

Based on Gin Web Framework v1.8.1

Ready to MySQL connection for reading data from MySQL directly with go-sql-driver/mysql v1.6.0

Configuring through .ini file with ini.v1 v1.66.6

Embeded Prometheus exporter go-gin-prometheus v0.1.0

You can construct custom SQL queries for sending it to Prometheus server from your project.

Installation:

- Clone repo ```bash# git clone https://github.com/dn-smirnov/golang-labs.git```

- Enter project folder e.t. ```bash# cd golang-labs/dts```

- run: ```bash# go build```

- Correct dts.ini on listen and allowed_from

- configure MySQL on existing server and testing database in dts.ini

- run ```bash#./dts```

- done!

- Type in your browser http://$YOUR_IP:8080/metrics or empty test point http://$YOUR_IP:8080/ping 

- You can modify or adding inside Routes section in bottom of main.go wish you want
