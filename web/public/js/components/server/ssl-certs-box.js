Vue.component("ssl-certs-box", {
	props: [
		"v-certs", // 证书列表
		"v-protocol", // 协议：https|tls
		"v-view-size", // 弹窗尺寸
		"v-single-mode" // 单证书模式
	],
	data: function () {
		let certs = this.vCerts
		if (certs == null) {
			certs = []
		}

		return {
			certs: certs
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
			teaweb.confirm("确定删除此证书吗？证书数据仍然保留，只是当前服务不再使用此证书。", function () {
				that.certs.$remove(index)
			})
		},

		// 选择证书
		selectCert: function () {
			let that = this
			let width = "50em"
			let height = "30em"
			let viewSize = this.vViewSize
			if (viewSize == null) {
				viewSize = "normal"
			}
			if (viewSize == "mini") {
				width = "35em"
				height = "20em"
			}
			teaweb.popup("/servers/certs/selectPopup?viewSize=" + viewSize, {
				width: width,
				height: height,
				callback: function (resp) {
					that.certs.push(resp.data.cert)
				}
			})
		},

		// 上传证书
		uploadCert: function () {
			let that = this
			teaweb.popup("/servers/certs/uploadPopup", {
				height: "28em",
				callback: function (resp) {
					teaweb.success("上传成功", function () {
						that.certs.push(resp.data.cert)
					})
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
		<div class="ui label small" v-for="(cert, index) in certs">
			{{cert.name}} / {{cert.dnsNames}} / 有效至{{formatTime(cert.timeEndAt)}} &nbsp; <a href="" title="删除" @click.prevent="removeCert(index)"><i class="icon remove"></i></a>
		</div>
		<div class="ui divider" v-if="buttonsVisible()"></div>
	</div>
	<div v-else>
		<span class="red">选择或上传证书后<span v-if="vProtocol == 'https'">HTTPS</span><span v-if="vProtocol == 'tls'">TLS</span>服务才能生效。</span>
		<div class="ui divider" v-if="buttonsVisible()"></div>
	</div>
	<div v-if="buttonsVisible()">
		<button class="ui button tiny" type="button" @click.prevent="selectCert()">选择已有证书</button> &nbsp;
		<button class="ui button tiny" type="button" @click.prevent="uploadCert()">上传新证书</button> &nbsp;
	</div>
</div>`
})