package middleware

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"log"
	"runtime"
)

const ContextLoggerName = "ContextLogger"

func Logger() echo.MiddlewareFunc {
	logger := logrus.New()

	logger.Level = logrus.DebugLevel
	logrus.SetFormatter(&logrus.JSONFormatter{})
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(logger.Writer())
	logger.Level = logrus.InfoLevel

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logEntry := logrus.NewEntry(logger)
			id := c.Request().Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = c.Response().Header().Get(echo.HeaderXRequestID)
			}
			logEntry = logEntry.WithField("request_id", id)

			req := c.Request()
			c.SetRequest(req.WithContext(context.WithValue(req.Context(), ContextLoggerName, logEntry)))

			return next(c)
		}
	}
}

type CallkerHook struct{}

func (c *CallkerHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
	}
}
func (c *CallkerHook) Fire(entry *logrus.Entry) error {
	var ok bool
	_, file, line, ok := runtime.Caller(4)
	if !ok {
		file = "???"
		line = 0
	}
	entry.Data["caller"] = fmt.Sprintf("%s:%d", file, line)
	return nil
}
