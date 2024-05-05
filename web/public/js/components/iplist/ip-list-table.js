Vue.component("ip-list-table", {
	props: ["v-items", "v-keyword", "v-show-search-button", "v-total"/** total items >= items length **/],
	data: function () {
		let maxDeletes = 10000
		if (this.vTotal != null && this.vTotal > 0 && this.vTotal < maxDeletes) {
			maxDeletes = this.vTotal
		}

		return {
			items: this.vItems,
			keyword: (this.vKeyword != null) ? this.vKeyword : "",
			selectedAll: false,
			hasSelectedItems: false,

			MaxDeletes: maxDeletes
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
		},
		deleteCount: function () {
			let that = this
			teaweb.confirm("确定要批量删除当前列表中的" + this.MaxDeletes + "个IP吗？", function () {
				let query = window.location.search
				if (query.startsWith("?")) {
					query = query.substring(1)
				}
				Tea.action("/servers/iplists/deleteCount?" + query)
					.post()
					.params({count: that.MaxDeletes})
					.success(function () {
						teaweb.successToast("批量删除成功", 1200, teaweb.reload)
					})
			})
		},
		formatSeconds: function (seconds) {
			if (seconds < 60) {
				return seconds + "秒"
			}
			if (seconds < 3600) {
				return Math.ceil(seconds / 60) + "分钟"
			}
			if (seconds < 86400) {
				return Math.ceil(seconds / 3600) + "小时"
			}
			return Math.ceil(seconds / 86400) + "天"
		},
		cancelChecked: function () {
			this.hasSelectedItems = false
			this.selectedAll = false

			let boxes = this.$refs.itemCheckBox
			if (boxes == null) {
				return
			}
			boxes.forEach(function (box) {
				box.checked = false
			})
		}
	},
	template: `<div>
 <div v-show="hasSelectedItems">
 	<div class="ui divider"></div>
 	<button class="ui button basic" type="button" @click.prevent="deleteAll">批量删除所选</button>
 	&nbsp; &nbsp; 
 	<button class="ui button basic" type="button" @click.prevent="deleteCount" v-if="vTotal != null && vTotal >= MaxDeletes">批量删除所有搜索结果（{{MaxDeletes}}个）</button>
 	
 	&nbsp; &nbsp; 
 	<button class="ui button basic" type="button" @click.prevent="cancelChecked">取消选中</button>
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
                <th style="width: 6em">类型</th>
                <th style="width: 6em">级别</th>
                <th style="width: 12em">过期时间</th>
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
					<span v-if="item.type != 'all'" :class="{green: item.list != null && item.list.type == 'white', grey: item.list != null && item.list.type == 'grey'}">
						<span v-if="item.value != null && item.value.length > 0"><keyword :v-word="keyword">{{item.value}}</keyword></span>
						<span v-else>
							<keyword :v-word="keyword">{{item.ipFrom}}</keyword> <span> <span class="small red" v-if="item.isRead != null && !item.isRead">&nbsp;New&nbsp;</span>&nbsp;<a :href="'/servers/iplists?ip=' + item.ipFrom" v-if="vShowSearchButton" title="搜索此IP"><span><i class="icon search small" style="color: #ccc"></i></span></a></span>
							<span v-if="item.ipTo.length > 0"> - <keyword :v-word="keyword">{{item.ipTo}}</keyword></span>
						</span>
					</span>
					<span v-else class="disabled">*</span>
					
					<div v-if="item.region != null && item.region.length > 0">
						<span class="grey small">{{item.region}}</span>
						<span v-if="item.isp != null && item.isp.length > 0 && item.isp != '内网IP'" class="grey small"><span class="disabled">|</span> {{item.isp}}</span>
					</div>
					<div v-else-if="item.isp != null && item.isp.length > 0 && item.isp != '内网IP'"><span class="grey small">{{item.isp}}</span></div>
				
					<div v-if="item.createdTime != null">
						<span class="small grey">添加于 {{item.createdTime}}
							<span v-if="item.list != null && item.list.id > 0">
								@ 
								
								<a :href="'/servers/iplists/list?listId=' + item.list.id" v-if="item.policy.id == 0"><span>[<span v-if="item.list.type == 'black'">黑</span><span v-if="item.list.type == 'white'">白</span><span v-if="item.list.type == 'grey'">灰</span>名单：{{item.list.name}}]</span></a>
								<span v-else>[<span v-if="item.list.type == 'black'">黑</span><span v-if="item.list.type == 'white'">白</span><span v-if="item.list.type == 'grey'">灰</span>名单：{{item.list.name}}</span>
								
								<span v-if="item.policy.id > 0">
									<span v-if="item.policy.server != null">
										<a :href="'/servers/server/settings/waf/ipadmin/allowList?serverId=' + item.policy.server.id + '&firewallPolicyId=' + item.policy.id" v-if="item.list.type == 'white'">[网站：{{item.policy.server.name}}]</a>
										<a :href="'/servers/server/settings/waf/ipadmin/denyList?serverId=' + item.policy.server.id + '&firewallPolicyId=' + item.policy.id" v-if="item.list.type == 'black'">[网站：{{item.policy.server.name}}]</a>
										<a :href="'/servers/server/settings/waf/ipadmin/greyList?serverId=' + item.policy.server.id + '&firewallPolicyId=' + item.policy.id" v-if="item.list.type == 'grey'">[网站：{{item.policy.server.name}}]</a>
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
						<div v-if="item.isExpired && item.lifeSeconds == null" style="margin-top: 0.5em">
							<span class="ui label tiny basic red">已过期</span>
						</div>
						<div  v-if="item.lifeSeconds != null">
							<span class="small grey" v-if="item.lifeSeconds > 0">{{formatSeconds(item.lifeSeconds)}}</span>
							<span class="small red" v-if="item.lifeSeconds < 0">已过期</span>
						</div>
					</div>
					<span v-else class="disabled">不过期</span>
				</td>
				<td>
					<span v-if="item.reason.length > 0">{{item.reason}}</span>
					<span v-else class="disabled">-</span>
					
					<div v-if="item.sourceNode != null && item.sourceNode.id > 0" style="margin-top: 0.4em">
						<a :href="'/clusters/cluster/node?clusterId=' + item.sourceNode.clusterId + '&nodeId=' + item.sourceNode.id"><span class="small"><i class="icon cloud"></i>{{item.sourceNode.name}}</span></a>
					</div>
					<div style="margin-top: 0.4em" v-if="item.sourceServer != null && item.sourceServer.id > 0">
						<a :href="'/servers/server?serverId=' + item.sourceServer.id" style="border: 0"><span class="small "><i class="icon clone outline"></i>{{item.sourceServer.name}}</span></a>
					</div>
					<div v-if="item.sourcePolicy != null && item.sourcePolicy.id > 0" style="margin-top: 0.4em">
						<a :href="'/servers/components/waf/group?firewallPolicyId=' +  item.sourcePolicy.id + '&type=inbound&groupId=' + item.sourceGroup.id + '#set' + item.sourceSet.id" v-if="item.sourcePolicy.serverId == 0"><span class="small "><i class="icon shield"></i>{{item.sourcePolicy.name}} &raquo; {{item.sourceGroup.name}} &raquo; {{item.sourceSet.name}}</span></a>
						<a :href="'/servers/server/settings/waf/group?serverId=' + item.sourcePolicy.serverId + '&firewallPolicyId=' + item.sourcePolicy.id + '&type=inbound&groupId=' + item.sourceGroup.id + '#set' + item.sourceSet.id" v-if="item.sourcePolicy.serverId > 0"><span class="small "><i class="icon shield"></i> {{item.sourcePolicy.name}} &raquo; {{item.sourceGroup.name}} &raquo; {{item.sourceSet.name}}</span></a>
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