package main

import (
	"content-management/app/media"
	"content-management/app/order"
	"content-management/cmd/content-server/build"
	"content-management/cmd/content-server/config"
	"content-management/pkg/application"
	intzipkin "content-management/pkg/zipkin"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/opentracing/opentracing-go"

	"gopkg.in/yaml.v2"

	consulAPI "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

func init() {
	//go intLog.LogStashRegister()
}

func main() {
	var cfgCh = make(chan config.Config, 1)
	watcher, err := registerWatcher("key", os.Getenv("CONSUL_CONFIG_KEY_VALUE"))
	defer watcher.Stop()

	watcher.Handler = func(index uint64, data interface{}) {
		if pair, ok := data.(*consulAPI.KVPair); ok {
			var cfg config.Config
			err = yaml.Unmarshal(pair.Value, &cfg)
			if err != nil {
				panic(err)
			}
			cfgCh <- cfg
		}
	}
	go func() {
		// Chạy goroutine chỗ này để đảm bảo việc hàm này luôn chạy bất đồng bộ
		if err = watcher.Run(os.Getenv("CONSUL_IP_ADDRESS") + `:` + os.Getenv("CONSUL_PORT")); err != nil {
			log.Fatal(err)
		}
	}()

	var s *http.Server
	for {
		cfg := <-cfgCh
		if s != nil {
			err = s.Shutdown(context.Background())
			if err != nil {
				panic(err)
			}
		}
		app, err := build.InitApp(cfg)
		if err != nil {
			panic(err)
		}

		reporter := intzipkin.NewReporter(os.Getenv("ZIPKIN_URL"))
		defer reporter.Close()

		err = intzipkin.InitTracer(reporter)
		if err != nil {
			log.Fatalf("unable to init tracer: %+v\n", err)
		}
		// Middlewares
		app.App.Router.Use(middlewareZipkin)

		configureApp(&cfg, app.App)

		go func() {
			s = &http.Server{
				Addr:           fmt.Sprintf(":%v", cfg.Port),
				Handler:        app.App.NewServeMux(),
				MaxHeaderBytes: 1 << 16,
				ReadTimeout:    10 * time.Second,
				WriteTimeout:   10 * time.Second,
			}
			fmt.Println("HTTP server content management listening on %v", s.Addr)
			err = s.ListenAndServe()
			switch err {
			case nil, http.ErrServerClosed:
				err = nil
			default:
				fmt.Errorf("HTTP server content management error: %v", err)
			}
			shutdownGracefully(s)
		}()
	}
}

func middlewareZipkin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := intzipkin.Tracer.StartSpan(r.Method + " " + r.URL.String())
		defer span.Finish()
		span.SetTag("Method", r.Method)
		span.SetTag("URL", r.URL.String())
		ctx := opentracing.ContextWithSpan(r.Context(), span)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func registerWatcher(key string, valueOfKey string) (watcher *watch.Plan, error error) {
	var params = make(map[string]interface{})
	params["type"] = key
	params["key"] = valueOfKey
	watcher, err := watch.Parse(params)
	if err != nil {
		return nil, err
	}

	return watcher, nil
}

func shutdownGracefully(s *http.Server) {
	signChan := make(chan os.Signal, 1)
	// Thiết lập một channel để lắng nghe tín hiệu dừng từ hệ điều hành,
	// ở đây chúng ta lưu ý 2 tín hiệu (signal) là SIGINT và SIGTERM
	signal.Notify(signChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	<-signChan

	// Thiết lập một khoản thời gian (Timeout) để dừng hoàn toàn ứng dụng và đóng tất cả kết nối.
	timeWait := 15 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeWait)
	defer func() {
		log.Println("Close another connection")
		cancel()
	}()

	if err := s.Shutdown(ctx); err == context.DeadlineExceeded {
		log.Print("Halted active connections")
	}
	close(signChan)
}

func configureApp(cfg *config.Config, app *application.Application) {
	app.Use(
		order.New(&order.Config{
			Prefix:             "",
			AwsS3Region:        cfg.S3.AwsS3Region,
			AwsS3Bucket:        cfg.S3.AwsS3Bucket,
			AwsAccessKey:       cfg.S3.AwsAccessKey,
			AwsSecretAccessKey: cfg.S3.AwsSecretAccessKey,
		}),
	)
	app.Use(
		media.New(&media.Config{
			Prefix:          "",
			Endpoint:        cfg.CMCClound.Endpoint,
			SecretAccessKey: cfg.CMCClound.SecretAccessKey,
			AccessKey:       cfg.CMCClound.AccessKey,
			BucketName:      cfg.CMCClound.BucketName,
		}),
	)
}
