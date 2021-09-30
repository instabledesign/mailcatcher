package http

import (
	"fmt"
	"net/http"

	"github.com/gol4ng/logger"
	stop_dispatcher "github.com/gol4ng/stop-dispatcher"
)

func StartServerShutdownEmitter(httpServer *http.Server, l logger.LoggerInterface) stop_dispatcher.Emitter {
	return func(stop func(reason stop_dispatcher.Reason)) {
		l.Info("starting http server", logger.Any("addr", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			stop(fmt.Errorf("http server[%s] : %w", httpServer.Addr, err))
		}
	}
}
