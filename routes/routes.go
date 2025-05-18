package routes

import (
	"fadlan/backend-api/controllers"
	"fadlan/backend-api/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	//initialize gin
	router := gin.Default()

	// route register
	router.POST("/api/register", controllers.Register)

	// route login
	router.POST("/api/login", controllers.Login)

	// route users
	router.GET("/api/users", middlewares.AuthMiddleware(), controllers.FindUsers)

	// route create user
	router.POST("/api/users", middlewares.AuthMiddleware(), controllers.CreateUser)

	return router
}