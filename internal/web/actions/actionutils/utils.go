package actionutils

import (
	"fmt"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

// 提示服务器错误信息
func Fail(action actions.ActionWrapper, err error) {
	if err != nil {
		logs.Println("[" + reflect.TypeOf(action).String() + "]" + findStack(err.Error()))
	}
	action.Object().Fail(teaconst.ErrServer)
}

// 提示页面错误信息
func FailPage(action actions.ActionWrapper, err error) {
	if err != nil {
		logs.Println("[" + reflect.TypeOf(action).String() + "]" + findStack(err.Error()))
	}
	action.Object().WriteString(teaconst.ErrServer)
}

// 判断动作的文件路径是否相当
func MatchPath(action *actions.ActionObject, path string) bool {
	return action.Request.URL.Path == path
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

			err += "\n\t\t" + string(filename) + ":" + fmt.Sprintf("%d", lineNo)

			break
		}
	}

	return err
}
