package routes

import (
	"github.com/davidddw/gopj/ginMvc/controllers"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

func RegisterRoutes(router *gin.Engine, db *xorm.Engine) {
	router.GET("/", controllers.IndexHome)
	router.GET("/index", controllers.IndexHome)
	router.GET("/users/:id", controllers.UserGet(db))
	router.GET("/users", controllers.UserGetList(db))
	router.POST("/users", controllers.UserPost(db))
	router.PUT("/users/:id", controllers.UserPut(db))
	// router.PATCH("/users/:id", controllers.UserPut)
	router.DELETE("/users/:id", controllers.UserDelete(db))
}
