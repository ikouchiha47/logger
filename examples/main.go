package main

import (
	"context"
	"errors"
	"sync"

	"github.com/ikouchiha47/logger"
)

func defaultLog() {
	logger.Infoln("default log")
}

func textLog() {
	logger.WithTextLogger()
	logger.Infoln("text log")

	logger.WithJSONLogger()
}

func errorWithContextLog() {
	err := errors.New("test error")
	ctx := context.WithValue(context.Background(), logger.LogParamKey, map[string]interface{}{"a": 1, "b": "2"})

	logger.WithError(err).Error("failed")
	logger.WithContext(ctx).WithError(err).Error("failed with context")
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defaultLog()
		wg.Done()
	}()

	go func() {
		textLog()
		wg.Done()
	}()

	go func() {
		errorWithContextLog()
		wg.Done()
	}()

	wg.Wait()
}
