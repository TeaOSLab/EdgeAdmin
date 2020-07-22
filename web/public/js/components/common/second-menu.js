/**
 * 二级菜单
 */
Vue.component("second-menu", {
	template: ' \
		<div class="second-menu"> \
			<div class="ui menu text blue small">\
				<slot></slot>\
			</div> \
			<div class="ui divider"></div> \
		</div>'
});