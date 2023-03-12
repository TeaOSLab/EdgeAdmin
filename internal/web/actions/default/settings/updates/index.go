// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package updates

import (
	"encoding/json"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
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

func (this *IndexAction) RunGet(params struct{}) {
	this.Data["version"] = teaconst.Version

	valueResp, err := this.RPC().SysSettingRPC().ReadSysSetting(this.AdminContext(), &pb.ReadSysSettingRequest{Code: systemconfigs.SettingCodeCheckUpdates})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var valueJSON = valueResp.ValueJSON
	var config = &systemconfigs.CheckUpdatesConfig{AutoCheck: false}
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
	}

	var apiResponse = &Response{}
	err = json.Unmarshal(data, apiResponse)
	if err != nil {
		this.Data["result"] = maps.Map{
			"isOk":    false,
			"message": "解析更新信息失败：" + err.Error(),
		}
		this.Success()
	}

	if apiResponse.Code != 200 {
		this.Data["result"] = maps.Map{
			"isOk":    false,
			"message": "解析更新信息失败：" + apiResponse.Message,
		}
		this.Success()
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
					this.Data["result"] = maps.Map{
						"isOk":    true,
						"message": "有最新的版本v" + types.String(latestVersion) + "可以更新",
						"hasNew":  true,
						"dlURL":   dlHost + vMap.GetString("url"),
					}
					this.Success()
				} else {
					this.Data["result"] = maps.Map{
						"isOk":    true,
						"message": "你已安装最新版本，无需更新",
					}
					this.Success()
				}
			}
		}
	}

	this.Data["result"] = maps.Map{
		"isOk":    false,
		"message": "找不到更新信息",
	}

	this.Success()
}
