package main

import (
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

	var bestN, bestOffset, bestSize int64

	for offset := int64(0); offset <= maxOffset; offset++ {
		for size := int64(minSize); size <= maxSize; size++ {
			n := repeats(r, offset, size)
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

	prefix := make([]byte, bestOffset)
	must(r.ReadAt(prefix, 0))

	repeat := make([]byte, bestSize)
	must(r.ReadAt(repeat, bestOffset))

	suffix := make([]byte, int64(r.Len())-bestOffset-bestN)
	must(r.ReadAt(suffix, bestOffset+bestN))

	log.Printf("len: %d, offset: %d, repeats: %d*%d=%d", r.Len(), bestOffset, bestSize, bestN/bestSize, bestN)

	log.Printf("prefix: %x", prefix)
	log.Printf("repeat: %x", repeat)
	log.Printf("suffix: %x", suffix)
}

func repeats(r *mmap.ReaderAt, offset, size int64) int64 {
	first := make([]byte, size)
	next := make([]byte, size)
	_, err := r.ReadAt(first, offset)
	if err != nil {
		return 0
	}

	max := int64(r.Len())
	n := size
	for {
		if offset+n+size > max {
			break
		}
		_, err := r.ReadAt(next, offset+n)
		if err != nil {
			return 0
		}
		if !reflect.DeepEqual(first, next) {
			break
		}
		n += size
	}
	return n
}
