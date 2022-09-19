package reader

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/headers"
)

func TestName(t *testing.T) {
	f, err := os.Open("testdata/100 акк.txt")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	defer f.Close()

	bots, errs := ParseUsersList(context.Background(), f)
	if len(errs) != 0 {
		t.Fatalf("got %d errors: %v", len(errs), errs)
	}

	if len(bots) != 100 {
		t.Fatalf("got %d bots instead of 100", len(bots))
	}

}

func TestParseBotAccount(t *testing.T) {
	var testString = "michellemagana598:fMS7ZbA7Uu|Instagram 248.0.0.17.109 Android (29/10; 540dpi; 1440x2400; LGE; LG-P690; gelato_tmb-sk; qcom; ru-RU; 239490569)|android-0d735e1f4db26782;fab80e64-2b3f-44c8-8916-703e6b7a91de;d23ccbb6-ca3b-4fe5-8b23-fd0163ba0ce5;c7f9fc1c-cdff-4962-a57d-125a99e81545|X-MID=;IG-U-DS-USER-ID=55063899557;IG-U-RUR=ODN,55063899557,1693496495:01f73e106e7e6c02e0414f5a6787745fad80bff6af73b01eee0e15b7e5c186447d6a8d62;Authorization=Bearer IGT:2:eyJkc191c2VyX2lkIjoiNTUwNjM4OTk1NTciLCJzZXNzaW9uaWQiOiI1NTA2Mzg5OTU1NyUzQUN0RGRybU1wek8zMDBiJTNBMyUzQUFZZnhld2dLaVVzU25WekFZZjhoSUFqSkJTMkUyeGI4empYSUotZkZfdyJ9;X-IG-WWW-Claim=hmac.AR2dDsO3wL_piE7dQKKv-ZjEwYU0vo-nxZ0hRuMFby-L0fFY"

	bot := domain.BotAccount{}

	err := bot.Parse(strings.Split(testString, "|"))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	wantedBot := domain.BotAccount{
		Bot: dbmodel.BotAccount{
			Username:  "michellemagana598",
			Password:  "fMS7ZbA7Uu",
			UserAgent: "Instagram 248.0.0.17.109 Android (29/10; 540dpi; 1440x2400; LGE; LG-P690; gelato_tmb-sk; qcom; ru-RU; 239490569)",
			DeviceData: headers.DeviceSettings{
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
			Session: headers.Session{
				DeviceID:      "android-0d735e1f4db26782",
				UUID:          uuid.MustParse("fab80e64-2b3f-44c8-8916-703e6b7a91de"),
				PhoneID:       uuid.MustParse("d23ccbb6-ca3b-4fe5-8b23-fd0163ba0ce5"),
				AdvertisingID: uuid.MustParse("c7f9fc1c-cdff-4962-a57d-125a99e81545"),
			},
			Headers: headers.Base{
				Mid:           "",
				DsUserID:      "55063899557",
				Rur:           "ODN,55063899557,1693496495:01f73e106e7e6c02e0414f5a6787745fad80bff6af73b01eee0e15b7e5c186447d6a8d62",
				Authorization: "Bearer IGT:2:eyJkc191c2VyX2lkIjoiNTUwNjM4OTk1NTciLCJzZXNzaW9uaWQiOiI1NTA2Mzg5OTU1NyUzQUN0RGRybU1wek8zMDBiJTNBMyUzQUFZZnhld2dLaVVzU25WekFZZjhoSUFqSkJTMkUyeGI4empYSUotZkZfdyJ9",
				WWWClaim:      "hmac.AR2dDsO3wL_piE7dQKKv-ZjEwYU0vo-nxZ0hRuMFby-L0fFY",
				AuthData: headers.AuthorizationData{
					DsUserID:  "55063899557",
					SessionID: "55063899557:CtDdrmMpzO300b:3:AYfxewgKiUsSnVzAYf8hIAjJBS2E2xb8zjXIJ-fF_w",
					CSRFToken: "",
				},
			},
		},
	}

	var opts = []cmp.Option{
		cmpopts.IgnoreFields(domain.BotAccount{}, "ID", "TaskID", "StartedAt", "CreatedAt"),
	}

	if !cmp.Equal(bot, wantedBot, opts...) {
		t.Fatalf("bot account assert failed, diff: \n%s", cmp.Diff(bot, wantedBot, opts...))
	}
}
