package request

import "mime/multipart"

type CreateProductRequest struct {
	Name  string                `form:"name" binding:"required,max=100"`
	Image *multipart.FileHeader `form:"image" binding:"required,validImage"`
}

type UpdateProductRequest struct {
	Name  string                `form:"name" binding:"max=100"`
	Image *multipart.FileHeader `form:"image" binding:"validImage"`
}
