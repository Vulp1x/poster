package service

import (
	"context"
	_ "embed"
	"encoding/base64"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/internal/config"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/postgres"
	"github.com/inst-api/poster/internal/store/tasks"
	"github.com/inst-api/poster/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

func TestName(t *testing.T) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte("admin0"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	fmt.Println(string(hashPass))

	err = bcrypt.CompareHashAndPassword(hashPass, []byte("admin0"))
	fmt.Println(err)
}

//go:embed testdata/test_cat.jpeg
var image []byte

func TestCreateDraftTask(t *testing.T) {
	conf := &config.Config{}

	err := conf.ParseConfiguration(true)
	if err != nil {
		log.Fatal("Failed to parse configuration: ", err)
	}

	// loggerConf := conf.Logger
	//
	// conf.Logger =
	// err = logger.InitLogger()
	// if err != nil {
	// 	log.Fatal("Failed to create logger: ", err)
	//
	// 	return
	// }

	title, textTemplate := "title", "text_template"

	managerID := uuid.MustParse("b79d3bff-54c0-4f4a-b165-fa832e787648")
	ctx := AddUserIDToContext(context.Background(), managerID)
	// ctx, cancel := context.WithCancel(context.Background())

	dbTXFunc, err := postgres.NewDBTxFunc(ctx, conf.Postgres)
	if err != nil {
		logger.Fatalf(ctx, "Failed to connect to database: %v", err)
	}

	store := tasks.NewStore(5*time.Second, dbTXFunc, nil)

	service := NewTasksService(nil, store)

	taskID, err := service.CreateTaskDraft(ctx, &tasksservice.CreateTaskDraftPayload{
		Title:        title,
		TextTemplate: textTemplate,
		PostImage:    base64.StdEncoding.EncodeToString(image),
	})
	if err != nil {
		t.Fatalf("failed to create task draft: %v", err)
	}

	parsedID, err := uuid.Parse(taskID)
	if err != nil {
		t.Fatalf("failed to parse task id from %s: %v", taskID, err)
	}

	q := dbmodel.New(dbTXFunc(ctx))
	task, err := q.FindTaskByID(ctx, parsedID)
	if err != nil {
		t.Fatalf("failed to find task with id %s: %v", parsedID, err)
	}

	ignoreFields := cmpopts.IgnoreFields(dbmodel.Task{}, "StartedAt", "CreatedAt", "UpdatedAt", "DeletedAt")
	expectedTask := dbmodel.Task{
		ID:           parsedID,
		ManagerID:    managerID,
		TextTemplate: textTemplate,
		Image:        image,
		Status:       1,
		Title:        title,
	}

	if !cmp.Equal(expectedTask, task, ignoreFields) {
		t.Fatalf("got unexpected task: diff: %s", cmp.Diff(expectedTask, task, ignoreFields))
	}
}
