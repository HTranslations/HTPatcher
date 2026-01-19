package patcher

import (
	"bytes"
	_ "embed"
	"errors"
	"htpatcher/internal/util"
	"image"
	"image/draw"
	"image/png"
	"os"
	"strings"
)

//go:embed credits.png
var creditsPng []byte

// CreditsPatcher handles adding credits to images
type CreditsPatcher struct{}

// NewCreditsPatcher creates a new credits patcher
func NewCreditsPatcher() *CreditsPatcher {
	return &CreditsPatcher{}
}

// AddCreditsToResource adds credits overlay to a game image
func (c *CreditsPatcher) AddCreditsToResource(path string, encryptionKey string, creditsLocation string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	isEncrypted := !strings.HasSuffix(path, ".png")
	header := []byte{}
	if isEncrypted {
		header, data, err = util.DecryptPng(data)
		if err != nil {
			return err
		}
	}

	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return err
	}

	// Decode the credits image
	creditsImg, err := png.Decode(bytes.NewReader(creditsPng))
	if err != nil {
		return err
	}

	// Create a new RGBA image based on the original
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	// Draw the credits image on top of the original image
	creditsBounds := creditsImg.Bounds()
	var offset image.Point
	switch creditsLocation {
	case "bottom_left":
		offset = image.Pt(0, bounds.Max.Y-creditsBounds.Dy())
	case "bottom_right":
		offset = image.Pt(bounds.Max.X-creditsBounds.Dx(), bounds.Max.Y-creditsBounds.Dy())
	case "top_left":
		offset = image.Pt(0, 0)
	case "top_right":
		offset = image.Pt(bounds.Max.X-creditsBounds.Dx(), 0)
	default:
		return errors.New("invalid credits location")
	}
	draw.Draw(rgba, creditsBounds.Add(offset), creditsImg, creditsBounds.Min, draw.Over)

	// Encode the modified image
	var buf bytes.Buffer
	if err := png.Encode(&buf, rgba); err != nil {
		return err
	}

	var dataToWrite []byte
	if isEncrypted {
		dataToWrite = util.EncryptPng(buf.Bytes(), header)
	} else {
		dataToWrite = buf.Bytes()
	}

	return os.WriteFile(path, dataToWrite, 0644)
}




