package router

import (
	"github.com/LLIEPJIOK/forum/internal/controller"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
}

func New(ctrl *controller.Controller) *Router {
	eng := gin.Default()

	user := eng.Group("/user")
	user.POST("", ctrl.AddUser)
	user.GET(":id", ctrl.GetUser)
	user.GET("/list/", ctrl.GetAllUsers)
	user.PUT(":id", ctrl.UpdateUser)
	user.DELETE(":id", ctrl.DeleteUser)

	post := eng.Group("/post")
	post.POST("", ctrl.AddPost)
	post.GET(":id", ctrl.GetPost)
	post.GET("/list/", ctrl.GetAllPosts)
	post.PUT(":id", ctrl.UpdatePost)
	post.DELETE(":id", ctrl.DeletePost)

	message := eng.Group("/message")
	message.POST("", ctrl.AddMessage)
	message.GET(":id", ctrl.GetMessage)
	message.GET("/list/", ctrl.GetAllMessages)
	message.PUT(":id", ctrl.UpdateMessage)
	message.DELETE(":id", ctrl.DeleteMessage)

	return &Router{
		engine: eng,
	}
}

func (r *Router) Run(address string) {
	r.engine.Run(address)
}
