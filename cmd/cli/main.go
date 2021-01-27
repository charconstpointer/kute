package main

import (
	"bufio"
	"log"
	"net"
	"os"

	"github.com/charconstpointer/kute/pkg/kute"
)

type MyConn struct {
	conn net.Conn
}

func (c *MyConn) Write(b []byte) (n int, err error) {
	bh := make([]byte, 1024)
	h := kute.Header(bh)
	h.Encode(kute.PASS, uint32(len(b)), 1, b)
	return c.conn.Write(h)
}

func (c *MyConn) Read(b []byte) (n int, err error) {
	return c.conn.Read(b)
}

func main() {
	conn, err := net.Dial("tcp", ":9999")
	c := MyConn{conn}
	if err != nil {
		log.Fatal(err.Error())
	}

	sc := bufio.NewScanner(os.Stdin)
	for {
		sc.Scan()
		text := sc.Text()

		n, err := c.Write([]byte(text))
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Println(n)
		log.Println("--------------------------------")
	}
}
