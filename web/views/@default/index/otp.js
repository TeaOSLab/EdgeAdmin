Tea.context(function () {
	this.isSubmitting = false

	this.encodedFrom = window.encodeURIComponent(this.from)

	this.$delay(function () {
		this.$find("form input[name='otpCode']").focus()
	});

	// 更多选项
	this.moreOptionsVisible = false;
	this.showMoreOptions = function () {
		this.moreOptionsVisible = !this.moreOptionsVisible;
	};

	this.submitBefore = function () {
		this.isSubmitting = true;
	};

	this.submitDone = function () {
		this.isSubmitting = false;
	};

	this.submitSuccess = function (resp) {
		if (this.from.length == 0) {
			window.location = "/dashboard";
		} else {
			window.location = this.from;
		}
	};
});