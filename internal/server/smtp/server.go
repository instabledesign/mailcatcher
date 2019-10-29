package smtp

import (
	"context"
	"errors"
	"net"

	"github.com/emersion/go-smtp"
	"github.com/gol4ng/logger"
	"github.com/gol4ng/mailcatcher/config"
	"github.com/gol4ng/mailcatcher/internal"
	"github.com/gol4ng/mailcatcher/internal/correlation_id"
)

var ErrServerClosed = errors.New("smtp: Server closed")

// The Server implements SMTP server methods.
type Server struct {
	config     *config.Config
	backend    *Backend
	smtpServer *smtp.Server
	logger     logger.LoggerInterface

	idGenerator       *correlation_id.RandomIdGenerator
	mailHandler       internal.MailHandler
	userLoggerFactory func(string, *logger.Context) logger.LoggerInterface
}

func (s *Server) MustStart(ctx context.Context) {
	go func() {
		if err := s.Start(ctx); err != ErrServerClosed {
			_ = s.logger.Critical("smtp server error: %error%", logger.Ctx("error", err))
			panic(err)
		}
	}()
}

func (s *Server) Start(_ context.Context) error {
	_ = s.logger.Debug("starting smtp server...", nil)
	s.smtpServer = smtp.NewServer(s.backend)

	s.smtpServer.Addr = s.config.SMTPAddr
	s.smtpServer.Domain = s.config.SMTPDomain
	s.smtpServer.ReadTimeout = s.config.SMTPReadTimeout
	s.smtpServer.WriteTimeout = s.config.SMTPWriteTimeout
	s.smtpServer.MaxMessageBytes = s.config.SMTPMaxMessageBytes
	s.smtpServer.MaxRecipients = s.config.SMTPMaxRecipients
	s.smtpServer.AllowInsecureAuth = true

	err := s.smtpServer.ListenAndServe()

	if e, ok := err.(*net.OpError); ok && e.Err.Error() == "use of closed network connection" {
		return ErrServerClosed
	}
	return err
}

func (s *Server) Stop(_ context.Context) error {
	_ = s.logger.Debug("stopping smtp server...", nil)
	s.smtpServer.Close()
	return nil
}

func NewServer(cfg *config.Config, backend *Backend) *Server {
	return &Server{
		config:  cfg,
		backend: backend,
		logger:  logger.NewNopLogger(),
	}
}
