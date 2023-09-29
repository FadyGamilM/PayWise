package dtos

type GetAccByOwnerNameReq struct {
	OwnerName string `form:"owner" binding:"required"`
}

type GetAccByIdReq struct {
	ID int64 `uri:"id" binding:"required,min=0"`
}

type CreateAccReq struct {
	// OwnerName string `json:"owner_name" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=EUR USD"`
}

type PaginateAccountsReq struct {
	Limit  int16 `form:"limit" binding:"required,max=10,min=1"`
	Offset int16 `form:"offset" binding:"required,min=1"`
}

type UpdateAccountReq struct {
	ID      int64   `uri:"id" binding:"required,min=0"`
	Balance float64 `json:"balance" binding:"rqeuired,min=0"`
}

type UpdateAccountByOwnerNameReq struct {
	OwnerName string  `json:"owner_name" binding:"rqeuired"`
	Balance   float64 `json:"balance" binding:"rqeuired,min=0"`
}

type DeleteAccountReq struct {
	ID int64 `uri:"id" binding:"required,min=0"`
}

type DeleteAccountByOwnerNameReq struct {
	OwnerName string `json:"owner_name" binding:"rqeuired"`
}
