package provider

var providerMap = map[string]struct{}{
	"Topolo":           {},
	"Rond":             {},
	"Kildy":            {},
	"TransparentCalls": {},
	"E-Voice":          {},
	"JustPhone":        {},
}

func IsValid(provider string) bool {
	_, ok := providerMap[provider]
	return ok
}
