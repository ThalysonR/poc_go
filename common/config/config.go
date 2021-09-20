package config

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"github.com/thalysonr/poc_go/common/utils"
)

var (
	configInstance *configObservable
	lock           = &sync.Mutex{}
)

type configObservable struct {
	subscriptions sync.Map
}

func GetConfigObservable() (*configObservable, error) {
	lock.Lock()
	defer lock.Unlock()

	if configInstance == nil {
		configInstance = &configObservable{
			subscriptions: sync.Map{},
		}
		err := configInstance.setup()
		if err != nil {
			return nil, err
		}
	}
	return configInstance, nil
}

func (c *configObservable) Subscribe(configObj interface{}, cb func(func(cfgObj interface{}) error)) (*ConfigSubscription, error) {
	subscription := &ConfigSubscription{
		id: uuid.New(),
		cb: cb,
	}
	err := viper.Unmarshal(configObj)
	if err != nil {
		return nil, err
	}

	configInstance.subscriptions.Store(subscription.id, subscription)
	return subscription, nil
}

func (c *configObservable) setup() error {
	runMode := os.Getenv("RUN_MODE")
	serviceName := os.Getenv("SERVICE_NAME")
	etcdAddr := os.Getenv("ETCD_SERVER")
	viper.SetConfigType("yml")
	if runMode == "LOCAL" {
		err := c.localConfig()
		if err != nil {
			return err
		}
		defer viper.WatchConfig()
	} else {
		err := c.remoteConfig(etcdAddr, runMode, serviceName)
		if err != nil {
			return err
		}
		defer viper.WatchRemoteConfig()
	}
	debouncer := utils.NewDebouncer(time.Second)
	viper.OnConfigChange(func(in fsnotify.Event) {
		debouncer(func() {
			c.subscriptions.Range(func(key, value interface{}) bool {
				sub := value.(*ConfigSubscription)
				sub.cb(func(cfgObj interface{}) error {
					return viper.Unmarshal(cfgObj)
				})
				return true
			})
		})
	})
	return nil
}

type ConfigSubscription struct {
	id uuid.UUID
	cb func(func(cfgObj interface{}) error)
}

func (c *ConfigSubscription) Unsubscribe() {
	configInstance.subscriptions.Delete(c.id)
}

////////////////////////////////////////////////////////////////////////////////
///////                       AUXILIARY FUNCTIONS                        ///////
////////////////////////////////////////////////////////////////////////////////

func (c *configObservable) localConfig() error {
	viper.SetConfigName("local.properties")
	viper.AddConfigPath(".")
	return viper.ReadInConfig()
}

func (c *configObservable) remoteConfig(etcdServer, runMode, serviceName string) error {
	viper.AddRemoteProvider("etcd", etcdServer, fmt.Sprintf("/config/%s/%s.properties.yml", strings.ToLower(serviceName), strings.ToLower(runMode)))
	return viper.ReadRemoteConfig()
}
