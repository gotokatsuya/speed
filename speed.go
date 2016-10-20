package speed

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	enableSpeedLogger bool
)

func init() {
	flag.BoolVar(&enableSpeedLogger, "speed", false, "Enable speed logger")
}

// EnableLogger enable speed logger
func EnableLogger() {
	enableSpeedLogger = true
}

// Logger is speed logger obj
type Logger struct {
	loggerRO

	mu sync.RWMutex
}

type loggerRO struct {
	description string
	fileName    string
	caller      string
	info        []string
	beginAt     time.Time
	endAt       time.Time

	dispatcher *log.Logger
}

// NewLogger returns a new Logger
func NewLogger(descriptions ...string) *Logger {
	if !enableSpeedLogger {
		return &Logger{}
	}

	var description string
	switch {
	case len(descriptions) == 0:
		description = "-"
	default:
		description = descriptions[0]
	}
	return &Logger{
		loggerRO: loggerRO{
			description: description,
		},
	}
}

// Description set description to logger
func (l *Logger) Description(description string) *Logger {
	if !enableSpeedLogger {
		return l
	}

	l.description = description
	return l
}

// Copy returns a new Logger
func (l *Logger) Copy() *Logger {
	if !enableSpeedLogger {
		return l
	}

	return &Logger{
		loggerRO: l.loggerRO,
	}
}

// Begin begin trace
func (l *Logger) Begin() *Logger {
	if !enableSpeedLogger {
		return l
	}

	l.mu.Lock()

	l.info = make([]string, 0)
	l.fileName, l.caller = func() (string, string) {
		pc, fileName, line, ok := runtime.Caller(2)
		if ok {
			fn := runtime.FuncForPC(pc).Name()
			return fileName, fmt.Sprintf("%v (%vL)", fn, line)
		}
		return fileName, "no value"
	}()
	l.beginAt = time.Now()
	return l
}

// End end trace
func (l *Logger) End() {
	if !enableSpeedLogger {
		return
	}

	l.endAt = time.Now()
	l.append("description", l.description)
	l.append("file", l.fileName)
	l.append("begin_at", l.beginAt)
	l.append("end_at", l.endAt)
	l.append("caller", l.caller)
	l.append("seconds", fmt.Sprintf("%.6f", l.duration().Seconds()))
	l.append("milliseconds", fmt.Sprintf("%.6f", l.duration().Seconds()*1000.0))
	l.append("microseconds", fmt.Sprintf("%.6f", l.duration().Seconds()*1000000.0))
	l.getDispatcher().Println(strings.Join(l.info, "\t"))

	l.mu.Unlock()
}

func (l *Logger) getDispatcher() *log.Logger {
	const fileNameFormat = "speed-%s.log"
	fileName := fmt.Sprintf(fileNameFormat, time.Now().Format("20060102"))
	filePath := path.Join(os.TempDir(), fileName)
	if _, err := os.Stat(filePath); err == nil && l.dispatcher != nil {
		return l.dispatcher
	}
	handler, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v", err)
		handler = os.Stdout
	}
	l.dispatcher = log.New(handler, "", 0)
	return l.dispatcher
}

func (l *Logger) append(key string, val interface{}) *Logger {
	l.info = append(l.info, fmt.Sprintf("%s:%v", key, val))
	return l
}

func (l *Logger) duration() time.Duration {
	return l.endAt.Sub(l.beginAt)
}
