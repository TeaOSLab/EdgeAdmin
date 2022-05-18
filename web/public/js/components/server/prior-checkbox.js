Vue.component("prior-checkbox", {
	props: ["v-config", "description"],
	data: function () {
		let description = this.description
		if (description == null) {
			description = "打开后可以覆盖父级或子级配置"
		}
		return {
			isPrior: this.vConfig.isPrior,
			realDescription: description
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
			<p class="comment"><strong v-if="isPrior">[已打开]</strong> {{realDescription}}。</p>
		</td>
	</tr>
</tbody>`
})