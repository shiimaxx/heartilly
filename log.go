package main

import (
	"go.uber.org/zap"
)

type Logger struct {
	baseLogger *zap.Logger
}

func NewLogger() (*Logger, error) {
	l, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	return &Logger{
		baseLogger: l,
	}, nil
}

func (l *Logger) Info(id int, eyesOn, msg string) {
	l.baseLogger.Info(msg, zap.Int("id", id), zap.String("eyes_on", eyesOn))
}

func (l *Logger) Debug(id int, eyesOn, msg string) {
	l.baseLogger.Debug(msg, zap.Int("id", id), zap.String("eyes_on", eyesOn))
}
