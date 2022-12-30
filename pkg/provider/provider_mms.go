package provider

var providerMmsMap = map[string]struct{}{
	"Topolo": {},
	"Rond":   {},
	"Kildy":  {},
}

func IsValidMmsProvider(provider string) bool {
	_, ok := providerMmsMap[provider]
	return ok
}
