package incident

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Incident struct {
	Data []IncidentData
	Url  string
}

type IncidentData struct {
	Topic  string `json:"topic"`
	Status string `json:"status"`
}

func New(url string) *Incident {
	return &Incident{
		Data: make([]IncidentData, 0),
		Url:  url,
	}
}

func (i *Incident) Fetch() error {
	content, err := i.GetContent()
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &i.Data)
	if err != nil {
		return err
	}

	return nil
}

func (i *Incident) GetContent() ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, i.Url, nil)
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
