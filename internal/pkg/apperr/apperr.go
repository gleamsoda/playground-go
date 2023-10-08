package apperr

import "github.com/morikuni/failure"

const (
	Internal         failure.StringCode = "Internal"
	InvalidArgument  failure.StringCode = "InvalidArgument"
	NotFound         failure.StringCode = "NotFound"
	AlreadyExists    failure.StringCode = "AlreadyExists"
	Unauthenticated  failure.StringCode = "Unauthenticated"
	PermissionDenied failure.StringCode = "PermissionDenied"
)
