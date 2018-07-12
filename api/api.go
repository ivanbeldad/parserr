package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
	// APIMovieURL ...
	APIMovieURL = APIURL + "/movie"
	// StatusCompleted ...
	StatusCompleted = "Completed"
	// TrackedDownloadStatusWarning ...
	TrackedDownloadStatusWarning = "Warning"
	// MaxTime Max interval to check series and clean them
	MaxTime = time.Minute * 5
	// CheckInterval Time between requests to check if rescan is completed
	CheckInterval = time.Second * 5
)

// API ..
type API struct {
	URL            string
	APIKey         string
	DownloadFolder string
	Type           string
}

// NewAPI Return an instance of an API
func NewAPI(url, apiKey, downloadFolder, apiType string) API {
	return API{
		URL:            url,
		APIKey:         apiKey,
		DownloadFolder: downloadFolder,
		Type:           apiType,
	}
}

// GetQueue ...
func (a API) GetQueue() (queue []QueueElement, err error) {
	body, err := get(a.getURL(APIQueueURL).String())
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &queue)
	return
}

// DeleteQueueItem ...
func (a API) DeleteQueueItem(id int) (err error) {
	u := a.getURL(APIQueueURL + "/" + strconv.Itoa(id)).String()
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
func (a API) GetHistory(page int) (history History, err error) {
	u := a.getURL(APIHistoryURL)
	query := u.Query()
	query.Add("page", strconv.Itoa(page))
	query.Add("pageSize", "10")
	u.RawQuery = query.Encode()
	body, err := get(u.String())
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
func (a API) GetEpisode(id int) (episode Episode, err error) {
	u := a.getURL(APIEpisodeURL + "/" + strconv.Itoa(id))
	body, err := get(u.String())
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &episode)
	return
}

// GetMovie ...
func (a API) GetMovie(id int) (movie Movie, err error) {
	u := a.getURL(APIMovieURL + "/" + strconv.Itoa(id))
	body, err := get(u.String())
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &movie)
	return
}

// ExecuteCommand ...
func (a API) ExecuteCommand(c CommandBody) (cs CommandStatus, err error) {
	j, err := json.Marshal(c)
	if err != nil {
		return
	}
	body, err := post(a.getURL(APICommandURL).String(), bytes.NewReader(j))
	err = json.Unmarshal(body, &cs)
	return
}

// ExecuteCommandAndWait ...
func (a API) ExecuteCommandAndWait(c CommandBody) (cs CommandStatus, err error) {
	cs, err = a.ExecuteCommand(c)
	if err != nil {
		return
	}
	totalWait := CheckInterval
	for totalWait <= MaxTime {
		time.Sleep(CheckInterval)
		cs, err = a.GetCommandStatus(cs.ID)
		if err == nil {
			if cs.State == CommandStateCompleted {
				log.Printf("finished %s successfully", c.Name)
				return
			}
			log.Printf("waiting response from %s", c.Name)
		}
		totalWait += CheckInterval
	}
	return cs, fmt.Errorf("timeout checking command %s, not completed", c.Name)
}

// GetCommandStatus ...
func (a API) GetCommandStatus(id int) (cs CommandStatus, err error) {
	u := a.getURL(APICommandURL + "/" + strconv.Itoa(id))
	body, err := get(u.String())
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &cs)
	return
}

// get Wrapper for http.Get. Add authentication handling automatically.
func get(u string) (body []byte, err error) {
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

// post Wrapper for http.Post. Add authentication handling automatically.
func post(u string, bodyReq io.Reader) (body []byte, err error) {
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

func (a API) getURL(path string) *url.URL {
	u := &url.URL{
		Scheme: "http",
		Host:   a.URL,
		Path:   path,
	}
	q := u.Query()
	q.Set("apikey", a.APIKey)
	u.RawQuery = q.Encode()
	return u
}
