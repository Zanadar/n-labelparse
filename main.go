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
	// Possibly seek here?
	name, err := NameAt(r, offset)
	if err != nil {
		log.Fatal("Problem", err)
	}

	fmt.Println("Name:", name)

}

func NameAt(r io.Reader, off int) (string, error) {
	var name bytes.Buffer
	root := byte(0x00)
	all, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	buf := bytes.NewReader(all)
	buf.Seek(int64(off), 0)

	lengthHeader, err := buf.ReadByte()
	if err != nil {
		return "", err
	}

	fmt.Println("LenghtHeader:", lengthHeader)

	count := 0
	for lengthHeader != root {
		if count > 10 {
			return "Probelm", fmt.Errorf("Problem with somethign")
		}
		isPointer := FetchOffset(int(lengthHeader))
		if isPointer {
			offset, err := buf.ReadByte()
			if err != nil {
				return "", err
			}
			fmt.Println("Pointer! Readto:", offset)
			// buf.Reset(all)
			buf.Seek(int64(offset), 0)
			if err != nil {
				return "", err
			}
			lengthHeader, err = buf.ReadByte()
			if err != nil {
				return "", err
			}
			count++
			continue
		}

		buffer := make([]byte, lengthHeader)
		_, err := io.ReadFull(buf, buffer)
		// if n != readTo {
		// 	return "", fmt.Errorf("Something went wrong : %d", n)
		// }
		if err != nil {
			return "", fmt.Errorf("Something else went wrong : %s", err)
		}
		fmt.Println("Buffer", buffer)
		name.Write(buffer)
		name.WriteString(".") // how do we not do this on the last call
		fmt.Println("Buffer", string(buffer))
		lengthHeader, err = buf.ReadByte()
		if err != nil {
			return "", err
		}
	}
	return name.String(), nil
}

func FetchOffset(b int) bool {
	pointerMask := 0xC0
	return (b & pointerMask) == pointerMask
}
