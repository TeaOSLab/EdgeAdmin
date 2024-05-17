// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package clusterutils

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/iwind/TeaGo/Tea"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type installerFile struct {
	Name    string `json:"name"`
	OS      string `json:"os"`
	Arch    string `json:"arch"`
	Version string `json:"version"`
}

func ListInstallerFiles() []*installerFile {
	var dir = Tea.Root + "/edge-api/deploy"
	matches, err := filepath.Glob(dir + "/edge-node-*.zip")
	if err != nil {
		return nil
	}

	var result = []*installerFile{}
	var reg = regexp.MustCompile(`^edge-node-(\w+)-(\w+)-v([\w.]+)\.zip$`)
	for _, match := range matches {
		var baseName = filepath.Base(match)
		var subMatches = reg.FindStringSubmatch(baseName)
		if len(subMatches) >= 4 {
			var osName = subMatches[1]
			if len(osName) > 0 {
				osName = strings.ToUpper(osName[:1]) + osName[1:]
			}

			var arch = subMatches[2]
			if arch == "amd64" {
				arch = "x86_64"
			}

			var version = subMatches[3]
			if version != teaconst.Version { // 只能下载当前版本
				continue
			}

			result = append(result, &installerFile{
				Name:    subMatches[0],
				OS:      osName,
				Arch:    arch,
				Version: version,
			})
		}
	}

	// 排序，将x86_64排在最上面
	if len(result) > 0 {
		sort.Slice(result, func(i, j int) bool {
			return result[i].Arch == "x86_64"
		})
	}

	return result
}
