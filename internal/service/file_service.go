package service

import (
	"basic-trade/common"
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type FileService struct {
	cld *cloudinary.Cloudinary
}

type IFileService interface {
	UploadImage(publicID string, image *multipart.FileHeader) (string, error)
}

func NewFileService(cld *cloudinary.Cloudinary) *FileService {
	return &FileService{cld: cld}
}

func (s *FileService) UploadImage(publicID string, image *multipart.FileHeader) (string, error) {
	file, err := common.ConvertFileHeaderToReader(image)
	if err != nil {
		return "", err
	}

	result, err := s.cld.Upload.Upload(
		context.Background(),
		file,
		uploader.UploadParams{
			PublicID:       publicID,
			UniqueFilename: api.Bool(false),
			Overwrite:      api.Bool(true),
		},
	)
	if err != nil {
		return "", err
	}

	return result.SecureURL, nil
}
