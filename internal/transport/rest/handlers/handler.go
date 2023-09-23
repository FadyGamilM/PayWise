package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	// all services
}

type HandlerConfig struct {
	R *gin.Engine
	// all services
}

func New(hc *HandlerConfig) *Handler {
	h := &Handler{}
	accountHandler := hc.R.Group("/api/accounts")
	accountHandler.GET("", h.HandleGetAllAccounts)
	return h
}

func (h *Handler) HandleGetAllAccounts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": "all accounts",
	})
}

func (h *Handler) HandleGetAccountByID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": "all accounts",
	})
}



func (h *Handler) HandleGetAccountByOwnerName(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": "all accounts",
	})
}
