Vue.component("loading-message", {
	template: `<div class="ui message loading">
        <div class="ui active inline loader small"></div>  &nbsp; <slot></slot>
    </div>`
})