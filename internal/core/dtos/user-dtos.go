package dtos

type CreateUserDto struct {
	Username string `json:"username" binding:"required"`
	FullName string `json:"full_name"`
	Email    string `json:"email" binding:"required, email"`
	Password string `json:"password" binding:"required, min=8, max=20"`
}

type GetUserByUsernameDto struct {
	Username string `json:"username" binding:"required"`
}

type GetAllAccountsForUserDto struct {
	Username string `json:"username" binding:"required"`
}
