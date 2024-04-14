package actionutils

import (
	"context"
	"errors"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/index/loginutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"net/http"
	"strconv"
)

type ParentAction struct {
	actions.ActionObject

	rpcClient *rpc.RPCClient

	ctx context.Context
}

// Parent 可以调用自身的一个简便方法
func (this *ParentAction) Parent() *ParentAction {
	return this
}

func (this *ParentAction) ErrorPage(err error) {
	if err == nil {
		return
	}

	// 日志
	this.CreateLog(oplogs.LevelError, codes.AdminCommon_LogSystemError, err.Error())

	if this.Request.Method == http.MethodGet {
		FailPage(this, err)
	} else {
		Fail(this, err)
	}
}

func (this *ParentAction) ErrorText(err string) {
	this.ErrorPage(errors.New(err))
}

func (this *ParentAction) NotFound(name string, itemId int64) {
	if itemId > 0 {
		this.ErrorPage(errors.New(name + " id: '" + strconv.FormatInt(itemId, 10) + "' is not found"))
	} else {
		this.ErrorPage(errors.New(name + " is not found"))
	}
}

func (this *ParentAction) NewPage(total int64, size ...int64) *Page {
	if len(size) > 0 {
		return NewActionPage(this, total, size[0])
	}

	var pageSize int64 = 10
	adminConfig, err := configloaders.LoadAdminUIConfig()
	if err == nil && adminConfig.DefaultPageSize > 0 {
		pageSize = int64(adminConfig.DefaultPageSize)
	}

	return NewActionPage(this, total, pageSize)
}

func (this *ParentAction) Nav(mainMenu string, tab string, firstMenu string) {
	this.Data["mainMenu"] = mainMenu
	this.Data["mainTab"] = tab
	this.Data["firstMenuItem"] = firstMenu
}

func (this *ParentAction) FirstMenu(menuItem string) {
	this.Data["firstMenuItem"] = menuItem
}

func (this *ParentAction) SecondMenu(menuItem string) {
	this.Data["secondMenuItem"] = menuItem
}

func (this *ParentAction) TinyMenu(menuItem string) {
	this.Data["tinyMenuItem"] = menuItem
}

func (this *ParentAction) AdminId() int64 {
	return this.Context.GetInt64(teaconst.SessionAdminId)
}

func (this *ParentAction) CreateLog(level string, messageCode langs.MessageCode, args ...any) {
	var description = messageCode.For(this.LangCode())
	var desc = fmt.Sprintf(description, args...)
	if level == oplogs.LevelInfo {
		if this.Code != 200 {
			level = oplogs.LevelWarn
			if len(this.Message) > 0 {
				desc += " 失败：" + this.Message
			}
		}
	}
	err := dao.SharedLogDAO.CreateAdminLog(this.AdminContext(), level, this.Request.URL.Path, desc, loginutils.RemoteIP(&this.ActionObject), messageCode, args)
	if err != nil {
		utils.PrintError(err)
	}
}

func (this *ParentAction) CreateLogInfo(messageCode langs.MessageCode, args ...any) {
	this.CreateLog(oplogs.LevelInfo, messageCode, args...)
}

// RPC 获取RPC
func (this *ParentAction) RPC() *rpc.RPCClient {
	if this.rpcClient != nil {
		return this.rpcClient
	}

	// 所有集群
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		logs.Fatal(err)
		return nil
	}
	this.rpcClient = rpcClient

	return rpcClient
}

// AdminContext 获取Context
// 每个请求的context都必须是一个新的实例
func (this *ParentAction) AdminContext() context.Context {
	if this.rpcClient == nil {
		rpcClient, err := rpc.SharedRPC()
		if err != nil {
			logs.Fatal(err)
			return nil
		}
		this.rpcClient = rpcClient
	}
	this.ctx = this.rpcClient.Context(this.AdminId())
	return this.ctx
}

// ViewData 视图里可以使用的数据
func (this *ParentAction) ViewData() maps.Map {
	return this.Data
}

func (this *ParentAction) LangCode() string {
	return configloaders.FindAdminLangForAction(this)
}

func (this *ParentAction) Lang(messageCode langs.MessageCode, args ...any) string {
	return langs.Message(this.LangCode(), messageCode, args...)
}

func (this *ParentAction) FailLang(messageCode langs.MessageCode, args ...any) {
	this.Fail(langs.Message(this.LangCode(), messageCode, args...))
}

func (this *ParentAction) FailFieldLang(field string, messageCode langs.MessageCode, args ...any) {
	this.FailField(field, langs.Message(this.LangCode(), messageCode, args...))
}

func (this *ParentAction) FilterHTTPFamily() bool {
	if this.Data.GetString("serverFamily") == "http" {
		return false
	}

	this.ResponseWriter.WriteHeader(http.StatusNotFound)
	_, _ = this.ResponseWriter.Write([]byte("page not found"))

	return true
}
