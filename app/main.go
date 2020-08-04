package main

import (
	"flag"
)

func main() {
	flag.Parse()
	mode := flag.Arg(0)
	if mode == "web" {
		Web()
	}
	if mode == "task" {

	}
}
