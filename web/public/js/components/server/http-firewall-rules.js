// 通用Header长度
let defaultGeneralHeaders = ["Cache-Control", "Connection", "Date", "Pragma", "Trailer", "Transfer-Encoding", "Upgrade", "Via", "Warning"]
Vue.component("http-cond-general-header-length", {
	props: ["v-checkpoint"],
	data: function () {
		let headers = null
		let length = null

		if (window.parent.UPDATING_RULE != null) {
			let options = window.parent.UPDATING_RULE.checkpointOptions
			if (options.headers != null && Array.$isArray(options.headers)) {
				headers = options.headers
			}
			if (options.length != null) {
				length = options.length
			}
		}


		if (headers == null) {
			headers = defaultGeneralHeaders
		}

		if (length == null) {
			length = 128
		}

		let that = this
		setTimeout(function () {
			that.change()
		}, 100)

		return {
			headers: headers,
			length: length
		}
	},
	watch: {
		length: function (v) {
			let len = parseInt(v)
			if (isNaN(len)) {
				len = 0
			}
			if (len < 0) {
				len = 0
			}
			this.length = len
			this.change()
		}
	},
	methods: {
		change: function () {
			this.vCheckpoint.options = [
				{
					code: "headers",
					value: this.headers
				},
				{
					code: "length",
					value: this.length
				}
			]
		}
	},
	template: `<div>
	<table class="ui table">
		<tr>
			<td class="title">通用Header列表</td>
			<td>
				<values-box :values="headers" :placeholder="'Header'" @change="change"></values-box>
				<p class="comment">需要检查的Header列表。</p>
			</td>
		</tr>
		<tr>
			<td>Header值超出长度</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" name="" style="width: 5em" v-model="length" maxlength="6"/>
					<span class="ui label">字节</span>
				</div>
				<p class="comment">超出此长度认为匹配成功，0表示不限制。</p>
			</td>
		</tr>
	</table>
</div>`
})

// CC
Vue.component("http-firewall-checkpoint-cc", {
	props: ["v-checkpoint"],
	data: function () {
		let keys = []
		let period = 60
		let threshold = 1000

		let options = {}
		if (window.parent.UPDATING_RULE != null) {
			options = window.parent.UPDATING_RULE.checkpointOptions
		}

		if (options == null) {
			options = {}
		}
		if (options.keys != null) {
			keys = options.keys
		}
		if (keys.length == 0) {
			keys = ["${remoteAddr}", "${requestPath}"]
		}
		if (options.period != null) {
			period = options.period
		}
		if (options.threshold != null) {
			threshold = options.threshold
		}

		let that = this
		setTimeout(function () {
			that.change()
		}, 100)

		return {
			keys: keys,
			period: period,
			threshold: threshold,
			options: {},
			value: threshold
		}
	},
	watch: {
		period: function () {
			this.change()
		},
		threshold: function () {
			this.change()
		}
	},
	methods: {
		changeKeys: function (keys) {
			this.keys = keys
			this.change()
		},
		change: function () {
			let period = parseInt(this.period.toString())
			if (isNaN(period) || period <= 0) {
				period = 60
			}

			let threshold = parseInt(this.threshold.toString())
			if (isNaN(threshold) || threshold <= 0) {
				threshold = 1000
			}
			this.value = threshold

			this.vCheckpoint.options = [
				{
					code: "keys",
					value: this.keys
				},
				{
					code: "period",
					value: period,
				},
				{
					code: "threshold",
					value: threshold
				}
			]
		}
	},
	template: `<div>
	<input type="hidden" name="operator" value="gt"/>
	<input type="hidden" name="value" :value="value"/>
	<table class="ui table">
		<tr>
			<td class="title">统计对象组合 *</td>
			<td>
				<metric-keys-config-box :v-keys="keys" @change="changeKeys"></metric-keys-config-box>
			</td>
		</tr>
		<tr>
			<td>统计周期 *</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" v-model="period" style="width: 6em" maxlength="8"/>
					<span class="ui label">秒</span>
				</div>
			</td>
		</tr>
		<tr>
			<td>阈值 *</td>
			<td>
				<input type="text" v-model="threshold" style="width: 6em" maxlength="8"/>
			</td>
		</tr>
	</table>
</div>`
})