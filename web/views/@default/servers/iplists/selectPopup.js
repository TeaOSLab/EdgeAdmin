Tea.context(function () {
	this.selectList = function (list) {
		NotifyPopup({
			code: 200,
			data: {
				list: {
					id: list.id,
					name: list.name
				}
			}
		})
	}
})