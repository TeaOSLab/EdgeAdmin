Vue.component("http-firewall-provider-selector", {
	props: ["v-type", "v-providers"],
	data: function () {
		let providers = this.vProviders
		if (providers == null) {
			providers = []
		}

		return {
			listType: this.vType,
			providers: providers
		}
	},
	methods: {
		addProvider: function () {
			let selectedProviderIds = this.providers.map(function (provider) {
				return provider.id
			})
			let that = this
			teaweb.popup("/servers/server/settings/waf/ipadmin/selectProvidersPopup?type=" + this.listType + "&selectedProviderIds=" + selectedProviderIds.join(","), {
				width: "50em",
				height: "26em",
				callback: function (resp) {
					that.providers = resp.data.selectedProviders
					that.$forceUpdate()
					that.notifyChange()
				}
			})
		},
		removeProvider: function (index) {
			this.providers.$remove(index)
			this.notifyChange()
		},
		resetProviders: function () {
			this.providers = []
			this.notifyChange()
		},
		notifyChange: function () {
			this.$emit("change", {
				"providers": this.providers
			})
		}
	},
	template: `<div>
	<span v-if="providers.length == 0" class="disabled">暂时没有选择<span v-if="listType =='allow'">允许</span><span v-else>封禁</span>的运营商。</span>
	<div v-show="providers.length > 0">
		<div class="ui label tiny basic" v-for="(provider, index) in providers" style="margin-bottom: 0.5em">
			<input type="hidden" :name="listType + 'ProviderIds'" :value="provider.id"/>
			{{provider.name}} <a href="" @click.prevent="removeProvider(index)" title="删除"><i class="icon remove"></i></a>
		</div>
	</div>
	<div class="ui divider"></div>
	<button type="button" class="ui button tiny" @click.prevent="addProvider">修改</button> &nbsp; <button type="button" class="ui button tiny" v-show="providers.length > 0" @click.prevent="resetProviders">清空</button>
</div>`
})