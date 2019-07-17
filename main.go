package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 2 {
		log.Fatal("<USAGE> $ namecompression /path/to/binary OFFSET")
	}
	r, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	offset, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		log.Fatal(err)
	}

	name, err := NameAt(r, offset)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(name)
}

func NameAt(r io.Reader, off int) (string, error) {
	var name []string
	root := byte(0x00)
	all, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	buf := bytes.NewReader(all)
	buf.Seek(int64(off), 0)
	if err != nil {
		return "", err
	}

	lengthHeader, err := buf.ReadByte()
	if err != nil {
		return "", err
	}

	for lengthHeader != root {
		if IsPointer(int(lengthHeader)) {
			offset, err := buf.ReadByte()
			if err != nil {
				return "", err
			}

			buf.Seek(int64(offset), 0)
			if err != nil {
				return "", err
			}

			lengthHeader, err = buf.ReadByte()
			if err != nil {
				return "", err
			}
			continue
		}

		buffer := make([]byte, lengthHeader)
		_, err := io.ReadFull(buf, buffer)
		if err != nil {
			return "", err
		}

		name = append(name, string(buffer))

		lengthHeader, err = buf.ReadByte()
		if err != nil {
			return "", err
		}
	}

	return strings.Join(name, "."), nil
}

func IsPointer(b int) bool {
	pointerMask := 0xC0
	return (b & pointerMask) == pointerMask
}
