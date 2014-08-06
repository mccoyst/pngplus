// +build ignore
// © 2014 Steve McCoy under the MIT license. See LICENSE for details.

package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/mccoyst/pngplus"
)

func main() {
	r := image.Rect(0, 0, 32, 32)
	p := image.NewNRGBA(r)
	draw.Draw(p, r, image.White, image.Point{}, 0)
	err := png.Encode(os.Stdout, p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oops: %v\n", err)
		os.Exit(1)
	}
	err = pngplus.EncodeBinary(os.Stdout, []byte("hello"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "oops: %v\n", err)
		os.Exit(1)
	}
	err = pngplus.EncodeITXT(os.Stdout, false, "name", "en-US", "there")
	if err != nil {
		fmt.Fprintf(os.Stderr, "oops: %v\n", err)
		os.Exit(1)
	}
}
