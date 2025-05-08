package logger

import (
	"encoding/json"

	cfgmd "github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/http/middlewares/config"
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

func Handler(ctx *atreugo.RequestCtx) error {
	logger := logrus.New()
	cfg := cfgmd.FromContext(ctx)
	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		logger.Error("failed parsing log level, falling back to debug level")
		level = logrus.DebugLevel
	}
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(level)

	// add request data

	entry := logger.WithFields(logrus.Fields{
		"method": ctx.Method(),
		"path":   ctx.Path()})

	// add request payload
	body := ctx.RequestCtx.Request.Body()
	if len(body) > 0 {
		var reqBody = make(map[string]interface{})
		err = json.Unmarshal(body, &reqBody)
		if err == nil {
			entry = entry.WithField("payload", reqBody)
		}

		ctx.SetUserValue("body", body)
	}

	// add actor
	// a := actor.FromContext(ctx)
	// if a != nil && a.ID != 0 {
	// 	entry = entry.WithField("user_id", a.ID)
	// }

	ctx.RequestCtx.SetUserValue(&loggerKey{}, entry)
	return ctx.Next()
}
