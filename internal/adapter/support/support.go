package support

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Support struct {
	Data []SupportData
	Url  string
}

type SupportData struct {
	Topic         string `json:"topic"`
	ActiveTickets int    `json:"active_tickets"`
}

func New(url string) *Support {
	return &Support{
		Data: make([]SupportData, 0),
		Url:  url,
	}
}

func (s *Support) Fetch() error {
	content, err := s.GetContent()
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &s.Data)
	if err != nil {
		return err
	}

	return nil
}

func (s *Support) GetContent() ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, s.Url, nil)
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
