package main

import (
	"fmt"

	"github.com/ZiyanK/fraud-detection-system/consumers/db"
	"github.com/ZiyanK/fraud-detection-system/consumers/kafka"
	"github.com/ZiyanK/fraud-detection-system/consumers/logger"
	"github.com/ZiyanK/fraud-detection-system/consumers/route"
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

	reader := kafka.InitConsumer(config.KafkaUrl, config.KafkaTopic)
	defer reader.Close()
	log.Info("Connect and subscribed to topic successfully.")

	go kafka.ReadMessages(reader)

	router := route.AddRouter()
	err := router.Run(fmt.Sprintf(":%v", config.Port))
	if err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
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
	viper.AutomaticEnv()

	// Get configuration values
	config.Port = viper.GetString("PORT")
	config.DSN = viper.GetString("DSN")
	config.JWTSecret = viper.GetString("JWT_SECRET")
	config.KafkaUrl = viper.GetString("KAFKA_URL")
	config.KafkaTopic = viper.GetString("KAFKA_TOPIC")

	log.Info("config", zap.Any("config", config))

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("unable to decode into struct", zap.String("err", err.Error()))
		return err
	}

	return nil
}
