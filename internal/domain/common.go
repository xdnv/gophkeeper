// common configuration parts
package domain

// Protocol scheme to use in application
// WARNING: security risk, use https/SSL in production code
const PROTOCOL_SCHEME = "http"

// Endpoint default
const ENDPOINT = "localhost:8080"

// Default loglevel
const LOGLEVEL = "info"

// Structure to be filled by common function to react in http or grpc handler
type HandlerStatus struct {
	Message    string
	Err        error
	HTTPStatus int
}

// Global context key type not to mess with other packages
type ctxKey string

const (
	CtxApp      ctxKey = "app"
	CtxUsername ctxKey = "username"
)

// Remote command IDs
const S_CMD_SYNC = "synchronize"
const S_CMD_UPDATE = "update"
const S_CMD_DELETE = "delete"

// Remote command to be executed on server
type RemoteCommand struct {
	Command   string   `db:"command" json:"command"`     // Command name
	Arguments []string `db:"arguments" json:"arguments"` // Optional command arguments
	Data      []byte   `db:"data" json:"data"`           // Optional data structure serialized to byte array
}
