package main

import (
	"errors"
	"flag"
	"github.com/gari8/urchin"
	"log"
)



func main() {
	var content urchin.Content
	flag.Parse()
	content.SubCmd = flag.Arg(0)
	file := flag.Arg(1)
	content.FilePath = &file

	switch content.SubCmd {
	case "work":
		content.Work()
	case "init":
		content.Create()
	default:
		log.Fatal(errors.New("usage text"))
	}
}



