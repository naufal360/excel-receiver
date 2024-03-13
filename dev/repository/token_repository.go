package repository

import (
	"context"
	"database/sql"
	"errors"
	"excel-receiver/constant"
	"excel-receiver/entity"
	"excel-receiver/http/api/ierr"
	"excel-receiver/provider"

	"github.com/jmoiron/sqlx"
)

type TokenRepository struct {
	db  *sqlx.DB
	log provider.ILogger
}

func NewToken(db *sqlx.DB, log provider.ILogger) *TokenRepository {
	return &TokenRepository{
		db:  db,
		log: log,
	}
}

func (i *TokenRepository) GetTokenAuthentication(ctx context.Context, token string) (*entity.Token, error) {
	var (
		result        entity.Token
		getTokenQuery = `select id, token, expired_at, created_at FROM token WHERE token = ? limit 1`
	)

	err := i.db.GetContext(ctx, &result, getTokenQuery, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			i.log.Errorf(provider.DBLog, "token: %s is not found or not active", token)

			return nil, ierr.NewF(constant.APIUnauthorized, "")
		}

		i.log.Errorf(provider.DBLog, "failed to get token authentication")

		return nil, err
	}

	if result.IsExpired() {
		i.log.Errorf(provider.DBLog, "token: %s is expired", token)

		return nil, ierr.NewF(constant.APIUnauthorized, "")
	}

	return &result, nil
}
