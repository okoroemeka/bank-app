package worker

import (
	"context"
	"github.com/hibiken/asynq"
	db "github.com/okoroemeka/simple_bank/db/sqlc"
	"github.com/okoroemeka/simple_bank/mail"
	"github.com/rs/zerolog/log"
)

const (
	QueCritical = "critical"
	QueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
	mailer mail.EmailSender
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mailer mail.EmailSender) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			QueCritical: 10,
			QueDefault:  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().Err(err).Str("task", task.Type()).Bytes("payload", task.Payload()).Msg("error processing task")
		}),
		Logger: NewLogger(),
	})
	return &RedisTaskProcessor{server: server, store: store, mailer: mailer}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}
