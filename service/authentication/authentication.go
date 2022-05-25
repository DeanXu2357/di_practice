package authentication

import (
	"fmt"
)

type Authentication interface {
	Verify(accountID, pwd, otp string) (bool, error)
}

func New(ar AccountRepo, h HashPassword, op OtpProxy, n Notification) Authentication {
	return &authentication{
		accountRepo:  ar,
		otpProxy:     op,
		hash:         h,
		notification: n,
	}
}

type authentication struct {
	accountRepo  AccountRepo
	otpProxy     OtpProxy
	hash         HashPassword
	notification Notification
}

func (a *authentication) Verify(accountID, pwd, otp string) (bool, error) {
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
		return true, nil
	} else {
		// Notify -- slack
		if err := a.notification.Notify(accountID); err != nil {
			return false, fmt.Errorf("Notify fail %w", err)
		}

		return false, nil
	}
}
