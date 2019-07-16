package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 2 {
		log.Fatal("<USAGE> $ namecompression /path/to/binary OFFSET")
	}
	r, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal("Problem", err)
	}

	offset, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		log.Fatal("Problem", err)
	}
	name, err := NameAt(r, offset)
	if err != nil {
		log.Fatal("Problem", err)
	}

	fmt.Println("Name:", name)

}

func NameAt(r io.Reader, offset int) (string, error) {
	var name bytes.Buffer
	root := byte(0x00)
	buf := bufio.NewReader(r)
	n, err := buf.Discard(offset)
	if n != offset {
		return "", fmt.Errorf("Discard problem : %d", n)
	}
	if err != nil {
		return "", err
	}

	lengthHeader, err := buf.ReadByte()
	if err != nil {
		return "", err
	}

	fmt.Println("LenghtHeader:", lengthHeader)

	for lengthHeader != root {
		readTo, isPointer := FetchOffset(int(lengthHeader))
		if isPointer {
			fmt.Println("Pointer!")
		} else {
			fmt.Println("Not a pointer. Offset", readTo)
		}

		buffer := make([]byte, readTo)
		_, err := io.ReadFull(buf, buffer)
		// if n != readTo {
		// 	return "", fmt.Errorf("Something went wrong : %d", n)
		// }
		if err != nil {
			return "", fmt.Errorf("Something else went wrong : %s", err)
		}
		fmt.Println("Buffer", buffer)
		name.Write(buffer)
		fmt.Println("Buffer", string(buffer))
		lengthHeader, err = buf.ReadByte()
		if err != nil {
			return "", err
		}
	}
	return name.String(), nil
}

func FetchOffset(b int) (offset int, isPointer bool) {
	pointerMask := 0xC0
	pointerSet := (b & pointerMask) == pointerMask // clear 6 lowest bits, check what's left is the mask
	if !pointerSet {
		return b, false
	}

	return (b ^ pointerMask), true
}
