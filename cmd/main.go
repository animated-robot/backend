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
	server := socketFactory.New()

	go server.Serve()
	defer server.Close()

	RouteAndListen(server, port)
}
