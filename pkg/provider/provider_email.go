package provider

var providerEmailMap = map[string]struct{}{
	"Gmail":      {},
	"Yahoo":      {},
	"Hotmail":    {},
	"MSN":        {},
	"Orange":     {},
	"Comcast":    {},
	"AOL":        {},
	"Live":       {},
	"RediffMail": {},
	"GMX":        {},
	"Protonmail": {},
	"Yandex":     {},
	"Mail.ru":    {},
}

func IsValidEmailProvider(provider string) bool {
	_, ok := providerEmailMap[provider]
	return ok
}
