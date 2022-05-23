package main

import (
	"di_practice/service/authentication"
	"fmt"
)

func main() {
	result, err := authentication.New().Verify("poyu", "pa55w0rd", "123")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("Verify: %t", result)
}
