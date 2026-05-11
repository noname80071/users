package avatars

import (
	"bytes"
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

	bounds := srcImg.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width != height {
		return nil, fmt.Errorf("invalid skin dimensions: expected square image, got %dx%d", width, height)
	}

	// Размер головы 1/8 от размера всего скина
	headSize := width / 8
	headOffset := headSize

	cropRect := image.Rect(
		headOffset,          // min X
		headOffset,          // min Y
		headOffset+headSize, // max X
		headOffset+headSize, // max Y
	)

	// Проверяем, что область вырезания не выходит за границы
	if cropRect.Max.X > width || cropRect.Max.Y > height {
		return nil, fmt.Errorf("crop area is outside image bounds: rect %v, image size %dx%d", cropRect, width, height)
	}

	croppedImg := imaging.Crop(srcImg, cropRect)

	buf := new(bytes.Buffer)
	err = png.Encode(buf, croppedImg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %w", err)
	}

	return buf.Bytes(), nil
}
