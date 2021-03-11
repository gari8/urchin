package urchin

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	multi    = "multipart/form-data"
	appJson  = "application/json"
	formData = "application/x-www-form-urlencoded"
)

func (t *Task) Exe() (*string, error) {

	switch t.Method {
	case http.MethodGet:
		return t.get()
	case http.MethodPost:
		if t.ContentType == nil {
			return t.postForm()
		}
		if strings.Contains(*t.ContentType, multi) {
			return t.postMultipart()
		} else if strings.Contains(*t.ContentType, appJson) {
			return t.postJson()
		} else {
			return t.postForm()
		}
	case http.MethodPut:
		return nil, nil
	case http.MethodDelete:
		return nil, nil
	default:
		return nil, errors.New("invalid http method please confirm it")
	}
}

func (t *Task) get() (*string, error) {
	req, err := http.NewRequest(t.Method, t.ServerURL, nil)
	if err != nil {
		return nil, err
	}

	for _, h := range t.Headers {
		if h == nil {
			continue
		}
		req.Header.Set(*h.HType, *h.HBody)
	}

	params := req.URL.Query()
	for _, q := range t.Queries {
		// GET で upload 禁止
		if q == nil || q.QFile != nil {
			continue
		}
		params.Add(*q.QName, *q.QBody)
	}

	return t.clientEvent(req)
}

func fileReader(fileName string) ([]byte, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	byte, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return byte, nil
}

func (t *Task) postForm() (*string, error) {
	values := url.Values{}
	for _, q := range t.Queries {
		if q == nil {
			continue
		}
		if q.QFile != nil {
			byte, err := fileReader(*q.QFile)
			if err != nil {
				return nil, err
			}
			str := fmt.Sprintf("%s", byte)
			values.Add(*q.QName, str)
		} else {
			values.Add(*q.QName, *q.QBody)
		}
	}

	req, err := http.NewRequest(t.Method, t.ServerURL, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", formData)

	for _, h := range t.Headers {
		if h == nil {
			continue
		}
		req.Header.Set(*h.HType, *h.HBody)
	}

	return t.clientEvent(req)
}

func (t *Task) postJson() (*string, error) {
	var jsonV []byte

	if t.QJson != nil {
		jsonV = []byte(*t.QJson)
	}

	req, err := http.NewRequest(t.Method, t.ServerURL, bytes.NewBuffer(jsonV))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", appJson)

	for _, h := range t.Headers {
		if h == nil {
			continue
		}
		req.Header.Set(*h.HType, *h.HBody)
	}

	return t.clientEvent(req)
}

func (t *Task) postMultipart() (*string, error) {

	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)

	for _, q := range t.Queries {
		if q.QFile != nil && q.QName != nil {
			file, err := os.Open(*q.QFile)
			if err != nil {
				return nil, err
			}

			fw, err := mw.CreateFormFile(*q.QName, *q.QFile)
			_, err = io.Copy(fw, file)
			if err != nil {
				return nil, err
			}
			_ = file.Close()
		} else if q.QName != nil && q.QBody != nil {
			w, err := mw.CreateFormField(*q.QName)
			if err != nil {
				return nil, err
			}
			_, err = w.Write([]byte(*q.QBody))
			if err != nil {
				return nil, err
			}
		}
	}

	contentType := mw.FormDataContentType()

	_ = mw.Close()

	req, err := http.NewRequest(t.Method, t.ServerURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	for _, h := range t.Headers {
		if h == nil {
			continue
		}
		req.Header.Set(*h.HType, *h.HBody)
	}

	return t.clientEvent(req)
}

func (t *Task) clientEvent(req *http.Request) (*string, error) {
	if t.BasicAuth != nil {
		req.SetBasicAuth(t.BasicAuth.UserName, t.BasicAuth.Password)
	}

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	r := string(body)
	return &r, nil
}

func containsQueryFile(s []*Query) bool {
	for _, v := range s {
		if v == nil || v.QFile != nil {
			return true
		}
	}
	return false
}
