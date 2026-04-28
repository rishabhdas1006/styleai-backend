package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
)

type CloudinaryManager struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryManager() (*CloudinaryManager, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return nil, err
	}

	return &CloudinaryManager{cld: cld}, nil
}

func (c *CloudinaryManager) DeleteByPrefix(prefix string) error {
	_, err := c.cld.Admin.DeleteAssetsByPrefix(context.Background(), admin.DeleteAssetsByPrefixParams{
		Prefix:    []string{prefix},
		AssetType: "image",
	})
	if err != nil {
		return fmt.Errorf("cloudinary cleanup failed: %w", err)
	}
	return nil
}
