package internal

import (
	"context"
	"encoding/json"

	"github.com/morzhanov/go-event-sourcing-example/internal/mq"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type command struct {
	Name string
	Data string
}

type commandstore struct {
	coll *mongo.Collection
	mq   mq.MQ
}

type CommandStore interface {
	Add(command string, data interface{}) error
}

func (s *commandstore) sendEvent(command string, data interface{}) error {
	str, ok := data.(string)
	if !ok {
		b, err := json.Marshal(data)
		if err != nil {
			return err
		}
		str = string(b)
	}
	m := kafka.Message{
		Key:   []byte(command),
		Value: []byte(str),
	}
	if _, err := s.mq.Conn().WriteMessages(m); err != nil {
		return err
	}
	return nil
}

func (s *commandstore) Add(cmd string, data interface{}) error {
	c := command{Name: cmd}
	str, ok := data.(string)
	if !ok {
		b, err := json.Marshal(data)
		str = string(b)
		if err != nil {
			return err
		}
	}
	c.Data = str
	_, err := s.coll.InsertOne(context.Background(), c)
	if err != nil {
		return err
	}
	return s.sendEvent(cmd, data)
}

func NewCommandStore(coll *mongo.Collection, mq mq.MQ) CommandStore {
	return &commandstore{coll, mq}
}
