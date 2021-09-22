Vue.component("origin-list-box", {
	props: ["v-primary-origins", "v-backup-origins", "v-server-type", "v-params"],
	data: function () {
		return {
			primaryOrigins: this.vPrimaryOrigins,
			backupOrigins: this.vBackupOrigins
		}
	},
	methods: {
		createPrimaryOrigin: function () {
			teaweb.popup("/servers/server/settings/origins/addPopup?originType=primary&" + this.vParams, {
				height: "27em",
				callback: function (resp) {
					teaweb.success("保存成功", function () {
						window.location.reload()
					})
				}
			})
		},
		createBackupOrigin: function () {
			teaweb.popup("/servers/server/settings/origins/addPopup?originType=backup&" + this.vParams, {
				height: "27em",
				callback: function (resp) {
					teaweb.success("保存成功", function () {
						window.location.reload()
					})
				}
			})
		},
		updateOrigin: function (originId, originType) {
			teaweb.popup("/servers/server/settings/origins/updatePopup?originType=" + originType + "&" + this.vParams + "&originId=" + originId, {
				height: "27em",
				callback: function (resp) {
					teaweb.success("保存成功", function () {
						window.location.reload()
					})
				}
			})
		},
		deleteOrigin: function (originId, originType) {
			let that = this
			teaweb.confirm("确定要删除此源站吗？", function () {
				Tea.action("/servers/server/settings/origins/delete?" + that.vParams + "&originId=" + originId + "&originType=" + originType)
					.post()
					.success(function () {
						teaweb.success("保存成功", function () {
							window.location.reload()
						})
					})
			})
		}
	},
	template: `<div>
	<h3>主要源站 <a href="" @click.prevent="createPrimaryOrigin()">[添加主要源站]</a> </h3>
	<p class="comment" v-if="primaryOrigins.length == 0">暂时还没有主要源站。</p>
	<origin-list-table v-if="primaryOrigins.length > 0" :v-origins="vPrimaryOrigins" :v-origin-type="'primary'" @deleteOrigin="deleteOrigin" @updateOrigin="updateOrigin"></origin-list-table>

	<h3>备用源站 <a href="" @click.prevent="createBackupOrigin()">[添加备用源站]</a></h3>
	<p class="comment" v-if="backupOrigins.length == 0" :v-origins="primaryOrigins">暂时还没有备用源站。</p>
	<origin-list-table v-if="backupOrigins.length > 0" :v-origins="backupOrigins" :v-origin-type="'backup'" @deleteOrigin="deleteOrigin" @updateOrigin="updateOrigin"></origin-list-table>
</div>`
})

Vue.component("origin-list-table", {
	props: ["v-origins", "v-origin-type"],
	data: function () {
		return {}
	},
	methods: {
		deleteOrigin: function (originId) {
			this.$emit("deleteOrigin", originId, this.vOriginType)
		},
		updateOrigin: function (originId) {
			this.$emit("updateOrigin", originId, this.vOriginType)
		}
	},
	template: `
<table class="ui table selectable">
	<thead>
		<tr>
			<th>源站地址</th>
			<th>权重</th>
			<th class="width10">状态</th>
			<th class="two op">操作</th>
		</tr>	
	</thead>
	<tr v-for="origin in vOrigins">
		<td :class="{disabled:!origin.isOn}"><a href="" @click.prevent="updateOrigin(origin.id)">{{origin.addr}} &nbsp;<i class="icon clone outline small"></i></a>
			<div v-if="origin.name.length > 0" style="margin-top: 0.5em">
				<tiny-basic-label>{{origin.name}}</tiny-basic-label>
			</div>
			<div v-if="origin.domains != null && origin.domains.length > 0">
				<grey-label v-for="domain in origin.domains">{{domain}}</grey-label>
			</div>
		</td>
		<td :class="{disabled:!origin.isOn}">{{origin.weight}}</td>
		<td>
			<label-on :v-is-on="origin.isOn"></label-on>
		</td>
		<td>
			<a href="" @click.prevent="updateOrigin(origin.id)">修改</a> &nbsp;
			<a href="" @click.prevent="deleteOrigin(origin.id)">删除</a>
		</td>
	</tr>
</table>`
})