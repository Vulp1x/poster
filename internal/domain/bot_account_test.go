package domain

import (
	"strings"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestBotAccount_assignHeaders(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		init    func(t minimock.Tester) *BotAccount
		inspect func(r *BotAccount, t *testing.T) // inspects *Bot after execution of assignHeaders

		args func(t minimock.Tester) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) // use for more precise error evaluation
	}{
		{
			name: "OK",
			init: func(t minimock.Tester) *BotAccount {
				return &BotAccount{}
			},
			args: func(t minimock.Tester) args {
				return args{
					input: "X-MID=;IG-U-DS-USER-ID=55063899557;IG-U-RUR=ODN,55063899557,1693496495:01f73e106e7e6c02e0414f5a6787745fad80bff6af73b01eee0e15b7e5c186447d6a8d62;Authorization=Bearer IGT:2:eyJkc191c2VyX2lkIjoiNTUwNjM4OTk1NTciLCJzZXNzaW9uaWQiOiI1NTA2Mzg5OTU1NyUzQUN0RGRybU1wek8zMDBiJTNBMyUzQUFZZnhld2dLaVVzU25WekFZZjhoSUFqSkJTMkUyeGI4empYSUotZkZfdyJ9;X-IG-WWW-Claim=hmac.AR2dDsO3wL_piE7dQKKv-ZjEwYU0vo-nxZ0hRuMFby-L0fFY",
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			tArgs := tt.args(mc)
			receiver := tt.init(mc)

			err := receiver.assignHeaders(tArgs.input)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func BenchmarkParse(b *testing.B) {
	var testString = "michellemagana598:fMS7ZbA7Uu|Instagram 248.0.0.17.109 Android (29/10; 540dpi; 1440x2400; LGE; LG-P690; gelato_tmb-sk; qcom; ru-RU; 239490569)|android-0d735e1f4db26782;fab80e64-2b3f-44c8-8916-703e6b7a91de;d23ccbb6-ca3b-4fe5-8b23-fd0163ba0ce5;c7f9fc1c-cdff-4962-a57d-125a99e81545|X-MID=;IG-U-DS-USER-ID=55063899557;IG-U-RUR=ODN,55063899557,1693496495:01f73e106e7e6c02e0414f5a6787745fad80bff6af73b01eee0e15b7e5c186447d6a8d62;Authorization=Bearer IGT:2:eyJkc191c2VyX2lkIjoiNTUwNjM4OTk1NTciLCJzZXNzaW9uaWQiOiI1NTA2Mzg5OTU1NyUzQUN0RGRybU1wek8zMDBiJTNBMyUzQUFZZnhld2dLaVVzU25WekFZZjhoSUFqSkJTMkUyeGI4empYSUotZkZfdyJ9;X-IG-WWW-Claim=hmac.AR2dDsO3wL_piE7dQKKv-ZjEwYU0vo-nxZ0hRuMFby-L0fFY"
	inputFields := strings.Split(testString, "|")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bot := new(BotAccount)
		bot.Parse(inputFields)
	}
}
