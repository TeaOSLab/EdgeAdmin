Vue.component("http-websocket-box", {
	props: ["v-websocket-ref", "v-websocket-config", "v-is-location", "v-is-group"],
	data: function () {
		let websocketRef = this.vWebsocketRef
		if (websocketRef == null) {
			websocketRef = {
				isPrior: false,
				isOn: false,
				websocketId: 0
			}
		}

		let websocketConfig = this.vWebsocketConfig
		if (websocketConfig == null) {
			websocketConfig = {
				id: 0,
				isOn: false,
				handshakeTimeout: {
					count: 30,
					unit: "second"
				},
				allowAllOrigins: true,
				allowedOrigins: [],
				requestSameOrigin: true,
				requestOrigin: ""
			}
		} else {
			if (websocketConfig.handshakeTimeout == null) {
				websocketConfig.handshakeTimeout = {
					count: 30,
					unit: "second",
				}
			}
			if (websocketConfig.allowedOrigins == null) {
				websocketConfig.allowedOrigins = []
			}
		}

		return {
			websocketRef: websocketRef,
			websocketConfig: websocketConfig,
			handshakeTimeoutCountString: websocketConfig.handshakeTimeout.count.toString(),
			advancedVisible: false
		}
	},
	watch: {
		handshakeTimeoutCountString: function (v) {
			let count = parseInt(v)
			if (!isNaN(count) && count >= 0) {
				this.websocketConfig.handshakeTimeout.count = count
			} else {
				this.websocketConfig.handshakeTimeout.count = 0
			}
		}
	},
	methods: {
		isOn: function () {
			return ((!this.vIsLocation && !this.vIsGroup) || this.websocketRef.isPrior) && this.websocketRef.isOn
		},
		changeAdvancedVisible: function (v) {
			this.advancedVisible = v
		},
		createOrigin: function () {
			let that = this
			teaweb.popup("/servers/server/settings/websocket/createOrigin", {
				height: "12.5em",
				callback: function (resp) {
					that.websocketConfig.allowedOrigins.push(resp.data.origin)
				}
			})
		},
		removeOrigin: function (index) {
			this.websocketConfig.allowedOrigins.$remove(index)
		}
	},
	template: `<div>
	<input type="hidden" name="websocketRefJSON" :value="JSON.stringify(websocketRef)"/>
	<input type="hidden" name="websocketJSON" :value="JSON.stringify(websocketConfig)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="websocketRef" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="((!vIsLocation && !vIsGroup) || websocketRef.isPrior)">
			<tr>
				<td class="title">启用Websocket</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="websocketRef.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td class="color-border">允许所有来源域<em>（Origin）</em></td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="websocketConfig.allowAllOrigins"/>
						<label></label>
					</div>
					<p class="comment">选中表示允许所有的来源域。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn() && !websocketConfig.allowAllOrigins">
			<tr>
				<td class="color-border">允许的来源域列表<em>（Origin）</em></td>
				<td>
					<div v-if="websocketConfig.allowedOrigins.length > 0">
						<div class="ui label small basic" v-for="(origin, index) in websocketConfig.allowedOrigins">
							{{origin}} <a href="" title="删除" @click.prevent="removeOrigin(index)"><i class="icon remove small"></i></a>
						</div>
						<div class="ui divider"></div>
					</div>
					<button class="ui button tiny" type="button" @click.prevent="createOrigin()">+</button>
					<p class="comment">只允许在列表中的来源域名访问Websocket服务。</p>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-show="isOn()"></more-options-tbody>
		<tbody v-show="isOn() && advancedVisible">
			<tr>
				<td class="color-border">传递请求来源域</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="websocketConfig.requestSameOrigin"/>
						<label></label>
					</div>
					<p class="comment">选中后，表示把接收到的请求中的<code-label>Origin</code-label>字段传递到源站。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn() && advancedVisible && !websocketConfig.requestSameOrigin">
			<tr>
				<td class="color-border">指定传递的来源域</td>
				<td>
					<input type="text" v-model="websocketConfig.requestOrigin" maxlength="200"/>
					<p class="comment">指定向源站传递的<span class="ui label tiny">Origin</span>字段值。</p>
				</td>
			</tr>
		</tbody>
		<!-- TODO 这个选项暂时保留 -->
		<tbody v-show="isOn() && false">
			<tr>
				<td>握手超时时间<em>（Handshake）</em></td>
				<td>
					<div class="ui fields inline">
						<div class="ui field">
							<input type="text" maxlength="10" v-model="handshakeTimeoutCountString" style="width:6em"/>
						</div>
						<div class="ui field">
							秒
						</div>
					</div>
					<p class="comment">0表示使用默认的时间设置。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})