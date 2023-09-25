package auth

import (
	"context"
	"fmt"
	"log"
	"paywise/config"
	tokenConfig "paywise/internal/business/auth/token"
	userService "paywise/internal/business/user"
	"paywise/internal/core"
	"paywise/internal/core/dtos"
)

type authService struct {
	userRepo  core.UserRepo
	tokenAuth tokenConfig.TokenMaker
}

type AuthServiceConfig struct {
	UserRepo  core.UserRepo
	TokenAuth tokenConfig.TokenMaker
}

func New(asc *AuthServiceConfig) core.AuthService {
	return &authService{
		userRepo:  asc.UserRepo,
		tokenAuth: asc.TokenAuth,
	}
}

func (as *authService) Login(ctx context.Context, reqDto *dtos.LoginReq) (*dtos.LoginRes, error) {
	// check if there is a registered user in our system with this username
	registeredUser, err := as.userRepo.GetByUsername(ctx, reqDto.Username)
	if err != nil {
		// TODO => customize the error later
		return nil, err
	}

	log.Println("the registeredUser is ===> ", registeredUser.FullName)

	// check if the given password is the same as the stored hashed password

	IsCorrect := userService.CheckPassword(reqDto.Password, registeredUser.HashedPassword)
	if !IsCorrect {
		return nil, fmt.Errorf("invalid user credentials")
	}

	log.Println("the password is matched ==> ", IsCorrect)

	// load the secret key from the config.yaml
	configs, err := config.LoadPasetoTokenConfig("./config")
	if err != nil {
		fmt.Println("error trying to load config variables", err)
		return nil, err
	}

	log.Println("the username is | ", reqDto.Username)
	log.Println("the expiration is | ", configs.Paseto.Expiration)
	// if its a valid credential
	if as.tokenAuth == nil {
		log.Println("holy mother fucker")
	}
	Token, err := as.tokenAuth.Create(reqDto.Username, configs.Paseto.Expiration)
	if err != nil {
		log.Println("error is ===========> ", err)
		return nil, err
	}

	// construct the response and return it
	response := &dtos.LoginRes{
		AccessToken: Token,
		User: &dtos.UserResDto{
			ID:       registeredUser.ID,
			Username: registeredUser.Username,
			FullName: registeredUser.FullName,
			Email:    registeredUser.Email,
		},
	}

	return response, nil
}
