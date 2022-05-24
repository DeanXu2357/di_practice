package authentication

import (
	"errors"
	"fmt"
)

type Authentication interface {
	Verify(accountID, pwd, otp string) (bool, error)
}

func New(ar AccountRepo, h HashPassword, op OtpProxy, f FailedCounter, n Notification, l LogFailedCount) Authentication {
	return &authentication{
		accountRepo:    ar,
		otpProxy:       op,
		failedCounter:  f,
		hash:           h,
		notification:   n,
		logFailedCount: l,
	}
}

type authentication struct {
	accountRepo    AccountRepo
	otpProxy       OtpProxy
	failedCounter  FailedCounter
	hash           HashPassword
	notification   Notification
	logFailedCount LogFailedCount
}

func (a *authentication) Verify(accountID, pwd, otp string) (bool, error) {
	// check if locked
	isLocked, err := a.failedCounter.IsAccountLocked(accountID)
	if err != nil {
		return false, fmt.Errorf("IsAccountLocked failed %w", err)
	}
	if isLocked {
		return false, errors.New("account locked")
	}

	// get pwd from db
	pwdFromDB, err := a.accountRepo.GetPwdFromDB(accountID)
	if err != nil {
		return false, fmt.Errorf("get Pwd from db failed %w", err)
	}

	// sha256 pwd
	hashedPwd := a.hash.HashPassword(pwd)

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
		if err := a.logFailedCount.LogFailedCount(accountID); err != nil {
			return false, err
		}

		// Notify -- slack
		if err := a.notification.Notify(); err != nil {
			return false, fmt.Errorf("Notify fail %w", err)
		}

		return false, nil
	}
}
