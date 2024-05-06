// Copyright 2024 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package configs

import (
	"fmt"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/dbs"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
)

type SimpleDBConfig struct {
	User       string   `yaml:"user"`
	Password   string   `yaml:"password"`
	Database   string   `yaml:"database"`
	Host       string   `yaml:"host"`
	BoolFields []string `yaml:"boolFields,omitempty"`
}

func (this *SimpleDBConfig) GenerateOldConfig(targetFile string) error {
	var dbConfig = &dbs.DBConfig{
		Driver: "mysql",
		Dsn:    url.QueryEscape(this.User) + ":" + url.QueryEscape(this.Password) + "@tcp(" + this.Host + ")/" + url.PathEscape(this.Database) + "?charset=utf8mb4&timeout=30s&multiStatements=true",
		Prefix: "edge",
	}
	dbConfig.Models.Package = "internal/db/models"

	var config = &dbs.Config{
		DBs: map[string]*dbs.DBConfig{
			Tea.Env: dbConfig,
		},
	}
	config.Default.DB = Tea.Env
	config.Fields = map[string][]string{
		"bool": this.BoolFields,
	}

	oldConfigYAML, encodeErr := yaml.Marshal(config)
	if encodeErr != nil {
		return encodeErr
	}

	err := os.WriteFile(targetFile, oldConfigYAML, 0666)
	if err != nil {
		return fmt.Errorf("create database config file failed: %w", err)
	}

	return nil
}
