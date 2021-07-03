Tea.context(function () {
	this.format = function (v) {
		if (v == 0) {
			return "00"
		}
		if (v < 10) {
			return "0" + v
		}
		return v.toString()
	}

	let randValues = []
	let times = []
	let count = 6
	for (let i = 0; i < count; i++) {
		randValues.push(Math.ceil(Math.random() * 100))
		switch (this.item.periodUnit) {
			case "month": {
				let date = new Date()
				date.setMonth(date.getMonth() - (count - 1 - i))
				let month = date.getMonth() + 1
				times.push(date.getFullYear() + this.format(month))
			}
				break
			case "week": {
				let date = new Date()
				times.push(date.getFullYear() + this.format(50 + i - count))
			}
				break
			case "day": {
				let date = new Date()
				date.setDate(date.getDate() - (count - i - 1))
				let day = date.getDate()
				times.push(date.getFullYear() + this.format(date.getMonth() + 1) + this.format(day))
			}
				break
			case "hour": {
				let date = new Date()
				date.setHours(date.getHours() - (count - i - 1))
				times.push(date.getFullYear() + this.format(date.getMonth() + 1) + this.format(date.getDate()) + this.format(date.getHours()))
			}
				break
			case "minute": {
				let date = new Date()
				date.setMinutes(date.getMinutes() - (count - i - 1))
				times.push(date.getFullYear() + this.format(date.getMonth() + 1) + this.format(date.getDate()) + this.format(date.getHours()) + this.format(date.getMinutes()))
			}
				break
		}
	}
	let total = randValues.$sum()

	this.testingStats = []
	let that = this
	randValues.forEach(function (v, index) {
		that.testingStats.push({
			keys: ["对象" + (index + 1)],
			value: v,
			total: total,
			time: times[index]
		})
	})
})