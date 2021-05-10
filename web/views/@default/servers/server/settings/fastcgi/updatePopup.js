Tea.context(function () {
	this.params = this.fastcgi.params
	if (this.params == null) {
		this.params = []
	} else {
		this.params.forEach(function (v) {
			switch (v.name) {
				case "DOCUMENT_ROOT":
					v.nameZh = "文档目录"
					break;
				case "SCRIPT_FILENAME":
					v.nameZh = "脚本文件"
					break
			}
		})
	}


	this.addParam = function () {
		this.params.push({
			"name": "",
			"value": "",
			"nameZh": ""
		})
		this.$delay(function () {
			this.$find("form input[name='paramNames']").last().focus()
		})
	}

	this.removeParam = function (index) {
		this.params.$remove(index)
	}
})