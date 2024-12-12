package domain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Data structure of Keeper record in UI lists
type KeeperRecord struct {
	ListNR      string    `db:"" json:"list_nr,omitempty"`              // Dynamic unique list number, can be used in UI commands. Not stored in DB
	ShortID     string    `db:"" json:"short_id,omitempty"`             // Short unique record ID (first 8 chars) used in UI commands. Not stored in DB
	ID          string    `db:"id" json:"id,omitempty"`                 // Unique record ID
	UserID      string    `db:"user_id" json:"user_id,omitempty"`       // Unique user ID
	Name        string    `db:"name" json:"name"`                       // Short human-readable record name
	Description string    `db:"description" json:"description"`         // Record description/metadata
	SecretType  string    `db:"secret_type" json:"secret_type"`         // Record type (credentials/creditcard/text/binary/...)
	Created     time.Time `db:"created" json:"created,omitempty"`       // Date and time when the record was added
	Modified    time.Time `db:"modified" json:"modified,omitempty"`     // Last date and time when the record was modified
	IsDeleted   bool      `db:"is_deleted" json:"is_deleted,omitempty"` // Whether record was deleted by user
	Secret      string    `db:"secret" json:"secret,omitempty"`         // Secret part stored in JSON format. Not optimized for HUGE binary data.
}

func (k *KeeperRecord) Reference() string {
	return fmt.Sprintf("#%s. %s (%s)", k.ListNR, k.Name, k.ShortID)
}

// Keeper secret types
const SECRET_CREDENTIALS = "credentials"
const SECRET_CREDITCARD = "creditcard"
const SECRET_TEXT = "text"
const SECRET_BINARY = "binary"

// Keeper records to transparently use in serialization
type KeeperRecords []KeeperRecord

// Secret structure template
type SecretData interface {
	ToString() string
}

// universal human-readble output for Keeper* classes
func KeepReadable(s SecretData) string {
	return s.ToString()
}

// universal secret deserialzer
func KeepDeserialized(class string, data []byte) (*SecretData, error) {
	var object SecretData

	switch class {
	case SECRET_CREDENTIALS:
		object = new(SecretCredentials)
	case SECRET_CREDITCARD:
		object = new(SecretCreditcard)
	case SECRET_TEXT:
		object = new(SecretText)
	case SECRET_BINARY:
		object = new(KeeperBinary)
	default:
		return nil, fmt.Errorf("unknown secret type: %s", class)
	}

	err := json.Unmarshal(data, &object)
	if err != nil {
		return nil, err
	}
	return &object, nil
}

// Secret structure for datatype "credentials"
type SecretCredentials struct {
	Address  string `db:"address" json:"address"`   // URL to apply credentials
	Login    string `db:"login" json:"login"`       //
	Password string `db:"password" json:"password"` // URL to apply credentials
}

// String representation for "credentials" datatype
func (k SecretCredentials) ToString() string {
	return fmt.Sprintf("Address: %s\nLogin: %s\nPassword: %s", k.Address, k.Login, k.Password)
}

// Secret structure for datatype "creditcard"
type SecretCreditcard struct {
	CardNumber     string `db:"card_number" json:"card_number"`         // Card number in 4 blocks of digits
	ExpirationDate string `db:"expiration_date" json:"expiration_date"` // Expiration date in "MM/YY" format
	SecurityCode   string `db:"security_code" json:"security_code"`     // CVC/CVV code (3 digits)
}

// String representation for "creditcard" datatype
func (k SecretCreditcard) ToString() string {
	return fmt.Sprintf("Card Number: %s\nExpiration Date: %s\nCVV: %s", k.CardNumber, k.ExpirationDate, k.SecurityCode)
}

// Secret structure for datatype "text"
type SecretText struct {
	Text string `db:"text" json:"text"` // Secret text
}

// String representation for "text" datatype
func (k SecretText) ToString() string {
	return fmt.Sprintf("Text:\n%s", k.Text)
}

// Secret structure for datatype "binary"
type KeeperBinary struct {
	FileName  string `db:"file_name" json:"file_name"` // Binary file name
	Extension string `db:"extension" json:"extension"` // Binary file extension
	FileSize  int64  `db:"file_size" json:"file_size"` // Binary file size
	Data      []byte `db:"data" json:"data"`           // Binary file data
}

// String representation for "binary" datatype
func (k KeeperBinary) ToString() string {
	return fmt.Sprintf("File name:%s\nExtension:%s\nFile size: %v", k.FileName, k.Extension, k.FileSize)
}

// Create new KeepBinary and read file from disk
func NewBinarySecret(filePath string) (*KeeperBinary, error) {
	k := new(KeeperBinary)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return k, err
	}

	k.FileName = fileInfo.Name()
	k.Extension = filepath.Ext(k.FileName)
	k.FileSize = fileInfo.Size()
	k.Data, err = os.ReadFile(filePath)
	if err != nil {
		return k, err
	}
	return k, nil
}

// Dump binary file from
func DumpBinary(k *KeeperBinary, filePath string) error {
	if k.FileSize == 0 || len(k.Data) == 0 {
		return fmt.Errorf("empty binary data storage, nothing to save")
	}
	return os.WriteFile(filePath, k.Data, 0666)
}

// Remote command to be executed on server
type RemoteCommand struct {
	Command   string   `db:"command" json:"command"`     // Command name
	Arguments []string `db:"arguments" json:"arguments"` // Optional command arguments
	Data      []byte   `db:"data" json:"data"`           // Optional data structure serialized to byte array
}
