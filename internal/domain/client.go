package domain

import "crypto/rsa"

// client configuration
type ClientConfig struct {
	TransportMode        string         `json:"transport_mode,omitempty"`    // data exchange transport mode: http or grpc
	Endpoint             string         `json:"address,omitempty"`           // the address:port server endpoint to send metric data
	ReportInterval       int64          `json:"report_interval,omitempty"`   // metric reporting frequency in seconds
	PollInterval         int64          `json:"poll_interval,omitempty"`     // metric poll interval in seconds
	LogLevel             string         `json:"log_level,omitempty"`         // log verbosity (log level)
	APIVersion           string         `json:""`                            // API version to send metric data. Recent is v2
	UseCompression       bool           `json:""`                            // activate gzip compression
	BulkUpdate           bool           `json:""`                            // activate bulk JSON metric update
	MaxConnectionRetries uint64         `json:""`                            // Connection retries for retriable functions (does not include original request. 0 to disable)
	UseRateLimit         bool           `json:""`                            // flag option to enable or disable rate limiter
	RateLimit            int64          `json:"rate_limit,omitempty"`        // max simultaneous connections to server (rate limit)
	MessageSignature     string         `json:"message_signature,omitempty"` // key to use signed messaging, empty value disables signing
	AuthToken            string         `json:""`                            // session-wise auth token
	SessionSigningKey    *rsa.PublicKey `json:""`                            // RSA 512 session server-side key used for signing
	CryptoKeyPath        string         `json:"crypto_key,omitempty"`        // path to public crypto key (to encrypt messages to server)
	ConfigFilePath       string         `json:""`                            //path to JSON config file
}
