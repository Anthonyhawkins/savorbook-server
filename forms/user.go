package forms

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterForm struct {
	Username        string `json:"username" validate:"required,min=3,max=32"`
	Email           string `json:validate:"required,email,min=6,max=32"`
	Password        string `json:"password" validate:"required,min=3,max=32"`
	PasswordConfirm string `json:"passwordConfirm" validate:"required,min=3,max=32"`
}
