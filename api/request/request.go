package request

type GetDataByUUIDRequest struct {
	UUID string `uri:"uuid" binding:"required,validUUID"`
}

type PaginationRequest struct {
	Limit  int32 `form:"limit" binding:"required,min=5,max=10"`
	Offset int32 `form:"offset" binding:"gt=-1"`
}

type SearchRequest struct {
	Keyword string `form:"keyword" binding:"required"`
	Limit   int32  `form:"limit" binding:"required,min=5,max=10"`
	Offset  int32  `form:"offset" binding:"required,gt=-1"`
}
