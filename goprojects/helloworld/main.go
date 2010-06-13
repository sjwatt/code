package main

import (
	"./file"
	"fmt"
	"os"
)

func main() {
	hello := []byte("hello, world\n")
	file.Stdout.Write(hello)
	file, err := file.Open("/root/goprojects/helloworld/test.txt", 0, 0)
	if file == nil {
		fmt.Printf("can't open file; err=%s\n", err.String())
		os.Exit(1)
	}
}

