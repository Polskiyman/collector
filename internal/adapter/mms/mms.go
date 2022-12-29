package mms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Mms struct {
	Data []MmsData
	Url  string
}

type MmsData struct {
	// TODO: add fields
	Test int `json:"test"` // TODO: remove
}

func New(url string) *Mms {
	return &Mms{
		Data: make([]MmsData, 0),
		Url:  url,
	}
}

func (m *Mms) Fetch() error {
	// TODO: create an http client

	// TODO: create an http request

	// TODO: make request and get response

	// TODO: parse response properly - and fill m.Data slice

	return nil
}

// TODO: remove
func (m *Mms) TestFunc() (MmsData, error) {
	url := fmt.Sprintf("%s/data", m.Url)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return MmsData{}, err
	}
	request.Header.Add("Accept", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return MmsData{}, err
	}

	if response.StatusCode != 200 {
		return MmsData{}, fmt.Errorf("not success")
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return MmsData{}, err
	}

	v := MmsData{}
	err = json.Unmarshal(content, &v)
	if err != nil {
		return MmsData{}, err
	}
	return v, nil
}
