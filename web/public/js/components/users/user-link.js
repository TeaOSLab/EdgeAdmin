Vue.component("user-link", {
	props: ["v-user", "v-keyword"],
	data: function () {
		let user = this.vUser
		if (user == null) {
			user = {id: 0, "username": "", "fullname": ""}
		}
		return {
			user: user
		}
	},
	template: `<div style="display: inline-block">
	<span v-if="user.id > 0"><keyword :v-word="vKeyword">{{user.fullname}}</keyword><span class="small grey">（<keyword :v-word="vKeyword">{{user.username}}</keyword>）</span></span>
	<span v-else class="disabled">[已删除]</span>
</div>`
})