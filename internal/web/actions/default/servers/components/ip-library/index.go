package iplibrary

import (
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.FirstMenu("index")
}

func (this *IndexAction) RunGet(params struct {
	Type string
}) {
	if len(params.Type) == 0 {
		params.Type = serverconfigs.IPLibraryTypes[0].GetString("code")
	}

	this.Data["types"] = serverconfigs.IPLibraryTypes
	this.Data["selectedType"] = params.Type

	// 列表
	listResp, err := this.RPC().IPLibraryRPC().FindAllEnabledIPLibrariesWithType(this.AdminContext(), &pb.FindAllEnabledIPLibrariesWithTypeRequest{Type: params.Type})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	libraryMaps := []maps.Map{}
	for _, library := range listResp.IpLibraries {
		var fileMap maps.Map = nil
		if library.File != nil {
			fileMap = maps.Map{
				"id":       library.File.Id,
				"filename": library.File.Filename,
				"sizeMB":   fmt.Sprintf("%.2f", float64(library.File.Size)/1024/1024),
			}
		}

		libraryMaps = append(libraryMaps, maps.Map{
			"id":          library.Id,
			"file":        fileMap,
			"createdTime": timeutil.FormatTime("Y-m-d H:i:s", library.CreatedAt),
		})
	}
	this.Data["libraries"] = libraryMaps

	this.Show()
}
