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
func RegisterNewUser(data io.Reader) (*[]byte, *domain.HandlerStatus) {
	hs := new(domain.HandlerStatus)

	var u domain.UserAccountData
	var ur domain.UserAccountRecord

	//logger.Debug(fmt.Sprintf("RegisterNewUser body: %v", data)) //DEBUG

	if err := json.NewDecoder(data).Decode(&u); err != nil {
		hs.Message = fmt.Sprintf("json login decode error: %s", err.Error())
		hs.Err = err
		hs.HTTPStatus = http.StatusBadRequest
		return nil, hs
	}

	err := ValidateAccountData(&u)
	if err != nil {
		hs.Message = fmt.Sprintf("error checking user account data: %s", err.Error())
		hs.Err = err
		hs.HTTPStatus = http.StatusBadRequest
		return nil, hs
	}

	exists, _, err := Stor.IsUserExists(u.Login, true)
	if err != nil {
		hs.Message = fmt.Sprintf("error checking user record [%s]: %s", u.Login, err.Error())
		hs.Err = err
		hs.HTTPStatus = http.StatusNotFound
		return nil, hs
	}

	if exists {
		hs.Message = "User already exists"
		hs.Err = fmt.Errorf(hs.Message)
		hs.HTTPStatus = http.StatusConflict
		return nil, hs
	}

	ur.Login = u.Login
	ur.Password, err = HashPassword(u.Password)
	ur.Email = u.Email
	if err != nil {
		hs.Message = fmt.Sprintf("error processing password: %s", err.Error())
		hs.Err = err
		hs.HTTPStatus = http.StatusNotFound
		return nil, hs
	}

	uid, err := Stor.UpdateUserRecord(&ur)
	if err != nil {
		hs.Message = fmt.Sprintf("error creating user record [%s]: %s", u.Login, err.Error())
		hs.Err = err
		hs.HTTPStatus = http.StatusNotFound
		return nil, hs
	}

	logger.Debugf("RegisterNewUser, created [%s][%s][%s]", ur.Login, ur.Email, uid) //DEBUG

	// resp, err := json.Marshal(m)
	// if err != nil {
	// 	hs.Message = fmt.Sprintf("json metric encode error: %s", err.Error())
	// 	hs.Err = err
	// 	hs.HTTPStatus = http.StatusInternalServerError
	// 	return nil, hs
	// }

	// return &resp, hs

	return nil, hs
}
