package main

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type IServer interface {
	Run(port string)
}

type Server struct {
	startupTime time.Time
	middleware Middleware
	socket     *socketio.Server
	log        *logrus.Logger
}

func NewServer(middleware Middleware, socket *socketio.Server, logger *logrus.Logger) *Server {
	return &Server{
		startupTime: time.Now(),
		middleware: middleware,
		socket:     socket,
		log:        logger,
	}
}

func (s *Server) Run(port string) {
	go s.socket.Serve()
	defer s.socket.Close()

	http.HandleFunc("/healthcheck", s.healthcheck)
	http.Handle("/socket.io/", s.middleware(s.socket))
	s.log.Println("Serving at :"+ port +"...")
	s.log.Fatal(http.ListenAndServe(":" + port, nil))
}

func (s *Server) healthcheck(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Application started at " + s.startupTime.Format(time.RFC3339)))
}
