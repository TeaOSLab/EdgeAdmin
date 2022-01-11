// 访问日志搜索框
Vue.component("http-access-log-search-box", {
	props: ["v-ip", "v-domain", "v-keyword", "v-cluster-id", "v-node-id", "v-clusters"],
	mounted: function () {
		if (this.vClusterId >0) {
			this.changeCluster({
				value: this.vClusterId
			})
		}
	},
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
			nodes: []
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
		changeCluster: function (item) {
			this.nodes = []
			if (item != null) {
				let that = this
				Tea.action("/servers/logs/nodeOptions")
					.params({
						clusterId: item.value
					})
					.post()
					.success(function (resp) {
						that.nodes = resp.data.nodes
					})
			}
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
				<input type="text" name="domain" placeholder="xxx.com" size="15" v-model="domain"/>
				<a class="ui label basic" :class="{disabled: domain.length == 0}" @click.prevent="cleanDomain"><i class="icon remove small"></i></a>
			</div>
		</div>
		<div class="ui field">
			<div class="ui input left right labeled small">
				<span class="ui label basic" style="font-weight: normal">关键词</span>
				<input type="text" name="keyword" v-model="keyword" placeholder="路径、UserAgent等..." size="18"/>
				<a class="ui label basic" :class="{disabled: keyword.length == 0}" @click.prevent="cleanKeyword"><i class="icon remove small"></i></a>
			</div>
		</div>
		<slot></slot>
	</div>
	<div class="ui fields inline" style="margin-top: 0.5em">
		<div class="ui field" v-if="vClusters != null && vClusters.length > 0">
			<combo-box title="集群" name="clusterId" placeholder="集群名称" :v-items="vClusters" :v-value="vClusterId" @change="changeCluster"></combo-box>
		</div>
		<div class="ui field" v-if="nodes.length > 0">
			<combo-box title="节点" name="nodeId" placeholder="节点名称" :v-items="nodes" :v-value="vNodeId"></combo-box>
		</div>
		<div class="ui field">
			<button class="ui button small" type="submit">搜索日志</button>
		</div>
	</div>
</div>`
})