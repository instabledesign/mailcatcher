package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gol4ng/logger"
)

type Server struct {
	httpAddr    string
	httpHandler http.Handler
	httpServer  *http.Server
	logger      logger.LoggerInterface

	requestContextCancel context.CancelFunc
}

func (s *Server) MustStart(baseContext context.Context) {
	go func() {
		if err := s.Start(baseContext); err != http.ErrServerClosed {
			_ = s.logger.Critical("http server error: %error%", logger.Ctx("error", err))
			panic(err)
		}
	}()
}

func (s *Server) Start(baseContext context.Context) error {
	if s.httpServer != nil {
		err := errors.New("http server already started")
		_ = s.logger.Error("%error%", logger.NewContext().Add("error", err))
		return err
	}

	_ = s.logger.Debug("starting http server...", nil)
	ctx, cancel := context.WithCancel(baseContext)
	s.requestContextCancel = cancel
	s.httpServer = &http.Server{
		Addr:    s.httpAddr,
		Handler: s.httpHandler,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer == nil {
		_ = s.logger.Notice("http server already stopped", nil)
		return nil
	}
	if d, ok := ctx.Deadline(); ok {
		_ = s.logger.Debug("stopping http server...", logger.Ctx("shutdown_timeout", d.Format(time.RFC3339)))
	}

	defer func() {
		s.requestContextCancel()
		s.requestContextCancel = nil
		s.httpServer = nil
	}()

	return s.httpServer.Shutdown(ctx)
}

func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		httpAddr:    addr,
		httpHandler: handler,
		logger:      logger.NewNopLogger(),
	}
}
