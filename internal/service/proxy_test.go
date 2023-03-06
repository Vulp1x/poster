package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/postgres"
	"github.com/inst-api/poster/internal/transport"
	"github.com/inst-api/poster/pkg/logger"
)

func TestProxiesPerTask(t *testing.T) {
	conf := postgres.Configuration{}
	conf.Default()
	ctx := logger.ToContext(context.Background(), logger.Logger())
	db, err := postgres.NewConn(ctx, conf)
	if err != nil {
		t.Fatalf("failed to create db connection: %v", err)
	}

	q := dbmodel.New(db)

	// proxies, err := q.FindAllCheapProxies(ctx)
	proxies, err := q.FindCheapProxiesForTask(ctx, uuid.MustParse("54423f7a-2c0b-4d8e-b018-e1cea2c7455f"))
	if err != nil {
		t.Fatalf("failed to find proxies for task: %v", err)
	}

	inChan := make(chan *url.URL, len(proxies))
	outchan := make(chan *url.URL, len(proxies))

	logger.Infof(ctx, "got %d proxies to check", len(proxies))

	for _, proxy := range proxies {
		proxyUrl := &url.URL{
			Scheme: "http",
			User:   url.UserPassword(proxy.Login, proxy.Pass),
			Host:   fmt.Sprintf("%s:%d", proxy.Host, proxy.Port),
		}

		inChan <- proxyUrl
	}

	wg := &sync.WaitGroup{}
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go testProxies(ctx, i, inChan, outchan, wg)
	}

	close(inChan)
	wg.Wait()

	close(outchan)

	f, err := os.Create(fmt.Sprintf("best_proxies_%s", time.Now().Format("2006_01_02_15_04_05")))
	if err != nil {
		logger.Errorf(ctx, "failed to create file: %v", err)
	}

	for bestProxy := range outchan {
		logger.Infof(ctx, "proxy %s is good", bestProxy.String())
		_, err = f.WriteString(bestProxy.String())
		if err != nil {
			logger.Errorf(ctx, "failed to write proxy to file: %v", err)
		}
		f.WriteString("\n")
	}

	f.Close()
}

func testProxies(ctx context.Context, processNumber int, inputChan, outChan chan *url.URL, wg *sync.WaitGroup) {
	defer wg.Done()
	cli := transport.ProxyingHTTPClientWithTimeout(10 * time.Second)

	ctx = logger.WithKV(ctx, "process_number", processNumber)

	for proxyURL := range inputChan {
		proxyCtx := logger.WithKV(transport.ContextWithProxy(ctx, proxyURL), "proxy", proxyURL.String())

		// const testURL = "https://2ip.ru"
		const testURL = "https://www.instagram.com"
		req, err := http.NewRequestWithContext(proxyCtx, "GET", testURL, nil)
		if err != nil {
			logger.Errorf(ctx, "failed to create request: %v", err)
			continue
		}

		req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")

		resp, err := cli.Do(req)
		if err != nil {
			logger.Errorf(proxyCtx, "failed to make request: %v", err)
			continue
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Errorf(ctx, "failed to read response body: %v", err)
		}

		err = resp.Body.Close()
		if err != nil {
			logger.Errorf(ctx, "failed to close response body: %v", err)
		}

		index := bytes.Index(bodyBytes, []byte("IP:"))
		if index > 0 {
			proxyCtx = logger.WithKV(proxyCtx, "real_ip", bodyBytes[index:index+15])
		}

		logger.Info(proxyCtx, "proxy is working, adding it to best proxies array, response body len: %d", len(bodyBytes))
		outChan <- proxyURL
	}
}

func TestNOProxy(t *testing.T) {
	testIP(context.Background())
}
func testIP(ctx context.Context) {
	cli := transport.InitHTTPClient()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://2ip.ru", nil)
	if err != nil {
		logger.Errorf(ctx, "failed to create request: %v", err)
	}

	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")

	resp, err := cli.Do(req)
	if err != nil {
		logger.Errorf(ctx, "failed to make request: %v", err)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf(ctx, "failed to read response body: %v", err)
	}

	err = resp.Body.Close()
	if err != nil {
		logger.Errorf(ctx, "failed to close response body: %v", err)
	}

	logger.Info(ctx, "proxy is working, adding it to best proxies array, response body len: %d", len(bodyBytes))
}
