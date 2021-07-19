// 指标对象
Vue.component("metric-keys-config-box", {
	props: ["v-keys"],
	data: function () {
		let keys = this.vKeys
		if (keys == null) {
			keys = []
		}
		return {
			keys: keys,
			isAdding: false,
			key: "",
			subKey: "",
			keyDescription: "",

			keyDefs: window.METRIC_HTTP_KEYS
		}
	},
	watch: {
		keys: function () {
			this.$emit("change", this.keys)
		}
	},
	methods: {
		cancel: function () {
			this.key = ""
			this.subKey = ""
			this.keyDescription = ""
			this.isAdding = false
		},
		confirm: function () {
			if (this.key.length == 0) {
				return
			}

			if (this.key.indexOf(".NAME") > 0) {
				if (this.subKey.length == 0) {
					teaweb.warn("请输入参数值")
					return
				}
				this.key = this.key.replace(".NAME", "." + this.subKey)
			}
			this.keys.push(this.key)
			this.cancel()
		},
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				if (that.$refs.key != null) {
					that.$refs.key.focus()
				}
			}, 100)
		},
		remove: function (index) {
			this.keys.$remove(index)
		},
		changeKey: function () {
			if (this.key.length == 0) {
				return
			}
			let that = this
			let def = this.keyDefs.$find(function (k, v) {
				return v.code == that.key
			})
			if (def != null) {
				this.keyDescription = def.description
			}
		},
		keyName: function (key) {
			let that = this
			let subKey = ""
			let def = this.keyDefs.$find(function (k, v) {
				if (v.code == key) {
					return true
				}
				if (key.startsWith("${arg.") && v.code.startsWith("${arg.")) {
					subKey = that.getSubKey("arg.", key)
					return true
				}
				if (key.startsWith("${header.") && v.code.startsWith("${header.")) {
					subKey = that.getSubKey("header.", key)
					return true
				}
				if (key.startsWith("${cookie.") && v.code.startsWith("${cookie.")) {
					subKey = that.getSubKey("cookie.", key)
					return true
				}
				return false
			})
			if (def != null) {
				if (subKey.length > 0) {
					return def.name + ": " + subKey
				}
				return def.name
			}
			return key
		},
		getSubKey: function (prefix, key) {
			prefix = "${" + prefix
			let index = key.indexOf(prefix)
			if (index >= 0) {
				key = key.substring(index + prefix.length)
				key = key.substring(0, key.length - 1)
				return key
			}
			return ""
		}
	},
	template: `<div>
	<input type="hidden" name="keysJSON" :value="JSON.stringify(keys)"/>
	<div>
		<div v-for="(key, index) in keys" class="ui label small basic">
			{{keyName(key)}} &nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
	</div>
	<div v-if="isAdding" style="margin-top: 1em">
		<div class="ui fields inline">
			<div class="ui field">
				<select class="ui dropdown" v-model="key" @change="changeKey">
					<option value="">[选择对象]</option>
					<option v-for="def in keyDefs" :value="def.code">{{def.name}}</option>
				</select>
			</div>
			<div class="ui field" v-if="key == '\${arg.NAME}'">
				<input type="text" v-model="subKey" placeholder="参数名" size="15"/>
			</div>
			<div class="ui field" v-if="key == '\${header.NAME}'">
				<input type="text" v-model="subKey" placeholder="Header名" size="15">
			</div>
			<div class="ui field" v-if="key == '\${cookie.NAME}'">
				<input type="text" v-model="subKey" placeholder="Cookie名" size="15">
			</div>
			<div class="ui field">
				<button type="button" class="ui button tiny" @click.prevent="confirm">确定</button>
				<a href="" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
		<p class="comment" v-if="keyDescription.length > 0">{{keyDescription}}</p>
	</div>
	<div style="margin-top: 1em" v-if="!isAdding">
		<button type="button" class="ui button tiny" @click.prevent="add">+</button>
	</div>
</div>`
})