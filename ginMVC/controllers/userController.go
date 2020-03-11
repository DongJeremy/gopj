package controllers

import (
	"net/http"
	"strconv"

	"github.com/davidddw/gopj/ginMVC/models"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

// get one
func UserGet(db *xorm.Engine) func(c *gin.Context) {
	return func(ctx *gin.Context) {
		id, _ := strconv.Atoi(ctx.Param("id"))
		userModel := models.User{}

		data, err := userModel.UserGet(db, id)

		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"msg": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	}

}

// get list
func UserGetList(db *xorm.Engine) func(c *gin.Context) {
	return func(ctx *gin.Context) {

		page := ctx.DefaultQuery("page", "1")
		pageSize := ctx.DefaultQuery("page_size", "10")

		pageInt, _ := strconv.Atoi(page)
		pageSizeInt, _ := strconv.Atoi(pageSize)

		userModel := models.User{}

		users, err := userModel.UserGetList(db, pageInt, pageSizeInt)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"msg":  "success",
			"data": users,
		})
	}
}

// add one
func UserPost(db *xorm.Engine) func(c *gin.Context) {
	return func(ctx *gin.Context) {

		name := ctx.PostForm("name")
		address := ctx.PostForm("address")
		age, _ := strconv.Atoi(ctx.PostForm("ag"))
		userModel := models.User{Name: name, Address: address, Age: age}
		if err := ctx.ShouldBind(&userModel); nil != err {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		id, err := userModel.UserAdd(db)

		if nil != err {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"msg": "success",
			"uid": id,
		})
	}

}

// update
func UserPut(db *xorm.Engine) func(c *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		idInt, err := strconv.Atoi(id)

		if nil != err || 0 == idInt {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "resource identifier not found",
			})
			return
		}

		userModel := models.User{}

		if err := ctx.ShouldBind(&userModel); nil != err {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		_, err = userModel.UserUpdate(db, idInt)

		if nil != err {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		// 更新成功返回 204
		ctx.JSON(http.StatusNoContent, gin.H{})
	}

}

// delete
func UserDelete(db *xorm.Engine) func(c *gin.Context) {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		//id := ctx.PostForm("id")
		idInt, err := strconv.Atoi(id)

		if nil != err || 0 == idInt {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": "resource identifier not found",
			})
			return
		}

		userModel := models.User{}

		_, err = userModel.UserDelete(db, idInt)

		if nil != err {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		// 删除成功返回 204
		ctx.JSON(http.StatusNoContent, gin.H{})
	}

}
