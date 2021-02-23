package urchin

type Data struct {
	Tasks []Task `yaml:"tasks"`
	TaskInterval *int `yaml:"task_interval"`
	MaxTrialCnt *int `yaml:"max_trial_count"`
}

type Task struct {
	TaskName string `yaml:"task_name"`
	ServerURL string `yaml:"server_url"`
	Method string `yaml:"method"`
	ContentType string `yaml:"content_type"`
	TrialCnt *int `yaml:"trial_count"`
	Queries []Query
	Headers []Header
}

type Header struct {
	HType string `yaml:"h_type"`
}

type Query struct {
	QName *string `yaml:"q_name"`
	QBody *string `yaml:"q_body"`
	QFile *string `yaml:"q_file"`
}

type Content struct {
	SubCmd string
	FilePath *string
}
