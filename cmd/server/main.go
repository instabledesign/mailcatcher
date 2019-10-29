package main

import (
	"context"
	"os"
	"syscall"
	"time"

	"github.com/gol4ng/logger"
	"github.com/gol4ng/mailcatcher/config"
	"github.com/gol4ng/mailcatcher/internal/service"
	"github.com/gol4ng/signal"
)

func main() {
	mainCtx, mainCancel := context.WithCancel(context.Background())
	container := service.NewContainer(config.NewConfig())
	l := container.GetLogger()

	httpServer := container.GetHTTPServer()
	smtpServer := container.GetSMTPServer()

	httpServer.MustStart(mainCtx)
	smtpServer.MustStart(mainCtx)

	idleConnsClosed := make(chan struct{})
	defer signal.SubscribeWithKiller(func(signal os.Signal) {
		stoppingCtx, cancel := context.WithTimeout(mainCtx, 10*time.Second)
		defer cancel()
		_ = l.Info("signal (%signal%) received : gracefully stopping application", logger.Ctx("signal", signal))
		if err := smtpServer.Stop(stoppingCtx); err != nil {
			_ = l.Error("stopping smtp server error: %error%", logger.Ctx("error", err))
		}
		if err := httpServer.Stop(stoppingCtx); err != nil {
			_ = l.Error("stopping http server error: %error%", logger.Ctx("error", err))
		}

		close(idleConnsClosed)
	}, os.Interrupt, syscall.SIGTERM)()

	<-idleConnsClosed
	mainCancel()
	_ = l.Info("Application stopped", nil)
}
