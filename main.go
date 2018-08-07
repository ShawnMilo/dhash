package main

// This is a simplified function copy based on https://github.com/andybalholm/dhash/blob/master/hash.go

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
)

var maxSize = int64(1024 * 1024 * 10) // 10 MB

func main() {
	for _, n := range os.Args[1:] {
		s, err := os.Stat(n)
		if err != nil {
			log.Printf("unable to open %q: %s", n, err)
			continue
		}
		if s.Size() > maxSize {
			log.Printf("%q is too large", n)
			continue
		}
		f, err := os.Open(n)
		if err != nil {
			log.Printf("failed to read %q: %s", n, err)
			continue
		}
		defer f.Close()
		h, err := dhash(f)
		if err != nil {
			log.Printf("failed to hash %q: %s", n, err)
			continue
		}
		fmt.Println(h)
	}
}

// dhash returns a string representing an image's hash.
// Based on this algorithm:
// http://www.hackerfactor.com/blog/index.php?/archives/529-Kind-of-Like-That.html.
func dhash(reader io.Reader) (string, error) {

	img, _, err := image.Decode(reader)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Calculate the mean brightness of each block in an 9x9 grid.
	var blocks [9][9]int
	for i := 0; i < 9; i++ {
		left := bounds.Min.X + (width * i / 9)
		right := bounds.Min.X + (width * (i + 1) / 9)
		if right == left {
			right = left + 1
		}
		for j := 0; j < 9; j++ {
			top := bounds.Min.Y + (height * j / 9)
			bottom := bounds.Min.Y + (height * (j + 1) / 9)
			if bottom == top {
				bottom = top + 1
			}
			var total int64

			switch img := img.(type) {
			case *image.YCbCr:
				for y := top; y < bottom; y++ {
					rowStart := y * img.YStride
					for x := left; x < right; x++ {
						total += int64(img.Y[rowStart+x])
					}
				}
			default:
				for x := left; x < right; x++ {
					for y := top; y < bottom; y++ {
						r, g, b, _ := img.At(x, y).RGBA()
						total += int64(r+r+r+b+g+g+g+g) >> 3
					}
				}
			}
			blocks[i][j] = int(total / int64((right-left)*(bottom-top)))
		}
	}

	var result [2]uint64
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if blocks[i][j] > blocks[i][j+1] {
				result[0] |= 1 << uint(i*8+j)
			}
			if blocks[i][j] > blocks[i+1][j] {
				result[1] |= 1 << uint(i*8+j)
			}
		}
	}

	return fmt.Sprintf("%016x%016x", result[0], result[1]), nil
}
