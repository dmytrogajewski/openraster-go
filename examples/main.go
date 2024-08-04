package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/dmytrogajewski/openraster-go/pkg/ora"
)

func colorEquals(c1, c2 color.Color) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

func main() {
	f, err := os.Open("assets/map.ora")

	if err != nil {
		panic(err)
	}

	stats, err := f.Stat()

	if err != nil {
		panic(err)
	}

	o := ora.NewOra()
	err = o.Load(f, stats.Size())

	if err != nil {
		panic(err)
	}

	for _, oi := range o.Children {
		bounds := oi.Image.Bounds()

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				c := oi.Image.At(x, y)
				if colorEquals(c, color.Black) {
					fmt.Print("█")
				} else {
					fmt.Print(" ")
				}
			}
			fmt.Println()
		}
		fmt.Printf("%v\n", oi.Name())
	}

}
