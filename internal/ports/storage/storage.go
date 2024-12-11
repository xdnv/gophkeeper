package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"internal/adapters/logger"
	"internal/domain"
)

// universal storage
type UniStorage struct {
	config *domain.ServerConfig
	ctx    context.Context
	//stor    *MemStorage
	db      *DbStorage
	timeout time.Duration
}

// init storage
func NewUniStorage(cf *domain.ServerConfig) *UniStorage {

	var (
		conn *sql.DB
		err  error
	)

	if !cf.MockMode {
		conn, err = sql.Open("pgx", cf.DatabaseDSN)
		if err != nil {
			logger.Fatal(err.Error())
		}
	} else {
		conn = cf.MockConn
	}

	return &UniStorage{
		config:  cf,
		ctx:     context.Background(),
		db:      NewDbStorage(conn),
		timeout: 5 * time.Second,
	}
}

func (t UniStorage) Bootstrap() error {
	return t.db.Bootstrap(t.ctx)
}

func (t UniStorage) Close() {
	t.db.Close()
}

func (t UniStorage) Ping() error {
	dbctx, cancel := context.WithTimeout(t.ctx, t.timeout)
	defer cancel()

	errMsg := "UniStorage.Ping error"
	backoff := func(ctx context.Context) error {
		err := t.db.Ping(dbctx)
		return HandleRetriableDB(err, errMsg)
	}
	err := DoRetry(dbctx, t.config.MaxConnectionRetries, backoff)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", errMsg, err))
	}
	return err
}

// Checks if user exists without throwing an error. Returns bool status and password hash if succeded.
func (t UniStorage) IsUserExists(uuid string, isLogin bool) (bool, *domain.UserAccountRecord, error) {

	dbctx, cancel := context.WithTimeout(t.ctx, t.timeout)
	defer cancel()
	var ue bool
	var ur *domain.UserAccountRecord

	errMsg := "UniStorage.IsUserExists error"
	backoff := func(ctx context.Context) error {
		var err error
		ue, ur, err = t.db.IsUserExists(dbctx, uuid, isLogin)
		return HandleRetriableDB(err, errMsg)
	}
	err := DoRetry(dbctx, t.config.MaxConnectionRetries, backoff)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", errMsg, err))
		return false, nil, err
	}
	return ue, ur, err
}

// Retrieves existing User record
func (t UniStorage) GetUserRecord(uuid string, isLogin bool) (*domain.UserAccountRecord, error) {

	dbctx, cancel := context.WithTimeout(t.ctx, t.timeout)
	defer cancel()
	var ur *domain.UserAccountRecord

	errMsg := "UniStorage.GetUserRecord error"
	backoff := func(ctx context.Context) error {
		var err error
		ur, err = t.db.GetUserRecord(dbctx, uuid, isLogin)
		return HandleRetriableDB(err, errMsg)
	}
	err := DoRetry(dbctx, t.config.MaxConnectionRetries, backoff)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", errMsg, err))
		return nil, err
	}
	return ur, err
}

func (t UniStorage) UpdateUserRecord(ur *domain.UserAccountRecord) (string, error) {

	dbctx, cancel := context.WithTimeout(t.ctx, t.timeout)
	defer cancel()
	var result string

	errMsg := "UniStorage.UpdateUserRecord error"
	backoff := func(ctx context.Context) error {
		var err error
		result, err = t.db.UpdateUserRecord(dbctx, ur)
		return HandleRetriableDB(err, errMsg)
	}
	err := DoRetry(dbctx, t.config.MaxConnectionRetries, backoff)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", errMsg, err))
	}

	return result, nil
}

func (t UniStorage) GetSecrets(userId string) (*domain.KeeperRecords, error) {

	dbctx, cancel := context.WithTimeout(t.ctx, t.timeout)
	defer cancel()
	var kr *domain.KeeperRecords

	errMsg := "UniStorage.GetSecrets error"
	backoff := func(ctx context.Context) error {
		var err error
		kr, err = t.db.GetSecrets(dbctx, userId)
		return HandleRetriableDB(err, errMsg)
	}
	err := DoRetry(dbctx, t.config.MaxConnectionRetries, backoff)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", errMsg, err))
		return nil, err
	}
	return kr, err
}

func (t UniStorage) GetSecret(id string) (*domain.KeeperRecord, error) {

	dbctx, cancel := context.WithTimeout(t.ctx, t.timeout)
	defer cancel()
	var r *domain.KeeperRecord

	errMsg := "UniStorage.GetSecret error"
	backoff := func(ctx context.Context) error {
		var err error
		r, err = t.db.GetSecret(dbctx, id)
		return HandleRetriableDB(err, errMsg)
	}
	err := DoRetry(dbctx, t.config.MaxConnectionRetries, backoff)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", errMsg, err))
		return nil, err
	}
	return r, err
}

func (t UniStorage) AddSecret(r *domain.KeeperRecord) error {
	dbctx, cancel := context.WithTimeout(t.ctx, t.timeout)
	defer cancel()
	var id string

	errMsg := "UniStorage.AddSecret error"
	backoff := func(ctx context.Context) error {
		var err error
		id, err = t.db.AddSecret(dbctx, r)
		return HandleRetriableDB(err, errMsg)
	}
	err := DoRetry(dbctx, t.config.MaxConnectionRetries, backoff)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", errMsg, err))
		return err
	}
	r.ID = id
	return nil
}

func (t UniStorage) UpdateSecret(r *domain.KeeperRecord) error {
	dbctx, cancel := context.WithTimeout(t.ctx, t.timeout)
	defer cancel()

	errMsg := "UniStorage.UpdateSecret error"
	backoff := func(ctx context.Context) error {
		_, err := t.db.UpdateSecret(dbctx, r)
		return HandleRetriableDB(err, errMsg)
	}
	err := DoRetry(dbctx, t.config.MaxConnectionRetries, backoff)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", errMsg, err))
		return err
	}
	return nil
}

func (t UniStorage) DeleteSecret(id string) error {
	dbctx, cancel := context.WithTimeout(t.ctx, t.timeout)
	defer cancel()

	errMsg := "UniStorage.DeleteSecret error"
	backoff := func(ctx context.Context) error {
		err := t.db.DeleteSecret(dbctx, id)
		return HandleRetriableDB(err, errMsg)
	}
	err := DoRetry(dbctx, t.config.MaxConnectionRetries, backoff)
	if err != nil {
		logger.Error(fmt.Sprintf("%s: %s\n", errMsg, err))
		return err
	}
	return nil
}
