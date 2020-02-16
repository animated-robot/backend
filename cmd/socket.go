package main

import (
	"animated-robot/domain"
	"animated-robot/storage"
	"encoding/json"

	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
)

func NewSocket(socketStore storage.SocketStoreInterface, sessionStore storage.SessionStoreInterface, log *logrus.Logger) *socketio.Server {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.WithFields(logrus.Fields{
			"message": err.Error(),
		}).Fatal("NewSocket: Error creating socket")
	}

	///////////
	// FRONT //
	///////////

	frontNsp := "/front"
	server.OnConnect(frontNsp, func(s socketio.Conn) error {
		s.SetContext("")

		log.WithFields(logrus.Fields{
			"socketId": s.ID(),
		}).Info("FrontConnection: OnConnect: connection stablished")

		return nil
	})
	server.OnEvent(frontNsp, "create_session", func(s socketio.Conn, inputPlayerJson string) {
		id, err := socketStore.Store(s)
		if err != nil {
			log.WithFields(logrus.Fields{
				"socketId": id,
				"error":    err.Error(),
			}).Error("FrontConnection: OnConnect: Socket storing failed")
		}
		log.WithFields(logrus.Fields{
			"socketId": id,
		}).Trace("FrontConnection: OnConnect: socket stored")

		sessionCode, _ := sessionStore.Create(s.ID())
		log.WithFields(logrus.Fields{
			"sessionCode":   sessionCode,
			"frontSocketId": s.ID(),
		}).Info("FrontConnection: OnConnect: Session Created")

		s.Emit("created_session", sessionCode)
	})
	server.OnError(frontNsp, func(e error) {
		log.WithFields(logrus.Fields{
			"message": e.Error(),
		}).Error("FrontConnection: OnError: meet error", e)
	})
	server.OnDisconnect(frontNsp, func(s socketio.Conn, reason string) {
		log.WithFields(logrus.Fields{
			"reason": reason,
		}).Info("FrontConnection: OnDisconnect: disconnected")
	})

	///////////
	// INPUT //
	///////////

	inputNsp := "/input"
	server.OnConnect(inputNsp, func(s socketio.Conn) error {
		log.WithFields(logrus.Fields{
			"socketId": s.ID(),
		}).Info("InputConnection: OnConnect: connection stablished")

		return nil
	})
	server.OnEvent(inputNsp, "register_player", func(s socketio.Conn, inputPlayerJson string) {
		var registerPlayer domain.RegisterPlayer
		err := json.Unmarshal([]byte(inputPlayerJson), &registerPlayer)
		if err != nil {
			log.WithFields(logrus.Fields{
				"input": inputPlayerJson,
				"event": "register_player",
			}).Error("InputConnection: OnEvent: Could not parse from json string")
		}

		player := storage.NewPlayer(registerPlayer.PlayerName)
		err = sessionStore.AddPlayer(registerPlayer.SessionCode, player)
		if err != nil {
			log.WithFields(logrus.Fields{
				"player": player,
			}).Error("Could not store player")
		}
	})
	server.OnEvent(inputNsp, "context", func(s socketio.Conn, inputContextJson string) {
		var inputContext domain.InputContext
		err := json.Unmarshal([]byte(inputContextJson), &inputContext)
		if err != nil {
			log.WithFields(logrus.Fields{
				"input": inputContextJson,
				"event": "context",
			}).Error("InputConnection: OnEvent: Could not parse from json string")
		}

		session, err := sessionStore.Get(inputContext.SessionCode)
		if err != nil {
			log.WithFields(logrus.Fields{
				"message": err.Error(),
				"event":   "context",
			}).Error("InputConnection: OnEvent: Error getting session")
		}
		frontSocket, err := socketStore.Retrieve(session.FrontSocketId)
		if err != nil {
			log.WithFields(logrus.Fields{
				"socketId": session.FrontSocketId,
				"message":  err.Error(),
				"event":    "context",
			}).Error("InputConnection: OnEvent: Error getting socket")
		}

		frontSocket.Emit("input_context", inputContextJson)
		log.WithFields(logrus.Fields{
			"socketId":    frontSocket.ID(),
			"sessionCode": session.Code,
		}).Info("InputConnection: OnEvent: input context sent to front")
	})
	server.OnError(inputNsp, func(e error) {
		log.WithFields(logrus.Fields{
			"message": e.Error(),
		}).Error("InputConnection: OnError: meet error", e)
	})
	server.OnDisconnect(inputNsp, func(s socketio.Conn, reason string) {
		log.WithFields(logrus.Fields{
			"reason": reason,
		}).Info("InputConnection: OnDisconnect: disconnected")
	})

	return server
}
