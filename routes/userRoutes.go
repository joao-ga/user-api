package routes

import (
	"github.com/gin-gonic/gin"
	"user-api/controllers"
)

func UserRoutes(router *gin.Engine) {
	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/", controllers.GetAllUsers)
		userRoutes.GET("/:id", controllers.GetUserById)
		userRoutes.POST("/", controllers.CreateUser)
		userRoutes.PUT("/:id", controllers.UpdateUser)
		userRoutes.DELETE("/:id", controllers.DeleteUser)
		userRoutes.GET("/sendemails", controllers.TestSendBirthdayEmails)
	}
}
