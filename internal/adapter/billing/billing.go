package billing

import (
	"fmt"
	"os"
)

var (
	one           = uint8('1')
	errBadPath    = fmt.Errorf("bad file path")
	errLenBilling = fmt.Errorf("not contains 8 bits")
)

type Billing struct {
	Data BillingData
	Path string
}

type BillingData struct {
	CreateCustomer bool
	Purchase       bool
	Payout         bool
	Recurring      bool
	FraudControl   bool
	CheckoutPage   bool
}

func New(path string) *Billing {
	return &Billing{
		Data: BillingData{},
		Path: path,
	}
}

func (b *Billing) Parse() error {
	content, err := os.ReadFile(b.Path)
	if err != nil {
		fmt.Println(errBadPath, err)
		return errBadPath
	}

	s := string(content)
	if len(s) != 8 {
		err = errLenBilling
		return err
	}
	b.Data.CheckoutPage = s[2] == one
	b.Data.FraudControl = s[3] == one
	b.Data.Recurring = s[4] == one
	b.Data.Payout = s[5] == one
	b.Data.Purchase = s[6] == one
	b.Data.CreateCustomer = s[7] == one
	return nil
}
