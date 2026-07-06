package message

import (
	"time"
)

const (
	Success = iota // 0
	Error          // 1
)

type Message struct {
	Topic  uint64
	Data   []byte
	Offset uint64
}

type SendLog struct {
	Time    time.Time
	Message int
}


func NewMessage(topic uint64, data []byte) *Message {
	return &Message{
		Topic: topic,
		Data:  data,
	}
}
