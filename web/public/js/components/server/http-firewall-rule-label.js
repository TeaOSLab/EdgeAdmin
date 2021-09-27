// 显示WAF规则的标签
Vue.component("http-firewall-rule-label", {
	props: ["v-rule"],
	data: function () {
		return {
			rule: this.vRule
		}
	},
	methods: {
		showErr: function (err) {

			teaweb.popupTip("规则校验错误，请修正：<span class=\"red\">"  + teaweb.encodeHTML(err) + "</span>")
		},

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
		
		<a href="" v-if="rule.err != null && rule.err.length > 0" @click.prevent="showErr(rule.err)" style="color: #db2828; opacity: 1; border-bottom: 1px #db2828 dashed; margin-left: 0.5em">规则错误</a>
	</div>
</div>`
})