package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/morzhanov/go-event-sourcing-example/internal/mq"
	"github.com/morzhanov/go-event-sourcing-example/internal/orders"
	"github.com/segmentio/kafka-go"
)

type querystore struct {
	db *sqlx.DB
	mq mq.MQ
}

type QueryStore interface {
	Get(id string) (*orders.Order, error)
}

func (s *querystore) saveOrder(o orders.Order) {
	_, err := s.db.Query(
		`INSERT INTO orders (id, name, status)
				VALUES (?, ?, ?)`, o.Id, o.Name, o.Status,
	)
	if err != nil {
		fmt.Println(fmt.Errorf("error in query store saveOrder: %w", err))
	}
}

func (s *querystore) processOrder(id string) {
	_, err := s.db.Query("UPDATE orders SET status = 'processed' WHERE id = ?", id)
	if err != nil {
		fmt.Println(fmt.Errorf("error in processor processOrder: %w", err))
	}
}

func (s *querystore) processCommand(m *kafka.Message) {
	command := string(m.Key)
	data := m.Value
	if command == "create_order" {
		ent := orders.Order{}
		if err := json.Unmarshal(data, &ent); err != nil {
			fmt.Println(fmt.Errorf("error in query store processCommand: %w", err))
		}
		s.saveOrder(ent)
		return
	}
	s.processOrder(string(data))
}

func (s *querystore) listenCommands() {
	r := s.mq.CreateReader("qs")
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Println(fmt.Errorf("error in query store commands listener: %w", err))
			continue
		}
		go s.processCommand(&m)
	}

}

func (s *querystore) Get(id string) (*orders.Order, error) {
	e := orders.Order{}
	if err := s.db.QueryRow("SELECT * FROM orders WHERE id = ?", id).Scan(&e); err != nil {
		return nil, err
	}
	return &e, nil
}

func NewQueryStore(db *sqlx.DB, mq mq.MQ) QueryStore {
	qs := &querystore{db, mq}
	go qs.listenCommands()
	return qs
}
