package main

import (
	"fmt"
	"github.com/goushalk/chromaboard/internal/storage"
)

func main() {
	err := storage.ConfigCreator()
	fmt.Println(err)
}
