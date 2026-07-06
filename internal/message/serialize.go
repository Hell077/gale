package message

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

const HeaderSize = 8 + 8 + 4

var ErrInvalidMessage = errors.New("invalid message")

func (m *Message) MarshalBinary() ([]byte, error) {
	buf := make([]byte, HeaderSize+len(m.Data))
	if len(m.Data) > math.MaxUint32 {
		return nil, ErrInvalidMessage
	}
	binary.BigEndian.PutUint64(buf[0:8], m.Topic)
	binary.BigEndian.PutUint64(buf[8:16], m.Offset)
	binary.BigEndian.PutUint32(buf[16:20], uint32(len(m.Data)))

	copy(buf[20:], m.Data)

	return buf, nil
}

func UnmarshalBinary(buf []byte) (*Message, error) {
	if len(buf) < HeaderSize {
		return nil, ErrInvalidMessage
	}

	topic := binary.BigEndian.Uint64(buf[0:8])
	offset := binary.BigEndian.Uint64(buf[8:16])
	size := binary.BigEndian.Uint32(buf[16:20])

	if len(buf[20:]) < int(size) {
		return nil, ErrInvalidMessage
	}

	msg := &Message{
		Topic:  topic,
		Offset: offset,
		Data:   make([]byte, size),
	}

	copy(msg.Data, buf[20:20+size])

	return msg, nil
}

func ReadFrom(r io.Reader) (*Message, error) {
	header := make([]byte, HeaderSize)
	if _, err := io.ReadFull(r, header); err != nil {
		return nil, err
	}

	size := binary.BigEndian.Uint32(header[16:20])
	if size > math.MaxInt32 {
		return nil, ErrInvalidMessage
	}

	buf := make([]byte, HeaderSize+int(size))
	copy(buf, header)
	if _, err := io.ReadFull(r, buf[HeaderSize:]); err != nil {
		return nil, err
	}

	return UnmarshalBinary(buf)
}

func WriteTo(w io.Writer, msg *Message) error {
	binary, err := msg.MarshalBinary()
	if err != nil {
		return err
	}
	n, err := w.Write(binary)
	if err != nil {
		return err
	}
	if n != len(binary) {
		return io.ErrShortWrite
	}
	return nil
}
