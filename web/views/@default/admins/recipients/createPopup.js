Tea.context(function () {
	this.userDescription = ""

	this.changeInstance = function (instance) {
		if (instance != null) {
			this.userDescription = instance.media.userDescription
		} else {
			this.userDescription = ""
		}
	}

	/**
	 * 发送时间
	 */
	this.timeFromHour = ""
	this.timeFromMinute = ""
	this.timeFromSecond = ""
	this.timeToHour = ""
	this.timeToMinute = ""
	this.timeToSecond = ""

	this.clearTime = function () {
		this.timeFromHour = ""
		this.timeFromMinute = ""
		this.timeFromSecond = ""
		this.timeToHour = ""
		this.timeToMinute = ""
		this.timeToSecond = ""
	}
})