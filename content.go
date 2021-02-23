package urchin

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

const fileName = "urchin.yml"
const notExist = `
The specified path is incorrect or the urchin file does not exist.

please refer to the text displayed after entering the help subcommand.
`
const usageText = `
subcommand is not detected.
please refer to the text displayed after entering the help subcommand.
`
const helpText = `
* creating urchin.yml template
urchin init

* execute urchin by your urchin.yml
urchin work <your urchin.yml path>

* showing this help-message again
urchin help

* referring to the website
please enter: https://github.com/gari8/urchin

* getting newer version
go get -u github.com/gari8/urchin/cmd/urchin
`
const templates = `tasks:
  - task_name: "sample"
    server_url: "https://sample.com/xxxx/xxxx"
    method: "POST"
    trial_count: 2
    queries:
      - q_name: "user_id"
        q_body: "1"
      - q_name: "title"
        q_body: "sample"
  - task_name: "sample2"
    server_url: "https://sample.com/xxxx/xxxx"
    method: "POST"
    queries:
      - q_name: "user_id"
        q_body: "2"
      - q_name: "title2"
        q_body: "sample2"
`

func (c *Content) Work() {
	if c.FilePath == nil {
		fmt.Print(notExist)
		return
	}
	buf, err := ioutil.ReadFile(*c.FilePath+"/"+fileName)
	if err != nil {
		fmt.Print(notExist)
		return
	}

	var data Data
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

func (c *Content) Create() {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err = file.WriteString(templates); err != nil {
		panic(err)
	}

	fmt.Println("create: "+fileName)
}

func (c *Content) Usage() {
	fmt.Print(usageText)
}

func (c *Content) Help() {
	fmt.Print(helpText)
}
