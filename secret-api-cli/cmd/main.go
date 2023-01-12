package main

import (
	"fmt"

	"worachanon226/gophercises/secret-api-cli"
)

func main() {
	v := secret.Memory("my-fake-key")
	err := v.Set("demo_key", "some crazy value")
	if err != nil {
		panic(err)
	}
	plain, err := v.Get("demo_key")
	if err != nil {
		panic(err)
	}
	fmt.Println("Plain:", plain)
}
