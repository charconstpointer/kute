package kute

import (
	"context"
	"log"
	"net"

	"golang.org/x/sync/errgroup"
)

type Pipe interface {
	Send(ctx context.Context, msg Msg) error
	Run() error
}
type State int

const (
	NotReady State = iota
	Ready
)

type BasicPipe struct {
	name string

	in  Ending
	out Ending

	sendCh chan Msg

	addr string
	next string

	state State
}

func NewBasicPipe(addr string, next string, name string) (Pipe, error) {
	pipe := BasicPipe{
		name:   name,
		addr:   addr,
		next:   next,
		state:  NotReady,
		sendCh: make(chan Msg),
	}
	return &pipe, nil
}

func (p *BasicPipe) Send(ctx context.Context, msg Msg) error {
	log.Println(p.name, "sending new msg")
	return p.out.Send(ctx, msg)
}

func (p *BasicPipe) Run() error {
	g, _ := errgroup.WithContext(context.Background())

	if p.addr != "" {
		g.Go(p.recvPipe)
	}

	if p.next != "" {
		g.Go(p.connPipe)
	}

	err := g.Wait()
	if err != nil {
		log.Fatalf("could not initalize pipe %s", p.name)
	}
	p.state = Ready
	return nil
}

func (p *BasicPipe) recvPipe() error {
	log.Printf("waiting for pipe connection on %s", p.addr)
	listener, err := net.Listen("tcp", p.addr)
	conn, err := listener.Accept()
	log.Printf("new pipe connected on %s", p.addr)
	ending := SingleEnd{
		stream: &TCPStream{
			conn: conn,
		},
	}
	p.in = &ending
	return err
}

func (p *BasicPipe) connPipe() error {
	log.Printf("trying to connect to pipe on %s", p.next)
	conn, err := net.Dial("tcp", p.next)
	ending := SingleEnd{
		stream: &TCPStream{
			conn: conn,
		},
	}
	p.out = &ending
	return err
}

func RunPipes(ctx context.Context, pipes ...Pipe) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, pipe := range pipes {
		g.Go(pipe.Run)
	}
	return g.Wait()
}

type Stream interface {
	Recv(ctx context.Context) error
	Send(msg Msg) error
}

type TCPStream struct {
	conn net.Conn
}

type Msg struct {
	H Header
}
type Ending interface {
	Send(ctx context.Context, msg Msg) error
}

type SingleEnd struct {
	stream Stream
	pipe   Pipe
	sendCh chan []byte
}

func NewSingleEnd(s Stream) (Ending, error) {
	ending := SingleEnd{
		stream: s,
		sendCh: make(chan []byte),
	}
	go ending.recv()
	return &ending, nil
}

func (s *TCPStream) Recv(ctx context.Context) error {
	for {
		b := make([]byte, 1024)
		n, err := s.conn.Read(b)
		if err != nil {
			return err
		}
		log.Printf("read %d bytes", n)
	}
}
func (s *TCPStream) Send(msg Msg) error {
	n, err := s.conn.Write(msg.H)
	log.Printf("wrote %d bytes to stream", n)
	return err
}

func (e *SingleEnd) recv() error {
	for {
		select {
		case _ = <-e.sendCh:
			log.Println("end recv new msg from stream")
			e.pipe.Send(context.Background(), Msg{})
		}
	}
}

func (e *SingleEnd) Send(ctx context.Context, msg Msg) error {
	log.Println("new msg")
	return e.stream.Send(msg)
}
