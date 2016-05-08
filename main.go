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
	key := requireEnv("FLICKR_API_KEY")
	secret := requireEnv("FLICKR_API_SECRET")
	savecreds, dest := parseArgs()
	filestore := storage.NewFileStorage(dest)
	authenticator := auth.NewAuthenticator(key, secret, filestore, nil, nil)
	downloader := dl.NewDownloader(os.Stdout)
	exitcode := app.Run("https://api.flickr.com", savecreds, authenticator,
		downloader, filestore, os.Stdout, os.Stderr)
	os.Exit(exitcode)
}

func requireEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		fmt.Fprintf(os.Stderr, "Please set the %s environment variable\n", name)
		os.Exit(1)
	}
	return value
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
