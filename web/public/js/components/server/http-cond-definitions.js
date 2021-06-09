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

// 参数匹配
Vue.component("http-cond-params", {
	props: ["v-cond"],
	mounted: function () {
		let cond = this.vCond
		if (cond == null) {
			return
		}
		this.operator = cond.operator

		// stringValue
		if (["regexp", "not regexp", "eq", "not", "prefix", "suffix", "contains", "not contains", "eq ip", "gt ip", "gte ip", "lt ip", "lte ip", "ip range"].$contains(cond.operator)) {
			this.stringValue = cond.value
			return
		}

		// numberValue
		if (["eq int", "eq float", "gt", "gte", "lt", "lte", "mod 10", "ip mod 10", "mod 100", "ip mod 100"].$contains(cond.operator)) {
			this.numberValue = cond.value
			return
		}

		// modValue
		if (["mod", "ip mod"].$contains(cond.operator)) {
			let pieces = cond.value.split(",")
			this.modDivValue = pieces[0]
			if (pieces.length > 1) {
				this.modRemValue = pieces[1]
			}
			return
		}

		// stringValues
		let that = this
		if (["in", "not in", "file ext", "mime type"].$contains(cond.operator)) {
			try {
				let arr = JSON.parse(cond.value)
				if (arr != null && (arr instanceof Array)) {
					arr.forEach(function (v) {
						that.stringValues.push(v)
					})
				}
			} catch (e) {

			}
			return
		}

		// versionValue
		if (["version range"].$contains(cond.operator)) {
			let pieces = cond.value.split(",")
			this.versionRangeMinValue = pieces[0]
			if (pieces.length > 1) {
				this.versionRangeMaxValue = pieces[1]
			}
			return
		}
	},
	data: function () {
		let cond = {
			isRequest: true,
			param: "",
			operator: window.REQUEST_COND_OPERATORS[0].op,
			value: ""
		}
		if (this.vCond != null) {
			cond = this.vCond
		}
		return {
			cond: cond,
			operators: window.REQUEST_COND_OPERATORS,
			operator: window.REQUEST_COND_OPERATORS[0].op,
			operatorDescription: window.REQUEST_COND_OPERATORS[0].description,
			variables: window.REQUEST_VARIABLES,
			variable: "",

			// 各种类型的值
			stringValue: "",
			numberValue: "",

			modDivValue: "",
			modRemValue: "",

			stringValues: [],

			versionRangeMinValue: "",
			versionRangeMaxValue: ""
		}
	},
	methods: {
		changeVariable: function () {
			let v = this.cond.param
			if (v == null) {
				v = ""
			}
			this.cond.param = v + this.variable
		},
		changeOperator: function () {
			let that = this
			this.operators.forEach(function (v) {
				if (v.op == that.operator) {
					that.operatorDescription = v.description
				}
			})

			this.cond.operator = this.operator

			// 移动光标
			let box = document.getElementById("variables-value-box")
			if (box != null) {
				setTimeout(function () {
					let input = box.getElementsByTagName("INPUT")
					if (input.length > 0) {
						input[0].focus()
					}
				}, 100)
			}
		},
		changeStringValues: function (v) {
			this.stringValues = v
			this.cond.value = JSON.stringify(v)
		}
	},
	watch: {
		stringValue: function (v) {
			this.cond.value = v
		},
		numberValue: function (v) {
			// TODO 校验数字
			this.cond.value = v
		},
		modDivValue: function (v) {
			if (v.length == 0) {
				return
			}
			let div = parseInt(v)
			if (isNaN(div)) {
				div = 1
			}
			this.modDivValue = div
			this.cond.value = div + "," + this.modRemValue
		},
		modRemValue: function (v) {
			if (v.length == 0) {
				return
			}
			let rem = parseInt(v)
			if (isNaN(rem)) {
				rem = 0
			}
			this.modRemValue = rem
			this.cond.value = this.modDivValue + "," + rem
		},
		versionRangeMinValue: function (v) {
			this.cond.value = this.versionRangeMinValue + "," + this.versionRangeMaxValue
		},
		versionRangeMaxValue: function (v) {
			this.cond.value = this.versionRangeMinValue + "," + this.versionRangeMaxValue
		}
	},
	template: `<tbody>
	<tr>
		<td>参数值</td>
		<td>
			<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
			<div>
				<div class="ui fields inline">
					<div class="ui field">
						<input type="text" placeholder="\${xxx}" v-model="cond.param"/>
					</div>
					<div class="ui field">
						<select class="ui dropdown" style="width: 7em; color: grey" v-model="variable" @change="changeVariable">
							<option value="">[常用参数]</option>
							<option v-for="v in variables" :value="v.code">{{v.code}} - {{v.name}}</option>
						</select>
					</div>
				</div>			
			</div>
			<p class="comment">其中可以使用变量，类似于<code-label>\${requestPath}</code-label>，也可以是多个变量的组合。</p>
		</td>
	</tr>
	<tr>
		<td>操作符</td>
		<td>
			<div>
				<select class="ui dropdown auto-width" v-model="operator" @change="changeOperator">
					<option v-for="operator in operators" :value="operator.op">{{operator.name}}</option>
				</select>
				<p class="comment">{{operatorDescription}}</p>
			</div>
		</td>
	</tr>
	<tr v-show="!['file exist', 'file not exist'].$contains(cond.operator)">
		<td>对比值</td>
		<td id="variables-value-box">
			<!-- 正则表达式 -->
			<div v-if="['regexp', 'not regexp'].$contains(cond.operator)">
				<input type="text" v-model="stringValue"/>
				<p class="comment">要匹配的正则表达式，比如<code-label>^/static/(.+).js</code-label>。</p>
			</div>
			
			<!-- 数字相关 -->
			<div v-if="['eq int', 'eq float', 'gt', 'gte', 'lt', 'lte'].$contains(cond.operator)">
				<input type="text" maxlength="11" size="11" style="width: 5em" v-model="numberValue"/>
				<p class="comment">要对比的数字。</p>
			</div>
			
			<!-- 取模 -->
			<div v-if="['mod 10'].$contains(cond.operator)">
				<input type="text" maxlength="11" size="11" style="width: 5em" v-model="numberValue"/>
				<p class="comment">参数值除以10的余数，在0-9之间。</p>
			</div>
			<div v-if="['mod 100'].$contains(cond.operator)">
				<input type="text" maxlength="11" size="11" style="width: 5em" v-model="numberValue"/>
				<p class="comment">参数值除以100的余数，在0-99之间。</p>
			</div>
			<div v-if="['mod', 'ip mod'].$contains(cond.operator)">
				<div class="ui fields inline">
					<div class="ui field">除：</div>
					<div class="ui field">
						<input type="text" maxlength="11" size="11" style="width: 5em" v-model="modDivValue" placeholder="除数"/>
					</div>
					<div class="ui field">余：</div>
					<div class="ui field">
						<input type="text" maxlength="11" size="11" style="width: 5em" v-model="modRemValue" placeholder="余数"/>
					</div>
				</div>
			</div>
			
			<!-- 字符串相关 -->
			<div v-if="['eq', 'not', 'prefix', 'suffix', 'contains', 'not contains'].$contains(cond.operator)">
				<input type="text" v-model="stringValue"/>
				<p class="comment" v-if="cond.operator == 'eq'">和参数值一致的字符串。</p>
				<p class="comment" v-if="cond.operator == 'not'">和参数值不一致的字符串。</p>
				<p class="comment" v-if="cond.operator == 'prefix'">参数值的前缀。</p>
				<p class="comment" v-if="cond.operator == 'suffix'">参数值的后缀为此字符串。</p>
				<p class="comment" v-if="cond.operator == 'contains'">参数值包含此字符串。</p>
				<p class="comment" v-if="cond.operator == 'not contains'">参数值不包含此字符串。</p>
			</div>
			<div v-if="['in', 'not in', 'file ext', 'mime type'].$contains(cond.operator)">
				<values-box @change="changeStringValues" :values="stringValues" size="15"></values-box>
				<p class="comment" v-if="cond.operator == 'in'">添加参数值列表。</p>
				<p class="comment" v-if="cond.operator == 'not in'">添加参数值列表。</p>
				<p class="comment" v-if="cond.operator == 'file ext'">添加扩展名列表，比如<code-label>png</code-label>、<code-label>html</code-label>，不包括点。</p>
				<p class="comment" v-if="cond.operator == 'mime type'">添加MimeType列表，类似于<code-label>text/html</code-label>、<code-label>image/*</code-label>。</p>
			</div>
			<div v-if="['version range'].$contains(cond.operator)">
				<div class="ui fields inline">
					<div class="ui field"><input type="text" v-model="versionRangeMinValue" maxlength="200" placeholder="最小版本" style="width: 10em"/></div>
					<div class="ui field">-</div>
					<div class="ui field"><input type="text" v-model="versionRangeMaxValue" maxlength="200" placeholder="最大版本" style="width: 10em"/></div>
				</div>
			</div>
			
			<!-- IP相关 -->
			<div v-if="['eq ip', 'gt ip', 'gte ip', 'lt ip', 'lte ip', 'ip range'].$contains(cond.operator)">
				<input type="text" style="width: 10em" v-model="stringValue" placeholder="x.x.x.x"/>
				<p class="comment">要对比的IP。</p>
			</div>
			<div v-if="['ip mod 10'].$contains(cond.operator)">
				<input type="text" maxlength="11" size="11" style="width: 5em" v-model="numberValue"/>
				<p class="comment">参数中IP转换成整数后除以10的余数，在0-9之间。</p>
			</div>
			<div v-if="['ip mod 100'].$contains(cond.operator)">
				<input type="text" maxlength="11" size="11" style="width: 5em" v-model="numberValue"/>
				<p class="comment">参数中IP转换成整数后除以100的余数，在0-99之间。</p>
			</div>
		</td>
	</tr>
</tbody>`
})