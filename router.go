package router

import (
	"github.com/gin-gonic/gin"
	"my-app/api/src/handler"
)

func LandRouter(r *gin.Engine) {
	api := r.Group("/api")
	{
		user := api.Group("/user")
		{
			user.POST("/register", handler.Register)
			user.POST("/login", handler.Login)
		}
	}
}
