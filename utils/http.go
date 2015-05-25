package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jxwr/cc/frontend/api"
)

func do(method, url string, in, out interface{}, timeout time.Duration) (*api.Response, error) {
	reqJson, _ := json.Marshal(in)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	client.Timeout = timeout
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		var fail api.Response
		err = json.Unmarshal([]byte(body), &fail)
		return &fail, err
	} else {
		err = json.Unmarshal([]byte(body), out)
		return nil, err
	}
}

func HttpPost(url string, in, out interface{}, timeout time.Duration) (*api.Response, error) {
	return do("POST", url, in, out, timeout)
}

func HttpPut(url string, in, out interface{}, timeout time.Duration) (*api.Response, error) {
	return do("PUT", url, in, out, timeout)
}

func HttpGet(url string, in, out interface{}, timeout time.Duration) (*api.Response, error) {
	return do("GET", url, in, out, timeout)
}
