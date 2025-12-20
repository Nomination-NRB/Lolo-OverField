package config

import (
	"encoding/json"
	"fmt"
	"time"

	"gucooing/lolo/db"
)

type DB struct {
	Dev             bool      `toml:"dev"`
	DbType          db.DbType `json:"dbType"`
	Dsn             string    `json:"dsn"`
	MaxIdleConns    int       `json:"maxIdleConns"`
	MaxOpenConns    int       `json:"maxOpenConns"`
	ConnMaxLifetime Duration  `json:"connMaxLifetime"`
}

var defaultDB = &DB{
	Dev:          false,
	DbType:       "sqlite",
	Dsn:          "./db/lolo.db",
	MaxIdleConns: 20,
	MaxOpenConns: 40,
	ConnMaxLifetime: Duration{
		time.Hour,
	},
}

func GetDB() *DB {
	if GetConfig().DB == nil {
		GetConfig().DB = defaultDB
	}
	return GetConfig().DB
}

func (x *DB) GetOption() *db.Option {
	return &db.Option{
		Dev:             x.Dev,
		Type:            x.DbType,
		Dsn:             x.Dsn,
		MaxIdleConns:    x.MaxIdleConns,
		MaxOpenConns:    x.MaxOpenConns,
		ConnMaxLifetime: x.ConnMaxLifetime.Duration,
	}
}

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value) * time.Second
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		return err
	default:
		return fmt.Errorf("无效的持续时间类型: %T", value)
	}
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}
