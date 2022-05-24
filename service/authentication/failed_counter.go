package authentication

//go:generate mockgen -destination ../../mocks/failed_counter/mocks.go -source=./failed_counter.go -package=mockFailedCounter

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type FailedCounter interface {
	Add(accountID string) error
	Get(accountID string) (string, error)
	Reset(accountID string) error
	IsAccountLocked(accountID string) (bool, error)
}

func NewFailedCounter() FailedCounter {
	return &failedCounter{}
}

type failedCounter struct {
}

func (f *failedCounter) Add(accountID string) error {
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

func (f *failedCounter) Get(accountID string) (string, error) {
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

func (f *failedCounter) Reset(accountID string) error {
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

func (f *failedCounter) IsAccountLocked(accountID string) (bool, error) {
	res, err := http.Get(fmt.Sprintf("https://is_locked/%s", accountID))
	if err != nil {
		return false, fmt.Errorf("http error: %w", err)
	}
	rtnBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, fmt.Errorf("parse response error: %w", err)
	}

	if string(rtnBytes) == "true" {
		return true, nil
	}

	return false, nil
}
