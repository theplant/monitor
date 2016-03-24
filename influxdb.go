package monitor

import (
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	influxdb "github.com/influxdata/influxdb/client"
)

// initInfluxdbMonitor reads influxdb client configuration from
// environment and returns a InfluxdbMonitor that wraps the
// configuration.
//
// Will returns nil if configuration not found or is invalid.
func initInfluxdbMonitor() *InfluxdbMonitor {
	influxdbURL := os.Getenv("INFLUXDB_URL")

	if influxdbURL == "" {
		return nil
	}

	monitor, err := newInfluxdbMonitor(influxdbURL)

	if err != nil {
		log.Println("Monitor Error", err)
		return nil
	}

	return monitor
}

// newInfluxdbMonitor creates new monitoring influxdb client. monitorURL syntax is
// `https://<username>:<password>@<influx DB host>/<database>`
//
// Will returns a error if monitorURL is invalid.
func newInfluxdbMonitor(monitorURL string) (*InfluxdbMonitor, error) {
	u, err := url.Parse(monitorURL)
	if err != nil || !u.IsAbs() {
		log.Println("InfluxDB URL Parse Error", err)
		return nil, err
	}

	monitor := InfluxdbMonitor{
		database: strings.TrimLeft(u.Path, "/"),
		cfg: &influxdb.Config{
			URL: *u,
		},
	}

	// NewClient always returns a nil error
	client, _ := influxdb.NewClient(*monitor.cfg)

	// Ignore duration, version
	_, _, err = client.Ping()
	if err != nil {
		log.Println("Influx Error", err)
	} else {
		log.Printf("Logging operations to Influxdb %q database\n", monitor.database)
	}

	return &monitor, nil
}

// InfluxdbMonitor implements monitor.Monitor interface, it wraps
// the influxdb client configuration.
type InfluxdbMonitor struct {
	cfg      *influxdb.Config
	database string
}

// InsertRecord part of monitor.Monitor.
func (im InfluxdbMonitor) InsertRecord(measurement string, value interface{}, tags map[string]string, at time.Time) {
	// NewClient always returns a nil error
	client, _ := influxdb.NewClient(*im.cfg)

	// Ignore response, we only care about write errors
	_, err := client.Write(influxdb.BatchPoints{
		Database: im.database,
		Points: []influxdb.Point{
			{
				Measurement: measurement,
				Fields: map[string]interface{}{
					"value": value,
				},
				Tags: tags,
				Time: at,
			},
		},
	})

	if err != nil {
		log.Println("Influx Error:", err)
	}
}
