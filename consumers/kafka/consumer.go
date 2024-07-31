package kafka

import (
	"encoding/json"

	"github.com/ZiyanK/fraud-detection-system/consumers/logger"
	"github.com/ZiyanK/fraud-detection-system/consumers/model"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

var (
	log = logger.CreateLogger()
)

func InitConsumer(kafkaUrl, kafkaTopic string) (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaUrl,
		"group.id":          "detection_group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Error("Failed to create consumer: ", zap.Error(err))
		return nil, err
	}

	err = consumer.Subscribe(kafkaTopic, nil)
	if err != nil {
		log.Error("Failed to subscribe to topic.", zap.String("topic", kafkaTopic), zap.Error(err))
		return nil, err
	}

	return consumer, nil
}

func ReadMessages(consumer *kafka.Consumer) {
	for {
		msg, err := consumer.ReadMessage(-1)
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

		} else {
			log.Error("Consumer error: ", zap.Error(err))
		}
	}
}
