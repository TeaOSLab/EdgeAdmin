// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package tasks

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/events"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/setup"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"time"
)

func init() {
	events.On(events.EventStart, func() {
		task := NewAuthorityTask()
		go task.Start()
	})
}

type AuthorityTask struct {
}

func NewAuthorityTask() *AuthorityTask {
	return &AuthorityTask{}
}

func (this *AuthorityTask) Start() {
	// 从缓存中读取
	config := configs.ReadPlusConfig()
	if config != nil {
		teaconst.IsPlus = config.IsPlus
	}

	// 开始计时器
	ticker := time.NewTicker(10 * time.Minute)
	if Tea.IsTesting() {
		// 快速测试
		ticker = time.NewTicker(1 * time.Minute)
	}

	// 初始化的时候先获取一次
	timeout := time.NewTimer(3 * time.Second)
	<-timeout.C
	err := this.Loop()
	if err != nil {
		logs.Println("[TASK][AuthorityTask]" + err.Error())
	}

	// 定时获取
	for range ticker.C {
		err := this.Loop()
		if err != nil {
			logs.Println("[TASK][AuthorityTask]" + err.Error())
		}
	}
}

func (this *AuthorityTask) Loop() error {
	// 如果还没有安装直接返回
	if !setup.IsConfigured() {
		return nil
	}

	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		return err
	}
	resp, err := rpcClient.AuthorityKeyRPC().ReadAuthorityKey(rpcClient.Context(0), &pb.ReadAuthorityKeyRequest{})
	if err != nil {
		return err
	}
	var oldState = teaconst.IsPlus
	if resp.AuthorityKey != nil {
		teaconst.IsPlus = true
	} else {
		teaconst.IsPlus = false
	}

	if oldState != teaconst.IsPlus {
		_ = configs.WritePlusConfig(&configs.PlusConfig{IsPlus: teaconst.IsPlus})
	}

	return nil
}
