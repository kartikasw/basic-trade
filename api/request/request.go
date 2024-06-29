package request

type GetDataByUUIDRequest struct {
	UUID string `uri:"uuid" binding:"required,validUUID"`
}

type PaginationRequest struct {
	Keyword string `form:"keyword"`
	Limit   int32  `form:"limit" binding:"required,min=5,max=10"`
	Offset  int32  `form:"offset" binding:"gt=-1"`
}
