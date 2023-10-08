package wallet

import "context"

const SendVerifyEmailQueue = "task:send_verify_email"

type SendVerifyEmailPayload struct {
	Username string `json:"username" binding:"required,alphanum"`
}

type Dispatcher interface {
	SendVerifyEmail(ctx context.Context, payload *SendVerifyEmailPayload) error
}
