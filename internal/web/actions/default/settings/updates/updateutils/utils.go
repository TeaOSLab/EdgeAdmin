// Copyright 2024 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package updateutils

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"os"
	"os/exec"
)

func CheckLocalAPINode(rpcClient *rpc.RPCClient, ctx context.Context) (exePath string, ok bool) {
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

	ok = true
	return
}
