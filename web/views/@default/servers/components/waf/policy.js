Tea.context(function () {
	this.upgradeTemplate = function () {
		teaweb.confirm("确定要加入这些新规则吗？", function () {
			this.$post(".upgradeTemplate")
				.params({
					policyId: this.firewallPolicy.id
				})
				.refresh()
		})
	}
})