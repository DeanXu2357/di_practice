package authentication

//go:generate mockgen -destination ../../mocks/notification/mocks.go -source=./notify.go -package=mockNotification

import "fmt"

type Notification interface {
	Notify(accountID string) error
}

type slackAdapter struct {
}

func NewSlackNotification() Notification {
	return &slackAdapter{}
}

func (a *slackAdapter) Notify(string) error {
	fmt.Println("this is slack api post")
	return nil
}

func NewNotificationDecorator(a Authentication, n Notification) Authentication {
	return &notificationDecorator{a, n}
}

type notificationDecorator struct {
	authentication Authentication
	n              Notification
}

func (nd *notificationDecorator) Verify(accountID, pwd, otp string) (bool, error) {
	isValid, err := nd.authentication.Verify(accountID, pwd, otp)
	if err != nil {
		return false, err
	}
	if !isValid {
		return isValid, nd.n.Notify(accountID)
	}
	return isValid, nil
}
