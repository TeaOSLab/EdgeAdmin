Vue.component("http-header-policy-box", {
	props: ["v-request-header-policy", "v-request-header-ref", "v-response-header-policy", "v-response-header-ref", "v-params", "v-is-location"],
	data: function () {
		let type = "request"
		let hash = window.location.hash
		if (hash == "#response") {
			type = "response"
		}

		// ref
		let requestHeaderRef = this.vRequestHeaderRef
		if (requestHeaderRef == null) {
			requestHeaderRef = {
				isPrior: false,
				isOn: true,
				headerPolicyId: 0
			}
		}

		let responseHeaderRef = this.vResponseHeaderRef
		if (responseHeaderRef == null) {
			responseHeaderRef = {
				isPrior: false,
				isOn: true,
				headerPolicyId: 0
			}
		}

		// 请求相关
		let requestSettingHeaders = []
		let requestDeletingHeaders = []

		let requestPolicy = this.vRequestHeaderPolicy
		if (requestPolicy != null) {
			if (requestPolicy.setHeaders != null) {
				requestSettingHeaders = requestPolicy.setHeaders
			}
			if (requestPolicy.deleteHeaders != null) {
				requestDeletingHeaders = requestPolicy.deleteHeaders
			}
		}

		// 响应相关
		let responseSettingHeaders = []
		let responseDeletingHeaders = []

		let responsePolicy = this.vResponseHeaderPolicy
		if (responsePolicy != null) {
			if (responsePolicy.setHeaders != null) {
				responseSettingHeaders = responsePolicy.setHeaders
			}
			if (responsePolicy.deleteHeaders != null) {
				responseDeletingHeaders = responsePolicy.deleteHeaders
			}
		}

		return {
			type: type,
			typeName: (type == "request") ? "请求" : "响应",
			requestHeaderRef: requestHeaderRef,
			responseHeaderRef: responseHeaderRef,
			requestSettingHeaders: requestSettingHeaders,
			requestDeletingHeaders: requestDeletingHeaders,
			responseSettingHeaders: responseSettingHeaders,
			responseDeletingHeaders: responseDeletingHeaders
		}
	},
	methods: {
		selectType: function (type) {
			this.type = type
			window.location.hash = "#" + type
			window.location.reload()
		},
		addSettingHeader: function (policyId) {
			teaweb.popup("/servers/server/settings/headers/createSetPopup?" + this.vParams + "&headerPolicyId=" + policyId, {
				callback: function () {
					window.location.reload()
				}
			})
		},
		addDeletingHeader: function (policyId, type) {
			teaweb.popup("/servers/server/settings/headers/createDeletePopup?" + this.vParams + "&headerPolicyId=" + policyId + "&type=" + type, {
				callback: function () {
					window.location.reload()
				}
			})
		},
		updateSettingPopup: function (policyId, headerId) {
			teaweb.popup("/servers/server/settings/headers/updateSetPopup?" + this.vParams + "&headerPolicyId=" + policyId + "&headerId=" + headerId, {
				callback: function () {
					window.location.reload()
				}
			})
		},
		deleteDeletingHeader: function (policyId, headerName) {
			teaweb.confirm("确定要删除'" + headerName + "'吗？", function () {
				Tea.action("/servers/server/settings/headers/deleteDeletingHeader")
					.params({
						headerPolicyId: policyId,
						headerName: headerName
					})
					.post()
					.refresh()
			})
		},
		deleteHeader: function (policyId, type, headerId) {
			teaweb.confirm("确定要删除此Header吗？", function () {
					this.$post("/servers/server/settings/headers/delete")
						.params({
							headerPolicyId: policyId,
							type: type,
							headerId: headerId
						})
						.refresh()
				}
			)
		}
	},
	template: `<div>
	<first-menu>
		<a class="item" :class="{active:type == 'request'}" @click.prevent="selectType('request')">请求Header</a>
		<a class="item" :class="{active:type == 'response'}" @click.prevent="selectType('response')">响应Header</a>
	</first-menu>
	
	<div class="margin"></div>
	
	<input type="hidden" name="type" :value="type"/>
	
	<!-- 请求 -->
	<div v-if="vIsLocation && type == 'request'">
		<input type="hidden" name="requestHeaderJSON" :value="JSON.stringify(requestHeaderRef)"/>
		<table class="ui table definition selectable">
			<prior-checkbox :v-config="requestHeaderRef"></prior-checkbox>
		</table>
		<submit-btn></submit-btn>
	</div>
	
	<div v-if="(!vIsLocation || requestHeaderRef.isPrior) && type == 'request'">
		<h3>设置请求Header <a href="" @click.prevent="addSettingHeader(vRequestHeaderPolicy.id)">[添加新Header]</a></h3>
		<p class="comment" v-if="requestSettingHeaders.length == 0">暂时还没有Header。</p>
		<table class="ui table selectable" v-if="requestSettingHeaders.length > 0">
			<thead>
				<tr>
					<th>名称</th>
					<th>值</th>
					<th class="two op">操作</th>
				</tr>
			</thead>
			<tr v-for="header in requestSettingHeaders">
				<td class="five wide">{{header.name}}</td>
				<td>{{header.value}}</td>
				<td><a href="" @click.prevent="updateSettingPopup(vRequestHeaderPolicy.id, header.id)">修改</a> &nbsp; <a href="" @click.prevent="deleteHeader(vRequestHeaderPolicy.id, 'setHeader', header.id)">删除</a> </td>
			</tr>
		</table>
		
		<h3>删除请求Header</h3>
		<p class="comment">这里可以设置需要从请求中删除的Header。</p>
		
		<table class="ui table definition selectable">
			<td class="title">需要删除的Header</td>
			<td>
				<div v-if="requestDeletingHeaders.length > 0">
					<div class="ui label small" v-for="headerName in requestDeletingHeaders">{{headerName}} <a href=""><i class="icon remove" title="删除" @click.prevent="deleteDeletingHeader(vRequestHeaderPolicy.id, headerName)"></i></a> </div>
					<div class="ui divider" ></div>
				</div>
				<button class="ui button small" type="button" @click.prevent="addDeletingHeader(vRequestHeaderPolicy.id, 'request')">+</button>
			</td>
		</table>
	</div>
	
	<!-- 响应 -->
	<div v-if="vIsLocation && type == 'response'">
		<input type="hidden" name="responseHeaderJSON" :value="JSON.stringify(responseHeaderRef)"/>
		<table class="ui table definition selectable">
			<prior-checkbox :v-config="responseHeaderRef"></prior-checkbox>
		</table>
		<submit-btn></submit-btn>
	</div>
	
	<div v-if="type == 'response'">
		<h3>设置响应Header <a href="" @click.prevent="addSettingHeader(vResponseHeaderPolicy.id)">[添加新Header]</a></h3>
		<p class="comment" v-if="responseSettingHeaders.length == 0">暂时还没有Header。</p>
		<table class="ui table selectable" v-if="responseSettingHeaders.length > 0">
			<thead>
				<tr>
					<th>名称</th>
					<th>值</th>
					<th class="two op">操作</th>
				</tr>
			</thead>
			<tr v-for="header in responseSettingHeaders">
				<td class="five wide">{{header.name}}</td>
				<td>{{header.value}}</td>
				<td><a href="" @click.prevent="updateSettingPopup(vResponseHeaderPolicy.id, header.id)">修改</a> &nbsp; <a href="" @click.prevent="deleteHeader(vResponseHeaderPolicy.id, 'setHeader', header.id)">删除</a> </td>
			</tr>
		</table>
		
		<h3>删除响应Header</h3>
		<p class="comment">这里可以设置需要从响应中删除的Header。</p>
		
		<table class="ui table definition selectable">
			<td class="title">需要删除的Header</td>
			<td>
				<div v-if="responseDeletingHeaders.length > 0">
					<div class="ui label small" v-for="headerName in responseDeletingHeaders">{{headerName}} <a href=""><i class="icon remove" title="删除" @click.prevent="deleteDeletingHeader(vResponseHeaderPolicy.id, headerName)"></i></a> </div>
					<div class="ui divider" ></div>
				</div>
				<button class="ui button small" type="button" @click.prevent="addDeletingHeader(vResponseHeaderPolicy.id, 'response')">+</button>
			</td>
		</table>
	</div>
	<div class="margin"></div>
</div>`
})