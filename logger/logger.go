package logger

import (
	"errors"
	"github.com/hecatoncheir/Hecatoncheir/broker"
	"time"
)

type LogData struct {
	ApiVersion, Message, Service string
	Time                         time.Time
}

type LogWriter struct {
	LoggerTopic string
	bro         *broker.Broker
}

func New(topicForWriteLog string, broker *broker.Broker) *LogWriter {
	logger := LogWriter{LoggerTopic: topicForWriteLog, bro: broker}
	return &logger
}

var (
	ErrLogDataWithoutTime = errors.New("log data without time")
)

func (logger *LogWriter) Write(logData LogData) error {
	if logData.Time.IsZero() {
		return ErrLogDataWithoutTime
	}

	err := logger.bro.WriteToTopic(logger.LoggerTopic, logData)
	if err != nil {
		return err
	}

	return nil
}
