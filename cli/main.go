package main

import (
	"fmt"

	"github.com/tskinn/pomogo"
)

func main() {
	client := pomogo.Client{}
	err := client.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := client.SessionStatus("")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp)
}
