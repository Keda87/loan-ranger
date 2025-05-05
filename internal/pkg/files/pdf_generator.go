package files

import (
	"bytes"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

func GeneratePDFBuffer(text string) (*bytes.Buffer, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	x, y := 10.0, 10.0
	for _, line := range strings.Split(text, "\n") {
		pdf.Text(x, y, line)
		y += 10
	}

	buf := new(bytes.Buffer)
	err := pdf.Output(buf)
	return buf, err
}
