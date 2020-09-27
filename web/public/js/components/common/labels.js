Vue.component("label-on", {
	props: ["v-is-on"],
	template: '<div><span v-if="vIsOn" class="ui label tiny green">已启用</span><span v-if="!vIsOn" class="ui label tiny red">已关闭</span></div>'
})