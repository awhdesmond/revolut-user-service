package users

import (
	"context"
	"errors"
	"time"

	"github.com/awhdesmond/revolut-user-service/pkg/common"
	"github.com/upper/db/v4"
	"go.uber.org/zap"
)

const (
	dbtable    = "users"
	loggerName = "users.store"
)

var (
	ErrUserNotFound            = errors.New("username not found")
	ErrUnexpectedDatabaseError = errors.New("unexpected error")
)

type Store interface {
	Upsert(ctx context.Context, username string, dob time.Time) error
	Read(ctx context.Context, username string) (User, error)
}

type store struct {
	sess   db.Session
	logger *zap.Logger
}

func NewStore(sess db.Session, logger *zap.Logger) Store {
	return &store{sess, logger.Named(loggerName)}
}

func (store *store) Upsert(ctx context.Context, username string, dob time.Time) error {
	_, err := store.sess.WithContext(ctx).SQL().Query(`
		INSERT INTO users (username, date_of_birth)
		VALUES (?, ?)
		ON CONFLICT(username)
		DO UPDATE SET
			date_of_birth = EXCLUDED.date_of_birth
	`, username, dob)
	if err != nil {
		store.logger.Error("db error", zap.Error(err))
		return ErrUnexpectedDatabaseError
	}
	return nil
}

func (store *store) Read(ctx context.Context, username string) (User, error) {
	var usr User
	q := store.sess.WithContext(ctx).SQL().SelectFrom(dbtable).Where("username = ?", username)
	err := q.One(usr)

	if common.IsDBErrorNoRows(err) {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		store.logger.Error("db error", zap.Error(err))
		return User{}, ErrUnexpectedDatabaseError
	}
	return usr, nil
}
