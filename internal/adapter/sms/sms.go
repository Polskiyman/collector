package sms

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	errBadPath   = fmt.Errorf("bad file path")
	errLenFields = fmt.Errorf("line not contains 4 fields")
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
	// open file
	f, err := os.Open(s.Path)
	if err != nil {
		fmt.Println(errBadPath, err)
		return errBadPath
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
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

		// do something with read line
		d, err := createSMSData(rec)
		if err != nil {
			fmt.Println(err)
			continue
		}

		s.Data = append(s.Data, d)
	}

	return nil
}

// TODO: типизировать все ошибки при обработке полей строки файлы (см errBadPath = fmt.Errorf("bad file path")
func createSMSData(line []string) (res SMSData, err error) {
	// TODO: check len of line
	fields := strings.Split(line[0], ";")

	if len(fields) != 4 {
		err = errLenFields
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
