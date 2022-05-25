package authentication

//go:generate mockgen -destination ../../mocks/log_failed_count/mocks.go -source=./log_failed_count.go -package=mockLogFailedCount

import (
	"fmt"
	"log"
)

type LogFailedCount interface {
	LogFailedCount(accountID string) error
}

type logFailedCount struct {
	failedCounter FailedCounter
	logger        *log.Logger
}

func NewLogFailedCount(f FailedCounter, l *log.Logger) LogFailedCount {
	return &logFailedCount{failedCounter: f, logger: l}
}

func (a *logFailedCount) LogFailedCount(accountID string) error {
	failedCounts, err := a.failedCounter.Get(accountID)
	if err != nil {
		return fmt.Errorf("get failed count fail %w", err)
	}

	// log failed count
	a.logger.Printf("failed times: %s", failedCounts)
	return nil
}

func NewLogFailedCountDecorator(a Authentication, l LogFailedCount) Authentication {
	return &logFailedCountDecorator{
		authentication: a,
		logFailedCount: l,
	}
}

type logFailedCountDecorator struct {
	authentication Authentication
	logFailedCount LogFailedCount
}

func (l *logFailedCountDecorator) Verify(accountID, pwd, otp string) (bool, error) {
	isValid, err := l.authentication.Verify(accountID, pwd, otp)
	if err != nil {
		return false, err
	}
	if !isValid {
		return false, l.logFailedCount.LogFailedCount(accountID)
	}

	return isValid, nil
}
