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

	adminservice "github.com/inst-api/poster/gen/admin_service"
	authservice "github.com/inst-api/poster/gen/auth_service"
	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/internal/config"
	"github.com/inst-api/poster/internal/mw"
	"github.com/inst-api/poster/internal/postgres"
	"github.com/inst-api/poster/internal/service"
	"github.com/inst-api/poster/internal/store"
	"github.com/inst-api/poster/internal/store/bots"
	"github.com/inst-api/poster/internal/store/tasks"
	"github.com/inst-api/poster/internal/store/users"
	"github.com/inst-api/poster/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	taskStore := tasks.NewStore(5*time.Second, dbTXFunc, txFunc, conf.Instagrapi.Hostname)
	botsStore := bots.NewStore(dbTXFunc, txFunc, conf.Instagrapi.Hostname)

	// Initialize the services.
	authServiceSvc := service.NewAuthService(userStore, conf.Security)
	tasksService := service.NewTasksService(authServiceSvc, taskStore)

	conn, err := grpc.DialContext(
		ctx,
		conf.Listen.ParserURL,
		grpc.WithUnaryInterceptor(mw.UnaryClientLog()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatalf(ctx, "failed to connect to parser: %v", err)
	}

	adminsService := service.NewAdminService(authServiceSvc, userStore, botsStore, conn)

	// potentially running in different processes.
	authServiceEndpoints := authservice.NewEndpoints(authServiceSvc)
	tasksServiceEndpoints := tasksservice.NewEndpoints(tasksService)
	adminsServiceEndpoints := adminservice.NewEndpoints(adminsService)

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
		adminsServiceEndpoints,
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
