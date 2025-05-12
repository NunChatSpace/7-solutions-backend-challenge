package logger

import (
	"encoding/json"
	"time"

	cfgmd "github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/http/middlewares/config"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/di"
	"github.com/savsgio/atreugo/v11"
	"github.com/sirupsen/logrus"
)

type loggerKey struct{}

func FromContext(ctx *atreugo.RequestCtx, name string) logrus.FieldLogger {
	value := ctx.Value(&loggerKey{}) // GET
	if value == nil {
		return nil
	}
	logger := value.(*logrus.Entry)

	return logger.WithField("caller", name)
}

func Handler(ctx *atreugo.RequestCtx, dep *di.Dependency) error {
	start := time.Now()
	logger := logrus.New()
	cfg := cfgmd.FromContext(ctx)
	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		logger.Error("failed parsing log level, falling back to debug level")
		level = logrus.DebugLevel
	}
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(level)
	entry := logger.WithFields(logrus.Fields{
		"method": string(ctx.Method()),
		"path":   string(ctx.Path())})

	body := ctx.RequestCtx.Request.Body()
	if len(body) > 0 {
		var reqBody = make(map[string]interface{})
		err = json.Unmarshal(body, &reqBody)
		if err == nil {
			entry = entry.WithField("payload", reqBody)
		}

		ctx.SetUserValue("body", body)
	}

	dep.Logger = entry
	err = ctx.Next()

	// Log execution time
	duration := time.Since(start)
	entry = entry.WithField("duration", duration.String())
	if err != nil {
		entry.WithError(err).Error("request failed")
	} else {
		entry.Info("request completed")
	}

	return err
}
