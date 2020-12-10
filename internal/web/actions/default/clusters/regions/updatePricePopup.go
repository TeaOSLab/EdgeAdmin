package regions

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdatePricePopupAction struct {
	actionutils.ParentAction
}

func (this *UpdatePricePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdatePricePopupAction) RunGet(params struct {
	RegionId int64
	ItemId   int64
}) {
	// 区域
	regionResp, err := this.RPC().NodeRegionRPC().FindEnabledNodeRegion(this.AdminContext(), &pb.FindEnabledNodeRegionRequest{NodeRegionId: params.RegionId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	region := regionResp.NodeRegion
	if region == nil {
		this.NotFound("nodeRegion", params.RegionId)
		return
	}
	this.Data["region"] = maps.Map{
		"id":   region.Id,
		"isOn": region.IsOn,
		"name": region.Name,
	}

	// 当前价格
	pricesMap := map[string]float32{}
	if len(region.PricesJSON) > 0 {
		err = json.Unmarshal(region.PricesJSON, &pricesMap)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	this.Data["price"] = pricesMap[numberutils.FormatInt64(params.ItemId)]

	// 价格项
	itemResp, err := this.RPC().NodePriceItemRPC().FindEnabledNodePriceItem(this.AdminContext(), &pb.FindEnabledNodePriceItemRequest{NodePriceItemId: params.ItemId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	item := itemResp.NodePriceItem
	if item == nil {
		this.NotFound("nodePriceItem", params.ItemId)
		return
	}

	this.Data["item"] = maps.Map{
		"id":       item.Id,
		"name":     item.Name,
		"bitsFrom": item.BitsFrom,
		"bitsTo":   item.BitsTo,
	}

	this.Show()
}

func (this *UpdatePricePopupAction) RunPost(params struct {
	RegionId int64
	ItemId   int64
	Price    float32

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改区域 %d-价格项 %d 的价格", params.RegionId, params.ItemId)

	_, err := this.RPC().NodeRegionRPC().UpdateNodeRegionPrice(this.AdminContext(), &pb.UpdateNodeRegionPriceRequest{
		NodeRegionId: params.RegionId,
		NodeItemId:   params.ItemId,
		Price:        params.Price,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Success()
}
