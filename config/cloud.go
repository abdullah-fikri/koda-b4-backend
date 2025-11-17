package config

import (
	"context"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

func CloudinaryInit() (*cloudinary.Cloudinary, context.Context) {
	url := os.Getenv("CLOUDINARY_URL")
	if url == "" {
		return nil, context.Background() 
	}

	cld, err := cloudinary.NewFromURL(url)
	if err != nil {
		return nil, context.Background() 
	}

	return cld, context.Background()
}
