package sse

import (
	assert "github.com/stretchr/testify/assert"
	"testing"
)

func TestSse(t *testing.T) {
	sse := NewSse()
	assert.Equal(t, sse.String(), "<retry: 2000\n\n>")
}

func TestAddMessage(t *testing.T) {
	sse := NewSse()
	sse.AddMessage("ack", "ACK")
	sse.AddMessage("message", "Hello World")
	sse.AddMessage("eof", "FIN")
	assert.Equal(t, sse.String(), "<retry: 2000\n\nevent: ack\ndata: ACK\n\nevent: message\ndata: Hello World\n\nevent: eof\ndata: FIN\n\n>")
}

func TestSetEventId(t *testing.T) {
	sse := NewSse()
	sse.SetEventId("listen-event")
	assert.Equal(t, sse.String(), "<retry: 2000\n\nid: listen-event\n\n>")

	sse.Flush()
	sse.SetEventId("")
	assert.Equal(t, sse.String(), "<id\n\n>")
}

func TestRead(t *testing.T) {
	sse := NewSse()
	sse.AddMessage("ack", "ACK")
	sse.AddMessage("message", "Hello World")
	sse.AddMessage("eof", "FIN")
	assert.Equal(t, len(sse.Data()), 10)

	sse.AddMessage("test", "123")
	assert.Equal(t, len(sse.Data()), 3)
}
