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