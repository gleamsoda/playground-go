package usecase

import (
	"context"
	"fmt"

	"playground/app"
	"playground/app/mq"
)

func (u *Usecase) SendVerifyEmail(ctx context.Context, args *mq.SendVerifyEmailPayload) (*app.VerifyEmail, error) {
	usr, err := u.r.GetUser(ctx, args.Username)
	if err != nil {
		return nil, err
	}

	ve, err := u.r.CreateVerifyEmail(ctx, app.NewVerifyEmail(
		usr.Username,
		usr.Email,
		app.RandomString(32),
	))
	if err != nil {
		return nil, err
	}

	subject := "Welcome to Simple Bank"
	// TODO: replace this URL with an environment variable that points to a front-end page
	verifyUrl := fmt.Sprintf("http://localhost:8080/v1/verify_email?email_id=%d&secret_code=%s", ve.ID, ve.SecretCode)
	content := fmt.Sprintf(VerifyEmailContentFormat, usr.FullName, verifyUrl)
	to := []string{usr.Email}

	err = u.mailer.Send(subject, content, to, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to send verify email: %w", err)
	}
	return ve, nil
}

var VerifyEmailContentFormat = `Hello %s,<br/>
Thank you for registering with us!<br/>
Please <a href="%s">click here</a> to verify your email address.<br/>
`

func (u *Usecase) VerifyEmail(ctx context.Context, args *app.VerifyEmailParams) (*app.User, error) {
	usr, _, err := u.r.UpdateUserEmailVerified(ctx, args)
	return usr, err
}
