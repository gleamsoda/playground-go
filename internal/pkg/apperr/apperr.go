package apperr

import "github.com/morikuni/failure"

const (
	Internal         failure.StringCode = "Internal"
	InvalidArgument  failure.StringCode = "InvalidArgument"
	NotFound         failure.StringCode = "NotFound"
	AlreadyExists    failure.StringCode = "AlreadyExists"
	Unauthorized     failure.StringCode = "Unauthorized"
	PermissionDenied failure.StringCode = "PermissionDenied"
)
