Vue.component("label-on", {
	props: ["v-is-on"],
	template: '<div><span v-if="vIsOn" class="ui label tiny green basic">已启用</span><span v-if="!vIsOn" class="ui label tiny red basic">已停用</span></div>'
})