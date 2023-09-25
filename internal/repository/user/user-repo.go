package user

import (
	"context"
	"database/sql"
	"log"
	"paywise/internal/core"
	"paywise/internal/database/postgres"
	"paywise/internal/models"

	"github.com/lib/pq"
)

type userRepo struct {
	pg *postgres.PG
}

func New(pg postgres.DBTX) core.UserRepo {
	return &userRepo{
		pg: &postgres.PG{
			DB: pg,
		},
	}
}

func (ur *userRepo) Insert(ctx context.Context, user *models.User) (*models.User, error) {
	createdUser := new(models.User)
	err := ur.pg.DB.QueryRowContext(ctx, INSERT_USER_QUERY, user.Username, user.FullName, user.Email, user.HashedPassword).Scan(
		&createdUser.ID,
		&createdUser.Username,
		&createdUser.FullName,
		&createdUser.Email,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			// the code for duplicate key
			switch pqErr.Code {
			case "23505":
				return nil, core.DB_ERROR{Type: core.Duplicate_Value_Resource}
			case "23502":
				return nil, core.DB_ERROR{Type: core.Null_Value_Resource}
			}
		}
	}

	return createdUser, nil
}

func (ur *userRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user := new(models.User)
	err := ur.pg.DB.QueryRowContext(ctx, GET_BY_USERNAME_QUERY, username).Scan(&user.ID, &user.Username, &user.FullName, &user.Email, &user.HashedPassword)
	if err != nil {
		log.Println("i found an error =>", err)
		if err == sql.ErrNoRows {
			return nil, core.DB_ERROR{
				Type: core.Resource_Not_Found,
			}
		}
		return nil, core.DB_ERROR{
			Type: core.Internal_Db_Server,
		}
	}
	return user, nil
}

func (ur *userRepo) GetAllAccounts(ctx context.Context, username string) ([]*models.Account, error) {
	var accounts []*models.Account
	rows, err := ur.pg.DB.QueryContext(ctx, GET_ALL_ACCOUNTS_BY_USERNAME_QUERY, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core.DB_ERROR{
				Type: core.Internal_Db_Server,
			}
		}
	}
	for rows.Next() {
		acc := new(models.Account)
		err = rows.Scan(&acc.ID, &acc.OwnerName, &acc.Balance, &acc.Currency)
		if err != nil {
			return nil, core.DB_ERROR{
				Type: core.Internal_Db_Server,
			}
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

const (
	INSERT_USER_QUERY = `
		INSERT INTO users 
		(username, full_name, email, hashed_password)
		VALUES ($1, $2, $3, $4)
		RETURNING 
		id, username, full_name, email
	`

	GET_BY_USERNAME_QUERY = `
		SELECT id, username, full_name, email, hashed_password
		FROM users 
		WHERE username = $1
	`

	GET_ALL_ACCOUNTS_BY_USERNAME_QUERY = `
		SELECT id, owner_name, balance, currency 
		FROM accounts
		WHERE owner_name = $1
	`
)
