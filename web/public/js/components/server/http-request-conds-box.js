Vue.component("http-request-conds-box", {
	props: ["v-conds"],
	data: function () {
		let conds = this.vConds
		if (conds == null) {
			conds = {
				isOn: true,
				connector: "or",
				groups: []
			}
		}
		if (conds.groups == null) {
			conds.groups = []
		}
		return {
			conds: conds,
			components: window.REQUEST_COND_COMPONENTS
		}
	},
	methods: {
		change: function () {
			this.$emit("change", this.conds)
		},
		addGroup: function () {
			window.UPDATING_COND_GROUP = null

			let that = this
			teaweb.popup("/servers/server/settings/conds/addGroupPopup", {
				height: "30em",
				callback: function (resp) {
					that.conds.groups.push(resp.data.group)
					that.change()
				}
			})
		},
		updateGroup: function (groupIndex, group) {
			window.UPDATING_COND_GROUP = group
			let that = this
			teaweb.popup("/servers/server/settings/conds/addGroupPopup", {
				height: "30em",
				callback: function (resp) {
					Vue.set(that.conds.groups, groupIndex, resp.data.group)
					that.change()
				}
			})
		},
		removeGroup: function (groupIndex) {
			let that = this
			teaweb.confirm("确定要删除这一组条件吗？", function () {
				that.conds.groups.$remove(groupIndex)
				that.change()
			})
		},
		typeName: function (cond) {
			let c = this.components.$find(function (k, v) {
				return v.type == cond.type
			})
			if (c != null) {
				return c.name;
			}
			return cond.param + " " + cond.operator
		}
	},
	template: `<div>
		<input type="hidden" name="condsJSON" :value="JSON.stringify(conds)"/>
		<div v-if="conds.groups.length > 0">
			<table class="ui table">
				<tr v-for="(group, groupIndex) in conds.groups">
					<td class="title" :class="{'color-border':conds.connector == 'and'}" :style="{'border-bottom':(groupIndex < conds.groups.length-1) ? '1px solid rgba(34,36,38,.15)':''}">分组{{groupIndex+1}}</td>
					<td style="background: white; word-break: break-all" :style="{'border-bottom':(groupIndex < conds.groups.length-1) ? '1px solid rgba(34,36,38,.15)':''}">
						<var v-for="(cond, index) in group.conds" style="font-style: normal;display: inline-block; margin-bottom:0.5em">
							<span class="ui label tiny">
								<var v-if="cond.type.length == 0 || cond.type == 'params'" style="font-style: normal">{{cond.param}} <var>{{cond.operator}}</var></var>
								<var v-if="cond.type.length > 0 && cond.type != 'params'" style="font-style: normal">{{typeName(cond)}}: </var>
								{{cond.value}}
								<sup v-if="cond.isCaseInsensitive" title="不区分大小写"><i class="icon info small"></i></sup>
							</span>
							
							<var v-if="index < group.conds.length - 1"> {{group.connector}} &nbsp;</var>
						</var>
					</td>
					<td style="width: 5em; background: white" :style="{'border-bottom':(groupIndex < conds.groups.length-1) ? '1px solid rgba(34,36,38,.15)':''}">
						<a href="" title="修改分组" @click.prevent="updateGroup(groupIndex, group)"><i class="icon pencil small"></i></a> <a href="" title="删除分组" @click.prevent="removeGroup(groupIndex)"><i class="icon remove"></i></a>
					</td>
				</tr>
			</table>
			<div class="ui divider"></div>
		</div>
		
		<!-- 分组之间关系 -->
		<table class="ui table" v-if="conds.groups.length > 1">
			<tr>
				<td class="title">分组之间关系</td>
				<td>
					<select class="ui dropdown auto-width" v-model="conds.connector">
						<option value="and">和</option>
						<option value="or">或</option>
					</select>
					<p class="comment">
						<span v-if="conds.connector == 'or'">只要满足其中一个条件分组即可。</span>
						<span v-if="conds.connector == 'and'">需要满足所有条件分组。</span>
					</p>	
				</td>
			</tr>
		</table>
		
		<div>
			<button class="ui button tiny basic" type="button" @click.prevent="addGroup()">+添加分组</button>
		</div>
	</div>	
</div>`
})