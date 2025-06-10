package router

import (
	"github.com/gin-gonic/gin"
	"video_develop-main/video-api/api/handler"
)

func LandRouter(r *gin.Engine) {
	api := r.Group("/api")
	{
		user := api.Group("/user")
		{
			user.POST("/send", handler.Send)
			user.POST("/login", handler.Login)
		}
		video := api.Group("/video")
		{
			video.POST("/video-add", handler.VideoAdd)
			video.POST("/video-list", handler.VideoList)
		}
		comment := api.Group("/comment")
		{
			comment.POST("/add-comment", handler.CommentAdd)
			comment.POST("/update-comment", handler.CommentUpdate)
			comment.POST("/delete-comment", handler.CommentDelete)
			comment.POST("/list-comment", handler.CommentList)
			comment.POST("/anthor-check", handler.AnthorCheck)
		}
	}

}
