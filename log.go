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

func (l *Logger) Fatal(id int, eyesOn, msg string) {
	l.baseLogger.Fatal(msg, zap.Int("id", id), zap.String("eyes_on", eyesOn))
}

func (l *Logger) Error(id int, eyesOn, msg string) {
	l.baseLogger.Error(msg, zap.Int("id", id), zap.String("eyes_on", eyesOn))
}

func (l *Logger) Warn(id int, eyesOn, msg string) {
	l.baseLogger.Warn(msg, zap.Int("id", id), zap.String("eyes_on", eyesOn))
}

func (l *Logger) Info(id int, eyesOn, msg string) {
	l.baseLogger.Info(msg, zap.Int("id", id), zap.String("eyes_on", eyesOn))
}

func (l *Logger) Debug(id int, eyesOn, msg string) {
	l.baseLogger.Debug(msg, zap.Int("id", id), zap.String("eyes_on", eyesOn))
}
