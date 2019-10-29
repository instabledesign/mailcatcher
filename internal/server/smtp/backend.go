package smtp

import (
	"fmt"

	"github.com/emersion/go-smtp"
	"github.com/gol4ng/logger"
	"github.com/gol4ng/mailcatcher/internal"
	"github.com/gol4ng/mailcatcher/internal/correlation_id"
)

// The Server implements SMTP server methods.
type Backend struct {
	idGenerator       *correlation_id.RandomIdGenerator
	mailHandler       internal.MailHandler
	userLoggerFactory func(string, *logger.Context) logger.LoggerInterface
}

// Login handles a login command with username and password.
func (s *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	logR := s.userLoggerFactory(
		fmt.Sprintf("[%s@%s] ", username, state.RemoteAddr),
		logger.Ctx("remote_addr", state.RemoteAddr).Add("username", username),
	)
	//_ = logR.Debug(`login attempt`, logger.Ctx("remote_addr", state.RemoteAddr).Add("username", username))
	// GRANT ALL
	//if username != "username" || password != "password" {
	//	_ = logR.Notice("invalid username or password", logger.Ctx("remote_addr", state.RemoteAddr).Add("username", username).Add("password", password))
	//	return nil, errors.New("invalid username or password")
	//}
	_ = logR.Debug(`successfully logged`, logger.Ctx("password", password))
	u := &internal.LoggedUser{Username: username, Password: password}
	return &Session{
		user:        u,
		mail:        &internal.Mail{User: u},
		mailHandler: s.mailHandler,
		logger:      logR,
	}, nil
}

// AnonymousLogin requires clients to authenticate using SMTP AUTH before sending emails
func (s *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	u := &internal.AnonymousUser{Username: "anon_" + s.idGenerator.Generate(5)}
	logR := s.userLoggerFactory(
		fmt.Sprintf("[%s@%s] ", u.Username, state.RemoteAddr),
		logger.Ctx("remote_addr", state.RemoteAddr).Add("username", u.Username),
	)
	_ = logR.Debug(`anonymous logged`, nil)

	return &Session{
		user:        u,
		mail:        &internal.Mail{User: u},
		mailHandler: s.mailHandler,
		logger:      logR,
	}, nil
}

func NewBackend(idGenerator *correlation_id.RandomIdGenerator, mailHandler internal.MailHandler, userLoggerFactory func(string, *logger.Context) logger.LoggerInterface) *Backend {
	return &Backend{
		idGenerator:       idGenerator,
		mailHandler:       mailHandler,
		userLoggerFactory: userLoggerFactory,
	}
}
