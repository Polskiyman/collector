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

	err = country.IsValid(fields[0])
	if err != nil {
		return
	}

	err = provider.IsValidProviderVoiceCall(fields[3])
	if err != nil {
		return
	}

	connectionStability, err := strconv.ParseFloat(fields[4], 32)
	if err != nil {
		return
	}
	ttfb, err := strconv.Atoi(fields[5])
	if err != nil {
		return
	}
	voicePurity, err := strconv.Atoi(fields[6])
	if err != nil {
		return
	}
	medianOfCallsTime, err := strconv.Atoi(fields[7])
	if err != nil {
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
