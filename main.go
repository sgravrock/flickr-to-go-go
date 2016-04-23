package main

import (
	"fmt"

	"os"

	"log"

	"github.com/mrjones/oauth"
)

func main() {
	consumerKey := os.Getenv("FLICKR_API_KEY")
	consumerSecret := os.Getenv("FLICKR_API_SECRET")

	c := oauth.NewConsumer(
		consumerKey,
		consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://www.flickr.com/services/oauth/request_token",
			AuthorizeTokenUrl: "https://www.flickr.com/services/oauth/authorize",
			AccessTokenUrl:    "https://www.flickr.com/services/oauth/access_token",
		})
	c.AdditionalAuthorizationUrlParams = map[string]string{
		"perms": "read",
	}
	c.Debug(true)

	requestToken, url, err := c.GetRequestTokenAndUrl("oob")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Println("(1) Go to: " + url)
	fmt.Println("(2) Grant access, you should get back a verification code.")
	fmt.Println("(3) Enter that verification code here: ")

	verificationCode := ""
	fmt.Scanln(&verificationCode)

	accessToken, err := c.AuthorizeToken(requestToken, verificationCode)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Printf("Got access code: %s\n", accessToken)
}
