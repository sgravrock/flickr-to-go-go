package auth

import (
	"fmt"
	"os/exec"
)

type consoleUiAdapter struct{}

func NewConsoleUiAdapter() UiAdapter {
	return consoleUiAdapter{}
}

func (a consoleUiAdapter) PromptForAccessCode(url string) (string, error) {
	if exec.Command("open", url).Run() == nil {
		fmt.Println("The Flickr login page should have opened " +
			"in your browser. Please log in.")
	} else {
		fmt.Printf("Open this in your browser: %s\n", url)
	}

	fmt.Println("After you've finished logging in, enter the " +
		"code from your browser:")

	verificationCode := ""
	_, err := fmt.Scanln(&verificationCode)
	return verificationCode, err
}
