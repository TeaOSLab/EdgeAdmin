Vue.component("dns-resolvers-config-box", {
	props: ["value", "name"],
	data: function () {
		let resolvers = this.value
		if (resolvers == null) {
			resolvers = []
		}

		let name = this.name
		if (name == null || name.length == 0) {
			name = "dnsResolversJSON"
		}

		return {
			formName: name,
			resolvers: resolvers,

			host: "",

			isAdding: false
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.hostRef.focus()
			})
		},
		confirm: function () {
			let host = this.host.trim()
			if (host.length == 0) {
				let that = this
				setTimeout(function () {
					that.$refs.hostRef.focus()
				})
				return
			}
			this.resolvers.push({
				host: host,
				port: 0, // TODO
				protocol: "" // TODO
			})
			this.cancel()
		},
		cancel: function () {
			this.isAdding = false
			this.host = ""
			this.port = 0
			this.protocol = ""
		},
		remove: function (index) {
			this.resolvers.$remove(index)
		}
	},
	template: `<div>
	<input type="hidden" :name="formName" :value="JSON.stringify(resolvers)"/>
	<div v-if="resolvers.length > 0">
		<div v-for="(resolver, index) in resolvers" class="ui label basic small">
			<span v-if="resolver.protocol.length > 0">{{resolver.protocol}}</span>{{resolver.host}}<span v-if="resolver.port > 0">:{{resolver.port}}</span>
			&nbsp;
			<a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
	</div>
	
	<div v-if="isAdding" style="margin-top: 1em">
		<div class="ui fields inline">
			<div class="ui field">
				<input type="text" placeholder="x.x.x.x" @keyup.enter="confirm" @keypress.enter.prevent="1" ref="hostRef" v-model="host"/>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确认</button>
				&nbsp; <a href="" @click.prevent="cancel" title="取消"><i class="icon remove small"></i></a>
			</div>
		</div>
	</div>
	
	<div v-if="!isAdding" style="margin-top: 1em">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})