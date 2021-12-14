// 请求方法列表
Vue.component("http-methods-box", {
	props: ["v-methods"],
	data: function () {
		let methods = this.vMethods
		if (methods == null) {
			methods = []
		}
		return {
			methods: methods,
			isAdding: false,
			addingMethod: ""
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.addingMethod.focus()
			}, 100)
		},
		confirm: function () {
			let that = this

			// 删除其中的空格
			this.addingMethod = this.addingMethod.replace(/\s/g, "").toUpperCase()

			if (this.addingMethod.length == 0) {
				teaweb.warn("请输入要添加的请求方法", function () {
					that.$refs.addingMethod.focus()
				})
				return
			}

			// 是否已经存在
			if (this.methods.$contains(this.addingMethod)) {
				teaweb.warn("此请求方法已经存在，无需重复添加", function () {
					that.$refs.addingMethod.focus()
				})
				return
			}

			this.methods.push(this.addingMethod)
			this.cancel()
		},
		remove: function (index) {
			this.methods.$remove(index)
		},
		cancel: function () {
			this.isAdding = false
			this.addingMethod = ""
		}
	},
	template: `<div>
	<input type="hidden" name="methodsJSON" :value="JSON.stringify(methods)"/>
	<div v-if="methods.length > 0">
		<span class="ui label small basic" v-for="(method, index) in methods">
			{{method}}
			&nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</span>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding">
		<div class="ui fields">
			<div class="ui field">
				<input type="text" v-model="addingMethod" @keyup.enter="confirm()" @keypress.enter.prevent="1" ref="addingMethod" placeholder="如GET" size="10"/>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>
				&nbsp; <a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
		<p class="comment">格式为大写，比如<code-label>GET</code-label>、<code-label>POST</code-label>等。</p>
		<div class="ui divider"></div>
	</div>
	<div style="margin-top: 0.5em" v-if="!isAdding">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})