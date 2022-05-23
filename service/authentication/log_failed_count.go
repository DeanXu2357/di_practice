package authentication

import (
	"fmt"
	"log"
	"os"
)

type LogFailedCount interface {
	LogFailedCount(accountID string) error
}

type logFailedCount struct {
	failedCounter FailedCounter
}

func NewLogFailedCount(f FailedCounter) LogFailedCount {
	return &logFailedCount{failedCounter: f}
}

func (a *logFailedCount) LogFailedCount(accountID string) error {
	failedCounts, err := a.failedCounter.Get(accountID)
	if err != nil {
		return fmt.Errorf("get failed count fail %w", err)
	}

	// log failed count
	logger := log.New(os.Stderr, "[Debug] ", 0)
	logger.Printf("failed times: %s", failedCounts)
	return nil
}
