Vue.component("server-name-box", {
	props: ["v-server-names"],
	data: function () {
		let serverNames = this.vServerNames;
		if (serverNames == null) {
			serverNames = []
		}
		return {
			serverNames: serverNames,
			isSearching: false,
			keyword: ""
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
					setTimeout(that.submitForm, 100)
				}
			});
		},

		removeServerName: function (index) {
			this.serverNames.$remove(index)
		},

		updateServerName: function (index, serverName) {
			window.UPDATING_SERVER_NAME = teaweb.clone(serverName)
			let that = this
			teaweb.popup("/servers/addServerNamePopup", {
				callback: function (resp) {
					var serverName = resp.data.serverName
					Vue.set(that.serverNames, index, serverName)
					setTimeout(that.submitForm, 100)
				}
			});
		},
		showSearchBox: function () {
			this.isSearching = !this.isSearching
			if (this.isSearching) {
				let that = this
				setTimeout(function () {
					that.$refs.keywordRef.focus()
				}, 200)
			} else {
				this.keyword = ""
			}
		},
		allServerNames: function () {
			if (this.serverNames == null) {
				return []
			}
			let result = []
			this.serverNames.forEach(function (serverName) {
				if (serverName.subNames != null && serverName.subNames.length > 0) {
					serverName.subNames.forEach(function (subName) {
						if (subName != null && subName.length > 0) {
							if (!result.$contains(subName)) {
								result.push(subName)
							}
						}
					})
				} else if (serverName.name != null && serverName.name.length > 0) {
					if (!result.$contains(serverName.name)) {
						result.push(serverName.name)
					}
				}
			})
			return result
		},
		submitForm: function () {
			Tea.runActionOn(this.$refs.serverNamesRef.form)
		}
	},
	watch: {
		keyword: function (v) {
			this.serverNames.forEach(function (serverName) {
				if (v.length == 0) {
					serverName.isShowing = true
					return
				}
				if (serverName.subNames == null || serverName.subNames.length == 0) {
					if (!teaweb.match(serverName.name, v)) {
						serverName.isShowing = false
					}
				} else {
					let found = false
					serverName.subNames.forEach(function (subName) {
						if (teaweb.match(subName, v)) {
							found = true
						}
					})
					serverName.isShowing = found
				}
			})
		}
	},
	template: `<div>
	<input type="hidden" name="serverNames" :value="JSON.stringify(serverNames)" ref="serverNamesRef"/>
	<div v-if="serverNames.length > 0">
		<div v-for="(serverName, index) in serverNames" class="ui label small basic" :class="{hidden: serverName.isShowing === false}">
			<em v-if="serverName.type != 'full'">{{serverName.type}}</em>  
			<span v-if="serverName.subNames == null || serverName.subNames.length == 0" :class="{disabled: serverName.isShowing === false}">{{serverName.name}}</span>
			<span v-else :class="{disabled: serverName.isShowing === false}">{{serverName.subNames[0]}}等{{serverName.subNames.length}}个域名</span>
			<a href="" title="修改" @click.prevent="updateServerName(index, serverName)"><i class="icon pencil small"></i></a> <a href="" title="删除" @click.prevent="removeServerName(index)"><i class="icon remove"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div class="ui fields inline">
	    <div class="ui field"><a href="" @click.prevent="addServerName()">[添加域名绑定]</a></div>
	    <div class="ui field" v-if="serverNames.length > 0"><span class="grey">|</span> </div>
	    <div class="ui field" v-if="serverNames.length > 0">
	        <a href="" @click.prevent="showSearchBox()" v-if="!isSearching"><i class="icon search small"></i></a>
	        <a href="" @click.prevent="showSearchBox()" v-if="isSearching"><i class="icon close small"></i></a>
        </div>
        <div class="ui field" v-if="isSearching">
            <input type="text" placeholder="搜索域名" ref="keywordRef" class="ui input tiny" v-model="keyword"/>
        </div>
    </div>
</div>`
})