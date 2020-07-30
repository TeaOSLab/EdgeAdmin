Tea.context(function () {
	this.deleteNode = function (nodeId) {
		teaweb.confirm("确定要删除这个节点吗？", function () {
			this.$post(".delete")
				.params({
					nodeId: nodeId
				})
				.refresh();
		});
	};
});