package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/panggggg/golang-project/controllers"
)

func BlogRoute(router *gin.Engine) {
	router.GET("/blogs", controllers.GetAllBlogs())
	router.POST("/blog", controllers.CreateBlog())
	router.GET("/blog/:blogId", controllers.GetABlog())
	router.PUT("/blog/:blogId", controllers.UpdateBlog())
	router.DELETE("/blog/:blogId", controllers.DeleteABlog())
}