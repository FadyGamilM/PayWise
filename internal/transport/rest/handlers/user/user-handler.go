package user

import (
	"net/http"
	"paywise/internal/core"
	"paywise/internal/core/dtos"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service core.UserService
}

type UserHandlerConfig struct {
	R           *gin.Engine
	UserService core.UserService
}

func New(uhc *UserHandlerConfig) *UserHandler {
	h := &UserHandler{service: uhc.UserService}

	userRoutes := uhc.R.Group("/api/users")

	userRoutes.GET("/accounts", h.HandleGetAllUserAccounts)
	userRoutes.POST("", h.HandleCreateUser)

	return h
}

func (uh *UserHandler) HandleGetAllUserAccounts(c *gin.Context) {
	var reqDto dtos.GetAllAccountsForUserDto
	if err := c.ShouldBind(&reqDto); err != nil {
		appErr := core.NewBadRequestError()
		c.JSON(appErr.StatusCode(), gin.H{
			"error": appErr,
		})
	}
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

func (uh *UserHandler) HandleCreateUser(c *gin.Context) {
	var reqDto dtos.CreateUserDto
	if err := c.ShouldBind(&reqDto); err != nil {
		appErr := core.NewBadRequestError()
		c.JSON(appErr.StatusCode(), gin.H{
			"error": appErr,
		})
	}

	// TODO => make a more robust error handler for each layer to specify the error more accurate
	user, err := uh.service.Create(c, &reqDto)
	if err != nil {
		appErr := core.NewInternalServerError()
		c.JSON(appErr.StatusCode(), gin.H{
			"error": appErr,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}
