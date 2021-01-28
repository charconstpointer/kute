package main

import (
	"log"
	"time"

	"github.com/charconstpointer/kute/pkg/kute"
)

func main() {
	m := kute.NewMux()
	// if err := m.Dial("ec2-18-134-242-182.eu-west-2.compute.amazonaws.com:8000"); err != nil {
	if err := m.Dial(":8000"); err != nil {
		// if err := m.Dial("178.128.207.63:8000"); err != nil {
		log.Fatal(err.Error())
	}

	if err := m.Recv(); err != nil {
		log.Fatal(err.Error())
	}

	time.Sleep(time.Hour)
}
