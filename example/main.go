package main

import (
	"fmt"
	"github.com/Sjhzjxc/go_config"
	"log"
)

func main() {
	appConfig := map[string]interface{}{}
	configs := map[string]interface{}{
		"App": &appConfig,
	}
	vipConfig := go_config.VipConfig{
		ConfigName: "config",
		ConfigPath: "./example",
		ConfigType: "yml",
		ConfigRun:  go_config.DefaultConfigReloadRun,
		Configs:    configs,
	}
	Config, err := go_config.NewConfig(vipConfig)
	if err != nil {
		log.Panicln(err.Error())
	}
	for key, value := range appConfig {
		fmt.Println(key, value)
	}
	appConfig["runmode"] = "dev"
	err = Config.Save()
	if err != nil {
		log.Panicln(err.Error())
	}
	for key, value := range appConfig {
		fmt.Println(key, value)
	}

}
