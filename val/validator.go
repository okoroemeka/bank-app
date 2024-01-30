package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullname = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(value string, minLength, maxLength int) error {
	valueLength := len(value)
	if valueLength < minLength || valueLength > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}
	return nil
}

func ValidateUsername(username string) error {
	if err := ValidateString(username, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(username) {
		return fmt.Errorf("must contain only lowercase alphanumeric characters and underscore")
	}
	return nil
}

func ValidateFullName(fullName string) error {
	if err := ValidateString(fullName, 3, 100); err != nil {
		return err
	}
	if !isValidFullname(fullName) {
		return fmt.Errorf("must contain only letters or spaces")
	}
	return nil
}

func ValidatePassword(password string) error {
	return ValidateString(password, 6, 100)
}

func ValidateEmail(email string) error {
	if err := ValidateString(email, 5, 100); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("must be a valid email address")
	}
	return nil
}

func ValidateVerifyEmailId(arg int64) error {
	if arg <= 0 {
		return fmt.Errorf("must be a positive integer")
	}
	return nil
}
