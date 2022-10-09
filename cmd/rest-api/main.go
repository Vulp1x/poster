package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	authservice "github.com/inst-api/poster/gen/auth_service"
	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/internal/config"
	"github.com/inst-api/poster/internal/postgres"
	"github.com/inst-api/poster/internal/service"
	"github.com/inst-api/poster/internal/store"
	"github.com/inst-api/poster/internal/store/tasks"
	"github.com/inst-api/poster/internal/store/users"
	"github.com/inst-api/poster/pkg/logger"
)

func main() {
	loc := time.FixedZone("Moscow", 3*60*60)
	time.Local = loc

	// Define command line flags, add any other flag required to configure the
	// service.
	var (
		debugFlag = flag.Bool("debug", false, "Log request and response bodies")
	)
	flag.Parse()

	configMode := os.Getenv("CONFIG_MODE")

	conf := &config.Config{}

	fmt.Println(configMode, *debugFlag)

	err := conf.ParseConfiguration(configMode)
	if err != nil {
		log.Fatal("Failed to parse configuration: ", err)
	}

	err = logger.InitLogger(conf.Logger)
	if err != nil {
		log.Fatal("Failed to create logger: ", err)

		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	dbTXFunc, err := postgres.NewDBTxFunc(ctx, conf.Postgres)
	if err != nil {
		logger.Fatalf(ctx, "Failed to connect to database: %v", err)
	}

	txFunc, err := postgres.NewTxFunc(ctx, conf.Postgres)
	if err != nil {
		logger.Fatalf(ctx, "Failed to connect to transaction's database: %v", err)
	}

	userStore := users.NewStore(store.WithDBTXFunc(dbTXFunc), store.WithTxFunc(txFunc))
	taskStore := tasks.NewStore(5*time.Second, dbTXFunc, txFunc)

	// Initialize the services.
	authServiceSvc := service.NewAuthService(userStore, conf.Security)
	tasksService := service.NewTasksService(authServiceSvc, taskStore)

	// potentially running in different processes.
	authServiceEndpoints := authservice.NewEndpoints(authServiceSvc)
	tasksServiceEndpoints := tasksservice.NewEndpoints(tasksService)

	// Create channel used by both the signal handler and server goroutines
	// to notify the main goroutine when to stop the server.
	errc := make(chan error)

	// Setup interrupt handler. This optional step configures the process so
	// that SIGINT and SIGTERM signals cause the services to stop gracefully.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	var wg sync.WaitGroup

	// Start the servers and send errors (if any) to the error channel.
	handleHTTPServer(
		ctx,
		conf.Listen.BindIP,
		conf.Listen.Port,
		conf.Listen.KeyFile,
		conf.Listen.CertFile,
		authServiceEndpoints,
		tasksServiceEndpoints,
		&wg,
		errc,
		*debugFlag,
	)
	// Wait for signal.
	logger.Infof(ctx, "exiting from main: (%v)", <-errc)

	// Send cancellation signal to the goroutines.
	cancel()

	wg.Wait()
	logger.Info(ctx, "exited")
}
