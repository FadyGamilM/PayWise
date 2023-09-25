package dtos

type CreateTransferReq struct {
	ToAccountID   int64   `json:"to_account_id" binding:"required,min=0"`
	FromAccountID int64   `json:"from_account_id" binding:"required,min=0"`
	Amount        float64 `json:"amount" binding:"required,min=0"`
}

type GetTransferByIdReq struct {
	TransferID int64 `uri:"id" binding:"required,min=0"`
}

type GetTransfersFromAccountReq struct {
	FromAccountID int64 `uri:"from_account_id" binding:"required,min=0"`
	Limit         int16 `form:"limit" binding:"required,min=1,max=10"`
	Offset        int16 `form:"offset" binding:"required,min=0"`
}

type GetTransfersToAccountReq struct {
	ToAccountID int64 `uri:"to_account_id" binding:"required,min=0"`
	Limit       int16 `form:"limit" binding:"required,min=1,max=10"`
	Offset      int16 `form:"offset" binding:"required,min=0"`
}

type GetTransfersBetweenTwoAccountsReq struct {
	FromAccountID int64 `uri:"from_account_id" binding:"required,min=0"`
	ToAccountID   int64 `uri:"to_account_id" binding:"required,min=0"`
	Limit         int16 `form:"limit" binding:"required,min=1,max=10"`
	Offset        int16 `form:"offset" binding:"required,min=0"`
}
