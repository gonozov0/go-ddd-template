package main

import (
	"go-echo-ddd-template/internal"
)

func main() {
	if err := internal.Run(); err != nil {
		panic(err)
	}
}
