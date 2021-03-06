# urchin

## What is this?

A simple client tool made by go.

## What you can do with this

- You could send a request based on your own written urchin.yml file.
- You could send requests in parallel for each task.

## How to install

```
go get -u github.com/gari8/urchin/cmd/urchin
```

## How to use

```
// creating urchin.yml template
urchin init

// executing program by ./urchin.yml
urchin work .

// when you want to keep an execution log
urchin -L work .

// checking if your urchin.yml is in the correct format
urchin check .
```

## How to write

```urchin.yml
#(*) = must

tasks:
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
```

## Applied usage

- You need to make two requests to the same server many times every 10 seconds.
query or params

- req1
```bash
user_id: "test_user",
title: "sql_file_content",
sql_content: "<content by ./dump.sql file>"
```

- req2
```bash
user_id: "test_user",
title: "sql_file_content2",
sql_content: "<content by ./dump_2.sql file>"
```
> condition: !!! req2 is executed one second later than req1 !!!

\example
```urchin.yml
base: &base
  server_url: "https://xxx.xxx.jp/api/v1/accepting_req"
  method: "POST"
  basic_auth:
    user_name: "xxx"
    password: "xxxxxx"
test_user: &test_user
  q_name: "user_id"
  q_body: "test_user"
tasks:
  - task_name: "req1"
    <<: *base
    queries:
      - *test_user
      - q_name: "title"
        q_body: "sql_file_content"
      - q_name: "sql_content"
        q_file: "./dump.sql"
  - task_name: "req2"
    <<: *base
    delay_ms: 1000 # Delaying execution in milliseconds
    queries:
      - *test_user
      - q_name: "title"
        q_body: "sql_file_content2"
      - q_name: "sql_content"
        q_file: "./dump_2.sql"
task_interval: 10
```
