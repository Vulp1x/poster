package workers

import (
	"reflect"
	"testing"

	"github.com/inst-api/poster/internal/dbmodel"
)

func Test_preparePostCaption(t *testing.T) {
	type args struct {
		task           dbmodel.Task
		landingAccount string
		targets        []dbmodel.TargetUser
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 []int64
	}{
		{
			name: "no photo targets, targets less than needed",
			args: args{
				task: dbmodel.Task{
					TextTemplate:        "@account\n",
					NeedPhotoTags:       false,
					TargetsPerPost:      5,
					PhotoTargetsPerPost: 0,
				},
				landingAccount: "testim",
				targets:        []dbmodel.TargetUser{{Username: "1"}, {Username: "2"}, {Username: "3"}, {Username: "4"}},
			},
			want:  "@testim\n @1 @2 @3 @4",
			want1: []int64{},
		},
		{
			name: "no photo targets, targets equal than needed",
			args: args{
				task: dbmodel.Task{
					TextTemplate:        "@account\n",
					NeedPhotoTags:       false,
					TargetsPerPost:      4,
					PhotoTargetsPerPost: 0,
				},
				landingAccount: "testim",
				targets:        []dbmodel.TargetUser{{Username: "1"}, {Username: "2"}, {Username: "3"}, {Username: "4"}},
			},
			want:  "@testim\n @1 @2 @3 @4",
			want1: []int64{},
		},

		{
			name: "no photo targets, targets more than needed",
			args: args{
				task: dbmodel.Task{
					TextTemplate:        "@account\n",
					NeedPhotoTags:       false,
					TargetsPerPost:      3,
					PhotoTargetsPerPost: 0,
				},
				landingAccount: "testim",
				targets:        []dbmodel.TargetUser{{Username: "1"}, {Username: "2"}, {Username: "3"}, {Username: "4"}},
			},
			want:  "@testim\n @1 @2 @3",
			want1: []int64{},
		},
		{
			name: "photo targets, targets less than needed",
			args: args{
				task: dbmodel.Task{
					TextTemplate:        "@account\n",
					NeedPhotoTags:       true,
					TargetsPerPost:      3,
					PhotoTargetsPerPost: 2,
				},
				landingAccount: "testim",
				targets:        []dbmodel.TargetUser{{Username: "1"}, {Username: "2"}, {Username: "3"}, {Username: "4", UserID: 4}},
			},
			want:  "@testim\n @1 @2 @3",
			want1: []int64{4},
		},
		{
			name: "photo targets, targets equal than needed",
			args: args{
				task: dbmodel.Task{
					TextTemplate:        "@account\n",
					NeedPhotoTags:       true,
					TargetsPerPost:      2,
					PhotoTargetsPerPost: 2,
				},
				landingAccount: "testim",
				targets:        []dbmodel.TargetUser{{Username: "1"}, {Username: "2"}, {Username: "3", UserID: 3}, {Username: "4", UserID: 4}},
			},
			want:  "@testim\n @1 @2",
			want1: []int64{3, 4},
		},

		{
			name: "photo targets, targets more than needed",
			args: args{
				task: dbmodel.Task{
					TextTemplate:        "@account\n",
					NeedPhotoTags:       true,
					TargetsPerPost:      1,
					PhotoTargetsPerPost: 2,
				},
				landingAccount: "testim",
				targets:        []dbmodel.TargetUser{{Username: "1"}, {Username: "2", UserID: 2}, {Username: "3", UserID: 3}, {Username: "4", UserID: 4}},
			},
			want:  "@testim\n @1",
			want1: []int64{2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := preparePostCaption(tt.args.task, tt.args.landingAccount, tt.args.targets)
			if got != tt.want {
				t.Errorf("preparePostCaption() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("preparePostCaption() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
