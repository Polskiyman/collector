package mms

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"collector/pkg/country"
	"collector/pkg/provider"
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
	content, err := m.GetContent()
	if err != nil {
		return err
	}

	var d Mms
	err = json.Unmarshal(content, &d)
	if err != nil {
		return err
	}

	m.filterResponse(d.Data)

	return nil
}

func (m *Mms) GetContent() ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, m.Url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("not success")
	}

	content, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (m *Mms) filterResponse(data []MmsData) {
	for _, d := range data {
		if _, ok := country.ByCode(d.Country); !ok {
			fmt.Println(country.ErrInvalidCountry)
			continue
		}
		if !provider.IsValidMmsProvider(d.Provider) {
			fmt.Println(provider.ErrInvalidProvider)
			continue
		}
		m.Data = append(m.Data, d)
	}
}
