package transactions

import (
	"log"
	"net/http"
	"paywise/internal/business/auth/token"
	"paywise/internal/core"
	"paywise/internal/core/dtos"
	"paywise/internal/transport/rest/middlewares"

	"github.com/gin-gonic/gin"
)

type MoneyTxHandler struct {
	service     core.TransactionService
	userService core.UserService
}

type MoneyTxHandlerConfig struct {
	R             *gin.Engine
	Service       core.TransactionService
	UserService   core.UserService
	TokenProvider token.TokenMaker
}

func New(mthc *MoneyTxHandlerConfig) *MoneyTxHandler {
	h := &MoneyTxHandler{
		service:     mthc.Service,
		userService: mthc.UserService,
	}

	moneyTxRoutes := mthc.R.Group("/api/money_transaction").Use(middlewares.Authenticate(mthc.TokenProvider))

	moneyTxRoutes.POST("", h.HandleTransferMoney)

	return h
}

// TODO => if the balance of the from account is zero, we must return an error, but our handler returns the created transfer and the created entries but doesn't update the balances and doesn't persist the created transfer and entiries (which is a good behavior from our db transaction management)
func (h *MoneyTxHandler) HandleTransferMoney(c *gin.Context) {
	var reqDto dtos.TxTransferMoneyReq
	if err := c.ShouldBind(&reqDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": core.NewBadRequestError(),
		})
	}

	log.Println("[HANDLER LAYER] | dto => ", reqDto.ToAccountID, " , ", reqDto.Amount, " , ", reqDto.FromAccountID)

	// authorization ==>
	// get the fromAccount using the fromAccountID
	// validate that the ownerName of this fromAccount is the same username of the current logged-in user because user can only transfer money from his account
	// get the payload from the request context set by the middleware
	payload := c.MustGet(middlewares.AUTHORIZATION_PAYLOAD_CTX_KEY).(*token.Payload)
	userAccounts, err := h.userService.GetAllAccountsOfUserByUsername(c, &dtos.GetAllAccountsForUserDto{
		Username: payload.Username,
	})
	if err != nil {
		appErr := core.NewInternalServerError()
		c.JSON(appErr.StatusCode(), gin.H{
			"error": appErr,
		})
		return
	}
	canTransfer := false
	for _, acc := range userAccounts {
		if acc.ID == reqDto.FromAccountID {
			canTransfer = true
			break
		}
	}

	if !canTransfer {
		appErr := core.NewUnAuthorizedError("user cannot transfer money from account rather than his/her accounts, you are not authroized to perform this action")
		c.JSON(appErr.StatusCode(), gin.H{
			"error": appErr,
		})
		return
	}

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
