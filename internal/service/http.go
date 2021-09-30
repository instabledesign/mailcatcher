package service

import (
	"context"
	"net"
	"net/http"
	"net/http/pprof"
	"sync"

	"github.com/instabledesign/mailcatcher/internal/server/http/handler"
	"github.com/gorilla/mux"
)

var httpServerOnce sync.Once

func (container *Container) GetHTTPServer() *http.Server {
	httpServerOnce.Do(func() {
		container.httpServer = &http.Server{
			Handler:        container.getHTTPHandler(),
			Addr:           container.Cfg.HTTPAddr,
			MaxHeaderBytes: http.DefaultMaxHeaderBytes,
			BaseContext: func(_ net.Listener) context.Context {
				return container.baseContext
			},
		}
	})
	return container.httpServer
}

func (container *Container) getHTTPHandler() http.Handler {
	router := mux.NewRouter()
	router.Path("/").Methods("GET").HandlerFunc(handler.Home())
	router.Path("/mails").Methods("GET").HandlerFunc(handler.Mail(container.GetMailStorage()))
	router.Path("/event").Methods("GET").HandlerFunc(handler.Event(container.GetBroker()))

	if container.Cfg.Debug {
		pproffRouter := router.PathPrefix("/debug/pprof").Subrouter()
		pproffRouter.HandleFunc("/", pprof.Index)
		pproffRouter.HandleFunc("/{t:(?:allocs|block|heap|goroutine|mutex|threadcreate)}", pprof.Index)

		pproffRouter.HandleFunc("/cmdline", pprof.Cmdline)
		pproffRouter.HandleFunc("/profile", pprof.Profile)
		pproffRouter.HandleFunc("/symbol", pprof.Symbol)
		pproffRouter.HandleFunc("/trace", pprof.Trace)
	}

	return router
}
