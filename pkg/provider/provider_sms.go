package provider

var providerSMSMap = map[string]struct{}{
	"Topolo": {},
	"Rond":   {},
	"Kildy":  {},
}

func IsValidSmaProvider(provider string) bool {
	_, ok := providerSMSMap[provider]
	return ok
}
