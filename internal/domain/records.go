package domain

import (
	"fmt"
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

// Keeper records to transparently use in serialization
type KeeperRecords []KeeperRecord
