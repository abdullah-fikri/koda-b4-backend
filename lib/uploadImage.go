package lib

import (
	"backend/config"
	"io"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadImage(file io.Reader) (string, error) {
	cld, ctx, err := config.CloudinaryInit()

	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
