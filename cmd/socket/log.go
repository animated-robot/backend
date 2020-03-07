package socket

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
)


func NewLogEventWrapper(namespace ISocketNamespace, log *logrus.Logger) *LogEventWrapper {
	return &LogEventWrapper{
		ns: namespace,
		log: log,
	}
}

type LogEventWrapper struct {
	ns ISocketNamespace
	log *logrus.Logger
}

func (lew LogEventWrapper) GetNamespace() string {
	return lew.ns.GetNamespace()
}

func (lew LogEventWrapper) OnConnect(s socketio.Conn) error {
	lew.logEvent(s, "OnConnect", lew.GetNamespace())
	return lew.ns.OnConnect(s)
}

func (lew LogEventWrapper) OnDisconnect(s socketio.Conn, reason string) {
	lew.logEvent(s, "OnDisconnect", lew.GetNamespace())
	lew.ns.OnDisconnect(s, reason)
}

func (lew LogEventWrapper) OnError(s socketio.Conn, e error) {
	lew.logEvent(s, "OnError", lew.GetNamespace())
	lew.ns.OnError(s, e)
}

func (lew LogEventWrapper) OnEvents() map[string]OnEventHandler {
	events := make(map[string]OnEventHandler)
	namespace := lew.GetNamespace()
	for eventName, handler := range lew.ns.OnEvents() {
		events[eventName] = lew.wrap(handler, eventName, namespace)
	}
	return events
}

func (lew LogEventWrapper) wrap(handler OnEventHandler, eventName string, namespace string) OnEventHandler {
	return func(s socketio.Conn, msg string) {
		lew.logEvent(s, eventName, namespace)
		handler(s, msg)
	}
}

func (lew LogEventWrapper) logEvent(s socketio.Conn, eventName string, namespace string) {
	url := s.URL().Path + "?" + s.URL().RawQuery
	lew.log.WithFields(logrus.Fields{
		"url": url,
		"namespace": namespace,
		"event": eventName,
	}).Trace("Log Event")
}