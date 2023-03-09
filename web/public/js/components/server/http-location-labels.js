Vue.component("http-location-labels", {
	props: ["v-location-config", "v-server-id"],
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
		},
		url: function (path) {
			return "/servers/server/settings/locations" + path + "?serverId=" + this.vServerId + "&locationId=" + this.location.id
		}
	},
	template: `	<div class="labels-box">
	<!-- 基本信息 -->
	<http-location-labels-label v-if="location.name != null && location.name.length > 0" :class="'olive'" :href="url('/location')">{{location.name}}</http-location-labels-label>
	
	<!-- domains -->
	<div v-if="location.domains != null && location.domains.length > 0">
		<grey-label v-for="domain in location.domains">{{domain}}</grey-label>
	</div>
	
	<!-- break -->
	<http-location-labels-label v-if="location.isBreak" :href="url('/location')">BREAK</http-location-labels-label>
	
	<!-- redirectToHTTPS -->
	<http-location-labels-label v-if="location.web != null && configIsOn(location.web.redirectToHTTPS)" :href="url('/http')">自动跳转HTTPS</http-location-labels-label>
	
	<!-- Web -->
	<http-location-labels-label v-if="location.web != null && configIsOn(location.web.root)" :href="url('/web')">文档根目录</http-location-labels-label>
	
	<!-- 反向代理 -->
	<http-location-labels-label v-if="refIsOn(location.reverseProxyRef, location.reverseProxy)" :v-href="url('/reverseProxy')">源站</http-location-labels-label>
	
	<!-- UAM -->
	<http-location-labels-label v-if="location.web != null && location.web.uam != null && location.web.uam.isPrior"><span :class="{disabled: !location.web.uam.isOn, red:location.web.uam.isOn}">5秒盾</span></http-location-labels-label>
	
	<!-- CC -->
	<http-location-labels-label v-if="location.web != null && location.web.cc != null && location.web.cc.isPrior"><span :class="{disabled: !location.web.cc.isOn, red:location.web.cc.isOn}">CC防护</span></http-location-labels-label>
	
	<!-- WAF -->
	<!-- TODO -->
	
	<!-- Cache -->
	<http-location-labels-label v-if="location.web != null && configIsOn(location.web.cache)" :v-href="url('/cache')">CACHE</http-location-labels-label>
	
	<!-- Charset -->
	<http-location-labels-label v-if="location.web != null && configIsOn(location.web.charset) && location.web.charset.charset.length > 0" :href="url('/charset')">{{location.web.charset.charset}}</http-location-labels-label>
	
	<!-- 访问日志 -->
	<!-- TODO -->
	
	<!-- 统计 -->
	<!-- TODO -->
	
	<!-- Gzip -->
	<http-location-labels-label v-if="location.web != null && refIsOn(location.web.gzipRef, location.web.gzip) && location.web.gzip.level > 0" :href="url('/gzip')">Gzip:{{location.web.gzip.level}}</http-location-labels-label>
	
	<!-- HTTP Header -->
	<http-location-labels-label v-if="location.web != null && refIsOn(location.web.requestHeaderPolicyRef, location.web.requestHeaderPolicy) && (len(location.web.requestHeaderPolicy.addHeaders) > 0 || len(location.web.requestHeaderPolicy.setHeaders) > 0 || len(location.web.requestHeaderPolicy.replaceHeaders) > 0 || len(location.web.requestHeaderPolicy.deleteHeaders) > 0)" :href="url('/headers')">请求Header</http-location-labels-label>
	<http-location-labels-label v-if="location.web != null && refIsOn(location.web.responseHeaderPolicyRef, location.web.responseHeaderPolicy) && (len(location.web.responseHeaderPolicy.addHeaders) > 0 || len(location.web.responseHeaderPolicy.setHeaders) > 0 || len(location.web.responseHeaderPolicy.replaceHeaders) > 0 || len(location.web.responseHeaderPolicy.deleteHeaders) > 0)" :href="url('/headers')">响应Header</http-location-labels-label>
	
	<!-- Websocket -->
	<http-location-labels-label v-if="location.web != null && refIsOn(location.web.websocketRef, location.web.websocket)" :href="url('/websocket')">Websocket</http-location-labels-label>
	
	<!-- 请求脚本 -->
	<http-location-labels-label v-if="location.web != null && location.web.requestScripts != null && ((location.web.requestScripts.initGroup != null && location.web.requestScripts.initGroup.isPrior) || (location.web.requestScripts.requestGroup != null && location.web.requestScripts.requestGroup.isPrior))" :href="url('/requestScripts')">请求脚本</http-location-labels-label>
	
	<!-- 自定义访客IP地址 -->
	<http-location-labels-label v-if="location.web != null && location.web.remoteAddr != null && location.web.remoteAddr.isPrior" :href="url('/remoteAddr')" :class="{disabled: !location.web.remoteAddr.isOn}">访客IP地址</http-location-labels-label>
	
	<!-- 请求限制 -->
	<http-location-labels-label v-if="location.web != null && location.web.requestLimit != null && location.web.requestLimit.isPrior" :href="url('/requestLimit')" :class="{disabled: !location.web.requestLimit.isOn}">请求限制</http-location-labels-label>
		
	<!-- 自定义页面 -->
	<div v-if="location.web != null && location.web.pages != null && location.web.pages.length > 0">
		<div v-for="page in location.web.pages" :key="page.id"><http-location-labels-label :href="url('/pages')">PAGE [状态码{{page.status[0]}}] -&gt; {{page.url}}</http-location-labels-label></div>
	</div>
	<div v-if="location.web != null && configIsOn(location.web.shutdown)">
		<http-location-labels-label :v-class="'red'" :href="url('/pages')">临时关闭</http-location-labels-label>
	</div>
	
	<!-- 重写规则 -->
	<div v-if="location.web != null && location.web.rewriteRules != null && location.web.rewriteRules.length > 0">
		<div v-for="rewriteRule in location.web.rewriteRules">
			<http-location-labels-label :href="url('/rewrite')">REWRITE {{rewriteRule.pattern}} -&gt; {{rewriteRule.replace}}</http-location-labels-label>
		</div>
	</div>
</div>`
})

Vue.component("http-location-labels-label", {
	props: ["v-class", "v-href"],
	template: `<a :href="vHref" class="ui label tiny basic" :class="vClass" style="font-size:0.7em;padding:4px;margin-top:0.3em;margin-bottom:0.3em"><slot></slot></a>`
})