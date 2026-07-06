package message

import (
	"encoding/binary"
	"errors"
	"io"
	"time"
)

const HandshakeTimeout = 5 * time.Second

func WriteHandshake(w io.Writer) error {
	ts := time.Now().UnixMilli()
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(ts))

	n, err := w.Write(data)
	if err != nil {
		return err
	}
	if n != len(data) {
		return io.ErrShortWrite
	}
	return nil
}

func ReadHandshake(r io.Reader) error {
	buf := make([]byte, 8)

	_, err := io.ReadFull(r, buf)
	if err != nil {
		return err
	}

	ts := int64(binary.BigEndian.Uint64(buf))
	sentAt := time.UnixMilli(ts)

	if time.Since(sentAt) > HandshakeTimeout {
		return errors.New("handshake timeout")
	}

	return nil
}
