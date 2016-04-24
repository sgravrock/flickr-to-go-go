package main

import (
	"fmt"
	"os"

	"github.com/sgravrock/flickr-to-go-go/auth"

	"io/ioutil"

	"github.com/sgravrock/flickr-to-go-go/storage"
)

func main() {
	key := os.Getenv("FLICKR_API_KEY")
	secret := os.Getenv("FLICKR_API_SECRET")
	savecreds, dest := parseArgs()
	filestore := storage.FileStorage{dest}
	httpClient, err := auth.Authenticate(key, secret, filestore,
		savecreds, nil, nil)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	resp, err := httpClient.Get("https://api.flickr.com/services/rest/?method=flickr.test.login&format=json&nojsoncallback=1")
	if err != nil {
		fmt.Printf("Couldn't verify login: %s\n", err.Error())
		os.Exit(1)
	}

	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Couldn't verify login: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println(string(buf))
}

func parseArgs() (bool, string) {
	if len(os.Args) == 3 && os.Args[1] == "--savecreds" {
		return true, os.Args[2]
	} else if len(os.Args) == 2 {
		return false, os.Args[1]
	} else {
		fmt.Fprintf(os.Stderr, "Usage: %s [--savecreds] dest\n", os.Args[0])
		os.Exit(1)
		return false, "" // not reached
	}
}
