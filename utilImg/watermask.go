package utilImg

import (
	"github.com/golang/freetype"
	"image"
	"image/draw"
	"io/ioutil"
)

func WaterMarkString(dst draw.Image, markString string, font string, fontSize float64, x int, y int) error {
	data, err := ioutil.ReadFile(font)
	if err != nil {
		return err
	}
	f, err := freetype.ParseFont(data)
	if err != nil {
		return err
	}

	draw.Draw(dst, dst.Bounds(), image.White, image.ZP, draw.Src)

	c := freetype.NewContext()
	c.SetDst(dst)
	c.SetClip(dst.Bounds())
	c.SetSrc(image.Black)
	c.SetFont(f)
	c.SetFontSize(fontSize)

	_, err = c.DrawString(markString, freetype.Pt(x, y))
	if err != nil {
		return err
	}
	return nil
}
