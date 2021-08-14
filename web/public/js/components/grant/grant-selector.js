Vue.component("grant-selector", {
	props: ["vGrant"],
	data: function () {
		return {
			grantId: (this.vGrant == null) ? 0 : this.vGrant.id,
			grant: this.vGrant
		}
	},
	methods: {
		// 选择授权
		select: function () {
			let that = this;
			teaweb.popup("/clusters/grants/selectPopup", {
				callback: (resp) => {
					that.grantId = resp.data.grant.id;
					if (that.grantId > 0) {
						that.grant = resp.data.grant;
					}
					that.notifyUpdate()
				}
			})
		},

		// 创建授权
		create: function () {
			let that = this
			teaweb.popup("/clusters/grants/createPopup", {
				height: "26em",
				callback: (resp) => {
					that.grantId = resp.data.grant.id;
					if (that.grantId > 0) {
						that.grant = resp.data.grant;
					}
					that.notifyUpdate()
				}
			})
		},

		// 修改授权
		update: function () {
			if (this.grant == null) {
				window.location.reload();
				return;
			}
			let that = this
			teaweb.popup("/clusters/grants/updatePopup?grantId=" + this.grant.id, {
				height: "26em",
				callback: (resp) => {
					that.grant = resp.data.grant
					that.notifyUpdate()
				}
			})
		},

		// 删除已选择授权
		remove: function () {
			this.grant = null
			this.grantId = 0
			this.notifyUpdate()
		},
		notifyUpdate: function () {
			this.$emit("change", this.grant)
		}
	},
	template: `<div>
	<input type="hidden" name="grantId" :value="grantId"/>
	<div class="ui label small basic" v-if="grant != null">{{grant.name}}<span class="small">（{{grant.methodName}}）</span> <a href="" title="修改" @click.prevent="update()"><i class="icon pencil small"></i></a> <a href="" title="删除" @click.prevent="remove()"><i class="icon remove"></i></a> </div>
	<div v-if="grant == null">
		<a href="" @click.prevent="select()">[选择已有认证]</a> &nbsp; &nbsp; <a href="" @click.prevent="create()">[添加新认证]</a>
	</div>
</div>`
})