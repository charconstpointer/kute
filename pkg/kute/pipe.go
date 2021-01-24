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

type BasicPipe struct {
	name string

	in  Ending
	out Ending

	sendCh chan Msg

	addr string
	next string

	state State

	logger Logger
}

type State int

const (
	NotReady State = iota
	Ready
)

type Msg struct {
	H Header
}

func NewBasicPipe(addr string, next string, name string) (Pipe, error) {
	pipe := BasicPipe{
		name:   name,
		addr:   addr,
		next:   next,
		state:  NotReady,
		sendCh: make(chan Msg),
		logger: &PipeLogger{prefix: name},
	}
	return &pipe, nil
}

func (p *BasicPipe) Send(ctx context.Context, msg Msg) error {
	p.logger.Infof("sending new msg %s", msg)
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
	p.logger.Infof("waiting for pipe connection on %s", p.addr)
	listener, err := net.Listen("tcp", p.addr)
	conn, err := listener.Accept()
	p.logger.Infof("new pipe connected on %s", p.addr)
	stream, err := NewTCPStream(conn, p.logger)
	ending := SingleEnd{
		stream: stream,
		logger: p.logger,
	}
	p.in = &ending
	return err
}

func (p *BasicPipe) connPipe() error {
	p.logger.Infof("trying to connect to pipe on %s", p.next)
	conn, err := net.Dial("tcp", p.next)
	stream, err := NewTCPStream(conn, p.logger)
	ending := SingleEnd{
		stream: stream,
		logger: p.logger,
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
