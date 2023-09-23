package account

import (
	"net/http"
	"paywise/internal/core"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	service core.AccountService
}

type AccountHandlerConfig struct {
	R       *gin.Engine
	Service core.AccountService
}

func New(ahc *AccountHandlerConfig) *AccountHandler {
	accountRoutes := ahc.R.Group("/api/accounts")
	h := &AccountHandler{
		service: ahc.Service,
	}
	accountRoutes.GET("/", h.HandleGetAll)
	return h
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
