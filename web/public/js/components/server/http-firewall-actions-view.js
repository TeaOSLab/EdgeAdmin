// Action列表
Vue.component("http-firewall-actions-view", {
	props: ["v-actions"],
	template: `<div>
		<div v-for="action in vActions" style="margin-bottom: 0.3em">
			<span :class="{red: action.category == 'block', orange: action.category == 'verify', green: action.category == 'allow'}">{{action.name}} ({{action.code.toUpperCase()}})
			  	<div v-if="action.options != null">
			  		<span class="grey small" v-if="action.code.toLowerCase() == 'page'">[{{action.options.status}}]</span>
				</div>	
			</span>
		</div>             
</div>`
})