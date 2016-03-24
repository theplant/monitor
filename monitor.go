// Package monitor is a monitor "provider" that provides
// a way for monitoring. It uses influxdb monitor by default.
package monitor

import (
	"log"
	"time"
)

// Monitor defines an interface for inserting record.
type Monitor interface {
	InsertRecord(string, interface{}, map[string]string, time.Time)
}

// Monit stores a Monitor that is used in *all* log functions.
var Monit Monitor

func init() {
	influxdbMonitor := initInfluxdbMonitor()

	if influxdbMonitor != nil {
		Monit = influxdbMonitor
	} else {
		log.Println("[WARNING] No Influxdb Monitor. Please specify a Monitor to monitor.Monit.")
	}
}

// CountError logs a value in measurement, with the given error's
// message stored in an `error` tag.
func CountError(measurement string, value float64, err error) {
	data := map[string]string{"error": err.Error()}
	Count(measurement, value, data)
}

// CountSimple logs a value in measurement (with no tags).
func CountSimple(measurement string, value float64) {
	Count(measurement, value, nil)
}

// Count logs a value in measurement with given tags.
func Count(measurement string, value float64, tags map[string]string) {
	insertRecord(measurement, value, tags, time.Now())
}

// insertRecord makes all monitor log functions *safe*.
func insertRecord(measurement string, value interface{}, tags map[string]string, at time.Time) {
	if Monit != nil {
		Monit.InsertRecord(measurement, value, tags, at)
	}
	return
}
