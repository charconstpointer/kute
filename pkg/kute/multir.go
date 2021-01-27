package kute

import (
	"context"
	"io"
	"log"

	"golang.org/x/sync/errgroup"
)

type Multir struct {
	readers []io.Reader

	R chan []byte
}

func (m *Multir) BeginRead() error {
	g, _ := errgroup.WithContext(context.Background())
	for _, r := range m.readers {
		re := r
		g.Go(func() error {
			return m.recv(re)
		})
	}

	return nil
}

func (m *Multir) Read(b []byte) (n int, err error) {
	for {
		log.Println("readdddd")
		select {
		case r := <-m.R:
			copy(b, r)
			return len(b), nil
		}

	}
}

func (m *Multir) recv(r io.Reader) error {
	log.Println("rcv")
	for {
		b := make([]byte, 1024)
		n, err := r.Read(b)
		if err != nil {
			return err
		}
		if n > 0 {
			m.R <- b[:n]
		}
	}
}

func NewMultir(readers ...io.Reader) (*Multir, error) {
	m := Multir{
		readers: readers,
		R:       make(chan []byte),
	}
	return &m, nil
}
