// 绑定IP列表
Vue.component("ip-list-bind-box", {
	props: ["v-http-firewall-policy-id", "v-type"],
	mounted: function () {
		this.refresh()
	},
	data: function () {
		return {
			policyId: this.vHttpFirewallPolicyId,
			type: this.vType,
			lists: []
		}
	},
	methods: {
		bind: function () {
			let that = this
			teaweb.popup("/servers/iplists/bindHTTPFirewallPopup?httpFirewallPolicyId=" + this.policyId + "&type=" + this.type, {
				width: "50em",
				height: "34em",
				callback: function () {

				},
				onClose: function () {
					that.refresh()
				}
			})
		},
		remove: function (index, listId) {
			let that = this
			teaweb.confirm("确定要删除这个绑定的IP名单吗？", function () {
				Tea.action("/servers/iplists/unbindHTTPFirewall")
					.params({
						httpFirewallPolicyId: that.policyId,
						listId: listId
					})
					.post()
					.success(function (resp) {
						that.lists.$remove(index)
					})
			})
		},
		refresh: function () {
			let that = this
			Tea.action("/servers/iplists/httpFirewall")
				.params({
					httpFirewallPolicyId: this.policyId,
					type: this.vType
				})
				.post()
				.success(function (resp) {
					that.lists = resp.data.lists
				})
		}
	},
	template: `<div>
	<a href="" @click.prevent="bind()" style="color: rgba(0,0,0,.6)">绑定+</a> &nbsp; <span v-if="lists.length > 0"><span class="disabled small">|&nbsp;</span> 已绑定：</span>
	<div class="ui label basic small" v-for="(list, index) in lists">
		<a :href="'/servers/iplists/list?listId=' + list.id" title="点击查看详情" style="opacity: 1">{{list.name}}</a>
		<a href="" title="删除" @click.prevent="remove(index, list.id)"><i class="icon remove small"></i></a>
	</div>
</div>`
})