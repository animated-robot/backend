package storage

import (
	socketio "github.com/googollee/go-socket.io"
	"github.com/stretchr/testify/assert"
	"testing"
)


// TODO new(socketio.Conn returns nil connection struct
func Test_SocketStore(t *testing.T) {
	conn := new(socketio.Conn)
	socketStore := NewSocketStoreInMemory()

	_, err := socketStore.Store(nil)
	assert.Error(t, err)

	id, _ := socketStore.Store(*conn)
	assert.NotEmpty(t, id)

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