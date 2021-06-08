// URL扩展名条件
Vue.component("http-cond-url-extension", {
	props: ["v-cond"],
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPathExtension}",
			operator: "in",
			value: "[]"
		}
		if (this.vCond != null && this.vCond.param == cond.param) {
			cond.value = this.vCond.value
		}

		let extensions = []
		try {
			extensions = JSON.parse(cond.value)
		} catch (e) {

		}

		return {
			cond: cond,
			extensions: extensions, // TODO 可以拖动排序

			isAdding: false,
			addingExt: ""
		}
	},
	watch: {
		extensions: function () {
			this.cond.value = JSON.stringify(this.extensions)
		}
	},
	methods: {
		addExt: function () {
			this.isAdding = !this.isAdding

			if (this.isAdding) {
				let that = this
				setTimeout(function () {
					that.$refs.addingExt.focus()
				}, 100)
			}
		},
		cancelAdding: function () {
			this.isAdding = false
			this.addingExt = ""
		},
		confirmAdding: function () {
			// TODO 做更详细的校验
			// TODO 如果有重复的则提示之

			if (this.addingExt.length == 0) {
				return
			}
			if (this.addingExt[0] != ".") {
				this.addingExt = "." + this.addingExt
			}
			this.addingExt = this.addingExt.replace(/\s+/g, "").toLowerCase()
			this.extensions.push(this.addingExt)

			// 清除状态
			this.cancelAdding()
		},
		removeExt: function (index) {
			this.extensions.$remove(index)
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<div v-if="extensions.length > 0">
		<div class="ui label small" v-for="(ext, index) in extensions">{{ext}} <a href="" title="删除" @click.prevent="removeExt(index)"><i class="icon remove"></i></a></div>
		<div class="ui divider"></div>
	</div>
	<div class="ui fields inline" v-if="isAdding">
		<div class="ui field">
			<input type="text" size="6" maxlength="100" v-model="addingExt" ref="addingExt" placeholder=".xxx" @keyup.enter="confirmAdding" @keypress.enter.prevent="1" />
		</div>
		<div class="ui field">
			<button class="ui button tiny" type="button" @click.prevent="confirmAdding">确认</button>
			<a href="" title="取消" @click.prevent="cancelAdding"><i class="icon remove"></i></a>
		</div> 
	</div>
	<div style="margin-top: 1em">
		<button class="ui button tiny" type="button" @click.prevent="addExt()">+添加扩展名</button>
	</div>
	<p class="comment">扩展名需要包含点（.）符号，例如<span class="ui label tiny">.jpg</span>、<span class="ui label tiny">.png</span>之类。</p>
</div>`
})

// 根据URL前缀
Vue.component("http-cond-url-prefix", {
	props: ["v-cond"],
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPath}",
			operator: "prefix",
			value: ""
		}
		if (this.vCond != null && typeof (this.vCond.value) == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value"/>
	<p class="comment">URL前缀，有此前缀的URL都将会被匹配，通常以<code-label>/</code-label>开头，比如<code-label>/static</code-label>。</p>
</div>`
})

Vue.component("http-cond-url-not-prefix", {
	props: ["v-cond"],
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPath}",
			operator: "prefix",
			value: "",
			isReverse: true
		}
		if (this.vCond != null && typeof this.vCond.value == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value"/>
	<p class="comment">要排除的URL前缀，有此前缀的URL都将会被匹配，通常以<code-label>/</code-label>开头，比如<code-label>/static</code-label>。</p>
</div>`
})

// URL精准匹配
Vue.component("http-cond-url-eq", {
	props: ["v-cond"],
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPath}",
			operator: "eq",
			value: ""
		}
		if (this.vCond != null && typeof this.vCond.value == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value"/>
	<p class="comment">完整的URL路径，通常以<code-label>/</code-label>开头，比如<code-label>/static/ui.js</code-label>。</p>
</div>`
})

