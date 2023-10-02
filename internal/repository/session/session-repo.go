package session

import (
	"context"
	"log"
	"paywise/internal/core"
	"paywise/internal/database/postgres"
	"paywise/internal/models"
)

type sessionRepo struct {
	pg *postgres.PG
}

func New(pg postgres.DBTX) core.SessionRepo {
	return &sessionRepo{
		pg: &postgres.PG{
			DB: pg,
		},
	}
}

const (
	CREATE_SESSION_QUERY = `
		INSERT INTO sessions (id, username, refresh_token, is_blocked, expire_at) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, username, refresh_token, is_blocked, expire_at 
	`
	GET_SESSION_BY_ID_QUERY = `
		SELECT username, refresh_token, is_blocked, expire_at 
		FROM sessions 
		WHERE id = $1
	`
)

func (sr *sessionRepo) CreateSession(ctx context.Context, session *models.Session) (*models.Session, error) {
	createdSession := new(models.Session)
	err := sr.pg.DB.QueryRowContext(ctx, CREATE_SESSION_QUERY, session.ID, session.Username, session.RefreshToken, session.IsBlocked, session.ExpireAt).Scan(
		&createdSession.ID,
		&createdSession.Username,
		&createdSession.RefreshToken,
		&createdSession.IsBlocked,
		&createdSession.ExpireAt,
	)
	if err != nil {
		// TODO => customize db error level
		log.Printf("error creating new session in database : %v\n", err)
		return nil, err
	}
	return createdSession, nil
}

func (sr *sessionRepo) GetBySessionID(ctx context.Context, sessionID int64) (*models.Session, error) {
	session := new(models.Session)
	err := sr.pg.DB.QueryRowContext(ctx, GET_SESSION_BY_ID_QUERY, sessionID).Scan(
		&session.ID,
		&session.Username,
		&session.RefreshToken,
		&session.IsBlocked,
		&session.ExpireAt,
	)
	if err != nil {
		// TODO => customize db error level
		log.Printf("error retreiving session from database : %v\n", err)
		return nil, err
	}
	return session, nil
}
