Vue.component("server-name-box", {
	props: ["v-server-names"],
	data: function () {
		let serverNames = this.vServerNames;
		if (serverNames == null) {
			serverNames = []
		}
		return {
			serverNames: serverNames
		}
	},
	methods: {
		addServerName: function () {
			let that = this
			teaweb.popup("/servers/addServerNamePopup", {
				callback: function (resp) {
					var serverName = resp.data.serverName
					that.serverNames.push(serverName)
				}
			});
		},

		removeServerName: function (index) {
			this.serverNames.$remove(index)
		}
	},
	template: `<div>
	<input type="hidden" name="serverNames" :value="JSON.stringify(serverNames)"/>
	<div v-if="serverNames.length > 0">
		<div v-for="(serverName, index) in serverNames" class="ui label small">
			<em v-if="serverName.type != 'full'">{{serverName.type}}</em>  {{serverName.name}} <a href="" title="删除" @click.prevent="removeServerName(index)"><i class="icon remove"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<a href="" @click.prevent="addServerName()">[添加域名绑定]</a>
</div>`
})