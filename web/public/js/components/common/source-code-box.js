let sourceCodeBoxIndex = 0

Vue.component("source-code-box", {
	props: ["name", "type", "id", "read-only"],
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
			indentWithTabs: true
		})
		boxEditor.setValue(value)

		let info = CodeMirror.findModeByMIME(this.type)
		if (info != null) {
			boxEditor.setOption("mode", info.mode)
			CodeMirror.modeURL = "/codemirror/mode/%N/%N.js"
			CodeMirror.autoLoadMode(boxEditor, info.mode)
		}
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
	template: `<div class="source-code-box">
	<div style="display: none" :id="valueBoxId"><slot></slot></div>
	<textarea :id="'source-code-box-' + index" :name="name"></textarea>
</div>`
})