package internal

import (
	"context"
	"encoding/json"
	"github.com/morzhanov/go-event-sourcing-example/internal/mq"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type commandstore struct {
	coll *mongo.Collection
	mq   mq.MQ
}

type CommandStore interface {
	Add(command string, data interface{}) error
}

func (s *commandstore) sendEvent(command string, data interface{}) error {
	w := kafka.Writer{Balancer: &kafka.LeastBytes{}}
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	m := kafka.Message{
		Key:   []byte(command),
		Value: js,
	}
	if err := w.WriteMessages(context.Background(), m); err != nil {
		return err
	}
	return w.Close()
}

func (s *commandstore) Add(command string, data interface{}) error {
	_, err := s.coll.InsertOne(context.Background(), data)
	if err != nil {
		return err
	}
	return s.sendEvent(command, data)
}

func NewCommandStore(coll *mongo.Collection, mq mq.MQ) CommandStore {
	return &commandstore{coll, mq}
}
