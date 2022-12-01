Vue.component("traffic-limit-config-box", {
	props: ["v-traffic-limit"],
	data: function () {
		let config = this.vTrafficLimit
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
		if (config.dailySize == null) {
			config.dailySize = {
				count: -1,
				unit: "gb"
			}
		}
		if (config.monthlySize == null) {
			config.monthlySize = {
				count: -1,
				unit: "gb"
			}
		}
		if (config.totalSize == null) {
			config.totalSize = {
				count: -1,
				unit: "gb"
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
<title>Traffic Limit Exceeded Warning</title>
<body>

<h1>Traffic Limit Exceeded Warning</h1>
<p>The site traffic has exceeded the limit. Please contact with the site administrator.</p>
<address>Request ID: \${requestId}.</address>

</body>
</html>`
		}
	},
	template: `<div>
	<input type="hidden" name="trafficLimitJSON" :value="JSON.stringify(config)"/>
	<table class="ui table selectable definition">
		<tbody>
			<tr>
				<td class="title">启用流量限制</td>
				<td>
					<checkbox v-model="config.isOn"></checkbox>
					<p class="comment">注意：由于流量统计是每5分钟统计一次，所以超出流量限制后，对用户的提醒也会有所延迟。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="config.isOn">
			<tr>
				<td>日流量限制</td>
				<td>
					<size-capacity-box :v-value="config.dailySize"></size-capacity-box>
				</td>
			</tr>
			<tr>
				<td>月流量限制</td>
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
					<p class="comment"><a href="" @click.prevent="showBodyTemplate">[使用模板]</a>。当达到流量限制时网页显示的HTML内容，不填写则显示默认的提示内容，适用于网站类服务。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})