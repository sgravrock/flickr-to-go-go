package main

import (
	"fmt"
	"os"

	"github.com/sgravrock/flickr-to-go-go/auth"

	"github.com/sgravrock/flickr-to-go-go/storage"
)

func main() {
	key := os.Getenv("FLICKR_API_KEY")
	secret := os.Getenv("FLICKR_API_SECRET")
	savecreds, dest := parseArgs()
	filestore := storage.FileStorage{dest}
	accessToken, err := auth.Authenticate(key, secret, filestore,
		savecreds, nil, nil)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	fmt.Printf("Got access code: %s\n", accessToken)
}

func parseArgs() (bool, string) {
	if len(os.Args) == 3 && os.Args[1] == "--savecreds" {
		return true, os.Args[2]
	} else if len(os.Args) == 2 {
		return false, os.Args[1]
	} else {
		fmt.Fprintf(os.Stderr, "Usage: %s dest [-savecreds]\n", os.Args[0])
		os.Exit(1)
		return false, "" // not reached
	}
}
