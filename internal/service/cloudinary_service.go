package service

import (
	"fmt"
	"net/url"
	"os"
	"styleai-backend/pkg/utils"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/google/uuid"
)

type CloudinaryService struct {
	apiKey     string
	apiSecret  string
	cloudName  string
	cldManager *utils.CloudinaryManager
}

type SignatureResponse struct {
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
	APIKey    string `json:"apiKey"`
	CloudName string `json:"cloudName"`
	Folder    string `json:"folder"`
}

func NewCloudinaryService(cldManager *utils.CloudinaryManager) *CloudinaryService {
	return &CloudinaryService{
		apiKey:     os.Getenv("CLOUDINARY_API_KEY"),
		apiSecret:  os.Getenv("CLOUDINARY_API_SECRET"),
		cloudName:  os.Getenv("CLOUDINARY_CLOUD_NAME"),
		cldManager: cldManager,
	}
}

func (s *CloudinaryService) GenerateProductImageSignature() (*SignatureResponse, error) {
	timestamp := time.Now().Unix()

	uploadID := uuid.NewString()

	folder := fmt.Sprintf(
		"products/pending/%s",
		uploadID,
	)

	params := url.Values{}
	params.Add("timestamp", fmt.Sprintf("%d", timestamp))
	params.Add("folder", folder)

	signature, err := api.SignParameters(params, s.apiSecret)
	if err != nil {
		return nil, err
	}

	return &SignatureResponse{
		Timestamp: timestamp,
		Signature: signature,
		APIKey:    s.apiKey,
		CloudName: s.cloudName,
		Folder:    folder,
	}, nil
}

func (s *CloudinaryService) DeleteProductImageFolder(folder string) error {
	return s.cldManager.DeleteByPrefix(folder)
}

func (s *CloudinaryService) GenerateSignature(productID string, variantID string) (*SignatureResponse, error) {
	timestamp := time.Now().Unix()

	folder := fmt.Sprintf("products/%s/%s", productID, variantID)

	params := url.Values{}
	params.Add("timestamp", fmt.Sprintf("%d", timestamp))
	params.Add("folder", folder)

	signature, err := api.SignParameters(params, s.apiSecret)
	if err != nil {
		return nil, err
	}

	return &SignatureResponse{
		Timestamp: timestamp,
		Signature: signature,
		APIKey:    s.apiKey,
		CloudName: s.cloudName,
		Folder:    folder,
	}, nil
}

func (s *CloudinaryService) DeleteVariantImages(
	productID string,
	variantID string,
) error {
	prefix := fmt.Sprintf(
		"products/%s/%s",
		productID,
		variantID,
	)

	return s.cldManager.DeleteByPrefix(prefix)
}
