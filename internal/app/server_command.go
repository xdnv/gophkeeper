// PingDB implementation on application layer
package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"internal/adapters/logger"
	"internal/domain"
	"io"
	"net/http"
)

func ExecuteCommand(ctx context.Context, data io.Reader) (*[]byte, *domain.HandlerStatus) {
	var hs domain.HandlerStatus
	var rc domain.RemoteCommand
	var response interface{}
	var err error

	userName := ctx.Value(domain.CtxUsername).(string)
	logger.Infof("User [%s] requested command execution", userName)

	// convert username to user ID
	ur, err := Stor.GetUserRecord(userName, true)
	if err != nil {
		hs.Message = fmt.Sprintf("user account error for login [%s]: %s", userName, err.Error())
		hs.Err = err
		hs.HTTPStatus = http.StatusBadRequest
		return nil, &hs
	}
	userID := ur.ID

	if err := json.NewDecoder(data).Decode(&rc); err != nil {
		hs.Message = fmt.Sprintf("json command decode error: %s", err.Error())
		hs.Err = err
		hs.HTTPStatus = http.StatusBadRequest
		return nil, &hs
	}

	switch rc.Command {
	case domain.S_CMD_SYNC:
		response, err = Stor.GetSecrets(userID)
		if err != nil {
			hs.Message = fmt.Sprintf("command execution error [%s]: %s", rc.Command, err.Error())
			hs.Err = err
			hs.HTTPStatus = http.StatusInternalServerError
			return nil, &hs
		}

	case domain.S_CMD_DELETE:
		if len(rc.Arguments) < 1 {
			hs.Message = fmt.Sprintf("command execution error [%s]: no record ID provided", rc.Command)
			hs.Err = errors.New(hs.Message)
			hs.HTTPStatus = http.StatusInternalServerError
			return nil, &hs
		}
		recordID := rc.Arguments[0]
		err = Stor.DeleteSecret(recordID, userID)
		if err != nil {
			hs.Message = fmt.Sprintf("command execution error [%s]: %s", rc.Command, err.Error())
			hs.Err = err
			hs.HTTPStatus = http.StatusInternalServerError
			return nil, &hs
		}
		//sync after delete
		response, err = Stor.GetSecrets(userID)
		if err != nil {
			hs.Message = fmt.Sprintf("command execution error [%s]: %s", rc.Command, err.Error())
			hs.Err = err
			hs.HTTPStatus = http.StatusInternalServerError
			return nil, &hs
		}

	case domain.S_CMD_NEW:
		if rc.Data == nil {
			hs.Message = fmt.Sprintf("command execution error [%s]: empty Data filed", rc.Command)
			hs.Err = errors.New(hs.Message)
			hs.HTTPStatus = http.StatusInternalServerError
			return nil, &hs
		}

		var kr domain.KeeperRecord
		buf := bytes.NewBuffer(rc.Data)
		if err := json.NewDecoder(buf).Decode(&kr); err != nil {
			hs.Message = fmt.Sprintf("json command data record decode error: %s", err.Error())
			hs.Err = err
			hs.HTTPStatus = http.StatusBadRequest
			return nil, &hs
		}
		kr.UserID = userID

		err = Stor.AddSecret(&kr)
		if err != nil {
			hs.Message = fmt.Sprintf("command execution error [%s]: %s", rc.Command, err.Error())
			hs.Err = err
			hs.HTTPStatus = http.StatusInternalServerError
			return nil, &hs
		}
		//no JSON response here

	default:
		hs.Message = fmt.Sprintf("unknown command: %s", rc.Command)
		hs.Err = fmt.Errorf("unknown command: %s", rc.Command)
		hs.HTTPStatus = http.StatusBadRequest
		return nil, &hs
	}

	// processing command response (should be valid JSON)
	resp, err := json.Marshal(response)
	if err != nil {
		hs.Message = fmt.Sprintf("json command response encode error: %s", err.Error())
		hs.Err = err
		hs.HTTPStatus = http.StatusInternalServerError
		return nil, &hs
	}

	hs.Message = "OK"
	return &resp, &hs
}
