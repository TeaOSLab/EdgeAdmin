package servers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
	"sync"
)

type StatusAction struct {
	actionutils.ParentAction
}

func (this *StatusAction) RunPost(params struct {
	ServerIds []int64
}) {
	status := map[int64]maps.Map{}
	statusLocker := sync.Mutex{}

	if len(params.ServerIds) == 0 {
		this.Data["status"] = status
		this.Success()
	}

	// 读取全局配置
	var wg = sync.WaitGroup{}
	wg.Add(len(params.ServerIds))

	for _, serverId := range params.ServerIds {
		go func(serverId int64) {
			defer utils.Recover()
			defer wg.Done()

			m := maps.Map{
				"isOk":    false,
				"message": "",
				"todo":    "",
				"type":    "",
			}

			defer func() {
				statusLocker.Lock()
				defer statusLocker.Unlock()

				status[serverId] = m
			}()

			// 检查cname
			serverDNSResp, err := this.RPC().ServerRPC().FindEnabledServerDNS(this.AdminContext(), &pb.FindEnabledServerDNSRequest{ServerId: serverId})
			if err != nil {
				this.ErrorPage(err)

				m["type"] = "serverErr"
				m["message"] = "服务器错误"
				m["todo"] = "错误信息：FindEnabledServerDNS(): " + err.Error() + "，请联系管理员修复此问题"
				return
			}

			if len(serverDNSResp.DnsName) == 0 {
				m["type"] = "dnsNameEmpty"
				m["message"] = "CNAME为空"
				m["todo"] = "请删除后重新创建服务"
				return
			}

			if serverDNSResp.Domain == nil {
				m["type"] = "clusterDNSEmpty"
				m["message"] = "集群配置错误"
				m["todo"] = "所属集群没有配置DNS，请联系管理员修复此问题。服务ID：" + numberutils.FormatInt64(serverId)
				return
			}

			// 检查DNS是否已经设置
			serverNamesResp, err := this.RPC().ServerRPC().FindServerNames(this.AdminContext(), &pb.FindServerNamesRequest{ServerId: serverId})
			if err != nil {
				this.ErrorPage(err)

				m["type"] = "serverErr"
				m["message"] = "服务器错误"
				m["todo"] = "错误信息：FindServerNames(): " + err.Error() + "，请联系管理员修复此问题"
				return
			}
			if serverNamesResp.IsAuditing {
				m["type"] = "auditing"
				m["message"] = "审核中"

				auditingPromptResp, err := this.RPC().ServerRPC().FindServerAuditingPrompt(this.AdminContext(), &pb.FindServerAuditingPromptRequest{ServerId: serverId})
				if err != nil {
					this.ErrorPage(err)
					m["type"] = "serverErr"
					m["message"] = "服务器错误"
					m["todo"] = "错误信息：FindServerNames(): " + err.Error() + "，请联系管理员修复此问题"
					return
				}

				var auditingPrompt = auditingPromptResp.PromptText
				if len(auditingPrompt) > 0 {
					m["todo"] = auditingPrompt
				} else {
					m["todo"] = "域名正在审核中，请耐心等待"
				}
				return
			}
			if serverNamesResp.AuditingResult != nil && !serverNamesResp.AuditingResult.IsOk {
				m["type"] = "auditingFailed"
				m["message"] = "审核不通过"
				m["todo"] = "审核不通过，原因：" + serverNamesResp.AuditingResult.Reason
				return
			}

			serverNames := []*serverconfigs.ServerNameConfig{}
			if len(serverNamesResp.ServerNamesJSON) > 0 {
				err = json.Unmarshal(serverNamesResp.ServerNamesJSON, &serverNames)
				if err != nil {
					this.ErrorPage(err)

					m["type"] = "serverErr"
					m["message"] = "服务器错误"
					m["todo"] = "错误信息：解析域名时出错：" + err.Error() + "，请联系管理员修复此问题"
					return
				}

				cname := serverDNSResp.DnsName + "." + serverDNSResp.Domain.Name + "."
				for _, serverName := range serverNames {
					if len(serverName.SubNames) == 0 {
						// TODO 可以指定查找解析记录的DNSResolver
						result, err := utils.LookupCNAME(serverName.Name)
						if err != nil {
							m["type"] = "dnsResolveErr"
							m["message"] = "域名解析错误"
							m["todo"] = "错误信息：解析域名'" + serverName.Name + "' CNAME记录时出错：" + err.Error() + "，请修复此问题。如果已经修改，请等待一个小时后再试。如果长时间无法生效，请咨询你的域名DNS服务商。"
							return
						}
						if len(result) == 0 {
							m["type"] = "dnsResolveErr"
							m["message"] = "域名解析错误"
							m["todo"] = "错误信息：找不到域名'" + serverName.Name + "'的CNAME记录，请修复此问题。如果已经修改，请等待一个小时后再试。如果长时间无法生效，请咨询你的域名DNS服务商。"
							return
						}
						if result == serverName.Name+"." {
							m["type"] = "dnsResolveErr"
							m["message"] = "域名解析错误"
							m["todo"] = "错误信息：找不到域名'" + serverName.Name + "'的CNAME记录，请设置为'" + cname + "'。如果已经设置，请等待一个小时后再试。如果长时间无法生效，请咨询你的域名DNS服务商。"
							return
						}
						if result != cname {
							m["type"] = "dnsResolveErr"
							m["message"] = "域名解析错误"
							m["todo"] = "错误信息：解析域名'" + serverName.Name + "' CNAME记录时出错：当前的CNAME值为" + result + "，请修改为" + cname + "。如果已经修改，请等待一个小时后再试。如果长时间无法生效，请咨询你的域名DNS服务商。"
							return
						}
					} else {
						for _, subName := range serverName.SubNames {
							// TODO 可以指定查找解析记录的DNSResolver
							result, err := utils.LookupCNAME(subName)
							if err != nil {
								m["type"] = "dnsResolveErr"
								m["message"] = "域名解析错误"
								m["todo"] = "错误信息：解析域名'" + subName + "' CNAME记录时出错：" + err.Error() + "，请修复此问题。如果已经修改，请等待一个小时后再试。如果长时间无法生效，请咨询你的域名DNS服务商。"
								return
							}
							if len(result) == 0 {
								m["type"] = "dnsResolveErr"
								m["message"] = "域名解析错误"
								m["todo"] = "错误信息：找不到域名'" + subName + "'的CNAME记录，请修复此问题。如果已经修改，请等待一个小时后再试。如果长时间无法生效，请咨询你的域名DNS服务商。"
								return
							}
							if result == cname+"." {
								m["type"] = "dnsResolveErr"
								m["message"] = "域名解析错误"
								m["todo"] = "错误信息：找不到域名'" + serverName.Name + "'的CNAME记录，请设置为'" + cname + "'。如果已经设置，请等待一个小时后再试。如果长时间无法生效，请咨询你的域名DNS服务商。"
								return
							}
							if result != cname {
								m["type"] = "dnsResolveErr"
								m["message"] = "域名解析错误"
								m["todo"] = "错误信息：解析域名'" + subName + "' CNAME记录时出错：当前的CNAME值为" + result + "，请修改为" + cname + "。如果已经修改，请等待一个小时后再试。如果长时间无法生效，请咨询你的域名DNS服务商。"
								return
							}
						}
					}
				}
			}

			m["isOk"] = true
		}(serverId)
	}

	wg.Wait()

	this.Data["status"] = status

	this.Success()
}
