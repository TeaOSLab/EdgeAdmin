Vue.component("values-box", {
	props: ["values", "size", "maxlength", "name", "placeholder"],
	data: function () {
		let values = this.values;
		if (values == null) {
			values = [];
		}
		return {
			"vValues": values,
			"isUpdating": false,
			"isAdding": false,
			"index": 0,
			"value": "",
			isEditing: false
		}
	},
	methods: {
		create: function () {
			this.isAdding = true;
			var that = this;
			setTimeout(function () {
				that.$refs.value.focus();
			}, 200);
		},
		update: function (index) {
			this.cancel()
			this.isUpdating = true;
			this.index = index;
			this.value = this.vValues[index];
			var that = this;
			setTimeout(function () {
				that.$refs.value.focus();
			}, 200);
		},
		confirm: function () {
			if (this.value.length == 0) {
				return
			}

			if (this.isUpdating) {
				Vue.set(this.vValues, this.index, this.value);
			} else {
				this.vValues.push(this.value);
			}
			this.cancel()
			this.$emit("change", this.vValues)
		},
		remove: function (index) {
			this.vValues.$remove(index)
			this.$emit("change", this.vValues)
		},
		cancel: function () {
			this.isUpdating = false;
			this.isAdding = false;
			this.value = "";
		},
		updateAll: function (values) {
			this.vValeus = values
		},
		addValue: function (v) {
			this.vValues.push(v)
		},

		startEditing: function () {
			this.isEditing = !this.isEditing
		}
	},
	template: `<div>
	<div v-show="!isEditing && vValues.length > 0">
		<div class="ui label tiny basic" v-for="(value, index) in vValues" style="margin-top:0.4em;margin-bottom:0.4em">{{value}}</div>
		<a href="" @click.prevent="startEditing" style="font-size: 0.8em; margin-left: 0.2em">[修改]</a>
	</div>
	<div v-show="isEditing || vValues.length == 0">
		<div style="margin-bottom: 1em" v-if="vValues.length > 0">
			<div class="ui label tiny basic" v-for="(value, index) in vValues" style="margin-top:0.4em;margin-bottom:0.4em">{{value}}
				<input type="hidden" :name="name" :value="value"/>
				&nbsp; <a href="" @click.prevent="update(index)" title="修改"><i class="icon pencil small" ></i></a> 
				<a href="" @click.prevent="remove(index)" title="删除"><i class="icon remove"></i></a> 
			</div> 
			<div class="ui divider"></div>
		</div> 
		<!-- 添加|修改 -->
		<div v-if="isAdding || isUpdating">
			<div class="ui fields inline">
				<div class="ui field">
					<input type="text" :size="size" :maxlength="maxlength" :placeholder="placeholder" v-model="value" ref="value" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
				</div> 
				<div class="ui field">
					<button class="ui button small" type="button" @click.prevent="confirm()">确定</button> 
				</div>
				<div class="ui field">
					<a href="" @click.prevent="cancel()" title="取消"><i class="icon remove small"></i></a> 
				</div> 
			</div> 
		</div> 
		<div v-if="!isAdding && !isUpdating">
			<button class="ui button tiny" type="button" @click.prevent="create()">+</button> 
		</div>
	</div>	
</div>`
});