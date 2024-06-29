package repository

import (
	"basic-trade/common"
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type IFileRepository struct {
	cld *cloudinary.Cloudinary
}

// How to test:
// https://github.com/cloudinary/cloudinary-go/blob/main/api/admin/asset_acceptance_test.go
type FileRepository interface {
	UploadImage(ctx context.Context, publicID string, image *multipart.FileHeader) (string, error)
	DeleteImage(ctx context.Context, publicID string) error
}

func NewFileRepository(cld *cloudinary.Cloudinary) FileRepository {
	return &IFileRepository{cld: cld}
}

func (s *IFileRepository) UploadImage(ctx context.Context, publicID string, image *multipart.FileHeader) (string, error) {
	file, err := common.ConvertFileHeaderToReader(image)
	if err != nil {
		return "", err
	}

	result, err := s.cld.Upload.Upload(
		ctx,
		file,
		uploader.UploadParams{
			PublicID:       publicID,
			UniqueFilename: api.Bool(false),
			Overwrite:      api.Bool(true),
			Folder:         "products",
		},
	)
	if err != nil {
		return "", err
	}

	return result.SecureURL, nil
}

func (s *IFileRepository) DeleteImage(ctx context.Context, publicID string) error {
	_, err := s.cld.Upload.Destroy(
		ctx,
		uploader.DestroyParams{
			PublicID:   publicID,
			Invalidate: api.Bool(true),
		},
	)

	return err
}
