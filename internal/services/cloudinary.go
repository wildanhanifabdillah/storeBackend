package services

import (
	"context"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadImage(file multipart.File, folder string) (string, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return "", err
	}

	res, err := cld.Upload.Upload(
		context.Background(),
		file,
		uploader.UploadParams{
			Folder: folder,
		},
	)

	if err != nil {
		return "", err
	}

	return res.SecureURL, nil
}
