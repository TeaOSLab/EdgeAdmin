// Copyright 2023 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package nodes

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"strings"
)

type SessionManager struct {
	life      uint
	rpcClient *rpc.RPCClient
}

func NewSessionManager() (*SessionManager, error) {
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		return nil, err
	}
	return &SessionManager{
		rpcClient: rpcClient,
	}, nil
}

func (this *SessionManager) Init(config *actions.SessionConfig) {
	this.life = config.Life
}

func (this *SessionManager) Read(sid string) map[string]string {
	// 忽略OTP
	if strings.HasSuffix(sid, "_otp") {
		return map[string]string{}
	}

	var result = map[string]string{}

	resp, err := this.rpcClient.LoginSessionRPC().FindLoginSession(this.rpcClient.Context(0), &pb.FindLoginSessionRequest{Sid: sid})
	if err != nil {
		logs.Println("SESSION", "read '"+sid+"' failed: "+err.Error())
		return result
	}

	var session = resp.LoginSession
	if session == nil || len(session.ValuesJSON) == 0 {
		return result
	}

	err = json.Unmarshal(session.ValuesJSON, &result)
	if err != nil {
		logs.Println("SESSION", "decode '"+sid+"' values failed: "+err.Error())
	}

	return result
}

func (this *SessionManager) WriteItem(sid string, key string, value string) bool {
	// 忽略OTP
	if strings.HasSuffix(sid, "_otp") {
		return false
	}

	_, err := this.rpcClient.LoginSessionRPC().WriteLoginSessionValue(this.rpcClient.Context(0), &pb.WriteLoginSessionValueRequest{
		Sid:   sid,
		Key:   key,
		Value: value,
	})
	if err != nil {
		logs.Println("SESSION", "write sid:'"+sid+"' key:'"+key+"' failed: "+err.Error())
	}

	return true
}

func (this *SessionManager) Delete(sid string) bool {
	// 忽略OTP
	if strings.HasSuffix(sid, "_otp") {
		return false
	}

	_, err := this.rpcClient.LoginSessionRPC().DeleteLoginSession(this.rpcClient.Context(0), &pb.DeleteLoginSessionRequest{Sid: sid})
	if err != nil {
		logs.Println("SESSION", "delete '"+sid+"' failed: "+err.Error())
	}
	return true
}
