package common

import (
	"fmt"
	"strings"

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

func NewPostgresDBSession(cfg PostgresSQLConfig) (db.Session, error) {
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

	return session, nil
}

func IsDBErrorNoRows(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "no more rows in this result set")
}
