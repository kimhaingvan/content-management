package consul

import (
	"content-management/core/config"
	"content-management/pkg/httpreq"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	consulAPI "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

type Client struct {
	*consulAPI.Client
	cfg     *Config
	rClient *httpreq.Resty
	baseURL string
}

func New(cfg *Config) (*Client, error) {
	config := consulAPI.DefaultConfig()
	config.Token = cfg.ConsulAclToken
	config.Address = fmt.Sprintf("%v:%v", cfg.ConsulIP, cfg.ConsulPort)
	c, err := consulAPI.NewClient(config)

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	rcfg := httpreq.RestyConfig{Client: client}

	if err != nil {
		return nil, err
	}
	return &Client{
		Client:  c,
		cfg:     cfg,
		rClient: httpreq.NewResty(rcfg),
		baseURL: fmt.Sprintf("http://%v:%v", cfg.ConsulIP, cfg.ConsulPort),
	}, nil
}

type Config struct {
	ApplicationName string
	ConsulAclToken  string
	ConsulIP        string
	ConsulPort      string
}

// Kết nối với consul server
// Lấy giá trị config từ key ở consul server
func (c *Client) GetConfigFromConsulServer() (cfg *config.Config, err error) {
	kv := c.KV()
	pair, _, err := kv.Get(strings.TrimPrefix(`config/`+c.cfg.ApplicationName+`/data`, "/"), nil)
	if err != nil {
		return nil, err
	}
	if pair != nil {
		value := string(pair.Value[:])
		cfg = config.DefaultConfig()
		err = json.Unmarshal([]byte(value), &cfg)
		if err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

// Tạo một watcher đề ngóng việc thay đổi config từ consul server,
func (c *Client) RegisterWatcher(key string, valueOfKey string) (watcher *watch.Plan, error error) {
	var params = make(map[string]interface{})
	params["type"] = key
	params["key"] = valueOfKey
	watcher, err := watch.Parse(params)
	if err != nil {
		return nil, err
	}
	return watcher, nil
}

func (c *Client) RegisterServiceWatcher(key string, value string) (watcher *watch.Plan, error error) {
	var params = make(map[string]interface{})
	params["type"] = key
	params["service"] = value
	params["datacenter"] = "dc1"
	watcher, err := watch.Parse(params)
	if err != nil {
		return nil, err
	}
	return watcher, nil
}

func RegisterConsulWatcher(cfgCh chan config.Config, cfg *Config) *watch.Plan {
	consulClient, err := New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// Đăng ký một watcher - watcher này sẽ ngóng sự thay đổi từ consul server
	// khi có một sự thay đổi thì nó sẽ báo về cho consul client
	watcher, err := consulClient.RegisterWatcher("key", strings.TrimPrefix(`config/`+cfg.ApplicationName+`/data`, "/"))
	if err != nil {
		log.Fatal(err)
	}

	// Watcher handler
	watcher.Handler = func(index uint64, data interface{}) {
		// Khi consul server thay đổi cặp key/value thì lấy lại giá trị value mới, đồng thời cập nhật lại channel cfgCh
		if pair, ok := data.(*consulAPI.KVPair); ok {
			var config config.Config
			err = json.Unmarshal(pair.Value, &config)
			if err != nil {
				log.Fatal(err)
			}
			cfgCh <- config
		}
	}
	go func() {
		// Chạy goroutine chỗ này để đảm bảo việc hàm này luôn chạy bất đồng bộ
		if err = watcher.Run(cfg.ConsulIP + `:` + cfg.ConsulPort); err != nil {
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
	return watcher
}
