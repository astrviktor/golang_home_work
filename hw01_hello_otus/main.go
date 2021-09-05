package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

const helloStr = "Hello, OTUS!"

func main() {
	fmt.Println(stringutil.Reverse(helloStr))
}
