Vue.component("ddos-protection-ports-config-box", {
	props: ["v-ports"],
	data: function () {
		let ports = this.vPorts
		if (ports == null) {
			ports = []
		}
		return {
			ports: ports,
			isAdding: false,
			addingPort: {
				port: "",
				description: ""
			}
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.addingPortInput.focus()
			})
		},
		confirm: function () {
			let portString = this.addingPort.port
			if (portString.length == 0) {
				this.warn("请输入端口号")
				return
			}
			if (!/^\d+$/.test(portString)) {
				this.warn("请输入正确的端口号")
				return
			}
			let port = parseInt(portString, 10)
			if (port <= 0) {
				this.warn("请输入正确的端口号")
				return
			}
			if (port > 65535) {
				this.warn("请输入正确的端口号")
				return
			}

			let exists = false
			this.ports.forEach(function (v) {
				if (v.port == port) {
					exists = true
				}
			})
			if (exists) {
				this.warn("端口号已经存在")
				return
			}

			this.ports.push({
				port: port,
				description: this.addingPort.description
			})
			this.notifyChange()
			this.cancel()
		},
		cancel: function () {
			this.isAdding = false
			this.addingPort = {
				port: "",
				description: ""
			}
		},
		remove: function (index) {
			this.ports.$remove(index)
			this.notifyChange()
		},
		warn: function (message) {
			let that = this
			teaweb.warn(message, function () {
				that.$refs.addingPortInput.focus()
			})
		},
		notifyChange: function () {
			this.$emit("change", this.ports)
		}
	},
	template: `<div>
	<div v-if="ports.length > 0">
		<div class="ui label basic tiny" v-for="(portConfig, index) in ports">
			{{portConfig.port}} <span class="grey small" v-if="portConfig.description.length > 0">（{{portConfig.description}}）</span> <a href="" @click.prevent="remove(index)" title="删除"><i class="icon remove"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding">
		<div class="ui fields inline">
			<div class="ui field">
				<div class="ui input left labeled">
					<span class="ui label">端口</span>
					<input type="text" v-model="addingPort.port" ref="addingPortInput" maxlength="5" size="5" placeholder="端口号" @keyup.enter="confirm" @keypress.enter.prevent="1"/>
				</div>
			</div>
			<div class="ui field">
				<div class="ui input left labeled">
					<span class="ui label">备注</span>
					<input type="text" v-model="addingPort.description" maxlength="12" size="12" placeholder="备注（可选）" @keyup.enter="confirm" @keypress.enter.prevent="1"/>
				</div>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>
				&nbsp;<a href="" @click.prevent="cancel()">取消</a>
			</div>
		</div>
	</div>
	<div v-if="!isAdding">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`
})