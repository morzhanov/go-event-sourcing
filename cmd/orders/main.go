package main

import (
	"github.com/morzhanov/go-event-sourcing-example/internal"
	"github.com/morzhanov/go-event-sourcing-example/internal/config"
	"github.com/morzhanov/go-event-sourcing-example/internal/mongodb"
	"github.com/morzhanov/go-event-sourcing-example/internal/mq"
	"github.com/morzhanov/go-event-sourcing-example/internal/orders"
	"github.com/morzhanov/go-event-sourcing-example/internal/psql"
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
	coll, err := mongodb.NewMongoDB(c.MongoDBUri)
	handleErr(err)
	db, err := psql.NewDb(c.PsqlConnectionString)
	err = psql.RunMigrations(db, "orders")
	handleErr(err)
	handleErr(err)
	cs := internal.NewCommandStore(coll, msgQ)
	qs := internal.NewQueryStore(db, msgQ)
	service := orders.NewService(cs, qs)
	err = service.Listen()
	handleErr(err)
}
