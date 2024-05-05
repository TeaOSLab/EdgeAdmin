// 请求限制
Vue.component("http-request-limit-config-box", {
	props: ["v-request-limit-config", "v-is-group", "v-is-location"],
	data: function () {
		let config = this.vRequestLimitConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				maxConns: 0,
				maxConnsPerIP: 0,
				maxBodySize: {
					count: -1,
					unit: "kb"
				},
				outBandwidthPerConn: {
					count: -1,
					unit: "kb"
				}
			}
		}
		return {
			config: config,
			maxConns: config.maxConns,
			maxConnsPerIP: config.maxConnsPerIP
		}
	},
	watch: {
		maxConns: function (v) {
			let conns = parseInt(v, 10)
			if (isNaN(conns)) {
				this.config.maxConns = 0
				return
			}
			if (conns < 0) {
				this.config.maxConns = 0
			} else {
				this.config.maxConns = conns
			}
		},
		maxConnsPerIP: function (v) {
			let conns = parseInt(v, 10)
			if (isNaN(conns)) {
				this.config.maxConnsPerIP = 0
				return
			}
			if (conns < 0) {
				this.config.maxConnsPerIP = 0
			} else {
				this.config.maxConnsPerIP = conns
			}
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.config.isPrior) && this.config.isOn
		}
	},
	template: `<div>
	<input type="hidden" name="requestLimitJSON" :value="JSON.stringify(config)"/>
	<table class="ui table selectable definition">
		<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
			<tr>
				<td class="title">启用请求限制</td>
				<td>
					<checkbox v-model="config.isOn"></checkbox>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td>最大并发连接数</td>
				<td>
					<input type="text" maxlength="6" v-model="maxConns"/>
					<p class="comment">当前网站最大并发连接数，超出此限制则响应用户<code-label>429</code-label>代码。为0表示不限制。</p>
				</td>
			</tr>
			<tr>
				<td>单IP最大并发连接数</td>
				<td>
					<input type="text" maxlength="6" v-model="maxConnsPerIP"/>
					<p class="comment">单IP最大连接数，统计单个IP总连接数时不区分网站，超出此限制则响应用户<code-label>429</code-label>代码。为0表示不限制。<span v-if="maxConnsPerIP > 0 && maxConnsPerIP <= 3" class="red">当前设置的并发连接数过低，可能会影响正常用户访问，建议不小于3。</span></p>
				</td>
			</tr>
			<tr>
				<td>单连接带宽限制</td>
				<td>
					<size-capacity-box :v-value="config.outBandwidthPerConn" :v-supported-units="['byte', 'kb', 'mb']"></size-capacity-box>
					<p class="comment">客户端单个请求每秒可以读取的下行流量。</p>
				</td>
			</tr>
			<tr>
				<td>单请求最大尺寸</td>
				<td>
					<size-capacity-box :v-value="config.maxBodySize" :v-supported-units="['byte', 'kb', 'mb', 'gb']"></size-capacity-box>
					<p class="comment">单个请求能发送的最大内容尺寸。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})