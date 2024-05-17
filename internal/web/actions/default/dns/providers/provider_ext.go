// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .
//go:build !plus

package providers

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

func (this *ProviderAction) readEdgeDNS(provider *pb.DNSProvider, apiParams maps.Map) (maps.Map, error) {
	return maps.Map{}, nil
}
