package shutdown

import (
	"context"
	"io"
	"os"
	"os/signal"

	"github.com/inst-api/poster/pkg/logger"
)

// Gracefully ждет одного из сигналов из списка signalsToWaitFor,
// и в случае получения поочередно вызывает метод Close() для всех элементов из itemsToClose
func Gracefully(signalsToWaitFor []os.Signal, itemsToClose ...io.Closer) {
	ctx := context.Background()
	// Handle common process-killing signalsToWaitFor, so we can gracefully shut down:
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, signalsToWaitFor...)
	sig := <-sigc
	logger.Infof(ctx, "Caught signal %s: shutting down.", sig)

	for _, closer := range itemsToClose {
		err := closer.Close()
		if err != nil {
			logger.Errorf(ctx, "failed to close %v: %v", closer, err)
		}
	}
}
