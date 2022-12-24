package sms

import (
	"fmt"
)

var (
	errBadPath = fmt.Errorf("bad file path")
)

// Sms TODO: think more about naming
type Sms struct {
	Data []SMSData
	Path string
}

type SMSData struct {
	Country      string
	Bandwidth    string
	ResponseTime string
	Provider     string
}

func New(path string) *Sms {
	return &Sms{
		Data: make([]SMSData, 0),
		Path: path,
	}
}

func (s *Sms) Parse() error {
	s.Data = []SMSData{{}, {}, {}}
	return nil
}
