package services

import (
	"net"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seantheyahn/simple-wiki/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//LoggerCore zap base logger instance
var LoggerCore *zap.Logger

//Logger zap sugared logger instance
var Logger *zap.SugaredLogger

//LoggerMiddleware gin middleware to log the requests
func LoggerMiddleware() func(*gin.Context) {
	return func(c *gin.Context) {
		start := time.Now().UTC()
		//keep the unmodified values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		end := time.Now().UTC()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				Logger.Error(e)
			}
		} else {
			Logger.Debugf("%s %d %s%s - %v - %s", c.Request.Method, c.Writer.Status(), path, query, c.ClientIP(), latency.String())
		}
	}
}

//PanicRecoveryMiddleware gin middleware to recover from panic and log the error
func PanicRecoveryMiddleware() func(*gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			// from gin source: check for a broken connection, as it is not really a
			// condition that warrants a panic stack trace.
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				if brokenPipe {
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				Logger.Errorf("[Recovery from panic] %v", err)
				c.String(500, "internal-server-error")
				c.Abort()
			}
		}()
		c.Next()
	}
}

func initLogger() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = nil
	conf := zap.NewProductionConfig()
	conf.EncoderConfig = encoderConfig
	conf.Encoding = "console"
	conf.Level = config.Instance.Logging.Level

	l, _ := conf.Build()
	LoggerCore = l
	Logger = l.Sugar()
}

func finalizeLogger() {
	LoggerCore.Sync()
}
