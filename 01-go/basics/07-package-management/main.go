package main

import (
	"fmt"

	"github.com/dyson/backend-roadmap/go/internal/toolkit"
)

func main() {
	fmt.Println("toolkit.Repeat:", toolkit.Repeat("go ", 3))
	fmt.Println("request id:", toolkit.NewRequestID())
}
