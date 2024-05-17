// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package updates

import (
	"encoding/json"
	"fmt"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/iwind/TeaGo/maps"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"io"
	"net/http"
	"runtime"
	"strings"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "updates", "")
}

func (this *IndexAction) RunGet(params struct {
	DoCheck bool
}) {
	this.Data["version"] = teaconst.Version
	this.Data["doCheck"] = params.DoCheck

	// 是否正在升级
	this.Data["isUpgrading"] = isUpgrading
	this.Data["isUpgradingDB"] = isUpgradingDB
	this.Data["upgradeProgress"] = fmt.Sprintf("%.2f", upgradeProgress * 100)
	if isUpgrading {
		this.Data["doCheck"] = false
	}

	valueResp, err := this.RPC().SysSettingRPC().ReadSysSetting(this.AdminContext(), &pb.ReadSysSettingRequest{Code: systemconfigs.SettingCodeCheckUpdates})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var valueJSON = valueResp.ValueJSON
	var config = systemconfigs.NewCheckUpdatesConfig()
	if len(valueJSON) > 0 {
		err = json.Unmarshal(valueJSON, config)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	this.Data["config"] = config

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
}) {
	valueResp, err := this.RPC().SysSettingRPC().ReadSysSetting(this.AdminContext(), &pb.ReadSysSettingRequest{Code: systemconfigs.SettingCodeCheckUpdates})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var valueJSON = valueResp.ValueJSON
	var config = systemconfigs.NewCheckUpdatesConfig()
	if len(valueJSON) > 0 {
		err = json.Unmarshal(valueJSON, config)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	type Response struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	var apiURL = teaconst.UpdatesURL
	apiURL = strings.ReplaceAll(apiURL, "${os}", runtime.GOOS)
	apiURL = strings.ReplaceAll(apiURL, "${arch}", runtime.GOARCH)
	apiURL = strings.ReplaceAll(apiURL, "${version}", teaconst.Version)
	resp, err := http.Get(apiURL)
	if err != nil {
		this.Data["result"] = maps.Map{
			"isOk":    false,
			"message": "读取更新信息失败：" + err.Error(),
		}
		this.Success()
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		this.Data["result"] = maps.Map{
			"isOk":    false,
			"message": "读取更新信息失败：" + err.Error(),
		}
		this.Success()
		return
	}

	var apiResponse = &Response{}
	err = json.Unmarshal(data, apiResponse)
	if err != nil {
		this.Data["result"] = maps.Map{
			"isOk":    false,
			"message": "解析更新信息失败：" + err.Error(),
		}
		this.Success()
		return
	}

	if apiResponse.Code != 200 {
		this.Data["result"] = maps.Map{
			"isOk":    false,
			"message": "解析更新信息失败：" + apiResponse.Message,
		}
		this.Success()
		return
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
					// 是否已忽略
					if len(config.IgnoredVersion) > 0 && stringutil.VersionCompare(config.IgnoredVersion, latestVersion) >= 0 {
						continue
					}

					this.Data["result"] = maps.Map{
						"isOk":        true,
						"version":     latestVersion,
						"message":     "有最新的版本 v" + latestVersion + " 可以更新",
						"hasNew":      true,
						"dlURL":       dlHost + vMap.GetString("url"),
						"day":         vMap.GetString("day"),
						"description": vMap.GetString("description"),
						"docURL":      vMap.GetString("docURL"),
					}
					this.Success()
					return
				} else {
					this.Data["result"] = maps.Map{
						"isOk":    true,
						"message": "你已安装最新版本，无需更新",
					}
					this.Success()
					return
				}
			}
		}
	}

	this.Data["result"] = maps.Map{
		"isOk":    false,
		"message": "没有发现更新的版本",
	}

	this.Success()
}
