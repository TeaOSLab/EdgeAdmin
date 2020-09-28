Vue.component("http-rewrite-rule-list", {
	props: ["v-web-id", "v-rewrite-rules"],
	mounted: function () {
		setTimeout(this.sort, 1000)
	},
	data: function () {
		let rewriteRules = this.vRewriteRules
		if (rewriteRules == null) {
			rewriteRules = []
		}
		return {
			rewriteRules: rewriteRules
		}
	},
	methods: {
		updateRewriteRule: function (rewriteRuleId) {
			teaweb.popup("/servers/server/settings/rewrite/updatePopup?webId=" + this.vWebId + "&rewriteRuleId=" + rewriteRuleId, {
				height: "26em",
				callback: function () {
					window.location.reload()
				}
			})
		},
		deleteRewriteRule: function (rewriteRuleId) {
			let that = this
			teaweb.confirm("确定要删除此重写规则吗？", function () {
				Tea.action("/servers/server/settings/rewrite/delete")
					.params({
						webId: that.vWebId,
						rewriteRuleId: rewriteRuleId
					})
					.post()
					.refresh()
			})
		},
		// 排序
		sort: function () {
			if (this.rewriteRules.length == 0) {
				return
			}
			let that = this
			sortTable(function (rowIds) {
				Tea.action("/servers/server/settings/rewrite/sort")
					.post()
					.params({
						webId: that.vWebId,
						rewriteRuleIds: rowIds
					})
					.success(function () {
						teaweb.success("保存成功")
					})
			})
		}
	},
	template: `<div>
	<div class="margin"></div>
	<p class="comment" v-if="rewriteRules.length == 0">暂时还没有重写规则。</p>
	<table class="ui table selectable" v-if="rewriteRules.length > 0" id="sortable-table">
		<thead>
			<tr>
				<th style="width:1em"></th>
				<th>匹配规则</th>
				<th>转发目标</th>
				<th>转发方式</th>
				<th class="two wide">状态</th>
				<th class="two op">操作</th>
			</tr>
		</thead>
		<tbody v-for="rule in rewriteRules" :v-id="rule.id">
			<tr>
				<td><i class="icon bars grey handle"></i></td>
				<td>{{rule.pattern}}
				<br/>
					<http-rewrite-labels-label class="ui label tiny" v-if="rule.isBreak">BREAK</http-rewrite-labels-label>
					<http-rewrite-labels-label class="ui label tiny" v-if="rule.mode == 'redirect' && rule.redirectStatus != 307">{{rule.redirectStatus}}</http-rewrite-labels-label>
					<http-rewrite-labels-label class="ui label tiny" v-if="rule.proxyHost.length > 0">Host: {{rule.proxyHost}}</http-rewrite-labels-label>
				</td>
				<td>{{rule.replace}}</td>
				<td>
					<span v-if="rule.mode == 'proxy'">隐式</span>
					<span v-if="rule.mode == 'redirect'">显示</span>
				</td>
				<td>
					<label-on :v-is-on="rule.isOn"></label-on>
				</td>
				<td>
					<a href="" @click.prevent="updateRewriteRule(rule.id)">修改</a> &nbsp;
					<a href="" @click.prevent="deleteRewriteRule(rule.id)">删除</a>
				</td>
			</tr>
		</tbody>
	</table>
	<p class="comment" v-if="rewriteRules.length > 0">拖动左侧的<i class="icon bars grey"></i>图标可以对重写规则进行排序。</p>

</div>`
})

Vue.component("http-rewrite-labels-label", {
	props: ["v-class"],
	template: `<span class="ui label tiny" :class="vClass" style="font-size:0.7em;padding:4px;margin-top:0.3em;margin-bottom:0.3em"><slot></slot></span>`
})