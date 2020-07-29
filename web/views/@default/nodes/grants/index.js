Tea.context(function () {
	this.deleteGrant = function (grantId) {
		teaweb.confirm("确定要删除此认证吗？", function () {
			this.$post(".delete")
				.params({
					"grantId": grantId
				})
				.refresh();
		});
	};
});