// 访问日志搜索框
Vue.component("http-access-log-search-box", {
	props: ["v-ip", "v-domain", "v-keyword"],
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
			keyword: keyword
		}
	},
	template: `<div>
	<div class="margin"></div>
	<div class="ui fields inline">
		<div class="ui field">
			<div class="ui input left labeled small">
				<span class="ui label basic" style="font-weight: normal">IP</span>
				<input type="text" name="ip" placeholder="x.x.x.x" size="15" v-model="ip"/>
			</div>
		</div>
		<div class="ui field">
			<div class="ui input left labeled small">
				<span class="ui label basic" style="font-weight: normal">域名</span>
				<input type="text" name="domain" placeholder="xxx.com" size="15" v-model="domain"/>
			</div>
		</div>
		<div class="ui field">
			<div class="ui input left labeled small">
				<span class="ui label basic" style="font-weight: normal">关键词</span>
				<input type="text" name="keyword" v-model="keyword" placeholder="路径、UserAgent等..." size="18"/>
			</div>
		</div>
		<slot></slot>
		<div class="ui field">
			<button class="ui button small" type="submit">查找</button>
		</div>
	</div>
</div>`
})