package headers

import "github.com/google/uuid"

// Default returns default headers6 that should be set for all requests
func Default() map[string][]string {
	return map[string][]string{
		"Accept-Encoding":             {"gzip, deflate"},
		"Accept":                      {"*/"},
		"Connection":                  {"keep-alive"},
		"X-IG-App-Locale":             {"en_US"},
		"X-IG-Device-Locale":          {"en_US"},
		"X-IG-Mapped-Locale":          {"en_US"},
		"X-Pigeon-Session-Id":         {"UFS-" + uuid.NewString() + "-1"},
		"X-Pigeon-Rawclienttime":      {"1662410349.132"},
		"X-IG-Bandwidth-Speed-KBPS":   {"2612.404"},
		"X-IG-Bandwidth-TotalBytes-B": {"49827250"},
		"X-IG-Bandwidth-TotalTime-MS": {"4680"},
		"X-IG-App-Startup-Country":    {"US"},
		"X-IG-WWW-Claim":              {"0"},
		"X-Bloks-Is-Layout-RTL":       {"false"},
		"X-Bloks-Is-Panorama-Enabled": {"true"},
		// "X-IG-Timezone-Offset":        {"10800"},
		"X-IG-Connection-Type": {"WIFI"},
		"X-IG-Capabilities":    {"3brTv10="},
		"X-IG-App-ID":          {"567067343352427"}, // TODO update to use from response headers
		"Priority":             {"u=3"},
		"Accept-Language":      {"en-US"},
		"Host":                 {"i.instagram.com"},
		"X-FB-HTTP-Engine":     {"Liger"},
		"X-FB-Client-IP":       {"True"},
		"X-FB-Server-Cluster":  {"True"},
		"IG-INTENDED-USER-ID":  {"0"},
		"X-IG-Nav-Chain":       {"9MV:self_profile:2,ProfileMediaTabFragment:self_profile:3,9Xf:self_following:4"},
		"X-IG-SALT-IDS":        {"1061244479"},
		"Content-Type":         {"application/x-www-form-urlencoded; charset=UTF-8"},
		"Pragma":               {"no-cache"},
		"Cache-Control":        {"no-cache"},
	}
}
