package logging_service

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func NewZerologLogger(service string) *zerologLogger {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	return &zerologLogger{
		service: service,
	}
}

type zerologLogger struct {
	service string
}

func (zl *zerologLogger) Debug(msg string) {
	log.Debug().
		Str("service", zl.service).
		Msg(msg)
}

func (zl *zerologLogger) Debugf(template string, msg string) {
	log.Debug().
		Str("service", zl.service).
		Msgf(template, msg)
}

func (zl *zerologLogger) Info(msg string) {
	log.Info().
		Str("service", zl.service).
		Msg(msg)
}

func (zl *zerologLogger) Infof(template string, msg string) {
	log.Info().
		Str("service", zl.service).
		Msgf(template, msg)
}

func (zl *zerologLogger) Warn(msg string) {
	log.Warn().
		Str("service", zl.service).
		Msg(msg)
}

func (zl *zerologLogger) Warnf(template string, msg string) {
	log.Warn().
		Str("service", zl.service).
		Msgf(template, msg)
}

func (zl *zerologLogger) Error(err error) {
	log.Error().
		Str("service", zl.service).
		Err(err).
		Msg("")
}

func (zl *zerologLogger) Fatal(err error) {
	log.Error().
		Str("service", zl.service).
		Err(err).
		Msg("")
}
