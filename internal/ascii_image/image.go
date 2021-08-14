package ascii_image

import (
	"bytes"
	"image"
	"image/color"
	"reflect"

	"github.com/nfnt/resize"
)

type Model struct {
	Image  image.Image
	Height int
	Width  int
}

var ASCIISTR = "IMND8OZ$7I?+=~:,.."

// scaleImage resizes an image to the given width and height.
func (m Model) scaleImage() image.Image {
	img := resize.Resize(uint(m.Width), uint(m.Height), m.Image, resize.Lanczos3)

	return img
}

// convertToAscii converts an image to ASCII.
func (m Model) convertToAscii(scaledImage image.Image) string {
	table := []byte(ASCIISTR)
	buf := new(bytes.Buffer)

	for i := 0; i < m.Height; i++ {
		for j := 0; j < m.Width; j++ {
			g := color.GrayModel.Convert(m.Image.At(j, i))
			y := reflect.ValueOf(g).FieldByName("Y").Uint()
			pos := int(y * 16 / 255)
			_ = buf.WriteByte(table[pos])
		}
		_ = buf.WriteByte('\n')
	}

	return buf.String()
}

func (m *Model) SetContent(img image.Image) {
	m.Image = img
}

func (m *Model) SetSize(width, height int) {
	m.Width = width
	m.Height = height
}

func (m Model) View() string {
	return m.convertToAscii(m.scaleImage())
}
