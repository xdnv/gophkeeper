package app

import (
	"sync"
	"time"

	"internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	blockedLogins = make(domain.BlockedLogins)
	mu            sync.Mutex
)

// Hash password provided by user
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Compare password with hashed one
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	//logger.Error(fmt.Sprintf("pwd compare error: %v\n", err)) //DEBUG
	return err == nil
}

// Generates user-specific JWT token valid for specific time in hours
func GenerateJWT(username string, hours int64) (string, error) {
	expirationTime := time.Now().Add(time.Duration(hours) * time.Hour)

	claims := &domain.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	return token.SignedString(Sc.SessionCryptoKey)
}

// Get string to be used in login blocker map
func GetLoginKey(login string, ip string) string {
	return ip + ":" + login
}

func IsUSerBlocked(login string, ip string) bool {
	result := false
	key := GetLoginKey(login, ip)
	mu.Lock()
	if blockedUser, ok := blockedLogins[key]; ok {
		//Block period is still active
		blockDuration := time.Minute
		if blockedUser.FailedAttempts > 3 {
			blockDuration = time.Hour
		}

		if time.Since(blockedUser.LastAttemptAt) < time.Duration(blockedUser.FailedAttempts*blockedUser.FailedAttempts)*blockDuration {
			result = true
		} else {
			//Block period has ended, release user+ip from prison
			delete(blockedLogins, key)
		}
	}
	mu.Unlock()
	return result
}

func RegisterFailedAuth(login string, ip string) {
	key := GetLoginKey(login, ip)
	mu.Lock()
	if blockedUser, ok := blockedLogins[key]; ok {
		blockedUser.FailedAttempts++
		blockedUser.LastAttemptAt = time.Now()
	} else {
		blockedLogins[key] = domain.BlockedLogin{
			IP:             ip,
			Login:          login,
			FailedAttempts: 1,
			LastAttemptAt:  time.Now(),
		}
	}
	mu.Unlock()
}
