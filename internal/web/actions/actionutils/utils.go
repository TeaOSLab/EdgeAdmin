package actionutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	rpcerrors "github.com/TeaOSLab/EdgeCommon/pkg/rpc/errors"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/gosock/pkg/gosock"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"
)

// Fail 提示服务器错误信息
func Fail(actionPtr actions.ActionWrapper, err error) {
	if err == nil {
		err = errors.New("unknown error")
	}

	var langCode = configloaders.FindAdminLangForAction(actionPtr)
	var serverErrString = codes.AdminCommon_ServerError.For(langCode)

	logs.Println("[" + reflect.TypeOf(actionPtr).String() + "]" + findStack(err.Error()))

	_, _, isLocalAPI, issuesHTML := parseAPIErr(actionPtr, err)
	if isLocalAPI && len(issuesHTML) > 0 {
		actionPtr.Object().Fail(serverErrString + "（" + err.Error() + "；最近一次错误提示：" + issuesHTML + "）")
	} else {
		actionPtr.Object().Fail(serverErrString + "（" + err.Error() + "）")
	}

}

// FailPage 提示页面错误信息
func FailPage(actionPtr actions.ActionWrapper, err error) {
	if err == nil {
		err = errors.New("unknown error")
	}

	var langCode = configloaders.FindAdminLangForAction(actionPtr)
	var serverErrString = codes.AdminCommon_ServerError.For(langCode)

	logs.Println("[" + reflect.TypeOf(actionPtr).String() + "]" + findStack(err.Error()))

	actionPtr.Object().ResponseWriter.WriteHeader(http.StatusInternalServerError)

	if len(actionPtr.Object().Request.Header.Get("X-Requested-With")) > 0 {
		actionPtr.Object().WriteString(serverErrString)
	} else {
		apiNodeIsStarting, apiNodeProgress, _, issuesHTML := parseAPIErr(actionPtr, err)
		var html = `<!DOCTYPE html>
<html>
	<head>
		<title>正在处理...</title>
		<meta charset="UTF-8"/>
		<style type="text/css">
		hr { border-top: 1px #ccc solid; }
		.red { color: red; }
		</style>
	</head>
	<body>
	<div style="background: #eee; border: 1px #ccc solid; padding: 10px; font-size: 12px; line-height: 1.8">
	`
		if apiNodeIsStarting { // API节点正在启动
			html += "<div class=\"red\">API节点正在启动，请耐心等待完成"

			if len(apiNodeProgress) > 0 {
				html += "：" + apiNodeProgress + "（刷新当前页面查看最新状态）"
			}

			html += "</div>"
		} else {
			html += serverErrString + `
		<div>可以通过查看 <strong><em>$安装目录/logs/run.log</em></strong> 日志文件查看具体的错误提示。</div>
		<hr/>
		<div class="red">Error: ` + err.Error() + `</div>`

			if len(issuesHTML) > 0 {
				html += `        <hr/>
        <div class="red">` + issuesHTML + `</div>`
			}
		}

		actionPtr.Object().WriteString(html + `
		</div>
	</body>
</html>`)
	}
}

// MatchPath 判断动作的文件路径是否相当
func MatchPath(action *actions.ActionObject, path string) bool {
	return action.Request.URL.Path == path
}

// FindParentAction 查找父级Action
func FindParentAction(actionPtr actions.ActionWrapper) *ParentAction {
	action, ok := actionPtr.(interface {
		Parent() *ParentAction
	})
	if ok {
		return action.Parent()
	}

	return nil
}

func findStack(err string) string {
	_, currentFilename, _, currentOk := runtime.Caller(1)
	if currentOk {
		for i := 1; i < 32; i++ {
			_, filename, lineNo, ok := runtime.Caller(i)
			if !ok {
				break
			}

			if filename == currentFilename || filepath.Base(filename) == "parent_action.go" {
				continue
			}

			goPath := os.Getenv("GOPATH")
			if len(goPath) > 0 {
				absGoPath, err := filepath.Abs(goPath)
				if err == nil {
					filename = strings.TrimPrefix(filename, absGoPath)[1:]
				}
			} else if strings.Contains(filename, "src") {
				filename = filename[strings.Index(filename, "src"):]
			}

			err += "\n\t\t" + filename + ":" + fmt.Sprintf("%d", lineNo)

			break
		}
	}

	return err
}

// 分析API节点的错误信息
func parseAPIErr(action actions.ActionWrapper, err error) (apiNodeIsStarting bool, apiNodeProgress string, isLocalAPI bool, issuesHTML string) {
	// 当前API终端地址
	var apiEndpoints = []string{}
	apiConfig, apiConfigErr := configs.LoadAPIConfig()
	if apiConfigErr == nil && apiConfig != nil {
		apiEndpoints = append(apiEndpoints, apiConfig.RPCEndpoints...)
	}

	var isRPCConnError bool
	_, isRPCConnError = rpcerrors.HumanError(err, apiEndpoints, Tea.ConfigFile(configs.ConfigFileName))
	if isRPCConnError {
		// API节点是否正在启动
		var sock = gosock.NewTmpSock("edge-api")
		reply, err := sock.SendTimeout(&gosock.Command{
			Code:   "starting",
			Params: nil,
		}, 1*time.Second)
		if err == nil && reply != nil {
			var params = maps.NewMap(reply.Params)
			if params.GetBool("isStarting") {
				apiNodeIsStarting = true

				var progressMap = params.GetMap("progress")
				apiNodeProgress = progressMap.GetString("description")
			}
		}
	}

	// 本地的一些错误提示
	if isRPCConnError {
		host, _, hostErr := net.SplitHostPort(action.Object().Request.Host)
		if hostErr == nil {
			for _, endpoint := range apiEndpoints {
				if strings.HasPrefix(endpoint, "http://"+host) || strings.HasPrefix(endpoint, "https://"+host) || strings.HasPrefix(endpoint, host) {
					isLocalAPI = true
					break
				}
			}
		}
	}

	if isLocalAPI {
		// 读取本地API节点的issues
		issuesData, issuesErr := os.ReadFile(Tea.Root + "/edge-api/logs/issues.log")
		if issuesErr == nil {
			var issueMaps = []maps.Map{}
			issuesErr = json.Unmarshal(issuesData, &issueMaps)
			if issuesErr == nil && len(issueMaps) > 0 {
				var issueMap = issueMaps[0]
				issuesHTML = "本地API节点启动错误：" + issueMap.GetString("message") + "，处理建议：" + issueMap.GetString("suggestion")
			}
		}
	}

	return
}
