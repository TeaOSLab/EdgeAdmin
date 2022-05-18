Vue.component("ddos-protection-ip-list-config-box", {
	props: ["v-ip-list"],
	data: function () {
		let list = this.vIpList
		if (list == null) {
			list = []
		}
		return {
			list: list,
			isAdding: false,
			addingIP: {
				ip: "",
				description: ""
			}
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.addingIPInput.focus()
			})
		},
		confirm: function () {
			let ip = this.addingIP.ip
			if (ip.length == 0) {
				this.warn("请输入IP")
				return
			}

			let exists = false
			this.list.forEach(function (v) {
				if (v.ip == ip) {
					exists = true
				}
			})
			if (exists) {
				this.warn("IP '" + ip + "'已经存在")
				return
			}

			let that = this
			Tea.Vue.$post("/ui/validateIPs")
				.params({
					ips: [ip]
				})
				.success(function () {
					that.list.push({
						ip: ip,
						description: that.addingIP.description
					})
					that.notifyChange()
					that.cancel()
				})
				.fail(function () {
					that.warn("请输入正确的IP")
				})
		},
		cancel: function () {
			this.isAdding = false
			this.addingIP = {
				ip: "",
				description: ""
			}
		},
		remove: function (index) {
			this.list.$remove(index)
			this.notifyChange()
		},
		warn: function (message) {
			let that = this
			teaweb.warn(message, function () {
				that.$refs.addingIPInput.focus()
			})
		},
		notifyChange: function () {
			this.$emit("change", this.list)
		}
	},
	template: `<div>
	<div v-if="list.length > 0">
		<div class="ui label basic tiny" v-for="(ipConfig, index) in list">
			{{ipConfig.ip}} <span class="grey small" v-if="ipConfig.description.length > 0">（{{ipConfig.description}}）</span> <a href="" @click.prevent="remove(index)" title="删除"><i class="icon remove"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding">
		<div class="ui fields inline">
			<div class="ui field">
				<div class="ui input left labeled">
					<span class="ui label">IP</span>
					<input type="text" v-model="addingIP.ip" ref="addingIPInput" maxlength="40" size="20" placeholder="IP" @keyup.enter="confirm" @keypress.enter.prevent="1"/>
				</div>
			</div>
			<div class="ui field">
				<div class="ui input left labeled">
					<span class="ui label">备注</span>
					<input type="text" v-model="addingIP.description" maxlength="10" size="10" placeholder="备注（可选）" @keyup.enter="confirm" @keypress.enter.prevent="1"/>
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