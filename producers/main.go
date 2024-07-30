package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	topic = "transactions"
)

type Transaction struct {
	TransactionID int
	UserID        int
	Amount        float64
	Type          string
	Timestamp     time.Time
}

func generateTransaction() Transaction {
	TransactionTypes := []string{"purchase", "transfer", "payment"}

	return Transaction{
		TransactionID: rand.Intn(100000),
		UserID:        rand.Intn(10000),
		Amount:        rand.Float64() * 1000,
		Type:          TransactionTypes[rand.Intn(len(TransactionTypes))],
		Timestamp:     time.Now(),
	}
}

func sendToKafka(transaction Transaction, topic string) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:61104",
	})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	message, err := json.Marshal(transaction)
	if err != nil {
		log.Fatalf("Failed to marshal message : %s \n", err)
	}

	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)

	p.Flush(15 * 1000)
}

func main() {
	for {
		transaction := generateTransaction()
		log.Printf("Generated transaction: %+v \n", transaction)
		sendToKafka(transaction, topic)
		time.Sleep(time.Second * 1)
	}
}
