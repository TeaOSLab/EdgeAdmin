Tea.context(function () {
	this.success = NotifyPopup
	this.cacheRef = null

	let cachePolicyId = 0
	if (this.cachePolicies.length > 0) {
		cachePolicyId = this.cachePolicies[0].id
	}
	if (window.parent.UPDATING_CACHE_REF != null) {
		let cacheRef = window.parent.UPDATING_CACHE_REF
		this.cacheRef = cacheRef
		if (cacheRef.cachePolicy != null) {
			cachePolicyId = cacheRef.cachePolicy.id
		}
	}
	this.cachePolicyId = cachePolicyId
})