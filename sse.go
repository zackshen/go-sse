package sse

import (
	"fmt"
	"net/http"
	"strings"
)

type SSE struct {
	buffer []string
}

func NewSSE() *SSE {
	buf := make([]string, 0)
	SSE := SSE{buffer: buf}
	SSE.SetRetry(2000)
	return &SSE
}

func (s *SSE) AddMessage(event string, text string) {
	event = fmt.Sprintf("event: %s\n", event)
	s.buffer = append(s.buffer, event)
	for _, line := range strings.Split(text, "\n") {
		data := fmt.Sprintf("data: %s\n", line)
		s.buffer = append(s.buffer, data)
	}
	s.buffer = append(s.buffer, "\n")
}

func (s *SSE) SetEventId(eventId string) {
	if eventId != "" {
		eventId = fmt.Sprintf("id: %s\n\n", eventId)
	} else {
		eventId = fmt.Sprintf("id\n\n")
	}

	s.buffer = append(s.buffer, eventId)
}

func (s *SSE) SetRetry(retry int) {
	retryStr := fmt.Sprintf("retry: %d\n\n", retry)
	s.buffer = append(s.buffer, retryStr)
}

func (s *SSE) String() string {
	bufferStr := strings.Join(s.buffer, "")
	bufferStr = fmt.Sprintf("<%s>", bufferStr)
	s.Flush()
	return bufferStr
}

func (s *SSE) Flush() {
	s.buffer = make([]string, 0)
}

func (s *SSE) Data() []string {
	defer func() { s.Flush() }()
	return s.buffer
}

type Connection chan string
type Messages chan string

type SSEHandler struct {
	sse          *SSE
	connections  map[Connection]bool
	connChan     chan Connection
	disconnChan  chan Connection
	messagesChan Messages
}

func NewSSEHandler() *SSEHandler {
	handler := &SSEHandler{
		sse:          &SSE{buffer: make([]string, 0)},
		connections:  make(map[Connection]bool, 100),
		connChan:     make(chan Connection, 100),
		disconnChan:  make(chan Connection, 100),
		messagesChan: make(Messages, 100),
	}
	return handler
}

func (s *SSEHandler) HttpHandler(response http.ResponseWriter, request *http.Request) {
	flusher, ok := response.(http.Flusher)

	if !ok {
		http.Error(response, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	response.Header().Set("Content-Type", "text/event-stream")
	response.Header().Set("Cache-Control", "no-cache")
	response.Header().Set("Connection", "keep-alive")
	response.Header().Set("Access-Control-Allow-Origin", "*")

	s.sse.SetRetry(10000)
	retry := s.sse.String()
	response.Write([]byte(retry))

	connection := make(Connection)
	s.connChan <- connection

	closeNotify := response.(http.CloseNotifier).CloseNotify()

	go func() {
		<-closeNotify
		s.disconnChan <- connection
	}()

	for {
		select {
		case msg := <-connection:
			response.Write([]byte(msg))
			flusher.Flush()
		}
	}
}

func (s *SSEHandler) Broadcast(msg string) {
	for conn := range s.connections {
		s.sse.AddMessage("test", msg)
		conn <- s.sse.String()
	}
}

func (s *SSEHandler) RemoveConnection(conn Connection) {
	delete(s.connections, conn)
}

func (s *SSEHandler) AddConnection(conn Connection) {
	s.connections[conn] = true
}

func (s *SSEHandler) Listen() {
	go func() {
		for {
			select {
			case msg := <-s.messagesChan:
				s.Broadcast(msg)
			case conn := <-s.connChan:
				s.AddConnection(conn)
			case conn := <-s.disconnChan:
				s.RemoveConnection(conn)
			}
		}
	}()
}
