package routes

import (
	"fadlan/backend-api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	//initialize gin
	router := gin.Default()

	// route register
	router.POST("/api/register", controllers.Register)

	// route login
	router.POST("/api/login", controllers.Login)

	return router
}