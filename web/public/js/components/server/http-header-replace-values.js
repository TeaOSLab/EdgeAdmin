Vue.component("http-header-replace-values", {
	props: ["v-replace-values"],
	data: function () {
		let values = this.vReplaceValues
		if (values == null) {
			values = []
		}
		return {
			values: values,
			isAdding: false,
			addingValue: {"pattern": "", "replacement": "", "isCaseInsensitive": false, "isRegexp": false}
		}
	},
	methods: {
		add: function () {
			this.isAdding = true
			let that = this
			setTimeout(function () {
				that.$refs.pattern.focus()
			})
		},
		remove: function (index) {
			this.values.$remove(index)
		},
		confirm: function () {
			let that = this
			if (this.addingValue.pattern.length == 0) {
				teaweb.warn("替换前内容不能为空", function () {
					that.$refs.pattern.focus()
				})
				return
			}

			this.values.push(this.addingValue)
			this.cancel()
		},
		cancel: function () {
			this.isAdding = false
			this.addingValue = {"pattern": "", "replacement": "", "isCaseInsensitive": false, "isRegexp": false}
		}
	},
	template: `<div>
	<input type="hidden" name="replaceValuesJSON" :value="JSON.stringify(values)"/>
	<div>
		<div v-for="(value, index) in values" class="ui label small" style="margin-bottom: 0.5em">
			<var>{{value.pattern}}</var><sup v-if="value.isCaseInsensitive" title="不区分大小写"><i class="icon info tiny"></i></sup> =&gt; <var v-if="value.replacement.length > 0">{{value.replacement}}</var><var v-else><span class="small grey">[空]</span></var>
			<a href="" @click.prevent="remove(index)" title="删除"><i class="icon remove small"></i></a>
		</div>
	</div>
	<div v-if="isAdding">
		<table class="ui table">
			<tr>
				<td class="title">替换前内容 *</td>
				<td><input type="text" v-model="addingValue.pattern" placeholder="替换前内容" ref="pattern" @keyup.enter="confirm()" @keypress.enter.prevent="1"/></td>
			</tr>	
			<tr>
				<td>替换后内容</td>
				<td><input type="text" v-model="addingValue.replacement" placeholder="替换后内容" @keyup.enter="confirm()" @keypress.enter.prevent="1"/></td>
			</tr>
			<tr>
				<td>是否忽略大小写</td>
				<td>
					<checkbox v-model="addingValue.isCaseInsensitive"></checkbox>
				</td>
			</tr>
		</table>

		<div>
			<button type="button" class="ui button tiny" @click.prevent="confirm">确定</button> &nbsp;
			<a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
		</div>
	</div>
	<div v-if="!isAdding">
		<button type="button" class="ui button tiny" @click.prevent="add">+</button>
	</div>
</div>`
})