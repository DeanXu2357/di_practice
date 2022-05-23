package authentication

type Authentication interface {
	Verify(accountID, pwd, otp string) (bool, error)
}

func New() Authentication {
	return &authentication{}
}

type authentication struct{}

func (a *authentication) Verify(accountID, pwd, otp string) (bool, error) {
	return false, nil
}
