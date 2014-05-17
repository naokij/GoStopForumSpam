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

type SpamData struct {
	SearchData
	Evidence string `json:"evidence"`
}

type SearchResponseItem struct {
	Appears    int     `json:"appears"`
	Frequency  int     `json:"frequency"`
	LastSeen   Time    `json:"lastseen"`
	Confidence float64 `json:"confidence"`
}

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) error {
	var err error
	var tmpTime time.Time
	trimedData := data[1 : len(data)-1]
	tmpTime, err = time.Parse(`2006-01-02 15:04:05`, string(trimedData))
	*t = Time(tmpTime)
	return err
}

func (t *Time) ToStdTime() time.Time {
	stdTime := time.Time(*t)
	return stdTime
}

type SearchResponse struct {
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
	if response.Error != "" {
		return response, errors.New(response.Error)
	}
	return response, nil
}

func (c *Client) Add(spamData SpamData) (err error) {
	if spamData.Email == "" || spamData.Ip == "" || spamData.Username == "" || spamData.Evidence == "" {
		return errors.New("stopforumspam.Add error: spamData not complete")
	}
	postValues := url.Values{}
	postValues.Add("ip_addr", spamData.Ip)
	postValues.Add("email", spamData.Email)
	postValues.Add("username", spamData.Username)
	postValues.Add("evidence", spamData.Evidence)
	postValues.Add("api_key", c.apiKey)
	var resp *http.Response
	resp, err = http.PostForm(ADD_ENDPOINT, postValues)
	if err != nil {
		return err
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("stopforumspam.Add error: %d %s", resp.StatusCode, body)
		return err
	}
	return nil
}
