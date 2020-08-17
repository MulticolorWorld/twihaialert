package main

import (
	"app/module"
	"flag"
)

func main() {
	flag.Parse()
	mode := flag.Arg(0)
	if mode == "web" {
		module.Web()
	}
	if mode == "task" {
		module.Task()
	}
}
