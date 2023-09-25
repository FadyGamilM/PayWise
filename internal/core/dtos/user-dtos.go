package dtos

type CreateUserDto struct {
	Username string `json:"username" binding:"required,alphanum"`
	FullName string `json:"full_name" binding:"required,alpha"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

type GetUserByUsernameDto struct {
	Username string `json:"username" binding:"required"`
}

type GetAllAccountsForUserDto struct {
	Username string `json:"username" binding:"required"`
}

type UserResDto struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}
