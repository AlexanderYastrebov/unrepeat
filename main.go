package main

import (
	"bytes"
	"flag"
	"fmt"
	"reflect"

	"golang.org/x/exp/mmap"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func main() {
	var maxOffset, minSize, maxSize int64

	flag.Int64Var(&maxOffset, "max-offset", 100, "maximum offset from the start of the file to check")
	flag.Int64Var(&minSize, "min-size", 10, "minimum size of the repeating block")
	flag.Int64Var(&maxSize, "max-size", 100, "maximum size of the repeating block")

	flag.Parse()

	r := must(mmap.Open(flag.Arg(0)))
	defer r.Close()

	data := getData(r)

	var bestN, bestOffset, bestSize int64

	for offset := int64(0); offset <= maxOffset; offset++ {
		for size := int64(minSize); size <= maxSize; size++ {
			n := repeats(data, offset, size)
			if n > 1 && n*size > bestN*bestSize {
				bestN, bestOffset, bestSize = n, offset, size
			}
		}
	}

	if bestN < 2 {
		fmt.Println("No repeats found")
		return
	}

	prefix := data[0:bestOffset]
	repeat := data[bestOffset : bestOffset+bestSize]
	suffix := data[bestOffset+bestN*bestSize:]

	fmt.Printf("length: %d, offset: %d, repeats: %d*%d=%d\n", len(data), bestOffset, bestSize, bestN, bestN*bestSize)

	fmt.Printf("prefix: %x\n", prefix)
	fmt.Printf("repeat: %x\n", repeat)
	fmt.Printf("suffix: %x\n", suffix)
}

func repeats(data []byte, offset, size int64) int64 {
	max := int64(len(data))
	if offset+size > max {
		return 0
	}
	first := data[offset : offset+size]

	n := int64(1)
	for {
		start := offset + size*n
		end := start + size
		if end > max {
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
