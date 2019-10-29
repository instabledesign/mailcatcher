package smtp

import (
	"io"
	"io/ioutil"

	"github.com/gol4ng/logger"
	"github.com/gol4ng/mailcatcher/internal"
)

// A Session is returned after successful login.
type Session struct {
	user        internal.User
	mail        *internal.Mail
	mailHandler internal.MailHandler
	logger      logger.LoggerInterface
}

func (s *Session) resetMail() {
	s.mail = &internal.Mail{User: s.user}
}

func (s *Session) Mail(from string) error {
	s.mail.From = from
	_ = s.logger.Debug("mail from", logger.Ctx("from", s.mail.From))
	return nil
}

func (s *Session) Rcpt(to string) error {
	s.mail.To = to
	_ = s.logger.Debug("mail to", logger.Ctx("to", s.mail.To))
	return nil
}

func (s *Session) Data(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		_ = s.logger.Error("an error occured during data reception", logger.Ctx("data", string(b)).Add("err", err))
		s.resetMail()
		return err
	}
	s.mail.Data = string(b)
	_ = s.logger.Debug("data", logger.Ctx("data", s.mail.Data))
	return nil
}

func (s *Session) Reset() {
	s.mailHandler.Handle(s.mail)
	_ = s.logger.Info("push mail", logger.Ctx("mail", s.mail))
	s.resetMail()
	_ = s.logger.Debug("reset data", nil)
}

func (s *Session) Logout() error {
	_ = s.logger.Debug("logout", nil)
	return nil
}
