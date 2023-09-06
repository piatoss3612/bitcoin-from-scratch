package main

import (
	"bytes"
	"fmt"
)

func main() {
	b := []byte{0x00, 0x01, 0x02, 0x03}
	buf := bytes.NewBuffer(b)

	tmp := buf.Next(3)
	fmt.Println(tmp)
	fmt.Println(b)

	fmt.Println(int(0x100))
}
