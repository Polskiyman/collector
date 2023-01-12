package service

import (
	"collector/pkg"
	"sort"
	"sync"

	"collector/internal/adapter/billing"
	"collector/internal/adapter/email"
	"collector/internal/adapter/incident"
	"collector/internal/adapter/mms"
	"collector/internal/adapter/sms"
	"collector/internal/adapter/support"
	"collector/internal/adapter/voiceCall"
	"collector/pkg/country"
)

type CollectorInterface interface {
	GetSystemData() (res ResultT)
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

type Collector struct {
	sms       *sms.Sms
	mms       *mms.Mms
	voiceCall *voiceCall.VoiceCall
	email     *email.Email
	billing   *billing.Billing
	incident  *incident.Incident
	support   *support.Support
}

func New(config pkg.Config) *Collector {
	return &Collector{
		sms:       sms.New(config.SmsPath),
		mms:       mms.New(config.MmsUrl),
		voiceCall: voiceCall.New(config.VoiceCallPath),
		email:     email.New(config.EmailPath),
		billing:   billing.New(config.BillingPath),
		incident:  incident.New(config.IncidentUrl),
		support:   support.New(config.SupportUrl),
	}
}

func (c *Collector) GetSystemData() (res ResultT) {
	res.Data = ResultSetT{
		SMS:       make([][]sms.SMSData, 2, 2),
		MMS:       make([][]mms.MmsData, 2, 2),
		VoiceCall: make([]voiceCall.VoiceCallData, 0),
		Email:     make(map[string][][]email.EmailData),
		Incidents: make([]incident.IncidentData, 0),
		Support:   make([]int, 2),
	}

	var wg sync.WaitGroup
	wg.Add(7)
	res.Status = true

	go c.getSmsData(&wg, &res)

	go c.getMmsData(&wg, &res)

	go c.getVoiceCallData(&wg, &res)

	go c.getEmailData(&wg, &res)

	go c.getBillingData(&wg, &res)

	go c.getIncidentData(&wg, &res)

	go c.getSupportData(&wg, &res)

	wg.Wait()

	return
}

func (c *Collector) getSmsData(wg *sync.WaitGroup, res *ResultT) {
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

func (c *Collector) getMmsData(wg *sync.WaitGroup, res *ResultT) {
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

func (c *Collector) getVoiceCallData(wg *sync.WaitGroup, res *ResultT) {
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

func (c *Collector) getEmailData(wg *sync.WaitGroup, res *ResultT) {
	defer wg.Done()

	if !res.Status {
		return
	}

	err := c.email.Parse()
	if err != nil {
		res.statusError(err)
		return
	}
	var t []email.EmailData
	t = append(t, c.email.Data...)
	sort.SliceStable(t, func(i, j int) bool {
		return t[i].DeliveryTime < t[j].DeliveryTime
	})
	for _, v := range t {
		_, ok := res.Data.Email[v.Country]
		if !ok {
			res.Data.Email[v.Country] = make([][]email.EmailData, 2)
		}
		res.Data.Email[v.Country][0] = append(res.Data.Email[v.Country][0], v)
	}

	for i, _ := range res.Data.Email {
		// can use for both slices, they are equal
		l, n := len(res.Data.Email[i][0]), 3
		if l < n {
			n = l
		}
		res.Data.Email[i][1] = append(res.Data.Email[i][0][l-n : l])
		res.Data.Email[i][0] = append(res.Data.Email[i][0][:n])
	}
	return
}

func (c *Collector) getBillingData(wg *sync.WaitGroup, res *ResultT) {
	defer wg.Done()

	if !res.Status {
		return
	}
	err := c.billing.Parse()
	if err != nil {
		res.statusError(err)
		return
	}
	res.Data.Billing = c.billing.Data
	return
}

func (c *Collector) getIncidentData(wg *sync.WaitGroup, res *ResultT) {
	defer wg.Done()

	if !res.Status {
		return
	}
	err := c.incident.Fetch()
	if err != nil {
		res.statusError(err)
		return
	}
	sort.SliceStable(c.incident.Data, func(i, j int) bool {
		return c.incident.Data[i].Status < c.incident.Data[j].Status
	})
	res.Data.Incidents = c.incident.Data
	return
}

func (c *Collector) getSupportData(wg *sync.WaitGroup, res *ResultT) {
	defer wg.Done()

	if !res.Status {
		return
	}
	err := c.support.Fetch()
	if err != nil {
		res.statusError(err)
		return
	}
	const amountOfWorkPerHour = 3
	activTickets := 0
	load := 1
	for _, v := range c.support.Data {
		activTickets += v.ActiveTickets
	}

	if 9 <= activTickets && activTickets <= 16 {
		load = 2
	}
	if activTickets > 16 {
		load = 3
	}
	res.Data.Support[0] = load
	waitingTime := activTickets * amountOfWorkPerHour
	res.Data.Support[1] = waitingTime
	return
}
