// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package configs

import (
	"encoding/json"
	"github.com/iwind/TeaGo/Tea"
	"io/ioutil"
)

var plusConfigFile = "plus.cache.json"

type PlusConfig struct {
	IsPlus bool `json:"isPlus"`
}

func ReadPlusConfig() *PlusConfig {
	data, err := ioutil.ReadFile(Tea.ConfigFile(plusConfigFile))
	if err != nil {
		return &PlusConfig{IsPlus: false}
	}
	var config = &PlusConfig{IsPlus: false}
	err = json.Unmarshal(data, config)
	if err != nil {
		return config
	}
	return config
}

func WritePlusConfig(config *PlusConfig) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(Tea.ConfigFile(plusConfigFile), configJSON, 0777)
	if err != nil {
		return err
	}
	return nil
}
