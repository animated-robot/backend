package main

import (
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

func RouteAndListen(s *socketio.Server, port string) {
	router := gin.New()
	router.Use(GinMiddleware())
	router.GET("/socket.io/*any", gin.WrapH(s))
	router.POST("/socket.io/*any", gin.WrapH(s))

	router.Run(":" + port)
}
