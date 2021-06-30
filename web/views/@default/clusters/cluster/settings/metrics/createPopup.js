Tea.context(function () {
	this.addItem = function (item) {
		this.$post("$")
			.params({
				clusterId: this.clusterId,
				itemId: item.id
			})
			.success(function () {
				item.isChecked = true
			})
	}

	this.removeItem = function (item) {
		this.$post(".delete")
			.params({
				clusterId: this.clusterId,
				itemId: item.id
			})
			.success(function () {
				item.isChecked = false
			})
	}
})