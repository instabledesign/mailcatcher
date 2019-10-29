package service

import (
	"github.com/gol4ng/mailcatcher/internal/correlation_id"
	smtpSrv "github.com/gol4ng/mailcatcher/internal/server/smtp"
)

func (container *App) GetSMTPServer() *smtpSrv.Server {
	return smtpSrv.NewServer(
		container.Cfg,
		container.getSMTPBackend(),
	)
}

func (container *App) getSMTPBackend() *smtpSrv.Backend {
	return smtpSrv.NewBackend(
		correlation_id.DefaultIdGenerator,
		container.GetMailHandler(),
		container.GetUserLogger,
	)
}
