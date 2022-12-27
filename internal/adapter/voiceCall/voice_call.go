package voiceCall

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
	errLenFields = fmt.Errorf("line not contains 8 fields")
	errEmptyLine = fmt.Errorf("line is empty")
)

type VoiceCall struct {
	Data []VoiceCallData
	Path string
}

type VoiceCallData struct {
	Country             string
	Bandwidth           string
	ResponseTime        string
	Provider            string
	ConnectionStability float32
	TTFB                int
	VoicePurity         int
	MedianOfCallsTime   int
}

func New(path string) *VoiceCall {
	return &VoiceCall{
		Data: make([]VoiceCallData, 0),
		Path: path,
	}
}

func (s *VoiceCall) Parse() error {
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

		d, err := createVoiceCallData(rec)
		if err != nil {
			fmt.Println(err)
			continue
		}

		s.Data = append(s.Data, d)
	}

	return nil
}

func createVoiceCallData(line []string) (res VoiceCallData, err error) {
	if len(line) < 1 {
		err = errEmptyLine
		return
	}
	fields := strings.Split(line[0], ";")

	if len(fields) != 8 {
		err = errLenFields
		return
	}

	ok := country.IsValid(fields[0])
	if !ok {
		err = country.ErrInvalidCountry
		return
	}

	ok = provider.IsValidProviderVoiceCall(fields[3])
	if !ok {
		err = provider.ErrInvalidProvider
		return
	}

	connectionStability, err := strconv.ParseFloat(fields[4], 32)
	if err != nil {
		err = fmt.Errorf("can't parse ConnectionStability field, field[4]=%s, err: %s", fields[4], err.Error())
		return
	}
	ttfb, err := strconv.Atoi(fields[5])
	if err != nil {
		err = fmt.Errorf("can't parse TTFB field, field[5]=%s, err: %s", fields[5], err.Error())
		return
	}
	voicePurity, err := strconv.Atoi(fields[6])
	if err != nil {
		err = fmt.Errorf("can't parse VoicePurity field, field[6]=%s, err: %s", fields[6], err.Error())
		return
	}
	medianOfCallsTime, err := strconv.Atoi(fields[7])
	if err != nil {
		err = fmt.Errorf("can't parse MedianOfCallsTime field, field[7]=%s, err: %s", fields[7], err.Error())
		return
	}

	res = VoiceCallData{
		Country:             fields[0],
		Bandwidth:           fields[1],
		ResponseTime:        fields[2],
		Provider:            fields[3],
		ConnectionStability: float32(connectionStability),
		TTFB:                ttfb,
		VoicePurity:         voicePurity,
		MedianOfCallsTime:   medianOfCallsTime,
	}
	return
}
