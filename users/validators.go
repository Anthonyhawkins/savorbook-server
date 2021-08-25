package users

import (
	"github.com/go-playground/validator/v10"
)

type LoginValidator struct {
	Login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"login"`
	Model UserModel
}

type RegisterValidator struct {
	Registration struct {
		Username    string `json:"username" validate:"required,min=3,max=32"`
		DisplayName string `json:"displayName" validate:"max=32"`
		Email       string `json:"email" validate:"required,email,min=6,max=32"`
		Password    string `json:"password" validate:"required,min=3,max=32"`
	} `json:"registration"`
	Model UserModel
}

type UserValidator struct {
	User struct {
		Username    string `json:"username" validate:"required,min=3,max=32"`
		DisplayName string `json:"displayName" validate:"max=32"`
		Email       string `json:"email" validate:"required,email,min=6,max=32"`
	} `json:"user"`
	Model UserModel
}

type PasswordValidator struct {
	Password struct {
		Password string `json:"password" validate:"required,min=3,max=32"`
	} `json:"password"`
	Model UserModel
}

func NewLoginValidator() *LoginValidator {
	return &LoginValidator{}
}

func NewRegisterValidator() *RegisterValidator {
	return &RegisterValidator{}
}

func NewUserValidator() *UserValidator {
	return &UserValidator{}
}

func NewPasswordValidator() *PasswordValidator {
	return &PasswordValidator{}
}

func (v *RegisterValidator) Validate() ([]string, error) {
	var errors []string
	validate := validator.New()
	err := validate.Struct(v)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			message := err.Field() + " - " + err.Tag()
			errors = append(errors, message)
		}
	}
	return errors, err
}

func (v *LoginValidator) Validate() ([]string, error) {
	var errors []string
	validate := validator.New()
	err := validate.Struct(v)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			message := err.Field() + " - " + err.Tag()
			errors = append(errors, message)
		}
	}
	return errors, err
}

func (v *UserValidator) Validate() ([]string, error) {
	var errors []string
	validate := validator.New()
	err := validate.Struct(v)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			message := err.Field() + " - " + err.Tag()
			errors = append(errors, message)
		}
	}
	return errors, err
}

func (v *PasswordValidator) Validate() ([]string, error) {
	var errors []string
	validate := validator.New()
	err := validate.Struct(v)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			message := err.Field() + " - " + err.Tag()
			errors = append(errors, message)
		}
	}
	return errors, err
}

func (v *LoginValidator) BindModel() error {
	v.Model.Email = v.Login.Email
	return nil
}

func (v *RegisterValidator) BindModel() error {
	v.Model.Email = v.Registration.Email
	v.Model.Username = v.Registration.Username
	v.Model.DisplayName = v.Registration.DisplayName
	return nil
}

func (v *UserValidator) BindModel() error {
	v.Model.Email = v.User.Email
	v.Model.Username = v.User.Username
	v.Model.DisplayName = v.User.DisplayName
	return nil
}
