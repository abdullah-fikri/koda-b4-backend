package config

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

func CloudinaryInit() (*cloudinary.Cloudinary, context.Context, error) {
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		return nil, nil, fmt.Errorf("CLOUDINARY_URL not found in environment")
	}

	ctx := context.Background()
	return cld, ctx, nil
}
