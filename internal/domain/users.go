package domain

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// holds seriazable user account details used for registration purposes
type UserAccountData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// internal server UAC record
type UserAccountRecord struct {
	ID               string    `db:"id" json:"id,omitempty"`                     // Unique user UID
	Login            string    `db:"login" json:"login"`                         // Unique user login
	Password         string    `db:"password" json:"password"`                   // password hash
	Email            string    `db:"email" json:"email"`                         // optional user e-mail
	RegistrationDate time.Time `db:"registration_date" json:"registration_date"` // Date of registration
	IsBanned         bool      `db:"is_banned" json:"is_banned,omitempty"`       // Whether user is banned
}

// internal server UAC metadata
type UserLoginMetadata struct {
	IP string `db:"ip" json:"ip,omitempty"` // IP-address
}

// client-wise AuthResponse structure with JWT token and crypto key
type AuthResponse struct {
	Token     string        `json:"token"`
	PublicKey rsa.PublicKey `json:"public_key"`
}

// JWT custom Claims to be used in auth token
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// auth bruteforce protector record type
type BlockedLogin struct {
	Login          string
	IP             string
	FailedAttempts int
	LastAttemptAt  time.Time
}

// auth bruteforce protector storage type
type BlockedLogins map[string]BlockedLogin
