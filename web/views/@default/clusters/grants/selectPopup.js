Tea.context(function () {
	this.grantId = 0;

	this.selectGrant = function (grant) {
		NotifyPopup({
			code: 200,
			data: {
				grant: grant
			}
		})
	};
});