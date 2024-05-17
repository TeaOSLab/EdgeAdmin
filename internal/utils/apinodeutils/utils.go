// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package apinodeutils

import (
	"bytes"
	"errors"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/iwind/TeaGo/Tea"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

func CanUpgrade(apiVersion string, osName string, arch string) (canUpgrade bool, reason string) {
	if len(apiVersion) == 0 {
		return false, "current api version should not be empty"
	}

	if stringutil.VersionCompare(apiVersion, "0.6.4") < 0 {
		return false, "api node version must greater than or equal to 0.6.4"
	}

	if osName != runtime.GOOS {
		return false, "os not match: " + osName
	}
	if arch != runtime.GOARCH {
		return false, "arch not match: " + arch
	}

	stat, err := os.Stat(apiExe())
	if err != nil {
		return false, "stat error: " + err.Error()
	}
	if stat.IsDir() {
		return false, "is directory"
	}

	localVersion, err := lookupLocalVersion()
	if err != nil {
		return false, "lookup version failed: " + err.Error()
	}
	if localVersion != teaconst.APINodeVersion {
		return false, "not newest api node"
	}
	if stringutil.VersionCompare(localVersion, apiVersion) <= 0 {
		return false, "need not upgrade, local '" + localVersion + "' vs remote '" + apiVersion + "'"
	}

	return true, ""
}

func lookupLocalVersion() (string, error) {
	var cmd = exec.Command(apiExe(), "-V")
	var output = &bytes.Buffer{}
	cmd.Stdout = output
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	var localVersion = strings.TrimSpace(output.String())

	// 检查版本号
	var reg = regexp.MustCompile(`^[\d.]+$`)
	if !reg.MatchString(localVersion) {
		return "", errors.New("lookup version failed: " + localVersion)
	}

	return localVersion, nil
}

func apiExe() string {
	return Tea.Root + "/edge-api/bin/edge-api"
}
