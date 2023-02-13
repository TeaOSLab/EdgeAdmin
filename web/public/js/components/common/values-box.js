Vue.component("values-box", {
	props: ["values", "v-values", "size", "maxlength", "name", "placeholder", "v-allow-empty", "validator"],
	data: function () {
		let values = this.values;
		if (values == null) {
			values = [];
		}

		if (this.vValues != null && typeof this.vValues == "object") {
			values = this.vValues
		}

		return {
			"realValues": values,
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
			this.value = this.realValues[index];
			var that = this;
			setTimeout(function () {
				that.$refs.value.focus();
			}, 200);
		},
		confirm: function () {
			if (this.value.length == 0) {
				if (typeof(this.vAllowEmpty) != "boolean" || !this.vAllowEmpty) {
					return
				}
			}

			// validate
			if (typeof(this.validator) == "function") {
				let resp = this.validator.call(this, this.value)
				if (typeof resp == "object") {
					if (typeof resp.isOk == "boolean" && !resp.isOk) {
						if (typeof resp.message == "string") {
							let that = this
							teaweb.warn(resp.message, function () {
								that.$refs.value.focus();
							})
						}
						return
					}
				}
			}

			if (this.isUpdating) {
				Vue.set(this.realValues, this.index, this.value);
			} else {
				this.realValues.push(this.value);
			}
			this.cancel()
			this.$emit("change", this.realValues)
		},
		remove: function (index) {
			this.realValues.$remove(index)
			this.$emit("change", this.realValues)
		},
		cancel: function () {
			this.isUpdating = false;
			this.isAdding = false;
			this.value = "";
		},
		updateAll: function (values) {
			this.realValues = values
		},
		addValue: function (v) {
			this.realValues.push(v)
		},

		startEditing: function () {
			this.isEditing = !this.isEditing
		},
		allValues: function () {
			return this.realValues
		}
	},
	template: `<div>
	<div v-show="!isEditing && realValues.length > 0">
		<div class="ui label tiny basic" v-for="(value, index) in realValues" style="margin-top:0.4em;margin-bottom:0.4em">
			<span v-if="value.toString().length > 0">{{value}}</span>
			<span v-if="value.toString().length == 0" class="disabled">[空]</span>
		</div>
		<a href="" @click.prevent="startEditing" style="font-size: 0.8em; margin-left: 0.2em">[修改]</a>
	</div>
	<div v-show="isEditing || realValues.length == 0">
		<div style="margin-bottom: 1em" v-if="realValues.length > 0">
			<div class="ui label tiny basic" v-for="(value, index) in realValues" style="margin-top:0.4em;margin-bottom:0.4em">
				<span v-if="value.toString().length > 0">{{value}}</span>
				<span v-if="value.toString().length == 0" class="disabled">[空]</span>
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