package kafka

import (
	"context"
	"encoding/json"

	"github.com/ZiyanK/fraud-detection-system/consumers/logger"
	"github.com/ZiyanK/fraud-detection-system/consumers/model"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

var (
	log = logger.CreateLogger()
)

func InitConsumer(kafkaUrl, kafkaTopic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaUrl},
		GroupID: "detection_group",
		Topic:   kafkaTopic,
	})
}

func ReadMessages(reader *kafka.Reader) {
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err == nil {
			var transaction model.Transaction
			if err := json.Unmarshal(msg.Value, &transaction); err != nil {
				log.Error("Failed to unmarshal transaction: ", zap.Error(err))
				continue
			}

			transaction.CheckFraud()

			if transaction.IsFraud {
				model.SendAlert(transaction)
			}

			transaction.Store()
		} else {
			log.Error("Consumer error: ", zap.Error(err))
		}
	}
}
