package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	// APIURL ...
	APIURL = "/api"
	// APIQueueURL ...
	APIQueueURL = APIURL + "/queue"
	// APICommandURL ...
	APICommandURL = APIURL + "/command"
)

// GetQueue ...
func GetQueue() (queue []QueueElement, err error) {
	body, err := Get(APIQueueURL)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &queue)
	return
}

// ExecuteCommand ...
func ExecuteCommand(c Command) (err error) {
	j, err := json.Marshal(c)
	if err != nil {
		return
	}
	_, err = http.Post(GetURL(APIQueueURL).String(), "application/json", bytes.NewReader(j))
	return
}

// Get Wrapper for http.Get. Add authentication handling automatically.
func Get(u string) (body []byte, err error) {
	res, err := http.Get(u)
	if err != nil {
		return
	}
	if res.StatusCode == 401 {
		return nil, fmt.Errorf("authorization invalid")
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	return
}

// Post Wrapper for http.Post. Add authentication handling automatically.
func Post(u string, bodyReq io.Reader) (body []byte, err error) {
	res, err := http.Post(u, "application/json", bodyReq)
	if err != nil {
		return
	}
	if res.StatusCode == 401 {
		return nil, fmt.Errorf("authorization invalid")
	}
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	return
}
