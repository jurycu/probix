package gin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/fvbock/endless"
	ginlimits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	bearerPrefix = "Bearer "
)

//TunningWebServer tbd
func TunningWebServer() *gin.Engine {

	//gin配置
	router := gin.New()
	//基础中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(Metrics())
	router.Use(ginlimits.RequestSizeLimiter(33554432))
	router.NoMethod(NoMethodHandler())
	router.NoRoute(NoRouteHandler())

	return router
}

//ConfigureRoute tbd
func ConfigureRoute(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "health")
	})

	router.GET("/probix", GetMetrics)
	// metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

}

//StartWebServer
func StartWebServer(router *gin.Engine) error {
	//server微调
	endless.DefaultReadTimeOut = time.Duration(180) * time.Second
	endless.DefaultWriteTimeOut = time.Duration(180) * time.Second
	endless.DefaultMaxHeaderBytes = 33554432
	return endless.ListenAndServe(":"+strconv.Itoa(9098), router)
}

// NoRouteHandler 未找到请求路由的处理函数
func NoRouteHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(http.StatusBadRequest, "NoRouteHandler error")
	}
}

// NoMethodHandler 未找到请求方法的处理函数
func NoMethodHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(http.StatusBadRequest, "NoMethodHandler error")
	}
}
