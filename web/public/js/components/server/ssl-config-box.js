Vue.component("ssl-config-box", {
	props: [
		"v-ssl-policy",
		"v-protocol",
		"v-server-id",
		"v-support-http3"
	],
	created: function () {
		let that = this
		setTimeout(function () {
			that.sortableCipherSuites()
		}, 100)
	},
	data: function () {
		let policy = this.vSslPolicy
		if (policy == null) {
			policy = {
				id: 0,
				isOn: true,
				certRefs: [],
				certs: [],
				clientCARefs: [],
				clientCACerts: [],
				clientAuthType: 0,
				minVersion: "TLS 1.1",
				hsts: null,
				cipherSuitesIsOn: false,
				cipherSuites: [],
				http2Enabled: true,
				http3Enabled: false,
				ocspIsOn: false
			}
		} else {
			if (policy.certRefs == null) {
				policy.certRefs = []
			}
			if (policy.certs == null) {
				policy.certs = []
			}
			if (policy.clientCARefs == null) {
				policy.clientCARefs = []
			}
			if (policy.clientCACerts == null) {
				policy.clientCACerts = []
			}
			if (policy.cipherSuites == null) {
				policy.cipherSuites = []
			}
		}

		let hsts = policy.hsts
		let hstsMaxAgeString = "31536000"
		if (hsts == null) {
			hsts = {
				isOn: false,
				maxAge: 31536000,
				includeSubDomains: false,
				preload: false,
				domains: []
			}
		}
		if (hsts.maxAge != null) {
			hstsMaxAgeString = hsts.maxAge.toString()
		}

		return {
			policy: policy,

			// hsts
			hsts: hsts,
			hstsOptionsVisible: false,
			hstsDomainAdding: false,
			hstsMaxAgeString: hstsMaxAgeString,
			addingHstsDomain: "",
			hstsDomainEditingIndex: -1,

			// 相关数据
			allVersions: window.SSL_ALL_VERSIONS,
			allCipherSuites: window.SSL_ALL_CIPHER_SUITES.$copy(),
			modernCipherSuites: window.SSL_MODERN_CIPHER_SUITES,
			intermediateCipherSuites: window.SSL_INTERMEDIATE_CIPHER_SUITES,
			allClientAuthTypes: window.SSL_ALL_CLIENT_AUTH_TYPES,
			cipherSuitesVisible: false,

			// 高级选项
			moreOptionsVisible: false
		}
	},
	watch: {
		hsts: {
			deep: true,
			handler: function () {
				this.policy.hsts = this.hsts
			}
		}
	},
	methods: {
		// 删除证书
		removeCert: function (index) {
			let that = this
			teaweb.confirm("确定删除此证书吗？证书数据仍然保留，只是当前网站不再使用此证书。", function () {
				that.policy.certRefs.$remove(index)
				that.policy.certs.$remove(index)
			})
		},

		// 选择证书
		selectCert: function () {
			let that = this
			let selectedCertIds = []
			if (this.policy != null && this.policy.certs.length > 0) {
				this.policy.certs.forEach(function (cert) {
					selectedCertIds.push(cert.id.toString())
				})
			}
			let serverId = this.vServerId
			if (serverId == null) {
				serverId = 0
			}
			teaweb.popup("/servers/certs/selectPopup?selectedCertIds=" + selectedCertIds + "&serverId=" + serverId, {
				width: "50em",
				height: "30em",
				callback: function (resp) {
					if (resp.data.cert != null && resp.data.certRef != null) {
						that.policy.certRefs.push(resp.data.certRef)
						that.policy.certs.push(resp.data.cert)
					}
					if (resp.data.certs != null && resp.data.certRefs != null) {
						that.policy.certRefs.$pushAll(resp.data.certRefs)
						that.policy.certs.$pushAll(resp.data.certs)
					}
					that.$forceUpdate()
				}
			})
		},

		// 上传证书
		uploadCert: function () {
			let that = this
			let serverId = this.vServerId
			if (typeof serverId != "number" && typeof serverId != "string") {
				serverId = 0
			}
			teaweb.popup("/servers/certs/uploadPopup?serverId=" + serverId, {
				height: "30em",
				callback: function (resp) {
					teaweb.success("上传成功", function () {
						that.policy.certRefs.push(resp.data.certRef)
						that.policy.certs.push(resp.data.cert)
					})
				}
			})
		},

		// 批量上传
		uploadBatch: function () {
			let that = this
			let serverId = this.vServerId
			if (typeof serverId != "number" && typeof serverId != "string") {
				serverId = 0
			}
			teaweb.popup("/servers/certs/uploadBatchPopup?serverId=" + serverId, {
				callback: function (resp) {
					if (resp.data.cert != null) {
						that.policy.certRefs.push(resp.data.certRef)
						that.policy.certs.push(resp.data.cert)
					}
					if (resp.data.certs != null) {
						that.policy.certRefs.$pushAll(resp.data.certRefs)
						that.policy.certs.$pushAll(resp.data.certs)
					}
					that.$forceUpdate()
				}
			})
		},

		// 申请证书
		requestCert: function () {
			// 已经在证书中的域名
			let excludeServerNames = []
			if (this.policy != null && this.policy.certs.length > 0) {
				this.policy.certs.forEach(function (cert) {
					excludeServerNames.$pushAll(cert.dnsNames)
				})
			}

			let that = this
			teaweb.popup("/servers/server/settings/https/requestCertPopup?serverId=" + this.vServerId + "&excludeServerNames=" + excludeServerNames.join(","), {
				callback: function () {
					that.policy.certRefs.push(resp.data.certRef)
					that.policy.certs.push(resp.data.cert)
				}
			})
		},

		// 更多选项
		changeOptionsVisible: function () {
			this.moreOptionsVisible = !this.moreOptionsVisible
		},

		// 格式化时间
		formatTime: function (timestamp) {
			return new Date(timestamp * 1000).format("Y-m-d")
		},

		// 格式化加密套件
		formatCipherSuite: function (cipherSuite) {
			return cipherSuite.replace(/(AES|3DES)/, "<var style=\"font-weight: bold\">$1</var>")
		},

		// 添加单个套件
		addCipherSuite: function (cipherSuite) {
			if (!this.policy.cipherSuites.$contains(cipherSuite)) {
				this.policy.cipherSuites.push(cipherSuite)
			}
			this.allCipherSuites.$removeValue(cipherSuite)
		},

		// 删除单个套件
		removeCipherSuite: function (cipherSuite) {
			let that = this
			teaweb.confirm("确定要删除此套件吗？", function () {
				that.policy.cipherSuites.$removeValue(cipherSuite)
				that.allCipherSuites = window.SSL_ALL_CIPHER_SUITES.$findAll(function (k, v) {
					return !that.policy.cipherSuites.$contains(v)
				})
			})
		},

		// 清除所选套件
		clearCipherSuites: function () {
			let that = this
			teaweb.confirm("确定要清除所有已选套件吗？", function () {
				that.policy.cipherSuites = []
				that.allCipherSuites = window.SSL_ALL_CIPHER_SUITES.$copy()
			})
		},

		// 批量添加套件
		addBatchCipherSuites: function (suites) {
			var that = this
			teaweb.confirm("确定要批量添加套件？", function () {
				suites.$each(function (k, v) {
					if (that.policy.cipherSuites.$contains(v)) {
						return
					}
					that.policy.cipherSuites.push(v)
				})
			})
		},

		/**
		 * 套件拖动排序
		 */
		sortableCipherSuites: function () {
			var box = document.querySelector(".cipher-suites-box")
			Sortable.create(box, {
				draggable: ".label",
				handle: ".icon.handle",
				onStart: function () {

				},
				onUpdate: function (event) {

				}
			})
		},

		// 显示所有套件
		showAllCipherSuites: function () {
			this.cipherSuitesVisible = !this.cipherSuitesVisible
		},

		// 显示HSTS更多选项
		showMoreHSTS: function () {
			this.hstsOptionsVisible = !this.hstsOptionsVisible;
			if (this.hstsOptionsVisible) {
				this.changeHSTSMaxAge()
			}
		},

		// 监控HSTS有效期修改
		changeHSTSMaxAge: function () {
			var v = parseInt(this.hstsMaxAgeString)
			if (isNaN(v) || v < 0) {
				this.hsts.maxAge = 0
				this.hsts.days = "-"
				return
			}
			this.hsts.maxAge = v
			this.hsts.days = v / 86400
			if (this.hsts.days == 0) {
				this.hsts.days = "-"
			}
		},

		// 设置HSTS有效期
		setHSTSMaxAge: function (maxAge) {
			this.hstsMaxAgeString = maxAge.toString()
			this.changeHSTSMaxAge()
		},

		// 添加HSTS域名
		addHstsDomain: function () {
			this.hstsDomainAdding = true
			this.hstsDomainEditingIndex = -1
			let that = this
			setTimeout(function () {
				that.$refs.addingHstsDomain.focus()
			}, 100)
		},

		// 修改HSTS域名
		editHstsDomain: function (index) {
			this.hstsDomainEditingIndex = index
			this.addingHstsDomain = this.hsts.domains[index]
			this.hstsDomainAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.addingHstsDomain.focus()
			}, 100)
		},

		// 确认HSTS域名添加
		confirmAddHstsDomain: function () {
			this.addingHstsDomain = this.addingHstsDomain.trim()
			if (this.addingHstsDomain.length == 0) {
				return;
			}
			if (this.hstsDomainEditingIndex > -1) {
				this.hsts.domains[this.hstsDomainEditingIndex] = this.addingHstsDomain
			} else {
				this.hsts.domains.push(this.addingHstsDomain)
			}
			this.cancelHstsDomainAdding()
		},

		// 取消HSTS域名添加
		cancelHstsDomainAdding: function () {
			this.hstsDomainAdding = false
			this.addingHstsDomain = ""
			this.hstsDomainEditingIndex = -1
		},

		// 删除HSTS域名
		removeHstsDomain: function (index) {
			this.cancelHstsDomainAdding()
			this.hsts.domains.$remove(index)
		},

		// 选择客户端CA证书
		selectClientCACert: function () {
			let that = this
			teaweb.popup("/servers/certs/selectPopup?isCA=1", {
				width: "50em",
				height: "30em",
				callback: function (resp) {
					if (resp.data.cert != null && resp.data.certRef != null) {
						that.policy.clientCARefs.push(resp.data.certRef)
						that.policy.clientCACerts.push(resp.data.cert)
					}
					if (resp.data.certs != null && resp.data.certRefs != null) {
						that.policy.clientCARefs.$pushAll(resp.data.certRefs)
						that.policy.clientCACerts.$pushAll(resp.data.certs)
					}
					that.$forceUpdate()
				}
			})
		},

		// 上传CA证书
		uploadClientCACert: function () {
			let that = this
			teaweb.popup("/servers/certs/uploadPopup?isCA=1", {
				height: "28em",
				callback: function (resp) {
					teaweb.success("上传成功", function () {
						that.policy.clientCARefs.push(resp.data.certRef)
						that.policy.clientCACerts.push(resp.data.cert)
					})
				}
			})
		},

		// 删除客户端CA证书
		removeClientCACert: function (index) {
			let that = this
			teaweb.confirm("确定删除此证书吗？证书数据仍然保留，只是当前网站不再使用此证书。", function () {
				that.policy.clientCARefs.$remove(index)
				that.policy.clientCACerts.$remove(index)
			})
		}
	},
	template: `<div>
	<h4>SSL/TLS相关配置</h4>
	<input type="hidden" name="sslPolicyJSON" :value="JSON.stringify(policy)"/>
	<table class="ui table definition selectable">
		<tbody>
			<tr v-show="vProtocol == 'https'">
				<td class="title">启用HTTP/2</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="policy.http2Enabled"/>
						<label></label>
					</div>
				</td>
			</tr>
			<tr v-show="vProtocol == 'https' && vSupportHttp3">
				<td class="title">启用HTTP/3</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="policy.http3Enabled"/>
						<label></label>
					</div>
				</td>
			</tr>
			<tr>
				<td class="title">设置证书</td>
				<td>
					<div v-if="policy.certs != null && policy.certs.length > 0">
						<div class="ui label small basic" v-for="(cert, index) in policy.certs" style="margin-top: 0.2em">
							{{cert.name}} / {{cert.dnsNames}} / 有效至{{formatTime(cert.timeEndAt)}} &nbsp; <a href="" title="删除" @click.prevent="removeCert(index)"><i class="icon remove"></i></a>
						</div>
						<div class="ui divider"></div>
					</div>
					<div v-else>
						<span class="red">选择或上传证书后<span v-if="vProtocol == 'https'">HTTPS</span><span v-if="vProtocol == 'tls'">TLS</span>服务才能生效。</span>
						<div class="ui divider"></div>
					</div>
					<button class="ui button tiny" type="button" @click.prevent="selectCert()">选择已有证书</button> &nbsp;
					<span class="disabled">|</span> &nbsp;
					<button class="ui button tiny" type="button" @click.prevent="uploadCert()">上传新证书</button> &nbsp;
					<button class="ui button tiny" type="button" @click.prevent="uploadBatch()">批量上传证书</button> &nbsp;
					<span class="disabled">|</span> &nbsp;
					<button class="ui button tiny" type="button" @click.prevent="requestCert()" v-if="vServerId != null && vServerId > 0">申请免费证书</button>
				</td>
			</tr>
			<tr>
				<td>TLS最低版本</td>
				<td>
					<select v-model="policy.minVersion" class="ui dropdown auto-width">
						<option v-for="version in allVersions" :value="version">{{version}}</option>
					</select>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeOptionsVisible"></more-options-tbody>
		<tbody v-show="moreOptionsVisible">
			<!-- 加密套件 -->
			<tr>
				<td>加密算法套件<em>（CipherSuites）</em></td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="policy.cipherSuitesIsOn" />
						<label>是否要自定义</label>
					</div>
					<div v-show="policy.cipherSuitesIsOn">
						<div class="ui divider"></div>
						<div class="cipher-suites-box">
							已添加套件({{policy.cipherSuites.length}})：
							<div v-for="cipherSuite in policy.cipherSuites" class="ui label tiny basic" style="margin-bottom: 0.5em">
								<input type="hidden" name="cipherSuites" :value="cipherSuite"/>
								<span v-html="formatCipherSuite(cipherSuite)"></span> &nbsp; <a href="" title="删除套件" @click.prevent="removeCipherSuite(cipherSuite)"><i class="icon remove"></i></a>
								<a href="" title="拖动改变顺序"><i class="icon bars handle"></i></a>
							</div>
						</div>
						<div>
							<div class="ui divider"></div>
							<span v-if="policy.cipherSuites.length > 0"><a href="" @click.prevent="clearCipherSuites()">[清除所有已选套件]</a> &nbsp; </span>
							<a href="" @click.prevent="addBatchCipherSuites(modernCipherSuites)">[添加推荐套件]</a> &nbsp;
							<a href="" @click.prevent="addBatchCipherSuites(intermediateCipherSuites)">[添加兼容套件]</a>
							<div class="ui divider"></div>
						</div>
		
						<div class="cipher-all-suites-box">
							<a href="" @click.prevent="showAllCipherSuites()"><span v-if="policy.cipherSuites.length == 0">所有</span>可选套件({{allCipherSuites.length}}) <i class="icon angle" :class="{down:!cipherSuitesVisible, up:cipherSuitesVisible}"></i></a>
							<a href="" v-if="cipherSuitesVisible" v-for="cipherSuite in allCipherSuites" class="ui label tiny" title="点击添加到自定义套件中" @click.prevent="addCipherSuite(cipherSuite)" v-html="formatCipherSuite(cipherSuite)" style="margin-bottom:0.5em"></a>
						</div>
						<p class="comment" v-if="cipherSuitesVisible">点击可选套件添加。</p>
					</div>
				</td>
			</tr>
			
			<!-- HSTS -->
			<tr v-show="vProtocol == 'https'">
				<td :class="{'color-border':hsts.isOn}">开启HSTS</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" name="hstsOn" v-model="hsts.isOn" value="1"/>
						<label></label>
					</div>
					<p class="comment">
						 开启后，会自动在响应Header中加入
						 <span class="ui label small">Strict-Transport-Security:
							 <var v-if="!hsts.isOn">...</var>
							 <var v-if="hsts.isOn"><span>max-age=</span>{{hsts.maxAge}}</var>
							 <var v-if="hsts.isOn && hsts.includeSubDomains">; includeSubDomains</var>
							 <var v-if="hsts.isOn && hsts.preload">; preload</var>
						 </span>
						  <span v-if="hsts.isOn">
							<a href="" @click.prevent="showMoreHSTS()">修改<i class="icon angle" :class="{down:!hstsOptionsVisible, up:hstsOptionsVisible}"></i> </a>
						 </span>
					</p>
				</td>
			</tr>
			<tr v-show="hsts.isOn && hstsOptionsVisible">
				<td class="color-border">HSTS有效时间<em>（max-age）</em></td>
				<td>
					<div class="ui fields inline">
						<div class="ui field">
							<input type="text" name="hstsMaxAge" v-model="hstsMaxAgeString" maxlength="10" size="10" @input="changeHSTSMaxAge()"/>
						</div>
						<div class="ui field">
							秒
						</div>
						<div class="ui field">{{hsts.days}}天</div>
					</div>
					<p class="comment">
						<a href="" @click.prevent="setHSTSMaxAge(31536000)" :class="{active:hsts.maxAge == 31536000}">[1年/365天]</a> &nbsp; &nbsp;
						<a href="" @click.prevent="setHSTSMaxAge(15768000)" :class="{active:hsts.maxAge == 15768000}">[6个月/182.5天]</a> &nbsp;  &nbsp;
						<a href="" @click.prevent="setHSTSMaxAge(2592000)"  :class="{active:hsts.maxAge == 2592000}">[1个月/30天]</a>
					</p>
				</td>
			</tr>
			<tr v-show="hsts.isOn && hstsOptionsVisible">
				<td class="color-border">HSTS包含子域名<em>（includeSubDomains）</em></td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" name="hstsIncludeSubDomains" value="1" v-model="hsts.includeSubDomains"/>
						<label></label>
					</div>
				</td>
			</tr>
			<tr v-show="hsts.isOn && hstsOptionsVisible">
				<td class="color-border">HSTS预加载<em>（preload）</em></td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" name="hstsPreload" value="1" v-model="hsts.preload"/>
						<label></label>
					</div>
				</td>
			</tr>
			<tr v-show="hsts.isOn && hstsOptionsVisible">
				<td class="color-border">HSTS生效的域名</td>
				<td colspan="2">
					<div class="names-box">
					<span class="ui label tiny basic" v-for="(domain, arrayIndex) in hsts.domains" :class="{blue:hstsDomainEditingIndex == arrayIndex}">{{domain}}
						<input type="hidden" name="hstsDomains" :value="domain"/> &nbsp;
						<a href="" @click.prevent="editHstsDomain(arrayIndex)" title="修改"><i class="icon pencil"></i></a>
						<a href="" @click.prevent="removeHstsDomain(arrayIndex)" title="删除"><i class="icon remove"></i></a>
					</span>
					</div>
					<div class="ui fields inline" v-if="hstsDomainAdding" style="margin-top:0.8em">
						<div class="ui field">
							<input type="text" name="addingHstsDomain" ref="addingHstsDomain" style="width:16em" maxlength="100" placeholder="域名，比如example.com" @keyup.enter="confirmAddHstsDomain()" @keypress.enter.prevent="1" v-model="addingHstsDomain" />
						</div>
						<div class="ui field">
							<button class="ui button tiny" type="button" @click="confirmAddHstsDomain()">确定</button>
							&nbsp; <a href="" @click.prevent="cancelHstsDomainAdding()">取消</a>
						</div>
					</div>
					<div class="ui field" style="margin-top: 1em">
						<button class="ui button tiny" type="button" @click="addHstsDomain()">+</button>
					</div>
					<p class="comment">如果没有设置域名的话，则默认支持所有的域名。</p>
				</td>
			</tr>
			
			<!-- OCSP -->
			<tr>
				<td>OCSP Stapling</td>
				<td><checkbox name="ocspIsOn" v-model="policy.ocspIsOn"></checkbox>
					<p class="comment">选中表示启用OCSP Stapling。</p>
				</td>
			</tr>
			
			<!-- 客户端认证 -->
			<tr>
				<td>客户端认证方式</td>
				<td>
					<select name="clientAuthType" v-model="policy.clientAuthType" class="ui dropdown auto-width">
						<option v-for="authType in allClientAuthTypes" :value="authType.type">{{authType.name}}</option>
					</select>
				</td>
			</tr>
			<tr>
				<td>客户端认证CA证书</td>
				<td>
					<div v-if="policy.clientCACerts != null && policy.clientCACerts.length > 0">
						<div class="ui label small basic" v-for="(cert, index) in policy.clientCACerts">
							{{cert.name}} / {{cert.dnsNames}} / 有效至{{formatTime(cert.timeEndAt)}} &nbsp; <a href="" title="删除" @click.prevent="removeClientCACert()"><i class="icon remove"></i></a>
						</div>
						<div class="ui divider"></div>
					</div>
					<button class="ui button tiny" type="button" @click.prevent="selectClientCACert()">选择已有证书</button> &nbsp;
					<button class="ui button tiny" type="button" @click.prevent="uploadClientCACert()">上传新证书</button>
					<p class="comment">用来校验客户端证书以增强安全性，通常不需要设置。</p>
				</td>
			</tr>
		</tbody>	
	</table>
	<div class="ui margin"></div>
</div>`
})