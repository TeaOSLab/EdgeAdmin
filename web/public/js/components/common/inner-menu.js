/**
 * 二级菜单
 */
Vue.component("inner-menu", {
	template: `
		<div class="second-menu" style="width:80%;position: absolute;top:-8px;right:1em"> 
			<div class="ui menu text blue small">
				<slot></slot>
			</div> 
		</div>`
});