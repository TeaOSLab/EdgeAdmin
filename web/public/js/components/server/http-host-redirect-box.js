Vue.component("http-host-redirect-box", {
	props: ["v-redirects"],
	mounted: function () {
		let that = this
		sortTable(function (ids) {
			let newRedirects = []
			ids.forEach(function (id) {
				that.redirects.forEach(function (redirect) {
					if (redirect.id == id) {
						newRedirects.push(redirect)
					}
				})
			})
			that.updateRedirects(newRedirects)
		})
	},
	data: function () {
		let redirects = this.vRedirects
		if (redirects == null) {
			redirects = []
		}

		let id = 0
		redirects.forEach(function (v) {
			id++
			v.id = id
		})

		return {
			redirects: redirects,
			statusOptions: [
				{"code": 301, "text": "Moved Permanently"},
				{"code": 308, "text": "Permanent Redirect"},
				{"code": 302, "text": "Found"},
				{"code": 303, "text": "See Other"},
				{"code": 307, "text": "Temporary Redirect"}
			],
			id: id
		}
	},
	methods: {
		add: function () {
			let that = this
			window.UPDATING_REDIRECT = null

			teaweb.popup("/servers/server/settings/redirects/createPopup", {
				width: "50em",
				height: "36em",
				callback: function (resp) {
					that.id++
					resp.data.redirect.id = that.id
					that.redirects.push(resp.data.redirect)
					that.change()
				}
			})
		},
		update: function (index, redirect) {
			let that = this
			window.UPDATING_REDIRECT = redirect

			teaweb.popup("/servers/server/settings/redirects/createPopup", {
				width: "50em",
				height: "36em",
				callback: function (resp) {
					resp.data.redirect.id = redirect.id
					Vue.set(that.redirects, index, resp.data.redirect)
					that.change()
				}
			})
		},
		remove: function (index) {
			let that = this
			teaweb.confirm("确定要删除这条跳转规则吗？", function () {
				that.redirects.$remove(index)
				that.change()
			})
		},
		change: function () {
			let that = this
			setTimeout(function (){
				that.$emit("change", that.redirects)
			}, 100)
		},
		updateRedirects: function (newRedirects) {
			this.redirects = newRedirects
			this.change()
		}
	},
	template: `<div>
	<input type="hidden" name="hostRedirectsJSON" :value="JSON.stringify(redirects)"/>
	
	<first-menu>
		<menu-item @click.prevent="add">[创建]</menu-item>
	</first-menu>
	<div class="margin"></div>

	<p class="comment" v-if="redirects.length == 0">暂时还没有URL跳转规则。</p>
	<div v-show="redirects.length > 0">
		<table class="ui table celled selectable" id="sortable-table">
			<thead>
				<tr>
					<th style="width: 1em"></th>
					<th>跳转前</th>
					<th style="width: 1em"></th>
					<th>跳转后</th>
					<th>HTTP状态码</th>
					<th class="two wide">状态</th>
					<th class="two op">操作</th>
				</tr>
			</thead>
			<tbody v-for="(redirect, index) in redirects" :key="redirect.id" :v-id="redirect.id">
				<tr>
					<td style="text-align: center;"><i class="icon bars handle grey"></i> </td>
					<td>
						<div v-if="redirect.type == '' || redirect.type == 'url'">
							{{redirect.beforeURL}}
							<div style="margin-top: 0.4em">
								<grey-label><strong>URL跳转</strong></grey-label>
								<grey-label v-if="redirect.matchPrefix">匹配前缀</grey-label>
								<grey-label v-if="redirect.matchRegexp">正则匹配</grey-label>
								<grey-label v-if="!redirect.matchPrefix && !redirect.matchRegexp">精准匹配</grey-label>
								<grey-label v-if="redirect.exceptDomains != null && redirect.exceptDomains.length > 0" v-for="domain in redirect.exceptDomains">排除:{{domain}}</grey-label>
								<grey-label v-if="redirect.onlyDomains != null && redirect.onlyDomains.length > 0" v-for="domain in redirect.onlyDomains">仅限:{{domain}}</grey-label>
							</div>
						</div>
						<div v-if="redirect.type == 'domain'">
							<span v-if="redirect.domainsAll">所有域名</span>
							<span v-if="!redirect.domainsAll && redirect.domainsBefore != null">
								<span v-if="redirect.domainsBefore.length == 1">{{redirect.domainsBefore[0]}}</span>
								<span v-if="redirect.domainsBefore.length > 1">{{redirect.domainsBefore[0]}}等{{redirect.domainsBefore.length}}个域名</span>
							</span>
							<div style="margin-top: 0.4em">
								<grey-label><strong>域名跳转</strong></grey-label>
								<grey-label v-if="redirect.domainAfterScheme != null && redirect.domainAfterScheme.length > 0">{{redirect.domainAfterScheme}}</grey-label>
								<grey-label v-if="redirect.domainBeforeIgnorePorts">忽略端口</grey-label>
							</div>
						</div>
						<div v-if="redirect.type == 'port'">
							<span v-if="redirect.portsAll">所有端口</span>
							<span v-if="!redirect.portsAll && redirect.portsBefore != null">
								<span v-if="redirect.portsBefore.length <= 5">{{redirect.portsBefore.join(", ")}}</span>
								<span v-if="redirect.portsBefore.length > 5">{{redirect.portsBefore.slice(0, 5).join(", ")}}等{{redirect.portsBefore.length}}个端口</span>
							</span>
							<div style="margin-top: 0.4em">
								<grey-label><strong>端口跳转</strong></grey-label>
								<grey-label v-if="redirect.portAfterScheme != null && redirect.portAfterScheme.length > 0">{{redirect.portAfterScheme}}</grey-label>
							</div>
						</div>
						
						<div style="margin-top: 0.5em" v-if="redirect.conds != null && redirect.conds.groups != null && redirect.conds.groups.length > 0">
							<grey-label>匹配条件</grey-label>
						</div>
					</td>
					<td nowrap="">-&gt;</td>
					<td>
						<span v-if="redirect.type == '' || redirect.type == 'url'">{{redirect.afterURL}}</span>
						<span v-if="redirect.type == 'domain'">{{redirect.domainAfter}}</span>
						<span v-if="redirect.type == 'port'">{{redirect.portAfter}}</span>
					</td>
					<td>
						<span v-if="redirect.status > 0">{{redirect.status}}</span>
						<span v-else class="disabled">默认</span>
					</td>
					<td><label-on :v-is-on="redirect.isOn"></label-on></td>
					<td>
						<a href="" @click.prevent="update(index, redirect)">修改</a> &nbsp;
						<a href="" @click.prevent="remove(index)">删除</a>	
					</td>
				</tr>
			</tbody>
		</table>
		<p class="comment" v-if="redirects.length > 1">所有规则匹配顺序为从上到下，可以拖动左侧的<i class="icon bars"></i>排序。</p>
	</div>
	<div class="margin"></div>
</div>`
})