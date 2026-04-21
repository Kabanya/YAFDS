package utils

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	globalLogger *log.Logger
	globalCloser func() error
	globalOnce   sync.Once
	globalErr    error
)

// Singleton
func InitFileLogger(filename string) error {
	globalOnce.Do(func() {
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			globalErr = err
			return
		}
		globalLogger = log.New(file, "", log.LstdFlags)
		globalCloser = func() error {
			return file.Close()
		}
	})
	return globalErr
}

// return logger
func Logger() (*log.Logger, error) {
	if globalLogger == nil {
		return nil, fmt.Errorf("logger has not been initialized")
	}
	return globalLogger, nil
}

// close logger
func CloseLogger() error {
	if globalCloser == nil {
		return fmt.Errorf("need to initialize logger via InitFileLogger()")
	}
	err := globalCloser()
	globalCloser = nil
	return err
}
