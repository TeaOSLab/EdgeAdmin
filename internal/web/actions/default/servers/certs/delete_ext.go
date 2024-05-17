// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .
//go:build !plus

package certs

func (this *DeleteAction) filterDelete(certId int64) error {
	return nil
}
