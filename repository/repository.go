package repository

import (
	"context"
	"excel-receiver/entity"
)

type TokenInterface interface {
	GetTokenAuthentication(ctx context.Context, token string) (*entity.Token, error)
}

type RequestInterface interface {
	CreateRequest(ctx context.Context, payload *entity.Request) (status string, err error)
}

type QueueInterface interface {
	ProduceQueue(ctx context.Context, requestData *entity.Queue) (err error)
}
