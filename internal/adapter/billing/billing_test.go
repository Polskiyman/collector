package billing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBilling_Parse(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantRes BillingData
	}{
		{
			name: "BillingParse",
			data: "billing_data.txt",
			wantRes: BillingData{
				CreateCustomer: true,
				Purchase:       true,
				Payout:         false,
				Recurring:      false,
				FraudControl:   true,
				CheckoutPage:   false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b Billing
			b.Path = tt.data
			err := b.Parse()
			assert.Nil(t, err, nil)
			assert.Equal(t, tt.wantRes, b.Data)
		})
	}
}
