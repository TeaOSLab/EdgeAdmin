// 使用Icon的链接方式
Vue.component("link-icon", {
	props: ["href", "title", "target", "size"],
	data: function () {
		let realSize = this.size
		if (realSize == null || realSize.length == 0) {
			realSize = "small"
		}

		return {
			vTitle: (this.title == null) ? "打开链接" : this.title,
			realSize: realSize
		}
	},
	template: `<span><slot></slot>&nbsp;<a :href="href" :title="vTitle" class="link grey" :target="target"><i class="icon linkify" :class="realSize"></i></a></span>`
})

// 带有下划虚线的连接
Vue.component("link-red", {
	props: ["href", "title"],
	data: function () {
		let href = this.href
		if (href == null) {
			href = ""
		}
		return {
			vHref: href
		}
	},
	methods: {
		clickPrevent: function () {
			emitClick(this, arguments)

			if (this.vHref.length > 0) {
				window.location = this.vHref
			}
		}
	},
	template: `<a :href="vHref" :title="title" style="border-bottom: 1px #db2828 dashed" @click.prevent="clickPrevent"><span class="red"><slot></slot></span></a>`
})

// 会弹出窗口的链接
Vue.component("link-popup", {
	props: ["title"],
	methods: {
		clickPrevent: function () {
			emitClick(this, arguments)
		}
	},
	template: `<a href="" :title="title" @click.prevent="clickPrevent"><slot></slot></a>`
})

Vue.component("popup-icon", {
	props: ["title", "href", "height"],
	methods: {
		clickPrevent: function () {
			if (this.href != null && this.href.length > 0) {
				teaweb.popup(this.href, {
					height: this.height
				})
			}
		}
	},
	template: `<span><slot></slot>&nbsp;<a href="" :title="title" @click.prevent="clickPrevent"><i class="icon expand small"></i></a></span>`
})

// 小提示
Vue.component("tip-icon", {
	props: ["content"],
	methods: {
		showTip: function () {
			teaweb.popupTip(this.content)
		}
	},
	template: `<a href="" title="查看帮助" @click.prevent="showTip"><i class="icon question circle grey"></i></a>`
})

// 提交点击事件
function emitClick(obj, arguments) {
	let event = "click"
	let newArgs = [event]
	for (let i = 0; i < arguments.length; i++) {
		newArgs.push(arguments[i])
	}
	obj.$emit.apply(obj, newArgs)
}