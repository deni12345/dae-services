package logx

import (
	"io"
	"log/slog"
)

var logx *slog.Logger

func NewLoggerJSON(w io.Writer) {
	logx = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: slog.LevelInfo}))
}

func Info(msg string, args ...any) {
	if logx == nil {
		return
	}
	logx.Info(msg, args...)
}

func Debug(msg string, args ...any) {
	if logx == nil {
		return
	}
	logx.Debug(msg, args...)
}

func Fatalf(msg string, args ...any) {
	if logx == nil {
		return
	}
	logx.Error(msg, args...)
}
