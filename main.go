package main

import (
	"log"

	"github.com/yaches/habr_crawler/cmd"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapWriter func(string, ...zap.Field)

func (f zapWriter) Write(p []byte) (int, error) {
	f(string(p))
	return len(p), nil
}

func ReplaceStdLoggerWithGlobalZap() {
	log.SetFlags(0)
	log.SetOutput(zapWriter(zap.L().Debug))
}

func initLogger() {
	config := zap.NewDevelopmentConfig()
	config.DisableCaller = true
	config.DisableStacktrace = true
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	l, err := config.Build()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(l)
	ReplaceStdLoggerWithGlobalZap()

}

func main() {
	initLogger()
	cmd.Execute()
}
