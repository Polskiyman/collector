package provider

import "fmt"

var providerVoiceCallMap = map[string]struct{}{
	"TransparentCalls": {},
	"E-Voice":          {},
	"JustPhone":        {},
}

var ErrInvalidProvider = fmt.Errorf("incorrect provider")

func IsValidProviderVoiceCall(provider string) bool {
	_, ok := providerVoiceCallMap[provider]
	return ok
}
