package usecase

import (
	"context"
	"fmt"

	"playground/internal/app"
	"playground/internal/pkg/mail"
)

type (
	SendVerifyEmail struct {
		r      app.Repository
		mailer mail.Sender
	}
	VerifyEmail struct {
		r app.Repository
	}
)

func NewSendVerifyEmail(r app.Repository, mailer mail.Sender) *SendVerifyEmail {
	return &SendVerifyEmail{
		r:      r,
		mailer: mailer,
	}
}

func (u *SendVerifyEmail) Execute(ctx context.Context, args *app.SendVerifyEmailPayload) (*app.VerifyEmail, error) {
	usr, err := u.r.User().Get(ctx, args.Username)
	if err != nil {
		return nil, err
	}

	var ve *app.VerifyEmail
	err = u.r.Transaction().Run(ctx, func(ctx context.Context, r app.Repository) error {
		var err error
		if ve, err = r.User().CreateVerifyEmail(ctx, app.NewVerifyEmail(
			usr.Username,
			usr.Email,
			app.RandomString(32),
		)); err != nil {
			return err
		}

		subject := "Welcome to Simple Bank"
		// TODO: replace this URL with an environment variable that points to a front-end page
		verifyUrl := fmt.Sprintf("http://localhost:8080/v1/verify_email?email_id=%d&secret_code=%s", ve.ID, ve.SecretCode)
		content := fmt.Sprintf(VerifyEmailContentFormat, usr.FullName, verifyUrl)
		to := []string{usr.Email}

		err = u.mailer.Send(subject, content, to, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to send verify email: %w", err)
		}
		return nil
	})
	return ve, err
}

var VerifyEmailContentFormat = `Hello %s,<br/>
Thank you for registering with us!<br/>
Please <a href="%s">click here</a> to verify your email address.<br/>
`

func NewVerifyEmail(r app.Repository) *VerifyEmail {
	return &VerifyEmail{
		r: r,
	}
}

func (u *VerifyEmail) Execute(ctx context.Context, args *app.VerifyEmailParams) (*app.User, error) {
	usr, _, err := u.r.User().UpdateEmailVerified(ctx, args)
	return usr, err
}
