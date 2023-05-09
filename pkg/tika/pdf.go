package tika

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
)

func isSameSentence(t1, t2 pdf.Text) bool {
	if t1.Y == t2.Y || t1.X == t2.X {
		return true
	}
	return false
}

func readPdf2(path string) (string, error) {
	f, r, err := pdf.Open(path)
	// remember close file
	defer f.Close()
	if err != nil {
		return "", err
	}
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		var lastTextStyle pdf.Text
		texts := p.Content().Text
		for _, text := range texts {
			if isSameSentence(text, lastTextStyle) {
				lastTextStyle.S = lastTextStyle.S + text.S
			} else {
				fmt.Printf("Font: %s, Font-size: %f, x: %f, y: %f, content: %s \n", lastTextStyle.Font, lastTextStyle.FontSize, lastTextStyle.X, lastTextStyle.Y, lastTextStyle.S)
				lastTextStyle = text
			}
		}
	}
	return "", nil
}

func readPdf(path string) (string, error) {
	f, r, err := pdf.Open("uploadfile/" + path)
	// remember close file
	defer f.Close()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}

type FieldRect struct {
	rect       pdf.Rect
	texts      []pdf.Text
	lastPos    *pdf.Point
	resultText string
}

func ReadPd3f(filename string) ([]string, error) {
	f, r, err := pdf.Open("uploadfile/" + filename)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	totalPage := r.NumPage()
	var posttext string

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		c := p.Content()
		var pctexts []pdf.Text
		for _, t := range c.Text {

			pctexts = append(pctexts, t)

		}
		// these values can also be derived from font size to gain
		// even more robustness
		const NEWLINE_TOLERATION = 2
		// unfortunately the pdf you sent does not have proper font
		// width information, so this is the best we can get without
		// inferring width information from the glyph shape itself.
		const SPACE_TOLERATION = 11
		var resultText string
		var lastpos *pdf.Point
		// sort.Slice(pctexts, func(i, j int) bool {
		// 	deltaY := pctexts[i].Y - pctexts[j].Y
		// 	if math.Abs(deltaY) < math.Min(pctexts[j].FontSize, pctexts[i].FontSize) { // tolerate some vertical deviation
		// 		return pctexts[i].X < pctexts[j].X // on the same line
		// 	}
		// 	return deltaY > 0 // not on the same line
		// })
		for _, f := range pctexts {
			if lastpos != nil {
				if lastpos.Y-f.Y > f.FontSize { // new line
					resultText += "\n"
				}
				if f.X-lastpos.X > SPACE_TOLERATION { // space
					resultText += " "
				}
			}
			resultText += f.S
			lastpos = &pdf.Point{X: f.X, Y: f.Y}

		}
		posttext += resultText
	}
	// fmt.Print(posttext)

	list1 := strings.Split(posttext, "\n")
	var textbodylist []string
	var textbody string

	for i, v := range list1 {
		textbody += v
		if len(textbody) > 600 || i == len(list1)-1 {
			textbodylist = append(textbodylist, filename+"\n"+textbody)
			textbody = ""
		}
	}
	// file, err := os.Create("tmpfile/" + path + "output.txt")
	// if err != nil {
	// 	fmt.Println("Error creating file:", err)
	// 	return nil, err
	// }
	// defer file.Close()

	// writer := bufio.NewWriter(file)

	// for _, v := range textbodylist {
	// 	writer.WriteString(v)
	// 	writer.WriteString("\n------------------\n")
	// }

	// writer.Flush()
	return textbodylist, nil
}
