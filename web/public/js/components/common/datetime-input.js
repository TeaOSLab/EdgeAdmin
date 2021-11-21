Vue.component("datetime-input", {
	props: ["v-name", "v-timestamp"],
	mounted: function () {
		let that = this
		teaweb.datepicker(this.$refs.dayInput, function (v) {
			that.day = v
			that.hour = "23"
			that.minute = "59"
			that.second = "59"
			that.change()
		})
	},
	data: function () {
		let timestamp = this.vTimestamp
		if (timestamp != null) {
			timestamp = parseInt(timestamp)
			if (isNaN(timestamp)) {
				timestamp = 0
			}
		} else {
			timestamp = 0
		}

		let day = ""
		let hour = ""
		let minute = ""
		let second = ""

		if (timestamp > 0) {
			let date = new Date()
			date.setTime(timestamp * 1000)

			let year = date.getFullYear().toString()
			let month = this.leadingZero((date.getMonth() + 1).toString(), 2)
			day = year + "-" + month + "-" + this.leadingZero(date.getDate().toString(), 2)

			hour = this.leadingZero(date.getHours().toString(), 2)
			minute = this.leadingZero(date.getMinutes().toString(), 2)
			second = this.leadingZero(date.getSeconds().toString(), 2)
		}

		return {
			timestamp: timestamp,
			day: day,
			hour: hour,
			minute: minute,
			second: second,

			hasDayError: false,
			hasHourError: false,
			hasMinuteError: false,
			hasSecondError: false
		}
	},
	methods: {
		change: function () {
			let date = new Date()

			// day
			if (!/^\d{4}-\d{1,2}-\d{1,2}$/.test(this.day)) {
				this.hasDayError = true
				return
			}
			let pieces = this.day.split("-")
			let year = parseInt(pieces[0])
			date.setFullYear(year)

			let month = parseInt(pieces[1])
			if (month < 1 || month > 12) {
				this.hasDayError = true
				return
			}
			date.setMonth(month - 1)

			let day = parseInt(pieces[2])
			if (day < 1 || day > 32) {
				this.hasDayError = true
				return
			}
			date.setDate(day)

			this.hasDayError = false

			// hour
			if (!/^\d+$/.test(this.hour)) {
				this.hasHourError = true
				return
			}
			let hour = parseInt(this.hour)
			if (isNaN(hour)) {
				this.hasHourError = true
				return
			}
			if (hour < 0 || hour >= 24) {
				this.hasHourError = true
				return
			}
			this.hasHourError = false
			date.setHours(hour)

			// minute
			if (!/^\d+$/.test(this.minute)) {
				this.hasMinuteError = true
				return
			}
			let minute = parseInt(this.minute)
			if (isNaN(minute)) {
				this.hasMinuteError = true
				return
			}
			if (minute < 0 || minute >= 60) {
				this.hasMinuteError = true
				return
			}
			this.hasMinuteError = false
			date.setMinutes(minute)

			// second
			if (!/^\d+$/.test(this.second)) {
				this.hasSecondError = true
				return
			}
			let second = parseInt(this.second)
			if (isNaN(second)) {
				this.hasSecondError = true
				return
			}
			if (second < 0 || second >= 60) {
				this.hasSecondError = true
				return
			}
			this.hasSecondError = false
			date.setSeconds(second)

			this.timestamp = Math.floor(date.getTime() / 1000)
		},
		leadingZero: function (s, l) {
			if (l <= s.length) {
				return s
			}
			for (let i = 0; i < l - s.length; i++) {
				s = "0" + s
			}
			return s
		}
	},
	template: `<div>
	<input type="hidden" :name="vName" :value="timestamp"/>
	<div class="ui fields inline" style="padding: 0; margin:0">
		<div class="ui field" :class="{error: hasDayError}">
			<input type="text" v-model="day" placeholder="YYYY-MM-DD" style="width:8.6em" maxlength="10" @input="change" ref="dayInput"/>
		</div>
		<div class="ui field" :class="{error: hasHourError}"><input type="text" v-model="hour" maxlength="2" style="width:4em" placeholder="时" @input="change"/></div>
		<div class="ui field">:</div>
		<div class="ui field" :class="{error: hasMinuteError}"><input type="text" v-model="minute" maxlength="2" style="width:4em" placeholder="分" @input="change"/></div>
		<div class="ui field">:</div>
		<div class="ui field" :class="{error: hasSecondError}"><input type="text" v-model="second" maxlength="2" style="width:4em" placeholder="秒" @input="change"/></div>
	</div>
</div>`
})