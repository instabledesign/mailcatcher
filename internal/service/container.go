package service

import (
	"context"
	"github.com/emersion/go-smtp"
	"github.com/gol4ng/logger"
	"github.com/instabledesign/mailcatcher/config"
	"github.com/instabledesign/mailcatcher/internal"
	"github.com/instabledesign/mailcatcher/internal/server/http/handler"
	stop_dispatcher "github.com/gol4ng/stop-dispatcher"
	"net/http"
)

type Container struct {
	Cfg *config.Config

	baseContext context.Context

	mailHandler      internal.MailHandler
	mailChan         chan *internal.Mail
	mailChanNotifier *internal.MailChanNotifier
	mailStorage      *internal.MailStorage

	broker     *handler.Broker
	httpServer *http.Server
	smtpServer *smtp.Server

	logger            logger.LoggerInterface
	loggerMiddlewares logger.Middlewares

	stopDispatcher *stop_dispatcher.Dispatcher
}

func NewContainer(cfg *config.Config, ctx context.Context) *Container {
	return &Container{
		Cfg:         cfg,
		baseContext: ctx,
		mailChan:    make(chan *internal.Mail, cfg.NotifBufferSize),
	}
}
