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
	if len(os.Args) == 1 {
		r := image.Rect(0, 0, 32, 32)
		p := image.NewNRGBA(r)
		draw.Draw(p, r, image.White, image.Point{}, 0)
		err := png.Encode(os.Stdout, p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "png encode oops: %v\n", err)
			os.Exit(1)
		}
		err = pngplus.EncodeBinary(os.Stdout, []byte("hello"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "binary encode oops: %v\n", err)
			os.Exit(1)
		}
		err = pngplus.EncodeBinary(os.Stdout, []byte("hello again"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "binary encode oops: %v\n", err)
			os.Exit(1)
		}

	} else if len(os.Args) == 2 {
		_, err := png.Decode(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "png decode oops: %v\n", err)
			os.Exit(1)
		}
		b, err := pngplus.DecodeBinary(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "binary decode oops: %v\n", err)
			os.Exit(1)
		}
		os.Stdout.Write(b)
		os.Stdout.WriteString("\n")
	} else {
		_, err := png.Decode(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "png decode seek oops: %v\n", err)
			os.Exit(1)
		}
		for {
			b, err := pngplus.DecodeBinary(os.Stdin)
			if err == pngplus.ErrNotBinary {
				continue
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "binary decode seek oops: %v\n", err)
				os.Exit(1)
			}
			os.Stdout.Write(b)
			os.Stdout.WriteString("\n")
		}
	}
}
