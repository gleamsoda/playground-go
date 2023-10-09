package helper

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
)

var AuthCtxkey = &ctxkey{"auth"}
var MetadataCtxkey = &ctxkey{"metadata"}

type ctxkey struct {
	name string
}

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func ErrorResponse(err error) map[string]any {
	return map[string]any{"error": err.Error()}
}
