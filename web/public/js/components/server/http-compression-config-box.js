// 压缩配置
Vue.component("http-compression-config-box", {
	props: ["v-compression-config", "v-is-location", "v-is-group"],
	mounted: function () {
		let that = this
		sortLoad(function () {
			that.initSortableTypes()
		})
	},
	data: function () {
		let config = this.vCompressionConfig
		if (config == null) {
			config = {
				isPrior: false,
				isOn: false,
				useDefaultTypes: true,
				types: ["brotli", "gzip", "zstd", "deflate"],
				level: 0,
				decompressData: false,
				gzipRef: null,
				deflateRef: null,
				brotliRef: null,
				minLength: {count: 1, "unit": "kb"},
				maxLength: {count: 32, "unit": "mb"},
				mimeTypes: ["text/*", "application/javascript", "application/json", "application/atom+xml", "application/rss+xml", "application/xhtml+xml", "font/*", "image/svg+xml"],
				extensions: [".js", ".json", ".html", ".htm", ".xml", ".css", ".woff2", ".txt"],
				exceptExtensions: [".apk", ".ipa"],
				conds: null,
				enablePartialContent: false,
				onlyURLPatterns: [],
				exceptURLPatterns: []
			}
		}

		if (config.types == null) {
			config.types = []
		}
		if (config.mimeTypes == null) {
			config.mimeTypes = []
		}
		if (config.extensions == null) {
			config.extensions = []
		}

		let allTypes = [
			{
				name: "Gzip",
				code: "gzip",
				isOn: true
			},
			{
				name: "Deflate",
				code: "deflate",
				isOn: true
			},
			{
				name: "Brotli",
				code: "brotli",
				isOn: true
			},
			{
				name: "ZSTD",
				code: "zstd",
				isOn: true
			}
		]

		let configTypes = []
		config.types.forEach(function (typeCode) {
			allTypes.forEach(function (t) {
				if (typeCode == t.code) {
					t.isOn = true
					configTypes.push(t)
				}
			})
		})
		allTypes.forEach(function (t) {
			if (!config.types.$contains(t.code)) {
				t.isOn = false
				configTypes.push(t)
			}
		})

		return {
			config: config,
			moreOptionsVisible: false,
			allTypes: configTypes
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
		changeExceptExtensions: function (values) {
			values.forEach(function (v, k) {
				if (v.length > 0 && v[0] != ".") {
					values[k] = "." + v
				}
			})
			this.config.exceptExtensions = values
		},
		changeMimeTypes: function (values) {
			this.config.mimeTypes = values
		},
		changeAdvancedVisible: function () {
			this.moreOptionsVisible = !this.moreOptionsVisible
		},
		changeConds: function (conds) {
			this.config.conds = conds
		},
		changeType: function () {
			this.config.types = []
			let that = this
			this.allTypes.forEach(function (v) {
				if (v.isOn) {
					that.config.types.push(v.code)
				}
			})
		},
		initSortableTypes: function () {
			let box = document.querySelector("#compression-types-box")
			let that = this
			Sortable.create(box, {
				draggable: ".checkbox",
				handle: ".icon.handle",
				onStart: function () {

				},
				onUpdate: function (event) {
					let checkboxes = box.querySelectorAll(".checkbox")
					let codes = []
					checkboxes.forEach(function (checkbox) {
						let code = checkbox.getAttribute("data-code")
						codes.push(code)
					})
					that.config.types = codes
				}
			})
		}
	},
	template: `<div>
	<input type="hidden" name="compressionJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
			<tr>
				<td class="title">启用内容压缩</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="config.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td>支持的扩展名</td>
				<td>
					<values-box :values="config.extensions" @change="changeExtensions" placeholder="比如 .html"></values-box>
					<p class="comment">含有这些扩展名的URL将会被压缩，不区分大小写。</p>
				</td>
			</tr>
			<tr>
				<td>例外扩展名</td>
				<td>
					<values-box :values="config.exceptExtensions" @change="changeExceptExtensions" placeholder="比如 .html"></values-box>
					<p class="comment">含有这些扩展名的URL将<strong>不会</strong>被压缩，不区分大小写。</p>
				</td>
			</tr>
			<tr>
				<td>支持的MimeType</td>
				<td>
					<values-box :values="config.mimeTypes" @change="changeMimeTypes" placeholder="比如 text/*"></values-box>
					<p class="comment">响应的Content-Type里包含这些MimeType的内容将会被压缩。</p>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-if="isOn()"></more-options-tbody>
		<tbody v-show="isOn() && moreOptionsVisible">
			<tr>
				<td>压缩算法</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="config.useDefaultTypes" id="compression-use-default"/>
						<label v-if="config.useDefaultTypes" for="compression-use-default">使用默认顺序<span class="grey small">（brotli、gzip、 zstd、deflate）</span></label>
						<label v-if="!config.useDefaultTypes" for="compression-use-default">使用默认顺序</label>
					</div>
					<div v-show="!config.useDefaultTypes">
						<div class="ui divider"></div>
						<div id="compression-types-box">
							<div class="ui checkbox" v-for="t in allTypes" style="margin-right: 2em" :data-code="t.code">
								<input type="checkbox" v-model="t.isOn" :id="'compression-type-' + t.code" @change="changeType"/>
								<label :for="'compression-type-' + t.code">{{t.name}} &nbsp; <i class="icon list small grey handle"></i></label>
							</div>
						</div>
					</div>
					
					<p class="comment" v-show="!config.useDefaultTypes">选择支持的压缩算法和优先顺序，拖动<i class="icon list small grey"></i>图表排序。</p>
				</td>
			</tr>
			<tr>
				<td>支持已压缩内容</td>
				<td>
					<checkbox v-model="config.decompressData"></checkbox>
					<p class="comment">支持对已压缩内容尝试重新使用新的算法压缩；不选中表示保留当前的压缩格式。</p>
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
				<td>支持Partial<br/>Content</td>
				<td>
					<checkbox v-model="config.enablePartialContent"></checkbox>
					<p class="comment">支持对分片内容（PartialContent）的压缩；除非客户端有特殊要求，一般不需要启用。</p>
				</td>
			</tr>
				<tr>
				<td>例外URL</td>
				<td>
					<url-patterns-box v-model="config.exceptURLPatterns"></url-patterns-box>
					<p class="comment">如果填写了例外URL，表示这些URL跳过不做处理。</p>
				</td>
			</tr>
			<tr>
				<td>限制URL</td>
				<td>
					<url-patterns-box v-model="config.onlyURLPatterns"></url-patterns-box>
					<p class="comment">如果填写了限制URL，表示只对这些URL进行压缩处理；如果不填则表示支持所有的URL。</p>
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
	<div class="margin"></div>
</div>`
})