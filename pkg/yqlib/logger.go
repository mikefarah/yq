package yqlib

import (
	"fmt"
	"log/slog"
	"os"
)

// Logger wraps log/slog providing a printf-style interface used throughout yq.
type Logger struct {
	levelVar slog.LevelVar
	slogger  *slog.Logger
}

func newLogger() *Logger {
	l := &Logger{}
	l.levelVar.Set(slog.LevelWarn)
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: &l.levelVar})
	l.slogger = slog.New(handler)
	return l
}

// SetLevel sets the minimum log level.
func (l *Logger) SetLevel(level slog.Level) {
	l.levelVar.Set(level)
}

// GetLevel returns the current log level.
func (l *Logger) GetLevel() slog.Level {
	return l.levelVar.Level()
}

// IsEnabledFor returns true if the given level is enabled.
func (l *Logger) IsEnabledFor(level slog.Level) bool {
	return l.levelVar.Level() <= level
}

// SetSlogger replaces the underlying slog.Logger (e.g. to configure output format).
func (l *Logger) SetSlogger(logger *slog.Logger) {
	l.slogger = logger
}

func (l *Logger) Debug(msg string) {
	if l.IsEnabledFor(slog.LevelDebug) {
		l.slogger.Debug(msg)
	}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.IsEnabledFor(slog.LevelDebug) {
		l.slogger.Debug(fmt.Sprintf(format, args...))
	}
}

func (l *Logger) Info(msg string) {
	l.slogger.Info(msg)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.slogger.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Warning(msg string) {
	l.slogger.Warn(msg)
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.slogger.Warn(fmt.Sprintf(format, args...))
}

func (l *Logger) Error(msg string) {
	l.slogger.Error(msg)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.slogger.Error(fmt.Sprintf(format, args...))
}
