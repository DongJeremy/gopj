package server

import (
	"fmt"

	"github.com/davidddw/gopj/ginMVC/configs"
	"github.com/davidddw/gopj/ginMVC/routes"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

// Run 配置并动服务
func Run(httpServer *gin.Engine, db *xorm.Engine) {
	// 服务配置
	serverConfig := configs.SystemConfig.Server

	// gin 运行时 release debug test
	gin.SetMode(serverConfig.Env)

	httpServer = gin.Default()

	// 配置视图
	if "" != serverConfig.ViewPattern {
		httpServer.LoadHTMLGlob(serverConfig.ViewPattern)
	}

	if "" != serverConfig.StaticPattern {
		httpServer.Static("/vender/static", serverConfig.StaticPattern)
	}

	// 注册路由
	routes.RegisterRoutes(httpServer, db)

	serverAddr := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)

	// 启动服务
	err := httpServer.Run(serverAddr)
	if nil != err {
		panic("server n error: " + err.Error())
	}
}
