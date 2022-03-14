package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal(errors.New("arguments is not enough"))
		return
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
		return
	}

	code := RunCmd(os.Args[2:], env)
	if code != 0 {
		fmt.Println(code)
	}
}
