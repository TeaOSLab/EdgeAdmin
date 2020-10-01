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
})