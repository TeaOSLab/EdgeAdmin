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
					that.change()
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
			this.change()
		},
		methodName: function (methodType) {
			switch (methodType) {
				case "basicAuth":
					return "BasicAuth"
				case "subRequest":
					return "子请求"
				case "typeA":
					return "URL鉴权A"
				case "typeB":
					return "URL鉴权B"
				case "typeC":
					return "URL鉴权C"
				case "typeD":
					return "URL鉴权D"
			}
			return ""
		},
		change: function () {
			let that = this
			setTimeout(function () {
				// 延时通知，是为了让表单有机会变更数据
				that.$emit("change", this.authConfig)
			}, 100)
		}
	},
	template: `<div>
<input type="hidden" name="authJSON" :value="JSON.stringify(authConfig)"/> 
<table class="ui table selectable definition">
	<prior-checkbox :v-config="authConfig" v-if="vIsLocation"></prior-checkbox>
	<tbody v-show="!vIsLocation || authConfig.isPrior">
		<tr>
			<td class="title">启用鉴权</td>
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
<!-- 鉴权方式 -->
<div v-show="isOn()">
	<h4>鉴权方式</h4>
	<table class="ui table selectable celled" v-show="authConfig.policyRefs.length > 0">
		<thead>
			<tr>
				<th class="three wide">名称</th>
				<th class="three wide">鉴权方法</th>
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
					<span v-if="ref.authPolicy.type == 'typeA'">{{ref.authPolicy.params.signParamName}}/有效期{{ref.authPolicy.params.life}}秒</span>
					<span v-if="ref.authPolicy.type == 'typeB'">有效期{{ref.authPolicy.params.life}}秒</span>
					<span v-if="ref.authPolicy.type == 'typeC'">有效期{{ref.authPolicy.params.life}}秒</span>
					<span v-if="ref.authPolicy.type == 'typeD'">{{ref.authPolicy.params.signParamName}}/{{ref.authPolicy.params.timestampParamName}}/有效期{{ref.authPolicy.params.life}}秒</span>
					
					<div v-if="(ref.authPolicy.params.exts != null && ref.authPolicy.params.exts.length > 0) || (ref.authPolicy.params.domains != null && ref.authPolicy.params.domains.length > 0)">
						<grey-label v-if="ref.authPolicy.params.exts != null" v-for="ext in ref.authPolicy.params.exts">扩展名：{{ext}}</grey-label>
						<grey-label v-if="ref.authPolicy.params.domains != null" v-for="domain in ref.authPolicy.params.domains">域名：{{domain}}</grey-label>
					</div>
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
	<button class="ui button small" type="button" @click.prevent="add">+添加鉴权方式</button>
</div>
<div class="margin"></div>
</div>`
})