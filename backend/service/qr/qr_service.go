package qr

import (
	"fmt"

	qrcode "github.com/skip2/go-qrcode"
)

type QRServiceI interface {
	GeneratePNG(shortURL string) ([]byte, error)
}

type qrService struct{}

func NewQRService() QRServiceI {
	return &qrService{}
}

func (s *qrService) GeneratePNG(shortURL string) ([]byte, error) {
	png, err := qrcode.Encode(shortURL, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}
	return png, nil
}
