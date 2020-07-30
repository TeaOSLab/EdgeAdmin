Tea.context(function () {
	this.grantId = 0;

	this.selectGrant = function (grant) {
		this.grantId = grant.id;
	};

	this.success = NotifyPopup;
});