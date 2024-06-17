package request

type CreateVariantRequest struct {
	VariantName string `form:"variant_name" binding:"required,max=100"`
	Quantity    int32  `form:"quantity" binding:"required"`
	ProductID   string `form:"product_id" binding:"required,validUUID"`
}

type UpdateVariantRequest struct {
	VariantName string `form:"variant_name" binding:"max=100"`
	Quantity    int32  `form:"quantity" binding:"required"`
}