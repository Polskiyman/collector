package service

import (
	"collector/internal/adapter/billing"
	"collector/internal/adapter/email"
	"collector/internal/adapter/incident"
	"collector/internal/adapter/mms"
	"collector/internal/adapter/sms"
	"collector/internal/adapter/voiceCall"
	"sort"
)

type Collector interface {
	GetSystemData() (string, error)
}

type ResultT struct {
	Status bool       `json:"status"` // true, если все этапы сбора данных прошли успешно, false во всех остальных случаях
	Data   ResultSetT `json:"data"`   // заполнен, если все этапы сбора данных прошли успешно, nil во всех остальных случаях
	Error  string     `json:"error"`  // пустая строка если все этапы сбора данных прошли успешно, в случае ошибки заполнено текстом ошибки (детали ниже)
}

type ResultSetT struct {
	SMS       [][]sms.SMSData                `json:"sms"`
	MMS       [][]mms.MmsData                `json:"mms"`
	VoiceCall []voiceCall.VoiceCallData      `json:"voice_call"`
	Email     map[string][][]email.EmailData `json:"email"`
	Billing   billing.BillingData            `json:"billing"`
	Support   []int                          `json:"support"`
	Incidents []incident.IncidentData        `json:"incident"`
}

type collector struct {
	sms *sms.Sms
}

func New(path string) *collector {
	return &collector{
		sms: sms.New(path),
	}
}

func (c *collector) GetSystemData() (res ResultT) {
	res.Data = ResultSetT{
		SMS: make([][]sms.SMSData, 2),
	}
	var err error
	err = c.sms.Parse()
	if err != nil {
		return ResultT{
			Status: false,
			Data:   ResultSetT{},
			Error:  err.Error(),
		}
	}
	sort.SliceStable(c.sms.Data, func(i, j int) bool {
		return c.sms.Data[i].Country < c.sms.Data[j].Country
	})

	res.Data.SMS[0] = c.sms.Data

	temp := make([]sms.SMSData, 2)
	copy(temp, c.sms.Data)

	sort.SliceStable(temp, func(i, j int) bool {
		return temp[i].Provider < temp[j].Provider
	})

	res.Data.SMS[1] = temp
	return res
}
