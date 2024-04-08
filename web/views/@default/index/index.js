Tea.context(function () {
	this.username = ""
	this.password = ""
	this.passwordMd5 = ""
	this.encodedFrom = window.encodeURIComponent(this.from)

	if (this.isDemo) {
		this.username = "admin"
		this.password = "123456"
	}

	this.isSubmitting = false

	this.$delay(function () {
		this.$find("form input[name='username']").focus()
		this.changePassword()
	});

	this.changeUsername = function () {

	}

	this.changePassword = function () {
		this.passwordMd5 = md5(this.password.trim());
	};

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
		// store information to local
		localStorage.setItem("sid", resp.data.localSid)
		localStorage.setItem("ip", resp.data.ip)

		// redirect back
		this.$delay(function () {
			if (resp.data.requireOTP) {
				window.location = "/index/otp?sid=" + resp.data.sid + "&remember=" + (resp.data.remember ? 1 : 0) + "&from=" + window.encodeURIComponent(this.from)
				return
			}
			if (this.from.length == 0) {
				window.location = "/dashboard";
			} else {
				window.location = this.from;
			}
		})
	};
});