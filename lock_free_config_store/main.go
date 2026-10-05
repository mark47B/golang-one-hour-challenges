package main

import ( 
//"runtime"
// "time"
//  "fmt"
// "context"
"log"
"errors"
// _ "net/http/pprof"
// "net/http"
	"sync"
	"sync/atomic"
)
type Config struct {
    Host string
    Port int
}

type ConfigStore struct {
    cfg atomic.Pointer[Config]

	mu sync.Mutex
}

func NewConfigStore(cfg *Config) *ConfigStore {
	if cfg == nil {
		return nil
	}
	c := *cfg
	s := &ConfigStore{}
	s.cfg.Store(&c)
	return s
}

var (
	PanicLoader = errors.New("loader has been panic")
)

func (s *ConfigStore) Get() *Config {
	return s.cfg.Load()
}
func (s *ConfigStore) Load(loader func() (*Config, error)) (e error) {
	defer func(){
		if r := recover(); r != nil{
			log.Println("loader panic")
			e = PanicLoader
		}
	}()
	s.mu.Lock()
	defer s.mu.Unlock()
	newCfg, err := loader()
	if err != nil {
		return err
	}
	if newCfg == nil {
		return errors.New("loader returned nil config")
	}
	cpyCfg := *newCfg
	s.cfg.Store(&cpyCfg)
	return nil
}

func main(){
	cfgStore := NewConfigStore(&Config{Host: "localhost", Port:8080})
	_ = cfgStore
}
