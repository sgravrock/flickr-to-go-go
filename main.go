package main

import (
	"fmt"
	"os"

	"github.com/sgravrock/flickr-to-go-go/app"
	"github.com/sgravrock/flickr-to-go-go/auth"
	"github.com/sgravrock/flickr-to-go-go/dl"
	"github.com/sgravrock/flickr-to-go-go/storage"
)

func main() {
	key := os.Getenv("FLICKR_API_KEY")
	secret := os.Getenv("FLICKR_API_SECRET")
	savecreds, dest := parseArgs()
	filestore := storage.NewFileStorage(dest)
	authenticator := auth.NewAuthenticator(key, secret, filestore, nil, nil)
	downloader := dl.NewDownloader()
	exitcode := app.Run("https://api.flickr.com", savecreds, authenticator,
		downloader, filestore, os.Stdout, os.Stderr)
	os.Exit(exitcode)
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
