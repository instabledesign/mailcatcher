package service

import (
	"sync"
	"time"

	internal_stop_dispatcher "github.com/instabledesign/mailcatcher/internal/stop_dispatcher"
	"github.com/gol4ng/stop-dispatcher"
	"github.com/gol4ng/stop-dispatcher/stop_callback"
)

var stopDispatcherOnce sync.Once

func (container *Container) GetStopDispatcher() *stop_dispatcher.Dispatcher {
	stopDispatcherOnce.Do(func() {
		g := stop_dispatcher.NewDispatcher(stop_dispatcher.WithReasonHandler(internal_stop_dispatcher.Logger(container.GetLogger())))
		g.RegisterCallback(stop_callback.Timeout(10 * time.Second))
		container.stopDispatcher = g
	})

	return container.stopDispatcher
}
