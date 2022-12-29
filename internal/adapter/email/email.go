package email

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"collector/pkg/country"
	"collector/pkg/provider"
)

var (
	errBadPath   = fmt.Errorf("bad file path")
	errLenFields = fmt.Errorf("line not contains 3 fields")
	errEmptyLine = fmt.Errorf("line is empty")
)

type Email struct {
	Data []EmailData
	Path string
}

type EmailData struct {
	Country      string
	Provider     string
	DeliveryTime int
}

func New(path string) *Email {
	return &Email{
		Data: make([]EmailData, 0),
		Path: path,
	}
}

func (s *Email) Parse() error {
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

		d, err := createEmailData(rec)
		if err != nil {
			fmt.Println(err)
			continue
		}

		s.Data = append(s.Data, d)
	}

	return nil
}

func createEmailData(line []string) (res EmailData, err error) {
	if len(line) < 1 {
		err = errEmptyLine
		return
	}
	fields := strings.Split(line[0], ";")

	if len(fields) != 3 {
		err = errLenFields
		return
	}

	ok := country.IsValid(fields[0])
	if !ok {
		err = country.ErrInvalidCountry
		return
	}

	ok = provider.IsValidEmailProvider(fields[1])
	if !ok {
		err = provider.ErrInvalidProvider
		return
	}
	deliveryTime, err := strconv.Atoi(fields[2])
	if err != nil {
		err = fmt.Errorf("can't parse DeliveryTime field, field[2]=%s, err: %s", fields[2], err.Error())
		return
	}
	res = EmailData{
		Country:      fields[0],
		Provider:     fields[1],
		DeliveryTime: deliveryTime,
	}
	return
}
