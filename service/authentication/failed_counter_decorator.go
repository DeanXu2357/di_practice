package authentication

import (
	"errors"
	"fmt"
)

func NewFailedCounterDecorator(a Authentication, f FailedCounter) Authentication {
	return &failedCounterDecorator{failedCounter: f, authentication: a}
}

type failedCounterDecorator struct {
	authentication Authentication
	failedCounter  FailedCounter
}

func (f *failedCounterDecorator) Verify(accountID, pwd, otp string) (bool, error) {
	isLocked, err := f.failedCounter.IsAccountLocked(accountID)
	if err != nil {
		return false, fmt.Errorf("IsAccountLocked failed %w", err)
	}
	if isLocked {
		return false, errors.New("account locked")
	}

	valid, err := f.authentication.Verify(accountID, pwd, otp)
	if err != nil {
		return false, err
	}

	if valid {
		if err := f.failedCounter.Reset(accountID); err != nil {
			return false, fmt.Errorf("reset failed count fail %w", err)
		}

		return true, nil
	} else {
		if err := f.failedCounter.Add(accountID); err != nil {
			return false, fmt.Errorf("add failed count fail %w", err)
		}

		return false, nil
	}
}
