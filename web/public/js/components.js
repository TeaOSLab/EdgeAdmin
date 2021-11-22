// 显示节点的多个集群
Vue.component("node-clusters-labels", {
	props: ["v-primary-cluster", "v-secondary-clusters", "size"],
	data: function () {
		let cluster = this.vPrimaryCluster
		let secondaryClusters = this.vSecondaryClusters
		if (secondaryClusters == null) {
			secondaryClusters = []
		}

		let labelSize = this.size
		if (labelSize == null) {
			labelSize = "small"
		}
		return {
			cluster: cluster,
			secondaryClusters: secondaryClusters,
			labelSize: labelSize
		}
	},
	template: `<div>
	<a v-if="cluster != null" :href="'/clusters/cluster?clusterId=' + cluster.id" title="主集群" style="margin-bottom: 0.3em;">
		<span class="ui label basic grey" :class="labelSize" v-if="labelSize != 'tiny'">{{cluster.name}}</span>
		<grey-label v-if="labelSize == 'tiny'">{{cluster.name}}</grey-label>
	</a>
	<a v-for="c in secondaryClusters" :href="'/clusters/cluster?clusterId=' + c.id" :class="labelSize" title="从集群">
		<span class="ui label basic grey" :class="labelSize" v-if="labelSize != 'tiny'">{{c.name}}</span>
		<grey-label v-if="labelSize == 'tiny'">{{c.name}}</grey-label>
	</a>
</div>`
})

// 单个集群选择
Vue.component("cluster-selector", {
	props: ["v-cluster-id"],
	mounted: function () {
		let that = this

		Tea.action("/clusters/options")
			.post()
			.success(function (resp) {
				that.clusters = resp.data.clusters
			})
	},
	data: function () {
		let clusterId = this.vClusterId
		if (clusterId == null) {
			clusterId = 0
		}
		return {
			clusters: [],
			clusterId: clusterId
		}
	},
	template: `<div>
	<select class="ui dropdown" style="max-width: 10em" name="clusterId" v-model="clusterId">
		<option value="0">[选择集群]</option>
		<option v-for="cluster in clusters" :value="cluster.id">{{cluster.name}}</option>
	</select>
</div>`
})

// 一个节点的多个集群选择器
Vue.component("node-clusters-selector", {
	props: ["v-primary-cluster", "v-secondary-clusters"],
	data: function () {
		let primaryCluster = this.vPrimaryCluster

		let secondaryClusters = this.vSecondaryClusters
		if (secondaryClusters == null) {
			secondaryClusters = []
		}

		return {
			primaryClusterId: (primaryCluster == null) ? 0 : primaryCluster.id,
			secondaryClusterIds: secondaryClusters.map(function (v) {
				return v.id
			}),

			primaryCluster: primaryCluster,
			secondaryClusters: secondaryClusters
		}
	},
	methods: {
		addPrimary: function () {
			let that = this
			let selectedClusterIds = [this.primaryClusterId].concat(this.secondaryClusterIds)
			teaweb.popup("/clusters/selectPopup?selectedClusterIds=" + selectedClusterIds.join(",") + "&mode=single", {
				height: "30em",
				width: "50em",
				callback: function (resp) {
					if (resp.data.cluster != null) {
						that.primaryCluster = resp.data.cluster
						that.primaryClusterId = that.primaryCluster.id
						that.notifyChange()
					}
				}
			})
		},
		removePrimary: function () {
			this.primaryClusterId = 0
			this.primaryCluster = null
			this.notifyChange()
		},
		addSecondary: function () {
			let that = this
			let selectedClusterIds = [this.primaryClusterId].concat(this.secondaryClusterIds)
			teaweb.popup("/clusters/selectPopup?selectedClusterIds=" + selectedClusterIds.join(",") + "&mode=multiple", {
				height: "30em",
				width: "50em",
				callback: function (resp) {
					if (resp.data.cluster != null) {
						that.secondaryClusterIds.push(resp.data.cluster.id)
						that.secondaryClusters.push(resp.data.cluster)
						that.notifyChange()
					}
				}
			})
		},
		removeSecondary: function (index) {
			this.secondaryClusterIds.$remove(index)
			this.secondaryClusters.$remove(index)
			this.notifyChange()
		},
		notifyChange: function () {
			this.$emit("change", {
				clusterId: this.primaryClusterId
			})
		}
	},
	template: `<div>
	<input type="hidden" name="primaryClusterId" :value="primaryClusterId"/>
	<input type="hidden" name="secondaryClusterIds" :value="JSON.stringify(secondaryClusterIds)"/>
	<table class="ui table">
		<tr>
			<td class="title">主集群</td>
			<td>
				<div v-if="primaryCluster != null">
					<div class="ui label basic small">{{primaryCluster.name}} &nbsp; <a href="" title="删除" @click.prevent="removePrimary"><i class="icon remove small"></i></a> </div>
				</div>
				<div style="margin-top: 0.6em" v-if="primaryClusterId == 0">
					<button class="ui button tiny" type="button" @click.prevent="addPrimary">+</button>
				</div>
				<p class="comment">多个集群配置有冲突时，优先使用主集群配置。</p>
			</td>
		</tr>
		<tr>
			<td>从集群</td>
			<td>
				<div v-if="secondaryClusters.length > 0">
					<div class="ui label basic small" v-for="(cluster, index) in secondaryClusters"><span class="grey">{{cluster.name}}</span> &nbsp; <a href="" title="删除" @click.prevent="removeSecondary(index)"><i class="icon remove small"></i></a> </div>
				</div>
				<div style="margin-top: 0.6em">
					<button class="ui button tiny" type="button" @click.prevent="addSecondary">+</button>
				</div>
			</td>
		</tr>
	</table>
</div>`
})

Vue.component("message-media-selector", {
    props: ["v-media-type"],
    mounted: function () {
        let that = this
        Tea.action("/admins/recipients/mediaOptions")
            .post()
            .success(function (resp) {
                that.medias = resp.data.medias

                // 初始化简介
                if (that.mediaType.length > 0) {
                    let media = that.medias.$find(function (_, media) {
                        return media.type == that.mediaType
                    })
                    if (media != null) {
                        that.description = media.description
                    }
                }
            })
    },
    data: function () {
        let mediaType = this.vMediaType
        if (mediaType == null) {
            mediaType = ""
        }
        return {
            medias: [],
            description: "",
            mediaType: mediaType
        }
    },
    watch: {
        mediaType: function (v) {
            let media = this.medias.$find(function (_, media) {
                return media.type == v
            })
            if (media == null) {
                this.description = ""
            } else {
                this.description = media.description
            }
            this.$emit("change", media)
        },
    },
    template: `<div>
    <select class="ui dropdown auto-width" name="mediaType" v-model="mediaType">
        <option value="">[选择媒介类型]</option>
        <option v-for="media in medias" :value="media.type">{{media.name}}</option>
    </select>
    <p class="comment" v-html="description"></p>
</div>`
})

// 消息接收人设置
Vue.component("message-receivers-box", {
	props: ["v-node-cluster-id"],
	mounted: function () {
		let that = this
		Tea.action("/clusters/cluster/settings/message/selectedReceivers")
			.params({
				clusterId: this.clusterId
			})
			.post()
			.success(function (resp) {
				that.receivers = resp.data.receivers
			})
	},
	data: function () {
		let clusterId = this.vNodeClusterId
		if (clusterId == null) {
			clusterId = 0
		}
		return {
			clusterId: clusterId,
			receivers: []
		}
	},
	methods: {
		addReceiver: function () {
			let that = this
			let recipientIdStrings = []
			let groupIdStrings = []
			this.receivers.forEach(function (v) {
				if (v.type == "recipient") {
					recipientIdStrings.push(v.id.toString())
				} else if (v.type == "group") {
					groupIdStrings.push(v.id.toString())
				}
			})

			teaweb.popup("/clusters/cluster/settings/message/selectReceiverPopup?recipientIds=" + recipientIdStrings.join(",") + "&groupIds=" + groupIdStrings.join(","), {
				callback: function (resp) {
					that.receivers.push(resp.data)
				}
			})
		},
		removeReceiver: function (index) {
			this.receivers.$remove(index)
		}
	},
	template: `<div>
        <input type="hidden" name="receiversJSON" :value="JSON.stringify(receivers)"/>           
        <div v-if="receivers.length > 0">
            <div v-for="(receiver, index) in receivers" class="ui label basic small">
               <span v-if="receiver.type == 'group'">分组：</span>{{receiver.name}} <span class="grey small" v-if="receiver.subName != null && receiver.subName.length > 0">({{receiver.subName}})</span> &nbsp; <a href="" title="删除" @click.prevent="removeReceiver(index)"><i class="icon remove"></i></a>
            </div>
             <div class="ui divider"></div>
        </div>
      <button type="button" class="ui button tiny" @click.prevent="addReceiver">+</button>
</div>`
})

Vue.component("message-recipient-group-selector", {
    props: ["v-groups"],
    data: function () {
        let groups = this.vGroups
        if (groups == null) {
            groups = []
        }
        let groupIds = []
        if (groups.length > 0) {
            groupIds = groups.map(function (v) {
                return v.id.toString()
            }).join(",")
        }

        return {
            groups: groups,
            groupIds: groupIds
        }
    },
    methods: {
        addGroup: function () {
            let that = this
            teaweb.popup("/admins/recipients/groups/selectPopup?groupIds=" + this.groupIds, {
                callback: function (resp) {
                    that.groups.push(resp.data.group)
                    that.update()
                }
            })
        },
        removeGroup: function (index) {
            this.groups.$remove(index)
            this.update()
        },
        update: function () {
            let groupIds = []
            if (this.groups.length > 0) {
                this.groups.forEach(function (v) {
                    groupIds.push(v.id)
                })
            }
            this.groupIds = groupIds.join(",")
        }
    },
    template: `<div>
    <input type="hidden" name="groupIds" :value="groupIds"/>
    <div v-if="groups.length > 0">
        <div>
            <div v-for="(group, index) in groups" class="ui label small basic">
                {{group.name}} &nbsp; <a href="" title="删除" @click.prevent="removeGroup(index)"><i class="icon remove"></i></a>
            </div>
        </div>
        <div class="ui divider"></div>
    </div>   
    <button class="ui button tiny" type="button" @click.prevent="addGroup()">+</button>
</div>`
})

Vue.component("message-media-instance-selector", {
    props: ["v-instance-id"],
    mounted: function () {
        let that = this
        Tea.action("/admins/recipients/instances/options")
            .post()
            .success(function (resp) {
                that.instances = resp.data.instances

                // 初始化简介
                if (that.instanceId > 0) {
                    let instance = that.instances.$find(function (_, instance) {
                        return instance.id == that.instanceId
                    })
                    if (instance != null) {
                        that.description = instance.description
                        that.update(instance.id)
                    }
                }
            })
    },
    data: function () {
        let instanceId = this.vInstanceId
        if (instanceId == null) {
            instanceId = 0
        }
        return {
            instances: [],
            description: "",
            instanceId: instanceId
        }
    },
    watch: {
        instanceId: function (v) {
            this.update(v)
        }
    },
    methods: {
        update: function (v) {
            let instance = this.instances.$find(function (_, instance) {
                return instance.id == v
            })
            if (instance == null) {
                this.description = ""
            } else {
                this.description = instance.description
            }
            this.$emit("change", instance)
        }
    },
    template: `<div>
    <select class="ui dropdown auto-width" name="instanceId" v-model="instanceId">
        <option value="0">[选择媒介]</option>
        <option v-for="instance in instances" :value="instance.id">{{instance.name}} ({{instance.media.name}})</option>
    </select>
    <p class="comment" v-html="description"></p>
</div>`
})

Vue.component("message-row", {
	props: ["v-message", "v-can-close"],
	data: function () {
		let paramsJSON = this.vMessage.params
		let params = null
		if (paramsJSON != null && paramsJSON.length > 0) {
			params = JSON.parse(paramsJSON)
		}

		return {
			message: this.vMessage,
			params: params,
			isClosing: false
		}
	},
	methods: {
		viewCert: function (certId) {
			teaweb.popup("/servers/certs/certPopup?certId=" + certId, {
				height: "28em",
				width: "48em"
			})
		},
		readMessage: function (messageId) {
			let that = this

			Tea.action("/messages/readPage")
				.params({"messageIds": [messageId]})
				.post()
				.success(function () {
					// 刷新父级页面Badge
					if (window.parent.Tea != null && window.parent.Tea.Vue != null) {
						window.parent.Tea.Vue.checkMessagesOnce()
					}

					// 刷新当前页面
					if (that.vCanClose && typeof (NotifyPopup) != "undefined") {
						that.isClosing = true
						setTimeout(function () {
							NotifyPopup({})
						}, 1000)
					} else {
						teaweb.reload()
					}
				})
		}
	},
	template: `<div>
<table class="ui table selectable" v-if="!isClosing">
	<tr :class="{error: message.level == 'error', positive: message.level == 'success', warning: message.level == 'warning'}">
		<td style="position: relative">
			<strong>{{message.datetime}}</strong>
			<span v-if="message.cluster != null && message.cluster.id != null">
				<span> | </span>
				<a :href="'/clusters/cluster?clusterId=' + message.cluster.id" target="_top" v-if="message.role == 'node'">集群：{{message.cluster.name}}</a>
				<a :href="'/ns/clusters/cluster?clusterId=' + message.cluster.id" target="_top" v-if="message.role == 'dns'">DNS集群：{{message.cluster.name}}</a>
			</span>
			<span v-if="message.node != null && message.node.id != null">
				<span> | </span>
				<a :href="'/clusters/cluster/node?clusterId=' + message.cluster.id + '&nodeId=' + message.node.id" target="_top" v-if="message.role == 'node'">节点：{{message.node.name}}</a>
				<a :href="'/ns/clusters/cluster/node?clusterId=' + message.cluster.id + '&nodeId=' + message.node.id" target="_top" v-if="message.role == 'dns'">DNS节点：{{message.node.name}}</a>
			</span>
			<a href=""  style="position: absolute; right: 1em" @click.prevent="readMessage(message.id)" title="标为已读"><i class="icon check"></i></a>
		</td>
	</tr>
	<tr :class="{error: message.level == 'error', positive: message.level == 'success', warning: message.level == 'warning'}">
		<td>
			{{message.body}}
			
			<!-- 健康检查 -->
			<div v-if="message.type == 'HealthCheckFailed'" style="margin-top: 0.8em">
				<a :href="'/clusters/cluster/node?clusterId=' + message.cluster.id + '&nodeId=' + param.node.id" v-for="param in params" class="ui label small basic" style="margin-bottom: 0.5em" target="_top">{{param.node.name}}: {{param.error}}</a>
			</div>
			
			<!-- 集群DNS设置 -->
			<div v-if="message.type == 'ClusterDNSSyncFailed'" style="margin-top: 0.8em">
				<a :href="'/dns/clusters/cluster?clusterId=' + message.cluster.id" target="_top">查看问题 &raquo;</a>
			</div>
			
			<!-- 证书即将过期 -->
			<div v-if="message.type == 'SSLCertExpiring'" style="margin-top: 0.8em">
				<a href="" @click.prevent="viewCert(params.certId)" target="_top">查看证书</a> &nbsp;|&nbsp; <a :href="'/servers/certs/acme'" v-if="params != null && params.acmeTaskId > 0" target="_top">查看任务&raquo;</a>
			</div>
			
			<!-- 证书续期成功 -->
			<div v-if="message.type == 'SSLCertACMETaskSuccess'" style="margin-top: 0.8em">
				<a href="" @click.prevent="viewCert(params.certId)" target="_top">查看证书</a> &nbsp;|&nbsp; <a :href="'/servers/certs/acme'" v-if="params != null && params.acmeTaskId > 0" target="_top">查看任务&raquo;</a>
			</div>
			
			<!-- 证书续期失败 -->
			<div v-if="message.type == 'SSLCertACMETaskFailed'" style="margin-top: 0.8em">
				<a href="" @click.prevent="viewCert(params.certId)" target="_top">查看证书</a> &nbsp;|&nbsp; <a :href="'/servers/certs/acme'" v-if="params != null && params.acmeTaskId > 0" target="_top">查看任务&raquo;</a>
			</div>
		</td>
	</tr>
</table>
<div class="margin"></div>
</div>`
})

// 选择多个线路
Vue.component("ns-routes-selector", {
	props: ["v-routes"],
	mounted: function () {
		let that = this
		Tea.action("/ns/routes/options")
			.post()
			.success(function (resp) {
				that.routes = resp.data.routes
			})
	},
	data: function () {
		let selectedRoutes = this.vRoutes
		if (selectedRoutes == null) {
			selectedRoutes = []
		}

		return {
			routeCode: "default",
			routes: [],
			isAdding: false,
			routeType: "default",
			selectedRoutes: selectedRoutes
		}
	},
	watch: {
		routeType: function (v) {
			this.routeCode = ""
			let that = this
			this.routes.forEach(function (route) {
				if (route.type == v && that.routeCode.length == 0) {
					that.routeCode = route.code
				}
			})
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			this.routeType = "default"
			this.routeCode = "default"
		},
		cancel: function () {
			this.isAdding = false
		},
		confirm: function () {
			if (this.routeCode.length == 0) {
				return
			}

			let that = this
			this.routes.forEach(function (v) {
				if (v.code == that.routeCode) {
					that.selectedRoutes.push(v)
				}
			})
			this.cancel()
		},
		remove: function (index) {
			this.selectedRoutes.$remove(index)
		}
	}
	,
	template: `<div>
	<div>
		<div class="ui label basic text small" v-for="(route, index) in selectedRoutes" style="margin-bottom: 0.3em">
			<input type="hidden" name="routeCodes" :value="route.code"/>
			{{route.name}} &nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding" style="margin-bottom: 1em">
		<div class="ui fields inline">
			<div class="ui field">
				<select class="ui dropdown" v-model="routeType">
					<option value="default">[默认线路]</option>
					<option value="user">自定义线路</option>
					<option value="isp">运营商</option>
					<option value="china">中国省市</option>
					<option value="world">全球国家地区</option>
				</select>
			</div>
			
			<div class="ui field">
				<select class="ui dropdown" v-model="routeCode" style="width: 10em">
					<option v-for="route in routes" :value="route.code" v-if="route.type == routeType">{{route.name}}</option>
				</select>
			</div>
			
			<div class="ui field">
				<button type="button" class="ui button tiny" @click.prevent="confirm">确定</button>
				&nbsp; <a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
	</div>
	<button class="ui button tiny" type="button" @click.prevent="add">+</button>
</div>`
})

// 递归DNS设置
Vue.component("ns-recursion-config-box", {
	props: ["v-recursion-config"],
	data: function () {
		let recursion = this.vRecursionConfig
		if (recursion == null) {
			recursion = {
				isOn: false,
				hosts: [],
				allowDomains: [],
				denyDomains: [],
				useLocalHosts: false
			}
		}
		if (recursion.hosts == null) {
			recursion.hosts = []
		}
		if (recursion.allowDomains == null) {
			recursion.allowDomains = []
		}
		if (recursion.denyDomains == null) {
			recursion.denyDomains = []
		}
		return {
			config: recursion,
			hostIsAdding: false,
			host: "",
			updatingHost: null
		}
	},
	methods: {
		changeHosts: function (hosts) {
			this.config.hosts = hosts
		},
		changeAllowDomains: function (domains) {
			this.config.allowDomains = domains
		},
		changeDenyDomains: function (domains) {
			this.config.denyDomains = domains
		},
		removeHost: function (index) {
			this.config.hosts.$remove(index)
		},
		addHost: function () {
			this.updatingHost = null
			this.host = ""
			this.hostIsAdding = !this.hostIsAdding
			if (this.hostIsAdding) {
				var that = this
				setTimeout(function () {
					let hostRef = that.$refs.hostRef
					if (hostRef != null) {
						hostRef.focus()
					}
				}, 200)
			}
		},
		updateHost: function (host) {
			this.updatingHost = host
			this.host = host.host
			this.hostIsAdding = !this.hostIsAdding

			if (this.hostIsAdding) {
				var that = this
				setTimeout(function () {
					let hostRef = that.$refs.hostRef
					if (hostRef != null) {
						hostRef.focus()
					}
				}, 200)
			}
		},
		confirmHost: function () {
			if (this.host.length == 0) {
				teaweb.warn("请输入DNS地址")
				return
			}

			// TODO 校验Host
			// TODO 可以输入端口号
			// TODO 可以选择协议

			this.hostIsAdding = false
			if (this.updatingHost == null) {
				this.config.hosts.push({
					host: this.host
				})
			} else {
				this.updatingHost.host = this.host
			}
		},
		cancelHost: function () {
			this.hostIsAdding = false
		}
	},
	template: `<div>
	<input type="hidden" name="recursionJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<tbody>
			<tr>
				<td class="title">是否启用</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" name="isOn" value="1" v-model="config.isOn"/>
						<label></label>
					</div>
					<p class="comment">启用后，如果找不到某个域名的解析记录，则向上一级DNS查找。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="config.isOn">
			<tr>
				<td>从节点本机读取<br/>上级DNS主机</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" name="useLocalHosts" value="1" v-model="config.useLocalHosts"/>
						<label></label>
					</div>
					<p class="comment">选中后，节点会试图从<code-label>/etc/resolv.conf</code-label>文件中读取DNS配置。 </p>
				</td>
			</tr>
			<tr v-show="!config.useLocalHosts">
				<td>上级DNS主机地址 *</td>
				<td>
					<div v-if="config.hosts.length > 0">
						<div v-for="(host, index) in config.hosts" class="ui label tiny basic">
							{{host.host}} &nbsp;
							<a href="" title="修改" @click.prevent="updateHost(host)"><i class="icon pencil tiny"></i></a>
							<a href="" title="删除" @click.prevent="removeHost(index)"><i class="icon remove small"></i></a>
						</div>
						<div class="ui divider"></div>
					</div>
					<div v-if="hostIsAdding">
						<div class="ui fields inline">
							<div class="ui field">
								<input type="text" placeholder="DNS主机地址" v-model="host" ref="hostRef" @keyup.enter="confirmHost" @keypress.enter.prevent="1"/>
							</div>
							<div class="ui field">
								<button class="ui button tiny" type="button" @click.prevent="confirmHost">确认</button> &nbsp; <a href="" title="取消" @click.prevent="cancelHost"><i class="icon remove small"></i></a>
							</div>
						</div>
					</div>
					<div style="margin-top: 0.5em">
						<button type="button" class="ui button tiny" @click.prevent="addHost">+</button>
					</div>
				</td>
			</tr>
			<tr>
				<td>允许的域名</td>
				<td><values-box name="allowDomains" :values="config.allowDomains" @change="changeAllowDomains"></values-box>
					<p class="comment">支持星号通配符，比如<code-label>*.example.org</code-label>。</p>
				</td>
			</tr>
			<tr>
				<td>不允许的域名</td>
				<td>
					<values-box name="denyDomains" :values="config.denyDomains" @change="changeDenyDomains"></values-box>
					<p class="comment">支持星号通配符，比如<code-label>*.example.org</code-label>。优先级比允许的域名高。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("ns-access-log-ref-box", {
	props: ["v-access-log-ref", "v-is-parent"],
	data: function () {
		let config = this.vAccessLogRef
		if (config == null) {
			config = {
				isOn: false,
				isPrior: false,
				logMissingDomains: false
			}
		}
		if (typeof (config.logMissingDomains) == "undefined") {
			config.logMissingDomains = false
		}
		return {
			config: config
		}
	},
	template: `<div>
	<input type="hidden" name="accessLogJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="config" v-if="!vIsParent"></prior-checkbox>
		<tbody v-show="vIsParent || config.isPrior">
			<tr>
				<td class="title">是否启用</td>
				<td>
					<checkbox name="isOn" value="1" v-model="config.isOn"></checkbox>
				</td>
			</tr>
			<tr>
				<td>记录所有访问</td>
				<td>
					<checkbox name="logMissingDomains" value="1" v-model="config.logMissingDomains"></checkbox>
					<p class="comment">包括对没有在系统里创建的域名访问。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("ns-route-ranges-box", {
	props: ["v-ranges"],
	data: function () {
		let ranges = this.vRanges
		if (ranges == null) {
			ranges = []
		}
		return {
			ranges: ranges,
			isAdding: false,

			// IP范围
			ipRangeFrom: "",
			ipRangeTo: ""
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.ipRangeFrom.focus()
			}, 100)
		},
		remove: function (index) {
			this.ranges.$remove(index)
		},
		cancelIPRange: function () {
			this.isAdding = false
			this.ipRangeFrom = ""
			this.ipRangeTo = ""
		},
		confirmIPRange: function () {
			// 校验IP
			let that = this
			this.ipRangeFrom = this.ipRangeFrom.trim()
			if (!this.validateIP(this.ipRangeFrom)) {
				teaweb.warn("开始IP填写错误", function () {
					that.$refs.ipRangeFrom.focus()
				})
				return
			}

			this.ipRangeTo = this.ipRangeTo.trim()
			if (!this.validateIP(this.ipRangeTo)) {
				teaweb.warn("结束IP填写错误", function () {
					that.$refs.ipRangeTo.focus()
				})
				return
			}

			this.ranges.push({
				type: "ipRange",
				params: {
					ipFrom: this.ipRangeFrom,
					ipTo: this.ipRangeTo
				}
			})
			this.cancelIPRange()
		},
		validateIP: function (ip) {
			if (!ip.match(/^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$/)) {
				return false
			}
			let pieces = ip.split(".")
			let isOk = true
			pieces.forEach(function (v) {
				let v1 = parseInt(v)
				if (v1 > 255) {
					isOk = false
				}
			})
			return isOk
		}
	},
	template: `<div>
	<input type="hidden" name="rangesJSON" :value="JSON.stringify(ranges)"/>
	<div v-if="ranges.length > 0">
		<div class="ui label tiny basic" v-for="(range, index) in ranges">
			<span v-if="range.type == 'ipRange'">IP范围：</span>
			{{range.params.ipFrom}} - {{range.params.ipTo}} &nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	
	<!-- IP 范围 -->
	<div style="margin-bottom: 1em" v-show="isAdding">
		<div class="ui fields inline">
			<div class="ui field">
				<input type="text" placeholder="开始IP" maxlength="15" size="15" v-model="ipRangeFrom" ref="ipRangeFrom"  @keyup.enter="confirmIPRange" @keypress.enter.prevent="1"/>
			</div>
			<div class="ui field">-</div>
			<div class="ui field">
				<input type="text" placeholder="结束IP" maxlength="15" size="15" v-model="ipRangeTo" ref="ipRangeTo" @keyup.enter="confirmIPRange" @keypress.enter.prevent="1"/>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirmIPRange">确定</button> &nbsp;
				<a href="" @click.prevent="cancelIPRange" title="取消"><i class="icon remove small"></i></a>
			</div>
		</div>
	</div>
	
	<button class="ui button tiny" type="button" @click.prevent="add">+</button>
</div>`
})

// 选择单一线路
Vue.component("ns-route-selector", {
	props: ["v-route-code"],
	mounted: function () {
		let that = this
		Tea.action("/ns/routes/options")
			.post()
			.success(function (resp) {
				that.routes = resp.data.routes
			})
	},
	data: function () {
		let routeCode = this.vRouteCode
		if (routeCode == null) {
			routeCode = ""
		}
		return {
			routeCode: routeCode,
			routes: []
		}
	},
	template: `<div>
	<div v-if="routes.length > 0">
		<select class="ui dropdown" name="routeCode" v-model="routeCode">
			<option value="">[线路]</option>
			<option v-for="route in routes" :value="route.code">{{route.name}}</option>
		</select>
	</div>
</div>`
})

Vue.component("ns-user-selector", {
	mounted: function () {
		let that = this

		Tea.action("/ns/users/options")
			.post()
			.success(function (resp) {
				that.users = resp.data.users
			})
	},
	props: ["v-user-id"],
	data: function () {
		let userId = this.vUserId
		if (userId == null) {
			userId = 0
		}
		return {
			users: [],
			userId: userId
		}
	},
	template: `<div>
	<select class="ui dropdown auto-width" name="userId" v-model="userId">
		<option value="0">[选择用户]</option>
		<option v-for="user in users" :value="user.id">{{user.fullname}} ({{user.username}})</option>
	</select>
</div>`
})

Vue.component("ns-access-log-box", {
	props: ["v-access-log", "v-keyword"],
	data: function () {
		let accessLog = this.vAccessLog
		return {
			accessLog: accessLog
		}
	},
	methods: {
		showLog: function () {
			let that = this
			let requestId = this.accessLog.requestId
			this.$parent.$children.forEach(function (v) {
				if (v.deselect != null) {
					v.deselect()
				}
			})
			this.select()

			teaweb.popup("/ns/clusters/accessLogs/viewPopup?requestId=" + requestId, {
				width: "50em",
				height: "24em",
				onClose: function () {
					that.deselect()
				}
			})
		},
		select: function () {
			this.$refs.box.parentNode.style.cssText = "background: rgba(0, 0, 0, 0.1)"
		},
		deselect: function () {
			this.$refs.box.parentNode.style.cssText = ""
		}
	},
	template: `<div class="access-log-row" :style="{'color': (!accessLog.isRecursive && (accessLog.nsRecordId == null || accessLog.nsRecordId == 0) || (accessLog.isRecursive && accessLog.recordValue != null && accessLog.recordValue.length == 0)) ? '#dc143c' : ''}" ref="box">
	<span v-if="accessLog.region != null && accessLog.region.length > 0" class="grey">[{{accessLog.region}}]</span> <keyword :v-word="vKeyword">{{accessLog.remoteAddr}}</keyword> [{{accessLog.timeLocal}}] [{{accessLog.networking}}] <em>{{accessLog.questionType}} <keyword :v-word="vKeyword">{{accessLog.questionName}}</keyword></em> -&gt; <em>{{accessLog.recordType}} <keyword :v-word="vKeyword">{{accessLog.recordValue}}</keyword></em><!-- &nbsp; <a href="" @click.prevent="showLog" title="查看详情"><i class="icon expand"></i></a>-->
	<div v-if="(accessLog.nsRoutes != null && accessLog.nsRoutes.length > 0) || accessLog.isRecursive" style="margin-top: 0.3em">
		<span class="ui label tiny basic grey" v-for="route in accessLog.nsRoutes">线路: {{route.name}}</span>
		<span class="ui label tiny basic grey" v-if="accessLog.isRecursive">递归DNS</span>
	</div>
	<div v-if="accessLog.error != null && accessLog.error.length > 0" style="color:#dc143c">
		<i class="icon warning circle"></i>错误：[{{accessLog.error}}]
	</div>
</div>`
})

Vue.component("ns-cluster-selector", {
	props: ["v-cluster-id"],
	mounted: function () {
		let that = this

		Tea.action("/ns/clusters/options")
			.post()
			.success(function (resp) {
				that.clusters = resp.data.clusters
			})
	},
	data: function () {
		let clusterId = this.vClusterId
		if (clusterId == null) {
			clusterId = 0
		}
		return {
			clusters: [],
			clusterId: clusterId
		}
	},
	template: `<div>
	<select class="ui dropdown auto-width" name="clusterId" v-model="clusterId">
		<option value="0">[选择集群]</option>
		<option v-for="cluster in clusters" :value="cluster.id">{{cluster.name}}</option>
	</select>
</div>`
})

Vue.component("plan-user-selector", {
	mounted: function () {
		let that = this

		Tea.action("/plans/users/options")
			.post()
			.success(function (resp) {
				that.users = resp.data.users
			})
	},
	props: ["v-user-id"],
	data: function () {
		let userId = this.vUserId
		if (userId == null) {
			userId = 0
		}
		return {
			users: [],
			userId: userId
		}
	},
	watch: {
		userId: function (v) {
			this.$emit("change", v)
		}
	},
	template: `<div>
	<select class="ui dropdown auto-width" name="userId" v-model="userId">
		<option value="0">[选择用户]</option>
		<option v-for="user in users" :value="user.id">{{user.fullname}} ({{user.username}})</option>
	</select>
</div>`
})

Vue.component("plan-price-view", {
	props: ["v-plan"],
	data: function () {
		return {
			plan: this.vPlan
		}
	},
	template: `<div>
	 <span v-if="plan.priceType == 'period'">
		<span v-if="plan.monthlyPrice > 0">月度：￥{{plan.monthlyPrice}}元<br/></span>
		<span v-if="plan.seasonallyPrice > 0">季度：￥{{plan.seasonallyPrice}}元<br/></span>
		<span v-if="plan.yearlyPrice > 0">年度：￥{{plan.yearlyPrice}}元</span>
	</span>
	<span v-if="plan.priceType == 'traffic'">
		基础价格：￥{{plan.trafficPrice.base}}元/GB
	</span>
</div>`
})

// 套餐价格配置
Vue.component("plan-price-config-box", {
	props: ["v-price-type", "v-monthly-price", "v-seasonally-price", "v-yearly-price", "v-traffic-price"],
	data: function () {
		let priceType = this.vPriceType
		if (priceType == null) {
			priceType = "period"
		}

		let monthlyPriceNumber = 0
		let monthlyPrice = this.vMonthlyPrice
		if (monthlyPrice == null || monthlyPrice <= 0) {
			monthlyPrice = ""
		} else {
			monthlyPrice = monthlyPrice.toString()
			monthlyPriceNumber = parseFloat(monthlyPrice)
			if (isNaN(monthlyPriceNumber)) {
				monthlyPriceNumber = 0
			}
		}

		let seasonallyPriceNumber = 0
		let seasonallyPrice = this.vSeasonallyPrice
		if (seasonallyPrice == null || seasonallyPrice <= 0) {
			seasonallyPrice = ""
		} else {
			seasonallyPrice = seasonallyPrice.toString()
			seasonallyPriceNumber = parseFloat(seasonallyPrice)
			if (isNaN(seasonallyPriceNumber)) {
				seasonallyPriceNumber = 0
			}
		}

		let yearlyPriceNumber = 0
		let yearlyPrice = this.vYearlyPrice
		if (yearlyPrice == null || yearlyPrice <= 0) {
			yearlyPrice = ""
		} else {
			yearlyPrice = yearlyPrice.toString()
			yearlyPriceNumber = parseFloat(yearlyPrice)
			if (isNaN(yearlyPriceNumber)) {
				yearlyPriceNumber = 0
			}
		}

		let trafficPrice = this.vTrafficPrice
		let trafficPriceBaseNumber = 0
		if (trafficPrice != null) {
			trafficPriceBaseNumber = trafficPrice.base
		} else {
			trafficPrice = {
				base: 0
			}
		}
		let trafficPriceBase = ""
		if (trafficPriceBaseNumber > 0) {
			trafficPriceBase = trafficPriceBaseNumber.toString()
		}

		return {
			priceType: priceType,
			monthlyPrice: monthlyPrice,
			seasonallyPrice: seasonallyPrice,
			yearlyPrice: yearlyPrice,

			monthlyPriceNumber: monthlyPriceNumber,
			seasonallyPriceNumber: seasonallyPriceNumber,
			yearlyPriceNumber: yearlyPriceNumber,

			trafficPriceBase: trafficPriceBase,
			trafficPrice: trafficPrice
		}
	},
	watch: {
		monthlyPrice: function (v) {
			let price = parseFloat(v)
			if (isNaN(price)) {
				price = 0
			}
			this.monthlyPriceNumber = price
		},
		seasonallyPrice: function (v) {
			let price = parseFloat(v)
			if (isNaN(price)) {
				price = 0
			}
			this.seasonallyPriceNumber = price
		},
		yearlyPrice: function (v) {
			let price = parseFloat(v)
			if (isNaN(price)) {
				price = 0
			}
			this.yearlyPriceNumber = price
		},
		trafficPriceBase: function (v) {
			let price = parseFloat(v)
			if (isNaN(price)) {
				price = 0
			}
			this.trafficPrice.base = price
		}
	},
	template: `<div>
	<input type="hidden" name="priceType" :value="priceType"/>
	<input type="hidden" name="monthlyPrice" :value="monthlyPriceNumber"/>
	<input type="hidden" name="seasonallyPrice" :value="seasonallyPriceNumber"/>
	<input type="hidden" name="yearlyPrice" :value="yearlyPriceNumber"/>
	<input type="hidden" name="trafficPriceJSON" :value="JSON.stringify(trafficPrice)"/>
	
	<div>
		<radio :v-value="'period'" :value="priceType" v-model="priceType">&nbsp;按时间周期</radio> &nbsp; &nbsp;
		<radio :v-value="'traffic'" :value="priceType" v-model="priceType">&nbsp;按流量</radio>
	</div>
	
	<!-- 按时间周期 -->
	<div v-show="priceType == 'period'">
		<div class="ui divider"></div>
		<table class="ui table">
			<tr>
				<td class="title">月度价格</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 7em" maxlength="10" v-model="monthlyPrice"/>
						<span class="ui label">元</span>
					</div>
				</td>
			</tr>
			<tr>
				<td class="title">季度价格</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 7em" maxlength="10" v-model="seasonallyPrice"/>
						<span class="ui label">元</span>
					</div>
				</td>
			</tr>
			<tr>
				<td class="title">年度价格</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 7em" maxlength="10" v-model="yearlyPrice"/>
						<span class="ui label">元</span>
					</div>
				</td>
			</tr>
		</table>
	</div>
	
	<!-- 按流量 -->
	<div v-show="priceType =='traffic'">
		<div class="ui divider"></div>
		<table class="ui table">
			<tr>
				<td class="title">基础流量费用</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" v-model="trafficPriceBase" maxlength="10" style="width: 7em"/>
						<span class="ui label">元/GB</span>
					</div>
				</td>
			</tr>
		</table>
	</div>
</div>`
})

Vue.component("http-stat-config-box", {
	props: ["v-stat-config", "v-is-location", "v-is-group"],
	data: function () {
		let stat = this.vStatConfig
		if (stat == null) {
			stat = {
				isPrior: false,
				isOn: false
			}
		}
		return {
			stat: stat
		}
	},
	template: `<div>
	<input type="hidden" name="statJSON" :value="JSON.stringify(stat)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="stat" v-if="vIsLocation || vIsGroup" ></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || stat.isPrior">
			<tr>
				<td class="title">是否开启统计</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="stat.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
	</table>
<div class="margin"></div>
</div>`
})

Vue.component("http-request-conds-box", {
	props: ["v-conds"],
	data: function () {
		let conds = this.vConds
		if (conds == null) {
			conds = {
				isOn: true,
				connector: "or",
				groups: []
			}
		}
		return {
			conds: conds,
			components: window.REQUEST_COND_COMPONENTS
		}
	},
	methods: {
		change: function () {
			this.$emit("change", this.conds)
		},
		addGroup: function () {
			window.UPDATING_COND_GROUP = null

			let that = this
			teaweb.popup("/servers/server/settings/conds/addGroupPopup", {
				height: "30em",
				callback: function (resp) {
					that.conds.groups.push(resp.data.group)
					that.change()
				}
			})
		},
		updateGroup: function (groupIndex, group) {
			window.UPDATING_COND_GROUP = group
			let that = this
			teaweb.popup("/servers/server/settings/conds/addGroupPopup", {
				height: "30em",
				callback: function (resp) {
					Vue.set(that.conds.groups, groupIndex, resp.data.group)
					that.change()
				}
			})
		},
		removeGroup: function (groupIndex) {
			let that = this
			teaweb.confirm("确定要删除这一组条件吗？", function () {
				that.conds.groups.$remove(groupIndex)
				that.change()
			})
		},
		typeName: function (cond) {
			let c = this.components.$find(function (k, v) {
				return v.type == cond.type
			})
			if (c != null) {
				return c.name;
			}
			return cond.param + " " + cond.operator
		}
	},
	template: `<div>
		<input type="hidden" name="condsJSON" :value="JSON.stringify(conds)"/>
		<div v-if="conds.groups.length > 0">
			<table class="ui table">
				<tr v-for="(group, groupIndex) in conds.groups">
					<td class="title" :class="{'color-border':conds.connector == 'and'}" :style="{'border-bottom':(groupIndex < conds.groups.length-1) ? '1px solid rgba(34,36,38,.15)':''}">分组{{groupIndex+1}}</td>
					<td style="background: white;" :style="{'border-bottom':(groupIndex < conds.groups.length-1) ? '1px solid rgba(34,36,38,.15)':''}">
						<var v-for="(cond, index) in group.conds" style="font-style: normal;display: inline-block; margin-bottom:0.5em">
							<span class="ui label tiny">
								<var v-if="cond.type.length == 0 || cond.type == 'params'" style="font-style: normal">{{cond.param}} <var>{{cond.operator}}</var></var>
								<var v-if="cond.type.length > 0 && cond.type != 'params'" style="font-style: normal">{{typeName(cond)}}: </var>
								{{cond.value}}
							</span>
							
							<var v-if="index < group.conds.length - 1"> {{group.connector}} &nbsp;</var>
						</var>
					</td>
					<td style="width: 5em; background: white" :style="{'border-bottom':(groupIndex < conds.groups.length-1) ? '1px solid rgba(34,36,38,.15)':''}">
						<a href="" title="修改分组" @click.prevent="updateGroup(groupIndex, group)"><i class="icon pencil small"></i></a> <a href="" title="删除分组" @click.prevent="removeGroup(groupIndex)"><i class="icon remove"></i></a>
					</td>
				</tr>
			</table>
			<div class="ui divider"></div>
		</div>
		
		<!-- 分组之间关系 -->
		<table class="ui table" v-if="conds.groups.length > 1">
			<tr>
				<td class="title">分组之间关系</td>
				<td>
					<select class="ui dropdown auto-width" v-model="conds.connector">
						<option value="and">和</option>
						<option value="or">或</option>
					</select>
					<p class="comment">
						<span v-if="conds.connector == 'or'">只要满足其中一个条件分组即可。</span>
						<span v-if="conds.connector == 'and'">需要满足所有条件分组。</span>
					</p>	
				</td>
			</tr>
		</table>
		
		<div>
			<button class="ui button tiny" type="button" @click.prevent="addGroup()">+添加分组</button>
		</div>
	</div>	
</div>`
})

Vue.component("ssl-config-box", {
	props: ["v-ssl-policy", "v-protocol", "v-server-id"],
	created: function () {
		let that = this
		setTimeout(function () {
			that.sortableCipherSuites()
		}, 100)
	},
	data: function () {
		let policy = this.vSslPolicy
		if (policy == null) {
			policy = {
				id: 0,
				isOn: true,
				certRefs: [],
				certs: [],
				clientCARefs: [],
				clientCACerts: [],
				clientAuthType: 0,
				minVersion: "TLS 1.1",
				hsts: null,
				cipherSuitesIsOn: false,
				cipherSuites: [],
				http2Enabled: true
			}
		} else {
			if (policy.certRefs == null) {
				policy.certRefs = []
			}
			if (policy.certs == null) {
				policy.certs = []
			}
			if (policy.clientCARefs == null) {
				policy.clientCARefs = []
			}
			if (policy.clientCACerts == null) {
				policy.clientCACerts = []
			}
			if (policy.cipherSuites == null) {
				policy.cipherSuites = []
			}
		}

		let hsts = policy.hsts
		if (hsts == null) {
			hsts = {
				isOn: false,
				maxAge: 0,
				includeSubDomains: false,
				preload: false,
				domains: []
			}
		}

		return {
			policy: policy,

			// hsts
			hsts: hsts,
			hstsOptionsVisible: false,
			hstsDomainAdding: false,
			addingHstsDomain: "",
			hstsDomainEditingIndex: -1,

			// 相关数据
			allVersions: window.SSL_ALL_VERSIONS,
			allCipherSuites: window.SSL_ALL_CIPHER_SUITES.$copy(),
			modernCipherSuites: window.SSL_MODERN_CIPHER_SUITES,
			intermediateCipherSuites: window.SSL_INTERMEDIATE_CIPHER_SUITES,
			allClientAuthTypes: window.SSL_ALL_CLIENT_AUTH_TYPES,
			cipherSuitesVisible: false,

			// 高级选项
			moreOptionsVisible: false
		}
	},
	watch: {
		hsts: {
			deep: true,
			handler: function () {
				this.policy.hsts = this.hsts
			}
		}
	},
	methods: {
		// 删除证书
		removeCert: function (index) {
			let that = this
			teaweb.confirm("确定删除此证书吗？证书数据仍然保留，只是当前服务不再使用此证书。", function () {
				that.policy.certRefs.$remove(index)
				that.policy.certs.$remove(index)
			})
		},

		// 选择证书
		selectCert: function () {
			let that = this
			let selectedCertIds = []
			if (this.policy != null && this.policy.certs.length > 0) {
				this.policy.certs.forEach(function (cert) {
					selectedCertIds.push(cert.id.toString())
				})
			}
			teaweb.popup("/servers/certs/selectPopup?selectedCertIds=" + selectedCertIds, {
				width: "50em",
				height: "30em",
				callback: function (resp) {
					that.policy.certRefs.push(resp.data.certRef)
					that.policy.certs.push(resp.data.cert)
				}
			})
		},

		// 上传证书
		uploadCert: function () {
			let that = this
			teaweb.popup("/servers/certs/uploadPopup", {
				height: "28em",
				callback: function (resp) {
					teaweb.success("上传成功", function () {
						that.policy.certRefs.push(resp.data.certRef)
						that.policy.certs.push(resp.data.cert)
					})
				}
			})
		},

		// 申请证书
		requestCert: function () {
			// 已经在证书中的域名
			let excludeServerNames = []
			if (this.policy != null && this.policy.certs.length > 0) {
				this.policy.certs.forEach(function (cert) {
					excludeServerNames.$pushAll(cert.dnsNames)
				})
			}

			let that = this
			teaweb.popup("/servers/server/settings/https/requestCertPopup?serverId=" + this.vServerId + "&excludeServerNames=" + excludeServerNames.join(","), {
				callback: function () {
					that.policy.certRefs.push(resp.data.certRef)
					that.policy.certs.push(resp.data.cert)
				}
			})
		},

		// 更多选项
		changeOptionsVisible: function () {
			this.moreOptionsVisible = !this.moreOptionsVisible
		},

		// 格式化时间
		formatTime: function (timestamp) {
			return new Date(timestamp * 1000).format("Y-m-d")
		},

		// 格式化加密套件
		formatCipherSuite: function (cipherSuite) {
			return cipherSuite.replace(/(AES|3DES)/, "<var style=\"font-weight: bold\">$1</var>")
		},

		// 添加单个套件
		addCipherSuite: function (cipherSuite) {
			if (!this.policy.cipherSuites.$contains(cipherSuite)) {
				this.policy.cipherSuites.push(cipherSuite)
			}
			this.allCipherSuites.$removeValue(cipherSuite)
		},

		// 删除单个套件
		removeCipherSuite: function (cipherSuite) {
			let that = this
			teaweb.confirm("确定要删除此套件吗？", function () {
				that.policy.cipherSuites.$removeValue(cipherSuite)
				that.allCipherSuites = window.SSL_ALL_CIPHER_SUITES.$findAll(function (k, v) {
					return !that.policy.cipherSuites.$contains(v)
				})
			})
		},

		// 清除所选套件
		clearCipherSuites: function () {
			let that = this
			teaweb.confirm("确定要清除所有已选套件吗？", function () {
				that.policy.cipherSuites = []
				that.allCipherSuites = window.SSL_ALL_CIPHER_SUITES.$copy()
			})
		},

		// 批量添加套件
		addBatchCipherSuites: function (suites) {
			var that = this
			teaweb.confirm("确定要批量添加套件？", function () {
				suites.$each(function (k, v) {
					if (that.policy.cipherSuites.$contains(v)) {
						return
					}
					that.policy.cipherSuites.push(v)
				})
			})
		},

		/**
		 * 套件拖动排序
		 */
		sortableCipherSuites: function () {
			var box = document.querySelector(".cipher-suites-box")
			Sortable.create(box, {
				draggable: ".label",
				handle: ".icon.handle",
				onStart: function () {

				},
				onUpdate: function (event) {

				}
			})
		},

		// 显示所有套件
		showAllCipherSuites: function () {
			this.cipherSuitesVisible = !this.cipherSuitesVisible
		},

		// 显示HSTS更多选项
		showMoreHSTS: function () {
			this.hstsOptionsVisible = !this.hstsOptionsVisible;
			if (this.hstsOptionsVisible) {
				this.changeHSTSMaxAge()
			}
		},

		// 监控HSTS有效期修改
		changeHSTSMaxAge: function () {
			var v = this.hsts.maxAge
			if (isNaN(v)) {
				this.hsts.days = "-"
				return
			}
			this.hsts.days = parseInt(v / 86400)
			if (isNaN(this.hsts.days)) {
				this.hsts.days = "-"
			} else if (this.hsts.days < 0) {
				this.hsts.days = "-"
			}
		},

		// 设置HSTS有效期
		setHSTSMaxAge: function (maxAge) {
			this.hsts.maxAge = maxAge
			this.changeHSTSMaxAge()
		},

		// 添加HSTS域名
		addHstsDomain: function () {
			this.hstsDomainAdding = true
			this.hstsDomainEditingIndex = -1
			let that = this
			setTimeout(function () {
				that.$refs.addingHstsDomain.focus()
			}, 100)
		},

		// 修改HSTS域名
		editHstsDomain: function (index) {
			this.hstsDomainEditingIndex = index
			this.addingHstsDomain = this.hsts.domains[index]
			this.hstsDomainAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.addingHstsDomain.focus()
			}, 100)
		},

		// 确认HSTS域名添加
		confirmAddHstsDomain: function () {
			this.addingHstsDomain = this.addingHstsDomain.trim()
			if (this.addingHstsDomain.length == 0) {
				return;
			}
			if (this.hstsDomainEditingIndex > -1) {
				this.hsts.domains[this.hstsDomainEditingIndex] = this.addingHstsDomain
			} else {
				this.hsts.domains.push(this.addingHstsDomain)
			}
			this.cancelHstsDomainAdding()
		},

		// 取消HSTS域名添加
		cancelHstsDomainAdding: function () {
			this.hstsDomainAdding = false
			this.addingHstsDomain = ""
			this.hstsDomainEditingIndex = -1
		},

		// 删除HSTS域名
		removeHstsDomain: function (index) {
			this.cancelHstsDomainAdding()
			this.hsts.domains.$remove(index)
		},

		// 选择客户端CA证书
		selectClientCACert: function () {
			let that = this
			teaweb.popup("/servers/certs/selectPopup?isCA=1", {
				width: "50em",
				height: "30em",
				callback: function (resp) {
					that.policy.clientCARefs.push(resp.data.certRef)
					that.policy.clientCACerts.push(resp.data.cert)
				}
			})
		},

		// 上传CA证书
		uploadClientCACert: function () {
			let that = this
			teaweb.popup("/servers/certs/uploadPopup?isCA=1", {
				height: "28em",
				callback: function (resp) {
					teaweb.success("上传成功", function () {
						that.policy.clientCARefs.push(resp.data.certRef)
						that.policy.clientCACerts.push(resp.data.cert)
					})
				}
			})
		},

		// 删除客户端CA证书
		removeClientCACert: function (index) {
			let that = this
			teaweb.confirm("确定删除此证书吗？证书数据仍然保留，只是当前服务不再使用此证书。", function () {
				that.policy.clientCARefs.$remove(index)
				that.policy.clientCACerts.$remove(index)
			})
		}
	},
	template: `<div>
	<h4>SSL/TLS相关配置</h4>
	<input type="hidden" name="sslPolicyJSON" :value="JSON.stringify(policy)"/>
	<table class="ui table definition selectable">
		<tbody>
			<tr v-show="vProtocol == 'https'">
				<td class="title">启用HTTP/2</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="policy.http2Enabled"/>
						<label></label>
					</div>
				</td>
			</tr>
			<tr>
				<td class="title">选择证书</td>
				<td>
					<div v-if="policy.certs != null && policy.certs.length > 0">
						<div class="ui label small" v-for="(cert, index) in policy.certs">
							{{cert.name}} / {{cert.dnsNames}} / 有效至{{formatTime(cert.timeEndAt)}} &nbsp; <a href="" title="删除" @click.prevent="removeCert(index)"><i class="icon remove"></i></a>
						</div>
						<div class="ui divider"></div>
					</div>
					<div v-else>
						<span class="red">选择或上传证书后<span v-if="vProtocol == 'https'">HTTPS</span><span v-if="vProtocol == 'tls'">TLS</span>服务才能生效。</span>
						<div class="ui divider"></div>
					</div>
					<button class="ui button tiny" type="button" @click.prevent="selectCert()">选择已有证书</button> &nbsp;
					<button class="ui button tiny" type="button" @click.prevent="uploadCert()">上传新证书</button> &nbsp;
					<button class="ui button tiny" type="button" @click.prevent="requestCert()" v-if="vServerId != null && vServerId > 0">申请免费证书</button>
				</td>
			</tr>
			<tr>
				<td>TLS最低版本</td>
				<td>
					<select v-model="policy.minVersion" class="ui dropdown auto-width">
						<option v-for="version in allVersions" :value="version">{{version}}</option>
					</select>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeOptionsVisible"></more-options-tbody>
		<tbody v-show="moreOptionsVisible">
			<!-- 加密套件 -->
			<tr>
				<td>加密算法套件<em>（CipherSuites）</em></td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="policy.cipherSuitesIsOn" />
						<label>是否要自定义</label>
					</div>
					<div v-show="policy.cipherSuitesIsOn">
						<div class="ui divider"></div>
						<div class="cipher-suites-box">
							已添加套件({{policy.cipherSuites.length}})：
							<div v-for="cipherSuite in policy.cipherSuites" class="ui label tiny basic" style="margin-bottom: 0.5em">
								<input type="hidden" name="cipherSuites" :value="cipherSuite"/>
								<span v-html="formatCipherSuite(cipherSuite)"></span> &nbsp; <a href="" title="删除套件" @click.prevent="removeCipherSuite(cipherSuite)"><i class="icon remove"></i></a>
								<a href="" title="拖动改变顺序"><i class="icon bars handle"></i></a>
							</div>
						</div>
						<div>
							<div class="ui divider"></div>
							<span v-if="policy.cipherSuites.length > 0"><a href="" @click.prevent="clearCipherSuites()">[清除所有已选套件]</a> &nbsp; </span>
							<a href="" @click.prevent="addBatchCipherSuites(modernCipherSuites)">[添加推荐套件]</a> &nbsp;
							<a href="" @click.prevent="addBatchCipherSuites(intermediateCipherSuites)">[添加兼容套件]</a>
							<div class="ui divider"></div>
						</div>
		
						<div class="cipher-all-suites-box">
							<a href="" @click.prevent="showAllCipherSuites()"><span v-if="policy.cipherSuites.length == 0">所有</span>可选套件({{allCipherSuites.length}}) <i class="icon angle" :class="{down:!cipherSuitesVisible, up:cipherSuitesVisible}"></i></a>
							<a href="" v-if="cipherSuitesVisible" v-for="cipherSuite in allCipherSuites" class="ui label tiny" title="点击添加到自定义套件中" @click.prevent="addCipherSuite(cipherSuite)" v-html="formatCipherSuite(cipherSuite)" style="margin-bottom:0.5em"></a>
						</div>
						<p class="comment" v-if="cipherSuitesVisible">点击可选套件添加。</p>
					</div>
				</td>
			</tr>
			
			<!-- HSTS -->
			<tr v-show="vProtocol == 'https'">
				<td :class="{'color-border':hsts.isOn}">是否开启HSTS</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" name="hstsOn" v-model="hsts.isOn" value="1"/>
						<label></label>
					</div>
					<p class="comment">
						 开启后，会自动在响应Header中加入
						 <span class="ui label small">Strict-Transport-Security:
							 <var v-if="!hsts.isOn">...</var>
							 <var v-if="hsts.isOn"><span>max-age=</span>{{hsts.maxAge}}</var>
							 <var v-if="hsts.isOn && hsts.includeSubDomains">; includeSubDomains</var>
							 <var v-if="hsts.isOn && hsts.preload">; preload</var>
						 </span>
						  <span v-if="hsts.isOn">
							<a href="" @click.prevent="showMoreHSTS()">修改<i class="icon angle" :class="{down:!hstsOptionsVisible, up:hstsOptionsVisible}"></i> </a>
						 </span>
					</p>
				</td>
			</tr>
			<tr v-show="hsts.isOn && hstsOptionsVisible">
				<td class="color-border">HSTS包含子域名<em>（includeSubDomains）</em></td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" name="hstsIncludeSubDomains" value="1" v-model="hsts.includeSubDomains"/>
						<label></label>
					</div>
				</td>
			</tr>
			<tr v-show="hsts.isOn && hstsOptionsVisible">
				<td class="color-border">HSTS预加载<em>（preload）</em></td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" name="hstsPreload" value="1" v-model="hsts.preload"/>
						<label></label>
					</div>
				</td>
			</tr>
			<tr v-show="hsts.isOn && hstsOptionsVisible">
				<td class="color-border">HSTS生效的域名</td>
				<td colspan="2">
					<div class="names-box">
					<span class="ui label tiny basic" v-for="(domain, arrayIndex) in hsts.domains" :class="{blue:hstsDomainEditingIndex == arrayIndex}">{{domain}}
						<input type="hidden" name="hstsDomains" :value="domain"/> &nbsp;
						<a href="" @click.prevent="editHstsDomain(arrayIndex)" title="修改"><i class="icon pencil"></i></a>
						<a href="" @click.prevent="removeHstsDomain(arrayIndex)" title="删除"><i class="icon remove"></i></a>
					</span>
					</div>
					<div class="ui fields inline" v-if="hstsDomainAdding" style="margin-top:0.8em">
						<div class="ui field">
							<input type="text" name="addingHstsDomain" ref="addingHstsDomain" style="width:16em" maxlength="100" placeholder="域名，比如example.com" @keyup.enter="confirmAddHstsDomain()" @keypress.enter.prevent="1" v-model="addingHstsDomain" />
						</div>
						<div class="ui field">
							<button class="ui button tiny" type="button" @click="confirmAddHstsDomain()">确定</button>
							&nbsp; <a href="" @click.prevent="cancelHstsDomainAdding()">取消</a>
						</div>
					</div>
					<div class="ui field" style="margin-top: 1em">
						<button class="ui button tiny" type="button" @click="addHstsDomain()">+</button>
					</div>
					<p class="comment">如果没有设置域名的话，则默认支持所有的域名。</p>
				</td>
			</tr>
			
			<!-- 客户端认证 -->
			<tr>
				<td>客户端认证方式</td>
				<td>
					<select name="clientAuthType" v-model="policy.clientAuthType" class="ui dropdown auto-width">
						<option v-for="authType in allClientAuthTypes" :value="authType.type">{{authType.name}}</option>
					</select>
				</td>
			</tr>
			<tr>
				<td>客户端认证CA证书</td>
				<td>
					<div v-if="policy.clientCACerts != null && policy.clientCACerts.length > 0">
						<div class="ui label small" v-for="(cert, index) in policy.clientCACerts">
							{{cert.name}} / {{cert.dnsNames}} / 有效至{{formatTime(cert.timeEndAt)}} &nbsp; <a href="" title="删除" @click.prevent="removeClientCACert()"><i class="icon remove"></i></a>
						</div>
						<div class="ui divider"></div>
					</div>
					<button class="ui button tiny" type="button" @click.prevent="selectClientCACert()">选择已有证书</button> &nbsp;
					<button class="ui button tiny" type="button" @click.prevent="uploadClientCACert()">上传新证书</button>
					<p class="comment">用来校验客户端证书以增强安全性，通常不需要设置。</p>
				</td>
			</tr>
		</tbody>	
	</table>
	<div class="ui margin"></div>
</div>`
})

// Action列表
Vue.component("http-firewall-actions-view", {
	props: ["v-actions"],
	template: `<div>
		<div v-for="action in vActions" style="margin-bottom: 0.3em">
			<span :class="{red: action.category == 'block', orange: action.category == 'verify', green: action.category == 'allow'}">{{action.name}} ({{action.code.toUpperCase()}})</span>
		</div>             
</div>`
})

// 显示WAF规则的标签
Vue.component("http-firewall-rule-label", {
	props: ["v-rule"],
	data: function () {
		return {
			rule: this.vRule
		}
	},
	methods: {
		showErr: function (err) {

			teaweb.popupTip("规则校验错误，请修正：<span class=\"red\">"  + teaweb.encodeHTML(err) + "</span>")
		},

	},
	template: `<div>
	<div class="ui label tiny basic">
		{{rule.name}}[{{rule.param}}] 

		<!-- cc2 -->
		<span v-if="rule.param == '\${cc2}'">
			{{rule.checkpointOptions.period}}秒/{{rule.checkpointOptions.threshold}}请求
		</span>

		<!-- refererBlock -->
		<span v-if="rule.param == '\${refererBlock}'">
			{{rule.checkpointOptions.allowDomains}}
		</span>

		<span v-else>
			<span v-if="rule.paramFilters != null && rule.paramFilters.length > 0" v-for="paramFilter in rule.paramFilters"> | {{paramFilter.code}}</span> 
		<var :class="{dash:rule.isCaseInsensitive}" :title="rule.isCaseInsensitive ? '大小写不敏感':''" v-if="!rule.isComposed">{{rule.operator}}</var> 
		{{rule.value}}
		</span>
		
		<a href="" v-if="rule.err != null && rule.err.length > 0" @click.prevent="showErr(rule.err)" style="color: #db2828; opacity: 1; border-bottom: 1px #db2828 dashed; margin-left: 0.5em">规则错误</a>
	</div>
</div>`
})

// 缓存条件列表
Vue.component("http-cache-refs-box", {
	props: ["v-cache-refs"],
	data: function () {
		let refs = this.vCacheRefs
		if (refs == null) {
			refs = []
		}
		return {
			refs: refs
		}
	},
	methods: {
		timeUnitName: function (unit) {
			switch (unit) {
				case "ms":
					return "毫秒"
				case "second":
					return "秒"
				case "minute":
					return "分钟"
				case "hour":
					return "小时"
				case "day":
					return "天"
				case "week":
					return "周 "
			}
			return unit
		}
	},
	template: `<div>
	<input type="hidden" name="refsJSON" :value="JSON.stringify(refs)"/>
	
	<p class="comment" v-if="refs.length == 0">暂时还没有缓存条件。</p>
	<div v-show="refs.length > 0">
		<table class="ui table selectable celled">
			<thead>
				<tr>
					<th>缓存条件</th>
					<th class="two wide">分组关系</th>
					<th class="width10">缓存时间</th>
				</tr>
				<tr v-for="(cacheRef, index) in refs">
					<td :class="{'color-border': cacheRef.conds.connector == 'and'}" :style="{'border-left':cacheRef.isReverse ? '1px #db2828 solid' : ''}">
						<http-request-conds-view :v-conds="cacheRef.conds"></http-request-conds-view>
						<grey-label v-if="cacheRef.minSize != null && cacheRef.minSize.count > 0">
							{{cacheRef.minSize.count}}{{cacheRef.minSize.unit}}
							<span v-if="cacheRef.maxSize != null && cacheRef.maxSize.count > 0">- {{cacheRef.maxSize.count}}{{cacheRef.maxSize.unit}}</span>
						</grey-label>
						<grey-label v-else-if="cacheRef.maxSize != null && cacheRef.maxSize.count > 0">0 - {{cacheRef.maxSize.count}}{{cacheRef.maxSize.unit}}</grey-label>
						<grey-label v-if="cacheRef.status != null && cacheRef.status.length > 0 && (cacheRef.status.length > 1 || cacheRef.status[0] != 200)">状态码：{{cacheRef.status.map(function(v) {return v.toString()}).join(", ")}}</grey-label>
					</td>
					<td>
						<span v-if="cacheRef.conds.connector == 'and'">和</span>
						<span v-if="cacheRef.conds.connector == 'or'">或</span>
					</td>
					<td>
						<span v-if="!cacheRef.isReverse">{{cacheRef.life.count}} {{timeUnitName(cacheRef.life.unit)}}</span>
						<span v-else class="red">不缓存</span>
					</td>
				</tr>
			</thead>
		</table>
	</div>
	<div class="margin"></div>
</div>`
})

Vue.component("ssl-certs-box", {
	props: [
		"v-certs", // 证书列表
		"v-protocol", // 协议：https|tls
		"v-view-size", // 弹窗尺寸
		"v-single-mode" // 单证书模式
	],
	data: function () {
		let certs = this.vCerts
		if (certs == null) {
			certs = []
		}

		return {
			certs: certs
		}
	},
	methods: {
		certIds: function () {
			return this.certs.map(function (v) {
				return v.id
			})
		},
		// 删除证书
		removeCert: function (index) {
			let that = this
			teaweb.confirm("确定删除此证书吗？证书数据仍然保留，只是当前服务不再使用此证书。", function () {
				that.certs.$remove(index)
			})
		},

		// 选择证书
		selectCert: function () {
			let that = this
			let width = "50em"
			let height = "30em"
			let viewSize = this.vViewSize
			if (viewSize == null) {
				viewSize = "normal"
			}
			if (viewSize == "mini") {
				width = "35em"
				height = "20em"
			}
			teaweb.popup("/servers/certs/selectPopup?viewSize=" + viewSize, {
				width: width,
				height: height,
				callback: function (resp) {
					that.certs.push(resp.data.cert)
				}
			})
		},

		// 上传证书
		uploadCert: function () {
			let that = this
			teaweb.popup("/servers/certs/uploadPopup", {
				height: "28em",
				callback: function (resp) {
					teaweb.success("上传成功", function () {
						that.certs.push(resp.data.cert)
					})
				}
			})
		},

		// 格式化时间
		formatTime: function (timestamp) {
			return new Date(timestamp * 1000).format("Y-m-d")
		},

		// 判断是否显示选择｜上传按钮
		buttonsVisible: function () {
			return this.vSingleMode == null || !this.vSingleMode || this.certs == null || this.certs.length == 0
		}
	},
	template: `<div>
	<input type="hidden" name="certIdsJSON" :value="JSON.stringify(certIds())"/>
	<div v-if="certs != null && certs.length > 0">
		<div class="ui label small" v-for="(cert, index) in certs">
			{{cert.name}} / {{cert.dnsNames}} / 有效至{{formatTime(cert.timeEndAt)}} &nbsp; <a href="" title="删除" @click.prevent="removeCert(index)"><i class="icon remove"></i></a>
		</div>
		<div class="ui divider" v-if="buttonsVisible()"></div>
	</div>
	<div v-else>
		<span class="red">选择或上传证书后<span v-if="vProtocol == 'https'">HTTPS</span><span v-if="vProtocol == 'tls'">TLS</span>服务才能生效。</span>
		<div class="ui divider" v-if="buttonsVisible()"></div>
	</div>
	<div v-if="buttonsVisible()">
		<button class="ui button tiny" type="button" @click.prevent="selectCert()">选择已有证书</button> &nbsp;
		<button class="ui button tiny" type="button" @click.prevent="uploadCert()">上传新证书</button> &nbsp;
	</div>
</div>`
})

Vue.component("http-host-redirect-box", {
	props: ["v-redirects"],
	mounted: function () {
		let that = this
		sortTable(function (ids) {
			let newRedirects = []
			ids.forEach(function (id) {
				that.redirects.forEach(function (redirect) {
					if (redirect.id == id) {
						newRedirects.push(redirect)
					}
				})
			})
			that.updateRedirects(newRedirects)
		})
	},
	data: function () {
		let redirects = this.vRedirects
		if (redirects == null) {
			redirects = []
		}

		let id = 0
		redirects.forEach(function (v) {
			id++
			v.id = id
		})

		return {
			redirects: redirects,
			statusOptions: [
				{"code": 301, "text": "Moved Permanently"},
				{"code": 308, "text": "Permanent Redirect"},
				{"code": 302, "text": "Found"},
				{"code": 303, "text": "See Other"},
				{"code": 307, "text": "Temporary Redirect"}
			],
			id: id
		}
	},
	methods: {
		add: function () {
			let that = this
			window.UPDATING_REDIRECT = null

			teaweb.popup("/servers/server/settings/redirects/createPopup", {
				width: "50em",
				height: "30em",
				callback: function (resp) {
					that.id++
					resp.data.redirect.id = that.id
					that.redirects.push(resp.data.redirect)
					that.change()
				}
			})
		},
		update: function (index, redirect) {
			let that = this
			window.UPDATING_REDIRECT = redirect

			teaweb.popup("/servers/server/settings/redirects/createPopup", {
				width: "50em",
				height: "30em",
				callback: function (resp) {
					resp.data.redirect.id = redirect.id
					Vue.set(that.redirects, index, resp.data.redirect)
					that.change()
				}
			})
		},
		remove: function (index) {
			let that = this
			teaweb.confirm("确定要删除这条跳转规则吗？", function () {
				that.redirects.$remove(index)
				that.change()
			})
		},
		change: function () {
			let that = this
			setTimeout(function (){
				that.$emit("change", that.redirects)
			}, 100)
		},
		updateRedirects: function (newRedirects) {
			this.redirects = newRedirects
			this.change()
		}
	},
	template: `<div>
	<input type="hidden" name="hostRedirectsJSON" :value="JSON.stringify(redirects)"/>
	
	<first-menu>
		<menu-item @click.prevent="add">[创建]</menu-item>
	</first-menu>
	<div class="margin"></div>

	<p class="comment" v-if="redirects.length == 0">暂时还没有URL跳转规则。</p>
	<div v-show="redirects.length > 0">
		<table class="ui table celled selectable" id="sortable-table">
			<thead>
				<tr>
					<th style="width: 1em"></th>
					<th>跳转前URL</th>
					<th style="width: 1em"></th>
					<th>跳转后URL</th>
					<th>匹配模式</th>
					<th>HTTP状态码</th>
					<th class="two wide">状态</th>
					<th class="two op">操作</th>
				</tr>
			</thead>
			<tbody v-for="(redirect, index) in redirects" :key="redirect.id" :v-id="redirect.id">
				<tr>
					<td style="text-align: center;"><i class="icon bars handle grey"></i> </td>
					<td>
						{{redirect.beforeURL}}
						<div style="margin-top: 0.5em" v-if="redirect.conds != null && redirect.conds.groups != null && redirect.conds.groups.length > 0">
							<span class="ui label text basic tiny">匹配条件</span>
						</div>
					</td>
					<td nowrap="">-&gt;</td>
					<td>{{redirect.afterURL}}</td>
					<td>
						<span v-if="redirect.matchPrefix">匹配前缀</span>
						<span v-if="redirect.matchRegexp">正则匹配</span>
						<span v-if="!redirect.matchPrefix && !redirect.matchRegexp">精准匹配</span>
					</td>
					<td>
						<span v-if="redirect.status > 0">{{redirect.status}}</span>
						<span v-else class="disabled">默认</span>
					</td>
					<td><label-on :v-is-on="redirect.isOn"></label-on></td>
					<td>
						<a href="" @click.prevent="update(index, redirect)">修改</a> &nbsp;
						<a href="" @click.prevent="remove(index)">删除</a>	
					</td>
				</tr>
			</tbody>
		</table>
		<p class="comment" v-if="redirects.length > 1">所有规则匹配顺序为从上到下，可以拖动左侧的<i class="icon bars"></i>排序。</p>
	</div>
	<div class="margin"></div>
</div>`
})

// 单个缓存条件设置
Vue.component("http-cache-ref-box", {
	props: ["v-cache-ref", "v-is-reverse"],
	data: function () {
		let ref = this.vCacheRef
		if (ref == null) {
			ref = {
				isOn: true,
				cachePolicyId: 0,
				key: "${scheme}://${host}${requestURI}",
				life: {count: 2, unit: "hour"},
				status: [200],
				maxSize: {count: 32, unit: "mb"},
				minSize: {count: 0, unit: "kb"},
				skipCacheControlValues: ["private", "no-cache", "no-store"],
				skipSetCookie: true,
				enableRequestCachePragma: false,
				conds: null,
				allowChunkedEncoding: true,
				isReverse: this.vIsReverse
			}
		}
		if (ref.life == null) {
			ref.life = {count: 2, unit: "hour"}
		}
		if (ref.maxSize == null) {
			ref.maxSize = {count: 32, unit: "mb"}
		}
		if (ref.minSize == null) {
			ref.minSize = {count: 0, unit: "kb"}
		}
		return {
			ref: ref,
			moreOptionsVisible: false
		}
	},
	methods: {
		changeOptionsVisible: function (v) {
			this.moreOptionsVisible = v
		},
		changeLife: function (v) {
			this.ref.life = v
		},
		changeMaxSize: function (v) {
			this.ref.maxSize = v
		},
		changeMinSize: function (v) {
			this.ref.minSize = v
		},
		changeConds: function (v) {
			this.ref.conds = v
		},
		changeStatusList: function (list) {
			let result = []
			list.forEach(function (status) {
				let statusNumber = parseInt(status)
				if (isNaN(statusNumber) || statusNumber < 100 || statusNumber > 999) {
					return
				}
				result.push(statusNumber)
			})
			this.ref.status = result
		}
	},
	template: `<tbody>
	<tr>
		<td class="title">匹配条件分组 *</td>
		<td>
			<http-request-conds-box :v-conds="ref.conds" @change="changeConds"></http-request-conds-box>
			
			<input type="hidden" name="cacheRefJSON" :value="JSON.stringify(ref)"/>
		</td>
	</tr>
	<tr v-show="!vIsReverse">
		<td>缓存有效期 *</td>
		<td>
			<time-duration-box :v-value="ref.life" @change="changeLife"></time-duration-box>
		</td>
	</tr>
	<tr v-show="!vIsReverse">
		<td>缓存Key *</td>
		<td>
			<input type="text" v-model="ref.key"/>
			<p class="comment">用来区分不同缓存内容的唯一Key。</p>
		</td>
	</tr>
	<tr v-show="!vIsReverse">
		<td colspan="2"><more-options-indicator @change="changeOptionsVisible"></more-options-indicator></td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>可缓存的最大内容尺寸</td>
		<td>
			<size-capacity-box :v-value="ref.maxSize" @change="changeMaxSize"></size-capacity-box>
			<p class="comment">内容尺寸如果高于此值则不缓存。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>可缓存的最小内容尺寸</td>
		<td>
			<size-capacity-box :v-value="ref.minSize" @change="changeMinSize"></size-capacity-box>
			<p class="comment">内容尺寸如果低于此值则不缓存。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>支持分片内容</td>
		<td>
			<checkbox name="allowChunkedEncoding" value="1" v-model="ref.allowChunkedEncoding"></checkbox>
			<p class="comment">选中后，Gzip和Chunked内容可以直接缓存，无需检查内容长度。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>状态码列表</td>
		<td>
			<values-box name="statusList" size="3" maxlength="3" :values="ref.status" @change="changeStatusList"></values-box>
			<p class="comment">允许缓存的HTTP状态码列表。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>跳过的Cache-Control值</td>
		<td>
			<values-box name="skipResponseCacheControlValues" size="10" maxlength="100" :values="ref.skipCacheControlValues"></values-box>
			<p class="comment">当响应的Cache-Control为这些值时不缓存响应内容，而且不区分大小写。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>跳过Set-Cookie</td>
		<td>
			<div class="ui checkbox">
				<input type="checkbox" value="1" v-model="ref.skipSetCookie"/>
				<label></label>
			</div>
			<p class="comment">选中后，当响应的Header中有Set-Cookie时不缓存响应内容。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>支持请求no-cache刷新</td>
		<td>
			<div class="ui checkbox">
				<input type="checkbox" name="enableRequestCachePragma" value="1" v-model="ref.enableRequestCachePragma"/>
				<label></label>
			</div>
			<p class="comment">选中后，当请求的Header中含有Pragma: no-cache或Cache-Control: no-cache时，会跳过缓存直接读取源内容。</p>
		</td>
	</tr>	
</tbody>`
})

// 浏览条件列表
Vue.component("http-request-conds-view", {
	props: ["v-conds"],
	data: function () {
		let conds = this.vConds
		if (conds == null) {
			conds = {
				isOn: true,
				connector: "or",
				groups: []
			}
		}

		let that = this
		conds.groups.forEach(function (group) {
			group.conds.forEach(function (cond) {
				cond.typeName = that.typeName(cond)
			})
		})

		return {
			initConds: conds,
			version: 0 // 为了让组件能及时更新加入此变量
		}
	},
	computed: {
		// 之所以使用computed，是因为需要动态更新
		conds: function () {
			return this.initConds
		}
	},
	methods: {
		typeName: function (cond) {
			let c = window.REQUEST_COND_COMPONENTS.$find(function (k, v) {
				return v.type == cond.type
			})
			if (c != null) {
				return c.name;
			}
			return cond.param + " " + cond.operator
		},
		notifyChange: function () {
			this.version++
			let that = this
			this.initConds.groups.forEach(function (group) {
				group.conds.forEach(function (cond) {
					cond.typeName = that.typeName(cond)
				})
			})
		}
	},
	template: `<div>
		<span v-if="version < 0">{{version}}</span>
		<div v-if="conds.groups.length > 0">
			<div v-for="(group, groupIndex) in conds.groups">
				<var v-for="(cond, index) in group.conds" style="font-style: normal;display: inline-block; margin-bottom:0.5em">
					<span class="ui label small basic" style="line-height: 1.5">
						<var v-if="cond.type.length == 0 || cond.type == 'params'" style="font-style: normal">{{cond.param}} <var>{{cond.operator}}</var></var>
						<var v-if="cond.type.length > 0 && cond.type != 'params'" style="font-style: normal">{{cond.typeName}}: </var>
						{{cond.value}}
					</span>
					
					<var v-if="index < group.conds.length - 1"> {{group.connector}} &nbsp;</var>
				</var>
				<div class="ui divider" v-if="groupIndex != conds.groups.length - 1" style="margin-top:0.3em;margin-bottom:0.5em"></div>
				<div>
					<span class="ui label tiny olive" v-if="group.description != null && group.description.length > 0">{{group.description}}</span>
				</div>
			</div>	
		</div>
	</div>	
</div>`
})

Vue.component("http-firewall-config-box", {
	props: ["v-firewall-config", "v-is-location", "v-is-group", "v-firewall-policy"],
	data: function () {
		let firewall = this.vFirewallConfig
		if (firewall == null) {
			firewall = {
				isPrior: false,
				isOn: false,
				firewallPolicyId: 0
			}
		}

		return {
			firewall: firewall
		}
	},
	template: `<div>
	<input type="hidden" name="firewallJSON" :value="JSON.stringify(firewall)"/>
	<table class="ui table selectable definition">
		<prior-checkbox :v-config="firewall" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || firewall.isPrior">
			<tr v-show="!vIsGroup">
				<td>WAF策略</td>
				<td>
					<div v-if="vFirewallPolicy != null">{{vFirewallPolicy.name}} <span v-if="vFirewallPolicy.modeInfo != null">&nbsp; <span :class="{green: vFirewallPolicy.modeInfo.code == 'defend', blue: vFirewallPolicy.modeInfo.code == 'observe', grey: vFirewallPolicy.modeInfo.code == 'bypass'}">[{{vFirewallPolicy.modeInfo.name}}]</span>&nbsp;</span> <link-icon :href="'/servers/components/waf/policy?firewallPolicyId=' + vFirewallPolicy.id"></link-icon>
						<p class="comment">使用当前服务所在集群的设置。</p>
					</div>
					<span v-else class="red">当前集群没有设置WAF策略，当前配置无法生效。</span>
				</td>
			</tr>
			<tr>
				<td class="title">是否启用WAF</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="firewall.isOn"/>
						<label></label>
					</div>
					<p class="comment">启用WAF之后，各项WAF设置才会生效。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

// 指标图表
Vue.component("metric-chart", {
	props: ["v-chart", "v-stats", "v-item"],
	mounted: function () {
		this.load()
	},
	data: function () {
		let stats = this.vStats
		if (stats == null) {
			stats = []
		}
		if (stats.length > 0) {
			let sum = stats.$sum(function (k, v) {
				return v.value
			})
			if (sum < stats[0].total) {
				if (this.vChart.type == "pie") {
					stats.push({
						keys: ["其他"],
						value: stats[0].total - sum,
						total: stats[0].total,
						time: stats[0].time
					})
				}
			}
		}
		if (this.vChart.maxItems > 0) {
			stats = stats.slice(0, this.vChart.maxItems)
		} else {
			stats = stats.slice(0, 10)
		}

		stats.$rsort(function (v1, v2) {
			return v1.value - v2.value
		})

		let widthPercent = 100
		if (this.vChart.widthDiv > 0) {
			widthPercent = 100 / this.vChart.widthDiv
		}

		return {
			chart: this.vChart,
			stats: stats,
			item: this.vItem,
			width: widthPercent + "%",
			chartId: "metric-chart-" + this.vChart.id,
			valueTypeName: (this.vItem != null && this.vItem.valueTypeName != null && this.vItem.valueTypeName.length > 0) ? this.vItem.valueTypeName : ""
		}
	},
	methods: {
		load: function () {
			var el = document.getElementById(this.chartId)
			if (el == null || el.offsetWidth == 0 || el.offsetHeight == 0) {
				setTimeout(this.load, 100)
			} else {
				this.render(el)
			}
		},
		render: function (el) {
			let chart = echarts.init(el)
			window.addEventListener("resize", function () {
				chart.resize()
			})
			switch (this.chart.type) {
				case "pie":
					this.renderPie(chart)
					break
				case "bar":
					this.renderBar(chart)
					break
				case "timeBar":
					this.renderTimeBar(chart)
					break
				case "timeLine":
					this.renderTimeLine(chart)
					break
				case "table":
					this.renderTable(chart)
					break
			}
		},
		renderPie: function (chart) {
			let values = this.stats.map(function (v) {
				return {
					name: v.keys[0],
					value: v.value
				}
			})
			let that = this
			chart.setOption({
				tooltip: {
					show: true,
					trigger: "item",
					formatter: function (data) {
						let stat = that.stats[data.dataIndex]
						let percent = 0
						if (stat.total > 0) {
							percent = Math.round((stat.value * 100 / stat.total) * 100) / 100
						}
						let value = stat.value
						switch (that.item.valueType) {
							case "byte":
								value = teaweb.formatBytes(value)
								break
						}
						return stat.keys[0] + ": " + value + "，占比：" + percent + "%"
					}
				},
				series: [
					{
						name: name,
						type: "pie",
						data: values,
						areaStyle: {},
						color: ["#9DD3E8", "#B2DB9E", "#F39494", "#FBD88A", "#879BD7"]
					}
				]
			})
		},
		renderTimeBar: function (chart) {
			this.stats.$sort(function (v1, v2) {
				return (v1.time < v2.time) ? -1 : 1
			})
			let values = this.stats.map(function (v) {
				return v.value
			})

			let axis = {unit: "", divider: 1}
			switch (this.item.valueType) {
				case "count":
					axis = teaweb.countAxis(values, function (v) {
						return v
					})
					break
				case "byte":
					axis = teaweb.bytesAxis(values, function (v) {
						return v
					})
					break
			}

			let that = this
			chart.setOption({
				xAxis: {
					data: this.stats.map(function (v) {
						return that.formatTime(v.time)
					})
				},
				yAxis: {
					axisLabel: {
						formatter: function (value) {
							return value + axis.unit
						}
					}
				},
				tooltip: {
					show: true,
					trigger: "item",
					formatter: function (data) {
						let stat = that.stats[data.dataIndex]
						let value = stat.value
						switch (that.item.valueType) {
							case "byte":
								value = teaweb.formatBytes(value)
								break
						}
						return that.formatTime(stat.time) + ": " + value
					}
				},
				grid: {
					left: 50,
					top: 10,
					right: 20,
					bottom: 25
				},
				series: [
					{
						name: name,
						type: "bar",
						data: values.map(function (v) {
							return v / axis.divider
						}),
						itemStyle: {
							color: "#9DD3E8"
						},
						areaStyle: {},
						barWidth: "20em"
					}
				]
			})
		},
		renderTimeLine: function (chart) {
			this.stats.$sort(function (v1, v2) {
				return (v1.time < v2.time) ? -1 : 1
			})
			let values = this.stats.map(function (v) {
				return v.value
			})

			let axis = {unit: "", divider: 1}
			switch (this.item.valueType) {
				case "count":
					axis = teaweb.countAxis(values, function (v) {
						return v
					})
					break
				case "byte":
					axis = teaweb.bytesAxis(values, function (v) {
						return v
					})
					break
			}

			let that = this
			chart.setOption({
				xAxis: {
					data: this.stats.map(function (v) {
						return that.formatTime(v.time)
					})
				},
				yAxis: {
					axisLabel: {
						formatter: function (value) {
							return value + axis.unit
						}
					}
				},
				tooltip: {
					show: true,
					trigger: "item",
					formatter: function (data) {
						let stat = that.stats[data.dataIndex]
						let value = stat.value
						switch (that.item.valueType) {
							case "byte":
								value = teaweb.formatBytes(value)
								break
						}
						return that.formatTime(stat.time) + ": " + value
					}
				},
				grid: {
					left: 50,
					top: 10,
					right: 20,
					bottom: 25
				},
				series: [
					{
						name: name,
						type: "line",
						data: values.map(function (v) {
							return v / axis.divider
						}),
						itemStyle: {
							color: "#9DD3E8"
						},
						areaStyle: {}
					}
				]
			})
		},
		renderBar: function (chart) {
			let values = this.stats.map(function (v) {
				return v.value
			})
			let axis = {unit: "", divider: 1}
			switch (this.item.valueType) {
				case "count":
					axis = teaweb.countAxis(values, function (v) {
						return v
					})
					break
				case "byte":
					axis = teaweb.bytesAxis(values, function (v) {
						return v
					})
					break
			}
			let bottom = 24
			let rotate = 0
			let result = teaweb.xRotation(chart, this.stats.map(function (v) {
				return v.keys[0]
			}))
			if (result != null) {
				bottom = result[0]
				rotate = result[1]
			}
			let that = this
			chart.setOption({
				xAxis: {
					data: this.stats.map(function (v) {
						return v.keys[0]
					}),
					axisLabel: {
						interval: 0,
						rotate: rotate
					}
				},
				tooltip: {
					show: true,
					trigger: "item",
					formatter: function (data) {
						let stat = that.stats[data.dataIndex]
						let percent = 0
						if (stat.total > 0) {
							percent = Math.round((stat.value * 100 / stat.total) * 100) / 100
						}
						let value = stat.value
						switch (that.item.valueType) {
							case "byte":
								value = teaweb.formatBytes(value)
								break
						}
						return stat.keys[0] + ": " + value + "，占比：" + percent + "%"
					}
				},
				yAxis: {
					axisLabel: {
						formatter: function (value) {
							return value + axis.unit
						}
					}
				},
				grid: {
					left: 40,
					top: 10,
					right: 20,
					bottom: bottom
				},
				series: [
					{
						name: name,
						type: "bar",
						data: values.map(function (v) {
							return v / axis.divider
						}),
						itemStyle: {
							color: "#9DD3E8"
						},
						areaStyle: {},
						barWidth: "20em"
					}
				]
			})

			if (this.item.keys != null) {
				// IP相关操作
				if (this.item.keys.$contains("${remoteAddr}")) {
					let that = this
					chart.on("click", function (args) {
						let index = that.item.keys.$indexesOf("${remoteAddr}")[0]
						let value = that.stats[args.dataIndex].keys[index]
						teaweb.popup("/servers/ipbox?ip=" + value, {
							width: "50em",
							height: "30em"
						})
					})
				}
			}
		},
		renderTable: function (chart) {
			let table = `<table class="ui table celled">
	<thead>
		<tr>
			<th>对象</th>
			<th>数值</th>
			<th>占比</th>
		</tr>
	</thead>`
			let that = this
			this.stats.forEach(function (v) {
				let value = v.value
				switch (that.item.valueType) {
					case "byte":
						value = teaweb.formatBytes(value)
						break
				}
				table += "<tr><td>" + v.keys[0] + "</td><td>" + value + "</td>"
				let percent = 0
				if (v.total > 0) {
					percent = Math.round((v.value * 100 / v.total) * 100) / 100
				}
				table += "<td><div class=\"ui progress blue\"><div class=\"bar\" style=\"min-width: 0; height: 4px; width: " + percent + "%\"></div></div>" + percent + "%</td>"
				table += "</tr>"
			})

			table += `</table>`
			document.getElementById(this.chartId).innerHTML = table
		},
		formatTime: function (time) {
			if (time == null) {
				return ""
			}
			switch (this.item.periodUnit) {
				case "month":
					return time.substring(0, 4) + "-" + time.substring(4, 6)
				case "week":
					return time.substring(0, 4) + "-" + time.substring(4, 6)
				case "day":
					return time.substring(0, 4) + "-" + time.substring(4, 6) + "-" + time.substring(6, 8)
				case "hour":
					return time.substring(0, 4) + "-" + time.substring(4, 6) + "-" + time.substring(6, 8) + " " + time.substring(8, 10)
				case "minute":
					return time.substring(0, 4) + "-" + time.substring(4, 6) + "-" + time.substring(6, 8) + " " + time.substring(8, 10) + ":" + time.substring(10, 12)
			}
			return time
		}
	},
	template: `<div style="float: left" :style="{'width': width}">
	<h4>{{chart.name}} <span>（{{valueTypeName}}）</span></h4>
	<div class="ui divider"></div>
	<div style="height: 20em; padding-bottom: 1em; " :id="chartId" :class="{'scroll-box': chart.type == 'table'}"></div>
</div>`
})

Vue.component("metric-board", {
	template: `<div><slot></slot></div>`
})

Vue.component("http-cache-config-box", {
	props: ["v-cache-config", "v-is-location", "v-is-group", "v-cache-policy"],
	data: function () {
		let cacheConfig = this.vCacheConfig
		if (cacheConfig == null) {
			cacheConfig = {
				isPrior: false,
				isOn: false,
				addStatusHeader: true,
				cacheRefs: [],
				purgeIsOn: false,
				purgeKey: ""
			}
		}
		if (cacheConfig.cacheRefs == null) {
			cacheConfig.cacheRefs = []
		}
		return {
			cacheConfig: cacheConfig,
			moreOptionsVisible: false
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.cacheConfig.isPrior) && this.cacheConfig.isOn
		},
		generatePurgeKey: function () {
			let r = Math.random().toString() + Math.random().toString()
			let s = r.replace(/0\./g, "")
				.replace(/\./g, "")
			let result = ""
			for (let i = 0; i < s.length; i++) {
				result += String.fromCharCode(parseInt(s.substring(i, i + 1)) + ((Math.random() < 0.5) ? "a" : "A").charCodeAt(0))
			}
			this.cacheConfig.purgeKey = result
		},
		showMoreOptions: function () {
			this.moreOptionsVisible = !this.moreOptionsVisible
		}
	},
	template: `<div>
	<input type="hidden" name="cacheJSON" :value="JSON.stringify(cacheConfig)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="cacheConfig" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || cacheConfig.isPrior">
			<tr v-show="!vIsGroup">
				<td>缓存策略</td>
				<td>
					<div v-if="vCachePolicy != null">{{vCachePolicy.name}} <link-icon :href="'/servers/components/cache/policy?cachePolicyId=' + vCachePolicy.id"></link-icon>
						<p class="comment">使用当前服务所在集群的设置。</p>
					</div>
					<span v-else class="red">当前集群没有设置缓存策略，当前配置无法生效。</span>
				</td>
			</tr>
			<tr>
				<td class="title">是否开启缓存</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="cacheConfig.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td colspan="2">
					<a href="" @click.prevent="showMoreOptions"><span v-if="moreOptionsVisible">收起选项</span><span v-else>更多选项</span><i class="icon angle" :class="{up: moreOptionsVisible, down:!moreOptionsVisible}"></i></a>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn() && moreOptionsVisible">
			<tr>
				<td>自动添加X-Cache Header</td>
				<td>
					<checkbox v-model="cacheConfig.addStatusHeader"></checkbox>
					<p class="comment">选中后自动在响应Header中增加<code-label>X-Cache: BYPASS|MISS|HIT</code-label>。</p>
				</td>
			</tr>
			<tr>
				<td>允许PURGE</td>
				<td>
					<checkbox v-model="cacheConfig.purgeIsOn"></checkbox>
					<p class="comment">允许使用PURGE方法清除某个URL缓存。</p>
				</td>
			</tr>
			<tr v-show="cacheConfig.purgeIsOn">
				<td>PURGE Key *</td>
				<td>
					<input type="text" maxlength="200" v-model="cacheConfig.purgeKey"/>
					<p class="comment"><a href="" @click.prevent="generatePurgeKey">[随机生成]</a>。需要在PURGE方法调用时加入<code-label>Edge-Purge-Key: {{cacheConfig.purgeKey}}</code-label> Header。只能包含字符、数字、下划线。</p>
				</td>
			</tr>
		</tbody>
	</table>
	
	<div v-show="isOn()">
		<h4>缓存条件</h4>
		<http-cache-refs-config-box :v-cache-config="cacheConfig" :v-cache-refs="cacheConfig.cacheRefs" ></http-cache-refs-config-box>
	</div>
	<div class="margin"></div>
</div>`
})

// 通用Header长度
let defaultGeneralHeaders = ["Cache-Control", "Connection", "Date", "Pragma", "Trailer", "Transfer-Encoding", "Upgrade", "Via", "Warning"]
Vue.component("http-cond-general-header-length", {
	props: ["v-checkpoint"],
	data: function () {
		let headers = null
		let length = null

		if (window.parent.UPDATING_RULE != null) {
			let options = window.parent.UPDATING_RULE.checkpointOptions
			if (options.headers != null && Array.$isArray(options.headers)) {
				headers = options.headers
			}
			if (options.length != null) {
				length = options.length
			}
		}


		if (headers == null) {
			headers = defaultGeneralHeaders
		}

		if (length == null) {
			length = 128
		}

		let that = this
		setTimeout(function () {
			that.change()
		}, 100)

		return {
			headers: headers,
			length: length
		}
	},
	watch: {
		length: function (v) {
			let len = parseInt(v)
			if (isNaN(len)) {
				len = 0
			}
			if (len < 0) {
				len = 0
			}
			this.length = len
			this.change()
		}
	},
	methods: {
		change: function () {
			this.vCheckpoint.options = [
				{
					code: "headers",
					value: this.headers
				},
				{
					code: "length",
					value: this.length
				}
			]
		}
	},
	template: `<div>
	<table class="ui table">
		<tr>
			<td class="title">通用Header列表</td>
			<td>
				<values-box :values="headers" :placeholder="'Header'" @change="change"></values-box>
				<p class="comment">需要检查的Header列表。</p>
			</td>
		</tr>
		<tr>
			<td>Header值超出长度</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" name="" style="width: 5em" v-model="length" maxlength="6"/>
					<span class="ui label">字节</span>
				</div>
				<p class="comment">超出此长度认为匹配成功，0表示不限制。</p>
			</td>
		</tr>
	</table>
</div>`
})

// CC
Vue.component("http-firewall-checkpoint-cc", {
	props: ["v-checkpoint"],
	data: function () {
		let keys = []
		let period = 60
		let threshold = 1000

		let options = {}
		if (window.parent.UPDATING_RULE != null) {
			options = window.parent.UPDATING_RULE.checkpointOptions
		}

		if (options == null) {
			options = {}
		}
		if (options.keys != null) {
			keys = options.keys
		}
		if (keys.length == 0) {
			keys = ["${remoteAddr}", "${requestPath}"]
		}
		if (options.period != null) {
			period = options.period
		}
		if (options.threshold != null) {
			threshold = options.threshold
		}

		let that = this
		setTimeout(function () {
			that.change()
		}, 100)

		return {
			keys: keys,
			period: period,
			threshold: threshold,
			options: {},
			value: threshold
		}
	},
	watch: {
		period: function () {
			this.change()
		},
		threshold: function () {
			this.change()
		}
	},
	methods: {
		changeKeys: function (keys) {
			this.keys = keys
			this.change()
		},
		change: function () {
			let period = parseInt(this.period.toString())
			if (isNaN(period) || period <= 0) {
				period = 60
			}

			let threshold = parseInt(this.threshold.toString())
			if (isNaN(threshold) || threshold <= 0) {
				threshold = 1000
			}
			this.value = threshold

			this.vCheckpoint.options = [
				{
					code: "keys",
					value: this.keys
				},
				{
					code: "period",
					value: period,
				},
				{
					code: "threshold",
					value: threshold
				}
			]
		}
	},
	template: `<div>
	<input type="hidden" name="operator" value="gt"/>
	<input type="hidden" name="value" :value="value"/>
	<table class="ui table">
		<tr>
			<td class="title">统计对象组合 *</td>
			<td>
				<metric-keys-config-box :v-keys="keys" @change="changeKeys"></metric-keys-config-box>
			</td>
		</tr>
		<tr>
			<td>统计周期 *</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" v-model="period" style="width: 6em" maxlength="8"/>
					<span class="ui label">秒</span>
				</div>
			</td>
		</tr>
		<tr>
			<td>阈值 *</td>
			<td>
				<input type="text" v-model="threshold" style="width: 6em" maxlength="8"/>
			</td>
		</tr>
	</table>
</div>`
})

// 防盗链
Vue.component("http-firewall-checkpoint-referer-block", {
	props: ["v-checkpoint"],
	data: function () {
		let allowEmpty = true
		let allowSameDomain = true
		let allowDomains = []

		let options = {}
		if (window.parent.UPDATING_RULE != null) {
			options = window.parent.UPDATING_RULE.checkpointOptions
		}

		if (options == null) {
			options = {}
		}
		if (typeof (options.allowEmpty) == "boolean") {
			allowEmpty = options.allowEmpty
		}
		if (typeof (options.allowSameDomain) == "boolean") {
			allowSameDomain = options.allowSameDomain
		}
		if (options.allowDomains != null && typeof (options.allowDomains) == "object") {
			allowDomains = options.allowDomains
		}

		let that = this
		setTimeout(function () {
			that.change()
		}, 100)

		return {
			allowEmpty: allowEmpty,
			allowSameDomain: allowSameDomain,
			allowDomains: allowDomains,
			options: {},
			value: 0
		}
	},
	watch: {
		allowEmpty: function () {
			this.change()
		},
		allowSameDomain: function () {
			this.change()
		}
	},
	methods: {
		changeAllowDomains: function (values) {
			this.allowDomains = values
			this.change()
		},
		change: function () {
			this.vCheckpoint.options = [
				{
					code: "allowEmpty",
					value: this.allowEmpty
				},
				{
					code: "allowSameDomain",
					value: this.allowSameDomain,
				},
				{
					code: "allowDomains",
					value: this.allowDomains
				}
			]
		}
	},
	template: `<div>
	<input type="hidden" name="operator" value="eq"/>
	<input type="hidden" name="value" :value="value"/>
	<table class="ui table">
		<tr>
			<td class="title">来源域名允许为空</td>
			<td>
				<checkbox v-model="allowEmpty"></checkbox>
				<p class="comment">允许不带来源的访问。</p>
			</td>
		</tr>
		<tr>
			<td>来源域名允许一致</td>
			<td>
				<checkbox v-model="allowSameDomain"></checkbox>
				<p class="comment">允许来源域名和当前访问的域名一致，相当于在站内访问。</p>
			</td>
		</tr>
		<tr>
			<td>允许的来源域名</td>
			<td>
				<values-box :values="allowDomains" @change="changeAllowDomains"></values-box>
				<p class="comment">允许的来源域名列表，比如<code-label>example.com</code-label>、<code-label>*.example.com</code-label>。单个星号<code-label>*</code-label>表示允许所有域名。</p>
			</td>
		</tr>
	</table>
</div>`
})

Vue.component("http-cache-refs-config-box", {
	props: ["v-cache-refs", "v-cache-config", "v-cache-policy-id"],
	mounted: function () {
		let that = this
		sortTable(function (ids) {
			let newRefs = []
			ids.forEach(function (id) {
				that.refs.forEach(function (ref) {
					if (ref.id == id) {
						newRefs.push(ref)
					}
				})
			})
			that.updateRefs(newRefs)
			that.change()
		})
	},
	data: function () {
		let refs = this.vCacheRefs
		if (refs == null) {
			refs = []
		}

		let id = 0
		refs.forEach(function (ref) {
			id++
			ref.id = id
		})
		return {
			refs: refs,
			id: id // 用来对条件进行排序
		}
	},
	methods: {
		addRef: function (isReverse) {
			window.UPDATING_CACHE_REF = null

			let width = window.innerWidth
			if (width > 1024) {
				width = 1024
			}
			let height = window.innerHeight
			if (height > 500) {
				height = 500
			}
			let that = this
			teaweb.popup("/servers/server/settings/cache/createPopup?isReverse=" + (isReverse ? 1 : 0), {
				width: width + "px",
				height: height + "px",
				callback: function (resp) {
					let newRef = resp.data.cacheRef
					if (newRef.conds == null) {
						return
					}

					that.id++
					newRef.id = that.id

					if (newRef.isReverse) {
						let newRefs = []
						let isAdded = false
						that.refs.forEach(function (v) {
							if (!v.isReverse && !isAdded) {
								newRefs.push(newRef)
								isAdded = true
							}
							newRefs.push(v)
						})
						if (!isAdded) {
							newRefs.push(newRef)
						}

						that.updateRefs(newRefs)
					} else {
						that.refs.push(newRef)
					}

					that.change()
				}
			})
		},
		updateRef: function (index, cacheRef) {
			window.UPDATING_CACHE_REF = cacheRef

			let width = window.innerWidth
			if (width > 1024) {
				width = 1024
			}
			let height = window.innerHeight
			if (height > 500) {
				height = 500
			}
			let that = this
			teaweb.popup("/servers/server/settings/cache/createPopup", {
				width: width + "px",
				height: height + "px",
				callback: function (resp) {
					resp.data.cacheRef.id = that.refs[index].id
					Vue.set(that.refs, index, resp.data.cacheRef)

					// 通知子组件更新
					that.$refs.cacheRef[index].notifyChange()

					that.change()
				}
			})
		},
		removeRef: function (index) {
			let that = this
			teaweb.confirm("确定要删除此缓存设置吗？", function () {
				that.refs.$remove(index)
				that.change()
			})
		},
		updateRefs: function (newRefs) {
			this.refs = newRefs
			if (this.vCacheConfig != null) {
				this.vCacheConfig.cacheRefs = newRefs
			}
		},
		timeUnitName: function (unit) {
			switch (unit) {
				case "ms":
					return "毫秒"
				case "second":
					return "秒"
				case "minute":
					return "分钟"
				case "hour":
					return "小时"
				case "day":
					return "天"
				case "week":
					return "周 "
			}
			return unit
		},
		change: function () {
			// 自动保存
			if (this.vCachePolicyId != null && this.vCachePolicyId > 0) {
				Tea.action("/servers/components/cache/updateRefs")
					.params({
						cachePolicyId: this.vCachePolicyId,
						refsJSON: JSON.stringify(this.refs)
					})
					.post()
			}
		}
	},
	template: `<div>
	<input type="hidden" name="refsJSON" :value="JSON.stringify(refs)"/>
	
	<div>
		<p class="comment" v-if="refs.length == 0">暂时还没有缓存条件。</p>
		<table class="ui table selectable celled" v-show="refs.length > 0" id="sortable-table">
			<thead>
				<tr>
					<th style="width:1em"></th>
					<th>缓存条件</th>
					<th class="two wide">分组关系</th>
					<th class="width10">缓存时间</th>
					<th class="two op">操作</th>
				</tr>
			</thead>	
			<tbody v-for="(cacheRef, index) in refs" :key="cacheRef.id" :v-id="cacheRef.id">
				<tr>
					<td style="text-align: center;"><i class="icon bars handle grey"></i> </td>
					<td :class="{'color-border': cacheRef.conds.connector == 'and'}" :style="{'border-left':cacheRef.isReverse ? '1px #db2828 solid' : ''}">
						<http-request-conds-view :v-conds="cacheRef.conds" ref="cacheRef"></http-request-conds-view>
						<grey-label v-if="cacheRef.minSize != null && cacheRef.minSize.count > 0">
							{{cacheRef.minSize.count}}{{cacheRef.minSize.unit}}
							<span v-if="cacheRef.maxSize != null && cacheRef.maxSize.count > 0">- {{cacheRef.maxSize.count}}{{cacheRef.maxSize.unit}}</span>
						</grey-label>
						<grey-label v-else-if="cacheRef.maxSize != null && cacheRef.maxSize.count > 0">0 - {{cacheRef.maxSize.count}}{{cacheRef.maxSize.unit}}</grey-label>
						<grey-label v-if="cacheRef.status != null && cacheRef.status.length > 0 && (cacheRef.status.length > 1 || cacheRef.status[0] != 200)">状态码：{{cacheRef.status.map(function(v) {return v.toString()}).join(", ")}}</grey-label>
					</td>
					<td>
						<span v-if="cacheRef.conds.connector == 'and'">和</span>
						<span v-if="cacheRef.conds.connector == 'or'">或</span>
					</td>
					<td>
						<span v-if="!cacheRef.isReverse">{{cacheRef.life.count}} {{timeUnitName(cacheRef.life.unit)}}</span>
						<span v-else class="red">不缓存</span>
					</td>
					<td>
						<a href="" @click.prevent="updateRef(index, cacheRef)">修改</a> &nbsp;
						<a href="" @click.prevent="removeRef(index)">删除</a>
					</td>
				</tr>
			</tbody>
		</table>
		<p class="comment" v-if="refs.length > 1">所有条件匹配顺序为从上到下，可以拖动左侧的<i class="icon bars"></i>排序。</p>
		
		<button class="ui button tiny" @click.prevent="addRef(false)">+添加缓存设置</button> &nbsp; &nbsp; <a href="" @click.prevent="addRef(true)">+添加不缓存设置</a>
	</div>
	<div class="margin"></div>
</div>`
})

Vue.component("origin-list-box", {
	props: ["v-primary-origins", "v-backup-origins", "v-server-type", "v-params"],
	data: function () {
		return {
			primaryOrigins: this.vPrimaryOrigins,
			backupOrigins: this.vBackupOrigins
		}
	},
	methods: {
		createPrimaryOrigin: function () {
			teaweb.popup("/servers/server/settings/origins/addPopup?originType=primary&" + this.vParams, {
				height: "27em",
				callback: function (resp) {
					teaweb.success("保存成功", function () {
						window.location.reload()
					})
				}
			})
		},
		createBackupOrigin: function () {
			teaweb.popup("/servers/server/settings/origins/addPopup?originType=backup&" + this.vParams, {
				height: "27em",
				callback: function (resp) {
					teaweb.success("保存成功", function () {
						window.location.reload()
					})
				}
			})
		},
		updateOrigin: function (originId, originType) {
			teaweb.popup("/servers/server/settings/origins/updatePopup?originType=" + originType + "&" + this.vParams + "&originId=" + originId, {
				height: "27em",
				callback: function (resp) {
					teaweb.success("保存成功", function () {
						window.location.reload()
					})
				}
			})
		},
		deleteOrigin: function (originId, originType) {
			let that = this
			teaweb.confirm("确定要删除此源站吗？", function () {
				Tea.action("/servers/server/settings/origins/delete?" + that.vParams + "&originId=" + originId + "&originType=" + originType)
					.post()
					.success(function () {
						teaweb.success("保存成功", function () {
							window.location.reload()
						})
					})
			})
		}
	},
	template: `<div>
	<h3>主要源站 <a href="" @click.prevent="createPrimaryOrigin()">[添加主要源站]</a> </h3>
	<p class="comment" v-if="primaryOrigins.length == 0">暂时还没有主要源站。</p>
	<origin-list-table v-if="primaryOrigins.length > 0" :v-origins="vPrimaryOrigins" :v-origin-type="'primary'" @deleteOrigin="deleteOrigin" @updateOrigin="updateOrigin"></origin-list-table>

	<h3>备用源站 <a href="" @click.prevent="createBackupOrigin()">[添加备用源站]</a></h3>
	<p class="comment" v-if="backupOrigins.length == 0" :v-origins="primaryOrigins">暂时还没有备用源站。</p>
	<origin-list-table v-if="backupOrigins.length > 0" :v-origins="backupOrigins" :v-origin-type="'backup'" @deleteOrigin="deleteOrigin" @updateOrigin="updateOrigin"></origin-list-table>
</div>`
})

Vue.component("origin-list-table", {
	props: ["v-origins", "v-origin-type"],
	data: function () {
		return {}
	},
	methods: {
		deleteOrigin: function (originId) {
			this.$emit("deleteOrigin", originId, this.vOriginType)
		},
		updateOrigin: function (originId) {
			this.$emit("updateOrigin", originId, this.vOriginType)
		}
	},
	template: `
<table class="ui table selectable">
	<thead>
		<tr>
			<th>源站地址</th>
			<th>权重</th>
			<th class="width10">状态</th>
			<th class="two op">操作</th>
		</tr>	
	</thead>
	<tr v-for="origin in vOrigins">
		<td :class="{disabled:!origin.isOn}"><a href="" @click.prevent="updateOrigin(origin.id)">{{origin.addr}} &nbsp;<i class="icon clone outline small"></i></a>
			<div v-if="origin.name.length > 0" style="margin-top: 0.5em">
				<tiny-basic-label>{{origin.name}}</tiny-basic-label>
			</div>
			<div v-if="origin.domains != null && origin.domains.length > 0">
				<grey-label v-for="domain in origin.domains">{{domain}}</grey-label>
			</div>
		</td>
		<td :class="{disabled:!origin.isOn}">{{origin.weight}}</td>
		<td>
			<label-on :v-is-on="origin.isOn"></label-on>
		</td>
		<td>
			<a href="" @click.prevent="updateOrigin(origin.id)">修改</a> &nbsp;
			<a href="" @click.prevent="deleteOrigin(origin.id)">删除</a>
		</td>
	</tr>
</table>`
})

Vue.component("http-firewall-policy-selector", {
	props: ["v-http-firewall-policy"],
	mounted: function () {
		let that = this
		Tea.action("/servers/components/waf/count")
			.post()
			.success(function (resp) {
				that.count = resp.data.count
			})
	},
	data: function () {
		let firewallPolicy = this.vHttpFirewallPolicy
		return {
			count: 0,
			firewallPolicy: firewallPolicy
		}
	},
	methods: {
		remove: function () {
			this.firewallPolicy = null
		},
		select: function () {
			let that = this
			teaweb.popup("/servers/components/waf/selectPopup", {
				callback: function (resp) {
					that.firewallPolicy = resp.data.firewallPolicy
				}
			})
		},
		create: function () {
			let that = this
			teaweb.popup("/servers/components/waf/createPopup", {
				height: "26em",
				callback: function (resp) {
					that.firewallPolicy = resp.data.firewallPolicy
				}
			})
		}
	},
	template: `<div>
	<div v-if="firewallPolicy != null" class="ui label basic">
		<input type="hidden" name="httpFirewallPolicyId" :value="firewallPolicy.id"/>
		{{firewallPolicy.name}} &nbsp; <a :href="'/servers/components/waf/policy?firewallPolicyId=' + firewallPolicy.id" target="_blank" title="修改"><i class="icon pencil small"></i></a>&nbsp; <a href="" @click.prevent="remove()" title="删除"><i class="icon remove small"></i></a>
	</div>
	<div v-if="firewallPolicy == null">
		<span v-if="count > 0"><a href="" @click.prevent="select">[选择已有策略]</a> &nbsp; &nbsp; </span><a href="" @click.prevent="create">[创建新策略]</a>
	</div>
</div>`
})

Vue.component("http-websocket-box", {
	props: ["v-websocket-ref", "v-websocket-config", "v-is-location", "v-is-group"],
	data: function () {
		let websocketRef = this.vWebsocketRef
		if (websocketRef == null) {
			websocketRef = {
				isPrior: false,
				isOn: false,
				websocketId: 0
			}
		}

		let websocketConfig = this.vWebsocketConfig
		if (websocketConfig == null) {
			websocketConfig = {
				id: 0,
				isOn: false,
				handshakeTimeout: {
					count: 30,
					unit: "second"
				},
				allowAllOrigins: true,
				allowedOrigins: [],
				requestSameOrigin: true,
				requestOrigin: ""
			}
		} else {
			if (websocketConfig.handshakeTimeout == null) {
				websocketConfig.handshakeTimeout = {
					count: 30,
					unit: "second",
				}
			}
			if (websocketConfig.allowedOrigins == null) {
				websocketConfig.allowedOrigins = []
			}
		}

		return {
			websocketRef: websocketRef,
			websocketConfig: websocketConfig,
			handshakeTimeoutCountString: websocketConfig.handshakeTimeout.count.toString(),
			advancedVisible: false
		}
	},
	watch: {
		handshakeTimeoutCountString: function (v) {
			let count = parseInt(v)
			if (!isNaN(count) && count >= 0) {
				this.websocketConfig.handshakeTimeout.count = count
			} else {
				this.websocketConfig.handshakeTimeout.count = 0
			}
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.websocketRef.isPrior) && this.websocketRef.isOn
		},
		changeAdvancedVisible: function (v) {
			this.advancedVisible = v
		},
		createOrigin: function () {
			let that = this
			teaweb.popup("/servers/server/settings/websocket/createOrigin", {
				height: "12.5em",
				callback: function (resp) {
					that.websocketConfig.allowedOrigins.push(resp.data.origin)
				}
			})
		},
		removeOrigin: function (index) {
			this.websocketConfig.allowedOrigins.$remove(index)
		}
	},
	template: `<div>
	<input type="hidden" name="websocketRefJSON" :value="JSON.stringify(websocketRef)"/>
	<input type="hidden" name="websocketJSON" :value="JSON.stringify(websocketConfig)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="websocketRef" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="((!vIsLocation && !vIsGroup) || websocketRef.isPrior)">
			<tr>
				<td class="title">是否启用配置</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="websocketRef.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td class="color-border">允许所有来源域<em>（Origin）</em></td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="websocketConfig.allowAllOrigins"/>
						<label></label>
					</div>
					<p class="comment">选中表示允许所有的来源域。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn() && !websocketConfig.allowAllOrigins">
			<tr>
				<td class="color-border">允许的来源域列表<em>（Origin）</em></td>
				<td>
					<div v-if="websocketConfig.allowedOrigins.length > 0">
						<div class="ui label tiny" v-for="(origin, index) in websocketConfig.allowedOrigins">
							{{origin}} <a href="" title="删除" @click.prevent="removeOrigin(index)"><i class="icon remove"></i></a>
						</div>
						<div class="ui divider"></div>
					</div>
					<button class="ui button tiny" type="button" @click.prevent="createOrigin()">+</button>
					<p class="comment">只允许在列表中的来源域名访问Websocket服务。</p>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-show="isOn()"></more-options-tbody>
		<tbody v-show="isOn() && advancedVisible">
			<tr>
				<td class="color-border">是否传递请求来源域</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="websocketConfig.requestSameOrigin"/>
						<label></label>
					</div>
					<p class="comment">选中表示把接收到的请求中的<span class="ui label tiny">Origin</span>字段传递到源站。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn() && advancedVisible && !websocketConfig.requestSameOrigin">
			<tr>
				<td class="color-border">指定传递的来源域</td>
				<td>
					<input type="text" v-model="websocketConfig.requestOrigin" maxlength="200"/>
					<p class="comment">指定向源站传递的<span class="ui label tiny">Origin</span>字段值。</p>
				</td>
			</tr>
		</tbody>
		<!-- TODO 这个选项暂时保留 -->
		<tbody v-show="isOn() && false">
			<tr>
				<td>握手超时时间<em>（Handshake）</em></td>
				<td>
					<div class="ui fields inline">
						<div class="ui field">
							<input type="text" maxlength="10" v-model="handshakeTimeoutCountString" style="width:6em"/>
						</div>
						<div class="ui field">
							秒
						</div>
					</div>
					<p class="comment">0表示使用默认的时间设置。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("http-rewrite-rule-list", {
	props: ["v-web-id", "v-rewrite-rules"],
	mounted: function () {
		setTimeout(this.sort, 1000)
	},
	data: function () {
		let rewriteRules = this.vRewriteRules
		if (rewriteRules == null) {
			rewriteRules = []
		}
		return {
			rewriteRules: rewriteRules
		}
	},
	methods: {
		updateRewriteRule: function (rewriteRuleId) {
			teaweb.popup("/servers/server/settings/rewrite/updatePopup?webId=" + this.vWebId + "&rewriteRuleId=" + rewriteRuleId, {
				height: "26em",
				callback: function () {
					window.location.reload()
				}
			})
		},
		deleteRewriteRule: function (rewriteRuleId) {
			let that = this
			teaweb.confirm("确定要删除此重写规则吗？", function () {
				Tea.action("/servers/server/settings/rewrite/delete")
					.params({
						webId: that.vWebId,
						rewriteRuleId: rewriteRuleId
					})
					.post()
					.refresh()
			})
		},
		// 排序
		sort: function () {
			if (this.rewriteRules.length == 0) {
				return
			}
			let that = this
			sortTable(function (rowIds) {
				Tea.action("/servers/server/settings/rewrite/sort")
					.post()
					.params({
						webId: that.vWebId,
						rewriteRuleIds: rowIds
					})
					.success(function () {
						teaweb.success("保存成功")
					})
			})
		}
	},
	template: `<div>
	<div class="margin"></div>
	<p class="comment" v-if="rewriteRules.length == 0">暂时还没有重写规则。</p>
	<table class="ui table selectable" v-if="rewriteRules.length > 0" id="sortable-table">
		<thead>
			<tr>
				<th style="width:1em"></th>
				<th>匹配规则</th>
				<th>转发目标</th>
				<th>转发方式</th>
				<th class="two wide">状态</th>
				<th class="two op">操作</th>
			</tr>
		</thead>
		<tbody v-for="rule in rewriteRules" :v-id="rule.id">
			<tr>
				<td><i class="icon bars grey handle"></i></td>
				<td>{{rule.pattern}}
				<br/>
					<http-rewrite-labels-label class="ui label tiny" v-if="rule.isBreak">BREAK</http-rewrite-labels-label>
					<http-rewrite-labels-label class="ui label tiny" v-if="rule.mode == 'redirect' && rule.redirectStatus != 307">{{rule.redirectStatus}}</http-rewrite-labels-label>
					<http-rewrite-labels-label class="ui label tiny" v-if="rule.proxyHost.length > 0">Host: {{rule.proxyHost}}</http-rewrite-labels-label>
				</td>
				<td>{{rule.replace}}</td>
				<td>
					<span v-if="rule.mode == 'proxy'">隐式</span>
					<span v-if="rule.mode == 'redirect'">显示</span>
				</td>
				<td>
					<label-on :v-is-on="rule.isOn"></label-on>
				</td>
				<td>
					<a href="" @click.prevent="updateRewriteRule(rule.id)">修改</a> &nbsp;
					<a href="" @click.prevent="deleteRewriteRule(rule.id)">删除</a>
				</td>
			</tr>
		</tbody>
	</table>
	<p class="comment" v-if="rewriteRules.length > 0">拖动左侧的<i class="icon bars grey"></i>图标可以对重写规则进行排序。</p>

</div>`
})

Vue.component("http-rewrite-labels-label", {
	props: ["v-class"],
	template: `<span class="ui label tiny" :class="vClass" style="font-size:0.7em;padding:4px;margin-top:0.3em;margin-bottom:0.3em"><slot></slot></span>`
})

Vue.component("server-name-box", {
    props: ["v-server-names"],
    data: function () {
        let serverNames = this.vServerNames;
        if (serverNames == null) {
            serverNames = []
        }
        return {
            serverNames: serverNames,
            isSearching: false,
            keyword: ""
        }
    },
    methods: {
        addServerName: function () {
            window.UPDATING_SERVER_NAME = null
            let that = this
            teaweb.popup("/servers/addServerNamePopup", {
                callback: function (resp) {
                    var serverName = resp.data.serverName
                    that.serverNames.push(serverName)
                }
            });
        },

        removeServerName: function (index) {
            this.serverNames.$remove(index)
        },

        updateServerName: function (index, serverName) {
            window.UPDATING_SERVER_NAME = serverName
            let that = this
            teaweb.popup("/servers/addServerNamePopup", {
                callback: function (resp) {
                    var serverName = resp.data.serverName
                    Vue.set(that.serverNames, index, serverName)
                }
            });
        },
        showSearchBox: function () {
            this.isSearching = !this.isSearching
            if (this.isSearching) {
                let that = this
                setTimeout(function () {
                    that.$refs.keywordRef.focus()
                }, 200)
            } else {
                this.keyword = ""
            }
        },
    },
    watch: {
        keyword: function (v) {
            this.serverNames.forEach(function (serverName) {
                if (v.length == 0) {
                    serverName.isShowing = true
                    return
                }
                if (serverName.subNames == null || serverName.subNames.length == 0) {
                    if (!teaweb.match(serverName.name, v)) {
                        serverName.isShowing = false
                    }
                } else {
                    let found = false
                    serverName.subNames.forEach(function (subName) {
                        if (teaweb.match(subName, v)) {
                            found = true
                        }
                    })
                    serverName.isShowing = found
                }
            })
        }
    },
    template: `<div>
	<input type="hidden" name="serverNames" :value="JSON.stringify(serverNames)"/>
	<div v-if="serverNames.length > 0">
		<div v-for="(serverName, index) in serverNames" class="ui label small basic">
			<em v-if="serverName.type != 'full'">{{serverName.type}}</em>  
			<span v-if="serverName.subNames == null || serverName.subNames.length == 0" :class="{disabled: serverName.isShowing === false}">{{serverName.name}}</span>
			<span v-else :class="{disabled: serverName.isShowing === false}">{{serverName.subNames[0]}}等{{serverName.subNames.length}}个域名</span>
			<a href="" title="修改" @click.prevent="updateServerName(index, serverName)"><i class="icon pencil small"></i></a> <a href="" title="删除" @click.prevent="removeServerName(index)"><i class="icon remove"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div class="ui fields inline">
	    <div class="ui field"><a href="" @click.prevent="addServerName()">[添加域名绑定]</a></div>
	    <div class="ui field" v-if="serverNames.length > 0"><span class="grey">|</span> </div>
	    <div class="ui field" v-if="serverNames.length > 0">
	        <a href="" @click.prevent="showSearchBox()" v-if="!isSearching"><i class="icon search small"></i></a>
	        <a href="" @click.prevent="showSearchBox()" v-if="isSearching"><i class="icon close small"></i></a>
        </div>
        <div class="ui field" v-if="isSearching">
            <input type="text" placeholder="搜索域名" ref="keywordRef" class="ui input tiny" v-model="keyword"/>
        </div>
    </div>
</div>`
})

// 域名列表
Vue.component("domains-box", {
	props: ["v-domains"],
	data: function () {
		let domains = this.vDomains
		if (domains == null) {
			domains = []
		}
		return {
			domains: domains,
			isAdding: false,
			addingDomain: ""
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.addingDomain.focus()
			}, 100)
		},
		confirm: function () {
			let that = this

			// 删除其中的空格
			this.addingDomain = this.addingDomain.replace(/\s/g, "")

			if (this.addingDomain.length == 0) {
				teaweb.warn("请输入要添加的域名", function () {
					that.$refs.addingDomain.focus()
				})
				return
			}


			// 基本校验
			if (this.addingDomain[0] == "~") {
				let expr = this.addingDomain.substring(1)
				try {
					new RegExp(expr)
				} catch (e) {
					teaweb.warn("正则表达式错误：" + e.message, function () {
						that.$refs.addingDomain.focus()
					})
					return
				}
			}

			this.domains.push(this.addingDomain)
			this.cancel()
		},
		remove: function (index) {
			this.domains.$remove(index)
		},
		cancel: function () {
			this.isAdding = false
			this.addingDomain = ""
		}
	},
	template: `<div>
	<input type="hidden" name="domainsJSON" :value="JSON.stringify(domains)"/>
	<div v-if="domains.length > 0">
		<span class="ui label small basic" v-for="(domain, index) in domains">
			<span v-if="domain.length > 0 && domain[0] == '~'" class="grey" style="font-style: normal">[正则]</span>
			<span v-if="domain.length > 0 && domain[0] == '.'" class="grey" style="font-style: normal">[后缀]</span>
			<span v-if="domain.length > 0 && domain[0] == '*'" class="grey" style="font-style: normal">[泛域名]</span>
			{{domain}}
			&nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</span>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding">
		<div class="ui fields">
			<div class="ui field">
				<input type="text" v-model="addingDomain" @keyup.enter="confirm()" @keypress.enter.prevent="1" ref="addingDomain" placeholder="*.xxx.com" size="30"/>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>
				&nbsp; <a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
		<p class="comment">支持普通域名（<code-label>example.com</code-label>）、泛域名（<code-label>*.example.com</code-label>）、域名后缀（以点号开头，如<code-label>.example.com</code-label>）和正则表达式（以波浪号开头，如<code-label>~.*.example.com</code-label>）。</p>
		<div class="ui divider"></div>
	</div>
	<div style="margin-top: 0.5em" v-if="!isAdding">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})

Vue.component("http-redirect-to-https-box", {
	props: ["v-redirect-to-https-config", "v-is-location"],
	data: function () {
		let redirectToHttpsConfig = this.vRedirectToHttpsConfig
		if (redirectToHttpsConfig == null) {
			redirectToHttpsConfig = {
				isPrior: false,
				isOn: false,
				host: "",
				port: 0,
				status: 0,
				onlyDomains: [],
				exceptDomains: []
			}
		} else {
			if (redirectToHttpsConfig.onlyDomains == null) {
				redirectToHttpsConfig.onlyDomains = []
			}
			if (redirectToHttpsConfig.exceptDomains == null) {
				redirectToHttpsConfig.exceptDomains = []
			}
		}
		return {
			redirectToHttpsConfig: redirectToHttpsConfig,
			portString: (redirectToHttpsConfig.port > 0) ? redirectToHttpsConfig.port.toString() : "",
			moreOptionsVisible: false,
			statusOptions: [
				{"code": 301, "text": "Moved Permanently"},
				{"code": 308, "text": "Permanent Redirect"},
				{"code": 302, "text": "Found"},
				{"code": 303, "text": "See Other"},
				{"code": 307, "text": "Temporary Redirect"}
			]
		}
	},
	watch: {
		"redirectToHttpsConfig.status": function () {
			this.redirectToHttpsConfig.status = parseInt(this.redirectToHttpsConfig.status)
		},
		portString: function (v) {
			let port = parseInt(v)
			if (!isNaN(port)) {
				this.redirectToHttpsConfig.port = port
			} else {
				this.redirectToHttpsConfig.port = 0
			}
		}
	},
	methods: {
		changeMoreOptions: function (isVisible) {
			this.moreOptionsVisible = isVisible
		},
		changeOnlyDomains: function (values) {
			this.redirectToHttpsConfig.onlyDomains = values
			this.$forceUpdate()
		},
		changeExceptDomains: function (values) {
			this.redirectToHttpsConfig.exceptDomains = values
			this.$forceUpdate()
		}
	},
	template: `<div>
	<input type="hidden" name="redirectToHTTPSJSON" :value="JSON.stringify(redirectToHttpsConfig)"/>
	
	<!-- Location -->
	<table class="ui table selectable definition" v-if="vIsLocation">
		<prior-checkbox :v-config="redirectToHttpsConfig"></prior-checkbox>
		<tbody v-show="redirectToHttpsConfig.isPrior">
			<tr>
				<td class="title">自动跳转到HTTPS</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="redirectToHttpsConfig.isOn"/>
						<label></label>
					</div>
					<p class="comment">开启后，所有HTTP的请求都会自动跳转到对应的HTTPS URL上，<more-options-angle @change="changeMoreOptions"></more-options-angle></p>
					
					<!--  TODO 如果已经设置了特殊设置，需要在界面上显示 -->
					<table class="ui table" v-show="moreOptionsVisible">
						<tr>
							<td class="title">状态码</td>
							<td>
								<select class="ui dropdown auto-width" v-model="redirectToHttpsConfig.status">
									<option value="0">[使用默认]</option>
									<option v-for="option in statusOptions" :value="option.code">{{option.code}} {{option.text}}</option>
								</select>
							</td>
						</tr>
						<tr>
							<td>域名或IP地址</td>
							<td>
								<input type="text" name="host" v-model="redirectToHttpsConfig.host"/>
								<p class="comment">默认和用户正在访问的域名或IP地址一致。</p>
							</td>
						</tr>
						<tr>
							<td>端口</td>
							<td>
								<input type="text" name="port" v-model="portString" maxlength="5" style="width:6em"/>
								<p class="comment">默认端口为443。</p>
							</td>
						</tr>
					</table>
				</td>
			</tr>	
		</tbody>
	</table>
	
	<!-- 非Location -->
	<div v-if="!vIsLocation">
		<div class="ui checkbox">
			<input type="checkbox" v-model="redirectToHttpsConfig.isOn"/>
			<label></label>
		</div>
		<p class="comment">开启后，所有HTTP的请求都会自动跳转到对应的HTTPS URL上，<more-options-angle @change="changeMoreOptions"></more-options-angle></p>
		
		<!--  TODO 如果已经设置了特殊设置，需要在界面上显示 -->
		<table class="ui table" v-show="moreOptionsVisible">
			<tr>
				<td class="title">状态码</td>
				<td>
					<select class="ui dropdown auto-width" v-model="redirectToHttpsConfig.status">
						<option value="0">[使用默认]</option>
						<option v-for="option in statusOptions" :value="option.code">{{option.code}} {{option.text}}</option>
					</select>
				</td>
			</tr>
			<tr>
				<td>跳转后域名或IP地址</td>
				<td>
					<input type="text" name="host" v-model="redirectToHttpsConfig.host"/>
					<p class="comment">默认和用户正在访问的域名或IP地址一致，不填写就表示使用当前的域名。</p>
				</td>
			</tr>
			<tr>
				<td>端口</td>
				<td>
					<input type="text" name="port" v-model="portString" maxlength="5" style="width:6em"/>
					<p class="comment">默认端口为443。</p>
				</td>
			</tr>
			<tr>
				<td>允许的域名</td>
				<td>
					<values-box :values="redirectToHttpsConfig.onlyDomains" @change="changeOnlyDomains"></values-box>
					<p class="comment">如果填写了允许的域名，那么只有这些域名可以自动跳转。</p>
				</td>
			</tr>
			<tr>
				<td>排除的域名</td>
				<td>
					<values-box :values="redirectToHttpsConfig.exceptDomains" @change="changeExceptDomains"></values-box>
					<p class="comment">如果填写了排除的域名，那么这些域名将不跳转。</p>
				</td>
			</tr>
		</table>
	</div>
	<div class="margin"></div>
</div>`
})

// 动作选择
Vue.component("http-firewall-actions-box", {
	props: ["v-actions", "v-firewall-policy", "v-action-configs"],
	mounted: function () {
		let that = this
		Tea.action("/servers/iplists/levelOptions")
			.success(function (resp) {
				that.ipListLevels = resp.data.levels
			})
			.post()

		this.loadJS(function () {
			let box = document.getElementById("actions-box")
			Sortable.create(box, {
				draggable: ".label",
				handle: ".icon.handle",
				onStart: function () {
					that.cancel()
				},
				onUpdate: function (event) {
					let labels = box.getElementsByClassName("label")
					let newConfigs = []
					for (let i = 0; i < labels.length; i++) {
						let index = parseInt(labels[i].getAttribute("data-index"))
						newConfigs.push(that.configs[index])
					}
					that.configs = newConfigs
				}
			})
		})
	},
	data: function () {
		if (this.vFirewallPolicy.inbound == null) {
			this.vFirewallPolicy.inbound = {}
		}
		if (this.vFirewallPolicy.inbound.groups == null) {
			this.vFirewallPolicy.inbound.groups = []
		}

		let id = 0
		let configs = []
		if (this.vActionConfigs != null) {
			configs = this.vActionConfigs
			configs.forEach(function (v) {
				v.id = (id++)
			})
		}

		var defaultPageBody = `<!DOCTYPE html>
<html>
<body>
403 Forbidden
</body>
</html>`


		return {
			id: id,

			actions: this.vActions,
			configs: configs,
			isAdding: false,
			editingIndex: -1,

			action: null,
			actionCode: "",
			actionOptions: {},

			// IPList相关
			ipListLevels: [],

			// 动作参数
			blockTimeout: "",
			blockScope: "global",

			captchaLife: "",
			get302Life: "",
			post307Life: "",
			recordIPType: "black",
			recordIPLevel: "critical",
			recordIPTimeout: "",
			recordIPListId: 0,
			recordIPListName: "",

			tagTags: [],

			pageStatus: 403,
			pageBody: defaultPageBody,
			defaultPageBody: defaultPageBody,

			goGroupName: "",
			goGroupId: 0,
			goGroup: null,

			goSetId: 0,
			goSetName: ""
		}
	},
	watch: {
		actionCode: function (code) {
			this.action = this.actions.$find(function (k, v) {
				return v.code == code
			})
			this.actionOptions = {}
		},
		blockTimeout: function (v) {
			v = parseInt(v)
			if (isNaN(v)) {
				this.actionOptions["timeout"] = 0
			} else {
				this.actionOptions["timeout"] = v
			}
		},
		blockScope: function (v) {
			this.actionOptions["scope"] = v
		},
		captchaLife: function (v) {
			v = parseInt(v)
			if (isNaN(v)) {
				this.actionOptions["life"] = 0
			} else {
				this.actionOptions["life"] = v
			}
		},
		get302Life: function (v) {
			v = parseInt(v)
			if (isNaN(v)) {
				this.actionOptions["life"] = 0
			} else {
				this.actionOptions["life"] = v
			}
		},
		post307Life: function (v) {
			v = parseInt(v)
			if (isNaN(v)) {
				this.actionOptions["life"] = 0
			} else {
				this.actionOptions["life"] = v
			}
		},
		recordIPType: function (v) {
			this.recordIPListId = 0
		},
		recordIPTimeout: function (v) {
			v = parseInt(v)
			if (isNaN(v)) {
				this.actionOptions["timeout"] = 0
			} else {
				this.actionOptions["timeout"] = v
			}
		},
		goGroupId: function (groupId) {
			let group = this.vFirewallPolicy.inbound.groups.$find(function (k, v) {
				return v.id == groupId
			})
			this.goGroup = group
			if (group == null) {
				this.goGroupName = ""
			} else {
				this.goGroupName = group.name
			}
			this.goSetId = 0
			this.goSetName = ""
		},
		goSetId: function (setId) {
			if (this.goGroup == null) {
				return
			}
			let set = this.goGroup.sets.$find(function (k, v) {
				return v.id == setId
			})
			if (set == null) {
				this.goSetId = 0
				this.goSetName = ""
			} else {
				this.goSetName = set.name
			}
		}
	},
	methods: {
		add: function () {
			this.action = null
			this.actionCode = "block"
			this.isAdding = true
			this.actionOptions = {}

			// 动作参数
			this.blockTimeout = ""
			this.blockScope = "global"
			this.captchaLife = ""
			this.get302Life = ""
			this.post307Life = ""

			this.recordIPLevel = "critical"
			this.recordIPType = "black"
			this.recordIPTimeout = ""
			this.recordIPListId = 0
			this.recordIPListName = ""

			this.tagTags = []

			this.pageStatus = 403
			this.pageBody = this.defaultPageBody

			this.goGroupName = ""
			this.goGroupId = 0
			this.goGroup = null

			this.goSetId = 0
			this.goSetName = ""

			let that = this
			this.action = this.vActions.$find(function (k, v) {
				return v.code == that.actionCode
			})

			// 滚到界面底部
			this.scroll()
		},
		remove: function (index) {
			this.isAdding = false
			this.editingIndex = -1
			this.configs.$remove(index)
		},
		update: function (index, config) {
			if (this.isAdding && this.editingIndex == index) {
				this.cancel()
				return
			}

			this.add()

			this.isAdding = true
			this.editingIndex = index

			this.actionCode = config.code

			switch (config.code) {
				case "block":
					this.blockTimeout = ""
					if (config.options.timeout != null || config.options.timeout > 0) {
						this.blockTimeout = config.options.timeout.toString()
					}
					if (config.options.scope != null && config.options.scope.length > 0) {
						this.blockScope = config.options.scope
					} else {
						this.blockScope = "global" // 兼容先前版本遗留的默认值
					}
					break
				case "allow":
					break
				case "log":
					break
				case "captcha":
					this.captchaLife = ""
					if (config.options.life != null || config.options.life > 0) {
						this.captchaLife = config.options.life.toString()
					}
					break
				case "notify":
					break
				case "get_302":
					this.get302Life = ""
					if (config.options.life != null || config.options.life > 0) {
						this.get302Life = config.options.life.toString()
					}
					break
				case "post_307":
					this.post307Life = ""
					if (config.options.life != null || config.options.life > 0) {
						this.post307Life = config.options.life.toString()
					}
					break;
				case "record_ip":
					if (config.options != null) {
						this.recordIPLevel = config.options.level
						this.recordIPType = config.options.type
						if (config.options.timeout > 0) {
							this.recordIPTimeout = config.options.timeout.toString()
						}
						let that = this

						// VUE需要在函数执行完之后才会调用watch函数，这样会导致设置的值被覆盖，所以这里使用setTimeout
						setTimeout(function () {
							that.recordIPListId = config.options.ipListId
							that.recordIPListName = config.options.ipListName
						})
					}
					break
				case "tag":
					this.tagTags = []
					if (config.options.tags != null) {
						this.tagTags = config.options.tags
					}
					break
				case "page":
					this.pageStatus = 403
					this.pageBody = this.defaultPageBody
					if (config.options.status != null) {
						this.pageStatus = config.options.status
					}
					if (config.options.body != null) {
						this.pageBody = config.options.body
					}

					break
				case "go_group":
					if (config.options != null) {
						this.goGroupName = config.options.groupName
						this.goGroupId = config.options.groupId
						this.goGroup = this.vFirewallPolicy.inbound.groups.$find(function (k, v) {
							return v.id == config.options.groupId
						})
					}
					break
				case "go_set":
					if (config.options != null) {
						this.goGroupName = config.options.groupName
						this.goGroupId = config.options.groupId
						this.goGroup = this.vFirewallPolicy.inbound.groups.$find(function (k, v) {
							return v.id == config.options.groupId
						})

						// VUE需要在函数执行完之后才会调用watch函数，这样会导致设置的值被覆盖，所以这里使用setTimeout
						let that = this
						setTimeout(function () {
							that.goSetId = config.options.setId
							if (that.goGroup != null) {
								let set = that.goGroup.sets.$find(function (k, v) {
									return v.id == config.options.setId
								})
								if (set != null) {
									that.goSetName = set.name
								}
							}
						})
					}
					break
			}

			// 滚到界面底部
			this.scroll()
		},
		cancel: function () {
			this.isAdding = false
			this.editingIndex = -1
		},
		confirm: function () {
			if (this.action == null) {
				return
			}

			if (this.actionOptions == null) {
				this.actionOptions = {}
			}

			// record_ip
			if (this.actionCode == "record_ip") {
				let timeout = parseInt(this.recordIPTimeout)
				if (isNaN(timeout)) {
					timeout = 0
				}
				if (this.recordIPListId <= 0) {
					return
				}
				this.actionOptions = {
					type: this.recordIPType,
					level: this.recordIPLevel,
					timeout: timeout,
					ipListId: this.recordIPListId,
					ipListName: this.recordIPListName
				}
			} else if (this.actionCode == "tag") { // tag
				if (this.tagTags == null || this.tagTags.length == 0) {
					return
				}
				this.actionOptions = {
					tags: this.tagTags
				}
			} else if (this.actionCode == "page") {
				let pageStatus = this.pageStatus.toString()
				if (!pageStatus.match(/^\d{3}$/)) {
					pageStatus = 403
				} else {
					pageStatus = parseInt(pageStatus)
				}

				this.actionOptions = {
					status: pageStatus,
					body: this.pageBody
				}
			} else if (this.actionCode == "go_group") { // go_group
				let groupId = this.goGroupId
				if (typeof (groupId) == "string") {
					groupId = parseInt(groupId)
					if (isNaN(groupId)) {
						groupId = 0
					}
				}
				if (groupId <= 0) {
					return
				}
				this.actionOptions = {
					groupId: groupId.toString(),
					groupName: this.goGroupName
				}
			} else if (this.actionCode == "go_set") { // go_set
				let groupId = this.goGroupId
				if (typeof (groupId) == "string") {
					groupId = parseInt(groupId)
					if (isNaN(groupId)) {
						groupId = 0
					}
				}

				let setId = this.goSetId
				if (typeof (setId) == "string") {
					setId = parseInt(setId)
					if (isNaN(setId)) {
						setId = 0
					}
				}
				if (setId <= 0) {
					return
				}
				this.actionOptions = {
					groupId: groupId.toString(),
					groupName: this.goGroupName,
					setId: setId.toString(),
					setName: this.goSetName
				}
			}

			let options = {}
			for (let k in this.actionOptions) {
				if (this.actionOptions.hasOwnProperty(k)) {
					options[k] = this.actionOptions[k]
				}
			}
			if (this.editingIndex > -1) {
				this.configs[this.editingIndex] = {
					id: this.configs[this.editingIndex].id,
					code: this.actionCode,
					name: this.action.name,
					options: options
				}
			} else {
				this.configs.push({
					id: (this.id++),
					code: this.actionCode,
					name: this.action.name,
					options: options
				})
			}

			this.cancel()
		},
		removeRecordIPList: function () {
			this.recordIPListId = 0
		},
		selectRecordIPList: function () {
			let that = this
			teaweb.popup("/servers/iplists/selectPopup?type=" + this.recordIPType, {
				width: "50em",
				height: "30em",
				callback: function (resp) {
					that.recordIPListId = resp.data.list.id
					that.recordIPListName = resp.data.list.name
				}
			})
		},
		changeTags: function (tags) {
			this.tagTags = tags
		},
		loadJS: function (callback) {
			if (typeof Sortable != "undefined") {
				callback()
				return
			}

			// 引入js
			let jsFile = document.createElement("script")
			jsFile.setAttribute("src", "/js/sortable.min.js")
			jsFile.addEventListener("load", function () {
				callback()
			})
			document.head.appendChild(jsFile)
		},
		scroll: function () {
			setTimeout(function () {
				let mainDiv = document.getElementsByClassName("main")
				if (mainDiv.length > 0) {
					mainDiv[0].scrollTo(0, 1000)
				}
			}, 10)
		}
	},
	template: `<div>
	<input type="hidden" name="actionsJSON" :value="JSON.stringify(configs)"/>
	<div v-show="configs.length > 0" style="margin-bottom: 0.5em" id="actions-box"> 
		<div v-for="(config, index) in configs" :data-index="index" :key="config.id" class="ui label small basic" :class="{blue: index == editingIndex}" style="margin-bottom: 0.4em">
			{{config.name}} <span class="small">({{config.code.toUpperCase()}})</span> 
			
			<!-- block -->
			<span v-if="config.code == 'block' && config.options.timeout > 0">：有效期{{config.options.timeout}}秒</span>
			
			<!-- captcha -->
			<span v-if="config.code == 'captcha' && config.options.life > 0">：有效期{{config.options.life}}秒</span>
			
			<!-- get 302 -->
			<span v-if="config.code == 'get_302' && config.options.life > 0">：有效期{{config.options.life}}秒</span>
			
			<!-- post 307 -->
			<span v-if="config.code == 'post_307' && config.options.life > 0">：有效期{{config.options.life}}秒</span>
			
			<!-- record_ip -->
			<span v-if="config.code == 'record_ip'">：{{config.options.ipListName}}</span>
			
			<!-- tag -->
			<span v-if="config.code == 'tag'">：{{config.options.tags.join(", ")}}</span>
			
			<!-- page -->
			<span v-if="config.code == 'page'">：[{{config.options.status}}]</span>
			
			<!-- go_group -->
			<span v-if="config.code == 'go_group'">：{{config.options.groupName}}</span>
			
			<!-- go_set -->
			<span v-if="config.code == 'go_set'">：{{config.options.groupName}} / {{config.options.setName}}</span>
			
			<!-- 范围 -->
			<span v-if="config.options.scope != null && config.options.scope.length > 0" class="small grey">
				&nbsp; 
				<span v-if="config.options.scope == 'global'">[所有服务]</span>
				<span v-if="config.options.scope == 'service'">[当前服务]</span>
			</span>
			
			<!-- 操作按钮 -->
			 &nbsp; <a href="" title="修改" @click.prevent="update(index, config)"><i class="icon pencil small"></i></a> &nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a> &nbsp; <a href="" title="拖动改变顺序"><i class="icon bars handle"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div style="margin-bottom: 0.5em" v-if="isAdding">
		<table class="ui table" :class="{blue: editingIndex > -1}">
			<tr>
				<td class="title">动作类型 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="actionCode">
						<option v-for="action in actions" :value="action.code">{{action.name}} ({{action.code.toUpperCase()}})</option>
					</select>
					<p class="comment" v-if="action != null && action.description.length > 0">{{action.description}}</p>
				</td>
			</tr>
			
			<!-- block -->
			<tr v-if="actionCode == 'block'">
				<td>封锁时间</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 5em" maxlength="10" v-model="blockTimeout" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">秒</span>
					</div>
				</td>
			</tr>
			<tr v-if="actionCode == 'block'">
				<td>封锁范围</td>
				<td>
					<select class="ui dropdown auto-width" v-model="blockScope">
						<option value="service">当前服务</option>
						<option value="global">所有服务</option>
					</select>
				</td>
			</tr>
			
			<!-- captcha -->
			<tr v-if="actionCode == 'captcha'">
				<td>有效时间</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 5em" maxlength="10" v-model="captchaLife" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">秒</span>
					</div>
					<p class="comment">验证通过后在这个时间内不再验证，默认600秒。</p>
				</td>
			</tr>
			
			<!-- get_302 -->
			<tr v-if="actionCode == 'get_302'">
				<td>有效时间</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 5em" maxlength="10" v-model="get302Life" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">秒</span>
					</div>
					<p class="comment">验证通过后在这个时间内不再验证。</p>
				</td>
			</tr>
			
			<!-- post_307 -->
			<tr v-if="actionCode == 'post_307'">
				<td>有效时间</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 5em" maxlength="10" v-model="post307Life" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">秒</span>
					</div>
					<p class="comment">验证通过后在这个时间内不再验证。</p>
				</td>
			</tr>
			
			<!-- record_ip -->
			<tr v-if="actionCode == 'record_ip'">
				<td>IP名单类型 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="recordIPType">
					<option value="black">黑名单</option>
					<option value="white">白名单</option>
					</select>
				</td>
			</tr>
			<tr v-if="actionCode == 'record_ip'">
				<td>选择IP名单 *</td>
				<td>
					<div v-if="recordIPListId > 0" class="ui label basic small">{{recordIPListName}} <a href="" @click.prevent="removeRecordIPList"><i class="icon remove small"></i></a></div>
					<button type="button" class="ui button tiny" @click.prevent="selectRecordIPList">+</button>
					<p class="comment">如不选择，则自动添加到当前策略的IP名单中。</p>
				</td>
			</tr>
			<tr v-if="actionCode == 'record_ip'">
				<td>级别</td>
				<td>
					<select class="ui dropdown auto-width" v-model="recordIPLevel">
						<option v-for="level in ipListLevels" :value="level.code">{{level.name}}</option>
					</select>
				</td>
			</tr>
			<tr v-if="actionCode == 'record_ip'">
				<td>超时时间</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 6em" maxlength="10" v-model="recordIPTimeout" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">秒</span>
					</div>
					<p class="comment">0表示不超时。</p>
				</td>
			</tr>
			
			<!-- tag -->
			<tr v-if="actionCode == 'tag'">
				<td>标签 *</td>
				<td>
					<values-box @change="changeTags" :values="tagTags"></values-box>
				</td>
			</tr>
			
			<!-- page -->
			<tr v-if="actionCode == 'page'">
				<td>状态码 *</td>
				<td><input type="text" style="width: 4em" maxlength="3" v-model="pageStatus"/></td>
			</tr>
			<tr v-if="actionCode == 'page'">
				<td>网页内容</td>
				<td>
					<textarea v-model="pageBody"></textarea>
				</td>
			</tr>
			
			<!-- 规则分组 -->
			<tr v-if="actionCode == 'go_group'">
				<td>下一个分组 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="goGroupId">
						<option value="0">[选择分组]</option>
						<option v-for="group in vFirewallPolicy.inbound.groups" :value="group.id">{{group.name}}</option>
					</select>
				</td>
			</tr>
			
			<!-- 规则集 -->
			<tr v-if="actionCode == 'go_set'">
				<td>下一个分组 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="goGroupId">
						<option value="0">[选择分组]</option>
						<option v-for="group in vFirewallPolicy.inbound.groups" :value="group.id">{{group.name}}</option>
					</select>
				</td>
			</tr>
			<tr v-if="actionCode == 'go_set' && goGroup != null">
				<td>下一个规则集 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="goSetId">
						<option value="0">[选择规则集]</option>
						<option v-for="set in goGroup.sets" :value="set.id">{{set.name}}</option>
					</select>
				</td>
			</tr>
		</table>
		<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button> &nbsp;
		<a href="" @click.prevent="cancel" title="取消"><i class="icon remove small"></i></a>
	</div>
	<div v-if="!isAdding">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
	<p class="comment">系统总是会先执行记录日志、标签等不会修改请求的动作，再执行阻止、验证码等可能改变请求的动作。</p>
</div>`
})

// 认证设置
Vue.component("http-auth-config-box", {
	props: ["v-auth-config", "v-is-location"],
	data: function () {
		let authConfig = this.vAuthConfig
		if (authConfig == null) {
			authConfig = {
				isPrior: false,
				isOn: false
			}
		}
		if (authConfig.policyRefs == null) {
			authConfig.policyRefs = []
		}
		return {
			authConfig: authConfig
		}
	},
	methods: {
		isOn: function () {
			return (!this.vIsLocation || this.authConfig.isPrior) && this.authConfig.isOn
		},
		add: function () {
			let that = this
			teaweb.popup("/servers/server/settings/access/createPopup", {
				callback: function (resp) {
					that.authConfig.policyRefs.push(resp.data.policyRef)
				},
				height: "28em"
			})
		},
		update: function (index, policyId) {
			let that = this
			teaweb.popup("/servers/server/settings/access/updatePopup?policyId=" + policyId, {
				callback: function (resp) {
					teaweb.success("保存成功", function () {
						teaweb.reload()
					})
				},
				height: "28em"
			})
		},
		remove: function (index) {
			this.authConfig.policyRefs.$remove(index)
		},
		methodName: function (methodType) {
			switch (methodType) {
				case "basicAuth":
					return "BasicAuth"
				case "subRequest":
					return "子请求"
			}
			return ""
		}
	},
	template: `<div>
<input type="hidden" name="authJSON" :value="JSON.stringify(authConfig)"/> 
<table class="ui table selectable definition">
	<prior-checkbox :v-config="authConfig" v-if="vIsLocation"></prior-checkbox>
	<tbody v-show="!vIsLocation || authConfig.isPrior">
		<tr>
			<td class="title">启用认证</td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" v-model="authConfig.isOn"/>
					<label></label>
				</div>
			</td>
		</tr>
	</tbody>
</table>
<div class="margin"></div>
<!-- 认证方式 -->
<div v-show="isOn()">
	<h4>认证方式</h4>
	<table class="ui table selectable celled" v-show="authConfig.policyRefs.length > 0">
		<thead>
			<tr>
				<th class="three wide">名称</th>
				<th class="three wide">认证方法</th>
				<th>参数</th>
				<th class="two wide">状态</th>
				<th class="two op">操作</th>
			</tr>
		</thead>
		<tbody v-for="(ref, index) in authConfig.policyRefs" :key="ref.authPolicyId">
			<tr>
				<td>{{ref.authPolicy.name}}</td>
				<td>
					{{methodName(ref.authPolicy.type)}}
				</td>
				<td>
					<span v-if="ref.authPolicy.type == 'basicAuth'">{{ref.authPolicy.params.users.length}}个用户</span>
					<span v-if="ref.authPolicy.type == 'subRequest'">
						<span v-if="ref.authPolicy.params.method.length > 0" class="grey">[{{ref.authPolicy.params.method}}]</span>
						{{ref.authPolicy.params.url}}
					</span>
				</td>
				<td>
					<label-on :v-is-on="ref.authPolicy.isOn"></label-on>
				</td>
				<td>
					<a href="" @click.prevent="update(index, ref.authPolicyId)">修改</a> &nbsp;
					<a href="" @click.prevent="remove(index)">删除</a>
				</td>
			</tr>
		</tbody>
	</table>
	<button class="ui button small" type="button" @click.prevent="add">+添加认证方式</button>
</div>
<div class="margin"></div>
</div>`
})

Vue.component("user-selector", {
	mounted: function () {
		let that = this

		Tea.action("/servers/users/options")
			.post()
			.success(function (resp) {
				that.users = resp.data.users
			})
	},
	props: ["v-user-id"],
	data: function () {
		let userId = this.vUserId
		if (userId == null) {
			userId = 0
		}
		return {
			users: [],
			userId: userId
		}
	},
	watch: {
		userId: function (v) {
			this.$emit("change", v)
		}
	},
	template: `<div>
	<select class="ui dropdown auto-width" name="userId" v-model="userId">
		<option value="0">[选择用户]</option>
		<option v-for="user in users" :value="user.id">{{user.fullname}} ({{user.username}})</option>
	</select>
</div>`
})

Vue.component("http-header-policy-box", {
	props: ["v-request-header-policy", "v-request-header-ref", "v-response-header-policy", "v-response-header-ref", "v-params", "v-is-location", "v-is-group", "v-has-group-request-config", "v-has-group-response-config", "v-group-setting-url"],
	data: function () {
		let type = "request"
		let hash = window.location.hash
		if (hash == "#response") {
			type = "response"
		}

		// ref
		let requestHeaderRef = this.vRequestHeaderRef
		if (requestHeaderRef == null) {
			requestHeaderRef = {
				isPrior: false,
				isOn: true,
				headerPolicyId: 0
			}
		}

		let responseHeaderRef = this.vResponseHeaderRef
		if (responseHeaderRef == null) {
			responseHeaderRef = {
				isPrior: false,
				isOn: true,
				headerPolicyId: 0
			}
		}

		// 请求相关
		let requestSettingHeaders = []
		let requestDeletingHeaders = []

		let requestPolicy = this.vRequestHeaderPolicy
		if (requestPolicy != null) {
			if (requestPolicy.setHeaders != null) {
				requestSettingHeaders = requestPolicy.setHeaders
			}
			if (requestPolicy.deleteHeaders != null) {
				requestDeletingHeaders = requestPolicy.deleteHeaders
			}
		}

		// 响应相关
		let responseSettingHeaders = []
		let responseDeletingHeaders = []

		let responsePolicy = this.vResponseHeaderPolicy
		if (responsePolicy != null) {
			if (responsePolicy.setHeaders != null) {
				responseSettingHeaders = responsePolicy.setHeaders
			}
			if (responsePolicy.deleteHeaders != null) {
				responseDeletingHeaders = responsePolicy.deleteHeaders
			}
		}
		
		return {
			type: type,
			typeName: (type == "request") ? "请求" : "响应",
			requestHeaderRef: requestHeaderRef,
			responseHeaderRef: responseHeaderRef,
			requestSettingHeaders: requestSettingHeaders,
			requestDeletingHeaders: requestDeletingHeaders,
			responseSettingHeaders: responseSettingHeaders,
			responseDeletingHeaders: responseDeletingHeaders
		}
	},
	methods: {
		selectType: function (type) {
			this.type = type
			window.location.hash = "#" + type
			window.location.reload()
		},
		addSettingHeader: function (policyId) {
			teaweb.popup("/servers/server/settings/headers/createSetPopup?" + this.vParams + "&headerPolicyId=" + policyId, {
				callback: function () {
					window.location.reload()
				}
			})
		},
		addDeletingHeader: function (policyId, type) {
			teaweb.popup("/servers/server/settings/headers/createDeletePopup?" + this.vParams + "&headerPolicyId=" + policyId + "&type=" + type, {
				callback: function () {
					window.location.reload()
				}
			})
		},
		updateSettingPopup: function (policyId, headerId) {
			teaweb.popup("/servers/server/settings/headers/updateSetPopup?" + this.vParams + "&headerPolicyId=" + policyId + "&headerId=" + headerId, {
				callback: function () {
					window.location.reload()
				}
			})
		},
		deleteDeletingHeader: function (policyId, headerName) {
			teaweb.confirm("确定要删除'" + headerName + "'吗？", function () {
				Tea.action("/servers/server/settings/headers/deleteDeletingHeader")
					.params({
						headerPolicyId: policyId,
						headerName: headerName
					})
					.post()
					.refresh()
			})
		},
		deleteHeader: function (policyId, type, headerId) {
			teaweb.confirm("确定要删除此Header吗？", function () {
					this.$post("/servers/server/settings/headers/delete")
						.params({
							headerPolicyId: policyId,
							type: type,
							headerId: headerId
						})
						.refresh()
				}
			)
		}
	},
	template: `<div>
	<div class="ui menu tabular small">
		<a class="item" :class="{active:type == 'request'}" @click.prevent="selectType('request')">请求Header<span v-if="requestSettingHeaders.length > 0">({{requestSettingHeaders.length}})</span></a>
		<a class="item" :class="{active:type == 'response'}" @click.prevent="selectType('response')">响应Header<span v-if="responseSettingHeaders.length > 0">({{responseSettingHeaders.length}})</span></a>
	</div>
	
	<div class="margin"></div>
	
	<input type="hidden" name="type" :value="type"/>
	
	<!-- 请求 -->
	<div v-if="(vIsLocation || vIsGroup) && type == 'request'">
		<input type="hidden" name="requestHeaderJSON" :value="JSON.stringify(requestHeaderRef)"/>
		<table class="ui table definition selectable">
			<prior-checkbox :v-config="requestHeaderRef"></prior-checkbox>
		</table>
		<submit-btn></submit-btn>
	</div>
	
	<div v-if="((!vIsLocation && !vIsGroup) || requestHeaderRef.isPrior) && type == 'request'">
		<div v-if="vHasGroupRequestConfig">
        	<div class="margin"></div>
        	<warning-message>由于已经在当前<a :href="vGroupSettingUrl + '#request'">服务分组</a>中进行了对应的配置，在这里的配置将不会生效。</warning-message>
    	</div>
    	<div :class="{'opacity-mask': vHasGroupRequestConfig}">
		<h3>设置请求Header <a href="" @click.prevent="addSettingHeader(vRequestHeaderPolicy.id)">[添加新Header]</a></h3>
			<p class="comment" v-if="requestSettingHeaders.length == 0">暂时还没有Header。</p>
			<table class="ui table selectable celled" v-if="requestSettingHeaders.length > 0">
				<thead>
					<tr>
						<th>名称</th>
						<th>值</th>
						<th class="two op">操作</th>
					</tr>
				</thead>
				<tr v-for="header in requestSettingHeaders">
					<td class="five wide">{{header.name}}</td>
					<td>{{header.value}}</td>
					<td><a href="" @click.prevent="updateSettingPopup(vRequestHeaderPolicy.id, header.id)">修改</a> &nbsp; <a href="" @click.prevent="deleteHeader(vRequestHeaderPolicy.id, 'setHeader', header.id)">删除</a> </td>
				</tr>
			</table>
			
			<h3>删除请求Header</h3>
			<p class="comment">这里可以设置需要从请求中删除的Header。</p>
			
			<table class="ui table definition selectable">
				<td class="title">需要删除的Header</td>
				<td>
					<div v-if="requestDeletingHeaders.length > 0">
						<div class="ui label small basic" v-for="headerName in requestDeletingHeaders">{{headerName}} <a href=""><i class="icon remove" title="删除" @click.prevent="deleteDeletingHeader(vRequestHeaderPolicy.id, headerName)"></i></a> </div>
						<div class="ui divider" ></div>
					</div>
					<button class="ui button small" type="button" @click.prevent="addDeletingHeader(vRequestHeaderPolicy.id, 'request')">+</button>
				</td>
			</table>
		</div>			
	</div>
	
	<!-- 响应 -->
	<div v-if="(vIsLocation || vIsGroup) && type == 'response'">
		<input type="hidden" name="responseHeaderJSON" :value="JSON.stringify(responseHeaderRef)"/>
		<table class="ui table definition selectable">
			<prior-checkbox :v-config="responseHeaderRef"></prior-checkbox>
		</table>
		<submit-btn></submit-btn>
	</div>
	
	<div v-if="((!vIsLocation && !vIsGroup) || requestHeaderRef.isPrior) && type == 'response'">
		<div v-if="vHasGroupResponseConfig">
        	<div class="margin"></div>
        	<warning-message>由于已经在当前<a :href="vGroupSettingUrl + '#response'">服务分组</a>中进行了对应的配置，在这里的配置将不会生效。</warning-message>
    	</div>
    	<div :class="{'opacity-mask': vHasGroupResponseConfig}">
			<h3>设置响应Header <a href="" @click.prevent="addSettingHeader(vResponseHeaderPolicy.id)">[添加新Header]</a></h3>
			<p class="comment" style="margin-top: 0; padding-top: 0">将会覆盖已有的同名Header。</p>
			<p class="comment" v-if="responseSettingHeaders.length == 0">暂时还没有Header。</p>
			<table class="ui table selectable celled" v-if="responseSettingHeaders.length > 0">
				<thead>
					<tr>
						<th>名称</th>
						<th>值</th>
						<th class="two op">操作</th>
					</tr>
				</thead>
				<tr v-for="header in responseSettingHeaders">
					<td class="five wide">{{header.name}}</td>
					<td>{{header.value}}</td>
					<td><a href="" @click.prevent="updateSettingPopup(vResponseHeaderPolicy.id, header.id)">修改</a> &nbsp; <a href="" @click.prevent="deleteHeader(vResponseHeaderPolicy.id, 'setHeader', header.id)">删除</a> </td>
				</tr>
			</table>
			
			<h3>删除响应Header</h3>
			<p class="comment">这里可以设置需要从响应中删除的Header。</p>
			
			<table class="ui table definition selectable">
				<td class="title">需要删除的Header</td>
				<td>
					<div v-if="responseDeletingHeaders.length > 0">
						<div class="ui label small basic" v-for="headerName in responseDeletingHeaders">{{headerName}} <a href=""><i class="icon remove" title="删除" @click.prevent="deleteDeletingHeader(vResponseHeaderPolicy.id, headerName)"></i></a> </div>
						<div class="ui divider" ></div>
					</div>
					<button class="ui button small" type="button" @click.prevent="addDeletingHeader(vResponseHeaderPolicy.id, 'response')">+</button>
				</td>
			</table>
		</div>			
	</div>
	<div class="margin"></div>
</div>`
})

Vue.component("http-cache-policy-selector", {
	props: ["v-cache-policy"],
	mounted: function () {
		let that = this
		Tea.action("/servers/components/cache/count")
			.post()
			.success(function (resp) {
				that.count = resp.data.count
			})
	},
	data: function () {
		let cachePolicy = this.vCachePolicy
		return {
			count: 0,
			cachePolicy: cachePolicy
		}
	},
	methods: {
		remove: function () {
			this.cachePolicy = null
		},
		select: function () {
			let that = this
			teaweb.popup("/servers/components/cache/selectPopup", {
				callback: function (resp) {
					that.cachePolicy = resp.data.cachePolicy
				}
			})
		},
		create: function () {
			let that = this
			teaweb.popup("/servers/components/cache/createPopup", {
				height: "26em",
				callback: function (resp) {
					that.cachePolicy = resp.data.cachePolicy
				}
			})
		}
	},
	template: `<div>
	<div v-if="cachePolicy != null" class="ui label basic">
		<input type="hidden" name="cachePolicyId" :value="cachePolicy.id"/>
		{{cachePolicy.name}} &nbsp; <a :href="'/servers/components/cache/policy?cachePolicyId=' + cachePolicy.id" target="_blank" title="修改"><i class="icon pencil small"></i></a>&nbsp; <a href="" @click.prevent="remove()" title="删除"><i class="icon remove small"></i></a>
	</div>
	<div v-if="cachePolicy == null">
		<span v-if="count > 0"><a href="" @click.prevent="select">[选择已有策略]</a> &nbsp; &nbsp; </span><a href="" @click.prevent="create">[创建新策略]</a>
	</div>
</div>`
})

Vue.component("http-pages-and-shutdown-box", {
	props: ["v-pages", "v-shutdown-config", "v-is-location"],
	data: function () {
		let pages = []
		if (this.vPages != null) {
			pages = this.vPages
		}
		let shutdownConfig = {
			isPrior: false,
			isOn: false,
			bodyType: "url",
			url: "",
			body: "",
			status: 0
		}
		if (this.vShutdownConfig != null) {
			if (this.vShutdownConfig.body == null) {
				this.vShutdownConfig.body = ""
			}
			if (this.vShutdownConfig.bodyType == null) {
				this.vShutdownConfig.bodyType = "url"
			}
			shutdownConfig = this.vShutdownConfig
		}

		let shutdownStatus = ""
		if (shutdownConfig.status > 0) {
			shutdownStatus = shutdownConfig.status.toString()
		}

		return {
			pages: pages,
			shutdownConfig: shutdownConfig,
			shutdownStatus: shutdownStatus
		}
	},
	watch: {
		shutdownStatus: function (status) {
			let statusInt = parseInt(status)
			if (!isNaN(statusInt) && statusInt > 0 && statusInt < 1000) {
				this.shutdownConfig.status = statusInt
			} else {
				this.shutdownConfig.status = 0
			}
		}
	},
	methods: {
		addPage: function () {
			let that = this
			teaweb.popup("/servers/server/settings/pages/createPopup", {
				height: "26em",
				callback: function (resp) {
					that.pages.push(resp.data.page)
				}
			})
		},
		updatePage: function (pageIndex, pageId) {
			let that = this
			teaweb.popup("/servers/server/settings/pages/updatePopup?pageId=" + pageId, {
				height: "26em",
				callback: function (resp) {
					Vue.set(that.pages, pageIndex, resp.data.page)
				}
			})
		},
		removePage: function (pageIndex) {
			let that = this
			teaweb.confirm("确定要移除此页面吗？", function () {
				that.pages.$remove(pageIndex)
			})
		},
		addShutdownHTMLTemplate: function () {
			this.shutdownConfig.body  = `<!DOCTYPE html>
<html>
<head>
\t<title>升级中</title>
\t<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
</head>
<body>

<h3>网站升级中</h3>
<p>为了给您提供更好的服务，我们正在升级网站，请稍后重新访问。</p>

<footer>Powered by GoEdge.</footer>

</body>
</html>`
		}
	},
	template: `<div>
<input type="hidden" name="pagesJSON" :value="JSON.stringify(pages)"/>
<input type="hidden" name="shutdownJSON" :value="JSON.stringify(shutdownConfig)"/>
<table class="ui table selectable definition">
	<tr>
		<td class="title">特殊页面</td>
		<td>
			<div v-if="pages.length > 0">
				<div class="ui label small basic" v-for="(page,index) in pages">
					{{page.status}} -&gt; <span v-if="page.bodyType == 'url'">{{page.url}}</span><span v-if="page.bodyType == 'html'">[HTML内容]</span> <a href="" title="修改" @click.prevent="updatePage(index, page.id)"><i class="icon pencil small"></i></a> <a href="" title="删除" @click.prevent="removePage(index)"><i class="icon remove"></i></a>
				</div>
				<div class="ui divider"></div>
			</div>
			<div>
				<button class="ui button small" type="button" @click.prevent="addPage()">+</button>
			</div>
			<p class="comment">根据响应状态码返回一些特殊页面，比如404，500等错误页面。</p>
		</td>
	</tr>	
	<tr>
		<td>临时关闭页面</td>
		<td>
			<div>
				<table class="ui table selectable definition">
					<prior-checkbox :v-config="shutdownConfig" v-if="vIsLocation"></prior-checkbox>
					<tbody v-show="!vIsLocation || shutdownConfig.isPrior">
						<tr>
							<td class="title">是否开启</td>
							<td>
								<div class="ui checkbox">
									<input type="checkbox" value="1" v-model="shutdownConfig.isOn" />
									<label></label>
								</div>
							</td>
						</tr>
					</tbody>
					<tbody v-show="(!vIsLocation || shutdownConfig.isPrior) && shutdownConfig.isOn">
						<tr>
							<td>内容类型 *</td>
							<td>
								<select class="ui dropdown auto-width" v-model="shutdownConfig.bodyType">
									<option value="url">读取URL</option>
									<option value="html">HTML</option>
								</select>
							</td>
						</tr>
						<tr v-show="shutdownConfig.bodyType == 'url'">
							<td class="title">页面URL *</td>
							<td>
								<input type="text" v-model="shutdownConfig.url" placeholder="页面文件路径或一个完整URL"/>
								<p class="comment">页面文件是相对于节点安装目录的页面文件比如pages/40x.html，或者一个完整的URL。</p>
							</td>
						</tr>
						<tr v-show="shutdownConfig.bodyType == 'html'">
							<td>HTML *</td>
							<td>
								<textarea name="body" ref="shutdownHTMLBody" v-model="shutdownConfig.body"></textarea>
								<p class="comment"><a href="" @click.prevent="addShutdownHTMLTemplate">[使用模板]</a>。填写页面的HTML内容，支持请求变量。</p>
							</td>
						</tr>
						<tr>
							<td>状态码</td>
							<td><input type="text" size="3" maxlength="3" name="shutdownStatus" style="width:5.2em" placeholder="状态码" v-model="shutdownStatus"/></td>
						</tr>
					</tbody>
				</table>
				<p class="comment">开启临时关闭页面时，所有请求的响应都会显示此页面。可用于临时升级网站使用。</p>
			</div>
		</td>
	</tr>
</table>
<div class="ui margin"></div>
</div>`
})

// 压缩配置
Vue.component("http-compression-config-box", {
	props: ["v-compression-config", "v-is-location", "v-is-group"],
	mounted: function () {
		let that = this
		sortLoad(function () {
			that.initSortableTypes()
		})
	},
	data: function () {
		let config = this.vCompressionConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				useDefaultTypes: true,
				types: ["brotli", "gzip", "deflate"],
				level: 5,
				decompressData: false,
				gzipRef: null,
				deflateRef: null,
				brotliRef: null,
				minLength: {count: 0, "unit": "kb"},
				maxLength: {count: 0, "unit": "kb"},
				mimeTypes: ["text/*", "application/*", "font/*"],
				extensions: [".js", ".json", ".html", ".htm", ".xml", ".css", ".woff2", ".txt"],
				conds: null
			}
		}

		if (config.types == null) {
			config.types = []
		}
		if (config.mimeTypes == null) {
			config.mimeTypes = []
		}
		if (config.extensions == null) {
			config.extensions = []
		}

		let allTypes = [
			{
				name: "Gzip",
				code: "gzip",
				isOn: true
			},
			{
				name: "Deflate",
				code: "deflate",
				isOn: true
			},
			{
				name: "Brotli",
				code: "brotli",
				isOn: true
			}
		]

		let configTypes = []
		config.types.forEach(function (typeCode) {
			allTypes.forEach(function (t) {
				if (typeCode == t.code) {
					t.isOn = true
					configTypes.push(t)
				}
			})
		})
		allTypes.forEach(function (t) {
			if (!config.types.$contains(t.code)) {
				t.isOn = false
				configTypes.push(t)
			}
		})

		return {
			config: config,
			moreOptionsVisible: false,
			allTypes: configTypes
		}
	},
	watch: {
		"config.level": function (v) {
			let level = parseInt(v)
			if (isNaN(level)) {
				level = 1
			} else if (level < 1) {
				level = 1
			} else if (level > 10) {
				level = 10
			}
			this.config.level = level
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.config.isPrior) && this.config.isOn
		},
		changeExtensions: function (values) {
			values.forEach(function (v, k) {
				if (v.length > 0 && v[0] != ".") {
					values[k] = "." + v
				}
			})
			this.config.extensions = values
		},
		changeMimeTypes: function (values) {
			this.config.mimeTypes = values
		},
		changeAdvancedVisible: function () {
			this.moreOptionsVisible = !this.moreOptionsVisible
		},
		changeConds: function (conds) {
			this.config.conds = conds
		},
		changeType: function () {
			this.config.types = []
			let that = this
			this.allTypes.forEach(function (v) {
				if (v.isOn) {
					that.config.types.push(v.code)
				}
			})
		},
		initSortableTypes: function () {
			let box = document.querySelector("#compression-types-box")
			let that = this
			Sortable.create(box, {
				draggable: ".checkbox",
				handle: ".icon.handle",
				onStart: function () {

				},
				onUpdate: function (event) {
					let checkboxes = box.querySelectorAll(".checkbox")
					let codes = []
					checkboxes.forEach(function (checkbox) {
						let code = checkbox.getAttribute("data-code")
						codes.push(code)
					})
					that.config.types = codes
				}
			})
		}
	},
	template: `<div>
	<input type="hidden" name="compressionJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
			<tr>
				<td class="title">是否启用</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="config.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td>压缩级别</td>
				<td>
					<select class="ui dropdown auto-width" v-model="config.level">
						<option v-for="i in 10" :value="i">{{i}}</option>	
					</select>
					<p class="comment">级别越高，压缩比例越大。</p>
				</td>
			</tr>
			<tr>
				<td>支持的扩展名</td>
				<td>
					<values-box :values="config.extensions" @change="changeExtensions" placeholder="比如 .html"></values-box>
					<p class="comment">含有这些扩展名的URL将会被压缩，不区分大小写。</p>
				</td>
			</tr>
			<tr>
				<td>支持的MimeType</td>
				<td>
					<values-box :values="config.mimeTypes" @change="changeMimeTypes" placeholder="比如 text/*"></values-box>
					<p class="comment">响应的Content-Type里包含这些MimeType的内容将会被压缩。</p>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-if="isOn()"></more-options-tbody>
		<tbody v-show="isOn() && moreOptionsVisible">
			<tr>
				<td>压缩算法</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="config.useDefaultTypes" id="compression-use-default"/>
						<label v-if="config.useDefaultTypes" for="compression-use-default">使用默认顺序<span class="grey small">（brotli、gzip、deflate）</span></label>
						<label v-if="!config.useDefaultTypes" for="compression-use-default">使用默认顺序</label>
					</div>
					<div v-show="!config.useDefaultTypes">
						<div class="ui divider"></div>
						<div id="compression-types-box">
							<div class="ui checkbox" v-for="t in allTypes" style="margin-right: 2em" :data-code="t.code">
								<input type="checkbox" v-model="t.isOn" :id="'compression-type-' + t.code" @change="changeType"/>
								<label :for="'compression-type-' + t.code">{{t.name}} &nbsp; <i class="icon list small grey handle"></i></label>
							</div>
						</div>
					</div>
					
					<p class="comment">选择支持的压缩算法和优先顺序，拖动<i class="icon list small grey"></i>图表排序。</p>
				</td>
			</tr>
			<tr>
				<td>支持已压缩内容</td>
				<td>
					<checkbox v-model="config.decompressData"></checkbox>
					<p class="comment">支持对已压缩内容尝试重新使用新的算法压缩。</p>
				</td>
			</tr>
			<tr>
				<td>内容最小长度</td>
				<td>
					<size-capacity-box :v-name="'minLength'" :v-value="config.minLength" :v-unit="'kb'"></size-capacity-box>
					<p class="comment">0表示不限制，内容长度从文件尺寸或Content-Length中获取。</p>
				</td>
			</tr>
			<tr>
				<td>内容最大长度</td>
				<td>
					<size-capacity-box :v-name="'maxLength'" :v-value="config.maxLength" :v-unit="'mb'"></size-capacity-box>
					<p class="comment">0表示不限制，内容长度从文件尺寸或Content-Length中获取。</p>
				</td>
			</tr>
			<tr>
				<td>匹配条件</td>
				<td>
					<http-request-conds-box :v-conds="config.conds" @change="changeConds"></http-request-conds-box>
	</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("firewall-event-level-options", {
    props: ["v-value"],
    mounted: function () {
        let that = this
        Tea.action("/ui/eventLevelOptions")
            .post()
            .success(function (resp) {
                that.levels = resp.data.eventLevels
                that.change()
            })
    },
    data: function () {
        let value = this.vValue
        if (value == null || value.length == 0) {
            value = "" // 不要给默认值，因为黑白名单等默认值均有不同
        }

        return {
            levels: [],
            description: "",
            level: value
        }
    },
    methods: {
        change: function () {
            this.$emit("change")

            let that = this
            let l = this.levels.$find(function (k, v) {
                return v.code == that.level
            })
            if (l != null) {
                this.description = l.description
            } else {
                this.description = ""
            }
        }
    },
    template: `<div>
    <select class="ui dropdown auto-width" name="eventLevel" v-model="level" @change="change">
        <option v-for="level in levels" :value="level.code">{{level.name}}</option>
    </select>
    <p class="comment">{{description}}</p>
</div>`
})

Vue.component("prior-checkbox", {
	props: ["v-config"],
	data: function () {
		return {
			isPrior: this.vConfig.isPrior
		}
	},
	watch: {
		isPrior: function (v) {
			this.vConfig.isPrior = v
		}
	},
	template: `<tbody>
	<tr :class="{active:isPrior}">
		<td class="title">打开独立配置</td>
		<td>
			<div class="ui toggle checkbox">
				<input type="checkbox" v-model="isPrior"/>
				<label class="red"></label>
			</div>
			<p class="comment"><strong v-if="isPrior">[已打开]</strong> 打开后可以覆盖父级或子级配置。</p>
		</td>
	</tr>
</tbody>`
})

Vue.component("http-charsets-box", {
	props: ["v-usual-charsets", "v-all-charsets", "v-charset-config", "v-is-location", "v-is-group"],
	data: function () {
		let charsetConfig = this.vCharsetConfig
		if (charsetConfig == null) {
			charsetConfig = {
				isPrior: false,
				isOn: false,
				charset: "",
				isUpper: false
			}
		}
		return {
			charsetConfig: charsetConfig,
			advancedVisible: false
		}
	},
	methods: {
		changeAdvancedVisible: function (v) {
			this.advancedVisible = v
		}
	},
	template: `<div>
	<input type="hidden" name="charsetJSON" :value="JSON.stringify(charsetConfig)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="charsetConfig" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || charsetConfig.isPrior">
			<tr>
				<td class="title">是否启用</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="charsetConfig.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="((!vIsLocation && !vIsGroup) || charsetConfig.isPrior) && charsetConfig.isOn">	
			<tr>
				<td class="title">选择字符编码</td>
				<td><select class="ui dropdown" style="width:20em" name="charset" v-model="charsetConfig.charset">
						<option value="">[未选择]</option>
						<optgroup label="常用字符编码"></optgroup>
						<option v-for="charset in vUsualCharsets" :value="charset.charset">{{charset.charset}}（{{charset.name}}）</option>
						<optgroup label="全部字符编码"></optgroup>
						<option v-for="charset in vAllCharsets" :value="charset.charset">{{charset.charset}}（{{charset.name}}）</option>
					</select>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-if="((!vIsLocation && !vIsGroup) || charsetConfig.isPrior) && charsetConfig.isOn"></more-options-tbody>
		<tbody v-show="((!vIsLocation && !vIsGroup) || charsetConfig.isPrior) && charsetConfig.isOn && advancedVisible">
			<tr>
				<td>字符编码是否大写</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="charsetConfig.isUpper"/>
						<label></label>
					</div>
					<p class="comment">选中后将指定的字符编码转换为大写，比如默认为<span class="ui label tiny">utf-8</span>，选中后将改为<span class="ui label tiny">UTF-8</span>。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("http-access-log-box", {
	props: ["v-access-log", "v-keyword", "v-show-server-link"],
	data: function () {
		let accessLog = this.vAccessLog
		if (accessLog.header != null && accessLog.header.Upgrade != null && accessLog.header.Upgrade.values != null && accessLog.header.Upgrade.values.$contains("websocket")) {
			if (accessLog.scheme == "http") {
				accessLog.scheme = "ws"
			} else if (accessLog.scheme == "https") {
				accessLog.scheme = "wss"
			}
		}

		return {
			accessLog: accessLog
		}
	},
	methods: {
		formatCost: function (seconds) {
			var s = (seconds * 1000).toString();
			var pieces = s.split(".");
			if (pieces.length < 2) {
				return s;
			}

			return pieces[0] + "." + pieces[1].substr(0, 3);
		},
		showLog: function () {
			let that = this
			let requestId = this.accessLog.requestId
			this.$parent.$children.forEach(function (v) {
				if (v.deselect != null) {
					v.deselect()
				}
			})
			this.select()
			teaweb.popup("/servers/server/log/viewPopup?requestId=" + requestId, {
				width: "50em",
				height: "28em",
				onClose: function () {
					that.deselect()
				}
			})
		},
		select: function () {
			this.$refs.box.parentNode.style.cssText = "background: rgba(0, 0, 0, 0.1)"
		},
		deselect: function () {
			this.$refs.box.parentNode.style.cssText = ""
		}
	},
	template: `<div style="word-break: break-all" :style="{'color': (accessLog.status >= 400) ? '#dc143c' : ''}" ref="box">
	<a v-if="accessLog.node != null && accessLog.node.nodeCluster != null" :href="'/clusters/cluster/node?nodeId=' + accessLog.node.id + '&clusterId=' + accessLog.node.nodeCluster.id" title="点击查看节点详情" target="_top"><span class="grey">[{{accessLog.node.name}}<span v-if="!accessLog.node.name.endsWith('节点')">节点</span>]</span></a>
	<a :href="'/servers/server/log?serverId=' + accessLog.serverId" title="点击到网站服务" v-if="vShowServerLink"><span class="grey">[服务]</span></a>
	<span v-if="accessLog.region != null && accessLog.region.length > 0" class="grey">[{{accessLog.region}}]</span> <ip-box><keyword :v-word="vKeyword">{{accessLog.remoteAddr}}</keyword></ip-box> [{{accessLog.timeLocal}}] <em>&quot;<keyword :v-word="vKeyword">{{accessLog.requestMethod}}</keyword> {{accessLog.scheme}}://<keyword :v-word="vKeyword">{{accessLog.host}}</keyword><keyword :v-word="vKeyword">{{accessLog.requestURI}}</keyword> <a :href="accessLog.scheme + '://' + accessLog.host + accessLog.requestURI" target="_blank" title="新窗口打开" class="disabled"><i class="external icon tiny"></i> </a> {{accessLog.proto}}&quot; </em> <keyword :v-word="vKeyword">{{accessLog.status}}</keyword> <code-label v-if="accessLog.attrs != null && accessLog.attrs['cache.status'] == 'HIT'">cache hit</code-label> <code-label v-if="accessLog.firewallActions != null && accessLog.firewallActions.length > 0">waf {{accessLog.firewallActions}}</code-label> <span v-if="accessLog.tags != null && accessLog.tags.length > 0">- <code-label v-for="tag in accessLog.tags" :key="tag">{{tag}}</code-label></span> - 耗时:{{formatCost(accessLog.requestTime)}} ms <span v-if="accessLog.humanTime != null && accessLog.humanTime.length > 0" class="grey small">&nbsp; ({{accessLog.humanTime}})</span>
	&nbsp; <a href="" @click.prevent="showLog" title="查看详情"><i class="icon expand"></i></a>
</div>`
})

Vue.component("http-access-log-config-box", {
	props: ["v-access-log-config", "v-fields", "v-default-field-codes", "v-is-location", "v-is-group"],
	data: function () {
		let that = this

		// 初始化
		setTimeout(function () {
			that.changeFields()
		}, 100)

		let accessLog = {
			isPrior: false,
			isOn: false,
			fields: [],
			status1: true,
			status2: true,
			status3: true,
			status4: true,
			status5: true,

            firewallOnly: false
		}
		if (this.vAccessLogConfig != null) {
			accessLog = this.vAccessLogConfig
		}

		this.vFields.forEach(function (v) {
			if (that.vAccessLogConfig == null) { // 初始化默认值
				v.isChecked = that.vDefaultFieldCodes.$contains(v.code)
			} else {
				v.isChecked = accessLog.fields.$contains(v.code)
			}
		})

		return {
			accessLog: accessLog
		}
	},
	methods: {
		changeFields: function () {
			this.accessLog.fields = this.vFields.filter(function (v) {
				return v.isChecked
			}).map(function (v) {
				return v.code
			})
		}
	},
	template: `<div>
	<input type="hidden" name="accessLogJSON" :value="JSON.stringify(accessLog)"/>
	<table class="ui table definition selectable" :class="{'opacity-mask': this.accessLog.firewallOnly}">
		<prior-checkbox :v-config="accessLog" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || accessLog.isPrior">
			<tr>
				<td class="title">是否开启访问日志存储</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="accessLog.isOn"/>
						<label></label>
					</div>
					<p class="comment">关闭访问日志，并不影响统计的运行。</p>
				</td>
			</tr>
		</tbody>
		<tbody  v-show="((!vIsLocation && !vIsGroup) || accessLog.isPrior) && accessLog.isOn">
			<tr>
				<td>要存储的访问日志字段</td>
				<td>
					<div class="ui checkbox" v-for="field in vFields" style="width:10em;margin-bottom:0.8em">
						<input type="checkbox" v-model="field.isChecked" @change="changeFields"/>
						<label>{{field.name}}</label>
					</div>
				</td>
			</tr>
			<tr>
				<td>要存储的访问日志状态码</td>
				<td>
					<div class="ui checkbox" style="width:3.5em">
						<input type="checkbox" v-model="accessLog.status1"/>
						<label>1xx</label>
					</div>
					<div class="ui checkbox" style="width:3.5em">
						<input type="checkbox" v-model="accessLog.status2"/>
						<label>2xx</label>
					</div>
					<div class="ui checkbox" style="width:3.5em">
						<input type="checkbox" v-model="accessLog.status3"/>
						<label>3xx</label>
					</div>
					<div class="ui checkbox" style="width:3.5em">
						<input type="checkbox" v-model="accessLog.status4"/>
						<label>4xx</label>
					</div>
					<div class="ui checkbox" style="width:3.5em">
						<input type="checkbox" v-model="accessLog.status5"/>
						<label>5xx</label>
					</div>
				</td>
			</tr>
		</tbody>
	</table>
	
	<div v-show="((!vIsLocation && !vIsGroup) || accessLog.isPrior) && accessLog.isOn">
        <h4>WAF相关</h4>
        <table class="ui table definition selectable">
            <tr>
                <td class="title">是否只记录WAF相关日志</td>
                <td>
                    <checkbox v-model="accessLog.firewallOnly"></checkbox>
                    <p class="comment">选中后只记录WAF相关的日志。通过此选项可有效减少访问日志数量，降低网络带宽和存储压力。</p>
                </td>
            </tr>
        </table>
    </div>
	<div class="margin"></div>
</div>`
})

// 显示流量限制说明
Vue.component("traffic-limit-view", {
	props: ["v-traffic-limit"],
	data: function () {
		return {
			config: this.vTrafficLimit
		}
	},
	template: `<div>
	<div v-if="config.isOn">
		<span v-if="config.dailySize != null && config.dailySize.count > 0">日流量限制：{{config.dailySize.count}}{{config.dailySize.unit.toUpperCase()}}<br/></span>
		<span v-if="config.monthlySize != null && config.monthlySize.count > 0">月流量限制：{{config.monthlySize.count}}{{config.monthlySize.unit.toUpperCase()}}<br/></span>
	</div>
	<span v-else class="disabled">没有限制。</span>
</div>`
})

// 基本认证用户配置
Vue.component("http-auth-basic-auth-user-box", {
	props: ["v-users"],
	data: function () {
		let users = this.vUsers
		if (users == null) {
			users = []
		}
		return {
			users: users,
			isAdding: false,
			updatingIndex: -1,

			username: "",
			password: ""
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			this.username = ""
			this.password = ""

			let that = this
			setTimeout(function () {
				that.$refs.username.focus()
			}, 100)
		},
		cancel: function () {
			this.isAdding = false
			this.updatingIndex = -1
		},
		confirm: function () {
			let that = this
			if (this.username.length == 0) {
				teaweb.warn("请输入用户名", function () {
					that.$refs.username.focus()
				})
				return
			}
			if (this.password.length == 0) {
				teaweb.warn("请输入密码", function () {
					that.$refs.password.focus()
				})
				return
			}
			if (this.updatingIndex < 0) {
				this.users.push({
					username: this.username,
					password: this.password
				})
			} else {
				this.users[this.updatingIndex].username = this.username
				this.users[this.updatingIndex].password = this.password
			}
			this.cancel()
		},
		update: function (index, user) {
			this.updatingIndex = index

			this.isAdding = true
			this.username = user.username
			this.password = user.password

			let that = this
			setTimeout(function () {
				that.$refs.username.focus()
			}, 100)
		},
		remove: function (index) {
			this.users.$remove(index)
		}
	},
	template: `<div>
	<input type="hidden" name="httpAuthBasicAuthUsersJSON" :value="JSON.stringify(users)"/>
	<div v-if="users.length > 0">
		<div class="ui label small basic" v-for="(user, index) in users">
			{{user.username}} <a href="" title="修改" @click.prevent="update(index, user)"><i class="icon pencil tiny"></i></a>
			<a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div v-show="isAdding">
		<div class="ui fields inline">
			<div class="ui field">
				<input type="text" placeholder="用户名" v-model="username" size="15" ref="username"/>
			</div>
			<div class="ui field">
				<input type="password" placeholder="密码" v-model="password" size="15" ref="password"/>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>&nbsp;
				<a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
	</div>
	<div v-if="!isAdding" style="margin-top: 1em">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})

Vue.component("http-location-labels", {
	props: ["v-location-config", "v-server-id"],
	data: function () {
		return {
			location: this.vLocationConfig
		}
	},
	methods: {
		// 判断是否已启用某配置
		configIsOn: function (config) {
			return config != null && config.isPrior && config.isOn
		},

		refIsOn: function (ref, config) {
			return this.configIsOn(ref) && config != null && config.isOn
		},

		len: function (arr) {
			return (arr == null) ? 0 : arr.length
		},
		url: function (path) {
			return "/servers/server/settings/locations" + path + "?serverId=" + this.vServerId + "&locationId=" + this.location.id
		}
	},
	template: `	<div class="labels-box">
	<!-- 基本信息 -->
	<http-location-labels-label v-if="location.name != null && location.name.length > 0" :class="'olive'" :href="url('/location')">{{location.name}}</http-location-labels-label>
	<http-location-labels-label v-if="location.isBreak" :href="url('/location')">BREAK</http-location-labels-label>
	
	<!-- redirectToHTTPS -->
	<http-location-labels-label v-if="location.web != null && configIsOn(location.web.redirectToHTTPS)" :href="url('/http')">自动跳转HTTPS</http-location-labels-label>
	
	<!-- Web -->
	<http-location-labels-label v-if="location.web != null && configIsOn(location.web.root)" :href="url('/web')">文档根目录</http-location-labels-label>
	
	<!-- 反向代理 -->
	<http-location-labels-label v-if="refIsOn(location.reverseProxyRef, location.reverseProxy)" :v-href="url('/reverseProxy')">反向代理</http-location-labels-label>
	
	<!-- WAF -->
	<!-- TODO -->
	
	<!-- Cache -->
	<http-location-labels-label v-if="location.web != null && configIsOn(location.web.cache)" :v-href="url('/cache')">CACHE</http-location-labels-label>
	
	<!-- Charset -->
	<http-location-labels-label v-if="location.web != null && configIsOn(location.web.charset) && location.web.charset.charset.length > 0" :href="url('/charset')">{{location.web.charset.charset}}</http-location-labels-label>
	
	<!-- 访问日志 -->
	<!-- TODO -->
	
	<!-- 统计 -->
	<!-- TODO -->
	
	<!-- Gzip -->
	<http-location-labels-label v-if="location.web != null && refIsOn(location.web.gzipRef, location.web.gzip) && location.web.gzip.level > 0" :href="url('/gzip')">Gzip:{{location.web.gzip.level}}</http-location-labels-label>
	
	<!-- HTTP Header -->
	<http-location-labels-label v-if="location.web != null && refIsOn(location.web.requestHeaderPolicyRef, location.web.requestHeaderPolicy) && (len(location.web.requestHeaderPolicy.addHeaders) > 0 || len(location.web.requestHeaderPolicy.setHeaders) > 0 || len(location.web.requestHeaderPolicy.replaceHeaders) > 0 || len(location.web.requestHeaderPolicy.deleteHeaders) > 0)" :href="url('/headers')">请求Header</http-location-labels-label>
	<http-location-labels-label v-if="location.web != null && refIsOn(location.web.responseHeaderPolicyRef, location.web.responseHeaderPolicy) && (len(location.web.responseHeaderPolicy.addHeaders) > 0 || len(location.web.responseHeaderPolicy.setHeaders) > 0 || len(location.web.responseHeaderPolicy.replaceHeaders) > 0 || len(location.web.responseHeaderPolicy.deleteHeaders) > 0)" :href="url('/headers')">响应Header</http-location-labels-label>
	
	<!-- Websocket -->
	<http-location-labels-label v-if="location.web != null && refIsOn(location.web.websocketRef, location.web.websocket)" :href="url('/websocket')">Websocket</http-location-labels-label>
	
	<!-- 特殊页面 -->
	<div v-if="location.web != null && location.web.pages != null && location.web.pages.length > 0">
		<div v-for="page in location.web.pages" :key="page.id"><http-location-labels-label :href="url('/pages')">PAGE [状态码{{page.status[0]}}] -&gt; {{page.url}}</http-location-labels-label></div>
	</div>
	<div v-if="location.web != null && configIsOn(location.web.shutdown)">
		<http-location-labels-label :v-class="'red'" :href="url('/pages')">临时关闭</http-location-labels-label>
	</div>
	
	<!-- 重写规则 -->
	<div v-if="location.web != null && location.web.rewriteRules != null && location.web.rewriteRules.length > 0">
		<div v-for="rewriteRule in location.web.rewriteRules">
			<http-location-labels-label :href="url('/rewrite')">REWRITE {{rewriteRule.pattern}} -&gt; {{rewriteRule.replace}}</http-location-labels-label>
		</div>
	</div>
</div>`
})

Vue.component("http-location-labels-label", {
	props: ["v-class", "v-href"],
	template: `<a :href="vHref" class="ui label tiny" :class="vClass" style="font-size:0.7em;padding:4px;margin-top:0.3em;margin-bottom:0.3em"><slot></slot></a>`
})

Vue.component("http-gzip-box", {
	props: ["v-gzip-config", "v-gzip-ref", "v-is-location"],
	data: function () {
		let gzip = this.vGzipConfig
		if (gzip == null) {
			gzip = {
				isOn: true,
				level: 0,
				minLength: null,
				maxLength: null,
				conds: null
			}
		}

		return {
			gzip: gzip,
			advancedVisible: false
		}
	},
	methods: {
		isOn: function () {
			return (!this.vIsLocation || this.vGzipRef.isPrior) && this.vGzipRef.isOn
		},
		changeAdvancedVisible: function (v) {
			this.advancedVisible = v
		}
	},
	template: `<div>
<input type="hidden" name="gzipRefJSON" :value="JSON.stringify(vGzipRef)"/> 
<table class="ui table selectable definition">
	<prior-checkbox :v-config="vGzipRef" v-if="vIsLocation"></prior-checkbox>
	<tbody v-show="!vIsLocation || vGzipRef.isPrior">
		<tr>
			<td class="title">启用Gzip压缩</td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" v-model="vGzipRef.isOn"/>
					<label></label>
				</div>
			</td>
		</tr>
	</tbody>
	<tbody v-show="isOn()">
		<tr>
			<td class="title">压缩级别</td>
			<td>
				<select class="dropdown auto-width" name="level" v-model="gzip.level">
					<option value="0">不压缩</option>
					<option v-for="i in 9" :value="i">{{i}}</option>
				</select>
				<p class="comment">级别越高，压缩比例越大。</p>
			</td>
		</tr>
	</tbody>
	<more-options-tbody @change="changeAdvancedVisible" v-if="isOn()"></more-options-tbody>
	<tbody v-show="isOn() && advancedVisible">
		<tr>
			<td>Gzip内容最小长度</td>
			<td>
				<size-capacity-box :v-name="'minLength'" :v-value="gzip.minLength" :v-unit="'kb'"></size-capacity-box>
				<p class="comment">0表示不限制，内容长度从文件尺寸或Content-Length中获取。</p>
			</td>
		</tr>
		<tr>
			<td>Gzip内容最大长度</td>
			<td>
				<size-capacity-box :v-name="'maxLength'" :v-value="gzip.maxLength" :v-unit="'mb'"></size-capacity-box>
				<p class="comment">0表示不限制，内容长度从文件尺寸或Content-Length中获取。</p>
			</td>
		</tr>
		<tr>
			<td>匹配条件</td>
			<td>
				<http-request-conds-box :v-conds="gzip.conds"></http-request-conds-box>
</td>
		</tr>
	</tbody>
</table>
</div>`
})

Vue.component("ssl-certs-view", {
	props: ["v-certs"],
	data: function () {
		let certs = this.vCerts
		if (certs == null) {
			certs = []
		}
		return {
			certs: certs
		}
	},
	methods: {
		// 格式化时间
		formatTime: function (timestamp) {
			return new Date(timestamp * 1000).format("Y-m-d")
		},

		// 查看详情
		viewCert: function (certId) {
			teaweb.popup("/servers/certs/certPopup?certId=" + certId, {
				height: "28em",
				width: "48em"
			})
		}
	},
	template: `<div>
	<div v-if="certs != null && certs.length > 0">
		<div class="ui label small" v-for="(cert, index) in certs">
			{{cert.name}} / {{cert.dnsNames}} / 有效至{{formatTime(cert.timeEndAt)}} &nbsp;<a href="" title="查看" @click.prevent="viewCert(cert.id)"><i class="icon external alternate"></i></a>
		</div>
	</div>
</div>`
})

Vue.component("reverse-proxy-box", {
	props: ["v-reverse-proxy-ref", "v-reverse-proxy-config", "v-is-location", "v-is-group", "v-family"],
	data: function () {
		let reverseProxyRef = this.vReverseProxyRef
		if (reverseProxyRef == null) {
			reverseProxyRef = {
				isPrior: false,
				isOn: false,
				reverseProxyId: 0
			}
		}

		let reverseProxyConfig = this.vReverseProxyConfig
		if (reverseProxyConfig == null) {
			reverseProxyConfig = {
				requestPath: "",
				stripPrefix: "",
				requestURI: "",
				requestHost: "",
				requestHostType: 0,
				addHeaders: [],
				connTimeout: {count: 0, unit: "second"},
				readTimeout: {count: 0, unit: "second"},
				idleTimeout: {count: 0, unit: "second"},
				maxConns: 0,
				maxIdleConns: 0
			}
		}
		if (reverseProxyConfig.addHeaders == null) {
			reverseProxyConfig.addHeaders = []
		}
		if (reverseProxyConfig.connTimeout == null) {
			reverseProxyConfig.connTimeout = {count: 0, unit: "second"}
		}
		if (reverseProxyConfig.readTimeout == null) {
			reverseProxyConfig.readTimeout = {count: 0, unit: "second"}
		}
		if (reverseProxyConfig.idleTimeout == null) {
			reverseProxyConfig.idleTimeout = {count: 0, unit: "second"}
		}

		if (reverseProxyConfig.proxyProtocol == null) {
			// 如果直接赋值Vue将不会触发变更通知
			Vue.set(reverseProxyConfig, "proxyProtocol", {
				isOn: false,
				version: 1
			})
		}

		let forwardHeaders = [
			{
				name: "X-Real-IP",
				isChecked: false
			},
			{
				name: "X-Forwarded-For",
				isChecked: false
			},
			{
				name: "X-Forwarded-By",
				isChecked: false
			},
			{
				name: "X-Forwarded-Host",
				isChecked: false
			},
			{
				name: "X-Forwarded-Proto",
				isChecked: false
			}
		]
		forwardHeaders.forEach(function (v) {
			v.isChecked = reverseProxyConfig.addHeaders.$contains(v.name)
		})

		return {
			reverseProxyRef: reverseProxyRef,
			reverseProxyConfig: reverseProxyConfig,
			advancedVisible: false,
			family: this.vFamily,
			forwardHeaders: forwardHeaders
		}
	},
	watch: {
		"reverseProxyConfig.requestHostType": function (v) {
			let requestHostType = parseInt(v)
			if (isNaN(requestHostType)) {
				requestHostType = 0
			}
			this.reverseProxyConfig.requestHostType = requestHostType
		},
		"reverseProxyConfig.connTimeout.count": function (v) {
			let count = parseInt(v)
			if (isNaN(count) || count < 0) {
				count = 0
			}
			this.reverseProxyConfig.connTimeout.count = count
		},
		"reverseProxyConfig.readTimeout.count": function (v) {
			let count = parseInt(v)
			if (isNaN(count) || count < 0) {
				count = 0
			}
			this.reverseProxyConfig.readTimeout.count = count
		},
		"reverseProxyConfig.idleTimeout.count": function (v) {
			let count = parseInt(v)
			if (isNaN(count) || count < 0) {
				count = 0
			}
			this.reverseProxyConfig.idleTimeout.count = count
		},
		"reverseProxyConfig.maxConns": function (v) {
			let maxConns = parseInt(v)
			if (isNaN(maxConns) || maxConns < 0) {
				maxConns = 0
			}
			this.reverseProxyConfig.maxConns = maxConns
		},
		"reverseProxyConfig.maxIdleConns": function (v) {
			let maxIdleConns = parseInt(v)
			if (isNaN(maxIdleConns) || maxIdleConns < 0) {
				maxIdleConns = 0
			}
			this.reverseProxyConfig.maxIdleConns = maxIdleConns
		},
		"reverseProxyConfig.proxyProtocol.version": function (v) {
			let version = parseInt(v)
			if (isNaN(version)) {
				version = 1
			}
			this.reverseProxyConfig.proxyProtocol.version = version
		}
	},
	methods: {
		isOn: function () {
			if (this.vIsLocation || this.vIsGroup) {
				return this.reverseProxyRef.isPrior && this.reverseProxyRef.isOn
			}
			return this.reverseProxyRef.isOn
		},
		changeAdvancedVisible: function (v) {
			this.advancedVisible = v
		},
		changeAddHeader: function () {
			this.reverseProxyConfig.addHeaders = this.forwardHeaders.filter(function (v) {
				return v.isChecked
			}).map(function (v) {
				return v.name
			})
		}
	},
	template: `<div>
	<input type="hidden" name="reverseProxyRefJSON" :value="JSON.stringify(reverseProxyRef)"/>
	<input type="hidden" name="reverseProxyJSON" :value="JSON.stringify(reverseProxyConfig)"/>
	<table class="ui table selectable definition">
		<prior-checkbox :v-config="reverseProxyRef" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || reverseProxyRef.isPrior">
			<tr>
				<td class="title">是否启用反向代理</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="reverseProxyRef.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
			<tr v-show="family == null || family == 'http'">
				<td>回源主机名<em>（Host）</em></td>
				<td>	
					<radio :v-value="0" v-model="reverseProxyConfig.requestHostType">跟随代理服务</radio> &nbsp;
					<radio :v-value="1" v-model="reverseProxyConfig.requestHostType">跟随源站</radio> &nbsp;
					<radio :v-value="2" v-model="reverseProxyConfig.requestHostType">自定义</radio>
					<div v-show="reverseProxyConfig.requestHostType == 2" style="margin-top: 0.8em">
						<input type="text" placeholder="比如example.com" v-model="reverseProxyConfig.requestHost"/>
					</div>
					<p class="comment">请求源站时的Host，用于修改源站接收到的域名
					<span v-if="reverseProxyConfig.requestHostType == 0">，"跟随代理服务"是指源站接收到的域名和当前代理服务保持一致</span>
					<span v-if="reverseProxyConfig.requestHostType == 1">，"跟随源站"是指源站接收到的域名仍然是填写的源站地址中的信息，不随代理服务域名改变而改变</span>					
					<span v-if="reverseProxyConfig.requestHostType == 2">，自定义Host内容中支持请求变量</span>。</p>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-if="isOn()"></more-options-tbody>
		<tbody v-show="isOn() && advancedVisible">
		    <tr v-show="family == null || family == 'http'">
		        <td>自动添加的Header</td>
		        <td>
		            <div>
		                <div style="width: 14em; float: left; margin-bottom: 1em" v-for="header in forwardHeaders" :key="header.name">
		                    <checkbox v-model="header.isChecked" @input="changeAddHeader">{{header.name}}</checkbox>
                        </div>
                        <div style="clear: both"></div>
                    </div>
                    <p class="comment">选中后，会自动向源站请求添加这些Header。</p>
                </td> 
            </tr>
			<tr v-show="family == null || family == 'http'">
				<td>请求URI<em>（RequestURI）</em></td>
				<td>
					<input type="text" placeholder="\${requestURI}" v-model="reverseProxyConfig.requestURI"/>
					<p class="comment">\${requestURI}为完整的请求URI，可以使用类似于"\${requestURI}?arg1=value1&arg2=value2"的形式添加你的参数。</p>
				</td>
			</tr>
			<tr v-show="family == null || family == 'http'">
				<td>去除URL前缀<em>（StripPrefix）</em></td>
				<td>
					<input type="text" v-model="reverseProxyConfig.stripPrefix" placeholder="/PREFIX"/>
					<p class="comment">可以把请求的路径部分前缀去除后再查找文件，比如把 <span class="ui label tiny">/web/app/index.html</span> 去除前缀 <span class="ui label tiny">/web</span> 后就变成 <span class="ui label tiny">/app/index.html</span>。 </p>
				</td>
			</tr>
			<tr v-if="family == null || family == 'http'">
				<td>是否自动刷新缓存区<em>（AutoFlush）</em></td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="reverseProxyConfig.autoFlush"/>
						<label></label>
					</div>
					<p class="comment">开启后将自动刷新缓冲区数据到客户端，在类似于SSE（server-sent events）等场景下很有用。</p>
				</td>
			</tr>
			<tr v-if="family == null || family == 'http'">
                <td class="color-border">源站连接失败超时时间</td>
                <td>
                    <div class="ui fields inline">
                        <div class="ui field">
                            <input type="text" name="connTimeout" value="10" size="6" v-model="reverseProxyConfig.connTimeout.count"/>
                        </div>
                        <div class="ui field">
                            秒
                        </div>
                    </div>
                    <p class="comment">连接源站失败的最大超时时间，0表示不限制。</p>
                </td>
            </tr>
            <tr v-if="family == null || family == 'http'">
                <td class="color-border">源站读取超时时间</td>
                <td>
                    <div class="ui fields inline">
                        <div class="ui field">
                            <input type="text" name="readTimeout" value="0" size="6" v-model="reverseProxyConfig.readTimeout.count"/>
                        </div>
                        <div class="ui field">
                            秒
                        </div>
                    </div>
                    <p class="comment">读取内容时的最大超时时间，0表示不限制。</p>
                </td>
            </tr>
            <tr v-if="family == null || family == 'http'">
                <td class="color-border">源站最大并发连接数</td>
                <td>
                    <div class="ui fields inline">
                        <div class="ui field">
                            <input type="text" name="maxConns" value="0" size="6" maxlength="10" v-model="reverseProxyConfig.maxConns"/>
                        </div>
                    </div>
                    <p class="comment">源站可以接受到的最大并发连接数，0表示使用系统默认。</p>
                </td>
            </tr>
            <tr v-if="family == null || family == 'http'">
                <td class="color-border">源站最大空闲连接数</td>
                <td>
                    <div class="ui fields inline">
                        <div class="ui field">
                            <input type="text" name="maxIdleConns" value="0" size="6" maxlength="10" v-model="reverseProxyConfig.maxIdleConns"/>
                        </div>
                    </div>
                    <p class="comment">当没有请求时，源站保持等待的最大空闲连接数量，0表示使用系统默认。</p>
                </td>
            </tr>
            <tr v-if="family == null || family == 'http'">
                <td class="color-border">源站最大空闲超时时间</td>
                <td>
                    <div class="ui fields inline">
                        <div class="ui field">
                            <input type="text" name="idleTimeout" value="0" size="6" v-model="reverseProxyConfig.idleTimeout.count"/>
                        </div>
                        <div class="ui field">
                            秒
                        </div>
                    </div>
                    <p class="comment">源站保持等待的空闲超时时间，0表示使用默认时间。</p>
                </td>
            </tr>
            <tr v-show="family != 'unix'">
            	<td>PROXY Protocol</td>
            	<td>
            		<checkbox name="proxyProtocolIsOn" v-model="reverseProxyConfig.proxyProtocol.isOn"></checkbox>
            		<p class="comment">选中后表示启用PROXY Protocol，每次连接源站时都会在头部写入客户端地址信息。</p>
				</td>
			</tr>
			<tr v-show="family != 'unix' && reverseProxyConfig.proxyProtocol.isOn">
				<td>PROXY Protocol版本</td>
				<td>
					<select class="ui dropdown auto-width" name="proxyProtocolVersion" v-model="reverseProxyConfig.proxyProtocol.version">
						<option value="1">1</option>
						<option value="2">2</option>
					</select>
					<p class="comment" v-if="reverseProxyConfig.proxyProtocol.version == 1">发送类似于<code-label>PROXY TCP4 192.168.1.1 192.168.1.10 32567 443</code-label>的头部信息。</p>
					<p class="comment" v-if="reverseProxyConfig.proxyProtocol.version == 2">发送二进制格式的头部信息。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("http-firewall-param-filters-box", {
	props: ["v-filters"],
	data: function () {
		let filters = this.vFilters
		if (filters == null) {
			filters = []
		}

		return {
			filters: filters,
			isAdding: false,
			options: [
				{name: "MD5", code: "md5"},
				{name: "URLEncode", code: "urlEncode"},
				{name: "URLDecode", code: "urlDecode"},
				{name: "BASE64Encode", code: "base64Encode"},
				{name: "BASE64Decode", code: "base64Decode"},
				{name: "UNICODE编码", code: "unicodeEncode"},
				{name: "UNICODE解码", code: "unicodeDecode"},
				{name: "HTML实体编码", code: "htmlEscape"},
				{name: "HTML实体解码", code: "htmlUnescape"},
				{name: "计算长度", code: "length"},
				{name: "十六进制->十进制", "code": "hex2dec"},
				{name: "十进制->十六进制", "code": "dec2hex"},
				{name: "SHA1", "code": "sha1"},
				{name: "SHA256", "code": "sha256"}
			],
			addingCode: ""
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			this.addingCode = ""
		},
		confirm: function () {
			if (this.addingCode.length == 0) {
				return
			}
			let that = this
			this.filters.push(this.options.$find(function (k, v) {
				return (v.code == that.addingCode)
			}))
			this.isAdding = false
		},
		cancel: function () {
			this.isAdding = false
		},
		remove: function (index) {
			this.filters.$remove(index)
		}
	},
	template: `<div>
		<input type="hidden" name="paramFiltersJSON" :value="JSON.stringify(filters)" />
		<div v-if="filters.length > 0">
			<div v-for="(filter, index) in filters" class="ui label small basic">
				{{filter.name}} <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove"></i></a>
			</div>
			<div class="ui divider"></div>
		</div>
		<div v-if="isAdding">
			<div class="ui fields inline">
				<div class="ui field">
					<select class="ui dropdown auto-width" v-model="addingCode">
						<option value="">[请选择]</option>
						<option v-for="option in options" :value="option.code">{{option.name}}</option>
					</select>
				</div>
				<div class="ui field">
					<button class="ui button tiny" type="button" @click.prevent="confirm()">确定</button>
					&nbsp; <a href="" @click.prevent="cancel()" title="取消"><i class="icon remove"></i></a>
				</div>
			</div>
		</div>
		<div v-if="!isAdding">
			<button class="ui button tiny" type="button" @click.prevent="add">+</button>
		</div>
		<p class="comment">可以对参数值进行特定的编解码处理。</p>
</div>`
})

Vue.component("http-remote-addr-config-box", {
	props: ["v-remote-addr-config", "v-is-location", "v-is-group"],
	data: function () {
		let config = this.vRemoteAddrConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				value: "${rawRemoteAddr}",
				isCustomized: false
			}
		}

		let optionValue = ""
		if (!config.isCustomized && (config.value == "${remoteAddr}" || config.value == "${rawRemoteAddr}")) {
			optionValue = config.value
		}

		return {
			config: config,
			options: [
				{
					name: "直接获取",
					description: "用户直接访问边缘节点，即 \"用户 --> 边缘节点\" 模式，这时候可以直接从连接中读取到真实的IP地址。",
					value: "${rawRemoteAddr}"
				},
				{
					name: "从上级代理中获取",
					description: "用户和边缘节点之间有别的代理服务转发，即 \"用户 --> [第三方代理服务] --> 边缘节点\"，这时候只能从上级代理中获取传递的IP地址。",
					value: "${remoteAddr}"
				},
				{
					name: "[自定义]",
					description: "通过自定义变量来获取客户端真实的IP地址。",
					value: ""
				}
			],
			optionValue: optionValue
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.config.isPrior) && this.config.isOn
		},
		changeOptionValue: function () {
			if (this.optionValue.length > 0) {
				this.config.value = this.optionValue
				this.config.isCustomized = false
			} else {
				this.config.isCustomized = true
			}
		}
	},
	template: `<div>
	<input type="hidden" name="remoteAddrJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
			<tr>
				<td class="title">是否启用</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="config.isOn"/>
						<label></label>
					</div>
					<p class="comment">选中后表示使用自定义的请求变量获取客户端IP。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td>获取IP方式 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="optionValue" @change="changeOptionValue">
						<option v-for="option in options" :value="option.value">{{option.name}}</option>
					</select>
					<p class="comment" v-for="option in options" v-if="option.value == optionValue && option.description.length > 0">{{option.description}}</p>
				</td>
			</tr>
			<tr v-show="optionValue.length == 0">
				<td>读取IP变量值 *</td>
				<td>
					<input type="hidden" v-model="config.value" maxlength="100"/>
					<div v-if="optionValue == ''" style="margin-top: 1em">
						<input type="text" v-model="config.value" maxlength="100"/>
						<p class="comment">通过此变量获取用户的IP地址。具体可用的请求变量列表可参考官方网站文档。</p>
					</div>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>		
</div>`
})

// 访问日志搜索框
Vue.component("http-access-log-search-box", {
	props: ["v-ip", "v-domain", "v-keyword"],
	data: function () {
		let ip = this.vIp
		if (ip == null) {
			ip = ""
		}

		let domain = this.vDomain
		if (domain == null) {
			domain = ""
		}

		let keyword = this.vKeyword
		if (keyword == null) {
			keyword = ""
		}

		return {
			ip: ip,
			domain: domain,
			keyword: keyword
		}
	},
	methods: {
		cleanIP: function () {
			this.ip = ""
			this.submit()
		},
		cleanDomain: function () {
			this.domain = ""
			this.submit()
		},
		cleanKeyword: function () {
			this.keyword = ""
			this.submit()
		},
		submit: function () {
			let parent = this.$el.parentNode
			while (true) {
				if (parent == null) {
					break
				}
				if (parent.tagName == "FORM") {
					break
				}
				parent = parent.parentNode
			}
			if (parent != null) {
				setTimeout(function () {
					parent.submit()
				}, 500)
			}
		}
	},
	template: `<div>
	<div class="margin"></div>
	<div class="ui fields inline">
		<div class="ui field">
			<div class="ui input left right labeled small">
				<span class="ui label basic" style="font-weight: normal">IP</span>
				<input type="text" name="ip" placeholder="x.x.x.x" size="15" v-model="ip"/>
				<a class="ui label basic" :class="{disabled: ip.length == 0}" @click.prevent="cleanIP"><i class="icon remove small"></i></a>
			</div>
		</div>
		<div class="ui field">
			<div class="ui input left right labeled small" >
				<span class="ui label basic" style="font-weight: normal">域名</span>
				<input type="text" name="domain" placeholder="xxx.com" size="15" v-model="domain"/>
				<a class="ui label basic" :class="{disabled: domain.length == 0}" @click.prevent="cleanDomain"><i class="icon remove small"></i></a>
			</div>
		</div>
		<div class="ui field">
			<div class="ui input left right labeled small">
				<span class="ui label basic" style="font-weight: normal">关键词</span>
				<input type="text" name="keyword" v-model="keyword" placeholder="路径、UserAgent等..." size="18"/>
				<a class="ui label basic" :class="{disabled: keyword.length == 0}" @click.prevent="cleanKeyword"><i class="icon remove small"></i></a>
			</div>
		</div>
		<slot></slot>
		<div class="ui field">
			<button class="ui button small" type="submit">查找</button>
		</div>
	</div>
</div>`
})

// 显示指标对象名
Vue.component("metric-key-label", {
	props: ["v-key"],
	data: function () {
		return {
			keyDefs: window.METRIC_HTTP_KEYS
		}
	},
	methods: {
		keyName: function (key) {
			let that = this
			let subKey = ""
			let def = this.keyDefs.$find(function (k, v) {
				if (v.code == key) {
					return true
				}
				if (key.startsWith("${arg.") && v.code.startsWith("${arg.")) {
					subKey = that.getSubKey("arg.", key)
					return true
				}
				if (key.startsWith("${header.") && v.code.startsWith("${header.")) {
					subKey = that.getSubKey("header.", key)
					return true
				}
				if (key.startsWith("${cookie.") && v.code.startsWith("${cookie.")) {
					subKey = that.getSubKey("cookie.", key)
					return true
				}
				return false
			})
			if (def != null) {
				if (subKey.length > 0) {
					return def.name + ": " + subKey
				}
				return def.name
			}
			return key
		},
		getSubKey: function (prefix, key) {
			prefix = "${" + prefix
			let index = key.indexOf(prefix)
			if (index >= 0) {
				key = key.substring(index + prefix.length)
				key = key.substring(0, key.length - 1)
				return key
			}
			return ""
		}
	},
	template: `<div class="ui label basic small">
	{{keyName(this.vKey)}}
</div>`
})

// 指标对象
Vue.component("metric-keys-config-box", {
	props: ["v-keys"],
	data: function () {
		let keys = this.vKeys
		if (keys == null) {
			keys = []
		}
		return {
			keys: keys,
			isAdding: false,
			key: "",
			subKey: "",
			keyDescription: "",

			keyDefs: window.METRIC_HTTP_KEYS
		}
	},
	watch: {
		keys: function () {
			this.$emit("change", this.keys)
		}
	},
	methods: {
		cancel: function () {
			this.key = ""
			this.subKey = ""
			this.keyDescription = ""
			this.isAdding = false
		},
		confirm: function () {
			if (this.key.length == 0) {
				return
			}

			if (this.key.indexOf(".NAME") > 0) {
				if (this.subKey.length == 0) {
					teaweb.warn("请输入参数值")
					return
				}
				this.key = this.key.replace(".NAME", "." + this.subKey)
			}
			this.keys.push(this.key)
			this.cancel()
		},
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				if (that.$refs.key != null) {
					that.$refs.key.focus()
				}
			}, 100)
		},
		remove: function (index) {
			this.keys.$remove(index)
		},
		changeKey: function () {
			if (this.key.length == 0) {
				return
			}
			let that = this
			let def = this.keyDefs.$find(function (k, v) {
				return v.code == that.key
			})
			if (def != null) {
				this.keyDescription = def.description
			}
		},
		keyName: function (key) {
			let that = this
			let subKey = ""
			let def = this.keyDefs.$find(function (k, v) {
				if (v.code == key) {
					return true
				}
				if (key.startsWith("${arg.") && v.code.startsWith("${arg.")) {
					subKey = that.getSubKey("arg.", key)
					return true
				}
				if (key.startsWith("${header.") && v.code.startsWith("${header.")) {
					subKey = that.getSubKey("header.", key)
					return true
				}
				if (key.startsWith("${cookie.") && v.code.startsWith("${cookie.")) {
					subKey = that.getSubKey("cookie.", key)
					return true
				}
				return false
			})
			if (def != null) {
				if (subKey.length > 0) {
					return def.name + ": " + subKey
				}
				return def.name
			}
			return key
		},
		getSubKey: function (prefix, key) {
			prefix = "${" + prefix
			let index = key.indexOf(prefix)
			if (index >= 0) {
				key = key.substring(index + prefix.length)
				key = key.substring(0, key.length - 1)
				return key
			}
			return ""
		}
	},
	template: `<div>
	<input type="hidden" name="keysJSON" :value="JSON.stringify(keys)"/>
	<div>
		<div v-for="(key, index) in keys" class="ui label small basic">
			{{keyName(key)}} &nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
	</div>
	<div v-if="isAdding" style="margin-top: 1em">
		<div class="ui fields inline">
			<div class="ui field">
				<select class="ui dropdown" v-model="key" @change="changeKey">
					<option value="">[选择对象]</option>
					<option v-for="def in keyDefs" :value="def.code">{{def.name}}</option>
				</select>
			</div>
			<div class="ui field" v-if="key == '\${arg.NAME}'">
				<input type="text" v-model="subKey" placeholder="参数名" size="15"/>
			</div>
			<div class="ui field" v-if="key == '\${header.NAME}'">
				<input type="text" v-model="subKey" placeholder="Header名" size="15">
			</div>
			<div class="ui field" v-if="key == '\${cookie.NAME}'">
				<input type="text" v-model="subKey" placeholder="Cookie名" size="15">
			</div>
			<div class="ui field">
				<button type="button" class="ui button tiny" @click.prevent="confirm">确定</button>
				<a href="" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
		<p class="comment" v-if="keyDescription.length > 0">{{keyDescription}}</p>
	</div>
	<div style="margin-top: 1em" v-if="!isAdding">
		<button type="button" class="ui button tiny" @click.prevent="add">+</button>
	</div>
</div>`
})

Vue.component("http-web-root-box", {
	props: ["v-root-config", "v-is-location", "v-is-group"],
	data: function () {
		let rootConfig = this.vRootConfig
		if (rootConfig == null) {
			rootConfig = {
				isPrior: false,
				isOn: true,
				dir: "",
				indexes: [],
				stripPrefix: "",
				decodePath: false,
				isBreak: false
			}
		}
		if (rootConfig.indexes == null) {
			rootConfig.indexes = []
		}
		return {
			rootConfig: rootConfig,
			advancedVisible: false
		}
	},
	methods: {
		changeAdvancedVisible: function (v) {
			this.advancedVisible = v
		},
		addIndex: function () {
			let that = this
			teaweb.popup("/servers/server/settings/web/createIndex", {
				height: "10em",
				callback: function (resp) {
					that.rootConfig.indexes.push(resp.data.index)
				}
			})
		},
		removeIndex: function (i) {
			this.rootConfig.indexes.$remove(i)
		},
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.rootConfig.isPrior) && this.rootConfig.isOn
		}
	},
	template: `<div>
	<input type="hidden" name="rootJSON" :value="JSON.stringify(rootConfig)"/>
	<table class="ui table selectable definition">
		<prior-checkbox :v-config="rootConfig" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || rootConfig.isPrior">
			<tr>
				<td class="title">是否开启静态资源分发</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="rootConfig.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td class="title">静态资源根目录</td>
				<td>
					<input type="text" name="root" v-model="rootConfig.dir" ref="focus" placeholder="类似于 /home/www"/>
					<p class="comment">可以访问此根目录下的静态资源。</p>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-if="isOn()"></more-options-tbody>

		<tbody v-show="isOn() && advancedVisible">
			<tr>
				<td>首页文件</td>
				<td>
					<!-- TODO 支持排序 -->
					<div v-if="rootConfig.indexes.length > 0">
						<div v-for="(index, i) in rootConfig.indexes" class="ui label tiny">
							{{index}} <a href="" title="删除" @click.prevent="removeIndex(i)"><i class="icon remove"></i></a>
						</div>
						<div class="ui divider"></div>
					</div>
					<button class="ui button tiny" type="button" @click.prevent="addIndex()">+</button>
					<p class="comment">在URL中只有目录没有文件名时默认查找的首页文件。</p>
				</td>
			</tr>
			<tr>
				<td>去除URL前缀</td>
				<td>
					<input type="text" v-model="rootConfig.stripPrefix" placeholder="/PREFIX"/>
					<p class="comment">可以把请求的路径部分前缀去除后再查找文件，比如把 <span class="ui label tiny">/web/app/index.html</span> 去除前缀 <span class="ui label tiny">/web</span> 后就变成 <span class="ui label tiny">/app/index.html</span>。 </p>
				</td>
			</tr>
			<tr>
				<td>路径解码</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="rootConfig.decodePath"/>
						<label></label>	
					</div>
					<p class="comment">是否对请求路径进行URL解码，比如把 <span class="ui label tiny">/Web+App+Browser.html</span> 解码成 <span class="ui label tiny">/Web App Browser.html</span> 再查找文件。</p>
				</td>
			</tr>
			<tr>
				<td>是否终止请求</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="rootConfig.isBreak"/>
						<label></label>	
					</div>
					<p class="comment">在找不到要访问的文件的情况下是否终止请求并返回404，如果选择终止请求，则不再尝试反向代理等设置。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

Vue.component("http-webp-config-box", {
	props: ["v-webp-config", "v-is-location", "v-is-group"],
	data: function () {
		let config = this.vWebpConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				quality: 50,
				minLength: {count: 0, "unit": "kb"},
				maxLength: {count: 0, "unit": "kb"},
				mimeTypes: ["image/png", "image/jpeg", "image/bmp", "image/x-ico", "image/gif"],
				extensions: [".png", ".jpeg", ".jpg", ".bmp", ".ico"],
				conds: null
			}
		}

		if (config.mimeTypes == null) {
			config.mimeTypes = []
		}
		if (config.extensions == null) {
			config.extensions = []
		}

		return {
			config: config,
			moreOptionsVisible: false,
			quality: config.quality
		}
	},
	watch: {
		quality: function (v) {
			let quality = parseInt(v)
			if (isNaN(quality)) {
				quality = 90
			} else if (quality < 1) {
				quality = 1
			} else if (quality > 100) {
				quality = 100
			}
			this.config.quality = quality
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.config.isPrior) && this.config.isOn
		},
		changeExtensions: function (values) {
			values.forEach(function (v, k) {
				if (v.length > 0 && v[0] != ".") {
					values[k] = "." + v
				}
			})
			this.config.extensions = values
		},
		changeMimeTypes: function (values) {
			this.config.mimeTypes = values
		},
		changeAdvancedVisible: function () {
			this.moreOptionsVisible = !this.moreOptionsVisible
		},
		changeConds: function (conds) {
			this.config.conds = conds
		}
	},
	template: `<div>
	<input type="hidden" name="webpJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
			<tr>
				<td class="title">是否启用</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="config.isOn"/>
						<label></label>
					</div>
					<p class="comment">选中后表示开启自动WebP压缩。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td>图片质量</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" v-model="quality" style="width: 5em" maxlength="4"/>
						<span class="ui label">%</span>
					</div>
					<p class="comment">取值在0到100之间，数值越大生成的图像越清晰，同时文件尺寸也会越大。</p>
				</td>
			</tr>
			<tr>
				<td>支持的扩展名</td>
				<td>
					<values-box :values="config.extensions" @change="changeExtensions" placeholder="比如 .html"></values-box>
					<p class="comment">含有这些扩展名的URL将会被转成WebP，不区分大小写。</p>
				</td>
			</tr>
			<tr>
				<td>支持的MimeType</td>
				<td>
					<values-box :values="config.mimeTypes" @change="changeMimeTypes" placeholder="比如 text/*"></values-box>
					<p class="comment">响应的Content-Type里包含这些MimeType的内容将会被转成WebP。</p>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-if="isOn()"></more-options-tbody>
		<tbody v-show="isOn() && moreOptionsVisible">
				<tr>
					<td>内容最小长度</td>
				<td>
					<size-capacity-box :v-name="'minLength'" :v-value="config.minLength" :v-unit="'kb'"></size-capacity-box>
					<p class="comment">0表示不限制，内容长度从文件尺寸或Content-Length中获取。</p>
				</td>
			</tr>
			<tr>
				<td>内容最大长度</td>
				<td>
					<size-capacity-box :v-name="'maxLength'" :v-value="config.maxLength" :v-unit="'mb'"></size-capacity-box>
					<p class="comment">0表示不限制，内容长度从文件尺寸或Content-Length中获取。</p>
				</td>
			</tr>
			<tr>
				<td>匹配条件</td>
				<td>
					<http-request-conds-box :v-conds="config.conds" @change="changeConds"></http-request-conds-box>
	</td>
			</tr>
		</tbody>
	</table>			
	<div class="ui margin"></div>
</div>`
})

Vue.component("origin-scheduling-view-box", {
	props: ["v-scheduling", "v-params"],
	data: function () {
		let scheduling = this.vScheduling
		if (scheduling == null) {
			scheduling = {}
		}
		return {
			scheduling: scheduling
		}
	},
	methods: {
		update: function () {
			teaweb.popup("/servers/server/settings/reverseProxy/updateSchedulingPopup?" + this.vParams, {
				height: "21em",
				callback: function () {
					window.location.reload()
				},
			})
		}
	},
	template: `<div>
	<div class="margin"></div>
	<table class="ui table selectable definition">
		<tr>
			<td class="title">当前正在使用的算法</td>
			<td>
				{{scheduling.name}} &nbsp; <a href="" @click.prevent="update()"><span>[修改]</span></a>
				<p class="comment">{{scheduling.description}}</p>
			</td>
		</tr>
	</table>
</div>`
})

Vue.component("http-firewall-block-options", {
	props: ["v-block-options"],
	data: function () {
		return {
			blockOptions: this.vBlockOptions,
			statusCode: this.vBlockOptions.statusCode,
			timeout: this.vBlockOptions.timeout
		}
	},
	watch: {
		statusCode: function (v) {
			let statusCode = parseInt(v)
			if (isNaN(statusCode)) {
				this.blockOptions.statusCode = 403
			} else {
				this.blockOptions.statusCode = statusCode
			}
		},
		timeout: function (v) {
			let timeout = parseInt(v)
			if (isNaN(timeout)) {
				this.blockOptions.timeout = 0
			} else {
				this.blockOptions.timeout = timeout
			}
		}
	},
	template: `<div>
<input type="hidden" name="blockOptionsJSON" :value="JSON.stringify(blockOptions)"/>
	<table class="ui table">
		<tr>
			<td class="title">状态码</td>
			<td>
				<input type="text" v-model="statusCode" style="width:4.5em" maxlength="3"/>
			</td>
		</tr>
		<tr>
			<td>提示内容</td>
			<td>
				<textarea rows="3" v-model="blockOptions.body"></textarea>
			</td>
		</tr>
		<tr>
			<td>超时时间</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" v-model="timeout" style="width: 5em" maxlength="6"/>
					<span class="ui label">秒</span>
				</div>
				<p class="comment">触发阻止动作时，封锁客户端IP的时间。</p>
			</td>
		</tr>
	</table>
</div>	
`
})

Vue.component("http-firewall-rules-box", {
	props: ["v-rules", "v-type"],
	data: function () {
		let rules = this.vRules
		if (rules == null) {
			rules = []
		}
		return {
			rules: rules
		}
	},
	methods: {
		addRule: function () {
			window.UPDATING_RULE = null
			let that = this
			teaweb.popup("/servers/components/waf/createRulePopup?type=" + this.vType, {
				callback: function (resp) {
					that.rules.push(resp.data.rule)
				}
			})
		},
		updateRule: function (index, rule) {
			window.UPDATING_RULE = rule
			let that = this
			teaweb.popup("/servers/components/waf/createRulePopup?type=" + this.vType, {
				callback: function (resp) {
					Vue.set(that.rules, index, resp.data.rule)
				}
			})
		},
		removeRule: function (index) {
			let that = this
			teaweb.confirm("确定要删除此规则吗？", function () {
				that.rules.$remove(index)
			})
		}
	},
	template: `<div>
		<input type="hidden" name="rulesJSON" :value="JSON.stringify(rules)"/>
		<div v-if="rules.length > 0">
			<div v-for="(rule, index) in rules" class="ui label small basic" style="margin-bottom: 0.5em">
				{{rule.name}}[{{rule.param}}] 
				
				<!-- cc2 -->
				<span v-if="rule.param == '\${cc2}'">
					{{rule.checkpointOptions.period}}秒/{{rule.checkpointOptions.threshold}}请求
				</span>	
				
				<!-- refererBlock -->
				<span v-if="rule.param == '\${refererBlock}'">
					{{rule.checkpointOptions.allowDomains}}
				</span>
				
				<span v-else>
					<span v-if="rule.paramFilters != null && rule.paramFilters.length > 0" v-for="paramFilter in rule.paramFilters"> | {{paramFilter.code}}</span> <var>{{rule.operator}}</var> {{rule.value}}
				</span>
				
				<a href="" title="修改" @click.prevent="updateRule(index, rule)"><i class="icon pencil small"></i></a>
				<a href="" title="删除" @click.prevent="removeRule(index)"><i class="icon remove"></i></a>
			</div>
			<div class="ui divider"></div>
		</div>
		<button class="ui button tiny" type="button" @click.prevent="addRule()">+</button>
</div>`
})

Vue.component("http-fastcgi-box", {
	props: ["v-fastcgi-ref", "v-fastcgi-configs", "v-is-location"],
	data: function () {
		let fastcgiRef = this.vFastcgiRef
		if (fastcgiRef == null) {
			fastcgiRef = {
				isPrior: false,
				isOn: false,
				fastcgiIds: []
			}
		}
		let fastcgiConfigs = this.vFastcgiConfigs
		if (fastcgiConfigs == null) {
			fastcgiConfigs = []
		} else {
			fastcgiRef.fastcgiIds = fastcgiConfigs.map(function (v) {
				return v.id
			})
		}

		return {
			fastcgiRef: fastcgiRef,
			fastcgiConfigs: fastcgiConfigs,
			advancedVisible: false
		}
	},
	methods: {
		isOn: function () {
			return (!this.vIsLocation || this.fastcgiRef.isPrior) && this.fastcgiRef.isOn
		},
		createFastcgi: function () {
			let that = this
			teaweb.popup("/servers/server/settings/fastcgi/createPopup", {
				height: "26em",
				callback: function (resp) {
					teaweb.success("添加成功", function () {
						that.fastcgiConfigs.push(resp.data.fastcgi)
						that.fastcgiRef.fastcgiIds.push(resp.data.fastcgi.id)
					})
				}
			})
		},
		updateFastcgi: function (fastcgiId, index) {
			let that = this
			teaweb.popup("/servers/server/settings/fastcgi/updatePopup?fastcgiId=" + fastcgiId, {
				callback: function (resp) {
					teaweb.success("修改成功", function () {
						Vue.set(that.fastcgiConfigs, index, resp.data.fastcgi)
					})
				}
			})
		},
		removeFastcgi: function (index) {
			this.fastcgiRef.fastcgiIds.$remove(index)
			this.fastcgiConfigs.$remove(index)
		}
	},
	template: `<div>
	<input type="hidden" name="fastcgiRefJSON" :value="JSON.stringify(fastcgiRef)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="fastcgiRef" v-if="vIsLocation"></prior-checkbox>
		<tbody v-show="(!this.vIsLocation || this.fastcgiRef.isPrior)">
			<tr>
				<td class="title">是否启用配置</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="fastcgiRef.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-if="isOn()">
			<tr>
				<td>Fastcgi服务</td>
				<td>
					<div v-show="fastcgiConfigs.length > 0" style="margin-bottom: 0.5em">
						<div class="ui label basic small" :class="{disabled: !fastcgi.isOn}" v-for="(fastcgi, index) in fastcgiConfigs">
							{{fastcgi.address}} &nbsp; <a href="" title="修改" @click.prevent="updateFastcgi(fastcgi.id, index)"><i class="ui icon pencil small"></i></a> &nbsp; <a href="" title="删除" @click.prevent="removeFastcgi(index)"><i class="ui icon remove"></i></a>
						</div>
						<div class="ui divided"></div>
					</div>
					<button type="button" class="ui button tiny" @click.prevent="createFastcgi()">+</button>
				</td>
			</tr>
		</tbody>
	</table>	
	<div class="margin"></div>
</div>`
})

// URL扩展名条件
Vue.component("http-cond-url-extension", {
	props: ["v-cond"],
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPathExtension}",
			operator: "in",
			value: "[]"
		}
		if (this.vCond != null && this.vCond.param == cond.param) {
			cond.value = this.vCond.value
		}

		let extensions = []
		try {
			extensions = JSON.parse(cond.value)
		} catch (e) {

		}

		return {
			cond: cond,
			extensions: extensions, // TODO 可以拖动排序

			isAdding: false,
			addingExt: ""
		}
	},
	watch: {
		extensions: function () {
			this.cond.value = JSON.stringify(this.extensions)
		}
	},
	methods: {
		addExt: function () {
			this.isAdding = !this.isAdding

			if (this.isAdding) {
				let that = this
				setTimeout(function () {
					that.$refs.addingExt.focus()
				}, 100)
			}
		},
		cancelAdding: function () {
			this.isAdding = false
			this.addingExt = ""
		},
		confirmAdding: function () {
			// TODO 做更详细的校验
			// TODO 如果有重复的则提示之

			if (this.addingExt.length == 0) {
				return
			}
			if (this.addingExt[0] != ".") {
				this.addingExt = "." + this.addingExt
			}
			this.addingExt = this.addingExt.replace(/\s+/g, "").toLowerCase()
			this.extensions.push(this.addingExt)

			// 清除状态
			this.cancelAdding()
		},
		removeExt: function (index) {
			this.extensions.$remove(index)
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<div v-if="extensions.length > 0">
		<div class="ui label small" v-for="(ext, index) in extensions">{{ext}} <a href="" title="删除" @click.prevent="removeExt(index)"><i class="icon remove"></i></a></div>
		<div class="ui divider"></div>
	</div>
	<div class="ui fields inline" v-if="isAdding">
		<div class="ui field">
			<input type="text" size="6" maxlength="100" v-model="addingExt" ref="addingExt" placeholder=".xxx" @keyup.enter="confirmAdding" @keypress.enter.prevent="1" />
		</div>
		<div class="ui field">
			<button class="ui button tiny" type="button" @click.prevent="confirmAdding">确认</button>
			<a href="" title="取消" @click.prevent="cancelAdding"><i class="icon remove"></i></a>
		</div> 
	</div>
	<div style="margin-top: 1em">
		<button class="ui button tiny" type="button" @click.prevent="addExt()">+添加扩展名</button>
	</div>
	<p class="comment">扩展名需要包含点（.）符号，例如<span class="ui label tiny">.jpg</span>、<span class="ui label tiny">.png</span>之类。</p>
</div>`
})

// 根据URL前缀
Vue.component("http-cond-url-prefix", {
	props: ["v-cond"],
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPath}",
			operator: "prefix",
			value: ""
		}
		if (this.vCond != null && typeof (this.vCond.value) == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value"/>
	<p class="comment">URL前缀，有此前缀的URL都将会被匹配，通常以<code-label>/</code-label>开头，比如<code-label>/static</code-label>。</p>
</div>`
})

Vue.component("http-cond-url-not-prefix", {
	props: ["v-cond"],
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPath}",
			operator: "prefix",
			value: "",
			isReverse: true
		}
		if (this.vCond != null && typeof this.vCond.value == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value"/>
	<p class="comment">要排除的URL前缀，有此前缀的URL都将会被匹配，通常以<code-label>/</code-label>开头，比如<code-label>/static</code-label>。</p>
</div>`
})

// URL精准匹配
Vue.component("http-cond-url-eq", {
	props: ["v-cond"],
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPath}",
			operator: "eq",
			value: ""
		}
		if (this.vCond != null && typeof this.vCond.value == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value"/>
	<p class="comment">完整的URL路径，通常以<code-label>/</code-label>开头，比如<code-label>/static/ui.js</code-label>。</p>
</div>`
})

Vue.component("http-cond-url-not-eq", {
	props: ["v-cond"],
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPath}",
			operator: "eq",
			value: "",
			isReverse: true
		}
		if (this.vCond != null && typeof this.vCond.value == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value"/>
	<p class="comment">要排除的完整的URL路径，通常以<code-label>/</code-label>开头，比如<code-label>/static/ui.js</code-label>。</p>
</div>`
})

// URL正则匹配
Vue.component("http-cond-url-regexp", {
	props: ["v-cond"],
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPath}",
			operator: "regexp",
			value: ""
		}
		if (this.vCond != null && typeof this.vCond.value == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value"/>
	<p class="comment">匹配URL的正则表达式，比如<code-label>^/static/(.*).js$</code-label>。</p>
</div>`
})

// 排除URL正则匹配
Vue.component("http-cond-url-not-regexp", {
	props: ["v-cond"],
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPath}",
			operator: "not regexp",
			value: ""
		}
		if (this.vCond != null && typeof this.vCond.value == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value"/>
	<p class="comment"><strong>不要</strong>匹配URL的正则表达式，意即只要匹配成功则排除此条件，比如<code-label>^/static/(.*).js$</code-label>。</p>
</div>`
})

// 根据MimeType
Vue.component("http-cond-mime-type", {
	props: ["v-cond"],
	data: function () {
		let cond = {
			isRequest: false,
			param: "${response.contentType}",
			operator: "mime type",
			value: "[]"
		}
		if (this.vCond != null && this.vCond.param == cond.param) {
			cond.value = this.vCond.value
		}
		return {
			cond: cond,
			mimeTypes: JSON.parse(cond.value), // TODO 可以拖动排序

			isAdding: false,
			addingMimeType: ""
		}
	},
	watch: {
		mimeTypes: function () {
			this.cond.value = JSON.stringify(this.mimeTypes)
		}
	},
	methods: {
		addMimeType: function () {
			this.isAdding = !this.isAdding

			if (this.isAdding) {
				let that = this
				setTimeout(function () {
					that.$refs.addingMimeType.focus()
				}, 100)
			}
		},
		cancelAdding: function () {
			this.isAdding = false
			this.addingMimeType = ""
		},
		confirmAdding: function () {
			// TODO 做更详细的校验
			// TODO 如果有重复的则提示之

			if (this.addingMimeType.length == 0) {
				return
			}
			this.addingMimeType = this.addingMimeType.replace(/\s+/g, "")
			this.mimeTypes.push(this.addingMimeType)

			// 清除状态
			this.cancelAdding()
		},
		removeMimeType: function (index) {
			this.mimeTypes.$remove(index)
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<div v-if="mimeTypes.length > 0">
		<div class="ui label small" v-for="(mimeType, index) in mimeTypes">{{mimeType}} <a href="" title="删除" @click.prevent="removeMimeType(index)"><i class="icon remove"></i></a></div>
		<div class="ui divider"></div>
	</div>
	<div class="ui fields inline" v-if="isAdding">
		<div class="ui field">
			<input type="text" size="16" maxlength="100" v-model="addingMimeType" ref="addingMimeType" placeholder="类似于image/png" @keyup.enter="confirmAdding" @keypress.enter.prevent="1" />
		</div>
		<div class="ui field">
			<button class="ui button tiny" type="button" @click.prevent="confirmAdding">确认</button>
			<a href="" title="取消" @click.prevent="cancelAdding"><i class="icon remove"></i></a>
		</div> 
	</div>
	<div style="margin-top: 1em">
		<button class="ui button tiny" type="button" @click.prevent="addMimeType()">+添加MimeType</button>
	</div>
	<p class="comment">服务器返回的内容的MimeType，比如<span class="ui label tiny">text/html</span>、<span class="ui label tiny">image/*</span>等。</p>
</div>`
})

// 参数匹配
Vue.component("http-cond-params", {
	props: ["v-cond"],
	mounted: function () {
		let cond = this.vCond
		if (cond == null) {
			return
		}
		this.operator = cond.operator

		// stringValue
		if (["regexp", "not regexp", "eq", "not", "prefix", "suffix", "contains", "not contains", "eq ip", "gt ip", "gte ip", "lt ip", "lte ip", "ip range"].$contains(cond.operator)) {
			this.stringValue = cond.value
			return
		}

		// numberValue
		if (["eq int", "eq float", "gt", "gte", "lt", "lte", "mod 10", "ip mod 10", "mod 100", "ip mod 100"].$contains(cond.operator)) {
			this.numberValue = cond.value
			return
		}

		// modValue
		if (["mod", "ip mod"].$contains(cond.operator)) {
			let pieces = cond.value.split(",")
			this.modDivValue = pieces[0]
			if (pieces.length > 1) {
				this.modRemValue = pieces[1]
			}
			return
		}

		// stringValues
		let that = this
		if (["in", "not in", "file ext", "mime type"].$contains(cond.operator)) {
			try {
				let arr = JSON.parse(cond.value)
				if (arr != null && (arr instanceof Array)) {
					arr.forEach(function (v) {
						that.stringValues.push(v)
					})
				}
			} catch (e) {

			}
			return
		}

		// versionValue
		if (["version range"].$contains(cond.operator)) {
			let pieces = cond.value.split(",")
			this.versionRangeMinValue = pieces[0]
			if (pieces.length > 1) {
				this.versionRangeMaxValue = pieces[1]
			}
			return
		}
	},
	data: function () {
		let cond = {
			isRequest: true,
			param: "",
			operator: window.REQUEST_COND_OPERATORS[0].op,
			value: ""
		}
		if (this.vCond != null) {
			cond = this.vCond
		}
		return {
			cond: cond,
			operators: window.REQUEST_COND_OPERATORS,
			operator: window.REQUEST_COND_OPERATORS[0].op,
			operatorDescription: window.REQUEST_COND_OPERATORS[0].description,
			variables: window.REQUEST_VARIABLES,
			variable: "",

			// 各种类型的值
			stringValue: "",
			numberValue: "",

			modDivValue: "",
			modRemValue: "",

			stringValues: [],

			versionRangeMinValue: "",
			versionRangeMaxValue: ""
		}
	},
	methods: {
		changeVariable: function () {
			let v = this.cond.param
			if (v == null) {
				v = ""
			}
			this.cond.param = v + this.variable
		},
		changeOperator: function () {
			let that = this
			this.operators.forEach(function (v) {
				if (v.op == that.operator) {
					that.operatorDescription = v.description
				}
			})

			this.cond.operator = this.operator

			// 移动光标
			let box = document.getElementById("variables-value-box")
			if (box != null) {
				setTimeout(function () {
					let input = box.getElementsByTagName("INPUT")
					if (input.length > 0) {
						input[0].focus()
					}
				}, 100)
			}
		},
		changeStringValues: function (v) {
			this.stringValues = v
			this.cond.value = JSON.stringify(v)
		}
	},
	watch: {
		stringValue: function (v) {
			this.cond.value = v
		},
		numberValue: function (v) {
			// TODO 校验数字
			this.cond.value = v
		},
		modDivValue: function (v) {
			if (v.length == 0) {
				return
			}
			let div = parseInt(v)
			if (isNaN(div)) {
				div = 1
			}
			this.modDivValue = div
			this.cond.value = div + "," + this.modRemValue
		},
		modRemValue: function (v) {
			if (v.length == 0) {
				return
			}
			let rem = parseInt(v)
			if (isNaN(rem)) {
				rem = 0
			}
			this.modRemValue = rem
			this.cond.value = this.modDivValue + "," + rem
		},
		versionRangeMinValue: function (v) {
			this.cond.value = this.versionRangeMinValue + "," + this.versionRangeMaxValue
		},
		versionRangeMaxValue: function (v) {
			this.cond.value = this.versionRangeMinValue + "," + this.versionRangeMaxValue
		}
	},
	template: `<tbody>
	<tr>
		<td>参数值</td>
		<td>
			<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
			<div>
				<div class="ui fields inline">
					<div class="ui field">
						<input type="text" placeholder="\${xxx}" v-model="cond.param"/>
					</div>
					<div class="ui field">
						<select class="ui dropdown" style="width: 7em; color: grey" v-model="variable" @change="changeVariable">
							<option value="">[常用参数]</option>
							<option v-for="v in variables" :value="v.code">{{v.code}} - {{v.name}}</option>
						</select>
					</div>
				</div>			
			</div>
			<p class="comment">其中可以使用变量，类似于<code-label>\${requestPath}</code-label>，也可以是多个变量的组合。</p>
		</td>
	</tr>
	<tr>
		<td>操作符</td>
		<td>
			<div>
				<select class="ui dropdown auto-width" v-model="operator" @change="changeOperator">
					<option v-for="operator in operators" :value="operator.op">{{operator.name}}</option>
				</select>
				<p class="comment">{{operatorDescription}}</p>
			</div>
		</td>
	</tr>
	<tr v-show="!['file exist', 'file not exist'].$contains(cond.operator)">
		<td>对比值</td>
		<td id="variables-value-box">
			<!-- 正则表达式 -->
			<div v-if="['regexp', 'not regexp'].$contains(cond.operator)">
				<input type="text" v-model="stringValue"/>
				<p class="comment">要匹配的正则表达式，比如<code-label>^/static/(.+).js</code-label>。</p>
			</div>
			
			<!-- 数字相关 -->
			<div v-if="['eq int', 'eq float', 'gt', 'gte', 'lt', 'lte'].$contains(cond.operator)">
				<input type="text" maxlength="11" size="11" style="width: 5em" v-model="numberValue"/>
				<p class="comment">要对比的数字。</p>
			</div>
			
			<!-- 取模 -->
			<div v-if="['mod 10'].$contains(cond.operator)">
				<input type="text" maxlength="11" size="11" style="width: 5em" v-model="numberValue"/>
				<p class="comment">参数值除以10的余数，在0-9之间。</p>
			</div>
			<div v-if="['mod 100'].$contains(cond.operator)">
				<input type="text" maxlength="11" size="11" style="width: 5em" v-model="numberValue"/>
				<p class="comment">参数值除以100的余数，在0-99之间。</p>
			</div>
			<div v-if="['mod', 'ip mod'].$contains(cond.operator)">
				<div class="ui fields inline">
					<div class="ui field">除：</div>
					<div class="ui field">
						<input type="text" maxlength="11" size="11" style="width: 5em" v-model="modDivValue" placeholder="除数"/>
					</div>
					<div class="ui field">余：</div>
					<div class="ui field">
						<input type="text" maxlength="11" size="11" style="width: 5em" v-model="modRemValue" placeholder="余数"/>
					</div>
				</div>
			</div>
			
			<!-- 字符串相关 -->
			<div v-if="['eq', 'not', 'prefix', 'suffix', 'contains', 'not contains'].$contains(cond.operator)">
				<input type="text" v-model="stringValue"/>
				<p class="comment" v-if="cond.operator == 'eq'">和参数值一致的字符串。</p>
				<p class="comment" v-if="cond.operator == 'not'">和参数值不一致的字符串。</p>
				<p class="comment" v-if="cond.operator == 'prefix'">参数值的前缀。</p>
				<p class="comment" v-if="cond.operator == 'suffix'">参数值的后缀为此字符串。</p>
				<p class="comment" v-if="cond.operator == 'contains'">参数值包含此字符串。</p>
				<p class="comment" v-if="cond.operator == 'not contains'">参数值不包含此字符串。</p>
			</div>
			<div v-if="['in', 'not in', 'file ext', 'mime type'].$contains(cond.operator)">
				<values-box @change="changeStringValues" :values="stringValues" size="15"></values-box>
				<p class="comment" v-if="cond.operator == 'in'">添加参数值列表。</p>
				<p class="comment" v-if="cond.operator == 'not in'">添加参数值列表。</p>
				<p class="comment" v-if="cond.operator == 'file ext'">添加扩展名列表，比如<code-label>png</code-label>、<code-label>html</code-label>，不包括点。</p>
				<p class="comment" v-if="cond.operator == 'mime type'">添加MimeType列表，类似于<code-label>text/html</code-label>、<code-label>image/*</code-label>。</p>
			</div>
			<div v-if="['version range'].$contains(cond.operator)">
				<div class="ui fields inline">
					<div class="ui field"><input type="text" v-model="versionRangeMinValue" maxlength="200" placeholder="最小版本" style="width: 10em"/></div>
					<div class="ui field">-</div>
					<div class="ui field"><input type="text" v-model="versionRangeMaxValue" maxlength="200" placeholder="最大版本" style="width: 10em"/></div>
				</div>
			</div>
			
			<!-- IP相关 -->
			<div v-if="['eq ip', 'gt ip', 'gte ip', 'lt ip', 'lte ip', 'ip range'].$contains(cond.operator)">
				<input type="text" style="width: 10em" v-model="stringValue" placeholder="x.x.x.x"/>
				<p class="comment">要对比的IP。</p>
			</div>
			<div v-if="['ip mod 10'].$contains(cond.operator)">
				<input type="text" maxlength="11" size="11" style="width: 5em" v-model="numberValue"/>
				<p class="comment">参数中IP转换成整数后除以10的余数，在0-9之间。</p>
			</div>
			<div v-if="['ip mod 100'].$contains(cond.operator)">
				<input type="text" maxlength="11" size="11" style="width: 5em" v-model="numberValue"/>
				<p class="comment">参数中IP转换成整数后除以100的余数，在0-99之间。</p>
			</div>
		</td>
	</tr>
</tbody>`
})

Vue.component("server-group-selector", {
	props: ["v-groups"],
	data: function () {
		let groups = this.vGroups
		if (groups == null) {
			groups = []
		}
		return {
			groups: groups
		}
	},
	methods: {
		selectGroup: function () {
			let that = this
			let groupIds = this.groups.map(function (v) {
				return v.id.toString()
			}).join(",")
			teaweb.popup("/servers/groups/selectPopup?selectedGroupIds=" + groupIds, {
				callback: function (resp) {
					that.groups.push(resp.data.group)
				}
			})
		},
		addGroup: function () {
			let that = this
			teaweb.popup("/servers/groups/createPopup", {
				callback: function (resp) {
					that.groups.push(resp.data.group)
				}
			})
		},
		removeGroup: function (index) {
			this.groups.$remove(index)
		},
		groupIds: function () {
			return this.groups.map(function (v) {
				return v.id
			})
		}
	},
	template: `<div>
	<div v-if="groups.length > 0">
		<div class="ui label small basic" v-if="groups.length > 0" v-for="(group, index) in groups">
			<input type="hidden" name="groupIds" :value="group.id"/>
			{{group.name}} &nbsp;<a href="" title="删除" @click.prevent="removeGroup(index)"><i class="icon remove"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div>
		<a href="" @click.prevent="selectGroup()">[选择分组]</a> &nbsp; <a href="" @click.prevent="addGroup()">[添加分组]</a>
	</div>
</div>`
})

// 指标周期设置
Vue.component("metric-period-config-box", {
	props: ["v-period", "v-period-unit"],
	data: function () {
		let period = this.vPeriod
		let periodUnit = this.vPeriodUnit
		if (period == null || period.toString().length == 0) {
			period = 1
		}
		if (periodUnit == null || periodUnit.length == 0) {
			periodUnit = "day"
		}
		return {
			periodConfig: {
				period: period,
				unit: periodUnit
			}
		}
	},
	watch: {
		"periodConfig.period": function (v) {
			v = parseInt(v)
			if (isNaN(v) || v <= 0) {
				v = 1
			}
			this.periodConfig.period = v
		}
	},
	template: `<div>
	<input type="hidden" name="periodJSON" :value="JSON.stringify(periodConfig)"/>
	<div class="ui fields inline">
		<div class="ui field">
			<input type="text" v-model="periodConfig.period" maxlength="4" size="4"/>
		</div>
		<div class="ui field">
			<select class="ui dropdown" v-model="periodConfig.unit">
				<option value="minute">分钟</option>
				<option value="hour">小时</option>
				<option value="day">天</option>
				<option value="week">周</option>
				<option value="month">月</option>
			</select>
		</div>
	</div>
	<p class="comment">在此周期内同一对象累积为同一数据。</p>
</div>`
})

Vue.component("traffic-limit-config-box", {
	props: ["v-traffic-limit"],
	data: function () {
		let config = this.vTrafficLimit
		if (config == null) {
			config = {
				isOn: false,
				dailySize: {
					count: -1,
					unit: "gb"
				},
				monthlySize: {
					count: -1,
					unit: "gb"
				},
				totalSize: {
					count: -1,
					unit: "gb"
				},
				noticePageBody: ""
			}
		}
		if (config.dailySize == null) {
			config.dailySize = {
				count: -1,
				unit: "gb"
			}
		}
		if (config.monthlySize == null) {
			config.monthlySize = {
				count: -1,
				unit: "gb"
			}
		}
		if (config.totalSize == null) {
			config.totalSize = {
				count: -1,
				unit: "gb"
			}
		}
		return {
			config: config
		}
	},
	methods: {
		showBodyTemplate: function () {
			this.config.noticePageBody = `<!DOCTYPE html>
<html>
<head>
<title>Traffic Limit Exceeded Warning</title>
<body>

The site traffic has exceeded the limit. Please contact with the site administrator.

</body>
</html>`
		}
	},
	template: `<div>
	<input type="hidden" name="trafficLimitJSON" :value="JSON.stringify(config)"/>
	<table class="ui table selectable definition">
		<tbody>
			<tr>
				<td class="title">是否启用</td>
				<td>
					<checkbox v-model="config.isOn"></checkbox>
					<p class="comment">注意：由于流量统计是每5分钟统计一次，所以超出流量限制后，对用户的提醒也会有所延迟。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="config.isOn">
			<tr>
				<td>日流量限制</td>
				<td>
					<size-capacity-box :v-value="config.dailySize"></size-capacity-box>
				</td>
			</tr>
			<tr>
				<td>月流量限制</td>
				<td>
					<size-capacity-box :v-value="config.monthlySize"></size-capacity-box>
				</td>
			</tr>
			<!--<tr>
				<td>总体限制</td>
				<td>
					<size-capacity-box :v-value="config.totalSize"></size-capacity-box>
					<p class="comment"></p>
				</td>
			</tr>-->
			<tr>
				<td>网页提示内容</td>
				<td>
					<textarea v-model="config.noticePageBody"></textarea>
					<p class="comment"><a href="" @click.prevent="showBodyTemplate">[使用模板]</a>。当达到流量限制时网页显示的HTML内容，不填写则显示默认的提示内容。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})

// TODO 支持关键词搜索
// TODO 改成弹窗选择
Vue.component("admin-selector", {
    props: ["v-admin-id"],
    mounted: function () {
        let that = this
        Tea.action("/admins/options")
            .post()
            .success(function (resp) {
                that.admins = resp.data.admins
            })
    },
    data: function () {
        let adminId = this.vAdminId
        if (adminId == null) {
            adminId = 0
        }
        return {
            admins: [],
            adminId: adminId
        }
    },
    template: `<div>
    <select class="ui dropdown auto-width" name="adminId" v-model="adminId">
        <option value="0">[选择系统用户]</option>
        <option v-for="admin in admins" :value="admin.id">{{admin.name}}（{{admin.username}}）</option>
    </select>
</div>`
})

// 绑定IP列表
Vue.component("ip-list-bind-box", {
	props: ["v-http-firewall-policy-id", "v-type"],
	mounted: function () {
		this.refresh()
	},
	data: function () {
		return {
			policyId: this.vHttpFirewallPolicyId,
			type: this.vType,
			lists: []
		}
	},
	methods: {
		bind: function () {
			let that = this
			teaweb.popup("/servers/iplists/bindHTTPFirewallPopup?httpFirewallPolicyId=" + this.policyId + "&type=" + this.type, {
				width: "50em",
				height: "34em",
				callback: function () {

				},
				onClose: function () {
					that.refresh()
				}
			})
		},
		remove: function (index, listId) {
			let that = this
			teaweb.confirm("确定要删除这个绑定的IP名单吗？", function () {
				Tea.action("/servers/iplists/unbindHTTPFirewall")
					.params({
						httpFirewallPolicyId: that.policyId,
						listId: listId
					})
					.post()
					.success(function (resp) {
						that.lists.$remove(index)
					})
			})
		},
		refresh: function () {
			let that = this
			Tea.action("/servers/iplists/httpFirewall")
				.params({
					httpFirewallPolicyId: this.policyId,
					type: this.vType
				})
				.post()
				.success(function (resp) {
					that.lists = resp.data.lists
				})
		}
	},
	template: `<div>
	<a href="" @click.prevent="bind()" style="color: rgba(0,0,0,.6)">绑定+</a> &nbsp; <span v-if="lists.length > 0"><span class="disabled small">|&nbsp;</span> 已绑定：</span>
	<div class="ui label basic small" v-for="(list, index) in lists">
		<a :href="'/servers/iplists/list?listId=' + list.id" title="点击查看详情" style="opacity: 1">{{list.name}}</a>
		<a href="" title="删除" @click.prevent="remove(index, list.id)"><i class="icon remove small"></i></a>
	</div>
</div>`
})

Vue.component("ip-list-table", {
	props: ["v-items", "v-keyword", "v-show-search-button"],
	data: function () {
		return {
			items: this.vItems,
			keyword: (this.vKeyword != null) ? this.vKeyword : "",
			selectedAll: false,
			hasSelectedItems: false
		}
	},
	methods: {
		updateItem: function (itemId) {
			this.$emit("update-item", itemId)
		},
		deleteItem: function (itemId) {
			this.$emit("delete-item", itemId)
		},
		viewLogs: function (itemId) {
			teaweb.popup("/servers/iplists/accessLogsPopup?itemId=" + itemId, {
				width: "50em",
				height: "30em"
			})
		},
		changeSelectedAll: function () {
			let boxes = this.$refs.itemCheckBox
			if (boxes == null) {
				return
			}

			let that = this
			boxes.forEach(function (box) {
				box.checked = that.selectedAll
			})

			this.hasSelectedItems = this.selectedAll
		},
		changeSelected: function (e) {
			let that = this
			that.hasSelectedItems = false
			let boxes = that.$refs.itemCheckBox
			if (boxes == null) {
				return
			}
			boxes.forEach(function (box) {
				if (box.checked) {
					that.hasSelectedItems = true
				}
			})
		},
		deleteAll: function () {
			let boxes = this.$refs.itemCheckBox
			if (boxes == null) {
				return
			}
			let itemIds = []
			boxes.forEach(function (box) {
				if (box.checked) {
					itemIds.push(box.value)
				}
			})
			if (itemIds.length == 0) {
				return
			}

			Tea.action("/servers/iplists/deleteItems")
				.post()
				.params({
					itemIds: itemIds
				})
				.success(function () {
					teaweb.successToast("批量删除成功", 1200, teaweb.reload)
				})
		}
	},
	template: `<div>
 <div v-show="hasSelectedItems">
 	<a href="" @click.prevent="deleteAll">[批量删除]</a>
</div>
 <table class="ui table selectable celled" v-if="items.length > 0">
        <thead>
            <tr>
            	<th style="width: 1em">
            		<div class="ui checkbox">
						<input type="checkbox" v-model="selectedAll" @change="changeSelectedAll"/>
						<label></label>
					</div>
				</th>
                <th style="width:18em">IP</th>
                <th>类型</th>
                <th>级别</th>
                <th>过期时间</th>
                <th>备注</th>
                <th class="three op">操作</th>
            </tr>
        </thead>
		<tbody v-for="item in items">
			<tr>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" :value="item.id" @change="changeSelected" ref="itemCheckBox"/>
						<label></label>
					</div>
				</td>
				<td>
					<span v-if="item.type != 'all'">
					<keyword :v-word="keyword">{{item.ipFrom}}</keyword> <span>&nbsp;<a :href="'/servers/iplists?ip=' + item.ipFrom" v-if="vShowSearchButton" title="搜索此IP"><span><i class="icon search small" style="color: #ccc"></i></span></a></span>
					<span v-if="item.ipTo.length > 0"> - <keyword :v-word="keyword">{{item.ipTo}}</keyword></span></span>
					<span v-else class="disabled">*</span>
					<div v-if="item.createdTime != null">
						<span class="small grey">添加于 {{item.createdTime}}
							<span v-if="item.list != null && item.list.id > 0">
								@ 
								
								<a :href="'/servers/iplists/list?listId=' + item.list.id" v-if="item.policy.id == 0"><span>[<span v-if="item.list.type == 'black'">黑</span><span v-if="item.list.type == 'white'">白</span>名单：{{item.list.name}}]</span></a>
								<span v-else>[<span v-if="item.list.type == 'black'">黑</span><span v-if="item.list.type == 'white'">白</span>名单：{{item.list.name}}</span>
								
								<span v-if="item.policy.id > 0">
									<span v-if="item.policy.server != null">
										<a :href="'/servers/server/settings/waf/ipadmin/allowList?serverId=' + item.policy.server.id + '&firewallPolicyId=' + item.policy.id" v-if="item.list.type == 'white'">[服务：{{item.policy.server.name}}]</a>
										<a :href="'/servers/server/settings/waf/ipadmin/denyList?serverId=' + item.policy.server.id + '&firewallPolicyId=' + item.policy.id" v-if="item.list.type == 'black'">[服务：{{item.policy.server.name}}]</a>
									</span>
									<span v-else>
										<a :href="'/servers/components/waf/ipadmin/lists?firewallPolicyId=' + item.policy.id +  '&type=' + item.list.type">[策略：{{item.policy.name}}]</a>
									</span>
								</span>
							</span>
						</span>
					</div>
				</td>
				<td>
					<span v-if="item.type.length == 0">IPv4</span>
					<span v-else-if="item.type == 'ipv4'">IPv4</span>
					<span v-else-if="item.type == 'ipv6'">IPv6</span>
					<span v-else-if="item.type == 'all'"><strong>所有IP</strong></span>
				</td>
				<td>
					<span v-if="item.eventLevelName != null && item.eventLevelName.length > 0">{{item.eventLevelName}}</span>
					<span v-else class="disabled">-</span>
				</td>
				<td>
					<div v-if="item.expiredTime.length > 0">
						{{item.expiredTime}}
						<div v-if="item.isExpired" style="margin-top: 0.5em">
							<span class="ui label tiny basic red">已过期</span>
						</div>
					</div>
					<span v-else class="disabled">不过期</span>
				</td>
				<td>
					<span v-if="item.reason.length > 0">{{item.reason}}</span>
					<span v-else class="disabled">-</span>
					
					<div style="margin-top: 0.4em" v-if="item.sourceServer != null && item.sourceServer.id > 0">
						<a :href="'/servers/server?serverId=' + item.sourceServer.id" style="border: 0"><span class="small "><i class="icon clone outline"></i>{{item.sourceServer.name}}</span></a>
					</div>
					<div v-if="item.sourcePolicy != null && item.sourcePolicy.id > 0" style="margin-top: 0.4em">
						<a :href="'/servers/components/waf/group?firewallPolicyId=' +  item.sourcePolicy.id + '&type=inbound&groupId=' + item.sourceGroup.id" v-if="item.sourcePolicy.serverId == 0"><span class="small "><i class="icon shield"></i>{{item.sourcePolicy.name}} &raquo; {{item.sourceGroup.name}} &raquo; {{item.sourceSet.name}}</span></a>
						<a :href="'/servers/server/settings/waf/group?serverId=' + item.sourcePolicy.serverId + '&firewallPolicyId=' + item.sourcePolicy.id + '&type=inbound&groupId=' + item.sourceGroup.id" v-if="item.sourcePolicy.serverId > 0"><span class="small "><i class="icon shield"></i> {{item.sourcePolicy.name}} &raquo; {{item.sourceGroup.name}} &raquo; {{item.sourceSet.name}}</span></a>
					</div>
				</td>
				<td>
					<a href="" @click.prevent="viewLogs(item.id)">日志</a> &nbsp;
					<a href="" @click.prevent="updateItem(item.id)">修改</a> &nbsp;
					<a href="" @click.prevent="deleteItem(item.id)">删除</a>
				</td>
			</tr>
        </tbody>
    </table>
</div>`
})

Vue.component("ip-item-text", {
    props: ["v-item"],
    template: `<span>
    <span v-if="vItem.type == 'all'">*</span>
    <span v-if="vItem.type == 'ipv4' || vItem.type.length == 0">
        {{vItem.ipFrom}}
        <span v-if="vItem.ipTo.length > 0">- {{vItem.ipTo}}</span>
    </span>
    <span v-if="vItem.type == 'ipv6'">{{vItem.ipFrom}}</span>
    <span v-if="vItem.eventLevelName != null && vItem.eventLevelName.length > 0">&nbsp; 级别：{{vItem.eventLevelName}}</span>
</span>`
})

Vue.component("ip-box", {
	props: [],
	methods: {
		popup: function () {
			let e = this.$refs.container
			let text = e.innerText
			if (text == null) {
				text = e.textContent
			}

			teaweb.popup("/servers/ipbox?ip=" + text, {
				width: "50em",
				height: "30em"
			})
		}
	},
	template: `<span @click.prevent="popup()" ref="container"><slot></slot></span>`
})

Vue.component("api-node-selector", {
	props: [],
	data: function () {
		return {}
	},
	template: `<div>
	暂未实现
</div>`
})

Vue.component("api-node-addresses-box", {
	props: ["v-addrs", "v-name"],
	data: function () {
		let addrs = this.vAddrs
		if (addrs == null) {
			addrs = []
		}
		return {
			addrs: addrs
		}
	},
	methods: {
		// 添加IP地址
		addAddr: function () {
			let that = this;
			teaweb.popup("/api/node/createAddrPopup", {
				height: "16em",
				callback: function (resp) {
					that.addrs.push(resp.data.addr);
				}
			})
		},

		// 修改地址
		updateAddr: function (index, addr) {
			let that = this;
			window.UPDATING_ADDR = addr
			teaweb.popup("/api/node/updateAddrPopup?addressId=", {
				callback: function (resp) {
					Vue.set(that.addrs, index, resp.data.addr);
				}
			})
		},

		// 删除IP地址
		removeAddr: function (index) {
			this.addrs.$remove(index);
		}
	},
	template: `<div>
	<input type="hidden" :name="vName" :value="JSON.stringify(addrs)"/>
	<div v-if="addrs.length > 0">
		<div>
			<div v-for="(addr, index) in addrs" class="ui label small">
				{{addr.protocol}}://{{addr.host.quoteIP()}}:{{addr.portRange}}</span>
				<a href="" title="修改" @click.prevent="updateAddr(index, addr)"><i class="icon pencil small"></i></a>
				<a href="" title="删除" @click.prevent="removeAddr(index)"><i class="icon remove"></i></a>
			</div>
		</div>
		<div class="ui divider"></div>
	</div>
	<div>
		<button class="ui button small" type="button" @click.prevent="addAddr()">+</button>
	</div>
</div>`
})

// 给Table增加排序功能
function sortTable(callback) {
	// 引入js
	let jsFile = document.createElement("script")
	jsFile.setAttribute("src", "/js/sortable.min.js")
	jsFile.addEventListener("load", function () {
		// 初始化
		let box = document.querySelector("#sortable-table")
		if (box == null) {
			return
		}
		Sortable.create(box, {
			draggable: "tbody",
			handle: ".icon.handle",
			onStart: function () {
			},
			onUpdate: function (event) {
				let rows = box.querySelectorAll("tbody")
				let rowIds = []
				rows.forEach(function (row) {
					rowIds.push(parseInt(row.getAttribute("v-id")))
				})
				callback(rowIds)
			}
		})
	})
	document.head.appendChild(jsFile)
}

function sortLoad(callback) {
	let jsFile = document.createElement("script")
	jsFile.setAttribute("src", "/js/sortable.min.js")
	jsFile.addEventListener("load", function () {
		if (typeof (callback) == "function") {
			callback()
		}
	})
	document.head.appendChild(jsFile)
}


Vue.component("page-box", {
	data: function () {
		return {
			page: ""
		}
	},
	created: function () {
		let that = this;
		setTimeout(function () {
			that.page = Tea.Vue.page;
		})
	},
	template: `<div>
	<div class="page" v-html="page"></div>
</div>`
})

Vue.component("network-addresses-box", {
	props: ["v-server-type", "v-addresses", "v-protocol", "v-name", "v-from", "v-support-range"],
	data: function () {
		let addresses = this.vAddresses
		if (addresses == null) {
			addresses = []
		}
		let protocol = this.vProtocol
		if (protocol == null) {
			protocol = ""
		}

		let name = this.vName
		if (name == null) {
			name = "addresses"
		}

		let from = this.vFrom
		if (from == null) {
			from = ""
		}

		return {
			addresses: addresses,
			protocol: protocol,
			name: name,
			from: from
		}
	},
	watch: {
		"vServerType": function () {
			this.addresses = []
		},
		"vAddresses": function () {
			if (this.vAddresses != null) {
				this.addresses = this.vAddresses
			}
		}
	},
	methods: {
		addAddr: function () {
			let that = this
			window.UPDATING_ADDR = null
			teaweb.popup("/servers/addPortPopup?serverType=" + this.vServerType + "&protocol=" + this.protocol + "&from=" + this.from + "&supportRange=" + (this.supportRange() ? 1 : 0), {
				height: "18em",
				callback: function (resp) {
					var addr = resp.data.address
					if (that.addresses.$find(function (k, v) {
						return addr.host == v.host && addr.portRange == v.portRange && addr.protocol == v.protocol
					}) != null) {
						teaweb.warn("要添加的网络地址已经存在")
						return
					}
					that.addresses.push(addr)
					if (["https", "https4", "https6"].$contains(addr.protocol)) {
						this.tlsProtocolName = "HTTPS"
					} else if (["tls", "tls4", "tls6"].$contains(addr.protocol)) {
						this.tlsProtocolName = "TLS"
					}

					// 发送事件
					that.$emit("change", that.addresses)
				}
			})
		},
		removeAddr: function (index) {
			this.addresses.$remove(index);

			// 发送事件
			this.$emit("change", this.addresses)
		},
		updateAddr: function (index, addr) {
			let that = this
			window.UPDATING_ADDR = addr
			teaweb.popup("/servers/addPortPopup?serverType=" + this.vServerType + "&protocol=" + this.protocol + "&from=" + this.from + "&supportRange=" + (this.supportRange() ? 1 : 0), {
				height: "18em",
				callback: function (resp) {
					var addr = resp.data.address
					Vue.set(that.addresses, index, addr)

					if (["https", "https4", "https6"].$contains(addr.protocol)) {
						this.tlsProtocolName = "HTTPS"
					} else if (["tls", "tls4", "tls6"].$contains(addr.protocol)) {
						this.tlsProtocolName = "TLS"
					}

					// 发送事件
					that.$emit("change", that.addresses)
				}
			})

			// 发送事件
			this.$emit("change", this.addresses)
		},
		supportRange: function () {
			return this.vSupportRange || (this.vServerType == "tcpProxy" || this.vServerType == "udpProxy")
		}
	},
	template: `<div>
	<input type="hidden" :name="name" :value="JSON.stringify(addresses)"/>
	<div v-if="addresses.length > 0">
		<div class="ui label small basic" v-for="(addr, index) in addresses">
			{{addr.protocol}}://<span v-if="addr.host.length > 0">{{addr.host.quoteIP()}}</span><span v-if="addr.host.length == 0">*</span>:<span v-if="addr.portRange.indexOf('-')<0">{{addr.portRange}}</span><span v-else style="font-style: italic">{{addr.portRange}}</span>
			<a href="" @click.prevent="updateAddr(index, addr)" title="修改"><i class="icon pencil small"></i></a>
			<a href="" @click.prevent="removeAddr(index)" title="删除"><i class="icon remove"></i></a> </div>
		<div class="ui divider"></div>
	</div>
	<a href="" @click.prevent="addAddr()">[添加端口绑定]</a>
</div>`
})

/**
 * 保存按钮
 */
Vue.component("submit-btn", {
	template: '<button class="ui button primary" type="submit"><slot>保存</slot></button>'
});

// 可以展示更多条目的角图表
Vue.component("more-items-angle", {
	props: ["v-data-url", "v-url"],
	data: function () {
		return {
			visible: false
		}
	},
	methods: {
		show: function () {
			this.visible = !this.visible
			if (this.visible) {
				this.showBox()
			} else {
				this.hideBox()
			}
		},
		showBox: function () {
			let that = this

			this.visible = true

			Tea.action(this.vDataUrl)
				.params({
					url: this.vUrl
				})
				.post()
				.success(function (resp) {
					let groups = resp.data.groups

					let boxLeft = that.$el.offsetLeft + 120;
					let boxTop = that.$el.offsetTop + 70;

					let box = document.createElement("div")
					box.setAttribute("id", "more-items-box")
					box.style.cssText = "z-index: 100; position: absolute; left: " + boxLeft + "px; top: " + boxTop + "px; max-height: 30em; overflow: auto; border-bottom: 1px solid rgba(34,36,38,.15)"
					document.body.append(box)

					let menuHTML = "<ul class=\"ui labeled menu vertical borderless\" style=\"padding: 0\">"
					groups.forEach(function (group) {
						menuHTML += "<div class=\"item header\">" + teaweb.encodeHTML(group.name) + "</div>"
						group.items.forEach(function (item) {
							menuHTML += "<a href=\"" + item.url + "\" class=\"item " + (item.isActive ? "active" : "") + "\" style=\"font-size: 0.9em;\">" + teaweb.encodeHTML(item.name) + "<i class=\"icon right angle\"></i></a>"
						})
					})
					menuHTML += "</ul>"
					box.innerHTML = menuHTML

					let listener = function (e) {
						if (e.target.tagName == "I") {
							return
						}

						if (!that.isInBox(box, e.target)) {
							document.removeEventListener("click", listener)
							that.hideBox()
						}
					}
					document.addEventListener("click", listener)
				})
		},
		hideBox: function () {
			let box = document.getElementById("more-items-box")
			if (box != null) {
				box.parentNode.removeChild(box)
			}
			this.visible = false
		},
		isInBox: function (parent, child) {
			while (true) {
				if (child == null) {
					break
				}
				if (child.parentNode == parent) {
					return true
				}
				child = child.parentNode
			}
			return false
		}
	},
	template: `<a href="" class="item" @click.prevent="show"><i class="icon angle" :class="{down: !visible, up: visible}"></i></a>`
})

/**
 * 菜单项
 */
Vue.component("menu-item", {
	props: ["href", "active", "code"],
	data: function () {
		let active = this.active
		if (typeof (active) == "undefined") {
			var itemCode = ""
			if (typeof (window.TEA.ACTION.data.firstMenuItem) != "undefined") {
				itemCode = window.TEA.ACTION.data.firstMenuItem
			}
			if (itemCode != null && itemCode.length > 0 && this.code != null && this.code.length > 0) {
				if (itemCode.indexOf(",") > 0) {
					active = itemCode.split(",").$contains(this.code)
				} else {
					active = (itemCode == this.code)
				}
			}
		}

		let href = (this.href == null) ? "" : this.href
		if (typeof (href) == "string" && href.length > 0 && href.startsWith(".")) {
			let qIndex = href.indexOf("?")
			if (qIndex >= 0) {
				href = Tea.url(href.substring(0, qIndex)) + href.substring(qIndex)
			} else {
				href = Tea.url(href)
			}
		}

		return {
			vHref: href,
			vActive: active
		}
	},
	methods: {
		click: function (e) {
			this.$emit("click", e)
		}
	},
	template: '\
		<a :href="vHref" class="item" :class="{active:vActive}" @click="click"><slot></slot></a> \
		'
});

// 使用Icon的链接方式
Vue.component("link-icon", {
	props: ["href", "title", "target"],
	data: function () {
		return {
			vTitle: (this.title == null) ? "打开链接" : this.title
		}
	},
	template: `<span><slot></slot>&nbsp;<a :href="href" :title="vTitle" class="link grey" :target="target"><i class="icon linkify small"></i></a></span>`
})

// 带有下划虚线的连接
Vue.component("link-red", {
	props: ["href", "title"],
	data: function () {
		let href = this.href
		if (href == null) {
			href = ""
		}
		return {
			vHref: href
		}
	},
	methods: {
		clickPrevent: function () {
			emitClick(this, arguments)
		}
	},
	template: `<a :href="vHref" :title="title" style="border-bottom: 1px #db2828 dashed" @click.prevent="clickPrevent"><span class="red"><slot></slot></span></a>`
})

// 会弹出窗口的链接
Vue.component("link-popup", {
	props: ["title"],
	methods: {
		clickPrevent: function () {
			emitClick(this, arguments)
		}
	},
	template: `<a href="" :title="title" @click.prevent="clickPrevent"><slot></slot></a>`
})

Vue.component("popup-icon", {
	props: ["title", "href", "height"],
	methods: {
		clickPrevent: function () {
			if (this.href != null && this.href.length > 0) {
				teaweb.popup(this.href, {
					height: this.height
				})
			}
		}
	},
	template: `<span><slot></slot>&nbsp;<a href="" :title="title" @click.prevent="clickPrevent"><i class="icon clone outline small"></i></a></span>`
})

// 小提示
Vue.component("tip-icon", {
	props: ["content"],
	methods: {
		showTip: function () {
			teaweb.popupTip(this.content)
		}
	},
	template: `<a href="" title="查看帮助" @click.prevent="showTip"><i class="icon question circle grey"></i></a>`
})

// 提交点击事件
function emitClick(obj, arguments) {
	let event = "click"
	let newArgs = [event]
	for (let i = 0; i < arguments.length; i++) {
		newArgs.push(arguments[i])
	}
	obj.$emit.apply(obj, newArgs)
}

Vue.component("countries-selector", {
	props: ["v-countries"],
	data: function () {
		let countries = this.vCountries
		if (countries == null) {
			countries = []
		}
		let countryIds = countries.$map(function (k, v) {
			return v.id
		})
		return {
			countries: countries,
			countryIds: countryIds
		}
	},
	methods: {
		add: function () {
			let countryStringIds = this.countryIds.map(function (v) {
				return v.toString()
			})
			let that = this
			teaweb.popup("/ui/selectCountriesPopup?countryIds=" + countryStringIds.join(","), {
				width: "48em",
				height: "23em",
				callback: function (resp) {
					that.countries = resp.data.countries
					that.change()
				}
			})
		},
		remove: function (index) {
			this.countries.$remove(index)
			this.change()
		},
		change: function () {
			this.countryIds = this.countries.$map(function (k, v) {
				return v.id
			})
		}
	},
	template: `<div>
	<input type="hidden" name="countryIdsJSON" :value="JSON.stringify(countryIds)"/>
	<div v-if="countries.length > 0" style="margin-bottom: 0.5em">
		<div v-for="(country, index) in countries" class="ui label tiny basic">{{country.name}} <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove"></i></a></div>
		<div class="ui divider"></div>
	</div>
	<div>
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})

Vue.component("more-options-tbody", {
	data: function () {
		return {
			isVisible: false
		}
	},
	methods: {
		show: function () {
			this.isVisible = !this.isVisible
			this.$emit("change", this.isVisible)
		}
	},
	template: `<tbody>
	<tr>
		<td colspan="2"><a href="" @click.prevent="show()"><span v-if="!isVisible">更多选项</span><span v-if="isVisible">收起选项</span><i class="icon angle" :class="{down:!isVisible, up:isVisible}"></i></a></td>
	</tr>
</tbody>`
})

Vue.component("download-link", {
	props: ["v-element", "v-file", "v-value"],
	created: function () {
		let that = this
		setTimeout(function () {
			that.url = that.composeURL()
		}, 1000)
	},
	data: function () {
		let filename = this.vFile
		if (filename == null || filename.length == 0) {
			filename = "unknown-file"
		}
		return {
			file: filename,
			url: this.composeURL()
		}
	},
	methods: {
		composeURL: function () {
			let text = ""
			if (this.vValue != null) {
				text = this.vValue
			} else {
				let e = document.getElementById(this.vElement)
				if (e == null) {
					teaweb.warn("找不到要下载的内容")
					return
				}
				text = e.innerText
				if (text == null) {
					text = e.textContent
				}
			}
			return Tea.url("/ui/download", {
				file: this.file,
				text: text
			})
		}
	},
	template: `<a :href="url" target="_blank" style="font-weight: normal"><slot></slot></a>`,
})

Vue.component("values-box", {
	props: ["values", "size", "maxlength", "name", "placeholder"],
	data: function () {
		let values = this.values;
		if (values == null) {
			values = [];
		}
		return {
			"vValues": values,
			"isUpdating": false,
			"isAdding": false,
			"index": 0,
			"value": "",
			isEditing: false
		}
	},
	methods: {
		create: function () {
			this.isAdding = true;
			var that = this;
			setTimeout(function () {
				that.$refs.value.focus();
			}, 200);
		},
		update: function (index) {
			this.cancel()
			this.isUpdating = true;
			this.index = index;
			this.value = this.vValues[index];
			var that = this;
			setTimeout(function () {
				that.$refs.value.focus();
			}, 200);
		},
		confirm: function () {
			if (this.value.length == 0) {
				return
			}

			if (this.isUpdating) {
				Vue.set(this.vValues, this.index, this.value);
			} else {
				this.vValues.push(this.value);
			}
			this.cancel()
			this.$emit("change", this.vValues)
		},
		remove: function (index) {
			this.vValues.$remove(index)
			this.$emit("change", this.vValues)
		},
		cancel: function () {
			this.isUpdating = false;
			this.isAdding = false;
			this.value = "";
		},
		updateAll: function (values) {
			this.vValeus = values
		},
		addValue: function (v) {
			this.vValues.push(v)
		},

		startEditing: function () {
			this.isEditing = !this.isEditing
		}
	},
	template: `<div>
	<div v-show="!isEditing && vValues.length > 0">
		<div class="ui label tiny basic" v-for="(value, index) in vValues" style="margin-top:0.4em;margin-bottom:0.4em">{{value}}</div>
		<a href="" @click.prevent="startEditing" style="font-size: 0.8em; margin-left: 0.2em">[修改]</a>
	</div>
	<div v-show="isEditing || vValues.length == 0">
		<div style="margin-bottom: 1em" v-if="vValues.length > 0">
			<div class="ui label tiny basic" v-for="(value, index) in vValues" style="margin-top:0.4em;margin-bottom:0.4em">{{value}}
				<input type="hidden" :name="name" :value="value"/>
				&nbsp; <a href="" @click.prevent="update(index)" title="修改"><i class="icon pencil small" ></i></a> 
				<a href="" @click.prevent="remove(index)" title="删除"><i class="icon remove"></i></a> 
			</div> 
			<div class="ui divider"></div>
		</div> 
		<!-- 添加|修改 -->
		<div v-if="isAdding || isUpdating">
			<div class="ui fields inline">
				<div class="ui field">
					<input type="text" :size="size" :maxlength="maxlength" :placeholder="placeholder" v-model="value" ref="value" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
				</div> 
				<div class="ui field">
					<button class="ui button small" type="button" @click.prevent="confirm()">确定</button> 
				</div>
				<div class="ui field">
					<a href="" @click.prevent="cancel()" title="取消"><i class="icon remove small"></i></a> 
				</div> 
			</div> 
		</div> 
		<div v-if="!isAdding && !isUpdating">
			<button class="ui button tiny" type="button" @click.prevent="create()">+</button> 
		</div>
	</div>	
</div>`
});

Vue.component("datetime-input", {
	props: ["v-name", "v-timestamp"],
	mounted: function () {
		let that = this
		teaweb.datepicker(this.$refs.dayInput, function (v) {
			that.day = v
			that.hour = "23"
			that.minute = "59"
			that.second = "59"
			that.change()
		})
	},
	data: function () {
		let timestamp = this.vTimestamp
		if (timestamp != null) {
			timestamp = parseInt(timestamp)
			if (isNaN(timestamp)) {
				timestamp = 0
			}
		} else {
			timestamp = 0
		}

		let day = ""
		let hour = ""
		let minute = ""
		let second = ""

		if (timestamp > 0) {
			let date = new Date()
			date.setTime(timestamp * 1000)

			let year = date.getFullYear().toString()
			let month = this.leadingZero((date.getMonth() + 1).toString(), 2)
			day = year + "-" + month + "-" + this.leadingZero(date.getDate().toString(), 2)

			hour = this.leadingZero(date.getHours().toString(), 2)
			minute = this.leadingZero(date.getMinutes().toString(), 2)
			second = this.leadingZero(date.getSeconds().toString(), 2)
		}

		return {
			timestamp: timestamp,
			day: day,
			hour: hour,
			minute: minute,
			second: second,

			hasDayError: false,
			hasHourError: false,
			hasMinuteError: false,
			hasSecondError: false
		}
	},
	methods: {
		change: function () {
			let date = new Date()

			// day
			if (!/^\d{4}-\d{1,2}-\d{1,2}$/.test(this.day)) {
				this.hasDayError = true
				return
			}
			let pieces = this.day.split("-")
			let year = parseInt(pieces[0])
			date.setFullYear(year)

			let month = parseInt(pieces[1])
			if (month < 1 || month > 12) {
				this.hasDayError = true
				return
			}
			date.setMonth(month - 1)

			let day = parseInt(pieces[2])
			if (day < 1 || day > 32) {
				this.hasDayError = true
				return
			}
			date.setDate(day)

			this.hasDayError = false

			// hour
			if (!/^\d+$/.test(this.hour)) {
				this.hasHourError = true
				return
			}
			let hour = parseInt(this.hour)
			if (isNaN(hour)) {
				this.hasHourError = true
				return
			}
			if (hour < 0 || hour >= 24) {
				this.hasHourError = true
				return
			}
			this.hasHourError = false
			date.setHours(hour)

			// minute
			if (!/^\d+$/.test(this.minute)) {
				this.hasMinuteError = true
				return
			}
			let minute = parseInt(this.minute)
			if (isNaN(minute)) {
				this.hasMinuteError = true
				return
			}
			if (minute < 0 || minute >= 60) {
				this.hasMinuteError = true
				return
			}
			this.hasMinuteError = false
			date.setMinutes(minute)

			// second
			if (!/^\d+$/.test(this.second)) {
				this.hasSecondError = true
				return
			}
			let second = parseInt(this.second)
			if (isNaN(second)) {
				this.hasSecondError = true
				return
			}
			if (second < 0 || second >= 60) {
				this.hasSecondError = true
				return
			}
			this.hasSecondError = false
			date.setSeconds(second)

			this.timestamp = Math.floor(date.getTime() / 1000)
		},
		leadingZero: function (s, l) {
			if (l <= s.length) {
				return s
			}
			for (let i = 0; i < l - s.length; i++) {
				s = "0" + s
			}
			return s
		}
	},
	template: `<div>
	<input type="hidden" :name="vName" :value="timestamp"/>
	<div class="ui fields inline" style="padding: 0; margin:0">
		<div class="ui field" :class="{error: hasDayError}">
			<input type="text" v-model="day" placeholder="YYYY-MM-DD" style="width:8.6em" maxlength="10" @input="change" ref="dayInput"/>
		</div>
		<div class="ui field" :class="{error: hasHourError}"><input type="text" v-model="hour" maxlength="2" style="width:4em" placeholder="时" @input="change"/></div>
		<div class="ui field">:</div>
		<div class="ui field" :class="{error: hasMinuteError}"><input type="text" v-model="minute" maxlength="2" style="width:4em" placeholder="分" @input="change"/></div>
		<div class="ui field">:</div>
		<div class="ui field" :class="{error: hasSecondError}"><input type="text" v-model="second" maxlength="2" style="width:4em" placeholder="秒" @input="change"/></div>
	</div>
</div>`
})

// 启用状态标签
Vue.component("label-on", {
	props: ["v-is-on"],
	template: '<div><span v-if="vIsOn" class="ui label tiny green basic">已启用</span><span v-if="!vIsOn" class="ui label tiny red basic">已停用</span></div>'
})

// 文字代码标签
Vue.component("code-label", {
	methods: {
		click: function (args) {
			this.$emit("click", args)
		}
	},
	template: `<span class="ui label basic tiny" style="padding: 3px;margin-left:2px;margin-right:2px" @click.prevent="click"><slot></slot></span>`
})

// tiny标签
Vue.component("tiny-label", {
	template: `<span class="ui label tiny" style="margin-bottom: 0.5em"><slot></slot></span>`
})

Vue.component("tiny-basic-label", {
	template: `<span class="ui label tiny basic" style="margin-bottom: 0.5em"><slot></slot></span>`
})

// 更小的标签
Vue.component("micro-basic-label", {
	template: `<span class="ui label tiny basic" style="margin-bottom: 0.5em; font-size: 0.7em; padding: 4px"><slot></slot></span>`
})


// 灰色的Label
Vue.component("grey-label", {
	template: `<span class="ui label basic grey tiny" style="margin-top: 0.4em; font-size: 0.7em; border: 1px solid #ddd!important; font-weight: normal;"><slot></slot></span>`
})


/**
 * 一级菜单
 */
Vue.component("first-menu", {
	props: [],
	template: ' \
		<div class="first-menu"> \
			<div class="ui menu text blue small">\
				<slot></slot>\
			</div> \
			<div class="ui divider"></div> \
		</div>'
});

/**
 * 更多选项
 */
Vue.component("more-options-indicator", {
	data: function () {
		return {
			visible: false
		}
	},
	methods: {
		changeVisible: function () {
			this.visible = !this.visible
			if (Tea.Vue != null) {
				Tea.Vue.moreOptionsVisible = this.visible
			}
			this.$emit("change", this.visible)
		}
	},
	template: '<a href="" style="font-weight: normal" @click.prevent="changeVisible()"><slot><span v-if="!visible">更多选项</span><span v-if="visible">收起选项</span></slot> <i class="icon angle" :class="{down:!visible, up:visible}"></i> </a>'
});

/**
 * 二级菜单
 */
Vue.component("second-menu", {
	template: ' \
		<div class="second-menu"> \
			<div class="ui menu text blue small">\
				<slot></slot>\
			</div> \
			<div class="ui divider"></div> \
		</div>'
});

Vue.component("more-options-angle", {
	data: function () {
		return {
			isVisible: false
		}
	},
	methods: {
		show: function () {
			this.isVisible = !this.isVisible
			this.$emit("change", this.isVisible)
		}
	},
	template: `<a href="" @click.prevent="show()"><span v-if="!isVisible">更多选项</span><span v-if="isVisible">收起选项</span><i class="icon angle" :class="{down:!isVisible, up:isVisible}"></i></a>`
})

/**
 * 菜单项
 */
Vue.component("inner-menu-item", {
	props: ["href", "active", "code"],
	data: function () {
		var active = this.active;
		if (typeof(active) =="undefined") {
			var itemCode = "";
			if (typeof (window.TEA.ACTION.data.firstMenuItem) != "undefined") {
				itemCode = window.TEA.ACTION.data.firstMenuItem;
			}
			active = (itemCode == this.code);
		}
		return {
			vHref: (this.href == null) ? "" : this.href,
			vActive: active
		};
	},
	template: '\
		<a :href="vHref" class="item right" style="color:#4183c4" :class="{active:vActive}">[<slot></slot>]</a> \
		'
});

Vue.component("health-check-config-box", {
	props: ["v-health-check-config"],
	data: function () {
		let healthCheckConfig = this.vHealthCheckConfig
		let urlProtocol = "http"
		let urlPort = ""
		let urlRequestURI = "/"
		let urlHost = ""

		if (healthCheckConfig == null) {
			healthCheckConfig = {
				isOn: false,
				url: "",
				interval: {count: 60, unit: "second"},
				statusCodes: [200],
				timeout: {count: 10, unit: "second"},
				countTries: 3,
				tryDelay: {count: 100, unit: "ms"},
				autoDown: true,
				countUp: 1,
				countDown: 3,
				userAgent: "",
				onlyBasicRequest: false
			}
			let that = this
			setTimeout(function () {
				that.changeURL()
			}, 500)
		} else {
			try {
				let url = new URL(healthCheckConfig.url)
				urlProtocol = url.protocol.substring(0, url.protocol.length - 1)

				// 域名
				urlHost = url.host
				if (urlHost == "%24%7Bhost%7D") {
					urlHost = "${host}"
				}
				let colonIndex = urlHost.indexOf(":")
				if (colonIndex > 0) {
					urlHost = urlHost.substring(0, colonIndex)
				}

				urlPort = url.port
				urlRequestURI = url.pathname
				if (url.search.length > 0) {
					urlRequestURI += url.search
				}
			} catch (e) {
			}

			if (healthCheckConfig.statusCodes == null) {
				healthCheckConfig.statusCodes = [200]
			}
			if (healthCheckConfig.interval == null) {
				healthCheckConfig.interval = {count: 60, unit: "second"}
			}
			if (healthCheckConfig.timeout == null) {
				healthCheckConfig.timeout = {count: 10, unit: "second"}
			}
			if (healthCheckConfig.tryDelay == null) {
				healthCheckConfig.tryDelay = {count: 100, unit: "ms"}
			}
			if (healthCheckConfig.countUp == null || healthCheckConfig.countUp < 1) {
				healthCheckConfig.countUp = 1
			}
			if (healthCheckConfig.countDown == null || healthCheckConfig.countDown < 1) {
				healthCheckConfig.countDown = 3
			}
		}
		return {
			healthCheck: healthCheckConfig,
			advancedVisible: false,
			urlProtocol: urlProtocol,
			urlHost: urlHost,
			urlPort: urlPort,
			urlRequestURI: urlRequestURI,
			urlIsEditing: healthCheckConfig.url.length == 0
		}
	},
	watch: {
		urlRequestURI: function () {
			if (this.urlRequestURI.length > 0 && this.urlRequestURI[0] != "/") {
				this.urlRequestURI = "/" + this.urlRequestURI
			}
			this.changeURL()
		},
		urlPort: function (v) {
			let port = parseInt(v)
			if (!isNaN(port)) {
				this.urlPort = port.toString()
			} else {
				this.urlPort = ""
			}
			this.changeURL()
		},
		urlProtocol: function () {
			this.changeURL()
		},
		urlHost: function () {
			this.changeURL()
		},
		"healthCheck.countTries": function (v) {
			let count = parseInt(v)
			if (!isNaN(count)) {
				this.healthCheck.countTries = count
			} else {
				this.healthCheck.countTries = 0
			}
		},
		"healthCheck.countUp": function (v) {
			let count = parseInt(v)
			if (!isNaN(count)) {
				this.healthCheck.countUp = count
			} else {
				this.healthCheck.countUp = 0
			}
		},
		"healthCheck.countDown": function (v) {
			let count = parseInt(v)
			if (!isNaN(count)) {
				this.healthCheck.countDown = count
			} else {
				this.healthCheck.countDown = 0
			}
		}
	},
	methods: {
		showAdvanced: function () {
			this.advancedVisible = !this.advancedVisible
		},
		changeURL: function () {
			let urlHost = this.urlHost
			if (urlHost.length == 0) {
				urlHost = "${host}"
			}
			this.healthCheck.url = this.urlProtocol + "://" + urlHost + ((this.urlPort.length > 0) ? ":" + this.urlPort : "") + this.urlRequestURI
		},
		changeStatus: function (values) {
			this.healthCheck.statusCodes = values.$map(function (k, v) {
				let status = parseInt(v)
				if (isNaN(status)) {
					return 0
				} else {
					return status
				}
			})
		},
		editURL: function () {
			this.urlIsEditing = !this.urlIsEditing
		}
	},
	template: `<div>
<input type="hidden" name="healthCheckJSON" :value="JSON.stringify(healthCheck)"/>
<table class="ui table definition selectable">
	<tbody>
		<tr>
			<td class="title">是否启用</td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" value="1" v-model="healthCheck.isOn"/>
					<label></label>
				</div>
			</td>
		</tr>
	</tbody>
	<tbody v-show="healthCheck.isOn">
		<tr>
			<td>URL *</td>
			<td>
				<div v-if="healthCheck.url.length > 0" style="margin-bottom: 1em"><code-label>{{healthCheck.url}}</code-label> &nbsp; <a href="" @click.prevent="editURL"><span class="small">修改 <i class="icon angle" :class="{down: !urlIsEditing, up: urlIsEditing}"></i></span></a> </div>
				<div v-show="urlIsEditing">
					<table class="ui table">
						 <tr>
							<td class="title">协议</td> 
							<td>
								<select class="ui dropdown auto-width" v-model="urlProtocol">
								<option value="http">http://</option>
								<option value="https">https://</option>
								</select>
							</td>
						</tr>
						<tr>
							<td>域名</td>
							<td>
								<input type="text" v-model="urlHost"/>
								<p class="comment">在此集群上可以访问到的一个域名。</p>
							</td>
						</tr>
						<tr>
							<td>端口</td>
							<td>
								<input type="text" maxlength="5" style="width:5.4em" placeholder="端口" v-model="urlPort"/>
							</td>
						</tr>
						<tr>
							<td>RequestURI</td>
							<td><input type="text" v-model="urlRequestURI" placeholder="/" style="width:20em"/></td>
						</tr>
					</table>
					<div class="ui divider"></div>
					<p class="comment" v-if="healthCheck.url.length > 0">拼接后的URL：<code-label>{{healthCheck.url}}</code-label>，其中\${host}指的是域名。</p>
				</div>
			</td>
		</tr>
		<tr>
			<td>检测时间间隔</td>
			<td>
				<time-duration-box :v-value="healthCheck.interval"></time-duration-box>
			</td>
		</tr>
		<tr>
			<td>是否自动下线</td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" value="1" v-model="healthCheck.autoDown"/>
					<label></label>
				</div>
				<p class="comment">选中后系统会根据健康检查的结果自动标记节点的上线/下线状态，并可能自动同步DNS设置。</p>
			</td>
		</tr>
		<tr v-show="healthCheck.autoDown">
			<td>连续上线次数</td>
			<td>
				<input type="text" v-model="healthCheck.countUp" style="width:5em" maxlength="6"/>
				<p class="comment">连续N次检查成功后自动恢复上线。</p>
			</td>
		</tr>
		<tr v-show="healthCheck.autoDown">
			<td>连续下线次数</td>
			<td>
				<input type="text" v-model="healthCheck.countDown" style="width:5em" maxlength="6"/>
				<p class="comment">连续N次检查失败后自动下线。</p>
			</td>
		</tr>
	</tbody>
	<tbody v-show="healthCheck.isOn">
		<tr>
			<td colspan="2"><more-options-angle @change="showAdvanced"></more-options-angle></td>
		</tr>
	</tbody>
	<tbody v-show="advancedVisible && healthCheck.isOn">
		<tr>
			<td>允许的状态码</td>
			<td>
				<values-box :values="healthCheck.statusCodes" maxlength="3" @change="changeStatus"></values-box>
			</td>
		</tr>
		<tr>
			<td>超时时间</td>
			<td>
				<time-duration-box :v-value="healthCheck.timeout"></time-duration-box>
			</td>	
		</tr>
		<tr>
			<td>连续尝试次数</td>
			<td>
				<input type="text" v-model="healthCheck.countTries" style="width: 5em" maxlength="2"/>
			</td>
		</tr>
		<tr>
			<td>每次尝试间隔</td>
			<td>
				<time-duration-box :v-value="healthCheck.tryDelay"></time-duration-box>
			</td>
		</tr>
		<tr>
			<td>终端信息<em>（User-Agent）</em></td>
			<td>
				<input type="text" v-model="healthCheck.userAgent" maxlength="200"/>
				<p class="comment">发送到服务器的User-Agent值，不填写表示使用默认值。</p>
			</td>
		</tr>
		<tr>
			<td>只基础请求</td>
			<td>
				<checkbox v-model="healthCheck.onlyBasicRequest"></checkbox>
				<p class="comment">只做基础的请求，不处理反向代理（不检查源站）、WAF等。</p>
			</td>
		</tr>
	</tbody>
</table>
<div class="margin"></div>
</div>`
})

Vue.component("time-duration-box", {
	props: ["v-name", "v-value", "v-count", "v-unit"],
	mounted: function () {
		this.change()
	},
	data: function () {
		let v = this.vValue
		if (v == null) {
			v = {
				count: this.vCount,
				unit: this.vUnit
			}
		}
		if (typeof (v["count"]) != "number") {
			v["count"] = -1
		}
		return {
			duration: v,
			countString: (v.count >= 0) ? v.count.toString() : ""
		}
	},
	watch: {
		"countString": function (newValue) {
			let value = newValue.trim()
			if (value.length == 0) {
				this.duration.count = -1
				return
			}
			let count = parseInt(value)
			if (!isNaN(count)) {
				this.duration.count = count
			}
			this.change()
		}
	},
	methods: {
		change: function () {
			this.$emit("change", this.duration)
		}
	},
	template: `<div class="ui fields inline" style="padding-bottom: 0; margin-bottom: 0">
	<input type="hidden" :name="vName" :value="JSON.stringify(duration)"/>
	<div class="ui field">
		<input type="text" v-model="countString" maxlength="11" size="11" @keypress.enter.prevent="1"/>
	</div>
	<div class="ui field">
		<select class="ui dropdown" v-model="duration.unit" @change="change">
			<option value="ms">毫秒</option>
			<option value="second">秒</option>
			<option value="minute">分钟</option>
			<option value="hour">小时</option>
			<option value="day">天</option>
		</select>
	</div>
</div>`
})

Vue.component("not-found-box", {
	props: ["message"],
	template: `<div style="text-align: center; margin-top: 5em;">
	<div style="font-size: 2em; margin-bottom: 1em"><i class="icon exclamation triangle large grey"></i></div>
	<p class="comment">{{message}}<slot></slot></p>
</div>`
})

// 警告消息
Vue.component("warning-message", {
	template: `<div class="ui icon message warning"><i class="icon warning circle"></i><div class="content"><slot></slot></div></div>`
})

let checkboxId = 0
Vue.component("checkbox", {
	props: ["name", "value", "v-value", "id", "checked"],
	data: function () {
		checkboxId++
		let elementId = this.id
		if (elementId == null) {
			elementId = "checkbox" + checkboxId
		}

		let elementValue = this.vValue
		if (elementValue == null) {
			elementValue = "1"
		}

		let checkedValue = this.value
        if (checkedValue == null && this.checked == "checked") {
            checkedValue = elementValue
        }

		return {
			elementId: elementId,
			elementValue: elementValue,
			newValue: checkedValue
		}
	},
	methods: {
		change: function () {
			this.$emit("input", this.newValue)
		}
	},
    watch: {
	    value: function (v) {
	        if (typeof v == "boolean") {
	            this.newValue = v
            }
        }
    },
	template: `<div class="ui checkbox">
	<input type="checkbox" :name="name" :value="elementValue" :id="elementId" @change="change" v-model="newValue"/>
	<label :for="elementId" style="font-size: 0.85em!important;"><slot></slot></label>
</div>`
})

Vue.component("network-addresses-view", {
	props: ["v-addresses"],
	template: `<div>
	<div class="ui label tiny basic" v-if="vAddresses != null" v-for="addr in vAddresses">
		{{addr.protocol}}://<span v-if="addr.host.length > 0">{{addr.host.quoteIP()}}</span><span v-else>*</span>:{{addr.portRange}}
	</div>
</div>`
})

Vue.component("size-capacity-view", {
	props:["v-default-text", "v-value"],
	template: `<div>
	<span v-if="vValue != null && vValue.count > 0">{{vValue.count}}{{vValue.unit.toUpperCase()}}</span>
	<span v-else>{{vDefaultText}}</span>
</div>`
})

// 信息提示窗口
Vue.component("tip-message-box", {
	props: ["code"],
	mounted: function () {
		let that = this
		Tea.action("/ui/showTip")
			.params({
				code: this.code
			})
			.success(function (resp) {
				that.visible = resp.data.visible
			})
			.post()
	},
	data: function () {
		return {
			visible: false
		}
	},
	methods: {
		close: function () {
			this.visible = false
			Tea.action("/ui/hideTip")
				.params({
					code: this.code
				})
				.post()
		}
	},
	template: `<div class="ui icon message" v-if="visible">
	<i class="icon info circle"></i>
	<i class="close icon" title="取消" @click.prevent="close" style="margin-top: 1em"></i>
	<div class="content">
		<slot></slot>
	</div>
</div>`
})

Vue.component("keyword", {
	props: ["v-word"],
	data: function () {
		let word = this.vWord
		if (word == null) {
			word = ""
		} else {
			word = word.replace(/\)/, "\\)")
			word = word.replace(/\(/, "\\(")
			word = word.replace(/\+/, "\\+")
			word = word.replace(/\^/, "\\^")
			word = word.replace(/\$/, "\\$")
		}

		let slot = this.$slots["default"][0]
		let text = this.encodeHTML(slot.text)
		if (word.length > 0) {
			text = text.replace(new RegExp("(" + word + ")", "ig"), "<span style=\"border: 1px #ccc dashed; color: #ef4d58\">$1</span>")
		}

		return {
			word: word,
			text: text
		}
	},
	methods: {
		encodeHTML: function (s) {
			s = s.replace("&", "&amp;")
			s = s.replace("<", "&lt;")
			s = s.replace(">", "&gt;")
			return s
		}
	},
	template: `<span><span style="display: none"><slot></slot></span><span v-html="text"></span></span>`
})

Vue.component("node-log-row", {
	props: ["v-log", "v-keyword"],
	data: function () {
		return {
			log: this.vLog,
			keyword: this.vKeyword
		}
	},
	template: `<div>
	<pre class="log-box" style="margin: 0; padding: 0"><span :class="{red:log.level == 'error', orange:log.level == 'warning', green: log.level == 'success'}"><span v-if="!log.isToday">[{{log.createdTime}}]</span><strong v-if="log.isToday">[{{log.createdTime}}]</strong><keyword :v-word="keyword">[{{log.tag}}]{{log.description}}</keyword></span> &nbsp; <span v-if="log.count > 1" class="ui label tiny" :class="{red:log.level == 'error', orange:log.level == 'warning'}">共{{log.count}}条</span> <span v-if="log.server != null && log.server.id > 0"><a :href="'/servers/server?serverId=' + log.server.id" class="ui label tiny basic">{{log.server.name}}</a></span></pre>
</div>`
})

Vue.component("provinces-selector", {
	props: ["v-provinces"],
	data: function () {
		let provinces = this.vProvinces
		if (provinces == null) {
			provinces = []
		}
		let provinceIds = provinces.$map(function (k, v) {
			return v.id
		})
		return {
			provinces: provinces,
			provinceIds: provinceIds
		}
	},
	methods: {
		add: function () {
			let provinceStringIds = this.provinceIds.map(function (v) {
				return v.toString()
			})
			let that = this
			teaweb.popup("/ui/selectProvincesPopup?provinceIds=" + provinceStringIds.join(","), {
				width: "48em",
				height: "23em",
				callback: function (resp) {
					that.provinces = resp.data.provinces
					that.change()
				}
			})
		},
		remove: function (index) {
			this.provinces.$remove(index)
			this.change()
		},
		change: function () {
			this.provinceIds = this.provinces.$map(function (k, v) {
				return v.id
			})
		}
	},
	template: `<div>
	<input type="hidden" name="provinceIdsJSON" :value="JSON.stringify(provinceIds)"/>
	<div v-if="provinces.length > 0" style="margin-bottom: 0.5em">
		<div v-for="(province, index) in provinces" class="ui label tiny basic">{{province.name}} <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove"></i></a></div>
		<div class="ui divider"></div>
	</div>
	<div>
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})

Vue.component("csrf-token", {
	created: function () {
		this.refreshToken()
	},
	mounted: function () {
		let that = this
		this.$refs.token.form.addEventListener("submit", function () {
			that.refreshToken()
		})

		// 自动刷新
		setInterval(function () {
			that.refreshToken()
		}, 10 * 60 * 1000)
	},
	data: function () {
		return {
			token: ""
		}
	},
	methods: {
		refreshToken: function () {
			let that = this
			Tea.action("/csrf/token")
				.get()
				.success(function (resp) {
					that.token = resp.data.token
				})
		}
	},
	template: `<input type="hidden" name="csrfToken" :value="token" ref="token"/>`
})


Vue.component("labeled-input", {
	props: ["name", "size", "maxlength", "label", "value"],
	template: '<div class="ui input right labeled"> \
	<input type="text" :name="name" :size="size" :maxlength="maxlength" :value="value"/>\
	<span class="ui label">{{label}}</span>\
</div>'
});

let radioId = 0
Vue.component("radio", {
	props: ["name", "value", "v-value", "id"],
	data: function () {
		radioId++
		let elementId = this.id
		if (elementId == null) {
			elementId = "radio" + radioId
		}
		return {
			"elementId": elementId
		}
	},
	methods: {
		change: function () {
			this.$emit("input", this.vValue)
		}
	},
	template: `<div class="ui checkbox radio">
	<input type="radio" :name="name" :value="vValue" :id="elementId" @change="change" :checked="(vValue == value)"/>
	<label :for="elementId"><slot></slot></label>
</div>`
})

Vue.component("copy-to-clipboard", {
	props: ["v-target"],
	created: function () {
		if (typeof ClipboardJS == "undefined") {
			let jsFile = document.createElement("script")
			jsFile.setAttribute("src", "/js/clipboard.min.js")
			document.head.appendChild(jsFile)
		}
	},
	methods: {
		copy: function () {
			new ClipboardJS('[data-clipboard-target]');
			teaweb.successToast("已复制到剪切板")
		}
	},
	template: `<a href="" title="拷贝到剪切板" :data-clipboard-target="'#' + vTarget" @click.prevent="copy"><i class="ui icon copy small"></i></em></a>`
})

// 节点角色名称
Vue.component("node-role-name", {
	props: ["v-role"],
	data: function () {
		let roleName = ""
		switch (this.vRole) {
			case "node":
				roleName = "边缘节点"
				break
			case "monitor":
				roleName = "监控节点"
				break
			case "api":
				roleName = "API节点"
				break
			case "user":
				roleName = "用户平台"
				break
			case "admin":
				roleName = "管理平台"
				break
			case "database":
				roleName = "数据库节点"
				break
			case "dns":
				roleName = "DNS节点"
				break
			case "report":
				roleName = "区域监控终端"
				break
		}
		return {
			roleName: roleName
		}
	},
	template: `<span>{{roleName}}</span>`
})

let sourceCodeBoxIndex = 0

Vue.component("source-code-box", {
	props: ["name", "type", "id", "read-only"],
	mounted: function () {
		let readOnly = this.readOnly
		if (typeof readOnly != "boolean") {
			readOnly = true
		}
		let box = document.getElementById("source-code-box-" + this.index)
		let valueBox = document.getElementById(this.valueBoxId)
		let value = ""
		if (valueBox.textContent != null) {
			value = valueBox.textContent
		} else if (valueBox.innerText != null) {
			value = valueBox.innerText
		}
		let boxEditor = CodeMirror.fromTextArea(box, {
			theme: "idea",
			lineNumbers: true,
			value: "",
			readOnly: readOnly,
			showCursorWhenSelecting: true,
			height: "auto",
			//scrollbarStyle: null,
			viewportMargin: Infinity,
			lineWrapping: true,
			highlightFormatting: false,
			indentUnit: 4,
			indentWithTabs: true
		})
		boxEditor.setValue(value)

		let info = CodeMirror.findModeByMIME(this.type)
		if (info != null) {
			boxEditor.setOption("mode", info.mode)
			CodeMirror.modeURL = "/codemirror/mode/%N/%N.js"
			CodeMirror.autoLoadMode(boxEditor, info.mode)
		}
	},
	data: function () {
		let index = sourceCodeBoxIndex++

		let valueBoxId = 'source-code-box-value-' + sourceCodeBoxIndex
		if (this.id != null) {
			valueBoxId = this.id
		}

		return {
			index: index,
			valueBoxId: valueBoxId
		}
	},
	template: `<div class="source-code-box">
	<div style="display: none" :id="valueBoxId"><slot></slot></div>
	<textarea :id="'source-code-box-' + index" :name="name"></textarea>
</div>`
})

Vue.component("size-capacity-box", {
	props: ["v-name", "v-value", "v-count", "v-unit", "size", "maxlength"],
	data: function () {
		let v = this.vValue
		if (v == null) {
			v = {
				count: this.vCount,
				unit: this.vUnit
			}
		}
		if (typeof (v["count"]) != "number") {
			v["count"] = -1
		}

		let vSize = this.size
		if (vSize == null) {
			vSize = 6
		}

		let vMaxlength = this.maxlength
		if (vMaxlength == null) {
			vMaxlength = 10
		}

		return {
			capacity: v,
			countString: (v.count >= 0) ? v.count.toString() : "",
			vSize: vSize,
			vMaxlength: vMaxlength
		}
	},
	watch: {
		"countString": function (newValue) {
			let value = newValue.trim()
			if (value.length == 0) {
				this.capacity.count = -1
				this.change()
				return
			}
			let count = parseInt(value)
			if (!isNaN(count)) {
				this.capacity.count = count
			}
			this.change()
		}
	},
	methods: {
		change: function () {
			this.$emit("change", this.capacity)
		}
	},
	template: `<div class="ui fields inline">
	<input type="hidden" :name="vName" :value="JSON.stringify(capacity)"/>
	<div class="ui field">
		<input type="text" v-model="countString" :maxlength="vMaxlength" :size="vSize"/>
	</div>
	<div class="ui field">
		<select class="ui dropdown" v-model="capacity.unit" @change="change">
			<option value="byte">字节</option>
			<option value="kb">KB</option>
			<option value="mb">MB</option>
			<option value="gb">GB</option>
			<option value="tb">TB</option>
			<option value="pb">PB</option>
		</select>
	</div>
</div>`
})

/**
 * 二级菜单
 */
Vue.component("inner-menu", {
	template: `
		<div class="second-menu" style="width:80%;position: absolute;top:-8px;right:1em"> 
			<div class="ui menu text blue small">
				<slot></slot>
			</div> 
		</div>`
});

Vue.component("datepicker", {
	props: ["v-name", "v-value", "v-bottom-left"],
	mounted: function () {
		let that = this
		teaweb.datepicker(this.$refs.dayInput, function (v) {
			that.day = v
			that.change()
		}, !!this.vBottomLeft)
	},
	data: function () {
		let name = this.vName
		if (name == null) {
			name = "day"
		}

		let day = this.vValue
		if (day == null) {
			day = ""
		}

		return {
			name: name,
			day: day
		}
	},
	methods: {
		change: function () {
			this.$emit("change", this.day)
		}
	},
	template: `<div style="display: inline-block">
	<input type="text" :name="name" v-model="day" placeholder="YYYY-MM-DD" style="width:8.6em" maxlength="10" @input="change" ref="dayInput" autocomplete="off"/>
</div>`
})

// 排序使用的箭头
Vue.component("sort-arrow", {
	props: ["name"],
	data: function () {
		let url = window.location.toString()
		let order = ""
		let newArgs = []
		if (window.location.search != null && window.location.search.length > 0) {
			let queryString = window.location.search.substring(1)
			let pieces = queryString.split("&")
			let that = this
			pieces.forEach(function (v) {
				let eqIndex = v.indexOf("=")
				if (eqIndex > 0) {
					let argName = v.substring(0, eqIndex)
					let argValue = v.substring(eqIndex + 1)
					if (argName == that.name) {
						order = argValue
					} else if (argValue != "asc" && argValue != "desc") {
						newArgs.push(v)
					}
				} else {
					newArgs.push(v)
				}
			})
		}
		if (order == "asc") {
			newArgs.push(this.name + "=desc")
		} else if (order == "desc") {
			newArgs.push(this.name + "=asc")
		} else {
			newArgs.push(this.name + "=desc")
		}

		let qIndex = url.indexOf("?")
		if (qIndex > 0) {
			url = url.substring(0, qIndex) + "?" + newArgs.join("&")
		} else {
			url = url + "?" + newArgs.join("&")
		}

		return {
			order: order,
			url: url
		}
	},
	template: `<a :href="url" title="排序">&nbsp; <i class="ui icon long arrow small" :class="{down: order == 'asc', up: order == 'desc', 'down grey': order == '' || order == null}"></i></a>`
})

Vue.component("user-link", {
	props: ["v-user", "v-keyword"],
	data: function () {
		let user = this.vUser
		if (user == null) {
			user = {id: 0, "username": "", "fullname": ""}
		}
		return {
			user: user
		}
	},
	template: `<div style="display: inline-block">
	<span v-if="user.id > 0"><keyword :v-word="vKeyword">{{user.fullname}}</keyword><span class="small grey">（<keyword :v-word="vKeyword">{{user.username}}</keyword>）</span></span>
	<span v-else class="disabled">[已删除]</span>
</div>`
})

// 监控节点分组选择
Vue.component("report-node-groups-selector", {
	props: ["v-group-ids"],
	mounted: function () {
		let that = this
		Tea.action("/clusters/monitors/groups/options")
			.post()
			.success(function (resp) {
				that.groups = resp.data.groups.map(function (group) {
					group.isChecked = that.groupIds.$contains(group.id)
					return group
				})
				that.isLoaded = true
			})
	},
	data: function () {
		var groupIds = this.vGroupIds
		if (groupIds == null) {
			groupIds = []
		}

		return {
			groups: [],
			groupIds: groupIds,
			isLoaded: false,
			allGroups: groupIds.length == 0
		}
	},
	methods: {
		check: function (group) {
			group.isChecked = !group.isChecked
			this.groupIds = []
			let that = this
			this.groups.forEach(function (v) {
				if (v.isChecked) {
					that.groupIds.push(v.id)
				}
			})
			this.change()
		},
		change: function () {
			let that = this
			let groups = []
			this.groupIds.forEach(function (groupId) {
				let group = that.groups.$find(function (k, v) {
					return v.id == groupId
				})
				if (group == null) {
					return
				}
				groups.push({
					id: group.id,
					name: group.name
				})
			})
			this.$emit("change", groups)
		}
	},
	watch: {
		allGroups: function (b) {
			if (b) {
				this.groupIds = []
				this.groups.forEach(function (v) {
					v.isChecked = false
				})
			}

			this.change()
		}
	},
	template: `<div>
	<input type="hidden" name="reportNodeGroupIdsJSON" :value="JSON.stringify(groupIds)"/>
	<span class="disabled" v-if="isLoaded && groups.length == 0">还没有分组。</span>
	<div v-if="groups.length > 0">
		<div>
			<div class="ui checkbox">
				<input type="checkbox" v-model="allGroups" id="all-group"/>
				<label for="all-group">全部分组</label>
			</div>
			<div class="ui divider" v-if="!allGroups"></div>
		</div>
		<div v-show="!allGroups">
			<div v-for="group in groups" :key="group.id" style="float: left; width: 7.6em; margin-bottom: 0.5em">
				<div class="ui checkbox">
					<input type="checkbox" v-model="group.isChecked" value="1" :id="'report-node-group-' + group.id" @click.prevent="check(group)"/>
					<label :for="'report-node-group-' + group.id">{{group.name}}</label>
				</div>
			</div>
		</div>
	</div>
</div>`
})

Vue.component("finance-user-selector", {
	mounted: function () {
		let that = this

		Tea.action("/finance/users/options")
			.post()
			.success(function (resp) {
				that.users = resp.data.users
			})
	},
	props: ["v-user-id"],
	data: function () {
		let userId = this.vUserId
		if (userId == null) {
			userId = 0
		}
		return {
			users: [],
			userId: userId
		}
	},
	watch: {
		userId: function (v) {
			this.$emit("change", v)
		}
	},
	template: `<div>
	<select class="ui dropdown auto-width" name="userId" v-model="userId">
		<option value="0">[选择用户]</option>
		<option v-for="user in users" :value="user.id">{{user.fullname}} ({{user.username}})</option>
	</select>
</div>`
})

// 节点登录推荐端口
Vue.component("node-login-suggest-ports", {
	data: function () {
		return {
			ports: [],
			availablePorts: [],
			autoSelected: false,
			isLoading: false
		}
	},
	methods: {
		reload: function (host) {
			let that = this
			this.autoSelected = false
			this.isLoading = true
			Tea.action("/clusters/cluster/suggestLoginPorts")
				.params({
					host: host
				})
				.success(function (resp) {
					if (resp.data.availablePorts != null) {
						that.availablePorts = resp.data.availablePorts
						if (that.availablePorts.length > 0) {
							that.autoSelectPort(that.availablePorts[0])
							that.autoSelected = true
						}
					}
					if (resp.data.ports != null) {
						that.ports = resp.data.ports
						if (that.ports.length > 0 && !that.autoSelected) {
							that.autoSelectPort(that.ports[0])
							that.autoSelected = true
						}
					}
				})
				.done(function () {
					that.isLoading = false
				})
				.post()
		},
		selectPort: function (port) {
			this.$emit("select", port)
		},
		autoSelectPort: function (port) {
			this.$emit("auto-select", port)
		}
	},
	template: `<span>
	<span v-if="isLoading">正在检查端口...</span>
	<span v-if="availablePorts.length > 0">
		可能端口：<a href="" v-for="port in availablePorts" @click.prevent="selectPort(port)" class="ui label tiny basic blue" style="border: 1px #2185d0 dashed; font-weight: normal">{{port}}</a>
		&nbsp; &nbsp;
	</span>
	<span v-if="ports.length > 0">
		常用端口：<a href="" v-for="port in ports" @click.prevent="selectPort(port)" class="ui label tiny basic blue" style="border: 1px #2185d0 dashed;  font-weight: normal">{{port}}</a>
	</span>
	<span v-if="ports.length == 0">常用端口有22等。</span>
	<span v-if="ports.length > 0" class="grey small">（可以点击要使用的端口）</span>
</span>`
})

Vue.component("node-group-selector", {
	props: ["v-cluster-id", "v-group"],
	data: function () {
		return {
			selectedGroup: this.vGroup
		}
	},
	methods: {
		selectGroup: function () {
			let that = this
			teaweb.popup("/clusters/cluster/groups/selectPopup?clusterId=" + this.vClusterId, {
				callback: function (resp) {
					that.selectedGroup = resp.data.group
				}
			})
		},
		addGroup: function () {
			let that = this
			teaweb.popup("/clusters/cluster/groups/createPopup?clusterId=" + this.vClusterId, {
				callback: function (resp) {
					that.selectedGroup = resp.data.group
				}
			})
		},
		removeGroup: function () {
			this.selectedGroup = null
		}
	},
	template: `<div>
	<div class="ui label small basic" v-if="selectedGroup != null">
		<input type="hidden" name="groupId" :value="selectedGroup.id"/>
		{{selectedGroup.name}} &nbsp;<a href="" title="删除" @click.prevent="removeGroup()"><i class="icon remove"></i></a>
	</div>
	<div v-if="selectedGroup == null">
		<a href="" @click.prevent="selectGroup()">[选择分组]</a> &nbsp; <a href="" @click.prevent="addGroup()">[添加分组]</a>
	</div>
</div>`
})

// 节点IP地址管理（标签形式）
Vue.component("node-ip-addresses-box", {
	props: ["v-ip-addresses", "role"],
	data: function () {
		return {
			ipAddresses: (this.vIpAddresses == null) ? [] : this.vIpAddresses,
			supportThresholds: this.role != "ns"
		}
	},
	methods: {
		// 添加IP地址
		addIPAddress: function () {
			window.UPDATING_NODE_IP_ADDRESS = null

			let that = this;
			teaweb.popup("/nodes/ipAddresses/createPopup?supportThresholds=" + (this.supportThresholds ? 1 : 0), {
				callback: function (resp) {
					that.ipAddresses.push(resp.data.ipAddress);
				},
				height: "24em",
				width: "44em"
			})
		},

		// 修改地址
		updateIPAddress: function (index, address) {
			window.UPDATING_NODE_IP_ADDRESS = address

			let that = this;
			teaweb.popup("/nodes/ipAddresses/updatePopup?supportThresholds=" + (this.supportThresholds ? 1 : 0), {
				callback: function (resp) {
					Vue.set(that.ipAddresses, index, resp.data.ipAddress);
				},
				height: "24em",
				width: "44em"
			})
		},

		// 删除IP地址
		removeIPAddress: function (index) {
			this.ipAddresses.$remove(index);
		},

		// 判断是否为IPv6
		isIPv6: function (ip) {
			return ip.indexOf(":") > -1
		}
	},
	template: `<div>
	<input type="hidden" name="ipAddressesJSON" :value="JSON.stringify(ipAddresses)"/>
	<div v-if="ipAddresses.length > 0">
		<div>
			<div v-for="(address, index) in ipAddresses" class="ui label tiny basic">
				<span v-if="isIPv6(address.ip)" class="grey">[IPv6]</span> {{address.ip}}
				<span class="small grey" v-if="address.name.length > 0">（{{address.name}}<span v-if="!address.canAccess">，不可访问</span>）</span>
				<span class="small grey" v-if="address.name.length == 0 && !address.canAccess">（不可访问）</span>
				<span class="small red" v-if="!address.isOn" title="未启用">[off]</span>
				<span class="small red" v-if="!address.isUp" title="已下线">[down]</span>
				<span class="small" v-if="address.thresholds != null && address.thresholds.length > 0">[{{address.thresholds.length}}个阈值]</span>
				&nbsp;
				<a href="" title="修改" @click.prevent="updateIPAddress(index, address)"><i class="icon pencil small"></i></a>
				<a href="" title="删除" @click.prevent="removeIPAddress(index)"><i class="icon remove"></i></a>
			</div>
		</div>
		<div class="ui divider"></div>
	</div>
	<div>
		<button class="ui button small" type="button" @click.prevent="addIPAddress()">+</button>
	</div>
</div>`
})

// 节点IP阈值
Vue.component("node-ip-address-thresholds-view", {
	props: ["v-thresholds"],
	data: function () {
		let thresholds = this.vThresholds
		if (thresholds == null) {
			thresholds = []
		} else {
			thresholds.forEach(function (v) {
				if (v.items == null) {
					v.items = []
				}
				if (v.actions == null) {
					v.actions = []
				}
			})
		}

		return {
			thresholds: thresholds,
			allItems: window.IP_ADDR_THRESHOLD_ITEMS,
			allOperators: [
				{
					"name": "小于等于",
					"code": "lte"
				},
				{
					"name": "大于",
					"code": "gt"
				},
				{
					"name": "不等于",
					"code": "neq"
				},
				{
					"name": "小于",
					"code": "lt"
				},
				{
					"name": "大于等于",
					"code": "gte"
				}
			],
			allActions: window.IP_ADDR_THRESHOLD_ACTIONS
		}
	},
	methods: {
		itemName: function (item) {
			let result = ""
			this.allItems.forEach(function (v) {
				if (v.code == item) {
					result = v.name
				}
			})
			return result
		},
		itemUnitName: function (itemCode) {
			let result = ""
			this.allItems.forEach(function (v) {
				if (v.code == itemCode) {
					result = v.unit
				}
			})
			return result
		},
		itemDurationUnitName: function (unit) {
			switch (unit) {
				case "minute":
					return "分钟"
				case "second":
					return "秒"
				case "hour":
					return "小时"
				case "day":
					return "天"
			}
			return unit
		},
		itemOperatorName: function (operator) {
			let result = ""
			this.allOperators.forEach(function (v) {
				if (v.code == operator) {
					result = v.name
				}
			})
			return result
		},
		actionName: function (actionCode) {
			let result = ""
			this.allActions.forEach(function (v) {
				if (v.code == actionCode) {
					result = v.name
				}
			})
			return result
		}
	},
	template: `<div>
	<!-- 已有条件 -->
	<div v-if="thresholds.length > 0">
		<div class="ui label basic small" v-for="(threshold, index) in thresholds" style="margin-bottom: 0.8em">
			<span v-for="(item, itemIndex) in threshold.items">
				<span>
					<span v-if="item.item != 'nodeHealthCheck'">
						[{{item.duration}}{{itemDurationUnitName(item.durationUnit)}}]
					</span>	 
					{{itemName(item.item)}}
					
					<span v-if="item.item == 'nodeHealthCheck'">
						<!-- 健康检查 -->
						<span v-if="item.value == 1">成功</span>
						<span v-if="item.value == 0">失败</span>
					</span>
					<span v-else>
						<!-- 连通性 -->
						<span v-if="item.item == 'connectivity' && item.options != null && item.options.groups != null && item.options.groups.length > 0">[<span v-for="(group, groupIndex) in item.options.groups">{{group.name}} <span v-if="groupIndex != item.options.groups.length - 1">&nbsp; </span></span>]</span>
						
						 <span class="grey">[{{itemOperatorName(item.operator)}}]</span> {{item.value}}{{itemUnitName(item.item)}} &nbsp;
					 </span>
				 </span>
				 <span v-if="itemIndex != threshold.items.length - 1" style="font-style: italic">AND &nbsp;</span></span>
				-&gt;
				<span v-for="(action, actionIndex) in threshold.actions">{{actionName(action.action)}}
				<span v-if="action.action == 'switch'">到{{action.options.ips.join(", ")}}</span>
				<span v-if="action.action == 'webHook'" class="small grey">({{action.options.url}})</span>
				 &nbsp;					 
				 <span v-if="actionIndex != threshold.actions.length - 1" style="font-style: italic">AND &nbsp;</span>
			 </span>
		</div>
	</div>
</div>`
})

// 节点IP阈值
Vue.component("node-ip-address-thresholds-box", {
	props: ["v-thresholds"],
	data: function () {
		let thresholds = this.vThresholds
		if (thresholds == null) {
			thresholds = []
		} else {
			thresholds.forEach(function (v) {
				if (v.items == null) {
					v.items = []
				}
				if (v.actions == null) {
					v.actions = []
				}
			})
		}

		return {
			editingIndex: -1,
			thresholds: thresholds,
			addingThreshold: {
				items: [],
				actions: []
			},
			isAdding: false,
			isAddingItem: false,
			isAddingAction: false,

			itemCode: "nodeAvgRequests",
			itemReportGroups: [],
			itemOperator: "lte",
			itemValue: "",
			itemDuration: "5",
			allItems: window.IP_ADDR_THRESHOLD_ITEMS,
			allOperators: [
				{
					"name": "小于等于",
					"code": "lte"
				},
				{
					"name": "大于",
					"code": "gt"
				},
				{
					"name": "不等于",
					"code": "neq"
				},
				{
					"name": "小于",
					"code": "lt"
				},
				{
					"name": "大于等于",
					"code": "gte"
				}
			],
			allActions: window.IP_ADDR_THRESHOLD_ACTIONS,

			actionCode: "up",
			actionBackupIPs: "",
			actionWebHookURL: ""
		}
	},
	methods: {
		add: function () {
			this.isAdding = !this.isAdding
		},
		cancel: function () {
			this.isAdding = false
			this.editingIndex = -1
			this.addingThreshold = {
				items: [],
				actions: []
			}
		},
		confirm: function () {
			if (this.addingThreshold.items.length == 0) {
				teaweb.warn("需要至少添加一个阈值")
				return
			}
			if (this.addingThreshold.actions.length == 0) {
				teaweb.warn("需要至少添加一个动作")
				return
			}

			if (this.editingIndex >= 0) {
				this.thresholds[this.editingIndex].items = this.addingThreshold.items
				this.thresholds[this.editingIndex].actions = this.addingThreshold.actions
			} else {
				this.thresholds.push({
					items: this.addingThreshold.items,
					actions: this.addingThreshold.actions
				})
			}

			// 还原
			this.cancel()
		},
		remove: function (index) {
			this.cancel()
			this.thresholds.$remove(index)
		},
		update: function (index) {
			this.editingIndex = index
			this.addingThreshold = {
				items: this.thresholds[index].items.$copy(),
				actions: this.thresholds[index].actions.$copy()
			}
			this.isAdding = true
		},

		addItem: function () {
			this.isAddingItem = !this.isAddingItem
			let that = this
			setTimeout(function () {
				that.$refs.itemValue.focus()
			}, 100)
		},
		cancelItem: function () {
			this.isAddingItem = false

			this.itemCode = "nodeAvgRequests"
			this.itmeOperator = "lte"
			this.itemValue = ""
			this.itemDuration = "5"
			this.itemReportGroups = []
		},
		confirmItem: function () {
			// 特殊阈值快速添加
			if (["nodeHealthCheck"].$contains(this.itemCode)) {
				if (this.itemValue.toString().length == 0) {
					teaweb.warn("请选择检查结果")
					return
				}

				let value = parseInt(this.itemValue)
				if (isNaN(value)) {
					value = 0
				} else if (value < 0) {
					value = 0
				} else if (value > 1) {
					value = 1
				}

				// 添加
				this.addingThreshold.items.push({
					item: this.itemCode,
					operator: this.itemOperator,
					value: value,
					duration: 0,
					durationUnit: "minute",
					options: {}
				})
				this.cancelItem()
				return
			}

			if (this.itemDuration.length == 0) {
				let that = this
				teaweb.warn("请输入统计周期", function () {
					that.$refs.itemDuration.focus()
				})
				return
			}
			let itemDuration = parseInt(this.itemDuration)
			if (isNaN(itemDuration) || itemDuration <= 0) {
				teaweb.warn("请输入正确的统计周期", function () {
					that.$refs.itemDuration.focus()
				})
				return
			}

			if (this.itemValue.length == 0) {
				let that = this
				teaweb.warn("请输入对比值", function () {
					that.$refs.itemValue.focus()
				})
				return
			}
			let itemValue = parseFloat(this.itemValue)
			if (isNaN(itemValue)) {
				teaweb.warn("请输入正确的对比值", function () {
					that.$refs.itemValue.focus()
				})
				return
			}


			let options = {}

			switch (this.itemCode) {
				case "connectivity": // 连通性校验
					if (itemValue > 100) {
						let that = this
						teaweb.warn("连通性对比值不能超过100", function () {
							that.$refs.itemValue.focus()
						})
						return
					}

					options["groups"] = this.itemReportGroups
					break
			}

			// 添加
			this.addingThreshold.items.push({
				item: this.itemCode,
				operator: this.itemOperator,
				value: itemValue,
				duration: itemDuration,
				durationUnit: "minute",
				options: options
			})

			// 还原
			this.cancelItem()
		},
		removeItem: function (index) {
			this.cancelItem()
			this.addingThreshold.items.$remove(index)
		},
		changeReportGroups: function (groups) {
			this.itemReportGroups = groups
		},
		itemName: function (item) {
			let result = ""
			this.allItems.forEach(function (v) {
				if (v.code == item) {
					result = v.name
				}
			})
			return result
		},
		itemUnitName: function (itemCode) {
			let result = ""
			this.allItems.forEach(function (v) {
				if (v.code == itemCode) {
					result = v.unit
				}
			})
			return result
		},
		itemDurationUnitName: function (unit) {
			switch (unit) {
				case "minute":
					return "分钟"
				case "second":
					return "秒"
				case "hour":
					return "小时"
				case "day":
					return "天"
			}
			return unit
		},
		itemOperatorName: function (operator) {
			let result = ""
			this.allOperators.forEach(function (v) {
				if (v.code == operator) {
					result = v.name
				}
			})
			return result
		},

		addAction: function () {
			this.isAddingAction = !this.isAddingAction
		},
		cancelAction: function () {
			this.isAddingAction = false
			this.actionCode = "up"
			this.actionBackupIPs = ""
			this.actionWebHookURL = ""
		},
		confirmAction: function () {
			this.doConfirmAction(false)
		},
		doConfirmAction: function (validated, options) {
			// 是否已存在
			let exists = false
			let that = this
			this.addingThreshold.actions.forEach(function (v) {
				if (v.action == that.actionCode) {
					exists = true
				}
			})
			if (exists) {
				teaweb.warn("此动作已经添加过了，无需重复添加")
				return
			}

			if (options == null) {
				options = {}
			}

			switch (this.actionCode) {
				case "switch":
					if (!validated) {
						Tea.action("/ui/validateIPs")
							.params({
								"ips": this.actionBackupIPs
							})
							.success(function (resp) {
								if (resp.data.ips.length == 0) {
									teaweb.warn("请输入备用IP", function () {
										that.$refs.actionBackupIPs.focus()
									})
									return
								}
								options["ips"] = resp.data.ips
								that.doConfirmAction(true, options)
							})
							.fail(function (resp) {
								teaweb.warn("输入的IP '" + resp.data.failIP + "' 格式不正确，请改正后提交", function () {
									that.$refs.actionBackupIPs.focus()
								})
							})
							.post()
						return
					}
					break
				case "webHook":
					if (this.actionWebHookURL.length == 0) {
						teaweb.warn("请输入WebHook URL", function () {
							that.$refs.webHookURL.focus()
						})
						return
					}
					if (!this.actionWebHookURL.match(/^(http|https):\/\//i)) {
						teaweb.warn("URL开头必须是http://或者https://", function () {
							that.$refs.webHookURL.focus()
						})
						return
					}
					options["url"] = this.actionWebHookURL
			}

			this.addingThreshold.actions.push({
				action: this.actionCode,
				options: options
			})

			// 还原
			this.cancelAction()
		},
		removeAction: function (index) {
			this.cancelAction()
			this.addingThreshold.actions.$remove(index)
		},
		actionName: function (actionCode) {
			let result = ""
			this.allActions.forEach(function (v) {
				if (v.code == actionCode) {
					result = v.name
				}
			})
			return result
		}
	},
	template: `<div>
	<input type="hidden" name="thresholdsJSON" :value="JSON.stringify(thresholds)"/>
		
	<!-- 已有条件 -->
	<div v-if="thresholds.length > 0">
		<div class="ui label basic small" v-for="(threshold, index) in thresholds">
			<span v-for="(item, itemIndex) in threshold.items">
				<span v-if="item.item != 'nodeHealthCheck'">
					[{{item.duration}}{{itemDurationUnitName(item.durationUnit)}}]
				</span> 
				{{itemName(item.item)}}
				
				<span v-if="item.item == 'nodeHealthCheck'">
					<!-- 健康检查 -->
					<span v-if="item.value == 1">成功</span>
					<span v-if="item.value == 0">失败</span>
				</span>
				<span v-else>
					<!-- 连通性 -->
					<span v-if="item.item == 'connectivity' && item.options != null && item.options.groups != null && item.options.groups.length > 0">[<span v-for="(group, groupIndex) in item.options.groups">{{group.name}} <span v-if="groupIndex != item.options.groups.length - 1">&nbsp; </span></span>]</span>
				
					<span  class="grey">[{{itemOperatorName(item.operator)}}]</span> &nbsp;{{item.value}}{{itemUnitName(item.item)}} 
			 	</span>
			 	&nbsp;<span v-if="itemIndex != threshold.items.length - 1" style="font-style: italic">AND &nbsp;</span>
			</span>
			-&gt;
			<span v-for="(action, actionIndex) in threshold.actions">{{actionName(action.action)}}
			<span v-if="action.action == 'switch'">到{{action.options.ips.join(", ")}}</span>
			<span v-if="action.action == 'webHook'" class="small grey">({{action.options.url}})</span>
			 &nbsp;<span v-if="actionIndex != threshold.actions.length - 1" style="font-style: italic">AND &nbsp;</span></span>
			&nbsp;
			<a href="" title="修改" @click.prevent="update(index)"><i class="icon pencil small"></i></a> 
			<a href="" title="删除" @click.prevent="remove(index)"><i class="icon small remove"></i></a>
		</div>
	</div>
	
	<!-- 新阈值 -->
	<div v-if="isAdding" style="margin-top: 0.5em">
		<table class="ui table celled">
			<thead>
				<tr>
					<td style="width: 50%; background: #f9fafb; border-bottom: 1px solid rgba(34,36,38,.1)">阈值</td>
					<th>动作</th>
				</tr>
			</thead>
			<tr>
				<td style="background: white">
					<!-- 已经添加的项目 -->
					<div>
						<div v-for="(item, index) in addingThreshold.items" class="ui label basic small" style="margin-bottom: 0.5em;">
							<span v-if="item.item != 'nodeHealthCheck'">
								[{{item.duration}}{{itemDurationUnitName(item.durationUnit)}}]
							</span> 
							{{itemName(item.item)}}
							
							<span v-if="item.item == 'nodeHealthCheck'">
								<!-- 健康检查 -->
								<span v-if="item.value == 1">成功</span>
								<span v-if="item.value == 0">失败</span>
							</span>
							<span v-else>
								<!-- 连通性 -->
								<span v-if="item.item == 'connectivity' && item.options != null && item.options.groups != null && item.options.groups.length > 0">[<span v-for="(group, groupIndex) in item.options.groups">{{group.name}} <span v-if="groupIndex != item.options.groups.length - 1">&nbsp; </span></span>]</span>
								 <span class="grey">[{{itemOperatorName(item.operator)}}]</span> {{item.value}}{{itemUnitName(item.item)}}
							 </span> 
							 &nbsp;
							<a href="" title="删除" @click.prevent="removeItem(index)"><i class="icon remove small"></i></a>
						</div>
					</div>
					
					<!-- 正在添加的项目 -->
					<div v-if="isAddingItem" style="margin-top: 0.8em">
						<table class="ui table">
							<tr>
								<td style="width: 6em">统计项目</td>
								<td>
									<select class="ui dropdown auto-width" v-model="itemCode">
									<option v-for="item in allItems" :value="item.code">{{item.name}}</option>
									</select>
									<p class="comment" style="font-weight: normal" v-for="item in allItems" v-if="item.code == itemCode">{{item.description}}</p>
								</td>
							</tr>
							<tr v-show="itemCode != 'nodeHealthCheck'">
								<td>统计周期</td>
								<td>
									<div class="ui input right labeled">
										<input type="text" v-model="itemDuration" style="width: 4em" maxlength="4" ref="itemDuration" @keyup.enter="confirmItem()" @keypress.enter.prevent="1"/>
										<span class="ui label">分钟</span>
									</div>
								</td>
							</tr>
							<tr v-show="itemCode != 'nodeHealthCheck'">
								<td>操作符</td>
								<td>
									<select class="ui dropdown auto-width" v-model="itemOperator">
										<option v-for="operator in allOperators" :value="operator.code">{{operator.name}}</option>
									</select>
								</td>
							</tr>
							<tr v-show="itemCode != 'nodeHealthCheck'">
								<td>对比值</td>
								<td>
									<div class="ui input right labeled">
										<input type="text" maxlength="20" style="width: 5em" v-model="itemValue" ref="itemValue" @keyup.enter="confirmItem()" @keypress.enter.prevent="1"/>
										<span class="ui label" v-for="item in allItems" v-if="item.code == itemCode">{{item.unit}}</span>
									</div>
								</td>
							</tr>
							<tr v-show="itemCode == 'nodeHealthCheck'">
								<td>检查结果</td>
								<td>
									<select class="ui dropdown auto-width" v-model="itemValue">
										<option value="1">成功</option>
										<option value="0">失败</option>
									</select>
									<p class="comment" style="font-weight: normal">只有状态发生改变的时候才会触发。</p>
								</td>
							</tr>
							
							<!-- 连通性 -->
							<tr v-if="itemCode == 'connectivity'">
								<td>终端分组</td>
								<td style="font-weight: normal">
									<div style="zoom: 0.8"><report-node-groups-selector @change="changeReportGroups"></report-node-groups-selector></div>
								</td>
							</tr>
						</table>
						<div style="margin-top: 0.8em">
							<button class="ui button tiny" type="button" @click.prevent="confirmItem">确定</button>							 &nbsp;
							<a href="" title="取消" @click.prevent="cancelItem"><i class="icon remove small"></i></a>
						</div>
					</div>
					<div style="margin-top: 0.8em" v-if="!isAddingItem">
						<button class="ui button tiny" type="button" @click.prevent="addItem">+</button>
					</div>
				</td>
				<td style="background: white">
					<!-- 已经添加的动作 -->
					<div>
						<div v-for="(action, index) in addingThreshold.actions" class="ui label basic small" style="margin-bottom: 0.5em">
							{{actionName(action.action)}} &nbsp;
							<span v-if="action.action == 'switch'">到{{action.options.ips.join(", ")}}</span>
							<span v-if="action.action == 'webHook'" class="small grey">({{action.options.url}})</span>
							<a href="" title="删除" @click.prevent="removeAction(index)"><i class="icon remove small"></i></a>
						</div>
					</div>
					
					<!-- 正在添加的动作 -->
					<div v-if="isAddingAction" style="margin-top: 0.8em">
						<table class="ui table">
							<tr>
								<td style="width: 6em">动作类型</td>
								<td>
									<select class="ui dropdown auto-width" v-model="actionCode">
										<option v-for="action in allActions" :value="action.code">{{action.name}}</option>
									</select>
									<p class="comment" v-for="action in allActions" v-if="action.code == actionCode">{{action.description}}</p>
								</td>
							</tr>
							
							<!-- 切换 -->
							<tr v-if="actionCode == 'switch'">
								<td>备用IP *</td>
								<td>
									<textarea rows="2" v-model="actionBackupIPs" ref="actionBackupIPs"></textarea>
									<p class="comment">每行一个备用IP。</p>
								</td>
							</tr>
							
							<!-- WebHook -->
							<tr v-if="actionCode == 'webHook'">
								<td>URL *</td>
								<td>
									<input type="text" maxlength="1000" placeholder="https://..." v-model="actionWebHookURL" ref="webHookURL" @keyup.enter="confirmAction()" @keypress.enter.prevent="1"/>
									<p class="comment">完整的URL，比如<code-label>https://example.com/webhook/api</code-label>，系统会在触发阈值的时候通过GET调用此URL。</p>
								</td>
							</tr>
						</table>
						<div style="margin-top: 0.8em">
							<button class="ui button tiny" type="button" @click.prevent="confirmAction">确定</button>	 &nbsp;
							<a href="" title="取消" @click.prevent="cancelAction"><i class="icon remove small"></i></a>
						</div>
					</div>
					
					<div style="margin-top: 0.8em" v-if="!isAddingAction">
						<button class="ui button tiny" type="button" @click.prevent="addAction">+</button>
					</div>	
				</td>
			</tr>
		</table>
		
		<!-- 添加阈值 -->
		<div>
			<button class="ui button tiny" :class="{disabled: (isAddingItem || isAddingAction)}" type="button" @click.prevent="confirm">确定</button> &nbsp;
			<a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
		</div>
	</div>
	
	<div v-if="!isAdding" style="margin-top: 0.5em">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})

Vue.component("node-region-selector", {
	props: ["v-region"],
	data: function () {
		return {
			selectedRegion: this.vRegion
		}
	},
	methods: {
		selectRegion: function () {
			let that = this
			teaweb.popup("/clusters/regions/selectPopup?clusterId=" + this.vClusterId, {
				callback: function (resp) {
					that.selectedRegion = resp.data.region
				}
			})
		},
		addRegion: function () {
			let that = this
			teaweb.popup("/clusters/regions/createPopup?clusterId=" + this.vClusterId, {
				callback: function (resp) {
					that.selectedRegion = resp.data.region
				}
			})
		},
		removeRegion: function () {
			this.selectedRegion = null
		}
	},
	template: `<div>
	<div class="ui label small basic" v-if="selectedRegion != null">
		<input type="hidden" name="regionId" :value="selectedRegion.id"/>
		{{selectedRegion.name}} &nbsp;<a href="" title="删除" @click.prevent="removeRegion()"><i class="icon remove"></i></a>
	</div>
	<div v-if="selectedRegion == null">
		<a href="" @click.prevent="selectRegion()">[选择区域]</a> &nbsp; <a href="" @click.prevent="addRegion()">[添加区域]</a>
	</div>
</div>`
})

Vue.component("dns-route-selector", {
	props: ["v-all-routes", "v-routes"],
	data: function () {
		let routes = this.vRoutes
		if (routes == null) {
			routes = []
		}
		routes.$sort(function (v1, v2) {
			if (v1.domainId == v2.domainId) {
				return v1.code < v2.code
			}
			return (v1.domainId < v2.domainId) ? 1 : -1
		})
		return {
			routes: routes,
			routeCodes: routes.$map(function (k, v) {
				return v.code + "@" + v.domainId
			}),
			isAdding: false,
			routeCode: "",
			keyword: "",
			searchingRoutes: this.vAllRoutes.$copy()
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			this.keyword = ""
			this.routeCode = ""

			let that = this
			setTimeout(function () {
				that.$refs.keywordRef.focus()
			}, 200)
		},
		cancel: function () {
			this.isAdding = false
		},
		confirm: function () {
			if (this.routeCode.length == 0) {
				return
			}
			if (this.routeCodes.$contains(this.routeCode)) {
				teaweb.warn("已经添加过此线路，不能重复添加")
				return
			}
			let that = this
			let route = this.vAllRoutes.$find(function (k, v) {
				return v.code + "@" + v.domainId == that.routeCode
			})
			if (route == null) {
				return
			}

			this.routeCodes.push(this.routeCode)
			this.routes.push(route)

			this.routes.$sort(function (v1, v2) {
				if (v1.domainId == v2.domainId) {
					return v1.code < v2.code
				}
				return (v1.domainId < v2.domainId) ? 1 : -1
			})

			this.routeCode = ""
			this.isAdding = false
		},
		remove: function (route) {
			this.routeCodes.$removeValue(route.code + "@" + route.domainId)
			this.routes.$removeIf(function (k, v) {
				return v.code + "@" + v.domainId == route.code + "@" + route.domainId
			})
		}
	},
	watch: {
		keyword: function (keyword) {
			if (keyword.length == 0) {
				this.searchingRoutes = this.vAllRoutes.$copy()
				this.routeCode = ""
				return
			}
			this.searchingRoutes = this.vAllRoutes.filter(function (route) {
				return teaweb.match(route.name, keyword) || teaweb.match(route.domainName, keyword)
			})
			if (this.searchingRoutes.length > 0) {
				this.routeCode = this.searchingRoutes[0].code + "@" + this.searchingRoutes[0].domainId
			} else {
				this.routeCode = ""
			}
		}
	},
	template: `<div>
	<input type="hidden" name="dnsRoutesJSON" :value="JSON.stringify(routeCodes)"/>
	<div v-if="routes.length > 0">
		<tiny-basic-label v-for="route in routes" :key="route.code + '@' + route.domainId">
			{{route.name}} <span class="grey small">（{{route.domainName}}）</span><a href="" @click.prevent="remove(route)"><i class="icon remove"></i></a>
		</tiny-basic-label>
		<div class="ui divider"></div>
	</div>
	<button type="button" class="ui button small" @click.prevent="add" v-if="!isAdding">+</button>
	<div v-if="isAdding">
		<div class="ui fields inline">
			<div class="ui field">
				<select class="ui dropdown" style="width: 18em" v-model="routeCode">
					<option value="" v-if="keyword.length == 0">[请选择]</option>
					<option v-for="route in searchingRoutes" :value="route.code + '@' + route.domainId">{{route.name}}（{{route.domainName}}）</option>
				</select>
			</div>
			<div class="ui field">
				<input type="text" placeholder="搜索..." size="10" v-model="keyword" ref="keywordRef" @keyup.enter="confirm" @keypress.enter.prevent="1"/>
			</div>
			
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>
			</div>
			<div class="ui field">
				<a href="" @click.prevent="cancel()"><i class="icon remove small"></i></a>
			</div>
		</div>
	</div>
</div>`
})

Vue.component("dns-domain-selector", {
	props: ["v-domain-id", "v-domain-name"],
	data: function () {
		let domainId = this.vDomainId
		if (domainId == null) {
			domainId = 0
		}
		let domainName = this.vDomainName
		if (domainName == null) {
			domainName = ""
		}
		return {
			domainId: domainId,
			domainName: domainName
		}
	},
	methods: {
		select: function () {
			let that = this
			teaweb.popup("/dns/domains/selectPopup", {
				callback: function (resp) {
					that.domainId = resp.data.domainId
					that.domainName = resp.data.domainName
					that.change()
				}
			})
		},
		remove: function() {
			this.domainId = 0
			this.domainName = ""
			this.change()
		},
		update: function () {
			let that = this
			teaweb.popup("/dns/domains/selectPopup?domainId=" + this.domainId, {
				callback: function (resp) {
					that.domainId = resp.data.domainId
					that.domainName = resp.data.domainName
					that.change()
				}
			})
		},
		change: function () {
			this.$emit("change", {
				id: this.domainId,
				name: this.domainName
			})
		}
	},
	template: `<div>
	<input type="hidden" name="dnsDomainId" :value="domainId"/>
	<div v-if="domainName.length > 0">
		<span class="ui label small basic">
			{{domainName}}
			<a href="" @click.prevent="update"><i class="icon pencil small"></i></a>
			<a href="" @click.prevent="remove()"><i class="icon remove"></i></a>
		</span>
	</div>
	<div v-if="domainName.length == 0">
		<a href="" @click.prevent="select()">[选择域名]</a>
	</div>
</div>`
})

Vue.component("grant-selector", {
	props: ["v-grant", "v-node-cluster-id", "v-ns-cluster-id"],
	data: function () {
		return {
			grantId: (this.vGrant == null) ? 0 : this.vGrant.id,
			grant: this.vGrant,
			nodeClusterId: (this.vNodeClusterId != null) ? this.vNodeClusterId : 0,
			nsClusterId: (this.vNsClusterId != null) ? this.vNsClusterId : 0
		}
	},
	methods: {
		// 选择授权
		select: function () {
			let that = this;
			teaweb.popup("/clusters/grants/selectPopup?nodeClusterId=" + this.nodeClusterId + "&nsClusterId=" + this.nsClusterId, {
				callback: (resp) => {
					that.grantId = resp.data.grant.id;
					if (that.grantId > 0) {
						that.grant = resp.data.grant;
					}
					that.notifyUpdate()
				},
				height: "26em"
			})
		},

		// 创建授权
		create: function () {
			let that = this
			teaweb.popup("/clusters/grants/createPopup", {
				height: "26em",
				callback: (resp) => {
					that.grantId = resp.data.grant.id;
					if (that.grantId > 0) {
						that.grant = resp.data.grant;
					}
					that.notifyUpdate()
				}
			})
		},

		// 修改授权
		update: function () {
			if (this.grant == null) {
				window.location.reload();
				return;
			}
			let that = this
			teaweb.popup("/clusters/grants/updatePopup?grantId=" + this.grant.id, {
				height: "26em",
				callback: (resp) => {
					that.grant = resp.data.grant
					that.notifyUpdate()
				}
			})
		},

		// 删除已选择授权
		remove: function () {
			this.grant = null
			this.grantId = 0
			this.notifyUpdate()
		},
		notifyUpdate: function () {
			this.$emit("change", this.grant)
		}
	},
	template: `<div>
	<input type="hidden" name="grantId" :value="grantId"/>
	<div class="ui label small basic" v-if="grant != null">{{grant.name}}<span class="small grey">（{{grant.methodName}}）</span><span class="small grey" v-if="grant.username != null && grant.username.length > 0">（{{grant.username}}）</span> <a href="" title="修改" @click.prevent="update()"><i class="icon pencil small"></i></a> <a href="" title="删除" @click.prevent="remove()"><i class="icon remove"></i></a> </div>
	<div v-if="grant == null">
		<a href="" @click.prevent="select()">[选择已有认证]</a> &nbsp; &nbsp; <a href="" @click.prevent="create()">[添加新认证]</a>
	</div>
</div>`
})

window.REQUEST_COND_COMPONENTS = [{"type":"url-extension","name":"URL扩展名","description":"根据URL中的文件路径扩展名进行过滤","component":"http-cond-url-extension","paramsTitle":"扩展名列表","isRequest":true},{"type":"url-prefix","name":"URL前缀","description":"根据URL中的文件路径前缀进行过滤","component":"http-cond-url-prefix","paramsTitle":"URL前缀","isRequest":true},{"type":"url-eq","name":"URL精准匹配","description":"检查URL中的文件路径是否一致","component":"http-cond-url-eq","paramsTitle":"URL完整路径","isRequest":true},{"type":"url-regexp","name":"URL正则匹配","description":"使用正则表达式检查URL中的文件路径是否一致","component":"http-cond-url-regexp","paramsTitle":"正则表达式","isRequest":true},{"type":"params","name":"参数匹配","description":"根据参数值进行匹配","component":"http-cond-params","paramsTitle":"参数配置","isRequest":true},{"type":"url-not-prefix","name":"排除：URL前缀","description":"根据URL中的文件路径前缀进行过滤","component":"http-cond-url-not-prefix","paramsTitle":"URL前缀","isRequest":true},{"type":"url-not-eq","name":"排除：URL精准匹配","description":"检查URL中的文件路径是否一致","component":"http-cond-url-not-eq","paramsTitle":"URL完整路径","isRequest":true},{"type":"url-not-regexp","name":"排除：URL正则匹配","description":"使用正则表达式检查URL中的文件路径是否一致，如果一致，则不匹配","component":"http-cond-url-not-regexp","paramsTitle":"正则表达式","isRequest":true},{"type":"mime-type","name":"内容MimeType","description":"根据服务器返回的内容的MimeType进行过滤。注意：当用于缓存条件时，此条件需要结合别的请求条件使用。","component":"http-cond-mime-type","paramsTitle":"MimeType列表","isRequest":false}]

window.REQUEST_COND_OPERATORS = [{"description":"判断是否正则表达式匹配","name":"正则表达式匹配","op":"regexp"},{"description":"判断是否正则表达式不匹配","name":"正则表达式不匹配","op":"not regexp"},{"description":"使用字符串对比参数值是否相等于某个值","name":"字符串等于","op":"eq"},{"description":"参数值包含某个前缀","name":"字符串前缀","op":"prefix"},{"description":"参数值包含某个后缀","name":"字符串后缀","op":"suffix"},{"description":"参数值包含另外一个字符串","name":"字符串包含","op":"contains"},{"description":"参数值不包含另外一个字符串","name":"字符串不包含","op":"not contains"},{"description":"使用字符串对比参数值是否不相等于某个值","name":"字符串不等于","op":"not"},{"description":"判断参数值在某个列表中","name":"在列表中","op":"in"},{"description":"判断参数值不在某个列表中","name":"不在列表中","op":"not in"},{"description":"判断小写的扩展名（不带点）在某个列表中","name":"扩展名","op":"file ext"},{"description":"判断MimeType在某个列表中，支持类似于image/*的语法","name":"MimeType","op":"mime type"},{"description":"判断版本号在某个范围内，格式为version1,version2","name":"版本号范围","op":"version range"},{"description":"将参数转换为整数数字后进行对比","name":"整数等于","op":"eq int"},{"description":"将参数转换为可以有小数的浮点数字进行对比","name":"浮点数等于","op":"eq float"},{"description":"将参数转换为数字进行对比","name":"数字大于","op":"gt"},{"description":"将参数转换为数字进行对比","name":"数字大于等于","op":"gte"},{"description":"将参数转换为数字进行对比","name":"数字小于","op":"lt"},{"description":"将参数转换为数字进行对比","name":"数字小于等于","op":"lte"},{"description":"对整数参数值取模，除数为10，对比值为余数","name":"整数取模10","op":"mod 10"},{"description":"对整数参数值取模，除数为100，对比值为余数","name":"整数取模100","op":"mod 100"},{"description":"对整数参数值取模，对比值格式为：除数,余数，比如10,1","name":"整数取模","op":"mod"},{"description":"将参数转换为IP进行对比","name":"IP等于","op":"eq ip"},{"description":"将参数转换为IP进行对比","name":"IP大于","op":"gt ip"},{"description":"将参数转换为IP进行对比","name":"IP大于等于","op":"gte ip"},{"description":"将参数转换为IP进行对比","name":"IP小于","op":"lt ip"},{"description":"将参数转换为IP进行对比","name":"IP小于等于","op":"lte ip"},{"description":"IP在某个范围之内，范围格式可以是英文逗号分隔的ip1,ip2，或者CIDR格式的ip/bits","name":"IP范围","op":"ip range"},{"description":"对IP参数值取模，除数为10，对比值为余数","name":"IP取模10","op":"ip mod 10"},{"description":"对IP参数值取模，除数为100，对比值为余数","name":"IP取模100","op":"ip mod 100"},{"description":"对IP参数值取模，对比值格式为：除数,余数，比如10,1","name":"IP取模","op":"ip mod"},{"description":"判断参数值解析后的文件是否存在","name":"文件存在","op":"file exist"},{"description":"判断参数值解析后的文件是否不存在","name":"文件不存在","op":"file not exist"}]

window.REQUEST_VARIABLES = [{"code":"${edgeVersion}","description":"","name":"边缘节点版本"},{"code":"${remoteAddr}","description":"会依次根据X-Forwarded-For、X-Real-IP、RemoteAddr获取，适合前端有别的反向代理服务时使用，存在伪造的风险","name":"客户端地址（IP）"},{"code":"${rawRemoteAddr}","description":"返回直接连接服务的客户端原始IP地址","name":"客户端地址（IP）"},{"code":"${remotePort}","description":"","name":"客户端端口"},{"code":"${remoteUser}","description":"","name":"客户端用户名"},{"code":"${requestURI}","description":"比如/hello?name=lily","name":"请求URI"},{"code":"${requestPath}","description":"比如/hello","name":"请求路径（不包括参数）"},{"code":"${requestURL}","description":"比如https://example.com/hello?name=lily","name":"完整的请求URL"},{"code":"${requestLength}","description":"","name":"请求内容长度"},{"code":"${requestMethod}","description":"比如GET、POST","name":"请求方法"},{"code":"${requestFilename}","description":"","name":"请求文件路径"},{"code":"${scheme}","description":"","name":"请求协议，http或https"},{"code":"${proto}","description:":"类似于HTTP/1.0","name":"包含版本的HTTP请求协议"},{"code":"${timeISO8601}","description":"比如2018-07-16T23:52:24.839+08:00","name":"ISO 8601格式的时间"},{"code":"${timeLocal}","description":"比如17/Jul/2018:09:52:24 +0800","name":"本地时间"},{"code":"${msec}","description":"比如1531756823.054","name":"带有毫秒的时间"},{"code":"${timestamp}","description":"","name":"unix时间戳，单位为秒"},{"code":"${host}","description":"","name":"主机名"},{"code":"${serverName}","description":"","name":"接收请求的服务器名"},{"code":"${serverPort}","description":"","name":"接收请求的服务器端口"},{"code":"${referer}","description":"","name":"请求来源URL"},{"code":"${referer.host}","description":"","name":"请求来源URL域名"},{"code":"${userAgent}","description":"","name":"客户端信息"},{"code":"${contentType}","description":"","name":"请求头部的Content-Type"},{"code":"${cookies}","description":"","name":"所有cookie组合字符串"},{"code":"${cookie.NAME}","description":"","name":"单个cookie值"},{"code":"${args}","description":"","name":"所有参数组合字符串"},{"code":"${arg.NAME}","description":"","name":"单个参数值"},{"code":"${headers}","description":"","name":"所有Header信息组合字符串"},{"code":"${header.NAME}","description":"","name":"单个Header值"}]

window.METRIC_HTTP_KEYS = [{"name":"客户端地址（IP）","code":"${remoteAddr}","description":"会依次根据X-Forwarded-For、X-Real-IP、RemoteAddr获取，适用于前端可能有别的反向代理的情形，存在被伪造的可能","icon":""},{"name":"直接客户端地址（IP）","code":"${rawRemoteAddr}","description":"返回直接连接服务的客户端原始IP地址","icon":""},{"name":"客户端用户名","code":"${remoteUser}","description":"通过基本认证填入的用户名","icon":""},{"name":"请求URI","code":"${requestURI}","description":"包含参数，比如/hello?name=lily","icon":""},{"name":"请求路径","code":"${requestPath}","description":"不包含参数，比如/hello","icon":""},{"name":"完整URL","code":"${requestURL}","description":"比如https://example.com/hello?name=lily","icon":""},{"name":"请求方法","code":"${requestMethod}","description":"比如GET、POST等","icon":""},{"name":"请求协议Scheme","code":"${scheme}","description":"http或https","icon":""},{"name":"文件扩展名","code":"${requestPathExtension}","description":"请求路径中的文件扩展名，包括点符号，比如.html、.png","icon":""},{"name":"主机名","code":"${host}","description":"通常是请求的域名","icon":""},{"name":"请求协议Proto","code":"${proto}","description":"包含版本的HTTP请求协议，类似于HTTP/1.0","icon":""},{"name":"HTTP协议","code":"${proto}","description":"包含版本的HTTP请求协议，类似于HTTP/1.0","icon":""},{"name":"URL参数值","code":"${arg.NAME}","description":"单个URL参数值","icon":""},{"name":"请求来源URL","code":"${referer}","description":"请求来源Referer URL","icon":""},{"name":"请求来源URL域名","code":"${referer.host}","description":"请求来源Referer URL域名","icon":""},{"name":"Header值","code":"${header.NAME}","description":"单个Header值，比如${header.User-Agent}","icon":""},{"name":"Cookie值","code":"${cookie.NAME}","description":"单个cookie值，比如${cookie.sid}","icon":""},{"name":"状态码","code":"${status}","description":"","icon":""},{"name":"响应的Content-Type值","code":"${response.contentType}","description":"","icon":""}]

window.IP_ADDR_THRESHOLD_ITEMS = [{"code":"nodeAvgRequests","description":"当前节点在单位时间内接收到的平均请求数。","name":"节点平均请求数","unit":"个"},{"code":"nodeAvgTrafficOut","description":"当前节点在单位时间内发送的下行流量。","name":"节点平均下行流量","unit":"M"},{"code":"nodeAvgTrafficIn","description":"当前节点在单位时间内接收的上行流量。","name":"节点平均上行流量","unit":"M"},{"code":"nodeHealthCheck","description":"当前节点健康检查结果。","name":"节点健康检查结果","unit":""},{"code":"connectivity","description":"通过区域监控得到的当前IP地址的连通性数值，取值在0和100之间。","name":"IP连通性","unit":"%"},{"code":"groupAvgRequests","description":"当前节点所在分组在单位时间内接收到的平均请求数。","name":"分组平均请求数","unit":"个"},{"code":"groupAvgTrafficOut","description":"当前节点所在分组在单位时间内发送的下行流量。","name":"分组平均下行流量","unit":"M"},{"code":"groupAvgTrafficIn","description":"当前节点所在分组在单位时间内接收的上行流量。","name":"分组平均上行流量","unit":"M"},{"code":"clusterAvgRequests","description":"当前节点所在集群在单位时间内接收到的平均请求数。","name":"集群平均请求数","unit":"个"},{"code":"clusterAvgTrafficOut","description":"当前节点所在集群在单位时间内发送的下行流量。","name":"集群平均下行流量","unit":"M"},{"code":"clusterAvgTrafficIn","description":"当前节点所在集群在单位时间内接收的上行流量。","name":"集群平均上行流量","unit":"M"}]

window.IP_ADDR_THRESHOLD_ACTIONS = [{"code":"up","description":"上线当前IP。","name":"上线"},{"code":"down","description":"下线当前IP。","name":"下线"},{"code":"notify","description":"发送已达到阈值通知。","name":"通知"},{"code":"switch","description":"在DNS中记录中将IP切换到指定的备用IP。","name":"切换"},{"code":"webHook","description":"调用外部的WebHook。","name":"WebHook"}]

