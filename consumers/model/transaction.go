package model

import (
	"time"

	"github.com/ZiyanK/fraud-detection-system/consumers/db"
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
	IsFraud       bool      `db:"is_fraud" json:"is_fraud"`
	Timestamp     time.Time `db:"timestamp" json:"timestamp"`
}

func (t *Transaction) CheckFraud() {
	if t.Amount > 500 {
		t.IsFraud = true
	}
}

func SendAlert(T Transaction) {
	log.Info("Fraud alert: %+v \n", zap.Any("transaction: ", T))
}

const (
	queryInsertTransaction = `
	INSERT INTO transactions (transaction_id, user_id, amount, type, is_fraud, timestamp)
	VALUES (:transaction_id, :user_id, :amount, :type, :is_fraud, :timestamp)`
)

func (T *Transaction) Store() error {
	_, err := db.DB.Sqlx.NamedExec(queryInsertTransaction, T)
	if err != nil {
		log.Error("Error inserting transaction in db: ", zap.Error(err))
		return err
	}

	return nil
}
