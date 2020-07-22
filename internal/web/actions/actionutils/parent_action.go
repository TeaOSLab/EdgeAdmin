package actionutils

import (
	"errors"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc/admin"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/iwind/TeaGo/actions"
	"net/http"
	"strconv"
)

type ParentAction struct {
	actions.ActionObject
}

func (this *ParentAction) ErrorPage(err error) {
	if err == nil {
		return
	}

	// 日志
	this.CreateLog(oplogs.LevelError, "系统发生错误：%s", err.Error())

	if this.Request.Method == http.MethodGet {
		FailPage(this, err)
	} else {
		Fail(this, err)
	}
}

func (this *ParentAction) ErrorText(err string) {
	this.ErrorPage(errors.New(err))
}

func (this *ParentAction) NotFound(name string, itemId int) {
	this.ErrorPage(errors.New(name + " id: '" + strconv.Itoa(itemId) + "' is not found"))
}

func (this *ParentAction) NewPage(total int, size ...int) *Page {
	if len(size) > 0 {
		return NewActionPage(this, total, size[0])
	}
	return NewActionPage(this, total, 10)
}

func (this *ParentAction) Nav(mainMenu string, tab string, firstMenu string) {
	this.Data["mainMenu"] = mainMenu
	this.Data["mainTab"] = tab
	this.Data["firstMenuItem"] = firstMenu
}

func (this *ParentAction) SecondMenu(menuItem string) {
	this.Data["secondMenuItem"] = menuItem
}

func (this *ParentAction) AdminId() int {
	return this.Context.GetInt("adminId")
}

func (this *ParentAction) CreateLog(level string, description string, args ...interface{}) {
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		utils.PrintError(err)
		return
	}
	_, err = rpcClient.AdminRPC().CreateLog(rpcClient.Context(this.AdminId()), &admin.CreateLogRequest{
		Level:       level,
		Description: fmt.Sprintf(description, args...),
		Action:      this.Request.URL.Path,
		Ip:          this.RequestRemoteIP(),
	})
	if err != nil {
		utils.PrintError(err)
	}
}
