package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ZiyanK/fraud-detection-system/consumers/logger"
	"github.com/ZiyanK/fraud-detection-system/consumers/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

var (
	log = logger.CreateLogger()
)

type PostTransactionInput struct {
	UserID int     `db:"user_id" json:"user_id"`
	Amount float64 `db:"amount" json:"amount"`
	Type   string  `db:"type" json:"type"`
}

func HandlerPostTransaction(c *gin.Context) {
	var body PostTransactionInput

	err := c.ShouldBindJSON(&body)
	if err != nil {
		log.Error("Error while reading request body for signup", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid body.",
		})
		return
	}

	valid := model.CheckValidTransactionType(body.Type)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid transaction type"})
		return
	}

	var transaction model.Transaction

	transaction.TransactionID = rand.Intn(100000)
	transaction.Source = "rest-api"
	transaction.UserID = body.UserID
	transaction.Amount = body.Amount
	transaction.Type = body.Type
	transaction.Timestamp = time.Now()

	err = transaction.Store()
	if err != nil {
		log.Error("Error while storing transaction: ", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

type Filters struct {
}

func HandlerGetTransactions(c *gin.Context) {
	// Get the limit query parameter from the URL
	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid limit value"})
		return
	}

	// Get the offset query parameter from the URL
	offsetStr := c.Query("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid offset value"})
		return
	}

	transactionType := c.Query("type")
	valid := model.CheckValidTransactionType(transactionType)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid transaction type"})
		return
	}

	userIDStr := c.Query("user_id")
	var userID int
	if len(userIDStr) > 0 {
		userID, err = strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid offset value"})
			return
		}
	}

	transactions, err := model.GetTransactions(limit, offset, userID, transactionType)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func HandlerGetTransaction(c *gin.Context) {
	transactionIDStr := c.Param("id")
	transactionID, err := strconv.Atoi(transactionIDStr)
	if err != nil {
		log.Info("Invalid transaction id: ", zap.Any("id", transactionIDStr))
		c.Status(http.StatusNotFound)
		return
	}

	var transaction model.Transaction
	transaction.TransactionID = transactionID

	err = transaction.Get()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func HandlerDeleteTransaction(c *gin.Context) {
	transactionIDStr := c.Param("id")
	transactionID, err := strconv.Atoi(transactionIDStr)
	if err != nil {
		log.Info("Invalid transaction id: ", zap.Any("id", transactionIDStr))
		c.Status(http.StatusNotFound)
		return
	}

	var transaction model.Transaction
	transaction.TransactionID = transactionID

	deleted, err := transaction.Delete()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if !deleted {
		c.Status(http.StatusNotFound)
		return
	}

	c.Status(http.StatusOK)
}
