package stopforumspam

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	SEARCH_ENDPOINT = "http://www.stopforumspam.com/api"
	ADD_ENDPOINT    = "http://www.stopforumspam.com/add.php"
)

type Client struct {
	apiKey string
}

type SearchData struct {
	Ip       string
	Email    string
	Username string
}

type SearchResponseItem struct {
	Appears    int       `json:"appears"`
	Frequency  int       `json:"frequency"`
	LastSeen   time.Time `json:"lastseen"`
	Confidence float64   `json:"confidence"`
}

type SearchResponse struct {
	Success  int                 `json:"success"`
	Ip       *SearchResponseItem `json:"ip"`
	Email    *SearchResponseItem `json:"email"`
	Username *SearchResponseItem `json:"username"`
	Error    string              `json:"error"`
}

func New(key string) *Client {
	return &Client{apiKey: key}
}

func (c *Client) Search(searchData SearchData) (response SearchResponse, err error) {
	queryValues := url.Values{}
	if searchData.Ip != "" {
		queryValues.Add("ip", searchData.Ip)
	}
	if searchData.Email != "" {
		queryValues.Add("email", searchData.Email)
	}
	if searchData.Username != "" {
		queryValues.Add("email", searchData.Username)
	}
	if len(queryValues) == 0 {
		return response, errors.New("stopforumspam.Search error no searchData provided")
	}
	queryValues.Add("f", "json")
	queryString := queryValues.Encode()
	resp, err := http.Get(SEARCH_ENDPOINT + "?" + queryString)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("stopforumspam.Search error: %d %s", resp.StatusCode, body)
		return response, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}
