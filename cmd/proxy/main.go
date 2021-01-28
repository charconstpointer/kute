package main

import (
	"log"

	"github.com/charconstpointer/kute/pkg/kute"
)

func main() {
	m := kute.NewMux()
	if err := m.ListenMux(":8000"); err != nil {
		log.Fatal(err.Error())
	}

	if err := m.Listen(":9000"); err != nil {
		log.Fatal(err.Error())
	}
}
