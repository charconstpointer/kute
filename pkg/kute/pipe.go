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
type Pipe struct {
	name string

	Next io.ReadWriter
	Prev io.ReadWriter

	sendCh chan Header
	buf    []byte
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

type Ending struct {
	Prev   io.ReadWriter
	b      []byte
	sendCh chan Header
}

func NewEnding() (*Ending, error) {
	return &Ending{
		b:      make([]byte, 1024),
		sendCh: make(chan Header),
	}, nil
}
func (e *Ending) Run() {
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
}
func (e *Ending) Read(b []byte) (n int, err error) {
	b = e.b
	return len(b), nil
}
func (e *Ending) Write(b []byte) (n int, err error) {
	h := Header(b)
	e.b = h
	h.Encode(REPL, h.Len(), h.ID(), []byte(strings.ToUpper(fmt.Sprintf("%s", h.Payload()))))
	e.sendCh <- h
	return len(b), nil
}
