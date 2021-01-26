package kute

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

type Runner interface {
	Run() error
}

//Pipe is an implementation of io.ReaderWriter, it allows you to send messages wrapped in Header struct
//You can combine as many pipes as you like
//Since net.Conn implements io.ReaderWriter as well, you are not limited to a single machine
//But should probably use io.ReaderWriterCloser for that
type Pipe struct {
	//Name of this pipe, helps debugging
	name string

	//Other pipes connected to this pipe
	//This is how imagine this in my head
	//prev -> [this] <- next -> [this] <- next
	Next io.ReadWriter
	Prev io.ReadWriter

	//SendCh handles async writes to a pipe, if you think about it, this is required step
	//Because Pipe invokes another Pipe's Write, this chain would never end
	sendCh chan Header

	//Buffer for storing messages
	buf []byte
}

func NewPipe(name string) (*Pipe, error) {
	return &Pipe{
		name:   name,
		sendCh: make(chan Header),
		buf:    make([]byte, 1024),
	}, nil
}

func (p *Pipe) Run() error {
	go func() {
		for {

			select {
			case msg := <-p.sendCh:
				switch msg.MessageType() {
				case PASS:
					_, err := p.Next.Write(msg)
					if err != nil {
						log.Fatal(err.Error())
					}

					break
				case REPL:
					if p.Prev == nil {
						p.buf = msg.Payload()
						continue
					}
					_, err := p.Prev.Write(msg)
					if err != nil {
						log.Fatal(err.Error())
					}
					break
				}
			}
		}
	}()
	return nil
}

func (p *Pipe) Read(b []byte) (n int, err error) {
	n = copy(b, p.buf)
	return n, nil
}

func (p *Pipe) Write(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	if p.Next == nil {
		return 0, errors.New("pipeline ended unexpectedly")
	}

	p.sendCh <- b

	return len(b), nil
}

//EchoEnding is different type of a pipe
//It modifies the content it receives
//In this example it uppercases the payload
type EchoEnding struct {
	Prev   io.ReadWriter
	b      []byte
	sendCh chan Header
}

func NewEchoEnding() (*EchoEnding, error) {
	return &EchoEnding{
		b:      make([]byte, 1024),
		sendCh: make(chan Header),
	}, nil
}

func (e *EchoEnding) Run() error {
	go func() {
		for {
			select {
			case msg := <-e.sendCh:
				if len(msg) > 0 {

					_, err := e.Prev.Write(msg)
					if err != nil {
						log.Fatal(err.Error())
					}
				}
			}
		}

	}()
	return nil
}

func (e *EchoEnding) Read(b []byte) (n int, err error) {
	b = e.b
	return len(b), nil
}

func (e *EchoEnding) Write(b []byte) (n int, err error) {
	h := Header(b)
	e.b = h
	h.Encode(REPL, h.Len(), h.ID(), []byte(strings.ToUpper(fmt.Sprintf("%s", h.Payload()))))
	e.sendCh <- h
	return len(b), nil
}

func RunAll(runners ...Runner) {
	for _, runner := range runners {
		runner.Run()
	}
}
