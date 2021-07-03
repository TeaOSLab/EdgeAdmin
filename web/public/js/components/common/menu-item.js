/**
 * 菜单项
 */
Vue.component("menu-item", {
	props: ["href", "active", "code"],
	data: function () {
		var active = this.active
		if (typeof (active) == "undefined") {
			var itemCode = ""
			if (typeof (window.TEA.ACTION.data.firstMenuItem) != "undefined") {
				itemCode = window.TEA.ACTION.data.firstMenuItem
			}
			if (itemCode != null && itemCode.length > 0 && this.code != null && this.code.length > 0) {
				if (itemCode.indexOf(",") > 0) {
					active = itemCode.split(",").$contains(this.code)
				} else {
					active = (itemCode == this.code)
				}
			}
		}
		return {
			vHref: (this.href == null) ? "" : this.href,
			vActive: active
		}
	},
	methods: {
		click: function (e) {
			this.$emit("click", e)
		}
	},
	template: '\
		<a :href="vHref" class="item" :class="{active:vActive}" @click="click"><slot></slot></a> \
		'
});