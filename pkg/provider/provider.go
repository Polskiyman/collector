package provider

var providerMap = map[string]struct{}{
	"Topolo": {},
	"Rond":   {},
	"Kildy":  {},
}

func IsValid(provider string) bool {
	_, ok := providerMap[provider]
	return ok
}
