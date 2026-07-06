package session

import (
	"time"

	"github.com/Hell077/gale/internal/message"
)

func (s *Conn) SendMessage(msg *message.Message) (message.SendLog, error) {
	if err := message.WriteTo(s.rw, msg); err != nil {
		return message.SendLog{
			Time:    time.Now(),
			Message: message.Error,
		}, err
	}
	s.queue.Push(msg.Topic)
	return message.SendLog{
		Time:    time.Now(),
		Message: message.Success,
	}, nil
}

func (s *Conn) ReceiveMessage() (*message.Message, error) {
	msg, err := message.ReadFrom(s.rw)
	if err != nil {
		return &message.Message{}, err
	}
	return msg, nil
}

func (s *Conn) RouteMessage() {

}
