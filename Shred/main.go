// main.go
package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	filename := flag.String("n", "", "filename to shred")
	verbose := flag.Bool("v", false, "print stats")

	flag.Parse()
	if *filename == "" {
		fmt.Fprintln(os.Stderr, "Shred -v -n filename")
		os.Exit(1)
	}
	start_time := time.Now()
	Shred(*filename)
	duration := time.Since(start_time)
	if true == *verbose {
		fmt.Println("time:", duration)
	}
}
