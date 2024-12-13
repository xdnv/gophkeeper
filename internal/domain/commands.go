// Client-server commands configuration
package domain

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
