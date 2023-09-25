package dtos

type LoginReq struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRes struct {
	AccessToken string      `json:"access_token"`
	User        *UserResDto `json:"user"`
}
