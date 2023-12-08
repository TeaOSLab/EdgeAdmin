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
		calculateParamName: function (param) {
			let paramName = ""
			if (param != null) {
				window.WAF_RULE_CHECKPOINTS.forEach(function (checkpoint) {
					if (param == "${" + checkpoint.prefix + "}" || param.startsWith("${" + checkpoint.prefix + ".")) {
						paramName = checkpoint.name
					}
				})
			}
			return paramName
		},
		calculateParamDescription: function (param) {
			let paramName = ""
			let paramDescription = ""
			if (param != null) {
				window.WAF_RULE_CHECKPOINTS.forEach(function (checkpoint) {
					if (param == "${" + checkpoint.prefix + "}" || param.startsWith("${" + checkpoint.prefix + ".")) {
						paramName = checkpoint.name
						paramDescription = checkpoint.description
					}
				})
			}
			return paramName + ": " + paramDescription
		},
		operatorName: function (operatorCode) {
			let operatorName = operatorCode
			if (typeof (window.WAF_RULE_OPERATORS) != null) {
				window.WAF_RULE_OPERATORS.forEach(function (v) {
					if (v.code == operatorCode) {
						operatorName = v.name
					}
				})
			}

			return operatorName
		},
		operatorDescription: function (operatorCode) {
			let operatorName = operatorCode
			let operatorDescription = ""
			if (typeof (window.WAF_RULE_OPERATORS) != null) {
				window.WAF_RULE_OPERATORS.forEach(function (v) {
					if (v.code == operatorCode) {
						operatorName = v.name
						operatorDescription = v.description
					}
				})
			}

			return operatorName + ": " + operatorDescription
		},
		operatorDataType: function (operatorCode) {
			let operatorDataType = "none"
			if (typeof (window.WAF_RULE_OPERATORS) != null) {
				window.WAF_RULE_OPERATORS.forEach(function (v) {
					if (v.code == operatorCode) {
						operatorDataType = v.dataType
					}
				})
			}

			return operatorDataType
		},
		isEmptyString: function (v) {
			return typeof v == "string" && v.length == 0
		}
	},
	template: `<div>
	<div class="ui label small basic" style="line-height: 1.5">
		{{rule.name}} <span :title="calculateParamDescription(rule.param)" class="hover">{{calculateParamName(rule.param)}}<span class="small grey"> {{rule.param}}</span></span>

		<!-- cc2 -->
		<span v-if="rule.param == '\${cc2}'">
			{{rule.checkpointOptions.period}}秒内请求数
		</span>

		<!-- refererBlock -->
		<span v-if="rule.param == '\${refererBlock}'">
			<span v-if="rule.checkpointOptions.allowDomains != null && rule.checkpointOptions.allowDomains.length > 0">允许{{rule.checkpointOptions.allowDomains}}</span>
			<span v-if="rule.checkpointOptions.denyDomains != null && rule.checkpointOptions.denyDomains.length > 0">禁止{{rule.checkpointOptions.denyDomains}}</span>
		</span>

		<span v-else>
			<span v-if="rule.paramFilters != null && rule.paramFilters.length > 0" v-for="paramFilter in rule.paramFilters"> | {{paramFilter.code}}</span> 
		<span class="hover" :class="{dash:!rule.isComposed && rule.isCaseInsensitive}" :title="operatorDescription(rule.operator) + ((!rule.isComposed && rule.isCaseInsensitive) ? '\\n[大小写不敏感] ':'')">&lt;{{operatorName(rule.operator)}}&gt;</span> 
			<span class="hover" v-if="!isEmptyString(rule.value)">{{rule.value}}</span>
			<span v-else-if="operatorDataType(rule.operator) != 'none'" class="disabled" style="font-weight: normal" title="空字符串">[空]</span>
		</span>
		
		<!-- description -->
		<span v-if="rule.description != null && rule.description.length > 0" class="grey small">（{{rule.description}}）</span>
		
		<a href="" v-if="rule.err != null && rule.err.length > 0" @click.prevent="showErr(rule.err)" style="color: #db2828; opacity: 1; border-bottom: 1px #db2828 dashed; margin-left: 0.5em">规则错误</a>
	</div>
</div>`
})