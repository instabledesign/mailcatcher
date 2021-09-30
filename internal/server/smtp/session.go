package smtp

import (
	"io"
	"io/ioutil"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/gol4ng/logger"
	"github.com/instabledesign/mailcatcher/internal"
)

// A Session is returned after successful login.
type Session struct {
	user        internal.User
	mail        *internal.Mail
	mailHandler internal.MailHandler
	logger      logger.LoggerInterface
}

func (s *Session) resetMail() {
	if s.mail == nil {
		s.mail = &internal.Mail{}
	}
	s.mail.Time = time.Now()
	s.mail.User = s.user
	s.mail.From = ""
	s.mail.Tos = []string{}
	s.mail.Data = ""
}

func (s *Session) Mail(from string, opts smtp.MailOptions) error {
	s.mail.From = from
	s.logger.Debug("mail from", logger.Any("from", s.mail.From))
	return nil
}

func (s *Session) Rcpt(to string) error {
	s.mail.Tos = append(s.mail.Tos, to)
	s.logger.Debug("mail to", logger.Any("to", to))
	return nil
}

func (s *Session) Data(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		s.logger.Error("an error occured during data reception", logger.Any("data", string(b)), logger.Error("err", err))
		s.resetMail()
		return err
	}
	s.mail.Data = string(b)
	s.logger.Debug("data", logger.Any("data", s.mail.Data))
	return nil
}

func (s *Session) Reset() {
	s.mailHandler.Handle(s.mail)
	s.logger.Info("push mail", logger.Any("mail", s.mail))
	s.resetMail()
	s.logger.Debug("reset data")
}

func (s *Session) Logout() error {
	s.logger.Debug("logout")
	return nil
}
