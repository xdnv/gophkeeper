package domain

import (
	"crypto/rsa"
	"database/sql"
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
)

// defines main session storage type based on server config given
type StorageType int

// session storage type
const (
	Memory StorageType = iota
	File
	Database
)

// return session storage type as string value
func (t StorageType) String() string {
	switch t {
	case Memory:
		return "Memory"
	case File:
		return "File"
	case Database:
		return "Database"
	}
	return fmt.Sprintf("Unknown (%d)", t)
}

// server configuration
type ServerConfig struct {
	TransportMode            string           `json:"transport_mode,omitempty"`   // data exchange transport mode: http or grpc
	Endpoint                 string           `json:"address,omitempty"`          // the address:port endpoint for server to listen
	MaxFileMemory            int64            `json:"max_file_memory,omitempty"`  // max memory size in MB to process files uploaded
	StorageMode              StorageType      `json:""`                           // session storage type
	MaxConnectionRetries     uint64           `json:""`                           // max connection retries to storage objects
	MockMode                 bool             `json:""`                           // pgSQL mock mode for test purposes
	Mock                     *sqlmock.Sqlmock `json:""`                           // pgSQL mock instance for test purposes
	MockConn                 *sql.DB          `json:""`                           // pgSQL mock connection for test purposes
	DatabaseDSN              string           `json:"database_dsn,omitempty"`     // database DSN (format: 'host=<host> [port=port] user=<user> password=<xxxx> dbname=<mydb> sslmode=disable')
	LogLevel                 string           `json:"log_level,omitempty"`        // log level
	CompressReplies          bool             `json:"compress_replies,omitempty"` // compress server replies, boolean
	CompressibleContentTypes []string         `json:""`                           // array of compressible mime types
	SessionCryptoKey         *rsa.PrivateKey  `json:""`                           // RSA 4096 session key used for encryption
	CryptoKeyPath            string           `json:"crypto_key,omitempty"`       // path to private crypto key (to decrypt messages from client)
	TrustedSubnet            string           `json:"trusted_subnet,omitempty"`   // trusted agent subnet in CIDR form. use empty value to disable security check.
	ConfigFilePath           string           `json:""`                           //path to JSON config file
}
