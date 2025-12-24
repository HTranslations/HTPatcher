package main

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"os"
	"strings"

	_ "embed"
)

//go:embed credits.png
var creditsPng []byte

var pngHeader = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52}

func AddCreditsToResource(path string, encryptionKey string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	isEncrypted := !strings.HasSuffix(path, ".png")
	header := []byte{}
	if isEncrypted {
		header, data, err = decryptPng(data)
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

	// Draw the credits image on top of the original image at the bottom-left corner
	creditsBounds := creditsImg.Bounds()
	offset := image.Pt(0, bounds.Max.Y-creditsBounds.Dy())
	draw.Draw(rgba, creditsBounds.Add(offset), creditsImg, creditsBounds.Min, draw.Over)

	// Encode the modified image
	var buf bytes.Buffer
	if err := png.Encode(&buf, rgba); err != nil {
		return err
	}

	var dataToWrite []byte
	if isEncrypted {
		dataToWrite = encryptPng(buf.Bytes(), header)
	} else {
		dataToWrite = buf.Bytes()
	}

	os.WriteFile(path, dataToWrite, 0644)
	return nil
}

func decryptPng(data []byte) ([]byte, []byte, error) {
	header := data[:32]
	data = append(pngHeader, data[32:]...)

	return header, data, nil
}

func encryptPng(data []byte, header []byte) []byte {
	return append(header, data[16:]...)
}
