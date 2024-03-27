package main

import (
	"go-echo-template/internal"
)

func main() {
	if err := internal.Run(); err != nil {
		panic(err)
	}
}
