Tea.context(function () {
	this.tasks.forEach(function (v) {
		switch (v.level) {
			case "good":
				v.color = "green"
				break
			case "normal":
				v.color = "blue"
				break
			case "bad":
				v.color = "orange"
				break
			case "broken":
				v.color = "red"
				break
		}
	})
})