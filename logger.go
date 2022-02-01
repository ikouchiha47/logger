package logger

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

const LogParamKey = "log_params"

var homeDir = os.Getenv("HOME") + "/"

var (
	std *logrus.Logger
	mu  sync.RWMutex
)

func init() {
	std = logrus.New()
	std.SetFormatter(&Formatter{ChildFormatter: &logrus.JSONFormatter{}})
}

func WithJSONLogger() {
	if std == nil {
		std = logrus.New()
	}

	mu.Lock()
	defer mu.Unlock()

	std.SetFormatter(&Formatter{ChildFormatter: &logrus.JSONFormatter{}})
}

func WithTextLogger() {
	if std == nil {
		std = logrus.New()
	}

	mu.Lock()
	defer mu.Unlock()

	std.SetFormatter(&Formatter{ChildFormatter: &logrus.TextFormatter{}})
}

func SetLogLevel(level string) {
	logLevel := logrus.ErrorLevel

	switch {
	case strings.EqualFold(level, "debug"):
		logLevel = logrus.DebugLevel
	case strings.EqualFold(level, "error"):
		logLevel = logrus.ErrorLevel
	case strings.EqualFold(level, "info"):
		logLevel = logrus.InfoLevel
	}

	std.SetLevel(logLevel)
}

func WithError(err error) *logrus.Entry {
	return std.WithError(err)
}

func WithContext(ctx context.Context) *logrus.Entry {
	if rmap, ok := ctx.Value(LogParamKey).(map[string]interface{}); ok {
		fields := logrus.Fields{}

		for k, v := range rmap {
			fields[k] = v
		}

		return std.WithFields(fields).WithContext(ctx)
	}
	return std.WithContext(ctx)
}

func Error(message ...interface{}) {
	std.Error(message...)
}

func Warnln(message ...interface{}) {
	std.Warnln(message...)
}

func Infoln(message ...interface{}) {
	std.Infoln(message...)
}

func Println(message ...interface{}) {
	std.Println(message...)
}

func Debugln(message ...interface{}) {
	std.Debugln(message...)
}

type Formatter struct {
	ChildFormatter logrus.Formatter
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	function, file, number := getDetails()
	file = strings.TrimPrefix(file, homeDir)
	data := logrus.Fields{"location": fmt.Sprintf("%s:%s %s", file, number, function)}

	for k, v := range entry.Data {
		data[k] = v
	}

	entry.Data = data

	return f.ChildFormatter.Format(entry)
}

func getDetails() (string, string, string) {
	skip := 3

start:
	pc, file, line, _ := runtime.Caller(skip)
	lineNumber := fmt.Sprintf("%d", line)
	function := runtime.FuncForPC(pc).Name()

	if strings.LastIndex(function, "sirupsen/logrus") != -1 ||
		strings.LastIndex(function, "ikouchiha47/logger") != -1 {
		skip++
		goto start
	}

	return function, file, lineNumber
}
