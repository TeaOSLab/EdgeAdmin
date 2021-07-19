// 显示WAF规则的标签
Vue.component("http-firewall-rule-label", {
	props: ["v-rule"],
	data: function () {
		return {
			rule: this.vRule
		}
	},
	template: `<div>
	<div class="ui label tiny basic">
		{{rule.name}}[{{rule.param}}] 

		<!-- cc2 -->
		<span v-if="rule.param == '\${cc2}'">
			{{rule.checkpointOptions.period}}秒/{{rule.checkpointOptions.threshold}}请求
		</span>
		<span v-else>
			<span v-if="rule.paramFilters != null && rule.paramFilters.length > 0" v-for="paramFilter in rule.paramFilters"> | {{paramFilter.code}}</span> 
		<var :class="{dash:rule.isCaseInsensitive}" :title="rule.isCaseInsensitive ? '大小写不敏感':''" v-if="!rule.isComposed">{{rule.operator}}</var> 
		{{rule.value}}
		</span>
	</div>
</div>`
})