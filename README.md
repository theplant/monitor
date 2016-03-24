# Monitor

[![Build Status](https://semaphoreci.com/api/v1/projects/8a3f5364-f25d-40cd-9613-d078cbd577f0/747395/badge.svg)](https://semaphoreci.com/theplant/monitor)

Generic monitor that integrated with [influxdb](https://influxdata.com/) for [Gin-backed](https://gin-gonic.github.io/gin/) web application.

# Quick Start

## InfluxDB configuration

Set up `INFLUXDB_URL` environment variable.

```sh
$ export INFLUXDB_URL="https://<INFLUX_USER>:<INFLUX_PW>@influxdb.com/"
```


## Set up Gin middleware

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/theplant/monitor"
)

func main() {
    r := gin.Default()

    // Set up operation middleware
    r.Use(monitor.OperationMonitor())

    r.GET("/ping", func(c *gin.Context) {

        // Custom measurement with monitor.Count
        monitor.Count("ping", 1, map[string]string{"ping": "pong"})

        c.JSON(200, gin.H{"message": "pong"})
    })
    r.Run() // listen and server on 0.0.0.0:8080
}
```

After you start with `go run main.go`. You'll see:

```sh
Logging operations to Influxdb "<your-influxdb>" database
```

Then it works. :)
