package gin

import (
	"bytes"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.MustRegister(RequestCounter)
	prometheus.MustRegister(RequestErrCounter)
	prometheus.MustRegister(RequestLatencies)
	prometheus.MustRegister(RequestLatenciesSummary)

}

func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		resource := c.FullPath()
		method := c.Request.Method
		c.Next()

		code := c.Writer.Status()
		if code == http.StatusOK {
			MonitorCount(resource, method, strconv.Itoa(code))
			MonitorHistogram(resource, method, strconv.Itoa(code), start)
			MonitorSummary(resource, method, strconv.Itoa(code), start)
		} else {
			MonitorCount(resource, method, strconv.Itoa(code))
			MonitorErrCount(resource, method, strconv.Itoa(code))
		}
	}
}

type SimplebodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w SimplebodyWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}
