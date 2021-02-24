# urchin

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
