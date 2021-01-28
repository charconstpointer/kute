package kute

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Mux struct {
	clients map[int]net.Conn
	next    net.Conn
}

func NewMux() *Mux {
	return &Mux{
		clients: make(map[int]net.Conn),
	}
}

func (m *Mux) Dial(addr string) error {
	if m.next != nil {
		return fmt.Errorf("cannot dial new mux as this mux is already connected with a different one on addr %s", m.next.RemoteAddr().String())
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return errors.Wrap(err, "cannot dial another mutex")
	}

	m.next = conn
	return err
}

func (m *Mux) ListenMux(addr string) error {
	l, err := net.Listen("tcp", addr)
	conn, err := l.Accept()
	m.next = conn
	go m.Recv()
	log.Printf("new mux connected %s", conn.RemoteAddr().String())
	return err
}

func (m *Mux) Listen(addr string) error {
	g, _ := errgroup.WithContext(context.Background())
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		id := getID(conn.RemoteAddr().String())
		m.clients[int(id)] = conn
		g.Go(func() error {
			return m.handleConn(conn, int(id))
		})
	}

}

func (m *Mux) Recv() error {
	for {
		b := make([]byte, 32*1024)
		n, err := m.next.Read(b)
		if err != nil {
			return err
		}

		if n > 0 {
			h := Header(b[:n])
			c, e := m.clients[int(h.ID())]
			if !e {
				conn, err := net.Dial("tcp", ":25565")
				if err != nil {
					log.Println("cannot dial minecraft", err.Error())
					return err
				}
				m.clients[int(h.ID())] = conn
				go m.handleConn(conn, int(h.ID()))
			}
			c = m.clients[int(h.ID())]
			n, err := c.Write(h.Payload()[:int(h.Len())])
			if err != nil {
				log.Println("cannot write payload to minecraft")
			}

			if n == 0 {
				log.Println("cannot write payload to minecraft")
			}
			log.Printf("wrote %d to minecraft server", n)
		}
	}
}

func (m *Mux) handleConn(conn net.Conn, id int) error {
	//id := getID(conn.RemoteAddr().String())
	log.Printf("handling new conn %s", conn.RemoteAddr().String())
	for {
		b := make([]byte, 32*1024)
		nr, err := conn.Read(b)

		if err != nil {
			return errors.Wrap(err, "cannot read from handled connection")
		}

		if nr > 0 {
			h := make(Header, HeaderSize)
			h.Encode(PASS, int32(id), b[:nr])
			nw, err := m.next.Write(h)
			if nr != int(h.Len()) {
				log.Println("nr != nw", nr, nw, h.Len())
			}
			if err != nil {
				log.Println(err.Error())
				return errors.Wrap(err, "cannot write message from handled connection")
			}
			if nw != len(h) {
				log.Println("len not right :L")
			}
			log.Printf("wrote %d bytes from %d", nw, id)
		}
	}
}

func getID(addr string) uint32 {
	i := strings.LastIndex(addr, ":")
	id := addr[i+1:]
	parsed, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err.Error())
	}

	return uint32(parsed)
}
