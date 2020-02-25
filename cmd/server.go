package main

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
	"net/http"
)

type IServer interface {
	Run(port string)
}

type Server struct {
	middleware Middleware
	socket     *socketio.Server
	log        *logrus.Logger
}

func NewServer(middleware Middleware, socket *socketio.Server, logger *logrus.Logger) *Server {
	return &Server{
		middleware: middleware,
		socket:     socket,
		log:        logger,
	}
}

func (s *Server) Run(port string) {
	go s.socket.Serve()
	defer s.socket.Close()

	http.Handle("/socket.io/", s.middleware(s.socket))
	s.log.Println("Serving at :"+ port +"...")
	s.log.Fatal(http.ListenAndServe(":" + port, nil))
}
