package kute

import (
	"context"
	"net"
)

type Stream interface {
	Recv(ctx context.Context) error
	Send(msg Msg) error
}

type TCPStream struct {
	conn   net.Conn
	logger Logger
}

func NewTCPStream(conn net.Conn, logger Logger) (Stream, error) {
	stream := TCPStream{
		conn:   conn,
		logger: logger,
	}
	go stream.Recv(context.Background())
	return &stream, nil
}

func (s *TCPStream) Recv(ctx context.Context) error {
	for {
		b := make([]byte, 1024)
		n, err := s.conn.Read(b)
		if err != nil {
			return err
		}
		s.logger.Infof("read %d bytes", n)
	}
}
func (s *TCPStream) Send(msg Msg) error {
	n, err := s.conn.Write(msg.H)
	s.logger.Infof("wrote %d bytes to stream", n)
	return err
}
