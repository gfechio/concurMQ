package main

import (
	hrotti "github.com/alsm/hrotti/broker"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func broker() {
	h := hrotti.NewHrotti(100, &hrotti.MemoryPersistence{})
	hrotti.INFO = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	h.AddListener("concurMQ test", hrotti.NewListenerConfig("tcp://0.0.0.0:1883"))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	h.Stop()
}
