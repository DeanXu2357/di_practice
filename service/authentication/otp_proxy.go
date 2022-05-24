package authentication

//go:generate mockgen -destination ../../mocks/otp_proxy/mocks.go -source=./otp_proxy.go -package=mockOtpProxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type OtpProxy interface {
	GetOtp(accountID string) (string, error)
}

func NewOtpProxy() OtpProxy {
	return &otpProxy{}
}

type otpProxy struct {
}

func (o *otpProxy) GetOtp(accountID string) (string, error) {
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
