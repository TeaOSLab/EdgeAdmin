Vue.component("http-location-labels", {
	props: ["v-location-config"],
	data: function () {
		return {
			location: this.vLocationConfig
		}
	},
	methods: {
		// 判断是否已启用某配置
		configIsOn: function (config) {
			return config != null && config.isPrior && config.isOn
		},

		refIsOn: function (ref, config) {
			return this.configIsOn(ref) && config != null && config.isOn
		},

		len: function (arr) {
			return (arr == null) ? 0 : arr.length
		}
	},
	template: `	<div class="labels-box">
	<!-- TODO 思考是否给各个标签加上链接 -->
	
	<!-- 基本信息 -->
	<http-location-labels-label v-if="location.name != null && location.name.length > 0" :class="'olive'">{{location.name}}</http-location-labels-label>
	<http-location-labels-label v-if="location.isBreak">BREAK</http-location-labels-label>
	
	<!-- redirectToHTTPS -->
	<http-location-labels-label v-if="location.web != null && configIsOn(location.web.redirectToHTTPS)">自动跳转HTTPS</http-location-labels-label>
	
	<!-- Web -->
	<http-location-labels-label v-if="location.web != null && configIsOn(location.web.root)">文档根目录</http-location-labels-label>
	
	<!-- 反向代理 -->
	<http-location-labels-label v-if="refIsOn(location.reverseProxyRef, location.reverseProxy)">反向代理</http-location-labels-label>
	
	<!-- WAF -->
	<!-- TODO -->
	
	<!-- Cache -->
	<!-- TODO -->
	
	<!-- Charset -->
	<http-location-labels-label v-if="location.web != null && configIsOn(location.web.charset) && location.web.charset.charset.length > 0">{{location.web.charset.charset}}</http-location-labels-label>
	
	<!-- 访问日志 -->
	<!-- TODO -->
	
	<!-- 统计 -->
	<!-- TODO -->
	
	<!-- Gzip -->
	<http-location-labels-label v-if="location.web != null && refIsOn(location.web.gzipRef, location.web.gzip) && location.web.gzip.level > 0">Gzip:{{location.web.gzip.level}}</http-location-labels-label>
	
	<!-- HTTP Header -->
	<http-location-labels-label v-if="location.web != null && refIsOn(location.web.requestHeaderPolicyRef, location.web.requestHeaderPolicy) && (len(location.web.requestHeaderPolicy.addHeaders) > 0 || len(location.web.requestHeaderPolicy.setHeaders) > 0 || len(location.web.requestHeaderPolicy.replaceHeaders) > 0 || len(location.web.requestHeaderPolicy.deleteHeaders) > 0)">请求Header</http-location-labels-label>
	<http-location-labels-label v-if="location.web != null && refIsOn(location.web.responseHeaderPolicyRef, location.web.responseHeaderPolicy) && (len(location.web.responseHeaderPolicy.addHeaders) > 0 || len(location.web.responseHeaderPolicy.setHeaders) > 0 || len(location.web.responseHeaderPolicy.replaceHeaders) > 0 || len(location.web.responseHeaderPolicy.deleteHeaders) > 0)">响应Header</http-location-labels-label>
	
	<!-- Websocket -->
	<http-location-labels-label v-if="location.web != null && refIsOn(location.web.websocketRef, location.web.websocket)">Websocket</http-location-labels-label>
	
	<!-- 特殊页面 -->
	<div v-if="location.web != null && location.web.pages != null && location.web.pages.length > 0">
		<div v-for="page in location.web.pages" :key="page.id"><http-location-labels-label>PAGE [状态码{{page.status[0]}}] -&gt; {{page.url}}</http-location-labels-label></div>
	</div>
	<div v-if="location.web != null && configIsOn(location.web.shutdown)">
		<http-location-labels-label :v-class="'red'">临时关闭</http-location-labels-label>
	</div>
	
	<!-- 重写规则 -->
	<!-- TODO -->
</div>`
})

Vue.component("http-location-labels-label", {
	props: ["v-class"],
	template: `<span class="ui label tiny" :class="vClass" style="font-size:0.7em;padding:4px;margin-top:0.3em;margin-bottom:0.3em"><slot></slot></span>`
})