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
```

## How to write

```urchin.yml
#(*) = must

tasks:
  - task_name: "" #(*)
    server_url: "" #(*)
    method: "POST" #(*)
    trial_count: 2
    content_type: ""
    # "application/json" or "application/x-www-form-urlencoded", "multipart/form-data" is WIP
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
        q_file: "./a.csv" # reading local file and then sending the content as string
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
    queries:
      - *test_user
      - q_name: "title"
        q_body: "sql_file_content2"
      - q_name: "sql_content"
        q_file: "./dump_2.sql"
task_interval: 10
```
