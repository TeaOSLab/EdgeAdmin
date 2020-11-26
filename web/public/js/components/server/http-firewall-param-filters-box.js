Vue.component("http-firewall-param-filters-box", {
	props: ["v-filters"],
	data: function () {
		let filters = this.vFilters
		if (filters == null) {
			filters = []
		}

		return {
			filters: filters,
			isAdding: false,
			options: [
				{name: "MD5", code: "md5"},
				{name: "URLEncode", code: "urlEncode"},
				{name: "URLDecode", code: "urlDecode"},
				{name: "BASE64Encode", code: "base64Encode"},
				{name: "BASE64Decode", code: "base64Decode"},
				{name: "UNICODE编码", code: "unicodeEncode"},
				{name: "UNICODE解码", code: "unicodeDecode"},
				{name: "HTML实体编码", code: "htmlEscape"},
				{name: "HTML实体解码", code: "htmlUnescape"},
				{name: "计算长度", code: "length"},
				{name: "十六进制->十进制", "code": "hex2dec"},
				{name: "十进制->十六进制", "code": "dec2hex"},
				{name: "SHA1", "code": "sha1"},
				{name: "SHA256", "code": "sha256"}
			],
			addingCode: ""
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			this.addingCode = ""
		},
		confirm: function () {
			if (this.addingCode.length == 0) {
				return
			}
			let that = this
			this.filters.push(this.options.$find(function (k, v) {
				return (v.code == that.addingCode)
			}))
			this.isAdding = false
		},
		cancel: function () {
			this.isAdding = false
		},
		remove: function (index) {
			this.filters.$remove(index)
		}
	},
	template: `<div>
		<input type="hidden" name="paramFiltersJSON" :value="JSON.stringify(filters)" />
		<div v-if="filters.length > 0">
			<div v-for="(filter, index) in filters" class="ui label small basic">
				{{filter.name}} <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove"></i></a>
			</div>
			<div class="ui divider"></div>
		</div>
		<div v-if="isAdding">
			<div class="ui fields inline">
				<div class="ui field">
					<select class="ui dropdown auto-width" v-model="addingCode">
						<option value="">[请选择]</option>
						<option v-for="option in options" :value="option.code">{{option.name}}</option>
					</select>
				</div>
				<div class="ui field">
					<button class="ui button tiny" type="button" @click.prevent="confirm()">确定</button>
					&nbsp; <a href="" @click.prevent="cancel()" title="取消"><i class="icon remove"></i></a>
				</div>
			</div>
		</div>
		<div v-if="!isAdding">
			<button class="ui button tiny" type="button" @click.prevent="add">+</button>
		</div>
		<p class="comment">可以对参数值进行特定的编解码处理。</p>
</div>`
})