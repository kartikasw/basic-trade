package validation

import (
	"basic-trade/common"
	"mime/multipart"
	"path/filepath"

	"github.com/go-playground/validator/v10"
)

var AllowedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".svg":  true,
}

var ValidImage validator.Func = func(fl validator.FieldLevel) bool {
	if fileHeader, ok := fl.Field().Interface().(multipart.FileHeader); ok {
		if fileHeader.Size > common.MAX_FILE_SIZE {
			return false
		}

		ext := filepath.Ext(fileHeader.Filename)
		if _, allowed := AllowedImageExtensions[ext]; !allowed {
			return false
		}

		return true
	}

	return false
}
