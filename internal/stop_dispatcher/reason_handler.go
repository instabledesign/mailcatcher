package stop_dispatcher

import (
	"fmt"
	"os"

	"github.com/gol4ng/logger"
	stop_dispatcher "github.com/gol4ng/stop-dispatcher"
)

func Logger(l logger.LoggerInterface) stop_dispatcher.ReasonHandler {
	return func(reason stop_dispatcher.Reason) {
		switch v := reason.(type) {
		case os.Signal:
			l.Info("received signal", logger.Any("signal", v.String()))
		case error:
			l.Error("Fatal error", logger.Error("error", v))
		default:
			l.Warning(fmt.Sprintf("received unexpected stop value : %v", v))
		}
	}
}
