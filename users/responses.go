package users

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	UserID      uint   `json:"userId"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}

type UserResponse struct {
	UserID      uint   `json:"userId"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Bio         string `gorm:"column:bio" json:"bio"`
}

func (r *LoginResponse) SerializeLogin(model *UserModel, accessToken string) {
	r.AccessToken = accessToken
	r.Username = model.Username
	r.UserID = model.ID
	r.DisplayName = model.DisplayName
}

func (r *UserResponse) SerializeUser(model *UserModel) {
	r.Username = model.Username
	r.UserID = model.ID
	r.DisplayName = model.DisplayName
	r.Bio = model.Bio
}
