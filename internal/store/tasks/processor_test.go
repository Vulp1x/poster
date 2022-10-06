package tasks

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/csv"
	"fmt"
	"image/jpeg"
	"testing"

	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/images"
	"github.com/inst-api/poster/internal/instagrapi"
)

//go:embed testdata/test_cat.jpeg
var catPhotoBytes []byte

//go:embed testdata/test_photo.webp
var webpPhotoBytes []byte

//go:embed testdata/100.txt
var accountsBytes []byte

func Test_worker_createPost(t *testing.T) {
	// os.Setenv("GODEBUG", "http2client=0")

	csvReader := csv.NewReader(bytes.NewBuffer(accountsBytes))
	csvReader.Comma = '|'
	firstLine, err := csvReader.Read()
	if err != nil {
		t.Fatalf("failed to read first line from bots: %v", err)
	}

	bots, uploadErrs := domain.ParseBotAccounts([]*tasksservice.BotAccountRecord{{Record: firstLine, LineNumber: 0}})
	if len(uploadErrs) != 0 {
		t.Fatalf("got %d errors when parsing bots: %+v", len(uploadErrs), uploadErrs)
	}

	bots_ := []domain.BotAccount(bots)
	// bots_[0].ResProxy = &dbmodel.Proxy{
	// 	Host: "192.168.1.19",
	// 	Port: 9090,
	// 	// Login: "dmitrijkholodkov7815",
	// 	// Pass:  "21e49b",
	// }

	// 109.248.7.220:10475:dmitrijkholodkov7815:21e49b
	// Instagram 252.0.0.17.111 Android (28/9; 320dpi; 720x1402; samsung; SM-S102DL; a10e; exynos7885; en_IN; 397702079)
	bots_[0].ResProxy = &dbmodel.Proxy{
		Host:  "109.248.7.220",
		Port:  10475,
		Login: "dmitrijkholodkov7815",
		Pass:  "21e49b",
	}

	fmt.Printf("using bot: %+v\n", bots_[0])

	img, err := jpeg.Decode(bytes.NewReader(catPhotoBytes))
	if err != nil {
		t.Fatalf("failed to decode image: %v", err)
	}

	fmt.Println(img.Bounds(), img.ColorModel(), len(catPhotoBytes)/1280)

	type fields struct {
		tasksQueue     chan *domain.BotWithTargets
		dbtxf          dbmodel.DBTXFunc
		cli            instagrapiClient
		generator      images.Generator
		processorIndex int64
		captionFormat  string
	}
	type args struct {
		ctx         context.Context
		botAccount  domain.BotAccount
		targetUsers []dbmodel.TargetUser
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			fields: fields{
				cli:       instagrapi.NewClient(),
				generator: images.NewDummyGenerator(catPhotoBytes),
			},
			args: args{
				ctx:        context.Background(),
				botAccount: bots_[0],
				targetUsers: []dbmodel.TargetUser{
					{
						Username: "velo__andrew",
						UserID:   2182798860,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &worker{
				botsQueue:      tt.fields.tasksQueue,
				dbtxf:          tt.fields.dbtxf,
				cli:            tt.fields.cli,
				generator:      tt.fields.generator,
				processorIndex: tt.fields.processorIndex,
				captionFormat:  tt.fields.captionFormat,
			}
			if err := p.createPost(tt.args.ctx, tt.args.botAccount, tt.args.targetUsers); (err != nil) != tt.wantErr {
				t.Errorf("createPost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
