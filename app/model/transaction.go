package model

import (
	"strconv"
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
	Type          string    `db:"type" json:"type"`
	IsFraud       bool      `db:"is_fraud" json:"is_fraud"`
	Source        string    `db:"source" json:"source"`
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
	INSERT INTO transactions (transaction_id, user_id, amount, type, is_fraud, source, timestamp)
	VALUES (:transaction_id, :user_id, :amount, :type, :is_fraud, :source, :timestamp)`

	queryGetTransaction = `
	SELECT transaction_id, user_id, amount, type, is_fraud, source, timestamp
	FROM transactions
	WHERE transaction_id = $1`

	queryDeleteTransaction = `
	DELETE FROM transactions
	WHERE transaction_id = $1`
)

var (
	TransactionTypes = []string{"purchase", "transfer", "payment", ""}
)

func (T *Transaction) Store() error {
	_, err := db.DB.Sqlx.NamedExec(queryInsertTransaction, T)
	if err != nil {
		log.Error("Error inserting transaction in db: ", zap.Error(err))
		return err
	}

	return nil
}

func (T *Transaction) Get() error {
	err := db.DB.Sqlx.Get(T, queryGetTransaction, T.TransactionID)
	if err != nil {
		log.Error("Error fetching transaction in db: ", zap.Error(err))
		return err
	}

	return nil
}

func (T *Transaction) Delete() (bool, error) {
	result, err := db.DB.Sqlx.Exec(queryDeleteTransaction, T.TransactionID)
	if err != nil {
		log.Error("Error deleting transaction in db: ", zap.Error(err))
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Error("Error checking rows affected: ", zap.Error(err))
		return false, err
	}

	return rows == 1, nil
}

func GetTransactions(limit, offset, userID int, transactionType string) ([]Transaction, error) {
	var filterValues []interface{}

	queryGetTransactions := `
	SELECT transaction_id, user_id, amount, type, is_fraud, source, timestamp
	FROM transactions	
	WHERE 1=1`

	if userID > 0 {
		filterValues = append(filterValues, userID)
		queryGetTransactions += ` AND user_id = $` + strconv.Itoa(len(filterValues))
	}
	if len(transactionType) > 0 {
		filterValues = append(filterValues, transactionType)
		queryGetTransactions += ` AND type = $` + strconv.Itoa(len(filterValues))
	}
	if limit > 0 {
		filterValues = append(filterValues, limit)
		queryGetTransactions += ` LIMIT $` + strconv.Itoa(len(filterValues))
	}
	if offset > 0 {
		filterValues = append(filterValues, offset)
		queryGetTransactions += ` OFFSET $` + strconv.Itoa(len(filterValues))
	}

	var transactions []Transaction
	err := db.DB.Sqlx.Select(&transactions, queryGetTransactions, filterValues...)
	if err != nil {
		log.Error("Error fetching transaction in db: ", zap.Error(err))
		return nil, err
	}

	return transactions, nil
}

func CheckValidTransactionType(transactionType string) bool {
	for _, v := range TransactionTypes {
		if v == transactionType {
			return true
		}
	}
	return false
}
