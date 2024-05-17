// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package nodes

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/ttlcache"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"strings"
	"time"
)

// SessionManager SESSION管理
type SessionManager struct {
	life uint
}

func NewSessionManager() (*SessionManager, error) {
	return &SessionManager{}, nil
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

	var cacheKey = "SESSION@" + sid
	var item = ttlcache.DefaultCache.Read(cacheKey)
	if item != nil && item.Value != nil {
		itemMap, ok := item.Value.(map[string]string)
		if ok {
			return itemMap
		}
	}

	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		return map[string]string{}
	}

	resp, err := rpcClient.LoginSessionRPC().FindLoginSession(rpcClient.Context(0), &pb.FindLoginSessionRequest{Sid: sid})
	if err != nil {
		logs.Println("SESSION", "read '"+sid+"' failed: "+err.Error())
		result["@error"] = err.Error()
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

	// Write to cache
	ttlcache.DefaultCache.Write(cacheKey, result, time.Now().Unix()+300 /** must not be too long **/)

	return result
}

func (this *SessionManager) WriteItem(sid string, key string, value string) bool {
	// 删除缓存
	defer ttlcache.DefaultCache.Delete("SESSION@" + sid)

	// 忽略OTP
	if strings.HasSuffix(sid, "_otp") {
		return false
	}

	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		return false
	}

	_, err = rpcClient.LoginSessionRPC().WriteLoginSessionValue(rpcClient.Context(0), &pb.WriteLoginSessionValueRequest{
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
	// 删除缓存
	defer ttlcache.DefaultCache.Delete("SESSION@" + sid)

	// 忽略OTP
	if strings.HasSuffix(sid, "_otp") {
		return false
	}

	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		return false
	}
	_, err = rpcClient.LoginSessionRPC().DeleteLoginSession(rpcClient.Context(0), &pb.DeleteLoginSessionRequest{Sid: sid})
	if err != nil {
		logs.Println("SESSION", "delete '"+sid+"' failed: "+err.Error())
	}
	return true
}
