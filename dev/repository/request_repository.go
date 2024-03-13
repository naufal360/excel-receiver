package repository

import (
	"context"
	"excel-receiver/constant"
	"excel-receiver/entity"
	"excel-receiver/provider"
	"time"

	"github.com/jmoiron/sqlx"
)

type RequestRepository struct {
	db  *sqlx.DB
	log provider.ILogger
}

func NewRequest(db *sqlx.DB, log provider.ILogger) *RequestRepository {
	return &RequestRepository{
		db:  db,
		log: log,
	}
}

func (r *RequestRepository) CreateRequest(ctx context.Context, payload *entity.Request) (status string, err error) {
	var (
		reqID    = ctx.Value(constant.RequestIDKey{}).(string)
		queryStr = `
			insert into request (
				request_id, status, file_path, created_at
			) values ( 
				:request_id, :status, :file_path, :created_at
			)
		`
	)

	r.log.InfoWithFields(provider.DBLog, map[string]interface{}{
		constant.ReqIDLog: reqID,
	}, "insert request")

	payload.CreatedAt = time.Now()

	query, args, err := r.db.BindNamed(queryStr, payload)
	if err != nil {
		r.log.ErrorWithFields(provider.DBLog, map[string]interface{}{
			constant.ReqIDLog: reqID,
		}, "failed to insert request: %s", err)
		return "failed", err
	}

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		r.log.ErrorWithFields(provider.DBLog, map[string]interface{}{
			constant.ReqIDLog: reqID,
		}, "failed to insert request: %s", err)

		return "failed", err
	}

	_, err = res.LastInsertId()
	if err != nil {
		r.log.ErrorWithFields(provider.DBLog, map[string]interface{}{
			constant.ReqIDLog: reqID,
		}, "failed to get last inserted id request: %s", err)

		return "failed", err
	}

	return "received", nil
}
