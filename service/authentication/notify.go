package authentication

import "fmt"

type Notification interface {
	Notify() error
}

type slackAdapter struct {
}

func NewSlackNotification() Notification {
	return &slackAdapter{}
}

func (a *slackAdapter) Notify() error {
	fmt.Println("this is slack api post")
	return nil
}
