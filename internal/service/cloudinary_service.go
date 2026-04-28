package service

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api"
)

type CloudinaryService struct {
	apiKey    string
	apiSecret string
	cloudName string
}

type SignatureResponse struct {
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
	APIKey    string `json:"apiKey"`
	CloudName string `json:"cloudName"`
	Folder    string `json:"folder"`
}

func NewCloudinaryService() *CloudinaryService {
	return &CloudinaryService{
		apiKey:    os.Getenv("CLOUDINARY_API_KEY"),
		apiSecret: os.Getenv("CLOUDINARY_API_SECRET"),
		cloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
	}
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
