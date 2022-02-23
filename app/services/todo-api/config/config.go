package config

import (
	"github.com/ardanlabs/conf/v3"
	"time"
)

type WebConfig struct {
	ReadTimeout     time.Duration `conf:"default:5s"`
	WriteTimeout    time.Duration `conf:"default:10s"`
	IdleTimeout     time.Duration `conf:"default:120s"`
	ShutdownTimeout time.Duration `conf:"default:20s"`
	APIHost         string        `conf:"default:0.0.0.0:3000"`
	DebugHost       string        `conf:"default:0.0.0.0:4000"`
}

type AuthConfig struct {
	KeysFolder string `conf:"default:zarf/keys/"`
	// DO NOT change this to default to "dev"
	AuthMode string `conf:"default:prod"`
}

type Config struct {
	conf.Version
	Web  WebConfig
	Auth AuthConfig
	DB   struct {
		User         string `conf:"default:postgres"`
		Password     string `conf:"default:postgres,mask"`
		Host         string `conf:"default:localhost"`
		Name         string `conf:"default:postgres"`
		MaxIdleConns int    `conf:"default:0"`
		MaxOpenConns int    `conf:"default:0"`
		DisableTLS   bool   `conf:"default:true"`
	}
}
