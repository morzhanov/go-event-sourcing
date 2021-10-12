package main

import (
	"fmt"
	"log"

	"github.com/morzhanov/go-event-sourcing-example/internal/config"
	"github.com/morzhanov/go-event-sourcing-example/internal/mq"
	"github.com/morzhanov/go-event-sourcing-example/internal/processor"
)

const initErr = "initialization error in step %s: %w"

func handleErr(step string, err error) {
	if err != nil {
		log.Fatal(fmt.Errorf(initErr, step, err))
	}
}

func main() {
	c, err := config.NewConfig("./config", ".env")
	handleErr("config", err)
	msgQ, err := mq.NewMq(c.KafkaUri, c.KafkaTopic)
	handleErr("mq", err)
	p := processor.NewProcessor(msgQ)
	p.Run()
}
