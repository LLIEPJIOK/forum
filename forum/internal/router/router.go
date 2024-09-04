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
	user.PUT(":id", ctrl.UpdateUser)
	user.DELETE(":id", ctrl.DeleteUser)

	return &Router{
		engine: eng,
	}
}

func (r *Router) Run(address string) {
	r.engine.Run(address)
}
