Tea.context(function () {
	this.success = NotifyPopup
	this.cacheRef = null

	if (window.parent.UPDATING_CACHE_REF != null) {
		this.cacheRef = window.parent.UPDATING_CACHE_REF
		this.isReverse = this.cacheRef.isReverse
	}
})