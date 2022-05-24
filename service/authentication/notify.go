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
