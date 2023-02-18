package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"reflect"

	"golang.org/x/exp/mmap"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func exit(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func main() {
	var maxPrefix, maxSuffix, minSize, maxSize int
	var text bool

	flag.IntVar(&maxPrefix, "max-prefix", 100, "maximum allowed prefix length")
	flag.IntVar(&maxSuffix, "max-suffix", 100, "maximum allowed suffix length")
	flag.IntVar(&minSize, "min-size", 10, "minimum size of the repeating block")
	flag.IntVar(&maxSize, "max-size", 100, "maximum size of the repeating block")
	flag.BoolVar(&text, "text", false, "print text instead of hexadecimal")

	flag.Parse()

	r := must(mmap.Open(flag.Arg(0)))
	defer r.Close()

	data := getData(r)

	var bestN, bestOffset, bestSize int

	for prefix := 0; prefix <= maxPrefix; prefix++ {
		for size := minSize; size <= maxSize; size++ {
			n := repeats(data, prefix, size)
			if n > 1 && n*size > bestN*bestSize {
				bestN, bestOffset, bestSize = n, prefix, size
			}
		}
	}

	if bestN < 2 {
		exit("No repeats found")
		return
	}

	prefix := data[0:bestOffset]
	repeat := data[bestOffset : bestOffset+bestSize]
	suffix := data[bestOffset+bestN*bestSize:]

	if len(suffix) > maxSuffix {
		exit("Suffix too large: %d\n", len(suffix))
		return
	}

	fmt.Printf("length: %d, offset: %d, repeats: %d*%d=%d\n", len(data), bestOffset, bestSize, bestN, bestN*bestSize)

	if text {
		fmt.Printf("prefix: %s\nrepeat: %s\nsuffix: %s\n", prefix, repeat, suffix)
	} else {
		fmt.Printf("prefix: %x\nrepeat: %x\nsuffix: %x\n", prefix, repeat, suffix)
	}
}

func repeats(data []byte, prefix, size int) int {
	if prefix+size > len(data) {
		return 0
	}
	first := data[prefix : prefix+size]

	n := 1
	for {
		start := prefix + size*n
		end := start + size
		if end > len(data) {
			break
		}
		if !bytes.Equal(first, data[start:end]) {
			break
		}
		n++
	}
	return n
}

func getData(r *mmap.ReaderAt) []byte {
	v := reflect.ValueOf(r)
	f := reflect.Indirect(v).FieldByName("data")
	return f.Bytes()
}
