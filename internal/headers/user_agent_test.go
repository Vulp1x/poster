package headers

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"testing"
)

func TestParseUserAgent(t *testing.T) {
	var tests = []struct {
		input   string
		want    DeviceSettings
		wantErr error
	}{
		{
			input: "Instagram 248.0.0.17.109 Android (29/10; 540dpi; 1440x2400; LGE; LG-P690; gelato_tmb-sk; qcom; ru-RU; 239490569)",
			want: DeviceSettings{
				AppVersion:     "248.0.0.17.109",
				AndroidVersion: 29,
				AndroidRelease: "10",
				Dpi:            "540dpi",
				Resolution:     "1440x2400",
				Manufacturer:   "LGE",
				Device:         "LG-P690",
				Model:          "gelato_tmb-sk",
				Cpu:            "qcom",
				VersionCode:    "239490569",
			},
		},
		{
			input: "Instagram 248.0.0.17.109 Android (28/9; 408dpi; 1440x3028; Xiaomi/xiaomi; Redmi Note 5A Prime; ugg; qcom; ru-RU; 239490569)",
			want: DeviceSettings{
				AppVersion:     "248.0.0.17.109",
				AndroidVersion: 28,
				AndroidRelease: "9",
				Dpi:            "408dpi",
				Resolution:     "1440x3028",
				Manufacturer:   "Xiaomi/xiaomi",
				Device:         "Redmi Note 5A Prime",
				Model:          "ugg",
				Cpu:            "qcom",
				VersionCode:    "239490569",
			},
		},
		{
			input: "Instagram 248.0.0.17.109 Android (29/10; 480dpi; 1440x2756; huawei/orange; orangeyumo; hwg740-l00; ru-RU; 239490569)",
			want: DeviceSettings{
				AppVersion:     "248.0.0.17.109",
				AndroidVersion: 29,
				AndroidRelease: "10",
				Dpi:            "480dpi",
				Resolution:     "1440x2756",
				Manufacturer:   "huawei/orange",
				Device:         "orangeyumo",
				Model:          "hwg740-l00",
				Cpu:            "qcom",
				VersionCode:    "239490569",
			},
		},
	}

	for _, test := range tests {
		got, err := NewDeviceSettings(test.input)
		if !errors.Is(err, test.wantErr) {
			t.Errorf("wanted error: %v, \ngot: %v", test.wantErr, err)
		}

		if test.want != got {
			t.Errorf("wanted '%#v', \ngot '%#v'", test.want, got)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewDeviceSettings("Instagram 248.0.0.17.109 Android (29/10; 540dpi; 1440x2400; LGE; LG-P690; gelato_tmb-sk; qcom; ru-RU; 239490569)")
	}
}

func TestName(t *testing.T) {
	pubKey := []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsaKXOX+2oFx9VN5C5SoB
LcoASF6SiR7wEB25UCa3shwBY0qkyg8DjaqQxDHeOfkphdRyeMd+35L1DRCFycCi
uzeDH5rGAPODRC4h/RwenaOHXmJvx77uS4UxQ7m/8ITP2WVsnG2TiMDKlW748PHq
GVsFNNv6yNGqgVhfL9P2rZ2TPM91lI0JYblR4ibuvS6Z0S4keKSHMh5N9pYutI6l
x5iU5r3pn48XMUmXJSHf60wQQVoeVbtvMrZipzRbzsB6X8di5YK4YhJB+1fewph1
VKFDY6kbndzDmj4hitLDxMvo9C4eBPVqPAapAIdVJKW3oCplucXg/y+mKKyET/Fh
xQIDAQAB
-----END PUBLIC KEY-----
`)

	block, _ := pem.Decode(pubKey)
	rsaKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		t.Errorf("failed to parse key: %v", err)
	}

	fmt.Printf("%#+v\n", rsaKey)
}
