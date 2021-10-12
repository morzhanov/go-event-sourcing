package internal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/morzhanov/go-event-sourcing-example/internal/mq"
	"github.com/morzhanov/go-event-sourcing-example/internal/orders/model"
	"github.com/segmentio/kafka-go"
)

type querystore struct {
	db *sqlx.DB
	mq mq.MQ
}

type QueryStore interface {
	Get(id string) (*model.Order, error)
	GetAll() ([]*model.Order, error)
}

func (s *querystore) saveOrder(o model.Order) {
	if _, err := s.db.Exec("INSERT INTO orders (id, name, status) VALUES ($1, $2, $3)", o.Id, o.Name, o.Status); err != nil {
		fmt.Println(fmt.Errorf("error in query store saveOrder: %w", err))
	}
}

func (s *querystore) processOrder(id string) {
	if _, err := s.db.Exec("UPDATE orders SET status = 'processed' WHERE id = $1", id); err != nil {
		fmt.Println(fmt.Errorf("error in processor processOrder: %w", err))
	}
}

func (s *querystore) processCommand(m *kafka.Message) {
	command := string(m.Key)
	data := m.Value
	if command == "create_order" {
		ent := model.Order{}
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

func (s *querystore) GetAll() ([]*model.Order, error) {
	rows, err := s.db.Query("SELECT * FROM orders")
	if err != nil {
		return nil, err
	}

	var res []*model.Order
	for rows.Next() {
		var id, name, status string
		if err := rows.Scan(&id, &name, &status); err != nil {
			return nil, err
		}
		res = append(res, &model.Order{id, name, status})
	}
	return res, nil
}

func (s *querystore) Get(id string) (*model.Order, error) {
	var name, status string
	if err := s.db.QueryRow("SELECT * FROM orders WHERE id = $1", id).Scan(&id, &name, &status); err != nil {
		return nil, err
	}
	return &model.Order{id, name, status}, nil
}

func NewQueryStore(db *sqlx.DB, mq mq.MQ) QueryStore {
	qs := &querystore{db, mq}
	go qs.listenCommands()
	return qs
}
