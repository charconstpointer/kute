package kute

import (
	"context"
)

type Ending interface {
	Send(ctx context.Context, msg Msg) error
}

type SingleEnd struct {
	stream Stream
	pipe   Pipe
	sendCh chan []byte

	logger Logger
}

func NewSingleEnd(pipe Pipe, s Stream, logger Logger) (Ending, error) {
	ending := SingleEnd{
		stream: s,
		pipe:   pipe,
		sendCh: make(chan []byte),
		logger: logger,
	}

	go ending.recv()
	return &ending, nil
}

func (e *SingleEnd) recv() error {
	for {
		select {
		case _ = <-e.sendCh:
			e.logger.Infof("end recv new msg from stream %v", e.stream)
			e.pipe.Send(context.Background(), Msg{})
		case reply := <-e.stream.Replies():
			e.logger.Infof("stream replied with %v", reply)
			e.pipe.Reply(context.Background(), reply)
		}

	}
}

func (e *SingleEnd) Send(ctx context.Context, msg Msg) error {
	e.logger.Infof("new msg %v", msg)
	return e.stream.Send(msg)
}
