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
				width: "45em",
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
				width: "45em",
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
				width: "45em",
				height: "27em",
				callback: function (resp) {
					teaweb.success("保存成功", function () {
						window.location.reload()
					})
				}
			})
		},
		deleteOrigin: function (originId, originAddr, originType) {
			let that = this
			teaweb.confirm("确定要删除此源站（" + originAddr + "）吗？", function () {
				Tea.action("/servers/server/settings/origins/delete?" + that.vParams + "&originId=" + originId + "&originType=" + originType)
					.post()
					.success(function () {
						teaweb.success("删除成功", function () {
							window.location.reload()
						})
					})
			})
		},
		updateOriginIsOn: function (originId, originAddr, isOn) {
			let message
			let resultMessage
			if (isOn) {
				message = "确定要启用此源站（" + originAddr + "）吗？"
				resultMessage = "启用成功"
			} else {
				message = "确定要停用此源站（" + originAddr + "）吗？"
				resultMessage = "停用成功"
			}
			let that = this
			teaweb.confirm(message, function () {
				Tea.action("/servers/server/settings/origins/updateIsOn?" + that.vParams + "&originId=" + originId + "&isOn=" + (isOn ? 1 : 0))
					.post()
					.success(function () {
						teaweb.success(resultMessage, function () {
							window.location.reload()
						})
					})
			})
		}
	},
	template: `<div>
	<h3>主要源站 <a href="" @click.prevent="createPrimaryOrigin()">[添加主要源站]</a> </h3>
	<p class="comment" v-if="primaryOrigins.length == 0">暂时还没有主要源站。</p>
	<origin-list-table v-if="primaryOrigins.length > 0" :v-origins="vPrimaryOrigins" :v-origin-type="'primary'" @delete-origin="deleteOrigin" @update-origin="updateOrigin" @update-origin-is-on="updateOriginIsOn"></origin-list-table>

	<h3>备用源站 <a href="" @click.prevent="createBackupOrigin()">[添加备用源站]</a></h3>
	<p class="comment" v-if="backupOrigins.length == 0">暂时还没有备用源站。</p>
	<origin-list-table v-if="backupOrigins.length > 0" :v-origins="backupOrigins" :v-origin-type="'backup'" @delete-origin="deleteOrigin" @update-origin="updateOrigin" @update-origin-is-on="updateOriginIsOn"></origin-list-table>
</div>`
})

Vue.component("origin-list-table", {
	props: ["v-origins", "v-origin-type"],
	data: function () {
		let hasMatchedDomains = false
		let origins = this.vOrigins
		if (origins != null && origins.length > 0) {
			origins.forEach(function (origin) {
				if (origin.domains != null && origin.domains.length > 0) {
					hasMatchedDomains = true
				}
			})
		}

		return {
			hasMatchedDomains: hasMatchedDomains
		}
	},
	methods: {
		deleteOrigin: function (originId, originAddr) {
			this.$emit("delete-origin", originId, originAddr, this.vOriginType)
		},
		updateOrigin: function (originId) {
			this.$emit("update-origin", originId, this.vOriginType)
		},
		updateOriginIsOn: function (originId, originAddr, isOn) {
			this.$emit("update-origin-is-on", originId, originAddr, isOn)
		}
	},
	template: `
<table class="ui table selectable">
	<thead>
		<tr>
			<th>源站地址</th>
			<th class="width5">权重</th>
			<th class="width6">状态</th>
			<th class="three op">操作</th>
		</tr>	
	</thead>
	<tbody>
		<tr v-for="origin in vOrigins">
			<td :class="{disabled:!origin.isOn}">
				<a href="" @click.prevent="updateOrigin(origin.id)" :class="{disabled:!origin.isOn}">{{origin.addr}} &nbsp;<i class="icon expand small"></i></a>
				<div style="margin-top: 0.3em">
					<tiny-basic-label class="grey border-grey" v-if="origin.isOSS"><i class="icon hdd outline"></i>对象存储</tiny-basic-label>
					<tiny-basic-label class="grey border-grey" v-if="origin.name.length > 0">{{origin.name}}</tiny-basic-label>
					<tiny-basic-label class="grey border-grey" v-if="origin.hasCert">证书</tiny-basic-label>
					<tiny-basic-label class="grey border-grey" v-if="origin.host != null && origin.host.length > 0">主机名: {{origin.host}}</tiny-basic-label>
					<tiny-basic-label class="grey border-grey" v-if="origin.followPort">端口跟随</tiny-basic-label>
					<tiny-basic-label class="grey border-grey" v-if="origin.addr != null && origin.addr.startsWith('https://') && origin.http2Enabled">HTTP/2</tiny-basic-label>
	
					<span v-if="origin.domains != null && origin.domains.length > 0"><tiny-basic-label class="grey border-grey" v-for="domain in origin.domains">匹配: {{domain}}</tiny-basic-label></span>
					<span v-else-if="hasMatchedDomains"><tiny-basic-label class="grey  border-grey">匹配: 所有域名</tiny-basic-label></span>
				</div>
			</td>
			<td :class="{disabled:!origin.isOn}">{{origin.weight}}</td>
			<td>
				<label-on :v-is-on="origin.isOn"></label-on>
			</td>
			<td>
				<a href="" @click.prevent="updateOrigin(origin.id)">修改</a> &nbsp;
				<a href="" v-if="origin.isOn" @click.prevent="updateOriginIsOn(origin.id, origin.addr, false)">停用</a><a href=""  v-if="!origin.isOn" @click.prevent="updateOriginIsOn(origin.id, origin.addr, true)"><span class="red">启用</span></a> &nbsp;
				<a href="" @click.prevent="deleteOrigin(origin.id, origin.addr)">删除</a>
			</td>
		</tr>
	</tbody>
</table>`
})