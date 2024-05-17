// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .
//go:build !plus

package certs

func (this *UploadBatchPopupAction) maxFiles() int {
	return 20
}
