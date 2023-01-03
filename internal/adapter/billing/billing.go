package billing

import (
	"fmt"
	"os"
)

var (
	value      = int32(49)
	errBadPath = fmt.Errorf("bad file path")
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

func (b *Billing) Parse() error {
	by, err := os.ReadFile(b.Path)
	if err != nil {
		fmt.Println(errBadPath, err)
		return errBadPath
	}

	s := string(by)
	for i, v := range s {
		if v == value && i == 2 {
			b.Data.CheckoutPage = true
		}
		if v == value && i == 3 {
			b.Data.FraudControl = true
		}
		if v == value && i == 4 {
			b.Data.Recurring = true
		}
		if v == value && i == 5 {
			b.Data.Payout = true
		}
		if v == value && i == 6 {
			b.Data.Purchase = true
		}
		if v == value && i == 7 {
			b.Data.CreateCustomer = true
		}
	}
	return nil
}
