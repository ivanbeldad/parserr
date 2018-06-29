package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
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
	// APIEpisodeURL ...
	APIEpisodeURL = APIURL + "/episode"
	// StatusCompleted ...
	StatusCompleted = "Completed"
	// TrackedDownloadStatusWarning ...
	TrackedDownloadStatusWarning = "Warning"
	// MaxTime Max interval to check series and clean them
	MaxTime = time.Minute * 5
	// CheckInterval Time between requests to check if rescan is completed
	CheckInterval = time.Second * 5
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

// DeleteQueueItem ...
func DeleteQueueItem(id int) (err error) {
	u := GetURL(APIQueueURL + "/" + strconv.Itoa(id)).String()
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", u, nil)
	if err != nil {
		return
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("error deleting item from queue, status code %d", res.StatusCode)
	}
	return nil
}

// GetHistory ...
func GetHistory(page int) (history History, err error) {
	u := GetURL(APIHistoryURL)
	query := u.Query()
	query.Add("page", strconv.Itoa(page))
	query.Add("pageSize", "10")
	u.RawQuery = query.Encode()
	body, err := Get(u.String())
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &history)
	if history.PageSize == 0 {
		return history, fmt.Errorf("history fetched 0 results, no more items")
	}
	return
}

// GetEpisode ...
func GetEpisode(id int) (episode Episode, err error) {
	u := GetURL(APIEpisodeURL + "/" + strconv.Itoa(id))
	body, err := Get(u.String())
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &episode)
	return
}

// ExecuteCommand ...
func ExecuteCommand(c CommandBody) (cs CommandStatus, err error) {
	j, err := json.Marshal(c)
	if err != nil {
		return
	}
	body, err := Post(GetURL(APICommandURL).String(), bytes.NewReader(j))
	err = json.Unmarshal(body, &cs)
	return
}

// ExecuteCommandAndWait ...
func ExecuteCommandAndWait(c CommandBody) (cs CommandStatus, err error) {
	cs, err = ExecuteCommand(c)
	if err != nil {
		return
	}
	totalWait := CheckInterval
	for totalWait <= MaxTime {
		time.Sleep(CheckInterval)
		cs, err = GetCommandStatus(cs.ID)
		if err == nil {
			if cs.State == CommandStateCompleted {
				log.Printf("finished %s successfully", c.Name)
				return
			}
		}
		totalWait += CheckInterval
	}
	return cs, fmt.Errorf("timeout checking command %s, not completed", c.Name)
}

// GetCommandStatus ...
func GetCommandStatus(id int) (cs CommandStatus, err error) {
	u := GetURL(APICommandURL + "/" + strconv.Itoa(id))
	body, err := Get(u.String())
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &cs)
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
