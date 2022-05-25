package main

import (
	"fmt"
	"log"
	"os"

	"di_practice/service/authentication"
)

func main() {
	ar := authentication.NewAccountRepo()
	op := authentication.NewOtpProxy()
	f := authentication.NewFailedCounter()
	h := authentication.NewSha256Hash()
	n := authentication.NewSlackNotification()
	logger := log.New(os.Stderr, "[Debug] ", 0)
	l := authentication.NewLogFailedCount(f, logger)

	svc := authentication.New(ar, h, op)
	svc = authentication.NewFailedCounterDecorator(svc, f)
	svc = authentication.NewLogFailedCountDecorator(svc, l)
	svc = authentication.NewNotificationDecorator(svc, n)
	result, err := svc.Verify("poyu", "pa55w0rd", "123")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("Verify: %t\n", result)
}
