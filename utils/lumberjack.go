package utils

import (
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"idea_server/global"
	"os"
	"path"
)

func GetWriteSyncer() (zapcore.WriteSyncer) {
	l := &lumberjack.Logger{
		Filename:   path.Join(global.IDEA_CONFIG.Zap.Director, "idea.log"),
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}
	if global.IDEA_CONFIG.Zap.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(l))
	}
	return zapcore.AddSync(l)
}
