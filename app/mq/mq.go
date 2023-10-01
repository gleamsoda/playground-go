package mq

import "context"

const SendVerifyEmailQueue = "task:send_verify_email"

type SendVerifyEmailPayload struct {
	Username string `json:"username" binding:"required,alphanum"`
}

type Producer interface {
	SendVerifyEmail(ctx context.Context, payload *SendVerifyEmailPayload) error
}
