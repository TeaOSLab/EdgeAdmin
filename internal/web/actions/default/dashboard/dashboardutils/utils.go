// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package dashboardutils

import (
	"bytes"
	"context"
	"encoding/json"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"github.com/shirou/gopsutil/v3/disk"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

// CheckDiskPartitions 检查服务器硬盘空间
func CheckDiskPartitions(thresholdPercent float64) (path string, usage uint64, usagePercent float64, shouldWarning bool) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return
	}
	if !lists.ContainsString([]string{"darwin", "linux", "freebsd"}, runtime.GOOS) {
		return
	}

	var rootFS = ""

	for _, p := range partitions {
		if p.Mountpoint == "/" {
			rootFS = p.Fstype
			break
		}
	}

	for _, p := range partitions {
		if p.Mountpoint == "/boot" {
			continue
		}
		if p.Fstype != rootFS {
			continue
		}

		// skip some specified partitions on macOS
		if runtime.GOOS == "darwin" {
			if strings.Contains(p.Mountpoint, "/Developer/") {
				continue
			}
		}

		stat, _ := disk.Usage(p.Mountpoint)
		if stat != nil {
			if stat.Used < (5<<30) || stat.Free > (100<<30) {
				continue
			}
			if stat.UsedPercent > thresholdPercent {
				path = stat.Path
				usage = stat.Used
				usagePercent = stat.UsedPercent
				shouldWarning = true
				break
			}
		}
	}

	return
}

// CheckLocalAPINode 检查本地的API节点
func CheckLocalAPINode(rpcClient *rpc.RPCClient, ctx context.Context) (exePath string, runtimeVersion string, fileVersion string, ok bool) {
	resp, err := rpcClient.APINodeRPC().FindCurrentAPINode(ctx, &pb.FindCurrentAPINodeRequest{})
	if err != nil {
		return
	}
	if resp.ApiNode == nil {
		return
	}
	var instanceCode = resp.ApiNode.InstanceCode
	if len(instanceCode) == 0 {
		return
	}
	var statusJSON = resp.ApiNode.StatusJSON
	if len(statusJSON) == 0 {
		return
	}

	var status = &nodeconfigs.NodeStatus{}
	err = json.Unmarshal(statusJSON, status)
	if err != nil {
		return
	}
	runtimeVersion = status.BuildVersion

	if len(runtimeVersion) == 0 {
		return
	}

	if stringutil.VersionCompare(runtimeVersion, teaconst.APINodeVersion) >= 0 {
		return
	}

	exePath = status.ExePath
	if len(exePath) == 0 {
		return
	}

	stat, err := os.Stat(exePath)
	if err != nil {
		return
	}
	if stat.IsDir() {
		return
	}

	// 实例信息
	{
		var outputBuffer = &bytes.Buffer{}
		var cmd = exec.Command(exePath, "instance")
		cmd.Stdout = outputBuffer
		err = cmd.Run()
		if err != nil {
			return
		}

		var outputBytes = outputBuffer.Bytes()
		if len(outputBytes) == 0 {
			return
		}

		var instanceMap = maps.Map{}
		err = json.Unmarshal(bytes.TrimSpace(outputBytes), &instanceMap)
		if err != nil {
			return
		}

		if instanceMap.GetString("code") != instanceCode {
			return
		}
	}

	// 文件版本
	{
		var outputBuffer = &bytes.Buffer{}
		var cmd = exec.Command(exePath, "-v")
		cmd.Stdout = outputBuffer
		err = cmd.Run()
		if err != nil {
			return
		}

		var outputString = outputBuffer.String()
		if len(outputString) == 0 {
			return
		}

		var subMatch = regexp.MustCompile(`\s+v([\d.]+)\s+`).FindStringSubmatch(outputString)
		if len(subMatch) == 0 {
			return
		}
		fileVersion = subMatch[1]

		// 文件版本是否为最新
		if fileVersion != teaconst.APINodeVersion {
			fileVersion = runtimeVersion
		}
	}

	ok = true
	return
}
