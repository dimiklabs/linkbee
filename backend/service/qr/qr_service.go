package qr

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

// ErrorCorrectionLevel maps the single-char string to the library constant.
var errorCorrectionMap = map[string]qrcode.RecoveryLevel{
	"L": qrcode.Low,
	"M": qrcode.Medium,
	"Q": qrcode.High,    // go-qrcode uses High for Q
	"H": qrcode.Highest, // go-qrcode uses Highest for H
}

// QROptions controls the appearance of the generated QR code.
type QROptions struct {
	// ForegroundHex and BackgroundHex are six-hex-digit colors without the '#' prefix.
	ForegroundHex string // default: "000000" (black)
	BackgroundHex string // default: "ffffff" (white)
	// Size is the output image width/height in pixels (clamped to 64–1024).
	Size int // default: 256
	// ErrorCorrection is one of "L", "M", "Q", "H" (default "M").
	ErrorCorrection string
}

func (o *QROptions) normalize() {
	if o.ForegroundHex == "" {
		o.ForegroundHex = "000000"
	}
	if o.BackgroundHex == "" {
		o.BackgroundHex = "ffffff"
	}
	if o.Size < 64 {
		o.Size = 64
	} else if o.Size > 1024 {
		o.Size = 1024
	}
	o.ErrorCorrection = strings.ToUpper(o.ErrorCorrection)
	if _, ok := errorCorrectionMap[o.ErrorCorrection]; !ok {
		o.ErrorCorrection = "M"
	}
}

type QRServiceI interface {
	GeneratePNG(shortURL string) ([]byte, error)
	GenerateCustomPNG(shortURL string, opts QROptions) ([]byte, error)
}

type qrService struct{}

func NewQRService() QRServiceI {
	return &qrService{}
}

func (s *qrService) GeneratePNG(shortURL string) ([]byte, error) {
	return s.GenerateCustomPNG(shortURL, QROptions{})
}

func (s *qrService) GenerateCustomPNG(shortURL string, opts QROptions) ([]byte, error) {
	opts.normalize()

	level := errorCorrectionMap[opts.ErrorCorrection]

	qr, err := qrcode.New(shortURL, level)
	if err != nil {
		return nil, fmt.Errorf("failed to create QR code: %w", err)
	}

	fg, err := hexToRGBA(opts.ForegroundHex)
	if err != nil {
		fg = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	}
	bg, err := hexToRGBA(opts.BackgroundHex)
	if err != nil {
		bg = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}

	qr.ForegroundColor = fg
	qr.BackgroundColor = bg

	png, err := qr.PNG(opts.Size)
	if err != nil {
		return nil, fmt.Errorf("failed to render QR code PNG: %w", err)
	}
	return png, nil
}

// hexToRGBA parses a 6-char hex string (without '#') into color.RGBA.
func hexToRGBA(hex string) (color.RGBA, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return color.RGBA{}, fmt.Errorf("invalid hex color: %q", hex)
	}
	r, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}
	g, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}
	b, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}, nil
}
