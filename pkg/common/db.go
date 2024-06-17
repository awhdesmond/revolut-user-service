package common

import (
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/upper/db/v4"
	postgresqladp "github.com/upper/db/v4/adapter/postgresql"
)

type PostgresSQLConfig struct {
	Host     string `mapstructure:"postgres-host"`
	Port     string `mapstructure:"postgres-port"`
	Database string `mapstructure:"postgres-database"`
	Username string `mapstructure:"postgres-username"`
	Password string `mapstructure:"postgres-password"`
}

func (c PostgresSQLConfig) Hostname() string {
	if c.Port == "" {
		return c.Host
	}
	return fmt.Sprintf("%s:%v", c.Host, c.Port)
}

func MakePostgresDBSession(cfg PostgresSQLConfig) (db.Session, error) {
	settings := postgresqladp.ConnectionURL{
		User:     cfg.Username,
		Password: cfg.Password,
		Host:     cfg.Hostname(),
		Database: cfg.Database,
	}

	session, err := postgresqladp.Open(settings)
	if err != nil {
		return nil, err
	}

	db.LC().SetLevel(db.LogLevelError)
	session.SetMaxIdleConns(db.DefaultSettings.MaxIdleConns())
	session.SetMaxOpenConns(db.DefaultSettings.MaxOpenConns())
	session.SetConnMaxLifetime(db.DefaultSettings.ConnMaxLifetime())
	session.SetConnMaxIdleTime(db.DefaultSettings.ConnMaxIdleTime())

	return session, nil
}

func IsDBErrorNoRows(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "no more rows in this result set")
}

// Redis
type RedisCfg struct {
	URI      string `mapstructure:"redis-uri"`
	Password string `mapstructure:"redis-password"`
}

func redisOptsToUnivOpts(opts *redis.Options, password string) *redis.UniversalOptions {
	return &redis.UniversalOptions{
		Addrs:    []string{opts.Addr},
		Password: password,
		DB:       opts.DB,
	}
}

func MakeRedisClient(cfg RedisCfg) (redis.UniversalClient, error) {
	opts, err := redis.ParseURL(cfg.URI)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewUniversalClient(redisOptsToUnivOpts(opts, cfg.Password))
	return rdb, nil
}
