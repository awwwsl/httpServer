package logging

import (
	"encoding/json"
	"fmt"
)

type LogLevel int

//goland:noinspection GoMixedReceiverTypes
func (l LogLevel) MarshalJSON() ([]byte, error) {
	return []byte(`"` + l.String() + `"`), nil
}

//goland:noinspection GoMixedReceiverTypes
func (l *LogLevel) UnmarshalJSON(bytes []byte) error {
	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}
	switch s {
	case "Lowest":
		*l = Lowest
	case "Trace":
		*l = Trace
	case "Verbose":
		*l = Verbose
	case "Debug":
		*l = Debug
	case "Information":
		*l = Information
	case "Warning":
		*l = Warning
	case "Error":
		*l = Error
	case "Fatal":
		*l = Fatal
	case "Highest":
		*l = Highest
	default:
		return fmt.Errorf("unknown log level: %s", s)
	}
	return nil
}

const (
	// None LogLevel is used for color terminal output clear and no logging
	None LogLevel = iota
	// Lowest LogLevel is the lowest log level, used for no logging
	Lowest
	// Trace is the lowest log level, used for logging all messages, including those sensitive to security
	Trace
	// Verbose is used for logging verbose messages, usually used for deeper debugging and should represent call stack and control flow
	Verbose
	// Debug is used for logging debug messages, usually used for debugging
	Debug
	// Information is used for logging information messages, usually used for logging information about the application
	Information
	// Warning is used for logging warning messages, usually used for logging warnings about the application
	Warning
	// Error is used for logging error messages, usually used for logging something may impact the application's stability
	Error
	// Fatal is used for logging fatal messages, usually used for logging something that will crash the application
	Fatal
	// Highest is the highest log level, used for always logging
	Highest
)

type ILogger interface {
	Trace(msg string, args ...interface{})
	Verbose(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Information(msg string, args ...interface{})
	Warning(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	Log(level LogLevel, msg string, skip int, args ...interface{})

	LogLevel() LogLevel
}

//goland:noinspection GoMixedReceiverTypes
func (l LogLevel) String() string {
	switch l {
	case Lowest:
		return "Lowest"
	case Trace:
		return "Trace"
	case Verbose:
		return "Verbose"
	case Debug:
		return "Debug"
	case Information:
		return "Information"
	case Warning:
		return "Warning"
	case Error:
		return "Error"
	case Fatal:
		return "Fatal"
	case Highest:
		return "Highest"
	default:
		return "Unknown"
	}
}
