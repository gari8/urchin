package urchin

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
	"time"
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
const invalidFile = `
your urchin.yml is invalid.

please refer to the text displayed after entering the help subcommand.
`
const helpText = `
* creating urchin.yml template
urchin init

* executing urchin by your urchin.yml
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
task_interval: 30
max_trial_count: 10
`
const boundaryText = "** -------------------- ** boundary ** -------------------- **"

func (c *Content) Work() {
	if c.FilePath == nil {
		fmt.Printf("\x1b[33m%s\x1b[0m\n", notExist)
		return
	}

	buf, err := ioutil.ReadFile(*c.FilePath+"/"+fileName)
	if err != nil {
		fmt.Printf("\x1b[33m%s\x1b[0m\n", notExist)
		return
	}

	var data Data
	if err = yaml.Unmarshal(buf, &data); err != nil {
		fmt.Printf("\x1b[33m%s\x1b[0m\n", invalidFile)
		return
	}

	if data.TaskInterval != nil {
		index := 1
		ticker := time.NewTicker(time.Millisecond * time.Duration(*data.TaskInterval*1000))
		defer ticker.Stop()
		taskRunner(data)
		for {
			select {
			case <-ticker.C:
				fmt.Printf("\x1b[36m%s\x1b[0m\n", boundaryText)
				taskRunner(data)
				index++
				if index == *data.MaxTrialCnt {
					fmt.Printf("\x1b[32m%s\x1b[0m\n", "Task completed " + strconv.Itoa(index) + " times in total " + "(" + strconv.Itoa(index*(*data.TaskInterval)) + "s)")
					return
				}
			}
		}
	} else {
		taskRunner(data)
	}
}

func taskRunner(data Data) {
	fmt.Println("")
	for _, task := range data.Tasks {
		if task.TrialCnt != nil {
			for i:=0; i<*task.TrialCnt; i++ {
				task.Exe()
			}
		} else {
			task.Exe()
		}
	}
	fmt.Println("")
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

	fmt.Printf("\x1b[32m%s\x1b[0m\n", "create: "+fileName)
}

func (c *Content) Usage() {
	fmt.Printf("\x1b[33m%s\x1b[0m\n", usageText)
}

func (c *Content) Help() {
	fmt.Printf("\x1b[35m%s\x1b[0m\n", helpText)
}
