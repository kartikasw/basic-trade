package validation

import (
	"basic-trade/common"
	"fmt"
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
		fmt.Println("file size = ", fileHeader.Size)
		if fileHeader.Size > common.MaxFileSize {
			return false
		}

		ext := filepath.Ext(fileHeader.Filename)
		fmt.Println("file ext = ", ext)
		if _, allowed := AllowedImageExtensions[ext]; !allowed {
			return false
		}

		return true
	}

	return false
}
