package main

import (
	"github.com/gin-gonic/gin"
	"github.com/panggggg/golang-project/configs"
	"github.com/panggggg/golang-project/routes"
)

func main() {
	router := gin.Default()

	//run database
	configs.ConnectDB()

	//routes
	routes.BlogRoute(router)

	router.Run("localhost:6000")
}