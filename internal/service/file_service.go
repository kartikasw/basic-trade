package service

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type FileService struct {
	cld *cloudinary.Cloudinary
}

type IFileService interface {
	UploadImage(image *multipart.FileHeader) (string, error)
}

func NewFileService(cld *cloudinary.Cloudinary) *FileService {
	return &FileService{cld: cld}
}

func (s *FileService) UploadImage(image *multipart.FileHeader) (string, error) {
	result, err := s.cld.Upload.Upload(
		context.Background(), 
		image.Filename, 
		uploader.UploadParams{},
	)
	if err != nil {
		return "", err
	}

	return result.SecureURL, nil
}
