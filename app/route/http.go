package route

import (
	"github.com/ZiyanK/fraud-detection-system/consumers/handler"
	"github.com/gin-gonic/gin"
)

const (
	pathPing = "/ping"

	pathTransaction   = "/transaction"
	pathTransactionID = "/transaction/:id"
)

func AddRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.GET(pathPing, func(c *gin.Context) {
		c.JSON(200, "pong")
	})

	router.POST(pathTransaction, handler.HandlerPostTransaction)
	router.GET(pathTransaction, handler.HandlerGetTransactions)
	router.GET(pathTransactionID, handler.HandlerGetTransaction)
	router.DELETE(pathTransactionID, handler.HandlerDeleteTransaction)

	return router
}
