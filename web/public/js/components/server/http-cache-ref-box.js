// 单个缓存条件设置
Vue.component("http-cache-ref-box", {
	props: ["v-cache-ref", "v-is-reverse"],
	data: function () {
		let ref = this.vCacheRef
		if (ref == null) {
			ref = {
				isOn: true,
				cachePolicyId: 0,
				key: "${scheme}://${host}${requestURI}",
				life: {count: 2, unit: "hour"},
				status: [200],
				maxSize: {count: 32, unit: "mb"},
				minSize: {count: 0, unit: "kb"},
				skipCacheControlValues: ["private", "no-cache", "no-store"],
				skipSetCookie: true,
				enableRequestCachePragma: false,
				conds: null,
				allowChunkedEncoding: true,
				isReverse: this.vIsReverse
			}
		}
		if (ref.life == null) {
			ref.life = {count: 2, unit: "hour"}
		}
		if (ref.maxSize == null) {
			ref.maxSize = {count: 32, unit: "mb"}
		}
		if (ref.minSize == null) {
			ref.minSize = {count: 0, unit: "kb"}
		}
		return {
			ref: ref,
			moreOptionsVisible: false
		}
	},
	methods: {
		changeOptionsVisible: function (v) {
			this.moreOptionsVisible = v
		},
		changeLife: function (v) {
			this.ref.life = v
		},
		changeMaxSize: function (v) {
			this.ref.maxSize = v
		},
		changeMinSize: function (v) {
			this.ref.minSize = v
		},
		changeConds: function (v) {
			this.ref.conds = v
		},
		changeStatusList: function (list) {
			let result = []
			list.forEach(function (status) {
				let statusNumber = parseInt(status)
				if (isNaN(statusNumber) || statusNumber < 100 || statusNumber > 999) {
					return
				}
				result.push(statusNumber)
			})
			this.ref.status = result
		}
	},
	template: `<tbody>
	<tr>
		<td class="title">匹配条件分组 *</td>
		<td>
			<http-request-conds-box :v-conds="ref.conds" @change="changeConds"></http-request-conds-box>
			
			<input type="hidden" name="cacheRefJSON" :value="JSON.stringify(ref)"/>
		</td>
	</tr>
	<tr v-show="!vIsReverse">
		<td>缓存有效期 *</td>
		<td>
			<time-duration-box :v-value="ref.life" @change="changeLife"></time-duration-box>
		</td>
	</tr>
	<tr v-show="!vIsReverse">
		<td>缓存Key *</td>
		<td>
			<input type="text" v-model="ref.key"/>
			<p class="comment">用来区分不同缓存内容的唯一Key。</p>
		</td>
	</tr>
	<tr v-show="!vIsReverse">
		<td colspan="2"><more-options-indicator @change="changeOptionsVisible"></more-options-indicator></td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>可缓存的最大内容尺寸</td>
		<td>
			<size-capacity-box :v-value="ref.maxSize" @change="changeMaxSize"></size-capacity-box>
			<p class="comment">内容尺寸如果高于此值则不缓存。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>可缓存的最小内容尺寸</td>
		<td>
			<size-capacity-box :v-value="ref.minSize" @change="changeMinSize"></size-capacity-box>
			<p class="comment">内容尺寸如果低于此值则不缓存。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>支持分片内容</td>
		<td>
			<checkbox name="allowChunkedEncoding" value="1" v-model="ref.allowChunkedEncoding"></checkbox>
			<p class="comment">选中后，Gzip和Chunked内容可以直接缓存，无需检查内容长度。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>状态码列表</td>
		<td>
			<values-box name="statusList" size="3" maxlength="3" :values="ref.status" @change="changeStatusList"></values-box>
			<p class="comment">允许缓存的HTTP状态码列表。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>跳过的Cache-Control值</td>
		<td>
			<values-box name="skipResponseCacheControlValues" size="10" maxlength="100" :values="ref.skipCacheControlValues"></values-box>
			<p class="comment">当响应的Cache-Control为这些值时不缓存响应内容，而且不区分大小写。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>跳过Set-Cookie</td>
		<td>
			<div class="ui checkbox">
				<input type="checkbox" value="1" v-model="ref.skipSetCookie"/>
				<label></label>
			</div>
			<p class="comment">选中后，当响应的Header中有Set-Cookie时不缓存响应内容。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>支持请求no-cache刷新</td>
		<td>
			<div class="ui checkbox">
				<input type="checkbox" name="enableRequestCachePragma" value="1" v-model="ref.enableRequestCachePragma"/>
				<label></label>
			</div>
			<p class="comment">选中后，当请求的Header中含有Pragma: no-cache或Cache-Control: no-cache时，会跳过缓存直接读取源内容。</p>
		</td>
	</tr>	
</tbody>`
})