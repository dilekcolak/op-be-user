package domain

import (
	"strings"
	"unicode"

	me "github.com/octoposprime/op-be-user/internal/domain/model/entity"
	mo "github.com/octoposprime/op-be-user/internal/domain/model/object"
)

// Service struct for user-related services
type Service struct {
}

// NewService creates a new instance of Service
func NewService() *Service {
	return &Service{}
}

// ValidateUser validates the user based on certain criteria
func (s *Service) ValidateUser(user *me.User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	return nil
}

// ValidatePassword checks the validity of a user's password
func (s *Service) ValidatePassword(userPassword *me.UserPassword) error {
	if err := userPassword.Validate(); err != nil {
		return err
	}
	return nil
}

// CheckUserNameRules checks the rules for usernames
func (s *Service) CheckUserNameRules(user *me.User) error {
	if user.UserName == "" {
		return mo.ErrorUserUsernameIsEmpty
	}
	if len(user.UserName) < 8 {
		return mo.ErrorUserUsernameIsTooShort
	}
	if len(user.UserName) > 20 {
		return mo.ErrorUserUsernameIsTooLong
	}
	if strings.Contains(user.UserName, " ") {
		return mo.ErrorUserUsernameIsNotValid
	}
	if !isValidUsername(user.UserName) {
		return mo.ErrorUserUsernameContainsSpecialChar
	}
	return nil
}

func isValidUsername(username string) bool {
	var hasLetterOrDigit bool
	for _, ch := range username {
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
			hasLetterOrDigit = true
		}
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' && ch != '.' {
			return false
		}
	}
	return hasLetterOrDigit // Make sure there is at least one alphanumeric character
}

// CheckEmailRules validates the email structure
func (s *Service) CheckEmailRules(user *me.User) error {
	if user.Email == "" {
		return mo.ErrorUserEmailIsEmpty
	}
	if !strings.Contains(user.Email, "@") || !strings.Contains(user.Email, ".") {
		return mo.ErrorUserEmailIsNotValid
	}
	if strings.Contains(user.Email, " ") {
		return mo.ErrorUserEmailIsNotValid
	}
	return nil
}

// CheckPasswordRules checks the rules for passwords
func (s *Service) CheckPasswordRules(userPassword *me.UserPassword) error {
	if userPassword.Password == "" {
		return mo.ErrorUserPasswordIsEmpty
	}
	if len(userPassword.Password) < 8 {
		return mo.ErrorUserPasswordIsTooShort
	}
	if len(userPassword.Password) > 20 {
		return mo.ErrorUserPasswordIsTooLong
	}
	if strings.Contains(userPassword.Password, " ") {
		return mo.ErrorUserPasswordIsNotValid
	}
	return nil
}
