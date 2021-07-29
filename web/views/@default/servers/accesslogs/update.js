Tea.context(function () {
	this.success = NotifySuccess("html:保存成功<br/>新的配置将会在1分钟之内生效", ".policy", {policyId: this.policy.id})

	this.type = this.policy.type


	/**
	 * syslog
	 */
	this.syslogProtocol = this.policy.options.protocol

	/**
	 * command
	 */
	if (this.policy.type == "command") {
		let args = this.policy.options.args
		if (args == null) {
			args = ""
		} else {
			args = args.join(" ")
		}
		this.policy.options.args = args
	}
})