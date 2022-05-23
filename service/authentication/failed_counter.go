package authentication

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type FailedCounter interface {
	addFailedCount(accountID string) error
	getFailedCount(accountID string) (string, error)
	resetFailedCount(accountID string) error
	isAccountLocked(accountID string) (string, error)
}

func NewFailedCounter() FailedCounter {

}

type failedCounter struct {
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

func (a *authentication) isAccountLocked(accountID string) (string, error) {
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
