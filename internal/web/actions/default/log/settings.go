package log

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
)

type SettingsAction struct {
	actionutils.ParentAction
}

func (this *SettingsAction) Init() {
	this.Nav("", "", "setting")
}

func (this *SettingsAction) RunGet(params struct{}) {
	config, err := configloaders.LoadLogConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["logConfig"] = config

	this.Show()
}

func (this *SettingsAction) RunPost(params struct {
	CanDelete    bool
	CanClean     bool
	CapacityJSON []byte
	Days         int

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	capacity := &shared.SizeCapacity{}
	err := json.Unmarshal(params.CapacityJSON, capacity)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	config, err := configloaders.LoadLogConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	config.CanDelete = params.CanDelete
	config.CanClean = params.CanClean
	config.Days = params.Days
	config.Capacity = capacity
	err = configloaders.UpdateLogConfig(config)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
