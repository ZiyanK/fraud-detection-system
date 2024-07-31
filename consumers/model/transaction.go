package model

import (
	"time"

	"github.com/ZiyanK/fraud-detection-system/consumers/logger"
	"go.uber.org/zap"
)

var (
	log = logger.CreateLogger()
)

type Transaction struct {
	TransactionID int       `db:"transaction_id" json:"transaction_id"`
	UserID        int       `db:"user_id" json:"user_id"`
	Amount        float64   `db:"amount" json:"amount"`
	Type          string    `db:"type" json:"json"`
	Timestamp     time.Time `db:"timestamp" json:"timestamp"`
	IsFraud       bool      `db:"is_fraud" json:"is_fraud"`
}

func (t *Transaction) CheckFraud() {
	if t.Amount > 500 {
		t.IsFraud = true
	}
}

func SendAlert(T Transaction) {
	log.Info("Fraud alert: %+v \n", zap.Any("transaction: ", T))
}
