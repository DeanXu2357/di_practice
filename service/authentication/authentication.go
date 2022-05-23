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
	res, err := http.Get(fmt.Sprintf("https://is_locked/%s", accountID))
	if err != nil {
		return false, fmt.Errorf("http error: %w", err)
	}
	rtnBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("parse response error: %w", err)
	}
	isLocked := string(rtnBytes)
	if isLocked == "true" {
		return false, errors.New("account locked")
	}

	// get pwd from db
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return false, fmt.Errorf("db error: %w", err)
	}

	var ac AccountModel
	db.First(&ac, accountID)

	pwdFromDB := ac.Password

	// sha256 pwd
	hash := sha256.Sum256([]byte(pwd))
	hashedPwd := hex.EncodeToString(hash[:])

	// get opt by http request
	res, err = http.Get(fmt.Sprintf("https://opt_service/%s", accountID))
	if err != nil {
		return false, fmt.Errorf("http error: %w", err)
	}
	rtnBytes, err = ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return false, fmt.Errorf("parse response error: %w", err)
	}
	if res.StatusCode != 200 {
		return false, fmt.Errorf("unexpected response status code: %d, body %s", res.StatusCode, rtnBytes)
	}
	currentOtp := string(rtnBytes)

	if otp == currentOtp && hashedPwd == pwdFromDB {
		res, err = http.Get(fmt.Sprintf("https://failed_count/%s/reset", accountID))
		if err != nil {
			return false, fmt.Errorf("http error: %w", err)
		}
		rtnBytes, err = ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			return false, fmt.Errorf("parse response error: %w", err)
		}
		if res.StatusCode != 200 {
			return false, fmt.Errorf("unexpected response status code: %d, body %s", res.StatusCode, rtnBytes)
		}

		return true, nil
	} else {
		// get failed count
		res, err = http.Get(fmt.Sprintf("https://failed_count/%s", accountID))
		if err != nil {
			return false, fmt.Errorf("http error: %w", err)
		}
		rtnBytes, err = ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			return false, fmt.Errorf("parse response error: %w", err)
		}
		if res.StatusCode != 200 {
			return false, fmt.Errorf("unexpected response status code: %d, body %s", res.StatusCode, rtnBytes)
		}

		// log failed count
		logger := log.New(os.Stderr, "[Debug] ", 0)
		logger.Printf("failed times: %s", string(rtnBytes))

		// notify slack
		fmt.Println("this is slack api post")

		// add failed count
		res, err = http.Get(fmt.Sprintf("https://failed_count/%s/add", accountID))
		if err != nil {
			return false, fmt.Errorf("http error: %w", err)
		}
		rtnBytes, err = ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			return false, fmt.Errorf("parse response error: %w", err)
		}
		if res.StatusCode != 200 {
			return false, fmt.Errorf("unexpected response status code: %d, body %s", res.StatusCode, rtnBytes)
		}

		return false, nil
	}
}
