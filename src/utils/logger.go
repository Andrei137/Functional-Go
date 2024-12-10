package functional

import (
    "time"
    "go.uber.org/zap"
    "github.com/fatih/color"
    "go.uber.org/zap/zapcore"
)

func customLevelEncoder() zapcore.LevelEncoder {
    return func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
        var levelStr string
        switch level {
            case zapcore.InfoLevel:
                levelStr = color.HiGreenString("[INFO]")
            case zapcore.ErrorLevel:
                levelStr = color.HiRedString("[ERROR]")
            default:
                levelStr = level.CapitalString()
        }
        enc.AppendString(levelStr)
    }
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
    enc.AppendString(t.Format("15:04:05.000"))
}

func InitLogger() *zap.Logger {
    config := zap.NewDevelopmentConfig()
    config.DisableCaller = true
    config.DisableStacktrace = true
    config.EncoderConfig.EncodeTime = customTimeEncoder
    config.EncoderConfig.EncodeLevel = customLevelEncoder()

    logger, err := config.Build()
    if err != nil {
        panic(err)
    }
    return logger
}
