package main

import (
	"fmt"
	"io"

	"github.com/AnimusPEXUS/goinmemfile"
)

func main() {
	t := goinmemfile.NewInMemFileFromBytes(
		[]byte("cool text"),
		0,
		false,
	)

	ro := t.RO()

	_, err := ro.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	txt, err := io.ReadAll(ro)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	fmt.Println("txt:", txt)

	_, err = ro.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	_, err = ro.Write([]byte{1, 2, 3})
	if err != nil {
		fmt.Println("err:", err)
		return
	}

	fmt.Println("ok")
}
