package users

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/awhdesmond/revolut-user-service/pkg/common"
	"github.com/redis/go-redis/v9"
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

	DefaultCacheTTL = 10 * time.Minute
)

type Store interface {
	Upsert(ctx context.Context, username string, dob time.Time) error
	Read(ctx context.Context, username string) (User, error)
}

type store struct {
	sess   db.Session
	rdb    redis.UniversalClient
	logger *zap.Logger
}

func NewStore(sess db.Session, rdb redis.UniversalClient, logger *zap.Logger) Store {
	return &store{sess, rdb, logger.Named(loggerName)}
}

func (store *store) rdbUserKey(username string) string {
	return fmt.Sprintf("revolut_user_service:username:%s", username)
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

	cmd := store.rdb.Del(ctx, store.rdbUserKey(username))
	if cmd.Err() != nil {
		store.logger.Error("cache error", zap.Error(err))
		return ErrUnexpectedDatabaseError
	}

	return nil
}

func (store *store) Read(ctx context.Context, username string) (User, error) {
	var usr User

	rdbUserKey := store.rdbUserKey(username)
	cmd := store.rdb.Get(ctx, rdbUserKey)
	if cmd.Err() != nil {
		// unexpected error
		if !errors.Is(cmd.Err(), redis.Nil) {
			return User{}, ErrUnexpectedDatabaseError
		}
		// else is key not found, so just fallthrough
	} else {
		// Key is found, return from cache
		if err := json.Unmarshal([]byte(cmd.Val()), &usr); err != nil {
			// dirty data in cache, refetch from db
			store.rdb.Del(ctx, rdbUserKey)
		} else {
			return usr, nil
		}
	}

	// Key is not found in cache, fetch from db
	q := store.sess.WithContext(ctx).SQL().SelectFrom(dbtable).Where("username = ?", username)
	err := q.One(&usr)

	if common.IsDBErrorNoRows(err) {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		store.logger.Error("db error", zap.Error(err))
		return User{}, ErrUnexpectedDatabaseError
	}

	go func() {
		// asynchronously save to cache
		data, err := json.Marshal(usr)
		if err != nil {
			store.logger.Error("redis marshal error", zap.Error(err))
			return
		}
		cmd := store.rdb.Set(context.Background(), rdbUserKey, string(data), DefaultCacheTTL)
		if cmd.Err() != nil {
			store.logger.Error("redis cache error", zap.Error(err))
		}
	}()
	return usr, nil
}
