Tea.context(function () {
	this.bind = function (list) {
		this.$post("$")
			.params({
				httpFirewallPolicyId: this.httpFirewallPolicyId,
				listId: list.id
			})
			.success(function () {
				list.isSelected = true
			})
	}

	this.unbind = function (list) {
		this.$post(".unbindHTTPFirewall")
			.params({
				httpFirewallPolicyId: this.httpFirewallPolicyId,
				listId: list.id
			})
			.success(function () {
				list.isSelected = false
			})
	}
})