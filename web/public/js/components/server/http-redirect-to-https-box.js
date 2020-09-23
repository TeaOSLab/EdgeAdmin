Vue.component("http-redirect-to-https-box", {
	props: ["v-redirect-to-https-config", "v-is-location"],
	data: function () {
		let redirectToHttpsConfig = this.vRedirectToHttpsConfig
		if (redirectToHttpsConfig == null) {
			redirectToHttpsConfig = {
				isPrior: false,
				isOn: false
			}
		}
		return {
			redirectToHttpsConfig: redirectToHttpsConfig
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
					<p class="comment">开启后，所有HTTP的请求都会自动跳转到对应的HTTPS URL上。</p>
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
		<p class="comment">开启后，所有HTTP的请求都会自动跳转到对应的HTTPS URL上。</p>
	</div>
	<div class="margin"></div>
</div>`
})