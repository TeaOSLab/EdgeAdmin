package ipadmin

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type SelectProvidersPopupAction struct {
	actionutils.ParentAction
}

func (this *SelectProvidersPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *SelectProvidersPopupAction) RunGet(params struct {
	Type                string
	SelectedProviderIds string
}) {
	this.Data["type"] = params.Type

	var selectedProviderIds = utils.SplitNumbers(params.SelectedProviderIds)

	providersResp, err := this.RPC().RegionProviderRPC().FindAllRegionProviders(this.AdminContext(), &pb.FindAllRegionProvidersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var providerMaps = []maps.Map{}
	for _, provider := range providersResp.RegionProviders {
		providerMaps = append(providerMaps, maps.Map{
			"id":        provider.Id,
			"name":      provider.Name,
			"isChecked": lists.ContainsInt64(selectedProviderIds, provider.Id),
		})
	}
	this.Data["providers"] = providerMaps

	this.Show()
}
