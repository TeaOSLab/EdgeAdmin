// 访问日志搜索框
Vue.component("http-access-log-search-box", {
	props: ["v-ip", "v-domain", "v-keyword", "v-cluster-id", "v-node-id"],
	data: function () {
		let ip = this.vIp
		if (ip == null) {
			ip = ""
		}

		let domain = this.vDomain
		if (domain == null) {
			domain = ""
		}

		let keyword = this.vKeyword
		if (keyword == null) {
			keyword = ""
		}

		return {
			ip: ip,
			domain: domain,
			keyword: keyword,
			clusterId: this.vClusterId
		}
	},
	methods: {
		cleanIP: function () {
			this.ip = ""
			this.submit()
		},
		cleanDomain: function () {
			this.domain = ""
			this.submit()
		},
		cleanKeyword: function () {
			this.keyword = ""
			this.submit()
		},
		submit: function () {
			let parent = this.$el.parentNode
			while (true) {
				if (parent == null) {
					break
				}
				if (parent.tagName == "FORM") {
					break
				}
				parent = parent.parentNode
			}
			if (parent != null) {
				setTimeout(function () {
					parent.submit()
				}, 500)
			}
		},
		changeCluster: function (clusterId) {
			this.clusterId = clusterId
		}
	},
	template: `<div style="z-index: 10">
	<div class="margin"></div>
	<div class="ui fields inline">
		<div class="ui field">
			<div class="ui input left right labeled small">
				<span class="ui label basic" style="font-weight: normal">IP</span>
				<input type="text" name="ip" placeholder="x.x.x.x" size="15" v-model="ip"/>
				<a class="ui label basic" :class="{disabled: ip.length == 0}" @click.prevent="cleanIP"><i class="icon remove small"></i></a>
			</div>
		</div>
		<div class="ui field">
			<div class="ui input left right labeled small" >
				<span class="ui label basic" style="font-weight: normal">域名</span>
				<input type="text" name="domain" placeholder="example.com" size="15" v-model="domain"/>
				<a class="ui label basic" :class="{disabled: domain.length == 0}" @click.prevent="cleanDomain"><i class="icon remove small"></i></a>
			</div>
		</div>
		<div class="ui field">
			<div class="ui input left right labeled small">
				<span class="ui label basic" style="font-weight: normal">关键词</span>
				<input type="text" name="keyword" v-model="keyword" placeholder="路径、UserAgent、请求ID等..." size="30"/>
				<a class="ui label basic" :class="{disabled: keyword.length == 0}" @click.prevent="cleanKeyword"><i class="icon remove small"></i></a>
			</div>
		</div>
		<div class="ui field"><tip-icon content="一些特殊的关键词：<br/>单个状态码：status:200<br/>状态码范围：status:500-504<br/>查询IP：ip:192.168.1.100<br/>查询URL：https://goedge.cloud/docs<br/>查询路径部分：requestPath:/hello/world<br/>查询协议版本：proto:HTTP/1.1<br/>协议：scheme:http<br/>请求方法：method:POST<br/>请求来源：referer:example.com"></tip-icon></div>
	</div>
	<div class="ui fields inline" style="margin-top: 0.5em">
		<div class="ui field">
			<node-cluster-combo-box :v-cluster-id="clusterId" @change="changeCluster"></node-cluster-combo-box>
		</div>
		<div class="ui field" v-if="clusterId > 0">
			<node-combo-box :v-cluster-id="clusterId" :v-node-id="vNodeId"></node-combo-box>
		</div>
		<slot></slot>
		<div class="ui field">
			<button class="ui button small" type="submit">搜索日志</button>
		</div>
	</div>
</div>`
})