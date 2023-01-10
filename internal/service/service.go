package service

import (
	"sort"
	"sync"

	"collector/internal/adapter/billing"
	"collector/internal/adapter/email"
	"collector/internal/adapter/incident"
	"collector/internal/adapter/mms"
	"collector/internal/adapter/sms"
	"collector/internal/adapter/voiceCall"
	"collector/pkg/country"
)

type Collector interface {
	GetSystemData() (string, error)
}

type ResultT struct {
	Status bool       `json:"status"` // true, если все этапы сбора данных прошли успешно, false во всех остальных случаях
	Data   ResultSetT `json:"data"`   // заполнен, если все этапы сбора данных прошли успешно, nil во всех остальных случаях
	Error  string     `json:"error"`  // пустая строка если все этапы сбора данных прошли успешно, в случае ошибки заполнено текстом ошибки (детали ниже)
	sync.Mutex
}

func (r *ResultT) statusError(err error) {
	r.Lock()
	r.Status = false
	r.Error += err.Error() + "; "
	r.Unlock()
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
	sms       *sms.Sms
	mms       *mms.Mms
	voiceCall *voiceCall.VoiceCall
}

func New(smsPath, mmsUrl, viceCallPath string) *collector {
	return &collector{
		sms:       sms.New(smsPath),
		mms:       mms.New(mmsUrl),
		voiceCall: voiceCall.New(viceCallPath),
	}
}

func (c *collector) GetSystemData() (res ResultT) {
	res.Data = ResultSetT{
		SMS:       make([][]sms.SMSData, 2, 2),
		MMS:       make([][]mms.MmsData, 2, 2),
		VoiceCall: make([]voiceCall.VoiceCallData, 0),
	}

	var wg sync.WaitGroup
	wg.Add(2)
	res.Status = true

	go c.getSmsData(&wg, &res)

	go c.getMmsData(&wg, &res)

	go c.getVoiceCallData(&wg, &res)

	wg.Wait()

	return
}

func (c *collector) getSmsData(wg *sync.WaitGroup, res *ResultT) {
	defer wg.Done()

	if !res.Status {
		return
	}

	err := c.sms.Parse()
	if err != nil {
		res.statusError(err)
		return
	}

	for i, v := range c.sms.Data {
		v.Country, _ = country.ByCode(v.Country)
		c.sms.Data[i].Country = v.Country
	}

	res.Data.SMS[0] = append(res.Data.SMS[0], c.sms.Data...)
	sort.SliceStable(res.Data.SMS[0], func(i, j int) bool {
		return res.Data.SMS[0][i].Country < res.Data.SMS[0][j].Country
	})

	res.Data.SMS[1] = append(res.Data.SMS[1], c.sms.Data...)
	sort.SliceStable(res.Data.SMS[1], func(i, j int) bool {
		return res.Data.SMS[1][i].Provider < res.Data.SMS[1][j].Provider
	})
	return
}

func (c *collector) getMmsData(wg *sync.WaitGroup, res *ResultT) {
	defer wg.Done()

	if !res.Status {
		return
	}

	err := c.mms.Fetch()
	if err != nil {
		res.statusError(err)
		return
	}

	for i, v := range c.mms.Data {
		v.Country, _ = country.ByCode(v.Country)
		c.mms.Data[i].Country = v.Country
	}

	res.Data.MMS[0] = append(res.Data.MMS[0], c.mms.Data...)
	sort.SliceStable(res.Data.MMS[0], func(i, j int) bool {
		return res.Data.MMS[0][i].Country < res.Data.MMS[0][j].Country
	})

	res.Data.MMS[1] = append(res.Data.MMS[1], c.mms.Data...)
	sort.SliceStable(res.Data.MMS[1], func(i, j int) bool {
		return res.Data.MMS[1][i].Provider < res.Data.MMS[1][j].Provider
	})
	return
}

func (c *collector) getVoiceCallData(wg *sync.WaitGroup, res *ResultT) {
	defer wg.Done()

	if !res.Status {
		return
	}

	err := c.voiceCall.Parse()
	if err != nil {
		res.statusError(err)
		return
	}

	res.Data.VoiceCall = c.voiceCall.Data
	return
}
