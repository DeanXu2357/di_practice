package authentication

import (
	"testing"

	mockAccountRepo "di_practice/mocks/account_repo"
	mockFailedCounter "di_practice/mocks/failed_counter"
	mockHashPassword "di_practice/mocks/hash_password"
	mockLogFailedCount "di_practice/mocks/log_failed_count"
	mockNotification "di_practice/mocks/notification"
	mockOtpProxy "di_practice/mocks/otp_proxy"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var repo *mockAccountRepo.MockAccountRepo
var otp *mockOtpProxy.MockOtpProxy
var fc *mockFailedCounter.MockFailedCounter
var hash *mockHashPassword.MockHashPassword
var n *mockNotification.MockNotification
var l *mockLogFailedCount.MockLogFailedCount

func Test_authentication_Verify_is_valid(t *testing.T) {
	createMocks(t)

	givenIsAccountLocked("poyu", false)
	givenPwdFromDB("poyu", "hashed_password")
	givenHashedPwd("pa55w0rd", "hashed_password")
	givenCurrentOtp("poyu", "abc")
	resetFailedCountSuccess("poyu")

	shouldBeValid(t, "poyu", "pa55w0rd", "abc")
}

func resetFailedCountSuccess(x string) *gomock.Call {
	return fc.EXPECT().
		Reset(gomock.Eq(x)).
		Return(nil)
}

func givenIsAccountLocked(x string, rets bool) *gomock.Call {
	return fc.EXPECT().
		IsAccountLocked(gomock.Eq(x)).
		Return(rets, nil)
}

func givenPwdFromDB(x string, rets string) *gomock.Call {
	return repo.EXPECT().
		GetPwdFromDB(gomock.Eq(x)).
		Return(rets, nil)
}

func givenHashedPwd(x string, rets string) *gomock.Call {
	return hash.EXPECT().
		HashPassword(gomock.Eq(x)).
		Return(rets)
}

func givenCurrentOtp(x string, rets string) *gomock.Call {
	return otp.EXPECT().
		GetOtp(gomock.Eq(x)).
		Return(rets, nil)
}

func shouldBeValid(t *testing.T, id string, pwd string, o string) {
	svc := New(repo, hash, otp, fc, n, l)
	actual, err := svc.Verify(id, pwd, o)

	assert.NoError(t, err)
	assert.True(t, actual)
}

// golang 中測試沒有 setUp 方法, 所以寫這個
func createMocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo = mockAccountRepo.NewMockAccountRepo(ctrl)
	otp = mockOtpProxy.NewMockOtpProxy(ctrl)
	fc = mockFailedCounter.NewMockFailedCounter(ctrl)
	hash = mockHashPassword.NewMockHashPassword(ctrl)
	n = mockNotification.NewMockNotification(ctrl)
	l = mockLogFailedCount.NewMockLogFailedCount(ctrl)
}
