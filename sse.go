package sse

import (
	"fmt"
	"strings"
)

type Sse struct {
	buffer []string
	data   chan string
}

func NewSse() *Sse {
	buf := make([]string, 0)
	dataChan := make(chan string)
	sse := Sse{buffer: buf, data: dataChan}
	sse.SetRetry(2000)
	return &sse
}

func (this *Sse) AddMessage(event string, text string) {
	event = fmt.Sprintf("event: %s\n", event)
	this.buffer = append(this.buffer, event)
	for _, line := range strings.Split(text, "\n") {
		data := fmt.Sprintf("data: %s\n", line)
		this.buffer = append(this.buffer, data)
	}
	this.buffer = append(this.buffer, "\n")
}

func (this *Sse) SetEventId(eventId string) {
	if eventId != "" {
		eventId = fmt.Sprintf("id: %s\n\n", eventId)
	} else {
		eventId = fmt.Sprintf("id\n\n")
	}

	this.buffer = append(this.buffer, eventId)
}

func (this *Sse) SetRetry(retry int) {
	retryStr := fmt.Sprintf("retry: %d\n\n", retry)
	this.buffer = append(this.buffer, retryStr)
}

func (this *Sse) String() string {
	bufferStr := strings.Join(this.buffer, "")
	return fmt.Sprintf("<%s>", bufferStr)
}

func (this *Sse) Flush() {
	this.buffer = make([]string, 0)
}

func (this *Sse) Data() []string {
	defer func() { this.Flush() }()
	return this.buffer
}
