package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	config configuration
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

// Generates a new transaction
func generateTransaction() Transaction {
	TransactionTypes := []string{"purchase", "transfer", "payment"}

	transaction := Transaction{
		TransactionID: rand.Intn(100000),
		UserID:        rand.Intn(10000),
		Amount:        math.Round(rand.Float64() * 1000),
		Type:          TransactionTypes[rand.Intn(len(TransactionTypes))],
		Source:        "message-broker",
		Timestamp:     time.Now(),
	}

	log.Printf("Generated transaction: %+v \n", transaction)

	return transaction
}

func newKafkaWriter(kafkaUrl, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaUrl),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func main() {
	// load config file
	if err := LoadConfig(); err != nil {
		panic(err)
	}

	// Initialize a kafka writer
	writer := newKafkaWriter(config.KafkaUrl, config.KafkaTopic)
	defer writer.Close()

	// Run a endless loop to generate transactions to send to the broker
	for {
		transaction := generateTransaction()

		transactionByte, err := json.Marshal(transaction)
		if err != nil {
			fmt.Println(err)
		}

		msg := kafka.Message{
			Value: transactionByte,
		}

		// Write the given message to the broker
		err = writer.WriteMessages(context.Background(), msg)
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(time.Second * 1)
	}
}

// configration is a struct used to get the environment variable
type configuration struct {
	KafkaUrl   string `mapstructure:"kafkaUrl"`
	KafkaTopic string `mapstructure:"kafkaTopic"`
}

// LoadConfig is used to load the configuration
func LoadConfig() error {
	viper.AutomaticEnv()

	config.KafkaUrl = viper.GetString("KAFKA_URL")
	config.KafkaTopic = viper.GetString("KAFKA_TOPIC")

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("unable to decode into struct", zap.String("err", err.Error()))
		return err
	}

	return nil
}
