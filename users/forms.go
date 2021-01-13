package users

/**
Data obtained POST requests and to be validated
*/
type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterForm struct {
	Username    string `json:"username" validate:"required,min=3,max=32"`
	DisplayName string `json:"displayName" validate:"max=32"`
	Email       string `json:"email" validate:"required,email,min=6,max=32"`
	Password    string `json:"password" validate:"required,min=3,max=32"`
}
