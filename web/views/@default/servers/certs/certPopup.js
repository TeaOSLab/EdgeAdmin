Tea.context(function () {
	// 打印缩进
	this.indent = function (index) {
		let indent = ""
		for (let i = 0; i < index; i++) {
			indent += " &nbsp; &nbsp; "
		}
		return indent
	}
})