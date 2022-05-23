package main

import (
	"fmt"

	"di_practice/service/authentication"
)

func main() {
	result, err := authentication.New().Verify("poyu", "pa55w0rd", "123")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("Verify: %t\n", result)
}
