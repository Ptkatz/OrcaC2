package screenopt

import (
	"github.com/vova616/screenshot"
	"image"
)

//Screenshot - Retrieve the screenshot of the active displays
func Screenshot() (*image.RGBA, error) {
	img, err := screenshot.CaptureScreen()
	if err != nil {
		return nil, err
	}
	return img, err
}
