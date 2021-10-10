package main

import (
	"github.com/morzhanov/go-event-sourcing-example/internal/config"
	"github.com/morzhanov/go-event-sourcing-example/internal/mq"
	"github.com/morzhanov/go-event-sourcing-example/internal/processor"
	"log"
)

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	c, err := config.NewConfig("./config", ".env")
	handleErr(err)
	msgQ, err := mq.NewMq(c.KafkaUri, c.KafkaTopic)
	handleErr(err)
	p := processor.NewProcessor(msgQ)
	p.Run()
}
