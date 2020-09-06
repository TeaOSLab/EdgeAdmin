/**
 * 菜单项
 */
Vue.component("inner-menu-item", {
	props: ["href", "active", "code"],
	data: function () {
		var active = this.active;
		if (typeof(active) =="undefined") {
			var itemCode = "";
			if (typeof (window.TEA.ACTION.data.firstMenuItem) != "undefined") {
				itemCode = window.TEA.ACTION.data.firstMenuItem;
			}
			active = (itemCode == this.code);
		}
		return {
			vHref: (this.href == null) ? "" : this.href,
			vActive: active
		};
	},
	template: '\
		<a :href="vHref" class="item right" style="color:#4183c4" :class="{active:vActive}">[<slot></slot>]</a> \
		'
});