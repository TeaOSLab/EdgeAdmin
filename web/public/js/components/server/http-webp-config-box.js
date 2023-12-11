Vue.component("http-webp-config-box", {
	props: ["v-webp-config", "v-is-location", "v-is-group", "v-require-cache"],
	data: function () {
		let config = this.vWebpConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				minLength: {count: 0, "unit": "kb"},
				maxLength: {count: 0, "unit": "kb"},
				mimeTypes: ["image/png", "image/jpeg", "image/bmp", "image/x-ico"],
				extensions: [".png", ".jpeg", ".jpg", ".bmp", ".ico"],
				conds: null
			}
		}

		if (config.mimeTypes == null) {
			config.mimeTypes = []
		}
		if (config.extensions == null) {
			config.extensions = []
		}

		return {
			config: config,
			moreOptionsVisible: false
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.config.isPrior) && this.config.isOn
		},
		changeExtensions: function (values) {
			values.forEach(function (v, k) {
				if (v.length > 0 && v[0] != ".") {
					values[k] = "." + v
				}
			})
			this.config.extensions = values
		},
		changeMimeTypes: function (values) {
			this.config.mimeTypes = values
		},
		changeAdvancedVisible: function () {
			this.moreOptionsVisible = !this.moreOptionsVisible
		},
		changeConds: function (conds) {
			this.config.conds = conds
		}
	},
	template: `<div>
	<input type="hidden" name="webpJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
			<tr>
				<td class="title">启用WebP压缩</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="config.isOn"/>
						<label></label>
					</div>
					<p class="comment">选中后表示开启自动WebP压缩；图片的宽和高均不能超过16383像素<span v-if="vRequireCache">；只有满足缓存条件的图片内容才会被转换</span>。</p>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-if="isOn()"></more-options-tbody>
		<tbody v-show="isOn() && moreOptionsVisible">
			<tr>
				<td>支持的扩展名</td>
				<td>
					<values-box :values="config.extensions" @change="changeExtensions" placeholder="比如 .html"></values-box>
					<p class="comment">含有这些扩展名的URL将会被转成WebP，不区分大小写。</p>
				</td>
			</tr>
			<tr>
				<td>支持的MimeType</td>
				<td>
					<values-box :values="config.mimeTypes" @change="changeMimeTypes" placeholder="比如 text/*"></values-box>
					<p class="comment">响应的Content-Type里包含这些MimeType的内容将会被转成WebP。</p>
				</td>
			</tr>
			<tr>
				<td>内容最小长度</td>
				<td>
					<size-capacity-box :v-name="'minLength'" :v-value="config.minLength" :v-unit="'kb'"></size-capacity-box>
					<p class="comment">0表示不限制，内容长度从文件尺寸或Content-Length中获取。</p>
				</td>
			</tr>
			<tr>
				<td>内容最大长度</td>
				<td>
					<size-capacity-box :v-name="'maxLength'" :v-value="config.maxLength" :v-unit="'mb'"></size-capacity-box>
					<p class="comment">0表示不限制，内容长度从文件尺寸或Content-Length中获取。</p>
				</td>
			</tr>
			<tr>
				<td>匹配条件</td>
				<td>
					<http-request-conds-box :v-conds="config.conds" @change="changeConds"></http-request-conds-box>
	</td>
			</tr>
		</tbody>
	</table>			
	<div class="ui margin"></div>
</div>`
})