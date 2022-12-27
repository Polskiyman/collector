package provider

import "fmt"

var providerSMSMap = map[string]struct{}{
	"Topolo": {},
	"Rond":   {},
	"Kildy":  {},
}

var errInvalidProviderSms = fmt.Errorf("incorrect provider")

func IsValidSmaProvider(provider string) (err error) {
	if _, ok := providerSMSMap[provider]; !ok {
		err = errInvalidProviderSms
		return
	}
	return nil
}
