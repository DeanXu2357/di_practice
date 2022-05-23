package authentication

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

	if hashedPwd != pwdFromDB {
		return false, nil
	}

	// get opt by http request
	res, err := http.Get(fmt.Sprintf("https://opt_service/%s", accountID))
	if err != nil {
		return false, fmt.Errorf("http error: %w", err)
	}
	defer res.Body.Close()
	rtnBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("parse response error: %w", err)
	}
	var resMap map[string]interface{}
	json.Unmarshal(rtnBytes, &resMap)
	currentOtp := resMap["otp"]

	if currentOtp != otp {
		return false, nil
	}

	return true, nil
}
