package apperrors

import "github.com/morikuni/failure"

const (
	NotFound        failure.StringCode = "NotFound"
	InvalidArgument failure.StringCode = "InvalidArgument"
	Internal        failure.StringCode = "Internal"
	Unauthorized    failure.StringCode = "Unauthorized"
)
