package gin

//Run
func Run() error {

	//web框架配置与插件
	router := TunningWebServer()
	//路由设置
	ConfigureRoute(router)
	//启动server
	return StartWebServer(router)

}
