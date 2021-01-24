package kute

import (
	"net"
)

type Stream interface {
	Replies() chan Msg
	Send(msg Msg) error
}

type TCPStream struct {
	conn   net.Conn
	R      chan Msg
	logger Logger
}

func NewTCPStream(conn net.Conn, logger Logger) (Stream, error) {
	stream := TCPStream{
		conn:   conn,
		logger: logger,
		R:      make(chan Msg),
	}
	go stream.recv()
	return &stream, nil
}

func (s *TCPStream) recv() error {
	for {
		b := make([]byte, 1024)
		n, err := s.conn.Read(b)
		if err != nil {
			return err
		}
		s.logger.Infof("read %d bytes", n)
		s.R <- Msg{H: make(Header, HeaderSize)}
		s.logger.Infof("sent reply")
	}
}

func (s *TCPStream) Replies() chan Msg {
	return s.R
}
func (s *TCPStream) Send(msg Msg) error {
	n, err := s.conn.Write(msg.H)
	s.logger.Infof("wrote %d bytes to stream", n)
	return err
}
