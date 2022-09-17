Vue.component("dns-domain-selector", {
	props: ["v-domain-id", "v-domain-name", "v-provider-name"],
	data: function () {
		let domainId = this.vDomainId
		if (domainId == null) {
			domainId = 0
		}
		let domainName = this.vDomainName
		if (domainName == null) {
			domainName = ""
		}

		let providerName = this.vProviderName
		if (providerName == null) {
			providerName = ""
		}

		return {
			domainId: domainId,
			domainName: domainName,
			providerName: providerName
		}
	},
	methods: {
		select: function () {
			let that = this
			teaweb.popup("/dns/domains/selectPopup", {
				callback: function (resp) {
					that.domainId = resp.data.domainId
					that.domainName = resp.data.domainName
					that.providerName = resp.data.providerName
					that.change()
				}
			})
		},
		remove: function() {
			this.domainId = 0
			this.domainName = ""
			this.change()
		},
		update: function () {
			let that = this
			teaweb.popup("/dns/domains/selectPopup?domainId=" + this.domainId, {
				callback: function (resp) {
					that.domainId = resp.data.domainId
					that.domainName = resp.data.domainName
					that.providerName = resp.data.providerName
					that.change()
				}
			})
		},
		change: function () {
			this.$emit("change", {
				id: this.domainId,
				name: this.domainName
			})
		}
	},
	template: `<div>
	<input type="hidden" name="dnsDomainId" :value="domainId"/>
	<div v-if="domainName.length > 0">
		<span class="ui label small basic">
			<span v-if="providerName != null && providerName.length > 0">{{providerName}} &raquo; </span> {{domainName}}
			<a href="" @click.prevent="update"><i class="icon pencil small"></i></a>
			<a href="" @click.prevent="remove()"><i class="icon remove"></i></a>
		</span>
	</div>
	<div v-if="domainName.length == 0">
		<a href="" @click.prevent="select()">[选择域名]</a>
	</div>
</div>`
})