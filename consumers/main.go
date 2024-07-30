package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Transaction struct {
	TransactionID int
	UserID        int
	Amount        float64
	Type          string
	Timestamp     time.Time
}

func detectFraud(T Transaction) bool {
	return T.Amount > 500
}

func sendAlert(T Transaction) {
	log.Printf("Fraud alert: %+v \n", T)
}

func main() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:61104",
		"group.id":          "detection_group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	err = c.Subscribe("transactions", nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connect and subscribed to topic successfully.")

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			var transaction Transaction
			if err := json.Unmarshal(msg.Value, &transaction); err != nil {
				log.Printf("Failed to unmarshal transaction: %s", err)
				continue
			}

			if detectFraud(transaction) {
				sendAlert(transaction)
			}

		} else {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
