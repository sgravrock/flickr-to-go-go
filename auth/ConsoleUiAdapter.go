package auth

import "fmt"

type consoleUiAdapter struct{}

func NewConsoleUiAdapter() UiAdapter {
	return consoleUiAdapter{}
}

func (a consoleUiAdapter) PromptForAccessCode(url string) (string, error) {
	fmt.Println("(1) Go to: " + url)
	fmt.Println("(2) Grant access. You should get back a verification code.")
	fmt.Println("(3) Enter that verification code here: ")

	verificationCode := ""
	_, err := fmt.Scanln(&verificationCode)
	return verificationCode, err
}
