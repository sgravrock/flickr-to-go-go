package main

import (
	"fmt"

	"os"

	"github.com/sgravrock/flickr-to-go-go/auth"
)

func main() {
	key := os.Getenv("FLICKR_API_KEY")
	secret := os.Getenv("FLICKR_API_SECRET")
	accessToken, err := auth.Authenticate(key, secret, nil, nil)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Printf("Got access code: %s\n", accessToken)
}
