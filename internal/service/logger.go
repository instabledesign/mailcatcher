package service

import (
	"os"
	"sync"

	"github.com/gol4ng/logger"
	"github.com/gol4ng/logger/handler"
	"github.com/gol4ng/logger/middleware"
	"github.com/instabledesign/mailcatcher/pkg/log"
)

var loggerOnce sync.Once

func (container *Container) GetLogger() logger.LoggerInterface {
	loggerOnce.Do(func() {
		container.logger = logger.NewLogger(
			container.getLoggerHandlerMiddleware().Decorate(
				handler.Stream(os.Stdout, log.NewDefaultFormatter(true, container.Cfg.LogVerboseLevel.Level())),
			),
		)
	})
	return container.logger
}

func (container *Container) GetUserLogger(prefix string, defaultContext *logger.Context) logger.LoggerInterface {
	middlewares := container.getLoggerHandlerMiddleware()
	if defaultContext != nil {
		middlewares = append(middlewares, middleware.Context(defaultContext))
	}

	return logger.NewLogger(middlewares.Decorate(
		handler.Stream(os.Stdout, log.NewPrefixFormatter(prefix, true, container.Cfg.LogVerboseLevel.Level())),
	))
}

func (container *Container) getLoggerHandler() logger.HandlerInterface {
	return container.getLoggerHandlerMiddleware().Decorate(
		handler.Stream(os.Stdout, log.NewDefaultFormatter(true, container.Cfg.LogVerboseLevel.Level())),
	)
}

var loggerHandlerMiddlewareOnce sync.Once

func (container *Container) getLoggerHandlerMiddleware() logger.Middlewares {
	loggerHandlerMiddlewareOnce.Do(func() {
		stack := logger.MiddlewareStack(
			middleware.Placeholder(),
			middleware.MinLevelFilter(container.Cfg.LogLevel.Level()),
		)
		if container.Cfg.Debug {
			stack = append(stack, middleware.Caller(3))
		}
		container.loggerMiddlewares = stack
	})

	return container.loggerMiddlewares
}
