Vue.component("reverse-proxy-box", {
	props: ["v-reverse-proxy-ref", "v-is-location"],
	data: function () {
		let reverseProxyRef = this.vReverseProxyRef
		if (reverseProxyRef == null) {
			reverseProxyRef = {
				isPrior: false,
				isOn: false,
				reverseProxyId: 0
			}
		}
		return {
			reverseProxyRef: reverseProxyRef
		}
	},
	template: `<div>
	<input type="hidden" name="reverseProxyRefJSON" :value="JSON.stringify(reverseProxyRef)"/>
	<table class="ui table selectable definition">
		<prior-checkbox :v-config="reverseProxyRef" v-if="vIsLocation"></prior-checkbox>
		<tbody v-show="!vIsLocation || reverseProxyRef.isPrior">
			<tr>
				<td class="title">是否启用反向代理</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="reverseProxyRef.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})