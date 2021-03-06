package urchin

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
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

* when you want to keep an execution log
urchin -L work <your urchin.yml path>

* showing this help-message again
urchin help

* checking if your urchin.yml is in the correct format
urchin check .

* referring to the website
please enter: https://github.com/gari8/urchin

* getting newer version
go get -u github.com/gari8/urchin/cmd/urchin
`
const templates = `tasks:
  - task_name: "" #(*)
    server_url: "" #(*)
    method: "POST" #(*)
    delay_ms: 1000 # Delaying execution in milliseconds
    trial_count: 2
    content_type: ""
    # "application/json" or "application/x-www-form-urlencoded", "multipart/form-data"
    headers:
      - h_type: ""
        h_body: ""
    basic_auth:
      user_name: ""
      password: "******"
    # query pattern (not "content-type": "application/json")
    queries:
      - q_name: "" # field_name
        q_body: "" # content
      - q_name: ""
        q_file: "./a.csv" #("content-type": "application/x-www-form-urlencoded") reading local file and then sending the content as string
                          #("content-type": "multipart/form-data") uploading local file to your oriented server
    # query pattern ("content-type": "application/json")
    q_json: "{\"user_name\": \"takashi\"}"
task_interval: 3 #(s)
max_trial_count: 3
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

func NowJST() string {
	now := time.Now()
	nowUTC := now.UTC()
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	nowJST := nowUTC.In(jst)

	return nowJST.Format(time.RFC3339)
}

func (c *Content) Work() {
	if c.FilePath == nil {
		handlingWarning(notExist)
		return
	}

	buf, err := ioutil.ReadFile(*c.FilePath + "/" + fileName)
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

	wg := new(sync.WaitGroup)

	if data.TaskInterval != nil {
		index := 1
		ticker := time.NewTicker(time.Millisecond * time.Duration(*data.TaskInterval*1000))
		defer ticker.Stop()
		c.taskRunner(data, wg)
		for {
			select {
			case <-ticker.C:
				c.taskRunner(data, wg)
				index++
				if data.MaxTrialCnt != nil && index == *data.MaxTrialCnt {
					wg.Wait()
					handlingSuccess("Task completed " + strconv.Itoa(index) + " times in total " + "(" + strconv.Itoa(index*(*data.TaskInterval)) + "s)")
					return
				}
			}
		}
	} else {
		c.taskRunner(data, wg)
	}

	wg.Wait()
}

func (c *Content) addLog(log string) {
	if !c.LogMode {
		return
	}
	name := "urchin.log"
	fp, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		return
	}
	defer fp.Close()
	_, _ = fmt.Fprintln(fp, log)
}

func (t *Task) newMessage(str *string) string {
	var ct string
	switch {
	case t.ContentType == nil:
		ct = formData
	case strings.Contains(*t.ContentType, multi):
		ct = multi
	case strings.Contains(*t.ContentType, appJson):
		ct = appJson
	}
	if str == nil {
		return fmt.Sprintf("[date: %s] [task_name: %s] [server_url: %s] [method: %s] [content-type: %s]", NowJST(), t.TaskName, t.ServerURL, t.Method, ct)
	}
	return fmt.Sprintf("[date: %s] [task_name: %s] [server_url: %s] [method: %s] [content-type: %s]\n => %s", NowJST(), t.TaskName, t.ServerURL, t.Method, ct, *str)
}

func (c *Content) taskRunner(data Data, wg *sync.WaitGroup) {
	wg.Add(len(data.Tasks))
	fmt.Println("")
	for _, task := range data.Tasks {
		go func(task Task) {
			if task.DelayMs != nil {
				time.Sleep(time.Millisecond * time.Duration(*task.DelayMs))
			}
			handlingAny(magenta, fmt.Sprintf("sent a request to task: %s", task.TaskName))
			if task.TrialCnt != nil {
				for i := 0; i < *task.TrialCnt; i++ {
					str, err := task.Exe()
					if err != nil || str == nil {
						m := fmt.Sprintf("%s", err)
						log := task.newMessage(&m)
						handlingAny(red, log)
						c.addLog(log)
					} else {
						log := task.newMessage(str)
						handlingAny(cyan, log)
						c.addLog(log)
					}
				}
			} else {
				str, err := task.Exe()
				if err != nil || str == nil {
					m := fmt.Sprintf("%s", err)
					log := task.newMessage(&m)
					handlingAny(red, log)
					c.addLog(log)
				} else {
					log := task.newMessage(str)
					handlingAny(cyan, log)
					c.addLog(log)
				}
			}
			wg.Done()
		}(task)
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
	handlingSuccess("create: " + fileName)
}

func (c *Content) Usage() {
	handlingWarning(usageText)
}

func (c *Content) Help() {
	handlingAny(magenta, helpText)
}

func (c *Content) Check() {
	if c.FilePath == nil {
		handlingWarning(notExist)
		return
	}

	buf, err := ioutil.ReadFile(*c.FilePath + "/" + fileName)
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
	handlingSuccess("OK, your urchin.yml has passed all the checks")
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
