package main

import (
	"io"
	"log"
	"net"

	"github.com/charconstpointer/kute/pkg/kute"
)

func main() {
	l, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal(err.Error())
	}
	conns := []io.Reader{}
	for i := 0; i < 2; i++ {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}
		conns = append(conns, conn)
	}

	mr, err := kute.NewMultir(conns...)
	go mr.BeginRead()
	for {
		b := make([]byte, 1024)
		n, err := mr.Read(b)
		if err != nil {
			log.Fatal(err.Error())
		}
		if n > 0 {
			h := kute.Header(b)
			log.Printf("[%d]{%d} => %s", h.ID(), h.Len(), h.Payload())
		}

	}
	// for {
	// 	select {
	// 	case r := <-mr.R:
	// 		h := kute.Header(r)
	// 		log.Printf("[%d]{%d} => %s", h.ID(), h.Len(), h.Payload())
	// 	}

	// }
}
