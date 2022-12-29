package mms

import (
	"collector/pkg/country"
	"collector/pkg/provider"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Mms struct {
	Data []MmsData
	Url  string
}

type MmsData struct {
	Country      string `json:"country"`
	Provider     string `json:"provider"`
	Bandwidth    string `json:"bandwidth"`
	ResponseTime string `json:"response_time"`
}

func New(url string) *Mms {
	return &Mms{
		Data: make([]MmsData, 0),
		Url:  url,
	}
}

func (m *Mms) Fetch() error {
	url := fmt.Sprintf("%s/data", m.Url)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	request.Header.Add("Accept", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("not success")
	}

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	data := make([]MmsData, 0)
	err = json.Unmarshal(content, &data)
	if err != nil {
		return err
	}

	m.filterResponse(data)

	return nil
}

func (m *Mms) filterResponse(data []MmsData) {
	for _, d := range data {
		if !country.IsValid(d.Country) {
			fmt.Println(country.ErrInvalidCountry)
			continue
		}
		if !provider.IsValidMmsProvider(d.Provider) {
			fmt.Println(country.ErrInvalidCountry)
			continue
		}
		m.Data = append(m.Data, d)
	}
}
