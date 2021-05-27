Vue.component("not-found-box", {
	props: ["message"],
	template: `<div style="text-align: center; margin-top: 5em;">
	<div style="font-size: 2em; margin-bottom: 1em"><i class="icon exclamation triangle large grey"></i></div>
	<p class="comment">{{message}}<slot></slot></p>
</div>`
})