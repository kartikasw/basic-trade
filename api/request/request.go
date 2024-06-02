package request

import "github.com/google/uuid"

type GetDataByUUIDRequest struct {
	UUID uuid.UUID `uri:"uuid" binding:"required,uuid4_rfc4122"`
}

type PaginationRequest struct {
	Limit  int32 `form:"limit" binding:"required,min=1"`
	Offset int32 `form:"offset" binding:"required,min=5,max=10"`
}

type SearchRequest struct {
	Keyword string `form:"keyword" binding:"required"`
	Limit   int32  `form:"limit" binding:"required,min=1"`
	Offset  int32  `form:"offset" binding:"required,min=5,max=10"`
}
