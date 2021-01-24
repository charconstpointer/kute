package main

import (
	"context"
	"log"
	"time"

	"github.com/charconstpointer/kute/pkg/kute"
)

func main() {
	msg := kute.Msg{
		H: make(kute.Header, kute.HeaderSize),
	}

	firstPipe, err := kute.NewBasicPipe("", ":9001", "first pipe")
	secondPipe, err := kute.NewBasicPipe(":9001", "", "next pipe")

	if err := kute.RunPipes(context.Background(), firstPipe, secondPipe); err != nil {
		log.Fatal("unable to configure and run given pipe structure")
	}

	err = firstPipe.Send(context.Background(), msg)
	if err != nil {
		log.Fatal(err.Error())
	}
	time.Sleep(time.Second * 10)
}
