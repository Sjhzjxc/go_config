package go_config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
)

type Config struct {
	vp      *viper.Viper
	configs map[string]interface{}
}

type VipConfig struct {
	ConfigName string
	ConfigPath string
	ConfigType string
	ConfigRun  func(*Config) func(in fsnotify.Event)
	Configs    map[string]interface{}
}

func (s *Config) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	if _, ok := s.configs[k]; !ok {
		s.configs[k] = v
	}
	return nil
}

func (s *Config) LoadAllSection() error {
	for k, v := range s.configs {
		err := s.ReadSection(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewConfig(config VipConfig) (*Config, error) {
	if config.ConfigName == "" {
		config.ConfigName = "config"
	}
	if config.ConfigPath == "" {
		config.ConfigPath = "."
	}
	if config.ConfigType == "" {
		config.ConfigType = "yml"
	}
	vp := viper.New()
	vp.SetConfigName(config.ConfigName)
	vp.AddConfigPath(config.ConfigPath)
	vp.SetConfigType(config.ConfigType)
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	vp.WatchConfig()
	s := &Config{vp: vp, configs: config.Configs}
	err = s.LoadConfig()
	if config.ConfigRun != nil {
		vp.OnConfigChange(config.ConfigRun(s))
	}
	return s, err
}

func DefaultConfig(configs map[string]interface{}) (*Config, error) {
	vipConfig := VipConfig{
		ConfigName: "config",
		ConfigPath: ".",
		ConfigType: "yml",
		Configs:    configs,
		ConfigRun:  DefaultConfigReloadRun,
	}
	return NewConfig(vipConfig)
}

func DefaultConfigReloadRun(s *Config) func(in fsnotify.Event) {
	return func(in fsnotify.Event) {
		err := s.LoadConfig()
		if err != nil {
			log.Fatalln(err.Error() + " config重载失败")
		}
	}
}

func (s *Config) LoadConfig() error {
	err := s.LoadAllSection()
	if err != nil {
		return err
	}
	return nil
}

func (s *Config) Save() error {
	for key, value := range s.configs {
		s.vp.Set(key, value)
	}
	err := s.vp.WriteConfig()
	if err != nil {
		return err
	}
	return err
}

func (s *Config) SaveAsMove() error {
	for key, value := range s.configs {
		s.vp.Set(key, value)
	}
	filename := s.vp.ConfigFileUsed()
	fileExt := path.Ext(filename)
	tempName := "temp" + fileExt
	err := s.vp.WriteConfigAs(tempName)
	if err != nil {
		return err
	}
	err = os.Rename(tempName, filename)
	return err
}

func (s *Config) SaveAs(name string) error {
	for key, value := range s.configs {
		s.vp.Set(key, value)
	}
	err := s.vp.WriteConfigAs(name)
	if err != nil {
		return err
	}
	return err
}
