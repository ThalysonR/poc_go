package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
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

func (c *configObservable) Subscribe(configObj interface{}, cb func(error)) (*ConfigSubscription, error) {
	subscription := &ConfigSubscription{
		id:        uuid.New(),
		configObj: configObj,
		cb:        cb,
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
	viper.SetConfigType("yml")
	if runMode == "LOCAL" {
		err := c.localConfig()
		if err != nil {
			return err
		}
		defer viper.WatchConfig()
	} else {
		err := c.remoteConfig(runMode, serviceName)
		if err != nil {
			return err
		}
		defer viper.WatchRemoteConfig()
	}
	viper.OnConfigChange(func(in fsnotify.Event) {
		c.subscriptions.Range(func(key, value interface{}) bool {
			sub := value.(*ConfigSubscription)
			sub.cb(viper.Unmarshal(sub.configObj))
			return true
		})
	})
	return nil
}

type ConfigSubscription struct {
	id        uuid.UUID
	configObj interface{}
	cb        func(error)
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

func (c *configObservable) remoteConfig(runMode, serviceName string) error {
	viper.AddRemoteProvider("etcd", "http://127.0.0.1:4001", fmt.Sprintf("/config/%s/%s.properties.yml", strings.ToLower(serviceName), strings.ToLower(runMode)))
	return viper.ReadRemoteConfig()
}
