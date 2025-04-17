package logging

import (
	"fmt"
	"runtime"
	"time"
)

type logger struct {
	level            LogLevel
	colorProvider    func(LogLevel) string
	useSourceContext bool
}

const (
	colorClear       = "\033[0m"
	colorTrace       = "\033[0;37m"
	colorVerbose     = "\033[0;36m"
	colorDebug       = "\033[0;32m"
	colorInformation = "\033[0;34m"
	colorWarning     = "\033[0;33m"
	colorError       = "\033[0;31m"
	colorFatal       = "\033[0;35m"
)

func (l logger) Log(level LogLevel, msg string, skip int, args ...interface{}) {
	if level < l.level {
		return
	}
	var sourceContext string
	if l.useSourceContext {
		_, file, line, ok := runtime.Caller(skip)
		if ok == false {
			// if we can't get the caller, we disable the source context
			l.Error("Can't get caller, disabling source context")
			l.useSourceContext = false
		} else {
			sourceContext = fmt.Sprintf(" [%s:%d]", file, line)
		}
	}
	if l.colorProvider != nil {
		// [2025-04-15 15:13:45 D] [main.go:15]: message
		fmt.Printf("%s[%s %c]%s: %s%s\n", l.colorProvider(level), time.Now().Format("2006-01-02 15:04:05"), level.String()[0], sourceContext, fmt.Sprintf(msg, args...), l.colorProvider(None))
	} else {
		fmt.Printf("[%s %c]%s: %s\n", time.Now().Format("2006-01-02 15:04:05"), level.String()[0], sourceContext, fmt.Sprintf(msg, args...))
	}
}

func NewLogger(level LogLevel) *logger {
	return &logger{
		level:            level,
		colorProvider:    DefaultColor,
		useSourceContext: true,
	}
}

func (l logger) Trace(msg string, args ...interface{}) {
	l.Log(Trace, msg, 2, args...)
}

func (l logger) Verbose(msg string, args ...interface{}) {
	l.Log(Verbose, msg, 2, args...)
}

func (l logger) Debug(msg string, args ...interface{}) {
	l.Log(Debug, msg, 2, args...)
}

func (l logger) Information(msg string, args ...interface{}) {
	l.Log(Information, msg, 2, args...)
}

func (l logger) Warning(msg string, args ...interface{}) {
	l.Log(Warning, msg, 2, args...)
}

func (l logger) Error(msg string, args ...interface{}) {
	l.Log(Error, msg, 2, args...)
}

func (l logger) Fatal(msg string, args ...interface{}) {
	l.Log(Fatal, msg, 2, args...)
}

func (l logger) LogLevel() LogLevel {
	return l.level
}
func DefaultColor(l LogLevel) string {
	switch l {
	case None:
		return colorClear
	case Trace:
		return colorTrace
	case Verbose:
		return colorVerbose
	case Debug:
		return colorDebug
	case Information:
		return colorInformation
	case Warning:
		return colorWarning
	case Error:
		return colorError
	case Fatal:
		return colorFatal
	default:
		return colorClear
	}
}
