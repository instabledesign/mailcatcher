package service

import (
	"sync"

	"github.com/emersion/go-smtp"
	"github.com/instabledesign/mailcatcher/internal/correlation_id"
	smtpSrv "github.com/instabledesign/mailcatcher/internal/server/smtp"
)

var smtpServerOnce sync.Once

func (container *Container) GetSMTPServer() *smtp.Server {
	smtpServerOnce.Do(func() {
		smtpServer := smtp.NewServer(container.getSMTPBackend())

		smtpServer.Addr = container.Cfg.SMTPAddr
		smtpServer.Domain = container.Cfg.SMTPDomain
		smtpServer.ReadTimeout = container.Cfg.SMTPReadTimeout
		smtpServer.WriteTimeout = container.Cfg.SMTPWriteTimeout
		smtpServer.MaxMessageBytes = container.Cfg.SMTPMaxMessageBytes
		smtpServer.MaxRecipients = container.Cfg.SMTPMaxRecipients
		smtpServer.AllowInsecureAuth = true

		container.smtpServer = smtpServer
	})

	return container.smtpServer
}

func (container *Container) getSMTPBackend() *smtpSrv.Backend {
	return smtpSrv.NewBackend(
		correlation_id.DefaultIdGenerator,
		container.GetMailHandler(),
		container.GetUserLogger,
	)
}
