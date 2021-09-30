package smtp

import (
	"errors"
	"fmt"
	"net"

	"github.com/emersion/go-smtp"
	"github.com/gol4ng/logger"
	stop_dispatcher "github.com/gol4ng/stop-dispatcher"
)

var ErrServerClosed = errors.New("smtp: Server closed")

func StartServerShutdownEmitter(smtpServer *smtp.Server, l logger.LoggerInterface) stop_dispatcher.Emitter {
	return func(stop func(stop_dispatcher.Reason)) {
		l.Info("starting smtp server", logger.Any("addr", smtpServer.Addr))
		err := smtpServer.ListenAndServe()

		if e, ok := err.(*net.OpError); ok && e.Err.Error() == "use of closed network connection" {
			stop(fmt.Errorf("smtp server[%s] : %w", smtpServer.Addr, ErrServerClosed))
		}
	}
}
