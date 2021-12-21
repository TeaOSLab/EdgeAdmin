// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package tasks

import (
	"encoding/json"
	"errors"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/events"
	"github.com/TeaOSLab/EdgeAdmin/internal/goman"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func init() {
	events.On(events.EventStart, func() {
		var task = NewCheckUpdatesTask()
		goman.New(func() {
			task.Start()
		})
	})
}

type CheckUpdatesTask struct {
	ticker *time.Ticker
}

func NewCheckUpdatesTask() *CheckUpdatesTask {
	return &CheckUpdatesTask{}
}

func (this *CheckUpdatesTask) Start() {
	this.ticker = time.NewTicker(12 * time.Hour)
	for range this.ticker.C {
		err := this.Loop()
		if err != nil {
			logs.Println("[TASK][CHECK_UPDATES_TASK]" + err.Error())
		}
	}
}

func (this *CheckUpdatesTask) Loop() error {
	type Response struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	// 目前支持Linux
	if runtime.GOOS != "linux" {
		return nil
	}

	var apiURL = teaconst.UpdatesURL
	apiURL = strings.ReplaceAll(apiURL, "${os}", runtime.GOOS)
	apiURL = strings.ReplaceAll(apiURL, "${arch}", runtime.GOARCH)
	resp, err := http.Get(apiURL)
	if err != nil {
		return errors.New("read api failed: " + err.Error())
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("read api failed: " + err.Error())
	}

	var apiResponse = &Response{}
	err = json.Unmarshal(data, apiResponse)
	if err != nil {
		return errors.New("decode version data failed: " + err.Error())
	}

	if apiResponse.Code != 200 {
		return errors.New("invalid response: " + apiResponse.Message)
	}

	var m = maps.NewMap(apiResponse.Data)
	var dlHost = m.GetString("host")
	var versions = m.GetSlice("versions")
	if len(versions) > 0 {
		for _, version := range versions {
			var vMap = maps.NewMap(version)
			if vMap.GetString("code") == "admin" {
				var latestVersion = vMap.GetString("version")
				if stringutil.VersionCompare(teaconst.Version, latestVersion) < 0 {
					teaconst.NewVersionCode = latestVersion
					teaconst.NewVersionDownloadURL = dlHost + vMap.GetString("url")
					return nil
				}
			}
		}
	}

	return nil
}
