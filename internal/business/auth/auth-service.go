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

	"github.com/google/uuid"
)

type authService struct {
	userRepo    core.UserRepo
	sessionRepo core.SessionRepo
	tokenAuth   tokenConfig.TokenMaker
}

type AuthServiceConfig struct {
	UserRepo    core.UserRepo
	SessionRepo core.SessionRepo
	TokenAuth   tokenConfig.TokenMaker
}

func New(asc *AuthServiceConfig) core.AuthService {
	return &authService{
		userRepo:    asc.UserRepo,
		sessionRepo: asc.SessionRepo,
		tokenAuth:   asc.TokenAuth,
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

	Token, accessTokenPayload, err := as.tokenAuth.Create(reqDto.Username, configs.Paseto.Access_token_expiration)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshTokenPayload, err := as.tokenAuth.Create(reqDto.Username, configs.Paseto.Refresh_token_expiration)
	if err != nil {
		return nil, err
	}

	sessionID, err := uuid.NewUUID()
	if err != nil {
		log.Printf("error trying to create a new session id : %v \n", err)
		return nil, err
	}

	as.sessionRepo.CreateSession(ctx, &models.Session{
		ID:           sessionID,
		Username:     reqDto.Username,
		RefreshToken: refreshToken,
		IsBlocked:    false,
		ExpireAt:     refreshTokenPayload.ExpireAt,
	})

	// construct the response and return it
	response := &dtos.LoginRes{
		SessionID:              sessionID,
		AccessToken:            Token,
		AccessTokenExpiration:  accessTokenPayload.ExpireAt,
		RefreshToken:           refreshToken,
		RefreshTokenExpiration: refreshTokenPayload.ExpireAt,
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

	Token, accessTokenPayload, err := as.tokenAuth.Create(reqDto.Username, configs.Paseto.Access_token_expiration)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshTokenPayload, err := as.tokenAuth.Create(reqDto.Username, configs.Paseto.Refresh_token_expiration)
	if err != nil {
		return nil, err
	}

	sessionID, err := uuid.NewUUID()
	if err != nil {
		log.Printf("error trying to create a new session id : %v \n", err)
		return nil, err
	}

	as.sessionRepo.CreateSession(ctx, &models.Session{
		ID:           sessionID,
		Username:     reqDto.Username,
		RefreshToken: refreshToken,
		IsBlocked:    false,
		ExpireAt:     refreshTokenPayload.ExpireAt,
	})

	// construct the response and return it
	response := &dtos.LoginRes{
		SessionID:              sessionID,
		AccessToken:            Token,
		AccessTokenExpiration:  accessTokenPayload.ExpireAt,
		RefreshToken:           refreshToken,
		RefreshTokenExpiration: refreshTokenPayload.ExpireAt,
		User: &dtos.UserResDto{
			ID:       createdUser.ID,
			Username: createdUser.Username,
			FullName: createdUser.FullName,
			Email:    createdUser.Email,
		},
	}
	return response, nil
}
