package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/gari8/urchin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Content struct {
	SubCmd string
	FilePath *string
}

func main() {
	var content Content
	flag.Parse()
	content.SubCmd = flag.Arg(0)
	file := flag.Arg(1)
	content.FilePath = &file

	switch content.SubCmd {
	case "work":
		content.work()
	case "init":
		content.create()
	default:
		log.Fatal(errors.New("usage text"))
	}
}

const fileName = "urchin.yml"
const notExist = `
The specified path is incorrect or the urchin file does not exist.

-h or help
`

func (c *Content) work() {
	if c.FilePath == nil {
		fmt.Print(notExist)
		return
	}
	buf, err := ioutil.ReadFile(*c.FilePath+"/"+fileName)
	if err != nil {
		fmt.Print(notExist)
		return
	}

	var data urchin.Data
	if err = yaml.Unmarshal(buf, &data); err != nil {
		log.Fatal(err)
	}
	for _, task := range data.Tasks {
		if task.TrialCnt != nil {
			for i:=0; i<*task.TrialCnt; i++ {
				task.Exe()
			}
		} else {
			task.Exe()
		}
	}
}

func (c *Content) create() {
	bytes, err := ioutil.ReadFile("_templates/"+fileName)
	if err != nil {
		panic(err)
	}

	text := string(bytes)

	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err = file.WriteString(text); err != nil {
		panic(err)
	}

	fmt.Println("create: "+fileName)
}
