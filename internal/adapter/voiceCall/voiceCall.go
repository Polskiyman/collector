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
	errBadPath         = fmt.Errorf("bad file path")
	errLenFields       = fmt.Errorf("line not contains 8 fields")
	errEmptyLine       = fmt.Errorf("line is empty")
	errInvalidCountry  = fmt.Errorf("incorrect country code")
	errInvalidProvider = fmt.Errorf("incorrect provider")
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
		err = errInvalidCountry
		return
	}

	ok = provider.IsValid(fields[3])
	if !ok {
		err = errInvalidProvider
		return
	}
	connectionStability, err := strconv.ParseFloat(fields[4], 32)
	ttfb, err := strconv.Atoi(fields[5])
	voicePurity, err := strconv.Atoi(fields[6])
	medianOfCallsTime, err := strconv.Atoi(fields[7])

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
