Vue.component("api-node-addresses-box", {
	props: ["v-addrs", "v-name"],
	data: function () {
		let addrs = this.vAddrs
		if (addrs == null) {
			addrs = []
		}
		return {
			addrs: addrs
		}
	},
	methods: {
		// 添加IP地址
		addAddr: function () {
			let that = this;
			teaweb.popup("/settings/api/node/createAddrPopup", {
				height: "16em",
				callback: function (resp) {
					that.addrs.push(resp.data.addr);
				}
			})
		},

		// 修改地址
		updateAddr: function (index, addr) {
			let that = this;
			window.UPDATING_ADDR = addr
			teaweb.popup("/settings/api/node/updateAddrPopup?addressId=", {
				callback: function (resp) {
					Vue.set(that.addrs, index, resp.data.addr);
				}
			})
		},

		// 删除IP地址
		removeAddr: function (index) {
			this.addrs.$remove(index);
		}
	},
	template: `<div>
	<input type="hidden" :name="vName" :value="JSON.stringify(addrs)"/>
	<div v-if="addrs.length > 0">
		<div>
			<div v-for="(addr, index) in addrs" class="ui label small basic">
				{{addr.protocol}}://{{addr.host.quoteIP()}}:{{addr.portRange}}</span>
				<a href="" title="修改" @click.prevent="updateAddr(index, addr)"><i class="icon pencil small"></i></a>
				<a href="" title="删除" @click.prevent="removeAddr(index)"><i class="icon remove"></i></a>
			</div>
		</div>
		<div class="ui divider"></div>
	</div>
	<div>
		<button class="ui button small" type="button" @click.prevent="addAddr()">+</button>
	</div>
</div>`
})