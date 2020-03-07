package socket

import (
	"animated-robot/domain"
	"animated-robot/storage"
	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)


func NewDummyLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Out = ioutil.Discard
	return logger
}

func TestFrontNamespace_GetNamespace(t *testing.T) {
	front := NewFrontNamespace(nil,nil,nil)
	assert.Equal(t, "/front", front.GetNamespace())
}

func NewMockSocketStore() *MockSocketStore {
	return &MockSocketStore{}
}

type MockSocketStore struct {
	storage.ISocketStore
	SocketId string
}

func (ms *MockSocketStore) Store(c socketio.Conn) error {
	ms.SocketId = c.ID()
	return nil
}

func NewMockSessionStore(sessionCode string) *MockSessionStore {
	return &MockSessionStore{
		SessionCode: sessionCode,
	}
}

type MockSessionStore struct {
	SessionCode  string
	CreateCalled bool
	SocketId     string
	storage.ISessionStore
}

func (ms *MockSessionStore) Create(socketId string) (string, error) {
	ms.SocketId = socketId
	ms.CreateCalled = true
	return ms.SessionCode, nil
}


func (ms *MockSessionStore) Get(sessionCode string) (domain.Session, error) {
	return domain.Session{
		FrontSocketId: ms.SocketId,
		SessionCode:   sessionCode,
		Players:       []domain.Player{},
	}, nil
}


type MockConn struct {
	socketio.Conn
	Id           string
	EmitCalled   bool
	EmitEventArg string
	EmitValueArg []interface{}
}

func NewMockConn(id string) *MockConn {
	return &MockConn{
		Id:   id,
	}
}

func (mc *MockConn) ID() string {
	return mc.Id
}

func (mc *MockConn) Emit(msg string, v ...interface{}) {
	mc.EmitCalled = true
	mc.EmitEventArg = msg
	mc.EmitValueArg = v
}

func TestFrontNamespace_EnterSessionEvent(t *testing.T) {
	sessionCode := "1234"
	socketStore := NewMockSocketStore()
	sessionStore := NewMockSessionStore(sessionCode)
	front := NewFrontNamespace(socketStore, sessionStore, NewDummyLogger())

	sockerId := "socketID"
	conn := NewMockConn(sockerId)

	front.createSession(conn, sessionCode)

	assert.Equal(t, sockerId, socketStore.SocketId)
	assert.True(t, sessionStore.CreateCalled)
	assert.Equal(t, sockerId,sessionStore.SocketId)
	assert.Equal(t, sessionCode, sessionStore.SessionCode)
	assert.True(t, conn.EmitCalled)
	assert.Equal(t, "session_created", conn.EmitEventArg)
	assert.Equal(t, `{"socketId":"` + sockerId + `","sessionCode":"` + sessionCode +`","players":[]}`, conn.EmitValueArg[0].(string))
}