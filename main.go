package main

import (
	"content-management/app/order"
	"content-management/cmd/content-server/build"
	"content-management/cmd/content-server/config"
	"content-management/pkg/application"
	"content-management/pkg/cmenv"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal("Error in loading config: ", err)
	}
	app, err := build.InitApp(*cfg)
	if err != nil {
		panic(err)
	}
	registerApp(cfg, app.App)
	s := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Port),
		Handler: app.App.NewServeMux(),
	}
	go func() {
		fmt.Println("HTTP server content management listening on %v", s.Addr)
		err = s.ListenAndServe()
		switch err {
		case nil, http.ErrServerClosed:
			err = nil
		default:
			fmt.Errorf("HTTP server content management error: %v", err)
		}
	}()

	// Gracefully shutdown
	shutdownGracefully(s)
}

func registerApp(cfg *config.Config, app *application.Application) {
	app.Use(order.New(&order.Config{
		Prefix:             "",
		AwsS3Region:        cfg.S3.AwsS3Region,
		AwsS3Bucket:        cfg.S3.AwsS3Bucket,
		AwsAccessKey:       cfg.S3.AwsAccessKey,
		AwsSecretAccessKey: cfg.S3.AwsSecretAccessKey,
	}))
}

func shutdownGracefully(s *http.Server) {
	signChan := make(chan os.Signal, 1)
	// Thiết lập một channel để lắng nghe tín hiệu dừng từ hệ điều hành,
	// ở đây chúng ta lưu ý 2 tín hiệu (signal) là SIGINT và SIGTERM
	signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)
	<-signChan
	// Thiết lập một khoản thời gian (Timeout) để dừng hoàn toàn ứng dụng và đóng tất cả kết nối.
	timeWait := 30 * time.Second
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

func loadConfig() (*config.Config, error) {
	config.InitFlags()
	config.ParseFlags()
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	cmenv.SetEnvironment("backend-server", cfg.Env)
	return &cfg, nil
}
