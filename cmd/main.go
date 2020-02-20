package main

import (
	"animated-robot/storage"
	"animated-robot/tools"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log := SetupDefaultLogger()
	sessionStore := storage.NewSessionStoreInMemory(tools.NewCodeGenerator())
	socketStore := storage.NewSocketStoreInMemory()

	socketFactory := NewSocketFactory(socketStore, sessionStore, log)
	socket := socketFactory.New()

	middleware := NewMiddlewarePipeline(log)
	server := NewServer(middleware, socket, log)

	server.Run(port)
}
