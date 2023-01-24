package log

import "github.com/sirupsen/logrus"

type Logger struct {
	logger *logrus.Logger
}

func NewLogger() *Logger {
	return &Logger{
		logger: logrus.New(),
	}
}
