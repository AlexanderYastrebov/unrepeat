package main

import (
	"bytes"
	"flag"
	"log"
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
			if n > size && n > bestN {
				bestN, bestOffset, bestSize = n, offset, size
				continue
			}
		}
	}

	if bestN == 0 {
		log.Printf("No repeats found")
		return
	}

	prefix := data[0:bestOffset]
	repeat := data[bestOffset : bestOffset+bestSize]
	suffix := data[bestOffset+bestN:]

	log.Printf("len: %d, offset: %d, repeats: %d*%d=%d", len(data), bestOffset, bestSize, bestN/bestSize, bestN)

	log.Printf("prefix: %x", prefix)
	log.Printf("repeat: %x", repeat)
	log.Printf("suffix: %x", suffix)
}

func repeats(data []byte, offset, size int64) int64 {
	max := int64(len(data))
	if offset+size > max {
		return 0
	}
	first := data[offset : offset+size]

	n := size
	for {
		if offset+n+size > max {
			break
		}
		if !bytes.Equal(first, data[offset+n:offset+n+size]) {
			break
		}
		n += size
	}
	return n
}

func getData(r *mmap.ReaderAt) []byte {
	v := reflect.ValueOf(r)
	f := reflect.Indirect(v).FieldByName("data")
	return f.Bytes()
}
