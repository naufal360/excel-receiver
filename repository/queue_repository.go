package repository

import (
	"context"
	"encoding/json"
	"excel-receiver/config"
	"excel-receiver/constant"
	"excel-receiver/entity"
	"excel-receiver/provider"

	"github.com/go-stomp/stomp/v3"
	"github.com/go-stomp/stomp/v3/frame"
)

type queueRepository struct {
	artemis *stomp.Conn
	log     provider.ILogger
}

func NewQueueArtemis(artemis *stomp.Conn, log provider.ILogger) *queueRepository {
	return &queueRepository{
		artemis: artemis,
		log:     log,
	}
}

func (q *queueRepository) ProduceQueue(ctx context.Context, requestData *entity.Queue) (err error) {
	var (
		address     = config.Configuration.Artemis.Address
		contentType = "application/json"
		reqID       = ctx.Value(constant.RequestIDKey{}).(string)
	)

	payload, err := json.Marshal(requestData)
	if err != nil {
		q.log.ErrorWithFields(provider.AppLog, map[string]interface{}{
			constant.ReqIDLog: reqID,
			"ERROR":           err,
		}, "error marshal json")

		return err
	}

	headers := []func(*frame.Frame) error{
		stomp.SendOpt.Header("destination-type", "ANYCAST"),
		stomp.SendOpt.Header("request-id", reqID),
		stomp.SendOpt.Header("persistent", "true"),
		stomp.SendOpt.Header("receipt", "true"),
	}

	err = q.artemis.Send(address, contentType, payload, headers...)

	if err != nil {
		q.log.ErrorWithFields(provider.AmqLog, map[string]interface{}{
			constant.ReqIDLog: reqID,
			"ERROR":           err,
		}, "failed to produce queue")

		return err
	}

	q.log.InfoWithFields(
		provider.AppLog,
		map[string]interface{}{
			constant.ReqIDLog: reqID,
			"QUEUE_TARGET":    address,
			"QUEUE_DATA":      string(payload),
		}, "success produce queue")

	return nil
}
