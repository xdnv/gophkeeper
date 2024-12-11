// RegisterNewUser implementation on application layer
package app

import (
	"encoding/json"
	"fmt"
	"internal/adapters/logger"
	"internal/domain"
	"io"
	"net/http"
)

// HTTP registration processing
func LoginUser(data io.Reader, metadata *domain.UserLoginMetadata) (*[]byte, *domain.HandlerStatus) {
	hs := new(domain.HandlerStatus)

	var u domain.UserAccountData

	//logger.Debugf("LoginUser body: %v", data) //DEBUG

	if err := json.NewDecoder(data).Decode(&u); err != nil {
		hs.Message = fmt.Sprintf("json login decode error: %s", err.Error())
		hs.Err = err
		hs.HTTPStatus = http.StatusBadRequest
		return nil, hs
	}

	err := ValidateAccountData(&u)
	if err != nil {
		hs.Message = fmt.Sprintf("error checking login data: %s", err.Error())
		hs.Err = err
		hs.HTTPStatus = http.StatusBadRequest
		return nil, hs
	}

	ip := metadata.IP
	logger.Infof("Auth attempt: [%s][%s]", u.Login, ip)

	if IsUSerBlocked(u.Login, ip) {
		hs.Message = "login rate limit excceded, please try later"
		hs.Err = fmt.Errorf(hs.Message)
		hs.HTTPStatus = http.StatusTooManyRequests
		return nil, hs
	}

	exists, ur, err := Stor.IsUserExists(u.Login, true)
	if err != nil {
		hs.Message = fmt.Sprintf("error checking user record [%s]: %s", u.Login, err.Error())
		hs.Err = err
		hs.HTTPStatus = http.StatusNotFound
		return nil, hs
	}

	if !exists || !CheckPasswordHash(u.Password, ur.Password) {
		RegisterFailedAuth(u.Login, ip)

		hs.Message = "wrong username or password"
		hs.Err = fmt.Errorf(hs.Message)
		hs.HTTPStatus = http.StatusUnauthorized
		return nil, hs
	}

	if ur.IsBanned {
		hs.Message = "this account is suspended"
		hs.Err = fmt.Errorf(hs.Message)
		hs.HTTPStatus = http.StatusUnauthorized
		return nil, hs
	}

	// generate AuthResponse structure with JWT token and crypto key
	tokenString, err := GenerateJWT(u.Login, 24)
	if err != nil {
		hs.Message = fmt.Sprintf("error generating JWT token [%s]: %s", u.Login, err.Error())
		hs.Err = err
		hs.HTTPStatus = http.StatusInternalServerError
		return nil, hs
	}

	response := domain.AuthResponse{
		Token:     tokenString,
		PublicKey: Sc.SessionCryptoKey.PublicKey,
	}

	//in case of empty response
	//return nil, hs

	// serialize response
	resp, err := json.Marshal(response)
	if err != nil {
		hs.Message = fmt.Sprintf("json auth response encode error: %s", err.Error())
		hs.Err = err
		hs.HTTPStatus = http.StatusInternalServerError
		return nil, hs
	}

	return &resp, hs
}
