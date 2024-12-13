package storage

import (
	"context"
	"database/sql"
	"fmt"
	"internal/adapters/logger"
	"internal/domain"
	"strings"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Main storage
type DbStorage struct {
	conn *sql.DB
}

// Init DB storage object
func NewDbStorage(conn *sql.DB) *DbStorage {
	return &DbStorage{conn: conn}
}

// Closes db connection
func (t DbStorage) Close() {
	t.conn.Close()
}

// prepare database
func (t DbStorage) Bootstrap(ctx context.Context) error {

	// begin transaction
	tx, err := t.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	logger.Info("BOOTSTRAP STARTED")

	//check config
	//tableName := "public.config"
	dbKey := "DBVersion"
	dbAppName := "GophKeeper"
	dbVersion := "20241201"

	//Important! pgx does not support sql.Named(), use pgx.NamedArgs{} instead

	// enabling support of crypto operations
	// usage example: pgp_sym_encrypt(data::text, 'secret_password') + select pgp_sym_decrypt(data, 'secret_password')
	logger.Info("init pgcrypto extension")
	_, err = tx.ExecContext(ctx, `
		create extension if not exists pgcrypto;
	`) //,
	if err != nil {
		return err
	}

	// config stores application-wide metadata
	//TODO: add version update procedure
	logger.Info("init config")
	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS public.config (
			key VARCHAR(128) NOT NULL PRIMARY KEY,
			value TEXT
		);
	`) //,
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO public.config (key, value)
			VALUES (@dbKey::text, @dbVersion::text)
		ON CONFLICT (key)
			DO UPDATE SET value = excluded.value;
	`,
		pgx.NamedArgs{
			"dbKey":     dbKey,
			"dbAppName": dbAppName,
			"dbVersion": dbVersion,
		},
	)
	if err != nil {
		return err
	}

	// users
	logger.Info("init users")
	_, err = tx.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS public.users (
			id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
            username VARCHAR(50) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			registration_date TIMESTAMP DEFAULT NOW(),
			is_banned BOOL NOT NULL
        );
    `)
	if err != nil {
		return err
	}
	//TODO: add optional user details used to recover lost access etc. May be subject to GDPR.
	//full_name VARCHAR(255),
	//mobile_number VARCHAR(255),

	// completely secret public secrets
	logger.Info("init secrets")
	_, err = tx.ExecContext(ctx, `
	    CREATE TABLE IF NOT EXISTS public.secrets (
			id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL,
			name VARCHAR(50) NOT NULL,
			description VARCHAR(1024) NOT NULL,
			secret_type VARCHAR(50) NOT NULL,
			created TIMESTAMP DEFAULT NOW(),
			modified TIMESTAMP DEFAULT NOW(),
			is_deleted BOOL NOT NULL,
			secret BYTEA,
			CONSTRAINT secret_fk_user FOREIGN KEY (user_id) REFERENCES public.users(id),
			CONSTRAINT secret_uniq_user UNIQUE (id, user_id)
	    );
	`)
	if err != nil {
		return err
	}

	logger.Info("BOOTSTRAP OK")

	// commit transaction
	return tx.Commit()
}

func (t DbStorage) Ping(ctx context.Context) error {
	return t.conn.PingContext(ctx)
}

func (t DbStorage) IsUserExists(ctx context.Context, uid string, isLogin bool) (bool, *domain.UserAccountRecord, error) {

	var u domain.UserAccountRecord

	query := `
		SELECT
			t.id,
			t.username,
			t.password,
			t.email,
			t.registration_date,
			t.is_banned
		FROM public.users t
		WHERE t.id = @uid;
	`
	//alter query condition to support search by login
	if isLogin {
		query = strings.Replace(query, "WHERE t.id = @uid;", "WHERE t.username = @uid;", 1)
	}

	err := t.conn.QueryRowContext(ctx, query,
		pgx.NamedArgs{
			"uid": uid,
		},
	).Scan(&u.ID, &u.Login, &u.Password, &u.Email, &u.RegistrationDate, &u.IsBanned)

	switch {
	case err == sql.ErrNoRows:
		return false, nil, nil
	case err != nil:
		logger.Errorf("IsUserExists: %s", err)
		return false, nil, err
	default:
		return true, &u, nil
	}
}

func (t DbStorage) GetUserRecord(ctx context.Context, uid string, isLogin bool) (*domain.UserAccountRecord, error) {

	var u domain.UserAccountRecord

	query := `
		SELECT
			t.id,
			t.username,
			t.password,
			t.email,
			t.registration_date,
			t.is_banned
		FROM public.users t
		WHERE t.id = @uid;
	`
	//alter query condition to support search by login
	if isLogin {
		query = strings.Replace(query, "WHERE t.id = @uid;", "WHERE t.username = @uid;", 1)
	}

	err := t.conn.QueryRowContext(ctx, query,
		pgx.NamedArgs{
			"uid": uid,
		},
	).Scan(&u.ID, &u.Login, &u.Password, &u.Email, &u.RegistrationDate, &u.IsBanned)

	if err != nil {
		logger.Errorf("GetUserRecord: %s", err)
		return nil, err
	}

	return &u, nil
}

func (t DbStorage) UpdateUserRecord(ctx context.Context, ur *domain.UserAccountRecord) (string, error) {

	var result string

	// begin transaction
	tx, err := t.conn.BeginTx(ctx, nil)
	if err != nil {
		return result, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO public.users (username, password, email, is_banned)
			VALUES (@username, @password, @email, @is_banned)
		ON CONFLICT (username) DO UPDATE
			SET
				password = EXCLUDED.password,
				email = EXCLUDED.email,
				is_banned = EXCLUDED.is_banned
		RETURNING id;
	`

	err = tx.QueryRowContext(ctx, query,
		pgx.NamedArgs{
			"username":  ur.Login,
			"password":  ur.Password,
			"email":     ur.Email,
			"is_banned": ur.IsBanned,
		},
	).Scan(&result)
	if err != nil {
		return result, err
	}

	ur.ID = result

	// commit transaction
	err = tx.Commit()
	return result, err
}

