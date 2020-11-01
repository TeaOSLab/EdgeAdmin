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
			window.UPDATING_SERVER_NAME = null
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
		},

		updateServerName: function (index, serverName) {
			window.UPDATING_SERVER_NAME = serverName
			let that = this
			teaweb.popup("/servers/addServerNamePopup", {
				callback: function (resp) {
					var serverName = resp.data.serverName
					Vue.set(that.serverNames, index, serverName)
				}
			});
		}
	},
	template: `<div>
	<input type="hidden" name="serverNames" :value="JSON.stringify(serverNames)"/>
	<div v-if="serverNames.length > 0">
		<div v-for="(serverName, index) in serverNames" class="ui label small">
			<em v-if="serverName.type != 'full'">{{serverName.type}}</em>  
			<span v-if="serverName.subNames == null || serverName.subNames.length == 0">{{serverName.name}}</span>
			<span v-else>{{serverName.subNames[0]}}等{{serverName.subNames.length}}个域名</span>
			<a href="" title="修改" @click.prevent="updateServerName(index, serverName)"><i class="icon pencil small"></i></a> <a href="" title="删除" @click.prevent="removeServerName(index)"><i class="icon remove"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<a href="" @click.prevent="addServerName()">[添加域名绑定]</a>
</div>`
})