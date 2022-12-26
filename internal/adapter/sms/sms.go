package sms

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"collector/pkg/country"
	"collector/pkg/provider"
)

var (
	errBadPath         = fmt.Errorf("bad file path")
	errLenFields       = fmt.Errorf("line not contains 4 fields")
	errEmptyLine       = fmt.Errorf("line is empty")
	errInvalidCountry  = fmt.Errorf("incorrect country code")
	errInvalidProvider = fmt.Errorf("incorrect provider")
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
	f, err := os.Open(s.Path)
	if err != nil {
		fmt.Println(errBadPath, err)
		return errBadPath
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(errBadPath, err)
			return errBadPath
		}

		d, err := createSMSData(rec)
		if err != nil {
			fmt.Println(err)
			continue
		}

		s.Data = append(s.Data, d)
	}

	return nil
}

func createSMSData(line []string) (res SMSData, err error) {
	if len(line) < 1 {
		err = errEmptyLine
		return
	}
	fields := strings.Split(line[0], ";")

	if len(fields) != 4 {
		err = errLenFields
		return
	}

	ok := country.IsValid(fields[0])
	if !ok {
		err = errInvalidCountry
		return
	}

	ok = provider.IsValid(fields[3])
	if !ok {
		err = errInvalidProvider
		return
	}

	res = SMSData{
		Country:      fields[0],
		Bandwidth:    fields[1],
		ResponseTime: fields[2],
		Provider:     fields[3],
	}
	return
}
