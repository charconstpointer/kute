package main

import (
	"log"
	"time"

	"github.com/charconstpointer/kute/pkg/kute"
)

func main() {
	m := kute.NewMux()
	if err := m.Dial(":8000"); err != nil {
		log.Fatal(err.Error())
	}

	if err := m.Recv(); err != nil {
		log.Fatal(err.Error())
	}

	time.Sleep(time.Hour)
}
