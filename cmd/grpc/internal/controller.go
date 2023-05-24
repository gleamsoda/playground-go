package internal

import (
	"context"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"playground/cmd/grpc/internal/boundary"
	"playground/domain"
	"playground/internal/validator"
)

type Controller struct {
	boundary.UnimplementedPlaygroundServer
	userUsecase domain.UserUsecase
}

func NewController(userUsecase domain.UserUsecase) *Controller {
	return &Controller{
		userUsecase: userUsecase,
	}
}

func (c *Controller) CreateUser(ctx context.Context, req *boundary.CreateUserRequest) (*boundary.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	args := domain.CreateUserInputParams{
		Username: req.GetUsername(),
		FullName: req.GetFullName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
	u, err := c.userUsecase.Create(ctx, args)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsp := &boundary.CreateUserResponse{
		User: convertUser(u),
	}
	return rsp, nil
}

func validateCreateUserRequest(req *boundary.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
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

func (c *Controller) LoginUser(ctx context.Context, req *boundary.LoginUserRequest) (*boundary.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	meta := c.extractMetadata(ctx)
	args := domain.LoginUserInputParams{
		Username:  req.GetUsername(),
		Password:  req.GetPassword(),
		UserAgent: meta.UserAgent,
		ClientIP:  meta.ClientIP,
	}
	r, err := c.userUsecase.Login(ctx, args)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	rsp := &boundary.LoginUserResponse{
		User:                  convertUser(r.User),
		SessionId:             r.SessionID.String(),
		AccessToken:           r.AccessToken,
		RefreshToken:          r.RefreshToken,
		AccessTokenExpiresAt:  timestamppb.New(r.AccessTokenExpiresAt),
		RefreshTokenExpiresAt: timestamppb.New(r.RefreshTokenExpiresAt),
	}
	return rsp, nil
}

func validateLoginUserRequest(req *boundary.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := validator.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}

func convertUser(user *domain.User) *boundary.User {
	return &boundary.User{
		Id:        user.ID,
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}

func (s *Controller) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if clientIPs := md.Get(xForwardedForHeader); len(clientIPs) > 0 {
			mtdt.ClientIP = clientIPs[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}

func fieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func invalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	badRequest := &errdetails.BadRequest{FieldViolations: violations}
	statusInvalid := status.New(codes.InvalidArgument, "invalid parameters")

	statusDetails, err := statusInvalid.WithDetails(badRequest)
	if err != nil {
		return statusInvalid.Err()
	}

	return statusDetails.Err()
}
