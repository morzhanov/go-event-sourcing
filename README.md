# Go Event Sourcing Example

Go CQRS + Event Sourcing example application.

## CQRS and Event Sourcing

<img src="https://i.ibb.co/XJ6FxYF/Untitled.png" alt="cqrs-es" />

CQRS stands for Command and Query Responsibility Segregation, a pattern that separates read and update operations for a data store. Implementing CQRS in your application can maximize its performance, scalability, and security. The flexibility created by migrating to CQRS allows a system to better evolve over time and prevents update commands from causing merge conflicts at the domain level.

Event sourcing persists the state of a business entity such an Order or a Customer as a sequence of state-changing events. Whenever the state of a business entity changes, a new event is appended to the list of events. Since saving an event is a single operation, it is inherently atomic. The application reconstructs an entity’s current state by replaying the events.

<a href="https://microservices.io/patterns/data/event-sourcing.html">Nice article on Event Sourcing.</a>
<a href="https://docs.microsoft.com/en-us/azure/architecture/patterns/cqrs"Aarticle on CQRS.</a>

## Project Architecture

<img src="https://i.ibb.co/5jvhLXh/IMG-0124.jpg" alt="architecture" />

- Orders - service for order creation and retrieval
    - Commands used for order creation and processing
    - Queries used for orders retrieval
- Processor - dumb service which mocks order processing
- Command store - persists commands (MongoDB)
- Query store - consumes events from kafka commands channel and saves/updates data for queries (PostgreSQL)
- Kafka - used as communication level between Command-Query stores and Order-Processor services


## Running

1. `cd /deploy`
2. `docker-compose up -d`
3. `cd ../`
4. `go run cmd/processor/main.go`
5. `go run cmd/orders/main.go`

After that you'l be able to run commands and queries.

## Commands

### Create Order
```bash
curl --location --request POST '127.0.0.1:8080' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": "someid",
    "name": "somename",
    "status": "new"
}'
```

### Process Order
```bash
curl --location --request POST '127.0.0.1:8080/someid' --data-raw ''
```

## Queries

### Get Orders
```bash
curl --location --request GET '127.0.0.1:8080'
```

### Get Order
```bash
curl --location --request GET '127.0.0.1:8080/someid'
```
