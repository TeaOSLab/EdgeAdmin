Tea.context(function () {
	this.selectCert = function (cert) {
		NotifyPopup({
			code: 200,
			data: {
				cert: cert,
				certRef: {
					isOn: true,
					certId: cert.id
				}
			}
		})
	}

	this.encodeURL = function (arg) {
		return window.encodeURIComponent(arg)
	}

	/**
	 * 复选框
	 */
	this.countChecked = 0

	this.certs.forEach(function (cert) {
		cert.isChecked = false
	})

	this.changeAll = function (b) {
		let that = this
		this.certs.forEach(function (cert) {
			cert.isChecked = b
		})

		if (b) {
			let countChecked = 0
			this.certs.forEach(function (cert, index) {
				if (cert.isChecked && !that.certInfos[index].isSelected) {
					countChecked++
				}
			})
			this.countChecked = countChecked
		} else {
			this.countChecked = 0
		}
	}

	this.changeCertChecked = function () {
		let countChecked = 0
		this.certs.forEach(function (cert) {
			if (cert.isChecked) {
				countChecked++
			}
		})
		this.countChecked = countChecked
	}

	this.confirmChecked = function () {
		let resultCerts = []
		let resultCertRefs = []
		this.certs.forEach(function (cert) {
			if (cert.isChecked) {
				resultCerts.push(cert)
				resultCertRefs.push({
					isOn: true,
					certId: cert.id
				})
			}
		})
		NotifyPopup({
			code: 200,
			data: {
				certs: resultCerts,
				certRefs: resultCertRefs
			}
		})
	}

	this.searchNoneUserCerts = function () {
		this.$refs.userSelector.clear()
		this.$delay(function () {
			this.$refs.searchForm.submit()
		}, 10)
	}
})