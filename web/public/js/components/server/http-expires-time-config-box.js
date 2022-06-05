Vue.component("http-expires-time-config-box", {
	props: ["v-expires-time"],
	data: function () {
		let expiresTime = this.vExpiresTime
		if (expiresTime == null) {
			expiresTime = {
				isPrior: false,
				isOn: false,
				overwrite: true,
				autoCalculate: true,
				duration: {count: -1, "unit": "hour"}
			}
		}
		return {
			expiresTime: expiresTime
		}
	},
	watch: {
		"expiresTime.isPrior": function () {
			this.notifyChange()
		},
		"expiresTime.isOn": function () {
			this.notifyChange()
		},
		"expiresTime.overwrite": function () {
			this.notifyChange()
		},
		"expiresTime.autoCalculate": function () {
			this.notifyChange()
		}
	},
	methods: {
		notifyChange: function () {
			this.$emit("change", this.expiresTime)
		}
	},
	template: `<div>
	<table class="ui table">
		<prior-checkbox :v-config="expiresTime"></prior-checkbox>
		<tbody v-show="expiresTime.isPrior">
			<tr>
				<td class="title">启用</td>
				<td><checkbox v-model="expiresTime.isOn"></checkbox>
					<p class="comment">启用后，将会在响应的Header中添加<code-label>Expires</code-label>字段，浏览器据此会将内容缓存在客户端；同时，在管理后台执行清理缓存时，也将无法清理客户端已有的缓存。</p>
				</td>
			</tr>
			<tr v-show="expiresTime.isPrior && expiresTime.isOn">
				<td>覆盖源站设置</td>
				<td>
					<checkbox v-model="expiresTime.overwrite"></checkbox>
					<p class="comment">选中后，会覆盖源站Header中已有的<code-label>Expires</code-label>字段。</p>
				</td>
			</tr>
			<tr v-show="expiresTime.isPrior && expiresTime.isOn">
				<td>自动计算时间</td>
				<td><checkbox v-model="expiresTime.autoCalculate"></checkbox>
					<p class="comment">根据已设置的缓存有效期进行计算。</p>
				</td>
			</tr>
			<tr v-show="expiresTime.isPrior && expiresTime.isOn && !expiresTime.autoCalculate">
				<td>强制缓存时间</td>
				<td>
					<time-duration-box :v-value="expiresTime.duration" @change="notifyChange"></time-duration-box>
					<p class="comment">从客户端访问的时间开始要缓存的时长。</p>
				</td>
			</tr>
		</tbody>
	</table>
</div>`
})