Vue.component("http-cond-url-not-eq", {
	props: ["v-cond"],
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPath}",
			operator: "eq",
			value: "",
			isReverse: true
		}
		if (this.vCond != null && typeof this.vCond.value == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value"/>
	<p class="comment">要排除的完整的URL路径，通常以<code-label>/</code-label>开头，比如<code-label>/static/ui.js</code-label>。</p>
</div>`
})

// URL正则匹配
Vue.component("http-cond-url-regexp", {
	props: ["v-cond"],
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPath}",
			operator: "regexp",
			value: ""
		}
		if (this.vCond != null && typeof this.vCond.value == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value"/>
	<p class="comment">匹配URL的正则表达式，比如<code-label>^/static/(.*).js$</code-label>。</p>
</div>`
})

// 排除URL正则匹配
Vue.component("http-cond-url-not-regexp", {
	props: ["v-cond"],
	data: function () {
		let cond = {
			isRequest: true,
			param: "${requestPath}",
			operator: "not regexp",
			value: ""
		}
		if (this.vCond != null && typeof this.vCond.value == "string") {
			cond.value = this.vCond.value
		}
		return {
			cond: cond
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value"/>
	<p class="comment"><strong>不要</strong>匹配URL的正则表达式，意即只要匹配成功则排除此条件，比如<code-label>^/static/(.*).js$</code-label>。</p>
</div>`
})

// 根据MimeType
Vue.component("http-cond-mime-type", {
	props: ["v-cond"],
	data: function () {
		let cond = {
			isRequest: false,
			param: "${response.contentType}",
			operator: "mime type",
			value: "[]"
		}
		if (this.vCond != null && this.vCond.param == cond.param) {
			cond.value = this.vCond.value
		}
		return {
			cond: cond,
			mimeTypes: JSON.parse(cond.value), // TODO 可以拖动排序

			isAdding: false,
			addingMimeType: ""
		}
	},
	watch: {
		mimeTypes: function () {
			this.cond.value = JSON.stringify(this.mimeTypes)
		}
	},
	methods: {
		addMimeType: function () {
			this.isAdding = !this.isAdding

			if (this.isAdding) {
				let that = this
				setTimeout(function () {
					that.$refs.addingMimeType.focus()
				}, 100)
			}
		},
		cancelAdding: function () {
			this.isAdding = false
			this.addingMimeType = ""
		},
		confirmAdding: function () {
			// TODO 做更详细的校验
			// TODO 如果有重复的则提示之

			if (this.addingMimeType.length == 0) {
				return
			}
			this.addingMimeType = this.addingMimeType.replace(/\s+/g, "")
			this.mimeTypes.push(this.addingMimeType)

			// 清除状态
			this.cancelAdding()
		},
		removeMimeType: function (index) {
			this.mimeTypes.$remove(index)
		}
	},
	template: `<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<div v-if="mimeTypes.length > 0">
		<div class="ui label small" v-for="(mimeType, index) in mimeTypes">{{mimeType}} <a href="" title="删除" @click.prevent="removeMimeType(index)"><i class="icon remove"></i></a></div>
		<div class="ui divider"></div>
	</div>
	<div class="ui fields inline" v-if="isAdding">
		<div class="ui field">
			<input type="text" size="16" maxlength="100" v-model="addingMimeType" ref="addingMimeType" placeholder="类似于image/png" @keyup.enter="confirmAdding" @keypress.enter.prevent="1" />
		</div>
		<div class="ui field">
			<button class="ui button tiny" type="button" @click.prevent="confirmAdding">确认</button>
			<a href="" title="取消" @click.prevent="cancelAdding"><i class="icon remove"></i></a>
		</div> 
	</div>
	<div style="margin-top: 1em">
		<button class="ui button tiny" type="button" @click.prevent="addMimeType()">+添加MimeType</button>
	</div>
	<p class="comment">服务器返回的内容的MimeType，比如<span class="ui label tiny">text/html</span>、<span class="ui label tiny">image/*</span>等。</p>
</div>`
})