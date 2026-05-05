package validators

import (
	"fmt"
	"net/mail"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ValidateEmail validates email format
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// ValidatePassword validates password requirements
func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	if len(password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}

	if len(password) > 128 {
		return fmt.Errorf("password must be less than 128 characters")
	}

	return nil
}

// ValidateObjectID validates if a string is a valid MongoDB ObjectID
func ValidateObjectID(id string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	_, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format")
	}

	return nil
}

// ValidateNonEmpty validates that a string is not empty
func ValidateNonEmpty(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s cannot be empty", fieldName)
	}
	return nil
}

// ValidateStringLength validates string length constraints
func ValidateStringLength(value, fieldName string, minLength, maxLength int) error {
	length := len(strings.TrimSpace(value))

	if length < minLength {
		return fmt.Errorf("%s must be at least %d characters", fieldName, minLength)
	}

	if length > maxLength {
		return fmt.Errorf("%s must be less than %d characters", fieldName, maxLength)
	}

	return nil
}
