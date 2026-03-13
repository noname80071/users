package avatars

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"

	"github.com/disintegration/imaging"
)

func CropAvatar(data []byte) ([]byte, error) {
	// Декодируем изображение
	srcImg, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	cropRect := image.Rect(8, 8, 16, 16)

	if cropRect.Max.X > srcImg.Bounds().Max.X || cropRect.Max.Y > srcImg.Bounds().Max.Y {
		return nil, errors.New("crop area is outside image bounds")
	}

	croppedImg := imaging.Crop(srcImg, cropRect)

	buf := new(bytes.Buffer)
	err = png.Encode(buf, croppedImg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %w", err)
	}

	return buf.Bytes(), nil
}
