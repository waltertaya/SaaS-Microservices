package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func InitLogger() {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./logs/api-gateway.log",
		MaxSize:    10, // megabytes
		MaxBackups: 5,
		MaxAge:     28,   // days
		Compress:   true, // gzip
	}

	writeSyncer := zapcore.AddSync(lumberJackLogger)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		zapcore.InfoLevel,
	)

	Logger = zap.New(core)
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		traceID := uuid.New().String()
		c.Set("trace_id", traceID)

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		userID := c.GetString("user_id")

		Logger.Info("request",
			zap.String("trace_id", traceID),
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("duration", duration),
			zap.String("client_ip", clientIP),
			zap.String("user_id", userID),
		)
	}
}
