package core

import (
	"fmt"
	"os"
	"runtime/trace"

	rl "github.com/gen2brain/raylib-go/raylib"
	"go.uber.org/zap"
)

var logger, _ = zap.NewDevelopment()

var Log = NewLogger("Core")

func NewLogger(name string) *zap.Logger {
	return logger.Named(name)
}

func AttachRaylibLogger(quiet bool) {
	if quiet {
		rl.SetTraceLogLevel(rl.LogNone)
		return
	}
	rl.SetTraceLogCallback(func(logType int, text string) {
		switch logType {
		case int(rl.LogDebug):
			Log.Debug(text)
			break
		case int(rl.LogInfo):
			Log.Info(text)
			break
		case int(rl.LogWarning):
			Log.Warn(text)
			break
		case int(rl.LogError):
			Log.Error(text)
			break
		default:
			Log.Info(text)
			break
		}
	})
}

func AttachEventLogger(quiet bool, verbose bool) {
	if quiet {
		return
	}
	OnEvent("*", func(evt Event[any]) error {
		msg := fmt.Sprintf("%s", evt.EventId)
		if verbose {
			msg += fmt.Sprintf(" -> %#v\n", evt.Data)
		}
		Log.Info(msg)
		return nil
	})
}

func AttachTraceLogger(quiet bool) {
	if quiet {
		return
	}
	f, _ := os.Create("trace.out")
	trace.Start(f)
}
