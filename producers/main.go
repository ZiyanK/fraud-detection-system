package main

import (
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	topic = "transaction"
)

type Transaction struct {
	TransactionID int       `db:"transaction_id" json:"transaction_id"`
	UserID        int       `db:"user_id" json:"user_id"`
	Amount        float64   `db:"amount" json:"amount"`
	Type          string    `db:"type" json:"type"`
	IsFraud       bool      `db:"is_fraud" json:"is_fraud"`
	Source        string    `db:"source" json:"source"`
	Timestamp     time.Time `db:"timestamp" json:"timestamp"`
}

func generateTransaction() Transaction {
	TransactionTypes := []string{"purchase", "transfer", "payment"}

	return Transaction{
		TransactionID: rand.Intn(100000),
		UserID:        rand.Intn(10000),
		Amount:        math.Round(rand.Float64() * 1000),
		Type:          TransactionTypes[rand.Intn(len(TransactionTypes))],
		Source:        "message-broker",
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
