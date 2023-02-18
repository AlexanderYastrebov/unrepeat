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
	var maxPrefix, maxSuffix, minSize, maxSize int64
	var text bool

	flag.Int64Var(&maxPrefix, "max-prefix", 100, "maximum allowed prefix length")
	flag.Int64Var(&maxSuffix, "max-suffix", 100, "maximum allowed suffix length")
	flag.Int64Var(&minSize, "min-size", 10, "minimum size of the repeating block")
	flag.Int64Var(&maxSize, "max-size", 100, "maximum size of the repeating block")
	flag.BoolVar(&text, "text", false, "print text instead of hexadecimal")

	flag.Parse()

	r := must(mmap.Open(flag.Arg(0)))
	defer r.Close()

	data := getData(r)

	var bestN, bestOffset, bestSize int64

	for prefix := int64(0); prefix <= maxPrefix; prefix++ {
		for size := int64(minSize); size <= maxSize; size++ {
			n := repeats(data, prefix, size)
			if n > 1 && n*size > bestN*bestSize {
				bestN, bestOffset, bestSize = n, prefix, size
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

	if int64(len(suffix)) > maxSuffix {
		fmt.Printf("Suffix too large: %d\n", len(suffix))
		return
	}

	fmt.Printf("length: %d, offset: %d, repeats: %d*%d=%d\n", len(data), bestOffset, bestSize, bestN, bestN*bestSize)

	if text {
		fmt.Printf("prefix: %s\nrepeat: %s\nsuffix: %s\n", prefix, repeat, suffix)
	} else {
		fmt.Printf("prefix: %x\nrepeat: %x\nsuffix: %x\n", prefix, repeat, suffix)
	}
}

func repeats(data []byte, prefix, size int64) int64 {
	max := int64(len(data))
	if prefix+size > max {
		return 0
	}
	first := data[prefix : prefix+size]

	n := int64(1)
	for {
		start := prefix + size*n
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
