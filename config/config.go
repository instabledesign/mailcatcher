package config

import (
	"context"
	"time"

	"github.com/gol4ng/logger"
	"github.com/instabledesign/mailcatcher/pkg/config"
)

type Config struct {
	Debug           bool               `config:"DEBUG,short=d"`
	LogLevel        logger.LevelString `config:"LOG_LEVEL,short=l"`
	LogVerboseLevel logger.LevelString `config:"LOG_VERBOSE_LEVEL,short=vl"`

	HTTPAddr            string        `config:"HTTP_ADDR"`
	HTTPShutdownTimeout time.Duration `config:"HTTP_SHUTDOWN_TIMEOUT"`

	SMTPAddr            string        `config:"SMTP_ADDR"`
	SMTPDomain          string        `config:"SMTP_DOMAIN"`
	SMTPReadTimeout     time.Duration `config:"SMTP_READ_TIMEOUT"`
	SMTPWriteTimeout    time.Duration `config:"SMTP_WRITE_TIMEOUT"`
	SMTPMaxMessageBytes int           `config:"SMTP_MAX_MESSAGE_BYTES"`
	SMTPMaxRecipients   int           `config:"SMTP_MAX_RECIPIENTS"`

	NotifBufferSize int `config:"NOTIF_BUFFER_SIZE"`
}

func NewConfig() *Config {
	cfg := &Config{
		Debug:           false,
		LogVerboseLevel: logger.LevelString(logger.DebugLevel.String()),
		LogLevel:        logger.LevelString(logger.InfoLevel.String()),

		HTTPAddr:            ":1080",
		HTTPShutdownTimeout: 3 * time.Second,

		SMTPAddr:            ":1025",
		SMTPDomain:          "localhost",
		SMTPReadTimeout:     20 * time.Second,
		SMTPWriteTimeout:    20 * time.Second,
		SMTPMaxMessageBytes: 1024 * 1024,
		SMTPMaxRecipients:   50,

		NotifBufferSize: 1000,
	}

	config.LoadOrFatal(context.Background(), cfg)
	return cfg
}
