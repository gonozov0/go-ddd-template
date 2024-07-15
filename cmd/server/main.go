package main

import (
	"os"

	"go-echo-ddd-template/internal"
)

func main() {
	if err := internal.Run(); err != nil {
		os.Exit(1)
	}
}
