// 节点角色名称
Vue.component("node-role-name", {
	props: ["v-role"],
	data: function () {
		let roleName = ""
		switch (this.vRole) {
			case "node":
				roleName = "边缘节点"
				break
			case "monitor":
				roleName = "监控节点"
				break
			case "api":
				roleName = "API节点"
				break
			case "user":
				roleName = "用户平台"
				break
			case "admin":
				roleName = "管理平台"
				break
			case "database":
				roleName = "数据库节点"
				break
			case "dns":
				roleName = "DNS节点"
				break
			case "report":
				roleName = "区域监控终端"
				break
		}
		return {
			roleName: roleName
		}
	},
	template: `<span>{{roleName}}</span>`
})