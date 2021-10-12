package processor

import (
	"context"
	"fmt"

	"github.com/morzhanov/go-event-sourcing-example/internal/mq"
)

type processor struct {
	mq mq.MQ
}

type Processor interface {
	Run()
}

func (s *processor) Run() {
	r := s.mq.CreateReader("processor")
	for {
		m, err := r.ReadMessage(context.Background())
		if string(m.Key) != "process_order" {
			continue
		}
		if err != nil {
			fmt.Println(fmt.Errorf("error in processor Run: %w", err))
			continue
		}
		fmt.Printf("Processing command: %s with data id = %s", m.Key, m.Value)
	}
}

func NewProcessor(mq mq.MQ) Processor {
	return &processor{mq}
}
