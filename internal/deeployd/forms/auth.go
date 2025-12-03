package forms

import (
	"github.com/deeploy-sh/deeploy/internal/shared/utils"
)

type RegisterForm struct {
	Email           string
	Password        string
	PasswordConfirm string
}

type RegisterErrors struct {
	Email           string
	Password        string
	PasswordConfirm string
	General         string
}

func (f *RegisterForm) Validate() RegisterErrors {
	var errors RegisterErrors
	if !utils.IsEmailValid(f.Email) {
		errors.Email = "Not a valid email"
	}
	if f.Email == "" {
		errors.Email = "Email is required"
	}
	if f.Password == "" {
		errors.Password = "Password is required"
	}
	if f.PasswordConfirm == "" {
		errors.PasswordConfirm = "Confirm your password"
	}
	if f.Password != f.PasswordConfirm {
		errors.Password = "Passwords do not match"
		errors.PasswordConfirm = "Passwords do not match"
	}
	return errors
}

func (e *RegisterErrors) HasErrors() bool {
	return e.Email != "" || e.Password != "" || e.PasswordConfirm != "" || e.General != ""
}

type LoginForm struct {
	Email    string
	Password string
}

type LoginErrors struct {
	Email    string
	Password string
	General  string
}

func (f *LoginForm) Validate() LoginErrors {
	var errors LoginErrors
	if !utils.IsEmailValid(f.Email) {
		errors.Email = "Not a valid email"
	}
	if f.Email == "" {
		errors.Email = "Email is required"
	}
	if f.Password == "" {
		errors.Password = "Password is required"
	}
	return errors
}

func (e *LoginErrors) HasErrors() bool {
	return e.Email != "" || e.Password != "" || e.General != ""
}
