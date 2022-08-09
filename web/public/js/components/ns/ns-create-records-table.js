Vue.component("ns-create-records-table", {
	props: ["v-types"],
	data: function () {
		let types = this.vTypes
		if (types == null) {
			types = []
		}
		return {
			types: types,
			records: [
				{
					name: "",
					type: "A",
					value: "",
					routeCodes: [],
					ttl: 600,
					index: 0
				}
			],
			lastIndex: 0,
			isAddingRoutes: false // 是否正在添加线路
		}
	},
	methods: {
		add: function () {
			this.records.push({
				name: "",
				type: "A",
				value: "",
				routeCodes: [],
				ttl: 600,
				index: ++this.lastIndex
			})
			let that = this
			setTimeout(function () {
				that.$refs.nameInputs.$last().focus()
			}, 100)
		},
		remove: function (index) {
			this.records.$remove(index)
		},
		addRoutes: function () {
			this.isAddingRoutes = true
		},
		cancelRoutes: function () {
			let that = this
			setTimeout(function () {
				that.isAddingRoutes = false
			}, 1000)
		},
		changeRoutes: function (record, routes) {
			if (routes == null) {
				record.routeCodes = []
			} else {
				record.routeCodes = routes.map(function (route) {
					return route.code
				})
			}
		}
	},
	template: `<div>
<input type="hidden" name="recordsJSON" :value="JSON.stringify(records)"/>
<table class="ui table selectable celled" style="max-width: 60em">
	<thead class="full-width">
		<tr>
			<th style="width:10em">记录名</th>
			<th style="width:7em">记录类型</th>
			<th>线路</th>
			<th v-if="!isAddingRoutes">记录值</th>
			<th v-if="!isAddingRoutes">TTL</th>
			<th class="one op" v-if="!isAddingRoutes">操作</th>
		</tr>
	</thead>
	<tr v-for="(record, index) in records" :key="record.index">
		<td>
			<input type="text" style="width:10em" v-model="record.name" ref="nameInputs"/>		
		</td>
		<td>
			<select class="ui dropdown auto-width" v-model="record.type">
				<option v-for="type in types" :value="type.type">{{type.type}}</option>
			</select>
		</td>
		<td>
			<ns-routes-selector @add="addRoutes" @cancel="cancelRoutes" @change="changeRoutes(record, $event)"></ns-routes-selector>
		</td>
		<td v-if="!isAddingRoutes">
		  <input type="text" style="width:10em" maxlength="512" v-model="record.value"/>
		</td>
		<td v-if="!isAddingRoutes">
			<div class="ui input right labeled">
				<input type="text" v-model="record.ttl" style="width:5em" maxlength="8"/>
				<span class="ui label">秒</span>
			</div>
		</td>
		<td v-if="!isAddingRoutes">
			<a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove"></i></a>
		</td>
	</tr>
</table>
<button class="ui button tiny" type="button" @click.prevent="add">+</button>
</div>`,
})