// Package logs defines log helper functions
package logs

import (
	"log"
	"math"

	"github.com/supinf/reinvent-sessions-api/app/config"
)

type Level int

const (
	FATAL Level = 1 + iota
	ERROR
	WARN
	INFO
	DEBUG
	TRACE
)

var level Level

func init() {
	level = Level(math.Min(float64(TRACE), math.Max(float64(FATAL), float64(config.NewConfig().LogLevel))))
}

func Trace(message string) {
	if level >= TRACE {
		log.Print(message)
	}
}

func Tracef(format string, v ...interface{}) {
	if level >= TRACE {
		log.Printf(format, v...)
	}
}

func Debug(message string) {
	if level >= DEBUG {
		log.Print(message)
	}
}

func Debugf(format string, v ...interface{}) {
	if level >= DEBUG {
		log.Printf(format, v...)
	}
}

func Info(message string) {
	if level >= INFO {
		log.Print(message)
	}
}

func Infof(format string, v ...interface{}) {
	if level >= INFO {
		log.Printf(format, v...)
	}
}

func Warn(message string) {
	if level >= WARN {
		log.Print(message)
	}
}

func Warnf(format string, v ...interface{}) {
	if level >= WARN {
		log.Printf(format, v...)
	}
}

func Error(message string) {
	if level >= ERROR {
		log.Print(message)
	}
}

func Errorf(format string, v ...interface{}) {
	if level >= ERROR {
		log.Printf(format, v...)
	}
}

func Fatal(v ...interface{}) {
	if level >= FATAL {
		log.Fatal(v)
	}
}

func Fatalf(format string, v ...interface{}) {
	if level >= FATAL {
		log.Fatalf(format, v...)
	}
}
