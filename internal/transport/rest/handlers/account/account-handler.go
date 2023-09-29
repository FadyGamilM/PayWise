package account

import (
	"log"
	"net/http"
	"paywise/internal/business/auth/token"
	"paywise/internal/core"
	"paywise/internal/core/dtos"
	"paywise/internal/transport/rest/middlewares"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	service core.AccountService
}

type AccountHandlerConfig struct {
	R             *gin.Engine
	Service       core.AccountService
	TokenProvider token.TokenMaker
}

func New(ahc *AccountHandlerConfig) *AccountHandler {
	accountRoutes := ahc.R.Group("/api/accounts").Use(middlewares.Authenticate(ahc.TokenProvider))
	h := &AccountHandler{
		service: ahc.Service,
	}
	accountRoutes.GET("", h.HandleGetAll)
	accountRoutes.POST("", h.HandleCreateAccount)
	return h
}

func (ah *AccountHandler) HandleCreateAccount(c *gin.Context) {
	reqDto := new(dtos.CreateAccReq)
	payload := c.MustGet(middlewares.AUTHORIZATION_PAYLOAD_CTX_KEY).(*token.Payload)
	acc, err := ah.service.Create(c, reqDto, payload.Username)
	if err != nil {
		log.Printf("err ==> ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "couldn't create an account for a user with username : " + payload.Username,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"data": acc,
	})
}

func (ah *AccountHandler) HandleGetAll(c *gin.Context) {

	accounts, err := ah.service.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": core.NewInternalServerError(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": accounts,
	})

}
