package dtos

type CreateEntryReq struct {
	AccountID int64   `json:"account_id" binding:"required,min=0"`
	Amount    float64 `json:"amount" binding:"required,min=0"`
}

type GetAllEntriesOfAccountReq struct {
	AccountID int64 `uri:"account_id" binding:"required,min=0"`
}

type GetEntryByIdReq struct {
	AccountID int64 `form:"account_id" binding:"required,min=0"`
	EntryID   int64 `form:"entry_id" binding:"required,min=0"`
}

type GetEntriesInPage struct {
	AccountID int64 `uri:"account_id" binding:"required,min=0"`
	Limit     int16 `form:"limit" binding:"required,max=10,min=1"`
	Offset    int16 `form:"offset" binding:"required,min=1"`
}
