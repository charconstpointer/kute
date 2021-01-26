package main

import (
	"log"
	"time"

	"github.com/charconstpointer/kute/pkg/kute"
)

func main() {

	end, _ := kute.NewEnding()
	middle, _ := kute.NewPipe("middle")
	start, _ := kute.NewPipe("start")

	start.Next = middle
	middle.Next = end
	end.Prev = middle
	middle.Prev = start

	kute.RunAll(start, middle, end)

	msg := make(kute.Header, kute.HeaderSize)
	msg.Encode(kute.PASS, kute.HeaderSize, 1, []byte("kute"))
	start.Write(msg)
	time.Sleep(time.Second)
	b := make([]byte, 1024)

	n, err := start.Read(b)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("%s", b[:n])
	time.Sleep(123312 * time.Second)
}
