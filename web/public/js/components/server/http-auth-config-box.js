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
				}
			})
		},
		update: function (index, policyId) {
			let that = this
			teaweb.popup("/servers/server/settings/access/updatePopup?policyId=" + policyId, {
				callback: function (resp) {
					Vue.set(that.authConfig.policyRefs, index, resp.data.policyRef)
				}
			})
		},
		delete: function (index) {
			that.authConfig.policyRefs.$remove(index)
		}
	},
	template: `<div>
<input type="text" name="authJSON" :value="JSON.stringify(authConfig)"/> 
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
<!-- 认证方法 -->
<div>
	<table class="ui table selectable celled" v-show="authConfig.policyRefs.length > 0">
		<thead>
			<tr>
				<th>认证方法</th>
				<th>参数</th>
				<th class="two wide">状态</th>
				<th class="two op">操作</th>
			</tr>
		</thead>
		<tbody v-for="ref in authConfig.policyRefs" :key="ref.authPolicyId">
			<tr>
				<td></td>
			</tr>
		</tbody>
	</table>
	<button class="ui button small" type="button" @click.prevent="add">+添加认证</button>
</div>
<div class="margin"></div>
</div>`
})