package handler

import (
	"context"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/timestamppb"

	"playground/internal/delivery/grpc/gen"
	"playground/internal/delivery/grpc/validator"
	"playground/internal/wallet"
)

func (c *Handler) CreateUser(ctx context.Context, req *gen.CreateUserRequest) (*gen.CreateUserResponse, error) {
	if violations := validateCreateUserRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	args := &wallet.CreateUserParams{
		Username: req.GetUsername(),
		FullName: req.GetFullName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
	u, err := c.w.CreateUser(ctx, args)
	if err != nil {
		return nil, err
	}

	rsp := &gen.CreateUserResponse{
		User: convertUser(u),
	}
	return rsp, nil
}

func validateCreateUserRequest(req *gen.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := validator.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	if err := validator.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}
	if err := validator.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}
	return violations
}

func (c *Handler) LoginUser(ctx context.Context, req *gen.LoginUserRequest) (*gen.LoginUserResponse, error) {
	if violations := validateLoginUserRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	meta := c.extractMetadata(ctx)
	args := &wallet.LoginUserParams{
		Username:  req.GetUsername(),
		Password:  req.GetPassword(),
		UserAgent: meta.UserAgent,
		ClientIP:  meta.ClientIP,
	}
	r, err := c.w.LoginUser(ctx, args)
	if err != nil {
		return nil, err
	}

	rsp := &gen.LoginUserResponse{
		User:                  convertUser(r.User),
		SessionId:             r.SessionID.String(),
		AccessToken:           r.AccessToken,
		RefreshToken:          r.RefreshToken,
		AccessTokenExpiresAt:  timestamppb.New(r.AccessTokenExpiresAt),
		RefreshTokenExpiresAt: timestamppb.New(r.RefreshTokenExpiresAt),
	}
	return rsp, nil
}

func validateLoginUserRequest(req *gen.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := validator.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	return violations
}

func (c *Handler) UpdateUser(ctx context.Context, req *gen.UpdateUserRequest) (*gen.UpdateUserResponse, error) {
	authPayload, err := c.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}
	if violations := validateUpdateUserRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	args := &wallet.UpdateUserParams{
		ReqUsername: authPayload.Username,
		Username:    req.GetUsername(),
		Password:    req.Password,
		FullName:    req.FullName,
		Email:       req.Email,
	}
	u, err := c.w.UpdateUser(ctx, args)
	if err != nil {
		return nil, err
	}

	rsp := &gen.UpdateUserResponse{
		User: convertUser(u),
	}
	return rsp, nil
}

func validateUpdateUserRequest(req *gen.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if req.Password != nil {
		if err := validator.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}
	if req.FullName != nil {
		if err := validator.ValidateFullName(req.GetFullName()); err != nil {
			violations = append(violations, fieldViolation("full_name", err))
		}
	}
	if req.Email != nil {
		if err := validator.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}
	return violations
}

func (c *Handler) VerifyEmail(ctx context.Context, req *gen.VerifyEmailRequest) (*gen.VerifyEmailResponse, error) {
	if violations := validateVerifyEmailRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	usr, err := c.w.VerifyEmail(ctx, &wallet.VerifyEmailParams{
		EmailID:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	})
	if err != nil {
		return nil, err
	}

	rsp := &gen.VerifyEmailResponse{
		IsVerified: usr.IsEmailVerified,
	}
	return rsp, nil
}

func validateVerifyEmailRequest(req *gen.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateEmailId(req.GetEmailId()); err != nil {
		violations = append(violations, fieldViolation("email_id", err))
	}
	if err := validator.ValidateSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, fieldViolation("secret_code", err))
	}
	return violations
}

func convertUser(user *wallet.User) *gen.User {
	return &gen.User{
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
