// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/events"
	"github.com/TeaOSLab/EdgeAdmin/internal/goman"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"io"
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
	// 检查是否开启
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		return err
	}
	valueResp, err := rpcClient.SysSettingRPC().ReadSysSetting(rpcClient.Context(0), &pb.ReadSysSettingRequest{Code: systemconfigs.SettingCodeCheckUpdates})
	if err != nil {
		return err
	}
	var valueJSON = valueResp.ValueJSON
	var config = systemconfigs.NewCheckUpdatesConfig()
	if len(valueJSON) > 0 {
		err = json.Unmarshal(valueJSON, config)
		if err != nil {
			return fmt.Errorf("decode config failed: %w", err)
		}
		if !config.AutoCheck {
			return nil
		}
	} else {
		return nil
	}

	// 开始检查
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
	apiURL = strings.ReplaceAll(apiURL, "${version}", teaconst.Version)
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("read api failed: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read api failed: %w", err)
	}

	var apiResponse = &Response{}
	err = json.Unmarshal(data, apiResponse)
	if err != nil {
		return fmt.Errorf("decode version data failed: %w", err)
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
				if stringutil.VersionCompare(teaconst.Version, latestVersion) < 0 && (len(config.IgnoredVersion) == 0 || stringutil.VersionCompare(latestVersion, config.IgnoredVersion) > 0) {
					teaconst.NewVersionCode = latestVersion
					teaconst.NewVersionDownloadURL = dlHost + vMap.GetString("url")
					return nil
				}
			}
		}
	}

	return nil
}
