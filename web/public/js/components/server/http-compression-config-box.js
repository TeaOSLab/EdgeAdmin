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
				types: ["brotli", "gzip", "deflate"],
				level: 5,
				decompressData: false,
				gzipRef: null,
				deflateRef: null,
				brotliRef: null,
				minLength: {count: 0, "unit": "kb"},
				maxLength: {count: 0, "unit": "kb"},
				mimeTypes: ["text/*", "application/*", "font/*"],
				extensions: [".js", ".json", ".html", ".htm", ".xml", ".css", ".woff2", ".txt"],
				conds: null
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
	watch: {
		"config.level": function (v) {
			let level = parseInt(v)
			if (isNaN(level)) {
				level = 1
			} else if (level < 1) {
				level = 1
			} else if (level > 10) {
				level = 10
			}
			this.config.level = level
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
				<td class="title">是否启用</td>
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
				<td>压缩级别</td>
				<td>
					<select class="ui dropdown auto-width" v-model="config.level">
						<option v-for="i in 10" :value="i">{{i}}</option>	
					</select>
					<p class="comment">级别越高，压缩比例越大。</p>
				</td>
			</tr>
			<tr>
				<td>支持的扩展名</td>
				<td>
					<values-box :values="config.extensions" @change="changeExtensions" placeholder="比如 .html"></values-box>
					<p class="comment">含有这些扩展名的URL将会被压缩，不区分大小写。</p>
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
						<label v-if="config.useDefaultTypes" for="compression-use-default">使用默认顺序<span class="grey small">（brotli、gzip、deflate）</span></label>
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
					
					<p class="comment">选择支持的压缩算法和优先顺序，拖动<i class="icon list small grey"></i>图表排序。</p>
				</td>
			</tr>
			<tr>
				<td>支持已压缩内容</td>
				<td>
					<checkbox v-model="config.decompressData"></checkbox>
					<p class="comment">支持对已压缩内容尝试重新使用新的算法压缩。</p>
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
	<div class="margin"></div>
</div>`
})