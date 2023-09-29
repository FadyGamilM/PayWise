package user

import (
	"context"
	"log"
	"paywise/internal/core"
	"paywise/internal/core/dtos"
	"paywise/internal/models"
)

type userService struct {
	repo core.UserRepo
}

type UserServiceConfig struct {
	UserRepository core.UserRepo
}

func New(usc *UserServiceConfig) core.UserService {
	return &userService{
		repo: usc.UserRepository,
	}
}

func (us *userService) GetUserByUsername(ctx context.Context, reqDto *dtos.GetUserByUsernameDto) (*models.User, error) {
	user, err := us.repo.GetByUsername(ctx, reqDto.Username)
	if err != nil {
		log.Printf("[User Service] | %v \n", err)
		return nil, err
	}
	return user, nil
}

func (us *userService) GetAllAccountsOfUserByUsername(ctx context.Context, reqDto *dtos.GetAllAccountsForUserDto) ([]*models.Account, error) {
	accounts, err := us.repo.GetAllAccounts(ctx, reqDto.Username)
	if err != nil {
		log.Printf("[User Service] | %v \n", err)
		return nil, err
	}
	return accounts, nil
}
