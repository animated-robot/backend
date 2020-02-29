package socket

import (
	"animated-robot/domain"
	"animated-robot/storage"
	"encoding/json"
	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
)

type FrontNamespace struct {
	socketStore storage.ISocketStore
	sessionStore storage.ISessionStore
	log *logrus.Logger
}

func NewFrontNamespace(socketStore storage.ISocketStore, sessionStore storage.ISessionStore, log *logrus.Logger) *FrontNamespace {
	return &FrontNamespace{
		socketStore: socketStore,
		sessionStore: sessionStore,
		log: log,
	}
}

func (n FrontNamespace) GetNamespace() string {
	return "/front"
}

func (n FrontNamespace) OnConnect(s socketio.Conn) error {
	s.SetContext("")

	n.log.WithFields(logrus.Fields{
		"socketId": s.ID(),
	}).Info("FrontConnection: OnConnect: connection stablished")

	return nil
}

func (n FrontNamespace) OnDisconnect(s socketio.Conn, reason string) {
	n.log.WithFields(logrus.Fields{
		"socketId": s.ID(),
		"reason": reason,
	}).Info("FrontConnection: OnDisconnect: disconnected")
}

func (n FrontNamespace) OnError(c socketio.Conn, e error) {
	n.log.WithFields(logrus.Fields{
		"socketId": c.ID(),
		"error": e.Error(),
	}).Error("FrontConnection: OnError: meet error")
}

func (n FrontNamespace) OnEvents() map[string]OnEventHandler {
	events := make(map[string]OnEventHandler)
	events["enter_session"] = n.enterSession
	events["create_session"] = n.createSession
	return events
}

func (n FrontNamespace) createSession(s socketio.Conn, str string) {
	n.log.Trace(str)
	id, err := n.socketStore.Store(s)
	if err != nil {
		n.log.WithFields(logrus.Fields{
			"socketId": id,
			"error":    err.Error(),
		}).Error("FrontConnection: OnConnect: Socket storing failed")
	}
	n.log.WithFields(logrus.Fields{
		"socketId": id,
	}).Trace("FrontConnection: OnConnect: socket stored")

	sessionCode, _ := n.sessionStore.Create(s.ID())
	n.log.WithFields(logrus.Fields{
		"sessionCode":   sessionCode,
		"frontSocketId": s.ID(),
	}).Info("FrontConnection: OnConnect: Session Created")

	_, sessionJson, err := n.getAndJSONParseSession(sessionCode)
	if err != nil {
		n.log.WithFields(logrus.Fields{
			"event": "create_session",
			"error": err.Error(),
			"sessionCode": sessionCode,
		}).Error("InputConnection: OnEvent: Could not get and parse session")
	}

	s.Emit("session_created", sessionJson)
	n.log.WithFields(logrus.Fields{
		"event": "create_session",
		"socketId":    s.ID(),
		"sessionCode": sessionCode,
		"session": sessionJson,
	}).Info("FrontConnection: OnEvent: session sent to front")
}

func (n FrontNamespace) enterSession(s socketio.Conn, sessionCode string) {
	socketId, err := n.socketStore.Store(s)
	if err != nil {
		n.log.WithFields(logrus.Fields{
			"socketId": s.ID(),
			"sessionCode": sessionCode,
			"event": "enter_session",
			"error": err.Error(),
		}).Error("FrontConnection: OnEvent: ")
	}

	err = n.sessionStore.UpdateSocketId(sessionCode, socketId)
	if err != nil {
		n.log.WithFields(logrus.Fields{
			"socketId":    socketId,
			"sessionCode": sessionCode,
			"event":       "enter_session",
			"error":     err.Error(),
		}).Error("FrontConnection: OnEvent: ")
	}

	_, sessionJson, err := n.getAndJSONParseSession(sessionCode)
	if err != nil {
		n.log.WithFields(logrus.Fields{
			"socketId":    socketId,
			"event":       "enter_session",
			"error":     err.Error(),
			"sessionCode": sessionCode,
		}).Error("FrontConnection: OnEvent: Could not get and parse session")
	}

	s.Emit("session_entered", sessionJson)
	n.log.WithFields(logrus.Fields{
		"event":       "enter_session",
		"socketId":    socketId,
		"sessionCode": sessionCode,
		"session": sessionJson,
	}).Info("FrontConnection: OnEvent: session sent to front")
}

func (n FrontNamespace) getAndJSONParseSession(sessionCode string) (domain.Session, string, error) {
	session, err := n.sessionStore.Get(sessionCode)
	if err != nil {
		return domain.Session{}, "", err
	}

	sessionJson, err := json.Marshal(session)
	if err != nil {
		return session, "", err
	}

	return session, string(sessionJson), nil
}