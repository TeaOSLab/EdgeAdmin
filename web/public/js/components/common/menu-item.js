/**
 * 菜单项
 */
Vue.component("menu-item", {
	props: ["href", "active", "code"],
	data: function () {
		let active = this.active
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

		let href = (this.href == null) ? "" : this.href
		if (typeof (href) == "string" && href.length > 0 && href.startsWith(".")) {
			let qIndex = href.indexOf("?")
			if (qIndex >= 0) {
				href = Tea.url(href.substring(0, qIndex)) + href.substring(qIndex)
			} else {
				href = Tea.url(href)
			}
		}

		return {
			vHref: href,
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