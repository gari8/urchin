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
const templates = `tasks: # must
  - task_name: "sample" # must
    server_url: "https://sample.com/xxxx/xxxx" # must
    method: "POST" # must
    trial_count: 2 # 1周で送信する回数
    content_type: "application/json" # content-typeはこちらにかける
    q_json: "{ \"user_id\": 1 }" # application/json　の場合jsonを書く
    headers: # header
      - h_type: # header名 
        h_body:
  - task_name: "sample2" # must
    server_url: "https://sample.com/xxxx/xxxx" # must
    method: "POST" # must
    basic_auth: # basic認証
      user_name: "user_1" # basic_auth がある場合 must
      password: "******" # basic_auth がある場合 must
    queries:
      - q_name: "user_id"
        q_body: "2"
      - q_name: "title2"
        q_file: "./index.html" # text file を読み込んで送信
task_interval: 3 # インターバルの秒数(s)　指定しなければ1周のみ
max_trial_count: 5 # 合計何周するか　指定しなければ止めるまでループ
`
const boundaryText = "** -------------------- ** boundary ** -------------------- **"

const (
	black = iota + 30
	red
	green
	yellow
	blue
	magenta
	cyan
	white
)

func (c *Content) Work() {
	if c.FilePath == nil {
		handlingWarning(notExist)
		return
	}

	buf, err := ioutil.ReadFile(*c.FilePath+"/"+fileName)
	if err != nil {
		handlingWarning(notExist)
		return
	}

	var data Data
	if err = yaml.Unmarshal(buf, &data); err != nil {
		// yml側の問題なのでwarning
		handlingWarning(invalidFile)
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
				handlingAny(cyan, boundaryText)
				taskRunner(data)
				index++
				if data.MaxTrialCnt != nil && index == *data.MaxTrialCnt {
					handlingSuccess("Task completed " + strconv.Itoa(index) + " times in total " + "(" + strconv.Itoa(index*(*data.TaskInterval)) + "s)")
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
				str, err := task.Exe()
				message := fmt.Sprintf("task_name: %s, server_url: %s, method: %s", task.TaskName, task.ServerURL, task.Method)
				handlingAny(magenta, message)
				if err != nil || str == nil {
					handlingError(err)
				} else {
					handlingAny(cyan, *str)
				}
			}
		} else {
			str, err := task.Exe()
			message := fmt.Sprintf("task_name: %s, server_url: %s, method: %s", task.TaskName, task.ServerURL, task.Method)
			handlingAny(magenta, message)
			if err != nil || str == nil {
				handlingError(err)
			} else {
				handlingAny(cyan, *str)
			}
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
		handlingError(err)
	}
	handlingSuccess("create: "+fileName)
}

func (c *Content) Usage() {
	handlingWarning(usageText)
}

func (c *Content) Help() {
	handlingAny(magenta, helpText)
}

func handlingError(err error) {
	fmt.Printf("\x1b[%dm%s\x1b[0m\n", red, err)
}

func handlingWarning(str string) {
	fmt.Printf("\x1b[%dm%s\x1b[0m\n", yellow, str)
}

func handlingSuccess(str string) {
	fmt.Printf("\x1b[%dm%s\x1b[0m\n", green, str)
}

func handlingAny(number int, str string) {
	fmt.Printf("\x1b[%dm%s\x1b[0m\n", number, str)
}
