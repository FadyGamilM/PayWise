package transactions

import (
	"log"
	"net/http"
	"paywise/internal/core"
	"paywise/internal/core/dtos"

	"github.com/gin-gonic/gin"
)

type MoneyTxHandler struct {
	service core.TransactionService
}

type MoneyTxHandlerConfig struct {
	R       *gin.Engine
	Service core.TransactionService
}

func New(mthc *MoneyTxHandlerConfig) *MoneyTxHandler {
	h := &MoneyTxHandler{
		service: mthc.Service,
	}

	moneyTxRoutes := mthc.R.Group("/api/money_transaction")

	moneyTxRoutes.POST("", h.HandleTransferMoney)

	return h
}

func (h *MoneyTxHandler) HandleTransferMoney(c *gin.Context) {
	var reqDto dtos.TxTransferMoneyReq
	if err := c.ShouldBind(&reqDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": core.NewBadRequestError(),
		})
	}

	log.Println("[HANDLER LAYER] | dto => ", reqDto.ToAccountID, " , ", reqDto.Amount, " , ", reqDto.ToAccountID)

	result, err := h.service.TransferMoneyTransaction(c, &reqDto)
	if err != nil {
		appErr := core.NewInternalServerError()
		c.JSON(appErr.StatusCode(), gin.H{
			"error": appErr,
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": result,
	})
}
