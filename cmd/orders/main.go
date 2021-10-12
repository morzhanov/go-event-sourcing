package main

import (
	"fmt"
	"log"

	"github.com/morzhanov/go-event-sourcing-example/internal"
	"github.com/morzhanov/go-event-sourcing-example/internal/config"
	"github.com/morzhanov/go-event-sourcing-example/internal/mongodb"
	"github.com/morzhanov/go-event-sourcing-example/internal/mq"
	"github.com/morzhanov/go-event-sourcing-example/internal/orders"
	"github.com/morzhanov/go-event-sourcing-example/internal/psql"
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

	coll, err := mongodb.NewMongoDB(c.MongoDBUri)
	handleErr("mongodb", err)

	db, err := psql.NewDb(c.PsqlConnectionString)
	handleErr("psql", err)
	err = psql.RunMigrations(db)
	handleErr("psql migrations", err)

	cs := internal.NewCommandStore(coll, msgQ)
	qs := internal.NewQueryStore(db, msgQ)
	service := orders.NewService(cs, qs)
	if err = service.Listen(); err != nil {
		log.Fatal(err)
	}
}
