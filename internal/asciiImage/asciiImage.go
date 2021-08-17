package asciiImage

import (
	"bytes"
	"image"
	"image/color"
	"reflect"

	"github.com/nfnt/resize"
)

// Model is a struct that contains all the properties of the ascii image.
type Model struct {
	Image   image.Image
	Content string
	Height  int
	Width   int
}

// asciiString is the ascii string that will be used to represent the image.
var asciiString = "IMND8OZ$7I?+=~:,.."

// ScaleImage resizes an image to the given width and height.
func ScaleImage(img image.Image, w, height int) (image.Image, int, int) {
	img = resize.Resize(uint(w), uint(height), img, resize.Lanczos3)

	return img, w, height
}

// ConvertToAscii converts an image to ASCII.
func ConvertToAscii(img image.Image, w, h int) string {
	table := []byte(asciiString)
	buf := new(bytes.Buffer)

	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			g := color.GrayModel.Convert(img.At(j, i))
			y := reflect.ValueOf(g).FieldByName("Y").Uint()
			pos := int(y * 16 / 255)
			_ = buf.WriteByte(table[pos])
		}
		_ = buf.WriteByte('\n')
	}

	return buf.String()
}

// SetContent sets the content of the ascii image.
func (m *Model) SetContent(content string) {
	m.Content = content
}

// SetImage sets the image of the ascii image.
func (m *Model) SetImage(img image.Image) {
	m.Image = img
}

// View returns a string representation of the ascii image.
func (m Model) View() string {
	return m.Content
}
