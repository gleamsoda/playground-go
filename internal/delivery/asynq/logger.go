package asynq

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type logger struct{}

var lgr = &logger{}

func (l *logger) Print(level zerolog.Level, args ...interface{}) {
	log.WithLevel(level).Msg(fmt.Sprint(args...))
}

func (l *logger) Printf(ctx context.Context, format string, v ...interface{}) {
	log.WithLevel(zerolog.DebugLevel).Msgf(format, v...)
}

func (l *logger) Debug(args ...interface{}) {
	l.Print(zerolog.DebugLevel, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.Print(zerolog.InfoLevel, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.Print(zerolog.WarnLevel, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.Print(zerolog.ErrorLevel, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.Print(zerolog.FatalLevel, args...)
}
