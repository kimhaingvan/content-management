package main

import (
	"content-management/app/media"
	"content-management/app/order"
	"content-management/cmd/content-server/build"
	"content-management/core/config"
	"content-management/pkg/application"
	"content-management/pkg/integration/consul"
	minio2 "content-management/pkg/integration/minio"
	intLog "content-management/pkg/log"
	"content-management/pkg/middleware"
	intzipkin "content-management/pkg/zipkin"
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
		if s != nil {
			err := s.Shutdown(context.Background())
			if err != nil {
				panic(err)
			}
		}
		app := build.BuildApplication(config.GetAppConfig())
		defer app.DB.Close()

		tracer, err := intzipkin.NewTracer()
		if err != nil {
			intLog.Fatal(err, nil, nil)
		}

		// Middlewares
		app.Router.Use(zipkinhttp.NewServerMiddleware(
			tracer,
			zipkinhttp.SpanName("Request")),
		)
		app.Router.Use(middleware.ZipkinMiddleware)
		app.Router.Use(middleware.APILoggingMiddleware)

		configureApp(config.GetAppConfig(), app)
		go func() {
			s = &http.Server{
				Addr:    fmt.Sprintf(":%v", config.GetAppConfig().ServerPort),
				Handler: app.NewServeMux(),
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
	// Thi???t l???p m???t channel ????? l???ng nghe t??n hi???u d???ng t??? h??? ??i???u h??nh,
	// ??? ????y ch??ng ta l??u ?? 2 t??n hi???u (signal) l?? SIGINT v?? SIGTERM
	signal.Notify(signChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	<-signChan

	// Thi???t l???p m???t kho???n th???i gian (Timeout) ????? d???ng ho??n to??n ???ng d???ng v?? ????ng t???t c??? k???t n???i.
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

func configureApp(cfg config.Config, app *application.Application) {
	minioClient := minio2.New(&minio2.Config{
		Endpoint:        cfg.Minio.Endpoint,
		SecretAccessKey: cfg.Minio.SecretAccessKey,
		AccessKey:       cfg.Minio.AccessKey,
		BucketName:      cfg.Minio.BucketName,
	})
	app.Use(
		media.New(&media.Config{
			Prefix:      "",
			MinioClient: minioClient,
		}),
	)
	app.Use(
		order.New(&order.Config{
			Prefix:      "",
			MinioClient: minioClient,
		}),
	)

}
