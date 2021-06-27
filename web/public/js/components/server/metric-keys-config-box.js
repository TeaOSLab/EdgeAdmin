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
			key: ""
		}
	},
	methods: {
		cancel: function () {
			this.key = ""
			this.isAdding = false
		},
		confirm: function () {
			if (this.key.length > 0) {
				this.keys.push(this.key)
				this.cancel()
			}
		},
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.key.focus()
			}, 100)
		},
		remove: function (index) {
			this.keys.$remove(index)
		}
	},
	template: `<div>
	<input type="hidden" name="keysJSON" :value="JSON.stringify(keys)"/>
	<div>
		<div v-for="(key, index) in keys" class="ui label small basic">
			{{key}} &nbsp; <a href="" title="删除" @click.prevent="remove"><i class="icon remove small"></i></a>
		</div>
	</div>
	<div v-if="isAdding" style="margin-top: 1em">
		<div class="ui fields inline">
			<div class="ui field">
				<input type="text" v-model="key" ref="key" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
			</div>
			<div class="ui field">
				<button type="button" class="ui button tiny" @click.prevent="confirm">确定</button>
				<a href="" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
	</div>
	<div style="margin-top: 1em" v-if="!isAdding">
		<button type="button" class="ui button tiny" @click.prevent="add">+</button>
	</div>
</div>`
})