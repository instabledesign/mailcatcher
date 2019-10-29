package service

import (
	"net/http"

	httpSrv "github.com/gol4ng/mailcatcher/internal/server/http"
	"github.com/gol4ng/mailcatcher/internal/server/http/handler"
	"github.com/gorilla/mux"
)

func (container *App) GetHTTPServer() *httpSrv.Server {
	return httpSrv.NewServer(
		container.Cfg.HTTPAddr,
		container.getHTTPHandler(),
	)
}

func (container *App) getHTTPHandler() http.Handler {
	router := mux.NewRouter()
	router.Path("/").Methods("GET").HandlerFunc(handler.Home())
	router.Path("/mails").Methods("GET").HandlerFunc(handler.Mail(container.GetMailStorage()))
	router.Path("/event").Methods("GET").HandlerFunc(handler.Event(container.GetMailChan()))
	return router
}
