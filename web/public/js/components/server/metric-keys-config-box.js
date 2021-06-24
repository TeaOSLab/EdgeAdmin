// 指标对象
Vue.component("metric-keys-config-box", {
	props: ["v-keys"],
	data: function () {
		let keys = this.vKeys
		if (keys == null) {
			keys = []
		}
		return {
			keys: keys,
			isAdding: false
		}
	},
	methods: {},
	template: `<div>
	<div>
		<div v-for="key in keys">
			
		</div>
	</div>
	<div>
		<div class="ui fields inline">
			<div class="ui field">
				<input type="text" placeholder=""/>
			</div>
			<div class="ui field">
				<button type="button" class="ui button tiny">确定</button>
				<a href="" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
	</div>
	<div>
		<button type="button" class="ui button tiny">+</button>
	</div>
</div>`
})