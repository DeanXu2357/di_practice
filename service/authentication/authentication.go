package authentication

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
)

type Authentication interface {
	Verify(accountID, pwd, otp string) (bool, error)
}

func New() Authentication {
	ar := NewAccountRepo()
	op := NewOtpProxy()
	f := NewFailedCounter()

	return &authentication{
		accountRepo:   ar,
		otpProxy:      op,
		failedCounter: f,
	}
}

type authentication struct {
	accountRepo   AccountRepo
	otpProxy      OtpProxy
	failedCounter FailedCounter
}

func (a *authentication) Verify(accountID, pwd, otp string) (bool, error) {
	// check if locked
	isLocked, err := a.failedCounter.IsAccountLocked(accountID)
	if err != nil {
		return false, fmt.Errorf("IsAccountLocked failed %w", err)
	}
	if isLocked == "true" {
		return false, errors.New("account locked")
	}

	// get pwd from db
	pwdFromDB, err := a.accountRepo.GetPwdFromDB(accountID)
	if err != nil {
		return false, fmt.Errorf("get Pwd from db failed %w", err)
	}

	// sha256 pwd
	hashedPwd := a.hashPassword(pwd)

	// get opt by http request
	currentOtp, err := a.otpProxy.GetOtp(accountID)
	if err != nil {
		return false, fmt.Errorf("get opt failed %w", err)
	}

	if otp == currentOtp && hashedPwd == pwdFromDB {
		if err := a.failedCounter.Reset(accountID); err != nil {
			return false, fmt.Errorf("reset failed count fail %w", err)
		}

		return true, nil
	} else {
		// add failed count
		if err := a.failedCounter.Add(accountID); err != nil {
			return false, fmt.Errorf("add failed count fail %w", err)
		}

		// get failed count & log failed count
		if err := a.logFailedCount(accountID); err != nil {
			return false, err
		}

		// notify slack
		if err := a.notify(); err != nil {
			return false, fmt.Errorf("notify fail %w", err)
		}

		return false, nil
	}
}

func (a *authentication) logFailedCount(accountID string) error {
	failedCounts, err := a.failedCounter.Get(accountID)
	if err != nil {
		return fmt.Errorf("get failed count fail %w", err)
	}

	// log failed count
	logger := log.New(os.Stderr, "[Debug] ", 0)
	logger.Printf("failed times: %s", failedCounts)
	return nil
}

func (a *authentication) notify() error {
	fmt.Println("this is slack api post")
	return nil
}

func (a *authentication) hashPassword(pwd string) string {
	hash := sha256.Sum256([]byte(pwd))
	hashedPwd := hex.EncodeToString(hash[:])
	return hashedPwd
}
