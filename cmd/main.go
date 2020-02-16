package main

import (
	"animated-robot/storage"
	"animated-robot/tools"
	"log"
	"net/http"
)

func main() {
	port := "8080"

	logger := SetupDefaultLogger()
	sessionStore := storage.NewSessionStoreInMemory(tools.NewCodeGenerator())
	socketStore := storage.NewSocketStoreInMemory()

	socket := NewSocket(socketStore, sessionStore, logger)
	go socket.Serve()
	defer socket.Close()

	http.Handle("/socket.io/", socket)
	http.Handle("/", http.FileServer(http.Dir("public")))
	log.Println("Serving at localhost:"+ port +"...")
	log.Fatal(http.ListenAndServe(":" + port, nil))
}