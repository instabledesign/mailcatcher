package service

import (
	"context"

	"github.com/gol4ng/logger"
	"github.com/gol4ng/mailcatcher/config"
	"github.com/gol4ng/mailcatcher/internal"
)

type DeamonStatus int
type Deamon interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Statut(ctx context.Context) DeamonStatus
}

type App struct {
	Cfg *config.Config

	mailHandler      internal.MailHandler
	mailChan         chan *internal.Mail
	mailChanNotifier *internal.MailChanNotifier
	mailStorage      *internal.MailStorage

	logger            logger.LoggerInterface
	loggerMiddlewares logger.Middlewares

	deamons []Deamon
}

func (container *App) GetMailHandler() internal.MailHandler {
	if container.mailHandler == nil {
		container.mailHandler = internal.NewMailHandlerComposite(
			container.GetMailChanNotifier(),
			container.GetMailStorage(),
		)
	}
	return container.mailHandler
}

func (container *App) GetMailChan() chan *internal.Mail {
	return container.mailChan
}

func (container *App) GetMailChanNotifier() *internal.MailChanNotifier {
	if container.mailChanNotifier == nil {
		container.mailChanNotifier = internal.NewMailChanNotifier(container.GetMailChan())
	}
	return container.mailChanNotifier
}

func (container *App) GetMailStorage() *internal.MailStorage {
	if container.mailStorage == nil {
		container.mailStorage = &internal.MailStorage{}
	}
	return container.mailStorage
}


func NewContainer(cfg *config.Config) *App {
	return &App{
		Cfg:      cfg,
		mailChan: make(chan *internal.Mail, cfg.NotifBufferSize),
	}
}
