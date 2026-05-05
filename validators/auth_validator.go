package validators

import (
	"fmt"
)

// RegisterValidatorImpl implements handlers.RegisterValidator
type RegisterValidatorImpl struct {
}

// ValidateRegisterRequest validates register request with format validation only
func (v *RegisterValidatorImpl) ValidateRegisterRequest(email, password, name string) error {
	// Format validation only - business logic check is handled by database unique constraint
	if err := ValidateEmail(email); err != nil {
		return err
	}

	if err := ValidatePassword(password); err != nil {
		return err
	}

	if err := ValidateNonEmpty(name, "name"); err != nil {
		return err
	}

	if err := ValidateStringLength(name, "name", 1, 50); err != nil {
		return err
	}

	return nil
}

// LoginValidatorImpl implements handlers.LoginValidator
type LoginValidatorImpl struct{}

// ValidateLoginRequest validates login request
func (v *LoginValidatorImpl) ValidateLoginRequest(email, password string) error {
	if err := ValidateEmail(email); err != nil {
		return err
	}

	if err := ValidatePassword(password); err != nil {
		return err
	}

	return nil
}

// CommentValidatorImpl implements handlers.CommentValidator
type CommentValidatorImpl struct{}

// ValidateCreateCommentRequest validates comment creation request
func (v *CommentValidatorImpl) ValidateCreateCommentRequest(content string) error {
	if err := ValidateNonEmpty(content, "content"); err != nil {
		return err
	}

	if err := ValidateStringLength(content, "content", 1, 500); err != nil {
		return err
	}

	return nil
}

// FavoriteValidatorImpl implements handlers.FavoriteValidator
type FavoriteValidatorImpl struct{}

// ValidateCreateFavoriteRequest validates favorite creation request
func (v *FavoriteValidatorImpl) ValidateCreateFavoriteRequest(cocktailID string) error {
	if err := ValidateObjectID(cocktailID); err != nil {
		return fmt.Errorf("invalid cocktailId: %w", err)
	}

	return nil
}

// ValidateDeleteFavoriteRequest validates favorite deletion request
func (v *FavoriteValidatorImpl) ValidateDeleteFavoriteRequest(cocktailID string) error {
	if err := ValidateObjectID(cocktailID); err != nil {
		return fmt.Errorf("invalid cocktailId: %w", err)
	}

	return nil
}
