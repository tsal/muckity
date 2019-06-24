package muckity

import (
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"sync"
)

// not exported as this should be super-simple to implement if you don't want to use muckity.yml.
type muckityConfig struct {
	config *viper.Viper
}

func (c *muckityConfig) Name() string {
	return "config"
}
func (c *muckityConfig) Type() string {
	return "muckity:config"
}

func (c *muckityConfig) Get(k string) interface{} {
	var v interface{}
	v = c.config.Get(k)
	return v
}

func (c *muckityConfig) Set(k string, v interface{}) {
	c.config.Set(k, v)
}

func (c *muckityConfig) BindEnv(input ...string) error {
	var err error
	err = c.config.BindEnv(input...)
	return err
}

func (c muckityConfig) Dump() string {
	conf := c.config.AllSettings()
	bs, err := yaml.Marshal(conf)
	if err != nil {
		panic(err)
	}
	return string(bs)
}

var _ Config = &muckityConfig{}
var _ System = &muckityConfig{}

var instance *muckityConfig

var once sync.Once

func newConfig(_ ...interface{}) *muckityConfig {
	var mc muckityConfig
	var err error
	mc.config = viper.New()
	mc.config.SetConfigName("muckity")
	mc.config.AddConfigPath("/etc/muckity")
	mc.config.AddConfigPath("$HOME/.config/muckity")
	mc.config.AddConfigPath(".")
	mc.config.SetEnvPrefix("muckity")
	err = mc.config.ReadInConfig()
	if err != nil {
		panic(err)
	}
	mc.config.WatchConfig() // TODO: see if there is a way to implement this with a websocket
	return &mc
}