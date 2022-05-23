package authentication

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Authentication interface {
	Verify(accountID, pwd, otp string) (bool, error)
}

func New() Authentication {
	return &authentication{}
}

type AccountModel struct {
	gorm.Model
	Password string
}

type authentication struct{}

func (a *authentication) Verify(accountID, pwd, otp string) (bool, error) {
	// check if locked
	isLocked, err := a.checkLocked(accountID)
	if err != nil {
		return false, fmt.Errorf("checkLocked failed %w", err)
	}
	if isLocked == "true" {
		return false, errors.New("account locked")
	}

	// get pwd from db
	pwdFromDB, err := a.getPwdFromDB(accountID)
	if err != nil {
		return false, fmt.Errorf("get Pwd from db failed %w", err)
	}

	// sha256 pwd
	hashedPwd := a.hashPassword(pwd)

	// get opt by http request
	currentOtp, err := a.getOtp(accountID)
	if err != nil {
		return false, fmt.Errorf("get opt failed %w", err)
	}

	if otp == currentOtp && hashedPwd == pwdFromDB {
		if err := a.resetFailedCount(accountID); err != nil {
			return false, fmt.Errorf("reset failed count fail %w", err)
		}

		return true, nil
	} else {
		// get failed count & log failed count
		if err := a.logFailedCount(accountID); err != nil {
			return false, err
		}

		// notify slack
		if err := a.notify(); err != nil {
			return false, fmt.Errorf("notify fail %w", err)
		}

		// add failed count
		if err := a.addFailedCount(accountID); err != nil {
			return false, fmt.Errorf("add failed count fail %w", err)
		}

		return false, nil
	}
}

func (a *authentication) logFailedCount(accountID string) error {
	failedCounts, err := a.getFailedCount(accountID)
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

func (a *authentication) addFailedCount(accountID string) error {
	res, err := http.Get(fmt.Sprintf("https://failed_count/%s/add", accountID))
	if err != nil {
		return fmt.Errorf("http error: %w", err)
	}
	rtnBytes, err := ioutil.ReadAll(res.Body)
	//defer res.Body.Close()
	if err != nil {
		return fmt.Errorf("parse response error: %w", err)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("unexpected response status code: %d, body %s", res.StatusCode, rtnBytes)
	}
	return nil
}

func (a *authentication) getFailedCount(accountID string) (string, error) {
	res, err := http.Get(fmt.Sprintf("https://failed_count/%s", accountID))
	if err != nil {
		return "", fmt.Errorf("http error: %w", err)
	}
	rtnBytes, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return "", fmt.Errorf("parse response error: %w", err)
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("unexpected response status code: %d, body %s", res.StatusCode, rtnBytes)
	}

	return string(rtnBytes), nil
}

func (a *authentication) resetFailedCount(accountID string) error {
	res, err := http.Get(fmt.Sprintf("https://failed_count/%s/reset", accountID))
	if err != nil {
		return fmt.Errorf("http error: %w", err)
	}
	rtnBytes, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return fmt.Errorf("parse response error: %w", err)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("unexpected response status code: %d, body %s", res.StatusCode, rtnBytes)
	}
	return nil
}

func (a *authentication) getOtp(accountID string) (string, error) {
	res, err := http.Get(fmt.Sprintf("https://opt_service/%s", accountID))
	if err != nil {
		return "", fmt.Errorf("http error: %w", err)
	}
	rtnBytes, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return "", fmt.Errorf("parse response error: %w", err)
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("unexpected response status code: %d, body %s", res.StatusCode, rtnBytes)
	}

	return string(rtnBytes), nil
}

func (a *authentication) getPwdFromDB(accountID string) (string, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return "", fmt.Errorf("db error: %w", err)
	}

	var ac AccountModel
	db.First(&ac, accountID)

	pwdFromDB := ac.Password
	return pwdFromDB, nil
}

func (a *authentication) hashPassword(pwd string) string {
	hash := sha256.Sum256([]byte(pwd))
	hashedPwd := hex.EncodeToString(hash[:])
	return hashedPwd
}

func (a *authentication) checkLocked(accountID string) (string, error) {
	res, err := http.Get(fmt.Sprintf("https://is_locked/%s", accountID))
	if err != nil {
		return "", fmt.Errorf("http error: %w", err)
	}
	rtnBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("parse response error: %w", err)
	}
	return string(rtnBytes), nil
}
