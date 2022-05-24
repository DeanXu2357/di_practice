package authentication

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
