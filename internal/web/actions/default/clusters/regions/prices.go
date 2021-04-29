package regions

import (
	"encoding/json"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/regions/regionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type PricesAction struct {
	actionutils.ParentAction
}

func (this *PricesAction) Init() {
	this.Nav("", "", "price")
}

func (this *PricesAction) RunGet(params struct{}) {
	// 所有价格项目
	itemsResp, err := this.RPC().NodePriceItemRPC().FindAllEnabledAndOnNodePriceItems(this.AdminContext(), &pb.FindAllEnabledAndOnNodePriceItemsRequest{Type: regionutils.PriceTypeTraffic})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	itemMaps := []maps.Map{}
	for _, item := range itemsResp.NodePriceItems {

		itemMaps = append(itemMaps, maps.Map{
			"id":             item.Id,
			"name":           item.Name,
			"bitsFromString": this.formatBits(item.BitsFrom),
			"bitsToString":   this.formatBits(item.BitsTo),
		})
	}
	this.Data["items"] = itemMaps

	// 所有区域
	regionsResp, err := this.RPC().NodeRegionRPC().FindAllEnabledNodeRegions(this.AdminContext(), &pb.FindAllEnabledNodeRegionsRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	regionMaps := []maps.Map{}
	for _, region := range regionsResp.NodeRegions {
		pricesMap := map[string]float32{}
		if len(region.PricesJSON) > 0 {
			err = json.Unmarshal(region.PricesJSON, &pricesMap)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}
		regionMaps = append(regionMaps, maps.Map{
			"id":     region.Id,
			"isOn":   region.IsOn,
			"name":   region.Name,
			"prices": pricesMap,
		})
	}
	this.Data["regions"] = regionMaps

	this.Show()
}

func (this *PricesAction) formatBits(bits int64) string {
	sizeHuman := ""
	if bits < 1000 {
		sizeHuman = numberutils.FormatInt64(bits) + "BPS"
	} else if bits < 1_000_000 {
		sizeHuman = fmt.Sprintf("%.2fKBPS", float64(bits)/1000)
	} else if bits < 1_000_000_000 {
		sizeHuman = fmt.Sprintf("%.2fMBPS", float64(bits)/1000/1000)
	} else if bits < 1_000_000_000_000 {
		sizeHuman = fmt.Sprintf("%.2fGBPS", float64(bits)/1000/1000/1000)
	} else if bits < 1_000_000_000_000_000 {
		sizeHuman = fmt.Sprintf("%.2fTBPS", float64(bits)/1000/1000/1000/1000)
	} else if bits < 1_000_000_000_000_000_000 {
		sizeHuman = fmt.Sprintf("%.2fPBPS", float64(bits)/1000/1000/1000/1000/1000)
	} else {
		sizeHuman = fmt.Sprintf("%.2fEBPS", float64(bits)/1000/1000/1000/1000/1000/1000)
	}
	return sizeHuman
}
