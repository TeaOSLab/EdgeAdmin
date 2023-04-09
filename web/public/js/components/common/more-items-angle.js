// 可以展示更多条目的角图表
Vue.component("more-items-angle", {
	props: ["v-data-url", "v-url"],
	data: function () {
		return {
			visible: false
		}
	},
	methods: {
		show: function () {
			this.visible = !this.visible
			if (this.visible) {
				this.showBox()
			} else {
				this.hideBox()
			}
		},
		showBox: function () {
			let that = this

			this.visible = true

			Tea.action(this.vDataUrl)
				.params({
					url: this.vUrl
				})
				.post()
				.success(function (resp) {
					let groups = resp.data.groups

					let boxLeft = that.$el.offsetLeft + 120;
					let boxTop = that.$el.offsetTop + 70;

					let box = document.createElement("div")
					box.setAttribute("id", "more-items-box")
					box.style.cssText = "z-index: 100; position: absolute; left: " + boxLeft + "px; top: " + boxTop + "px; max-height: 30em; overflow: auto; border-bottom: 1px solid rgba(34,36,38,.15)"
					document.body.append(box)

					let menuHTML = "<ul class=\"ui labeled menu vertical borderless\" style=\"padding: 0\">"
					groups.forEach(function (group) {
						menuHTML += "<div class=\"item header\">" + teaweb.encodeHTML(group.name) + "</div>"
						group.items.forEach(function (item) {
							menuHTML += "<a href=\"" + item.url + "\" class=\"item " + (item.isActive ? "active" : "") + "\" style=\"font-size: 0.9em;\">" + teaweb.encodeHTML(item.name) + "<i class=\"icon right angle\"></i></a>"
						})
					})
					menuHTML += "</ul>"
					box.innerHTML = menuHTML

					let listener = function (e) {
						if (e.target.tagName == "I") {
							return
						}

						if (!that.isInBox(box, e.target)) {
							document.removeEventListener("click", listener)
							that.hideBox()
						}
					}
					document.addEventListener("click", listener)
				})
		},
		hideBox: function () {
			let box = document.getElementById("more-items-box")
			if (box != null) {
				box.parentNode.removeChild(box)
			}
			this.visible = false
		},
		isInBox: function (parent, child) {
			while (true) {
				if (child == null) {
					break
				}
				if (child.parentNode == parent) {
					return true
				}
				child = child.parentNode
			}
			return false
		}
	},
	template: `<a href="" class="item" @click.prevent="show" style="padding-right: 0"><span style="font-size: 0.8em">切换</span><i class="icon angle" :class="{down: !visible, up: visible}"></i></a>`
})