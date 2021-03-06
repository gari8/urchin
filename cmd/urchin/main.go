package main

import (
	"flag"
	"github.com/gari8/urchin"
)



func main() {
	var content urchin.Content
	flag.BoolVar(&content.LogMode, "L", false, "add log mode")
	flag.Parse()
	content.SubCmd = flag.Arg(0)
	file := flag.Arg(1)
	content.FilePath = &file

	switch content.SubCmd {
	case "work":
		content.Work()
	case "init":
		content.Create()
	case "help":
		content.Help()
	case "check":
		content.Check()
	default:
		content.Usage()
	}
}



