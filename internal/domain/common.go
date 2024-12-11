// common configuration parts
package domain

// Protocol scheme to use in application
//WARNING: security risk, use https/SSL in production code
const PROTOCOL_SCHEME = "http"

// Endpoint default
const ENDPOINT = "localhost:8080"

// Default loglevel
const LOGLEVEL = "info"

//structure to be filled by common function to react in http or grpc handler
type HandlerStatus struct {
	Message    string
	Err        error
	HTTPStatus int
}

//global context key type not to mess with other packages
type ctxKey string

const (
	CtxApp      ctxKey = "app"
	CtxUsername ctxKey = "username"
)
