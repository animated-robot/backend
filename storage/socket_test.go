package storage

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/stretchr/testify/assert"
	"testing"
)

type EmptyConn struct {
	id string
	socketio.Conn
}

func (c EmptyConn) ID() string {
	return c.id
}

func Test_SocketStore(t *testing.T) {
	id := "myID"
	conn := EmptyConn{
		id: id,
	}
	socketStore := NewSocketStoreInMemory()

	err := socketStore.Store(nil)
	assert.Error(t, err)

	err = socketStore.Store(conn)
	assert.Nil(t, err)

	_, err = socketStore.Retrieve("")
	assert.Error(t, err)

	c, _ := socketStore.Retrieve(id)
	assert.NotNil(t, c)
	assert.Equal(t, conn, c)

	err = socketStore.Remove(id)
	assert.Nil(t, err)

	err = socketStore.Remove(id)
	assert.NotNil(t, err)
}