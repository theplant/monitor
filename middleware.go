package monitor

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

// OperationMonitor returns Gin middleware that logs requests into
// an InfluxDB database.
func OperationMonitor() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		defer func() {
			interval := time.Now().Sub(start)
			localCtx := ctx.Copy()
			go func() {
				tags := tagsForContext(localCtx)
				recordOperation("request", start, interval, tags)
			}()
		}()

		ctx.Next()
	}
}

// recordOperation is used in OperationMonitor for logging
// http requests into monitor.
func recordOperation(measurement string, start time.Time, duration time.Duration, tags map[string]string) {
	insertRecord(measurement, float64(duration/time.Millisecond), tags, start)
}

func tagsForContext(ctx *gin.Context) map[string]string {
	return map[string]string{
		"path":           scrubPath(ctx.Request.URL.Path),
		"request_method": ctx.Request.Method,
		"response_code":  fmt.Sprint(ctx.Writer.Status()),
	}
}

var idScrubber = regexp.MustCompile("[0-9]+")

func scrubPath(path string) string {
	return idScrubber.ReplaceAllString(path, ":id")
}
