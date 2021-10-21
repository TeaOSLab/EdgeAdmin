Vue.component("bandwidth-limit-config-box", {
	props: ["v-bandwidth-limit"],
	data: function () {
		let config = this.vBandwidthLimit
		if (config == null) {
			config = {
				isOn: false,
				dailySize: {
					count: -1,
					unit: "gb"
				},
				monthlySize: {
					count: -1,
					unit: "gb"
				},
				totalSize: {
					count: -1,
					unit: "gb"
				},
				noticePageBody: ""
			}
		}
		return {
			config: config
		}
	},
	methods: {
		showBodyTemplate: function () {
			this.config.noticePageBody = `<!DOCTYPE html>
<html>
<head>
<title>Bandwidth Limit Exceeded Warning</title>
<body>

The site bandwidth has exceeded the limit. Please contact with the site administrator.

</body>
</html>`
		}
	},
	template: `<div>
	<input type="hidden" name="bandwidthLimitJSON" :value="JSON.stringify(config)"/>
	<table class="ui table selectable definition">
		<tbody>
			<tr>
				<td class="title">是否启用</td>
				<td>
					<checkbox v-model="config.isOn"></checkbox>
					<p class="comment">注意：由于带宽统计是每5分钟统计一次，所以超出带宽限制后，对用户的提醒也会有所延迟。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="config.isOn">
			<tr>
				<td>日带宽限制</td>
				<td>
					<size-capacity-box :v-value="config.dailySize"></size-capacity-box>
				</td>
			</tr>
			<tr>
				<td>月带宽限制</td>
				<td>
					<size-capacity-box :v-value="config.monthlySize"></size-capacity-box>
				</td>
			</tr>
			<!--<tr>
				<td>总体限制</td>
				<td>
					<size-capacity-box :v-value="config.totalSize"></size-capacity-box>
					<p class="comment"></p>
				</td>
			</tr>-->
			<tr>
				<td>网页提示内容</td>
				<td>
					<textarea v-model="config.noticePageBody"></textarea>
					<p class="comment"><a href="" @click.prevent="showBodyTemplate">[使用模板]</a>。当达到带宽限制时网页显示的HTML内容，不填写则显示默认的提示内容。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})