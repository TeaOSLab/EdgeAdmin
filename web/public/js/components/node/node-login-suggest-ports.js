// 节点登录推荐端口
Vue.component("node-login-suggest-ports", {
	data: function () {
		return {
			ports: [],
			availablePorts: [],
			autoSelected: false,
			isLoading: false
		}
	},
	methods: {
		reload: function (host) {
			let that = this
			this.autoSelected = false
			this.isLoading = true
			Tea.action("/clusters/cluster/suggestLoginPorts")
				.params({
					host: host
				})
				.success(function (resp) {
					if (resp.data.availablePorts != null) {
						that.availablePorts = resp.data.availablePorts
						if (that.availablePorts.length > 0) {
							that.autoSelectPort(that.availablePorts[0])
							that.autoSelected = true
						}
					}
					if (resp.data.ports != null) {
						that.ports = resp.data.ports
						if (that.ports.length > 0 && !that.autoSelected) {
							that.autoSelectPort(that.ports[0])
							that.autoSelected = true
						}
					}
				})
				.done(function () {
					that.isLoading = false
				})
				.post()
		},
		selectPort: function (port) {
			this.$emit("select", port)
		},
		autoSelectPort: function (port) {
			this.$emit("auto-select", port)
		}
	},
	template: `<span>
	<span v-if="isLoading">正在检查端口...</span>
	<span v-if="availablePorts.length > 0">
		可能端口：<a href="" v-for="port in availablePorts" @click.prevent="selectPort(port)" class="ui label tiny basic blue" style="border: 1px #2185d0 dashed; font-weight: normal">{{port}}</a>
		&nbsp; &nbsp;
	</span>
	<span v-if="ports.length > 0">
		常用端口：<a href="" v-for="port in ports" @click.prevent="selectPort(port)" class="ui label tiny basic blue" style="border: 1px #2185d0 dashed;  font-weight: normal">{{port}}</a>
	</span>
	<span v-if="ports.length == 0">常用端口有22等。</span>
	<span v-if="ports.length > 0" class="grey small">（可以点击要使用的端口）</span>
</span>`
})