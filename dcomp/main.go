package main

// This is a simplified function copy based on https://github.com/andybalholm/dhash/blob/master/hash.go

import (
	"fmt"
	"log"
	"math/bits"
	"os"
	"strconv"
)

var maxSize = int64(1024 * 1024 * 10) // 10 MB

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Two arguments required; got %d", len(os.Args)-1)
	}

	a := os.Args[len(os.Args)-2]
	b := os.Args[len(os.Args)-1]

	if len(a) != 32 {
		log.Fatalf("%q is not 32 characters", a)
	}
	if len(b) != 32 {
		log.Fatalf("%q is not 32 characters", a)
	}

	score := getScore(a[:16], b[:16]) + getScore(a[16:], b[16:])
	fmt.Printf("%.3f\n", float64(1) - float64(score) / float64(128))

}

func getScore(a, b string) int {
	xor := getUint(a) ^ getUint(b)
	return bits.OnesCount64(xor)
}

func getUint(s string) uint64 {
	u, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		log.Fatalf("Failed to do something: %s\n", err)
	}
	return u
}
