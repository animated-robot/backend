package socket

import (
	"animated-robot/domain"
	"animated-robot/storage"
	"animated-robot/tools"
	"encoding/json"
	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
)

type InputNamespace struct {
	uuidGenerator tools.UUIDGenerator
	socketStore storage.ISocketStore
	sessionStore storage.ISessionStore
	log *logrus.Logger
}

func (n InputNamespace) GetNamespace() string {
	return "/input"
}

func (n InputNamespace) OnConnect(s socketio.Conn) error {
	n.log.WithFields(logrus.Fields{
		"socketId": s.ID(),
	}).Info("InputConnection: OnConnect: connection stablished")

	return nil
}

func (n InputNamespace) OnDisconnect(s socketio.Conn, reason string) {
	n.log.WithFields(logrus.Fields{
		"socketId": s.ID(),
		"reason": reason,
	}).Info("InputConnection: OnDisconnect: disconnected")
}

func (n InputNamespace) OnError(s socketio.Conn, e error) {
	n.log.WithFields(logrus.Fields{
		"socketId": s.ID(),
		"error": e.Error(),
	}).Error("InputConnection: OnError: meet error")
}

func (n InputNamespace) OnEvents() map[string]OnEventHandler {
	events := make(map[string]OnEventHandler)
	events["register_player"] = n.onRegisterPlayer
	events["context"] = n.onInputContext
	return events
}

func NewInputNamespace(socketStore storage.ISocketStore, sessionStore storage.ISessionStore, log *logrus.Logger, uuidGenerator tools.UUIDGenerator) *InputNamespace {
	return &InputNamespace{
		uuidGenerator: uuidGenerator,
		socketStore: socketStore,
		sessionStore: sessionStore,
		log: log,
	}
}

func (n InputNamespace) onRegisterPlayer(s socketio.Conn, inputPlayerJson string) {
	var registerPlayer domain.RegisterPlayer
	err := json.Unmarshal([]byte(inputPlayerJson), &registerPlayer)
	if err != nil {
		n.log.WithFields(logrus.Fields{
			"input": inputPlayerJson,
			"event": "register_player",
		}).Error("InputConnection: OnEvent: Could not parse from json string")
	}

	playerId, _ := n.uuidGenerator.Generate()
	registerPlayer.Player["id"] = playerId
	err = n.sessionStore.AddPlayer(registerPlayer.SessionCode, registerPlayer.Player)
	if err != nil {
		player, _ := json.Marshal(registerPlayer.Player)
		n.log.WithFields(logrus.Fields{
			"player": player,
			"event": "register_player",
			"error": err.Error(),
		}).Error("InputConnection: OnEvent: Could not store player")
	}
	// input
	s.Emit("player_registered", playerId.String())

	// front
	session, sessionJson, err := n.getAndJSONParseSession(registerPlayer.SessionCode)
	if err != nil {
		n.log.WithFields(logrus.Fields{
			"event": "create_session",
			"error": err.Error(),
			"sessionCode": registerPlayer.SessionCode,
		}).Error("InputConnection: OnEvent: Could not get and parse session")
	}

	socket, err := n.socketStore.Retrieve(session.FrontSocketId)

	socket.Emit("session_changed", sessionJson)
}

func (n InputNamespace) onInputContext(s socketio.Conn, inputContextJson string) {
	var inputContext domain.InputContext
	err := json.Unmarshal([]byte(inputContextJson), &inputContext)
	if err != nil {
		n.log.WithFields(logrus.Fields{
			"input": inputContextJson,
			"event": "context",
			"error": err.Error(),
		}).Error("InputConnection: OnEvent: Could not parse from json string")
	}

	session, err := n.sessionStore.Get(inputContext.SessionCode)
	if err != nil {
		n.log.WithFields(logrus.Fields{
			"error": err.Error(),
			"event":   "context",
		}).Error("InputConnection: OnEvent: Error getting session")
	}
	frontSocket, err := n.socketStore.Retrieve(session.FrontSocketId)
	if err != nil {
		n.log.WithFields(logrus.Fields{
			"socketId": session.FrontSocketId,
			"error":  err.Error(),
			"event":    "context",
		}).Error("InputConnection: OnEvent: Error getting socket")
	}

	frontSocket.Emit("input_context", inputContextJson)
	n.log.WithFields(logrus.Fields{
		"socketId":    frontSocket.ID(),
		"sessionCode": session.SessionCode,
		"context": inputContextJson,
	}).Info("InputConnection: OnEvent: input context sent to front")
}

func (n InputNamespace) getAndJSONParseSession(sessionCode string) (domain.Session, string, error) {
	session, err := n.sessionStore.Get(sessionCode)
	if err != nil {
		return domain.Session{}, "", err
	}

	sessionJson, _ := json.Marshal(session)
	if err != nil {
		return session, "", err
	}

	return session, string(sessionJson), nil
}
