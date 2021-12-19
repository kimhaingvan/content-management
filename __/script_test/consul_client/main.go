package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/consul/api/watch"

	consulAPI "github.com/hashicorp/consul/api"
	"github.com/k0kubun/pp"
)

const (
	CONSUL_IP   = "127.0.0.1"
	CONSUL_PORT = "8500"
)

func main() {
	config := consulAPI.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	consulClient, err := consulAPI.NewClient(config)
	if err != nil {
		pp.Println("watcher.Run")
		panic(err)
	}
	kv := consulClient.KV()
	pair, _, err := kv.Get("redis/config/minconns", nil)
	if err != nil {
		panic(err)
	}
	value := string(pair.Value[:])
	var dataEnvironment testmodel
	err = json.Unmarshal(pair.Value, &dataEnvironment)

	var ch = make(chan int, 1)
	watcher, err := RegisterWatcher("key", "redis/config/minconns")

	defer watcher.Stop()
	watcher.Handler = func(index uint64, data interface{}) {
		if pair, ok := data.(*consulAPI.KVPair); ok {
			value = string(pair.Value[:])
			ch <- 1
		}
	}
	go func() {
		// Chạy goroutine chỗ này để đảm bảo việc hàm này luôn chạy bất đồng bộ
		if err = watcher.Run(CONSUL_IP + `:` + CONSUL_PORT); err != nil {
			log.Fatal(err)
		}
	}()

	for {
		<-ch
		go func() {
			err = json.Unmarshal([]byte(value), &dataEnvironment)
			pp.Println(value)
			// check Health tại đây
			s := &http.Server{
				Handler: nil,
				Addr:    ":8080",
			}
			s.ListenAndServe()
			pp.Println("ts.ListenAndServe()")
			shutdownGracefully(s)
		}()
	}

	//urlCheckHealth := "http://localhost:8500/v1/health/checks/web"
	//responseHealth, _, err := callApi.CallAPI("GET", urlCheckHealth, nil)
	//if err != nil {
	//	panic(err)
	//}
	//responseMapHealth := StringArrayJsonToMapStringInterface(responseHealth)
	//for j := 0; j < len(responseMapHealth); j++ {
	//	nameService := fmt.Sprintf(`%v`, responseMapHealth[j]["Node"])
	//	status := fmt.Sprintf(`%v`, responseMapHealth[j]["Status"])
	//	if "inmacs-MBP" == nameService {
	//		if status == "passing" {
	//			return address
	//		}
	//	}
	//}
}

func shutdownGracefully(s *http.Server) {
	signChan := make(chan os.Signal, 1)
	// Thiết lập một channel để lắng nghe tín hiệu dừng từ hệ điều hành,
	// ở đây chúng ta lưu ý 2 tín hiệu (signal) là SIGINT và SIGTERM
	signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)
	<-signChan

	// Thiết lập một khoản thời gian (Timeout) để dừng hoàn toàn ứng dụng và đóng tất cả kết nối.
	timeWait := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeWait)
	defer func() {
		log.Println("Close another connection")
		cancel()
	}()

	if err := s.Shutdown(ctx); err == context.DeadlineExceeded {
		log.Print("Halted active connections")
	}
	pp.Println(s.Shutdown(ctx))
	close(signChan)
}

type testmodel struct {
	Abc int32 `json:"abc"`
}

func StringArrayJsonToMapStringInterface(jsonString string) []map[string]interface{} {
	bytes := []byte(jsonString)
	var jsonMap []map[string]interface{}
	json.Unmarshal(bytes, &jsonMap)
	return jsonMap
}

func RegisterWatcher(key string, valueOfKey string) (watcher *watch.Plan, error error) {
	var params = make(map[string]interface{})
	params["type"] = key
	params["key"] = valueOfKey
	watcher, err := watch.Parse(params)
	if err != nil {
		return nil, err
	}
	return watcher, nil
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	pp.Println("srv.Shutdown(ctx)")
	srv.Shutdown(ctx)
	pp.Println("exit na")
	os.Exit(0)
}
