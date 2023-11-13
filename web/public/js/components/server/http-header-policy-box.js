Vue.component("http-header-policy-box", {
	props: ["v-request-header-policy", "v-request-header-ref", "v-response-header-policy", "v-response-header-ref", "v-params", "v-is-location", "v-is-group", "v-has-group-request-config", "v-has-group-response-config", "v-group-setting-url"],
	data: function () {
		let type = "response"
		let hash = window.location.hash
		if (hash == "#request") {
			type = "request"
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
		let requestNonStandardHeaders = []

		let requestPolicy = this.vRequestHeaderPolicy
		if (requestPolicy != null) {
			if (requestPolicy.setHeaders != null) {
				requestSettingHeaders = requestPolicy.setHeaders
			}
			if (requestPolicy.deleteHeaders != null) {
				requestDeletingHeaders = requestPolicy.deleteHeaders
			}
			if (requestPolicy.nonStandardHeaders != null) {
				requestNonStandardHeaders = requestPolicy.nonStandardHeaders
			}
		}

		// 响应相关
		let responseSettingHeaders = []
		let responseDeletingHeaders = []
		let responseNonStandardHeaders = []

		let responsePolicy = this.vResponseHeaderPolicy
		if (responsePolicy != null) {
			if (responsePolicy.setHeaders != null) {
				responseSettingHeaders = responsePolicy.setHeaders
			}
			if (responsePolicy.deleteHeaders != null) {
				responseDeletingHeaders = responsePolicy.deleteHeaders
			}
			if (responsePolicy.nonStandardHeaders != null) {
				responseNonStandardHeaders = responsePolicy.nonStandardHeaders
			}
		}

		let responseCORS = {
			isOn: false
		}
		if (responsePolicy.cors != null) {
			responseCORS = responsePolicy.cors
		}

		return {
			type: type,
			typeName: (type == "request") ? "请求" : "响应",

			requestHeaderRef: requestHeaderRef,
			responseHeaderRef: responseHeaderRef,
			requestSettingHeaders: requestSettingHeaders,
			requestDeletingHeaders: requestDeletingHeaders,
			requestNonStandardHeaders: requestNonStandardHeaders,

			responseSettingHeaders: responseSettingHeaders,
			responseDeletingHeaders: responseDeletingHeaders,
			responseNonStandardHeaders: responseNonStandardHeaders,
			responseCORS: responseCORS
		}
	},
	methods: {
		selectType: function (type) {
			this.type = type
			window.location.hash = "#" + type
			window.location.reload()
		},
		addSettingHeader: function (policyId) {
			teaweb.popup("/servers/server/settings/headers/createSetPopup?" + this.vParams + "&headerPolicyId=" + policyId + "&type=" + this.type, {
				height: "22em",
				callback: function () {
					teaweb.successRefresh("保存成功")
				}
			})
		},
		addDeletingHeader: function (policyId, type) {
			teaweb.popup("/servers/server/settings/headers/createDeletePopup?" + this.vParams + "&headerPolicyId=" + policyId + "&type=" + type, {
				callback: function () {
					teaweb.successRefresh("保存成功")
				}
			})
		},
		addNonStandardHeader: function (policyId, type) {
			teaweb.popup("/servers/server/settings/headers/createNonStandardPopup?" + this.vParams + "&headerPolicyId=" + policyId + "&type=" + type, {
				callback: function () {
					teaweb.successRefresh("保存成功")
				}
			})
		},
		updateSettingPopup: function (policyId, headerId) {
			teaweb.popup("/servers/server/settings/headers/updateSetPopup?" + this.vParams + "&headerPolicyId=" + policyId + "&headerId=" + headerId + "&type=" + this.type, {
				height: "22em",
				callback: function () {
					teaweb.successRefresh("保存成功")
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
		deleteNonStandardHeader: function (policyId, headerName) {
			teaweb.confirm("确定要删除'" + headerName + "'吗？", function () {
				Tea.action("/servers/server/settings/headers/deleteNonStandardHeader")
					.params({
						headerPolicyId: policyId,
						headerName: headerName
					})
					.post()
					.refresh()
			})
		},
		deleteHeader: function (policyId, type, headerId) {
			teaweb.confirm("确定要删除此报头吗？", function () {
					this.$post("/servers/server/settings/headers/delete")
						.params({
							headerPolicyId: policyId,
							type: type,
							headerId: headerId
						})
						.refresh()
				}
			)
		},
		updateCORS: function (policyId) {
			teaweb.popup("/servers/server/settings/headers/updateCORSPopup?" + this.vParams + "&headerPolicyId=" + policyId + "&type=" + this.type, {
				height: "30em",
				callback: function () {
					teaweb.successRefresh("保存成功")
				}
			})
		}
	},
	template: `<div>
	<div class="ui menu tabular small">
		<a class="item" :class="{active:type == 'response'}" @click.prevent="selectType('response')">响应报头<span v-if="responseSettingHeaders.length > 0">({{responseSettingHeaders.length}})</span></a>
		<a class="item" :class="{active:type == 'request'}" @click.prevent="selectType('request')">请求报头<span v-if="requestSettingHeaders.length > 0">({{requestSettingHeaders.length}})</span></a>
	</div>
	
	<div class="margin"></div>
	
	<input type="hidden" name="type" :value="type"/>
	
	<!-- 请求 -->
	<div v-if="(vIsLocation || vIsGroup) && type == 'request'">
		<input type="hidden" name="requestHeaderJSON" :value="JSON.stringify(requestHeaderRef)"/>
		<table class="ui table definition selectable">
			<prior-checkbox :v-config="requestHeaderRef"></prior-checkbox>
		</table>
		<submit-btn></submit-btn>
	</div>
	
	<div v-if="((!vIsLocation && !vIsGroup) || requestHeaderRef.isPrior) && type == 'request'">
		<div v-if="vHasGroupRequestConfig">
        	<div class="margin"></div>
        	<warning-message>由于已经在当前<a :href="vGroupSettingUrl + '#request'">网站分组</a>中进行了对应的配置，在这里的配置将不会生效。</warning-message>
    	</div>
    	<div :class="{'opacity-mask': vHasGroupRequestConfig}">
		<h4>设置请求报头 &nbsp; <a href="" @click.prevent="addSettingHeader(vRequestHeaderPolicy.id)" style="font-size: 0.8em">[添加新报头]</a></h4>
			<p class="comment" v-if="requestSettingHeaders.length == 0">暂时还没有自定义报头。</p>
			<table class="ui table selectable celled" v-if="requestSettingHeaders.length > 0">
				<thead>
					<tr>
						<th>名称</th>
						<th>值</th>
						<th class="two op">操作</th>
					</tr>
				</thead>
				<tbody v-for="header in requestSettingHeaders">
					<tr>
						<td class="five wide">
							<a href="" @click.prevent="updateSettingPopup(vRequestHeaderPolicy.id, header.id)">{{header.name}} <i class="icon expand small"></i></a>
							<div>
								<span v-if="header.status != null && header.status.codes != null && !header.status.always"><grey-label v-for="code in header.status.codes" :key="code">{{code}}</grey-label></span>
								<span v-if="header.methods != null && header.methods.length > 0"><grey-label v-for="method in header.methods" :key="method">{{method}}</grey-label></span>
								<span v-if="header.domains != null && header.domains.length > 0"><grey-label v-for="domain in header.domains" :key="domain">{{domain}}</grey-label></span>
								<grey-label v-if="header.shouldAppend">附加</grey-label>
								<grey-label v-if="header.disableRedirect">跳转禁用</grey-label>
								<grey-label v-if="header.shouldReplace && header.replaceValues != null && header.replaceValues.length > 0">替换</grey-label>
							</div>
						</td>
						<td>{{header.value}}</td>
						<td><a href="" @click.prevent="updateSettingPopup(vRequestHeaderPolicy.id, header.id)">修改</a> &nbsp; <a href="" @click.prevent="deleteHeader(vRequestHeaderPolicy.id, 'setHeader', header.id)">删除</a> </td>
					</tr>
				</tbody>
			</table>
			
			<h4>其他设置</h4>
			
			<table class="ui table definition selectable">
				<tbody>
					<tr>
						<td class="title">删除报头 <tip-icon content="可以通过此功能删除转发到源站的请求报文中不需要的报头"></tip-icon></td>
						<td>
							<div v-if="requestDeletingHeaders.length > 0">
								<div class="ui label small basic" v-for="headerName in requestDeletingHeaders">{{headerName}} <a href=""><i class="icon remove" title="删除" @click.prevent="deleteDeletingHeader(vRequestHeaderPolicy.id, headerName)"></i></a> </div>
								<div class="ui divider" ></div>
							</div>
							<button class="ui button small" type="button" @click.prevent="addDeletingHeader(vRequestHeaderPolicy.id, 'request')">+</button>
						</td>
					</tr>
					<tr>
						<td class="title">非标报头 <tip-icon content="可以通过此功能设置转发到源站的请求报文中非标准的报头，比如hello_world"></tip-icon></td>
						<td>
							<div v-if="requestNonStandardHeaders.length > 0">
								<div class="ui label small basic" v-for="headerName in requestNonStandardHeaders">{{headerName}} <a href=""><i class="icon remove" title="删除" @click.prevent="deleteNonStandardHeader(vRequestHeaderPolicy.id, headerName)"></i></a> </div>
								<div class="ui divider" ></div>
							</div>
							<button class="ui button small" type="button" @click.prevent="addNonStandardHeader(vRequestHeaderPolicy.id, 'request')">+</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>			
	</div>
	
	<!-- 响应 -->
	<div v-if="(vIsLocation || vIsGroup) && type == 'response'">
		<input type="hidden" name="responseHeaderJSON" :value="JSON.stringify(responseHeaderRef)"/>
		<table class="ui table definition selectable">
			<prior-checkbox :v-config="responseHeaderRef"></prior-checkbox>
		</table>
		<submit-btn></submit-btn>
	</div>
	
	<div v-if="((!vIsLocation && !vIsGroup) || responseHeaderRef.isPrior) && type == 'response'">
		<div v-if="vHasGroupResponseConfig">
        	<div class="margin"></div>
        	<warning-message>由于已经在当前<a :href="vGroupSettingUrl + '#response'">网站分组</a>中进行了对应的配置，在这里的配置将不会生效。</warning-message>
    	</div>
    	<div :class="{'opacity-mask': vHasGroupResponseConfig}">
			<h4>设置响应报头 &nbsp; <a href="" @click.prevent="addSettingHeader(vResponseHeaderPolicy.id)" style="font-size: 0.8em">[添加新报头]</a></h4>
			<p class="comment" style="margin-top: 0; padding-top: 0">将会覆盖已有的同名报头。</p>
			<p class="comment" v-if="responseSettingHeaders.length == 0">暂时还没有自定义报头。</p>
			<table class="ui table selectable celled" v-if="responseSettingHeaders.length > 0">
				<thead>
					<tr>
						<th>名称</th>
						<th>值</th>
						<th class="two op">操作</th>
					</tr>
				</thead>
				<tbody v-for="header in responseSettingHeaders">
					<tr>
						<td class="five wide">
							<a href="" @click.prevent="updateSettingPopup(vResponseHeaderPolicy.id, header.id)">{{header.name}} <i class="icon expand small"></i></a>
							<div>
								<span v-if="header.status != null && header.status.codes != null && !header.status.always"><grey-label v-for="code in header.status.codes" :key="code">{{code}}</grey-label></span>
								<span v-if="header.methods != null && header.methods.length > 0"><grey-label v-for="method in header.methods" :key="method">{{method}}</grey-label></span>
								<span v-if="header.domains != null && header.domains.length > 0"><grey-label v-for="domain in header.domains" :key="domain">{{domain}}</grey-label></span>
								<grey-label v-if="header.shouldAppend">附加</grey-label>
								<grey-label v-if="header.disableRedirect">跳转禁用</grey-label>
								<grey-label v-if="header.shouldReplace && header.replaceValues != null && header.replaceValues.length > 0">替换</grey-label>
							</div>
							
							<!-- CORS -->
							<div v-if="header.name == 'Access-Control-Allow-Origin' && header.value == '*'">
								<span class="red small">建议使用当前页面下方的"CORS自适应跨域"功能代替Access-Control-*-*相关报头。</span>
							</div>
						</td>
						<td>{{header.value}}</td>
						<td><a href="" @click.prevent="updateSettingPopup(vResponseHeaderPolicy.id, header.id)">修改</a> &nbsp; <a href="" @click.prevent="deleteHeader(vResponseHeaderPolicy.id, 'setHeader', header.id)">删除</a> </td>
					</tr>
				</tbody>
			</table>
			
			<h4>其他设置</h4>
			
			<table class="ui table definition selectable">
				<tbody>
					<tr>
						<td class="title">删除报头 <tip-icon content="可以通过此功能删除响应报文中不需要的报头"></tip-icon></td>
						<td>
							<div v-if="responseDeletingHeaders.length > 0">
								<div class="ui label small basic" v-for="headerName in responseDeletingHeaders">{{headerName}} &nbsp; <a href=""><i class="icon remove small" title="删除" @click.prevent="deleteDeletingHeader(vResponseHeaderPolicy.id, headerName)"></i></a></div>
								<div class="ui divider" ></div>
							</div>
							<button class="ui button small" type="button" @click.prevent="addDeletingHeader(vResponseHeaderPolicy.id, 'response')">+</button>
						</td>
					</tr>
					<tr>
						<td>非标报头 <tip-icon content="可以通过此功能设置响应报文中非标准的报头，比如hello_world"></tip-icon></td>
						<td>
							<div v-if="responseNonStandardHeaders.length > 0">
								<div class="ui label small basic" v-for="headerName in responseNonStandardHeaders">{{headerName}} &nbsp; <a href=""><i class="icon remove small" title="删除" @click.prevent="deleteNonStandardHeader(vResponseHeaderPolicy.id, headerName)"></i></a></div>
								<div class="ui divider" ></div>
							</div>
							<button class="ui button small" type="button" @click.prevent="addNonStandardHeader(vResponseHeaderPolicy.id, 'response')">+</button>
						</td>
					</tr>
					<tr>
						<td class="title">CORS自适应跨域</td>
						<td>
							<span v-if="responseCORS.isOn" class="green">已启用</span><span class="disabled" v-else="">未启用</span> &nbsp; <a href="" @click.prevent="updateCORS(vResponseHeaderPolicy.id)">[修改]</a>
							<p class="comment"><span v-if="!responseCORS.isOn">启用后，服务器可以</span><span v-else>服务器会</span>自动生成<code-label>Access-Control-*-*</code-label>相关的报头。</p>
						</td>
					</tr>
				</tbody>
			</table>
		</div>		
	</div>
	<div class="margin"></div>
</div>`
})