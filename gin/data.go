package gin

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/klog"
)

func GetMetrics(c *gin.Context) {
	params := c.Request.URL.Query()
	target := params.Get("target")
	if target == "" {
		http.Error(c.Writer, "Target parameter is missing", http.StatusBadRequest)
		return
	}
	method := "GET"
	if m := params.Get("method"); m != "" {
		method = m
	}
	body := params.Get("body")
	metricsName := params.Get("metricsName")
	metricsHelp := params.Get("metricsHelp")
	if metricsName == "" || metricsHelp == "" {
		http.Error(c.Writer, "metricsName or metricsHelp is null", http.StatusBadRequest)
		klog.Error(fmt.Errorf("metricsName or metricsHelp is null"))
		return
	}

	probixSuccessGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "probix_success",
		Help: "Displays whether or not the probix was a success",
	},
		[]string{"method"},
	)
	probixDurationGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "probix_duration_seconds",
		Help: "Returns how long the probe took to complete in seconds",
	},
		[]string{"method"},
	)
	start := time.Now()
	registry := prometheus.NewRegistry()
	registry.MustRegister(probixSuccessGauge)
	registry.MustRegister(probixDurationGauge)
	success, err := requestHTTP(target, method, body, metricsName, metricsHelp, registry)
	duration := time.Since(start).Seconds()
	if err != nil {
		http.Error(c.Writer, "Target request is failed,err:"+err.Error(), http.StatusBadRequest)
		klog.Error(err)
		return
	}
	probixDurationGauge.With(prometheus.Labels{"method": method}).Set(duration)
	if success {
		probixSuccessGauge.With(prometheus.Labels{"method": method}).Set(1)
		klog.Info("probix success,duration seconds:", duration)
	} else {
		probixSuccessGauge.With(prometheus.Labels{"method": method}).Set(0)
		klog.Info("probix failed,duration seconds:", duration)
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(c.Writer, c.Request)
}
