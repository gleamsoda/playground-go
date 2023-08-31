package grpc

import (
	"context"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"playground/app"
	"playground/driver/grpc/gen"
	"playground/driver/grpc/validator"
)

func (c *Controller) CreateUser(ctx context.Context, req *gen.CreateUserRequest) (*gen.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	args := &app.CreateUserParams{
		Username: req.GetUsername(),
		FullName: req.GetFullName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
	u, err := c.u.CreateUser(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
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

func (c *Controller) LoginUser(ctx context.Context, req *gen.LoginUserRequest) (*gen.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	meta := c.extractMetadata(ctx)
	args := &app.LoginUserParams{
		Username:  req.GetUsername(),
		Password:  req.GetPassword(),
		UserAgent: meta.UserAgent,
		ClientIP:  meta.ClientIP,
	}
	r, err := c.u.LoginUser(ctx, args)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find user")
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

func convertUser(user *app.User) *gen.User {
	return &gen.User{
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
