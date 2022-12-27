package provider

import "fmt"

var providerVoiceCallMap = map[string]struct{}{
	"TransparentCalls": {},
	"E-Voice":          {},
	"JustPhone":        {},
}

var errInvalidProviderVoiceCall = fmt.Errorf("incorrect provider")

func IsValidProviderVoiceCall(provider string) (err error) {
	if _, ok := providerVoiceCallMap[provider]; !ok {
		err = errInvalidProviderVoiceCall
		return
	}
	return nil
}
