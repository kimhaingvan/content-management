package main

import (
	"content-management/app/order"
	"content-management/cmd/content-server/build"
	"content-management/core/config"
	"content-management/pkg/application"
	intLog "content-management/pkg/log"
	"content-management/pkg/middleware"
	intzipkin "content-management/pkg/zipkin"
	"content-management/registry"
	"content-management/thirdparty/consul"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"

	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
	go intLog.RegisterLogStash(os.Getenv("LOGSTASH_IP"), os.Getenv("LOGSTASH_PORT"), os.Getenv("APPLICATION_NAME"))
}

func main() {
	var cfgCh = make(chan config.Config, 1)
	watcher := consul.RegisterConsulWatcher(cfgCh, &consul.Config{
		ApplicationName: os.Getenv("APPLICATION_NAME"),
		ConsulAclToken:  os.Getenv("CONSUL_ACL_TOKEN"),
		ConsulIP:        os.Getenv("CONSUL_IP"),
		ConsulPort:      os.Getenv("CONSUL_PORT"),
	})
	defer watcher.Stop()

	var s *http.Server
	for {
		appCfg := <-cfgCh
		config.SetAppConfig(appCfg)
		intLog.SetLoglevel(appCfg.Log.Level)
		if s != nil {
			err := s.Shutdown(context.Background())
			if err != nil {
				panic(err)
			}
		}
		r, err := registry.New(config.GetAppConfig())
		if err != nil {
			log.Fatal(err)
		}
		defer r.DB.GormDB.Close()

		app := build.BuildApplication(r)

		tracer, err := intzipkin.NewTracer()
		if err != nil {
			intLog.Fatal(err, nil, nil)
		}

		// Middlewares
		app.Router.Use(zipkinhttp.NewServerMiddleware(tracer, zipkinhttp.SpanName("Request")))
		app.Router.Use(middleware.ZipkinMiddleware)
		app.Router.Use(middleware.APILoggingMiddleware)

		configureApp(app)
		go func() {
			s = &http.Server{
				Addr:    fmt.Sprintf(":%v", config.GetAppConfig().ServerPort),
				Handler: app,
			}

			fmt.Println("HTTP server content management listening on", config.GetAppConfig().ServerPort)
			err = s.ListenAndServe()
			switch err {
			case nil, http.ErrServerClosed:
				err = nil
			default:
				fmt.Errorf("HTTP server content management error: %v", err)
			}
		}()
	}
}

func shutdownGracefully(s *http.Server) {
	signChan := make(chan os.Signal, 1)
	// Thiết lập một channel để lắng nghe tín hiệu dừng từ hệ điều hành,
	// ở đây chúng ta lưu ý 2 tín hiệu (signal) là SIGINT và SIGTERM
	signal.Notify(signChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	<-signChan

	// Thiết lập một khoản thời gian (Timeout) để dừng hoàn toàn ứng dụng và đóng tất cả kết nối.
	timeWait := 3 * time.Second
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

func configureApp(app *application.Application) {
	// Order
	orderServer := order.NewOrderServer(&order.Config{
		Prefix:      "/order",
		Application: app,
	})

	application.Use(orderServer)
	app.Router.Mount("/admin", app.Admin.NewServeMux("/admin"))
}
