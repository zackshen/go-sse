package sse

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestSse(t *testing.T) {
	sse := NewSSE()
	assert.Equal(t, sse.String(), "<retry: 2000\n\n>")
}

func TestAddMessage(t *testing.T) {
	sse := NewSSE()
	sse.AddMessage("ack", "ACK")
	sse.AddMessage("message", "Hello World")
	sse.AddMessage("eof", "FIN")
	assert.Equal(t, sse.String(), "<retry: 2000\n\nevent: ack\ndata: ACK\n\nevent: message\ndata: Hello World\n\nevent: eof\ndata: FIN\n\n>")
}

func TestSetEventId(t *testing.T) {
	sse := NewSSE()
	sse.SetEventId("listen-event")
	assert.Equal(t, sse.String(), "<retry: 2000\n\nid: listen-event\n\n>")

	sse.Flush()
	sse.SetEventId("")
	assert.Equal(t, sse.String(), "<id\n\n>")
}

func TestRead(t *testing.T) {
	sse := NewSSE()
	sse.AddMessage("ack", "ACK")
	sse.AddMessage("message", "Hello World")
	sse.AddMessage("eof", "FIN")
	assert.Equal(t, len(sse.Data()), 10)

	sse.AddMessage("test", "123")
	assert.Equal(t, len(sse.Data()), 3)
}
