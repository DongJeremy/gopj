package main

import (
	"github.com/davidddw/gopj/ginMVC/configs"
	"github.com/davidddw/gopj/ginMvc/drivers"
	"github.com/davidddw/gopj/ginMvc/server"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

var httpServer *gin.Engine

func main() {

	var db *xorm.Engine
	serverConfig := configs.SystemConfig.Server
	switch serverConfig.Dbtype {
	case "mysql":
		db, _ = drivers.InitMySQL()
	case "oracle":
		db, _ = drivers.InitOracle()
	case "sqlite":
		db, _ = drivers.InitSqlite()
	}

	defer db.Close()
	// 启动服务
	server.Run(httpServer, db)
}
