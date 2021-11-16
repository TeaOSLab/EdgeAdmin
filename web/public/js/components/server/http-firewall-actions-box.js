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