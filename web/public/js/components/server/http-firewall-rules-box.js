Vue.component("http-firewall-rules-box", {
	props: ["v-rules", "v-type"],
	data: function () {
		let rules = this.vRules
		if (rules == null) {
			rules = []
		}
		return {
			rules: rules
		}
	},
	methods: {
		addRule: function () {
			window.UPDATING_RULE = null
			let that = this
			teaweb.popup("/servers/components/waf/createRulePopup?type=" + this.vType, {
				height: "30em",
				callback: function (resp) {
					that.rules.push(resp.data.rule)
				}
			})
		},
		updateRule: function (index, rule) {
			window.UPDATING_RULE = teaweb.clone(rule)
			let that = this
			teaweb.popup("/servers/components/waf/createRulePopup?type=" + this.vType, {
				height: "30em",
				callback: function (resp) {
					Vue.set(that.rules, index, resp.data.rule)
				}
			})
		},
		removeRule: function (index) {
			let that = this
			teaweb.confirm("确定要删除此规则吗？", function () {
				that.rules.$remove(index)
			})
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
		isEmptyString: function (v) {
			return typeof v == "string" && v.length == 0
		}
	},
	template: `<div>
		<input type="hidden" name="rulesJSON" :value="JSON.stringify(rules)"/>
		<div v-if="rules.length > 0">
			<div v-for="(rule, index) in rules" class="ui label small basic" style="margin-bottom: 0.5em; line-height: 1.5">
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
					<span v-if="rule.paramFilters != null && rule.paramFilters.length > 0" v-for="paramFilter in rule.paramFilters"> | {{paramFilter.code}}</span> <span class="hover" :title="operatorDescription(rule.operator) + ((!rule.isComposed && rule.isCaseInsensitive) ? '\\n[大小写不敏感] ':'')">&lt;{{operatorName(rule.operator)}}&gt;</span> 
						<span v-if="!isEmptyString(rule.value)" class="hover">{{rule.value}}</span>
						<span v-else-if="operatorDataType(rule.operator) != 'none'" class="disabled" style="font-weight: normal" title="空字符串">[空]</span>
				</span>
				
				<!-- description -->
				<span v-if="rule.description != null && rule.description.length > 0" class="grey small">（{{rule.description}}）</span>
				
				<a href="" title="修改" @click.prevent="updateRule(index, rule)"><i class="icon pencil small"></i></a>
				<a href="" title="删除" @click.prevent="removeRule(index)"><i class="icon remove"></i></a>
			</div>
			<div class="ui divider"></div>
		</div>
		<button class="ui button tiny" type="button" @click.prevent="addRule()">+</button>
</div>`
})