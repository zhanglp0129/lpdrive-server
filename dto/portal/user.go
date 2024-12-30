package portaldto

type UserLoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserChangePasswordDTO struct {
	ID          int64  `json:"-"`
	OldPassword string `json:"oldPassword"`
	Password    string `json:"password"`
}
