// PingDB implementation on application layer
package app

import (
	"context"
	"errors"
	"fmt"
	"internal/adapters/logger"
	"internal/domain"
	"net/http"
)

func PingDBServer(ctx context.Context) *domain.HandlerStatus {
	hs := new(domain.HandlerStatus)

	userName := ctx.Value(domain.CtxUsername).(string)
	logger.Infof("User [%s] requested ping", userName)

	if Sc.StorageMode != domain.Database {
		hs.Message = "cannot ping DB connection: server does not run in Database mode"
		hs.Err = errors.New(hs.Message)
		hs.HTTPStatus = http.StatusBadRequest
		return hs
	}

	if err := Stor.Ping(); err != nil {
		hs.Message = fmt.Sprintf("error pinging DB server: %s", err)
		hs.Err = err
		hs.HTTPStatus = http.StatusInternalServerError
		return hs
	}

	hs.Message = "Ping OK"
	return hs
}
