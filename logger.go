package zerolog_watermill_adapter

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/rs/zerolog"
	"io"
	"log"
)

// NewZerologLoggerAdapter creates ZerologLoggerAdapter which sends info to io.Writer and errors to zerolog.Logger
func NewZerologLoggerAdapter(out io.Writer, logger zerolog.Logger, debug bool, trace bool) watermill.LoggerAdapter {
	l := log.New(out, "[watermill] ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	a := &ZerologLoggerAdapter{InfoLogger: l, ErrorLogger: logger}

	if debug {
		a.DebugLogger = l
	}
	if trace {
		a.TraceLogger = l
	}

	return a
}

type ZerologLoggerAdapter struct {
	watermill.StdLoggerAdapter
	ErrorLogger zerolog.Logger
	InfoLogger  *log.Logger
	DebugLogger *log.Logger
	TraceLogger *log.Logger

	fields watermill.LogFields
}

func (l *ZerologLoggerAdapter) Error(msg string, err error, fields watermill.LogFields) {
	// prevent behavior when sentry merges all errors from handlers to one issue
	if msg == "Handler returned error" {
		l.ErrorLogger.Error().Interface("fields", fields).Err(err).Msg(err.Error())
		return
	}

	l.ErrorLogger.Error().Interface("fields", fields).Err(err).Msg(msg)
}
