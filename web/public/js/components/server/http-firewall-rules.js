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
		let ignoreCommonFiles = true
		let enableFingerprint = true

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
		if (options.ignoreCommonFiles != null && typeof (options.ignoreCommonFiles) == "boolean") {
			ignoreCommonFiles = options.ignoreCommonFiles
		}
		if (options.enableFingerprint != null && typeof (options.enableFingerprint) == "boolean") {
			enableFingerprint = options.enableFingerprint
		}

		let that = this
		setTimeout(function () {
			that.change()
		}, 100)

		return {
			keys: keys,
			period: period,
			threshold: threshold,
			ignoreCommonFiles: ignoreCommonFiles,
			enableFingerprint: enableFingerprint,
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
		},
		ignoreCommonFiles: function () {
			this.change()
		},
		enableFingerprint: function () {
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

			let ignoreCommonFiles = this.ignoreCommonFiles
			if (typeof ignoreCommonFiles != "boolean") {
				ignoreCommonFiles = false
			}

			let enableFingerprint = this.enableFingerprint
			if (typeof enableFingerprint != "boolean") {
				enableFingerprint = true
			}

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
				},
				{
					code: "ignoreCommonFiles",
					value: ignoreCommonFiles
				},
				{
					code: "enableFingerprint",
					value: enableFingerprint
				}
			]
		},
		thresholdTooLow: function () {
			let threshold = parseInt(this.threshold.toString())
			if (isNaN(threshold) || threshold <= 0) {
				threshold = 1000
			}
			return threshold > 0 && threshold < 5
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
				<p class="comment" v-if="thresholdTooLow()"><span class="red">对于网站类应用来说，当前阈值设置的太低，有可能会影响用户正常访问。</span></p>
			</td>
		</tr>
		<tr>
			<td>检查请求来源指纹</td>
			<td>
				<checkbox v-model="enableFingerprint"></checkbox>
				<p class="comment">在接收到HTTPS请求时尝试检查请求来源的指纹，用来检测代理服务和爬虫攻击；如果你在网站前面放置了别的反向代理服务，请取消此选项。</p>
			</td>
		</tr>
		<tr>
			<td>忽略常用文件</td>
			<td>
				<checkbox v-model="ignoreCommonFiles"></checkbox>
				<p class="comment">忽略js、css、jpg等常在网页里被引用的文件名，即对这些文件的访问不加入计数，可以减少误判几率。</p>
			</td>
		</tr>
	</table>
</div>`
})

// 防盗链
Vue.component("http-firewall-checkpoint-referer-block", {
	props: ["v-checkpoint"],
	data: function () {
		let allowEmpty = true
		let allowSameDomain = true
		let allowDomains = []
		let denyDomains = []
		let checkOrigin = true

		let options = {}
		if (window.parent.UPDATING_RULE != null) {
			options = window.parent.UPDATING_RULE.checkpointOptions
		}

		if (options == null) {
			options = {}
		}
		if (typeof (options.allowEmpty) == "boolean") {
			allowEmpty = options.allowEmpty
		}
		if (typeof (options.allowSameDomain) == "boolean") {
			allowSameDomain = options.allowSameDomain
		}
		if (options.allowDomains != null && typeof (options.allowDomains) == "object") {
			allowDomains = options.allowDomains
		}
		if (options.denyDomains != null && typeof (options.denyDomains) == "object") {
			denyDomains = options.denyDomains
		}
		if (typeof options.checkOrigin == "boolean") {
			checkOrigin = options.checkOrigin
		}

		let that = this
		setTimeout(function () {
			that.change()
		}, 100)

		return {
			allowEmpty: allowEmpty,
			allowSameDomain: allowSameDomain,
			allowDomains: allowDomains,
			denyDomains: denyDomains,
			checkOrigin: checkOrigin,
			options: {},
			value: 0
		}
	},
	watch: {
		allowEmpty: function () {
			this.change()
		},
		allowSameDomain: function () {
			this.change()
		},
		checkOrigin: function () {
			this.change()
		}
	},
	methods: {
		changeAllowDomains: function (values) {
			this.allowDomains = values
			this.change()
		},
		changeDenyDomains: function (values) {
			this.denyDomains = values
			this.change()
		},
		change: function () {
			this.vCheckpoint.options = [
				{
					code: "allowEmpty",
					value: this.allowEmpty
				},
				{
					code: "allowSameDomain",
					value: this.allowSameDomain,
				},
				{
					code: "allowDomains",
					value: this.allowDomains
				},
				{
					code: "denyDomains",
					value: this.denyDomains
				},
				{
					code: "checkOrigin",
					value: this.checkOrigin
				}
			]
		}
	},
	template: `<div>
	<input type="hidden" name="operator" value="eq"/>
	<input type="hidden" name="value" :value="value"/>
	<table class="ui table">
		<tr>
			<td class="title">来源域名允许为空</td>
			<td>
				<checkbox v-model="allowEmpty"></checkbox>
				<p class="comment">允许不带来源的访问。</p>
			</td>
		</tr>
		<tr>
			<td>来源域名允许一致</td>
			<td>
				<checkbox v-model="allowSameDomain"></checkbox>
				<p class="comment">允许来源域名和当前访问的域名一致，相当于在站内访问。</p>
			</td>
		</tr>
		<tr>
			<td>允许的来源域名</td>
			<td>
				<values-box :values="allowDomains" @change="changeAllowDomains"></values-box>
				<p class="comment">允许的来源域名列表，比如<code-label>example.com</code-label>（顶级域名)、<code-label>*.example.com</code-label>（example.com的所有二级域名）。单个星号<code-label>*</code-label>表示允许所有域名。</p>
			</td>
		</tr>
		<tr>
			<td>禁止的来源域名</td>
			<td>
				<values-box :values="denyDomains" @change="changeDenyDomains"></values-box>
				<p class="comment">禁止的来源域名列表，比如<code-label>example.org</code-label>（顶级域名）、<code-label>*.example.org</code-label>（example.org的所有二级域名）；除了这些禁止的来源域名外，其他域名都会被允许，除非限定了允许的来源域名。</p>
			</td>
		</tr>
		<tr>
			<td>同时检查Origin</td>
			<td>
				<checkbox v-model="checkOrigin"></checkbox>
				<p class="comment">如果请求没有指定Referer Header，则尝试检查Origin Header，多用于跨站调用。</p>
			</td>
		</tr>
	</table>
</div>`
})