func (t DbStorage) AddSecret(ctx context.Context, r *domain.KeeperRecord) (string, error) {

	var id string

	// begin transaction
	tx, err := t.conn.BeginTx(ctx, nil)
	if err != nil {
		return id, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO public.secrets (user_id, name, description, secret_type, is_deleted, secret)
			VALUES (@user_id::uuid, @name, @description, @secret_type, @is_deleted, pgp_sym_encrypt(@secret::text, @user_id::text))
		RETURNING id;
	`
	err = tx.QueryRowContext(ctx, query,
		pgx.NamedArgs{
			"user_id":     r.UserID,
			"name":        r.Name,
			"description": r.Description,
			"secret_type": r.SecretType,
			"is_deleted":  r.IsDeleted,
			"secret":      r.Secret,
		},
	).Scan(&id)
	if err != nil {
		logger.Errorf("InsertSecret: %s", err)
		return id, err
	}

	// commit transaction
	err = tx.Commit()

	return id, err
}

func (t DbStorage) UpdateSecret(ctx context.Context, r *domain.KeeperRecord) (string, error) {

	var id string

	// begin transaction
	tx, err := t.conn.BeginTx(ctx, nil)
	if err != nil {
		return id, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO public.secrets (id, user_id, name, description, secret_type, is_deleted, secret)
			VALUES (@id, @user_id::uuid, @name, @description, @secret_type, @is_deleted, pgp_sym_encrypt(@secret::text, @user_id::text))
		ON CONFLICT ON CONSTRAINT secret_uniq_user DO UPDATE
			SET
				name = EXCLUDED.name,
				description = EXCLUDED.description,
				secret_type = EXCLUDED.secret_type,
				modified = NOW(),
				is_deleted = EXCLUDED.is_deleted,
				secret = EXCLUDED.secret
		RETURNING id;
	`
	err = tx.QueryRowContext(ctx, query,
		pgx.NamedArgs{
			"id":          r.ID,
			"user_id":     r.UserID,
			"name":        r.Name,
			"description": r.Description,
			"secret_type": r.SecretType,
			"is_deleted":  r.IsDeleted,
			"secret":      r.Secret,
		},
	).Scan(&id)
	if err != nil {
		logger.Errorf("UpdateSecret: %s", err)
		return id, err
	}

	// commit transaction
	err = tx.Commit()

	return id, err
}

func (t DbStorage) GetSecret(ctx context.Context, id string) (*domain.KeeperRecord, error) {

	var r domain.KeeperRecord

	query := `
		SELECT
			t.id,
			t.user_id,
			t.name,
			t.description,
			t.secret_type,
			t.is_deleted,
			pgp_sym_decrypt(secret, user_id) AS secret
		FROM public.secrets t
		WHERE t.id = @id;
	`
	err := t.conn.QueryRowContext(ctx, query,
		pgx.NamedArgs{
			"id": id,
		},
	).Scan(&r.ID, &r.UserID, &r.Name, &r.Description, &r.SecretType, &r.IsDeleted, &r.Secret)

	if err != nil {
		logger.Errorf("GetSecret: %s", err)
		return nil, err
	}

	return &r, nil
}

func (t DbStorage) GetSecrets(ctx context.Context, userId string) (*domain.KeeperRecords, error) {

	var r domain.KeeperRecords

	query := `
		SELECT id, user_id, name, description, secret_type, is_deleted, pgp_sym_decrypt(secret, user_id::text)
		FROM public.secrets
		WHERE user_id = @user_id AND NOT is_deleted;
	`

	rows, err := t.conn.QueryContext(ctx, query,
		pgx.NamedArgs{
			"user_id": userId,
		},
	)
	if err != nil {
		logger.Errorf("GetSecrets: %s", err)
		return nil, err
	}
	defer rows.Close()

	var counter = 1
	for rows.Next() {
		kr := new(domain.KeeperRecord)
		if err := rows.Scan(&kr.ID, &kr.UserID, &kr.Name, &kr.Description, &kr.SecretType, &kr.IsDeleted, &kr.Secret); err != nil {
			logger.Errorf("GetSecrets: %s", err)
			return nil, err
		}

		kr.ListNR = fmt.Sprintf("%v", counter)
		kr.ShortID = kr.ID[:8]

		r = append(r, *kr)

		counter++
	}

	//logger.Infof("GetSecrets: got rows %s", len(r)) //DEBUG

	return &r, nil
}

// Delete secret by ID for specified userID (security measure)
func (t DbStorage) DeleteSecret(ctx context.Context, id string, userId string) error {

	// begin transaction
	tx, err := t.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE public.secrets
		SET is_deleted = TRUE
		WHERE id = @id AND user_id = @user_id;
	`

	_, err = tx.ExecContext(ctx, query,
		pgx.NamedArgs{
			"id":      id,
			"user_id": userId,
		},
	)
	if err != nil {
		return err
	}

	// commit transaction
	err = tx.Commit()
	return err
}
