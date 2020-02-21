package main

import (
	"animated-robot/domain"
	"animated-robot/storage"
	"animated-robot/tools"
	"encoding/json"
	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
)

type ISocketFactory interface {
	New() *socketio.Server
}

type SocketFactory struct {
	uuidGenerator tools.UUIDGenerator
	socketStore storage.ISocketStore
	sessionStore storage.ISessionStore
	log *logrus.Logger
}

func NewSocketFactory(socketStore storage.ISocketStore, sessionStore storage.ISessionStore, uuidGenerator tools.UUIDGenerator, log *logrus.Logger) ISocketFactory {
	return SocketFactory{
		uuidGenerator: uuidGenerator,
		socketStore:  socketStore,
		sessionStore: sessionStore,
		log:          log,
	}
}

func (sf SocketFactory) New() *socketio.Server {
	server, err := socketio.NewServer(nil)
	if err != nil {
		sf.log.WithFields(logrus.Fields{
			"message": err.Error(),
		}).Fatal("NewSocket: Error creating socket")
	}

	sf.setupFrontNamespace(server)
	sf.setupInputNamespace(server)

	return server
}

func (sf SocketFactory) setupInputNamespace(server *socketio.Server) {
	inputNsp := "/input"
	server.OnConnect(inputNsp, func(s socketio.Conn) error {
		sf.log.WithFields(logrus.Fields{
			"socketId": s.ID(),
		}).Info("InputConnection: OnConnect: connection stablished")

		return nil
	})
	server.OnEvent(inputNsp, "register_player", func(s socketio.Conn, inputPlayerJson string) {
		var registerPlayer domain.RegisterPlayer
		err := json.Unmarshal([]byte(inputPlayerJson), &registerPlayer)
		if err != nil {
			sf.log.WithFields(logrus.Fields{
				"input": inputPlayerJson,
				"event": "register_player",
			}).Error("InputConnection: OnEvent: Could not parse from json string")
		}

		playerId, _ := sf.uuidGenerator.Generate()
		err = sf.sessionStore.AddPlayer(registerPlayer.SessionCode, playerId)
		if err != nil {
			player, _ := json.Marshal(registerPlayer.Player)
			sf.log.WithFields(logrus.Fields{
				"player": player,
				"event": "register_player",
				"message": err.Error(),
			}).Error("InputConnection: OnEvent: Could not store player")
		}
		// input
		s.Emit("player_registered", playerId.String())

		// front
		session, sessionJson, err := sf.getAndJSONParseSession(registerPlayer.SessionCode)
		if err != nil {
			sf.log.WithFields(logrus.Fields{
				"event": "create_session",
				"message": err.Error(),
				"sessionCode": registerPlayer.SessionCode,
			}).Error("InputConnection: OnEvent: Could not get and parse session")
		}

		socket, err := sf.socketStore.Retrieve(session.FrontSocketId)

		socket.Emit("session_changed", sessionJson)
	})
	server.OnEvent(inputNsp, "context", func(s socketio.Conn, inputContextJson string) {
		var inputContext domain.InputContext
		err := json.Unmarshal([]byte(inputContextJson), &inputContext)
		if err != nil {
			sf.log.WithFields(logrus.Fields{
				"input": inputContextJson,
				"event": "context",
			}).Error("InputConnection: OnEvent: Could not parse from json string")
		}

		session, err := sf.sessionStore.Get(inputContext.SessionCode)
		if err != nil {
			sf.log.WithFields(logrus.Fields{
				"message": err.Error(),
				"event":   "context",
			}).Error("InputConnection: OnEvent: Error getting session")
		}
		frontSocket, err := sf.socketStore.Retrieve(session.FrontSocketId)
		if err != nil {
			sf.log.WithFields(logrus.Fields{
				"socketId": session.FrontSocketId,
				"message":  err.Error(),
				"event":    "context",
			}).Error("InputConnection: OnEvent: Error getting skt")
		}

		frontSocket.Emit("input_context", inputContextJson)
		sf.log.WithFields(logrus.Fields{
			"socketId":    frontSocket.ID(),
			"sessionCode": session.Code,
		}).Info("InputConnection: OnEvent: input context sent to front")
	})
	server.OnError(inputNsp, func(s socketio.Conn, e error) {
		sf.log.WithFields(logrus.Fields{
			"socketId": s.ID(),
			"message": e.Error(),
		}).Error("InputConnection: OnError: meet error", e)
	})
	server.OnDisconnect(inputNsp, func(s socketio.Conn, reason string) {
		sf.log.WithFields(logrus.Fields{
			"socketId": s.ID(),
			"reason": reason,
		}).Info("InputConnection: OnDisconnect: disconnected")
	})
}

func (sf SocketFactory) setupFrontNamespace(server *socketio.Server) {
	frontNsp := "/front"
	server.OnConnect(frontNsp, func(s socketio.Conn) error {
		s.SetContext("")

		sf.log.WithFields(logrus.Fields{
			"socketId": s.ID(),
		}).Info("FrontConnection: OnConnect: connection stablished")

		return nil
	})
	server.OnEvent(frontNsp, "create_session", func(s socketio.Conn, str string) {
		sf.log.Trace(str)
		id, err := sf.socketStore.Store(s)
		if err != nil {
			sf.log.WithFields(logrus.Fields{
				"socketId": id,
				"error":    err.Error(),
			}).Error("FrontConnection: OnConnect: Socket storing failed")
		}
		sf.log.WithFields(logrus.Fields{
			"socketId": id,
		}).Trace("FrontConnection: OnConnect: skt stored")

		sessionCode, _ := sf.sessionStore.Create(s.ID())
		sf.log.WithFields(logrus.Fields{
			"sessionCode":   sessionCode,
			"frontSocketId": s.ID(),
		}).Info("FrontConnection: OnConnect: Session Created")

		_, sessionJson, err := sf.getAndJSONParseSession(sessionCode)
		if err != nil {
			sf.log.WithFields(logrus.Fields{
				"event": "create_session",
				"message": err.Error(),
				"sessionCode": sessionCode,
			}).Error("InputConnection: OnEvent: Could not get and parse session")
		}

		s.Emit("session_created", sessionJson)
	})
	server.OnError(frontNsp, func(s socketio.Conn, e error)  {
		sf.log.WithFields(logrus.Fields{
			"socketId": s.ID(),
			"message": e.Error(),
		}).Error("FrontConnection: OnError: meet error")
	})
	server.OnDisconnect(frontNsp, func(s socketio.Conn, reason string) {
		sf.log.WithFields(logrus.Fields{
			"socketId": s.ID(),
			"reason": reason,
		}).Info("FrontConnection: OnDisconnect: disconnected")
	})
}

func (sf SocketFactory) getAndJSONParseSession(sessionCode string) (domain.Session, string, error) {
	session, err := sf.sessionStore.Get(sessionCode)
	if err != nil {
		return domain.Session{}, "", err
	}

	sessionJson, _ := json.Marshal(session)
	if err != nil {
		return session, "", err
	}

	return session, string(sessionJson), nil
}