package main

import (
	"animated-robot/cmd/socket"
	"animated-robot/storage"
	"animated-robot/tools"
)

func main() {
	config := MustGetEnvVars()

	log := SetupDefaultLogger(config.LOG_LEVEL, config.LOG_FILE_PATH)
	sessionStore := storage.NewSessionStoreInMemory(tools.NewCodeGenerator())
	socketStore := storage.NewSocketStoreInMemory()
	uuidGenerator := tools.NewUUIDGenerator()

	socketFactory := socket.NewSocketFactory(socketStore, sessionStore, uuidGenerator, log)
	socket := socketFactory.New()

	middleware := NewMiddlewarePipeline(log)
	server := NewServer(middleware, socket, log)

	server.Run(config.PORT)
}
