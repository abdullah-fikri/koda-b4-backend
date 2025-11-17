package lib

import (
	"backend/config"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)


func UploadImage(file multipart.File) (string, error) {
	cld, ctx := config.CloudinaryInit()

	if cld == nil {
		return UploadLocal(file)
	}

	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{})
	if err != nil {
		return UploadLocal(file)
	}

	return uploadResult.SecureURL, nil
}

func UploadLocal(src multipart.File) (string, error) {
	path := "uploads/"
	os.MkdirAll(path, 0755)

	filename := fmt.Sprintf("%d.jpg", time.Now().Unix())
	fullpath := path + filename

	dst, err := os.Create(fullpath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}
	return "/static/" + filename, nil
}
