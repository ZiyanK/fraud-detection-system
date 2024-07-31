package main

import (
	"github.com/ZiyanK/fraud-detection-system/consumers/db"
	"github.com/ZiyanK/fraud-detection-system/consumers/kafka"
	"github.com/ZiyanK/fraud-detection-system/consumers/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	log    = logger.CreateLogger()
	config configuration
)

func main() {
	// load config file
	if err := LoadConfig(); err != nil {
		panic(err)
	}

	// Databases Init
	if err := db.InitConn(config.DSN); err != nil {
		log.Fatal("Failed to conenct to the database", zap.Error(err))
	}

	consumer, err := kafka.InitConsumer(config.KafkaUrl, config.KafkaTopic)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()
	log.Info("Connect and subscribed to topic successfully.")

	kafka.ReadMessages(consumer)
}

// configration is a struct used to get the environment variable from the config.yaml file
type configuration struct {
	DSN        string `mapstructure:"dsn"`
	Port       string `mapstructure:"port"`
	JWTSecret  string `mapstructure:"jwtSecret"`
	Mode       string `mapstructure:"mode"`
	KafkaUrl   string `mapstructure:"kafkaUrl"`
	KafkaTopic string `mapstructure:"kafkaTopic"`
}

// LoadConfig is used to load the configuration
func LoadConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file", zap.String("err", err.Error()))
		return err
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("unable to decode into struct", zap.String("err", err.Error()))
		return err
	}

	return nil
}
