package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	// APIURL ...
	APIURL = "/api"
	// APIQueueURL ...
	APIQueueURL = APIURL + "/queue"
	// APICommandURL ...
	APICommandURL = APIURL + "/command"
	// APIHistoryURL ...
	APIHistoryURL = APIURL + "/history"
	// StatusCompleted ...
	StatusCompleted = "Completed"
	// TrackedDownloadStatusWarning ...
	TrackedDownloadStatusWarning = "Warning"
)

// GetQueue ...
func GetQueue() (queue []QueueElement, err error) {
	body, err := Get(GetURL(APIQueueURL).String())
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &queue)
	return
}

// GetHistory ...
func GetHistory(page int) (history History, err error) {
	u := GetURL(APIHistoryURL)
	query := u.Query()
	query.Add("page", strconv.Itoa(page))
	u.RawQuery = query.Encode()
	body, err := Get(u.String())
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &history)
	return
}

// ExecuteCommand ...
func ExecuteCommand(c Command) (err error) {
	j, err := json.Marshal(c)
	if err != nil {
		return
	}
	_, err = http.Post(GetURL(APICommandURL).String(), "application/json", bytes.NewReader(j))
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
