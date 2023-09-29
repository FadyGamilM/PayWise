package user

import (
	"net/http"
	"paywise/internal/business/auth/token"
	"paywise/internal/core"
	"paywise/internal/core/dtos"
	"paywise/internal/transport/rest/middlewares"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service core.UserService
}

type UserHandlerConfig struct {
	R             *gin.Engine
	UserService   core.UserService
	TokenProvider token.TokenMaker
}

func New(uhc *UserHandlerConfig) *UserHandler {
	h := &UserHandler{service: uhc.UserService}

	userRoutes := uhc.R.Group("/api/users").Use(middlewares.Authenticate(uhc.TokenProvider))

	userRoutes.GET("/accounts", h.HandleGetAllUserAccounts)

	return h
}

func (uh *UserHandler) HandleGetAllUserAccounts(c *gin.Context) {
	// => the username must be extracted from the authorization token not from the request body anymore ..
	var reqDto dtos.GetAllAccountsForUserDto
	// get the payload from the request context set by the middleware
	payload := c.MustGet(middlewares.AUTHORIZATION_PAYLOAD_CTX_KEY).(*token.Payload)
	// set the reqDto fields
	reqDto.Username = payload.Username
	// TODO => make a more robust error handler for each layer to specify the error more accurate
	accounts, err := uh.service.GetAllAccountsOfUserByUsername(c, &reqDto)
	if err != nil {
		appErr := core.NewInternalServerError()
		c.JSON(appErr.StatusCode(), gin.H{
			"error": appErr,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": accounts,
	})
}
