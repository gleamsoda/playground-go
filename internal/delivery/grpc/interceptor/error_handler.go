package interceptor

import (
	"context"
	"fmt"

	"playground/internal/pkg/apperr"

	"github.com/morikuni/failure"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorHandler(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	result, err := handler(ctx, req)
	if err != nil {
		if code, ok := failure.CodeOf(err); ok {
			cs, _ := failure.CallStackOf(err)
			log.Error().
				Str("line", fmt.Sprintf("%s", cs.HeadFrame())).
				Msg(err.Error())
			switch code {
			case apperr.Internal:
				err = status.Error(codes.Internal, errorResponse(err))
			case apperr.InvalidArgument:
				err = status.Error(codes.InvalidArgument, errorResponse(err))
			case apperr.NotFound:
				err = status.Error(codes.NotFound, errorResponse(err))
			case apperr.AlreadyExists:
				err = status.Error(codes.AlreadyExists, errorResponse(err))
			case apperr.Unauthenticated:
				err = status.Error(codes.Unauthenticated, errorResponse(err))
			case apperr.PermissionDenied:
				err = status.Error(codes.PermissionDenied, errorResponse(err))
			default:
			}
		}
	}
	return result, err
}

func errorResponse(err error) string {
	msg := "something went wrong"
	if m, ok := failure.MessageOf(err); ok {
		msg = m
	}
	return msg
}
