let sourceCodeBoxIndex = 0

Vue.component("source-code-box", {
	props: ["name", "type", "id", "read-only", "width", "height", "focus"],
	mounted: function () {
		let readOnly = this.readOnly
		if (typeof readOnly != "boolean") {
			readOnly = true
		}
		let box = document.getElementById("source-code-box-" + this.index)
		let valueBox = document.getElementById(this.valueBoxId)
		let value = ""
		if (valueBox.textContent != null) {
			value = valueBox.textContent
		} else if (valueBox.innerText != null) {
			value = valueBox.innerText
		}

		this.createEditor(box, value, readOnly)
	},
	data: function () {
		let index = sourceCodeBoxIndex++

		let valueBoxId = 'source-code-box-value-' + sourceCodeBoxIndex
		if (this.id != null) {
			valueBoxId = this.id
		}

		return {
			index: index,
			valueBoxId: valueBoxId
		}
	},
	methods: {
		createEditor: function (box, value, readOnly) {
			let boxEditor = CodeMirror.fromTextArea(box, {
				theme: "idea",
				lineNumbers: true,
				value: "",
				readOnly: readOnly,
				showCursorWhenSelecting: true,
				height: "auto",
				//scrollbarStyle: null,
				viewportMargin: Infinity,
				lineWrapping: true,
				highlightFormatting: false,
				indentUnit: 4,
				indentWithTabs: true,
			})
			let that = this
			boxEditor.on("change", function () {
				that.change(boxEditor.getValue())
			})
			boxEditor.setValue(value)

			if (this.focus) {
				boxEditor.focus()
			}

			let width = this.width
			let height = this.height
			if (width != null && height != null) {
				width = parseInt(width)
				height = parseInt(height)
				if (!isNaN(width) && !isNaN(height)) {
					if (width <= 0) {
						width = box.parentNode.offsetWidth
					}
					boxEditor.setSize(width, height)
				}
			} else if (height != null) {
				height = parseInt(height)
				if (!isNaN(height)) {
					boxEditor.setSize("100%", height)
				}
			}

			let info = CodeMirror.findModeByMIME(this.type)
			if (info != null) {
				boxEditor.setOption("mode", info.mode)
				CodeMirror.modeURL = "/codemirror/mode/%N/%N.js"
				CodeMirror.autoLoadMode(boxEditor, info.mode)
			}
		},
		change: function (code) {
			this.$emit("change", code)
		}
	},
	template: `<div class="source-code-box">
	<div style="display: none" :id="valueBoxId"><slot></slot></div>
	<textarea :id="'source-code-box-' + index" :name="name"></textarea>
</div>`
})