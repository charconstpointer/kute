package main

import (
	"log"
	"time"

	"github.com/charconstpointer/kute/pkg/kute"
)

func main() {

	start, _ := kute.NewPipe("start")
	middle, _ := kute.NewPipe("middle")
	end, _ := kute.NewEnding()
	//[start] -> [middle] -> [end]
	start.Next = middle

	middle.Next = end
	middle.Prev = start

	end.Prev = middle

	kute.RunAll(start, middle, end)

	//create message
	msg := make(kute.Header, kute.HeaderSize)
	msg.Encode(kute.PASS, kute.HeaderSize, 1, []byte("kute"))

	//send it down the pipe
	start.Write(msg)
	time.Sleep(time.Second)
	b := make([]byte, 1024)

	//message will have come back at this point and should be uppercase, because Ending end have modified it
	n, err := start.Read(b)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("%s", b[:n])
}
