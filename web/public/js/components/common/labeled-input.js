Vue.component("labeled-input", {
	props: ["name", "size", "maxlength", "label", "value"],
	template: '<div class="ui input right labeled"> \
	<input type="text" :name="name" :size="size" :maxlength="maxlength" :value="value"/>\
	<span class="ui label">{{label}}</span>\
</div>'
});