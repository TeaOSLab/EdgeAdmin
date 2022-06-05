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
				<td class="title">启用配置</td>
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