// Action列表
Vue.component("http-firewall-actions-view", {
	props: ["v-actions"],
	template: `<div>
		<div v-for="action in vActions" style="margin-bottom: 0.3em">
			<span :class="{red: action.category == 'block', orange: action.category == 'verify', green: action.category == 'allow'}">{{action.name}} ({{action.code.toUpperCase()}})
			  	<div v-if="action.options != null">
			  		<span class="grey small" v-if="action.code.toLowerCase() == 'page'">[{{action.options.status}}]</span>
			  		<span class="grey small" v-if="action.code.toLowerCase() == 'allow' && action.options != null && action.options.scope != null && action.options.scope.length > 0">
			  			<span v-if="action.options.scope == 'group'">[分组]</span>
						<span v-if="action.options.scope == 'server'">[网站]</span>
						<span v-if="action.options.scope == 'global'">[网站和策略]</span>	
					</span>
					<span class="grey small" v-if="action.code.toLowerCase() == 'record_ip'">
						<span v-if="action.options.type == 'black'" class="red">黑名单</span>
						<span v-if="action.options.type == 'white'" class="green">白名单</span>
						<span v-if="action.options.type == 'grey'" class="grey">灰名单</span>
					</span>
				</div>	
			</span>
		</div>             
</div>`
})