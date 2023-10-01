package auth

import (
	"context"
	"fmt"
	"log"
	"paywise/config"
	tokenConfig "paywise/internal/business/auth/token"
	"paywise/internal/core"
	"paywise/internal/core/dtos"
	"paywise/internal/models"
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

// TODO => this method must run within the same transaction ..
func (as *authService) Login(ctx context.Context, reqDto *dtos.LoginReq) (*dtos.LoginRes, error) {
	// check if there is a registered user in our system with this username
	registeredUser, err := as.userRepo.GetByUsername(ctx, reqDto.Username)
	if err != nil {
		// TODO => customize the error later
		return nil, err
	}

	// check if the given password is the same as the stored hashed password

	IsCorrect := CheckPassword(reqDto.Password, registeredUser.HashedPassword)
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

	log.Println("the expiration date is => ", configs.Paseto.Access_token_expiration)
	Token, err := as.tokenAuth.Create(reqDto.Username, configs.Paseto.Access_token_expiration)
	if err != nil {
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

func (as *authService) Signup(ctx context.Context, reqDto *dtos.CreateUserDto) (*dtos.LoginRes, error) {
	// hash the password
	hashedPass, err := HashPassword(reqDto.Password)
	if err != nil {
		log.Printf("[Auth Service] | %v \n", err)
		return nil, err
	}
	createdUser, err := as.userRepo.Insert(ctx, &models.User{
		Username:       reqDto.Username,
		Email:          reqDto.Email,
		FullName:       reqDto.FullName,
		HashedPassword: hashedPass,
	})
	if err != nil {
		return nil, err
	}

	// load the secret key from the config.yaml
	configs, err := config.LoadPasetoTokenConfig("./config")
	if err != nil {
		fmt.Println("error trying to load config variables", err)
		return nil, err
	}

	Token, err := as.tokenAuth.Create(reqDto.Username, configs.Paseto.Access_token_expiration)
	if err != nil {
		return nil, err
	}

	// construct the response and return it
	response := &dtos.LoginRes{
		AccessToken: Token,
		User: &dtos.UserResDto{
			ID:       createdUser.ID,
			Username: createdUser.Username,
			FullName: createdUser.FullName,
			Email:    createdUser.Email,
		},
	}
	return response, nil
}
