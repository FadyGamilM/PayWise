package auth

import (
	"log"
	"net/http"
	"paywise/internal/core"
	"paywise/internal/core/dtos"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service core.AuthService
}

type AuthHandlerConfig struct {
	R           *gin.Engine
	AuthService core.AuthService
}

func New(ahc *AuthHandlerConfig) *AuthHandler {
	h := &AuthHandler{service: ahc.AuthService}

	authRoutes := ahc.R.Group("/api/auth")

	authRoutes.POST("/login", h.HandleLogin)
	authRoutes.POST("/signup", h.HandleSignup)

	return h
}

func (ah *AuthHandler) HandleSignup(c *gin.Context) {
	signupDto := new(dtos.CreateUserDto)
	if err := c.ShouldBind(signupDto); err != nil {
		log.Printf("THE ERROR CATCHED IN HANDLER LAYER IS | %v \n", err)
		appErr := core.NewBadRequestError()
		c.JSON(appErr.StatusCode(), gin.H{
			"error": appErr,
		})
	}
	// call the service layer
	response, err := ah.service.Signup(c, signupDto)
	if err != nil {
		// TODO => for now i will return internal server error, but this must be customized later
		appErr := core.NewInternalServerError()
		c.JSON(appErr.StatusCode(), gin.H{
			"error": appErr,
		})
	}
	// return the response
	c.JSON(http.StatusCreated, gin.H{
		"data": response,
	})
}

func (ah *AuthHandler) HandleLogin(c *gin.Context) {
	// extract the request body
	loginDto := new(dtos.LoginReq)
	if err := c.ShouldBind(loginDto); err != nil {
		log.Printf("THE ERROR CATCHED IN HANDLER LAYER IS | %v \n", err)
		appErr := core.NewBadRequestError()
		c.JSON(appErr.StatusCode(), gin.H{
			"error": appErr,
		})
	}

	// call the service layer
	response, err := ah.service.Login(c, loginDto)
	if err != nil {
		// TODO => for now i will return internal server error, but this must be customized later
		appErr := core.NewInternalServerError()
		c.JSON(appErr.StatusCode(), gin.H{
			"error": appErr,
		})
	}

	// return the response
	c.JSON(http.StatusCreated, gin.H{
		"data": response,
	})
}
