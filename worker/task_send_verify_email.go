package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hibiken/asynq"
	custom_error "github.com/okoroemeka/simple_bank/custom-error"
	db "github.com/okoroemeka/simple_bank/db/sqlc"
	"github.com/okoroemeka/simple_bank/util"
	"github.com/rs/zerolog/log"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("cannot marshal payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)

	info, err := distributor.client.EnqueueContext(ctx, task)

	if err != nil {
		return fmt.Errorf("cannot enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("cannot unmarshal payload: %w", err)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if errors.Is(err, custom_error.ErrorNoRecordFound) {
			return fmt.Errorf("user does not exist: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("cannot get user: %w", err)
	}
	// TODO: send email to user

	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	})

	if err != nil {
		return fmt.Errorf("cannot create verify email: %w", err)
	}

	subject := "Welcome to Simple Bank"
	to := []string{user.Email}
	verificationUrl := fmt.Sprintf("http://localhost:8080/v1/verify_email?id=%d&code=%s", verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`Hi %s,<br/> 
		welcome to Simple Bank.<br/> 
		Please use the following link <a href="%s">click here</a> to verify your email`, user.FullName, verificationUrl)

	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)

	if err != nil {
		return fmt.Errorf("cannot send verify email: %w", err)
	}

	log.Info().Str("task", task.Type()).Bytes("payload", task.Payload()).Str("email", user.Email).Msgf("processing task")
	return nil
}
