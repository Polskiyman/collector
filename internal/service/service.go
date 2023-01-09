package service

import (
	"collector/pkg/country"
	"sort"

	"collector/internal/adapter/billing"
	"collector/internal/adapter/email"
	"collector/internal/adapter/incident"
	"collector/internal/adapter/mms"
	"collector/internal/adapter/sms"
	"collector/internal/adapter/voiceCall"
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
	mms *mms.Mms
}

func New(path, mmsUrl string) *collector {
	return &collector{
		sms: sms.New(path),
		mms: mms.New(mmsUrl),
	}
}

func (c *collector) GetSystemData() (res ResultT) {
	res.Data = ResultSetT{
		SMS: make([][]sms.SMSData, 2, 2),
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
	// tody c.smc.Data - заменит коды стран на страны

	for i, v := range c.sms.Data {
		v.Country, _ = country.ByCode(v.Country)
		c.sms.Data[i] = v
	}

	res.Data.SMS[0] = append(res.Data.SMS[0], c.sms.Data...)

	sort.SliceStable(res.Data.SMS[0], func(i, j int) bool {
		return res.Data.SMS[0][i].Country < res.Data.SMS[0][j].Country
	})

	res.Data.SMS[1] = append(res.Data.SMS[1], c.sms.Data...)

	sort.SliceStable(res.Data.SMS[1], func(i, j int) bool {
		return res.Data.SMS[1][i].Provider < res.Data.SMS[1][j].Provider
	})

	return res
}
