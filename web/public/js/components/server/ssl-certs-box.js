Vue.component("ssl-certs-box", {
	props: [
		"v-certs", // 证书列表
		"v-cert", // 单个证书
		"v-protocol", // 协议：https|tls
		"v-view-size", // 弹窗尺寸：normal, mini
		"v-single-mode", // 单证书模式
		"v-description", // 描述文字
		"v-domains", // 搜索的域名列表或者函数
		"v-user-id" // 用户ID
	],
	data: function () {
		let certs = this.vCerts
		if (certs == null) {
			certs = []
		}
		if (this.vCert != null) {
			certs.push(this.vCert)
		}

		let description = this.vDescription
		if (description == null || typeof (description) != "string") {
			description = ""
		}

		return {
			certs: certs,
			description: description
		}
	},
	methods: {
		certIds: function () {
			return this.certs.map(function (v) {
				return v.id
			})
		},
		// 删除证书
		removeCert: function (index) {
			let that = this
			teaweb.confirm("确定删除此证书吗？证书数据仍然保留，只是当前网站不再使用此证书。", function () {
				that.certs.$remove(index)
			})
		},

		// 选择证书
		selectCert: function () {
			let that = this
			let width = "54em"
			let height = "32em"
			let viewSize = this.vViewSize
			if (viewSize == null) {
				viewSize = "normal"
			}
			if (viewSize == "mini") {
				width = "35em"
				height = "20em"
			}

			let searchingDomains = []
			if (this.vDomains != null) {
				if (typeof this.vDomains == "function") {
					let resultDomains = this.vDomains()
					if (resultDomains != null && typeof resultDomains == "object" && (resultDomains instanceof Array)) {
						searchingDomains = resultDomains
					}
				} else if (typeof this.vDomains == "object" && (this.vDomains instanceof Array)) {
					searchingDomains = this.vDomains
				}
				if (searchingDomains.length > 10000) {
					searchingDomains = searchingDomains.slice(0, 10000)
				}
			}

			let selectedCertIds = this.certs.map(function (cert) {
				return cert.id
			})
			let userId = this.vUserId
			if (userId == null) {
				userId = 0
			}

			teaweb.popup("/servers/certs/selectPopup?viewSize=" + viewSize + "&searchingDomains=" + window.encodeURIComponent(searchingDomains.join(",")) + "&selectedCertIds=" + selectedCertIds.join(",") + "&userId=" + userId, {
				width: width,
				height: height,
				callback: function (resp) {
					if (resp.data.cert != null) {
						that.certs.push(resp.data.cert)
					}
					if (resp.data.certs != null) {
						that.certs.$pushAll(resp.data.certs)
					}
					that.$forceUpdate()
				}
			})
		},

		// 上传证书
		uploadCert: function () {
			let that = this
			let userId = this.vUserId
			if (typeof userId != "number" && typeof userId != "string") {
				userId = 0
			}
			teaweb.popup("/servers/certs/uploadPopup?userId=" + userId, {
				height: "28em",
				callback: function (resp) {
					teaweb.success("上传成功", function () {
						if (resp.data.cert != null) {
							that.certs.push(resp.data.cert)
						}
						if (resp.data.certs != null) {
							that.certs.$pushAll(resp.data.certs)
						}
						that.$forceUpdate()
					})
				}
			})
		},

		// 批量上传
		uploadBatch: function () {
			let that = this
			let userId = this.vUserId
			if (typeof userId != "number" && typeof userId != "string") {
				userId = 0
			}
			teaweb.popup("/servers/certs/uploadBatchPopup?userId=" + userId, {
				callback: function (resp) {
					if (resp.data.cert != null) {
						that.certs.push(resp.data.cert)
					}
					if (resp.data.certs != null) {
						that.certs.$pushAll(resp.data.certs)
					}
					that.$forceUpdate()
				}
			})
		},

		// 格式化时间
		formatTime: function (timestamp) {
			return new Date(timestamp * 1000).format("Y-m-d")
		},

		// 判断是否显示选择｜上传按钮
		buttonsVisible: function () {
			return this.vSingleMode == null || !this.vSingleMode || this.certs == null || this.certs.length == 0
		}
	},
	template: `<div>
	<input type="hidden" name="certIdsJSON" :value="JSON.stringify(certIds())"/>
	<div v-if="certs != null && certs.length > 0">
		<div class="ui label small basic" v-for="(cert, index) in certs">
			{{cert.name}} / {{cert.dnsNames}} / 有效至{{formatTime(cert.timeEndAt)}} &nbsp; <a href="" title="删除" @click.prevent="removeCert(index)"><i class="icon remove"></i></a>
		</div>
		<div class="ui divider" v-if="buttonsVisible()"></div>
	</div>
	<div v-else>
		<span class="red" v-if="description.length == 0">选择或上传证书后<span v-if="vProtocol == 'https'">HTTPS</span><span v-if="vProtocol == 'tls'">TLS</span>服务才能生效。</span>
		<span class="grey" v-if="description.length > 0">{{description}}</span>
		<div class="ui divider" v-if="buttonsVisible()"></div>
	</div>
	<div v-if="buttonsVisible()">
		<button class="ui button tiny" type="button" @click.prevent="selectCert()">选择已有证书</button> &nbsp;
		<span class="disabled">|</span> &nbsp;
		<button class="ui button tiny" type="button" @click.prevent="uploadCert()">上传新证书</button> &nbsp;
		<button class="ui button tiny" type="button" @click.prevent="uploadBatch()">批量上传证书</button> &nbsp;
	</div>
</div>`
})