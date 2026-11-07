package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("hello, backend roadmap")
	fmt.Printf("runtime:  %s\n", runtime.Version())
	fmt.Printf("platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("cpus:     %d\n", runtime.NumCPU())
	fmt.Printf("maxprocs: %d\n", runtime.GOMAXPROCS(0))
}
