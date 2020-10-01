/**
 * Tea.Date 对象
 *
 * @class Tea.Date
 */
/**
 * Tea.Date构造器。使用方法如：<br/>
 * var date = new Tea.Date();<br/>
 * var date = new Tea.Date("Y-m-d H:i:s");<br/>
 * var date = new Tea.Date("Y-m-d H:i:s", 1169226085);
 *
 * @constructor Tea.Date
 * @param String format 时间格式，为可选参数，目前支持O,r,Y,y,L,M,m,n,F,t,w,D,l,d,z,H,i,s,j,h,G,g,a,A等字符。
 * @param int time 时间戳，为可选参数
 */
Tea.Date = function (format, time) {
	var date = new Date();

	if (typeof(format) == "undefined") {
		format = "r";
	}

	if (typeof(time) != "undefined") {
		time = parseInt(time, 10);
		date.setTime(time);
	}

	//parse char
	this.get = function (chr) {
		if ((chr >= "a" && chr <= "z") || (chr >= "A" && chr <= "Z")) {
			var func = "_parse_" + chr;
			if (this[func]) {
				return this[func]();
			}
		}
		return chr;
	};

	/**
	 * 根据提供的格式取得对应的时间格式
	 *
	 * @method parse
	 * @param String format
	 */
	this.parse = function (format) {
		var result = "";
		if (format.length > 0) {
			for (var i=0; i<format.length; i++) {
				var chr = format.charAt(i);
				result += this.get(chr);
			}
		}
		return result;
	};

	/**
	 * 设置某一时间为某值
	 *
	 * @method set
	 * @param String type 时间选项，如 d 表示天,Y 表示年,H 表示小时，等等。
	 * @param int value 新的值
	 */
	this.set = function (type, value) {
		value = parseInt(value, 10);
		switch (type) {
			case "d":
				date.setDate(value);
				break;
			case "Y":
				date.setFullYear(value);
				break;
			case "H":
			case "G":
				date.setHours(value);
				break;
			case "i":
				date.setMinutes(value);
				break;
			case "s":
				date.setSeconds(value);
				break;
			case "m":
			case "n":
				date.setMonth(value - 1);
				break;
		}
	};

	//timezone
	this._parse_O = function () {
		var hours = (Math.abs(date.getTimezoneOffset()/60)).toString();
		if (hours.length == 1) {
			hours = "0" + hours;
		}
		return "+" + hours + "00";
	};

	this._parse_r = function () {
		return this.parse("D, d M Y H:i:s O");
	};

	//parse year
	this._parse_Y = function () {
		return date.getFullYear().toString();
	};

	this._parse_y = function () {
		var y = this._parse_Y();
		return y.substr(2);
	};

	this._parse_L = function () {
		var y = parseInt(this.parse("Y"));
		if (y%4 ==0  && (y%100 > 0 || y%400 == 0)) {
			return "1";
		}
		return "0";
	};

	//month
	this._parse_m = function () {
		var n = this._parse_n();
		if (n.length < 2) {
			n = "0" + n;
		}
		return n;
	};

	this._parse_n = function () {
		return (date.getMonth() + 1).toString();
	};

	this._parse_t = function () {
		var t = 32 - new Date(this.get("Y"), this.get("m") - 1 , 32).getDate();
		return t;
	};

	this._parse_F = function () {
		var n = parseInt(this.parse("n"));
		var months = ["", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"];
		return months[n];
	};

	this._parse_M = function () {
		var n = parseInt(this.parse("n"));
		var months = ["", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
		return months[n];
	};

	//week
	this._parse_w = function () {
		return date.getDay().toString();
	};

	this._parse_D = function () {
		var w = parseInt(this._parse_w());
		var days = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];
		return days[w];
	};

	this._parse_l = function () {
		var w = parseInt(this._parse_w());
		var days = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
		return days[w];
	};

	//day
	this._parse_d = function () {
		var j = this._parse_j();
		if (j.length < 2) {
			j = "0" + j;
		}
		return j;
	};

	this._parse_j = function () {
		return date.getDate().toString();
	};

	this._parse_W = function () {
		var _date = new Tea.Date();
		_date.set("m", 1);
		_date.set("d", 1);
		var w = parseInt(_date.parse("w"));
		var m = parseInt(this.parse("m"), 10);
		var total = 0;
		for (var i=1; i<m; i++) {
			var date2 = new Tea.Date();
			date2.set("m", i);
			var t = parseInt(date2.parse("t"));
			total += t;
		}
		total += parseInt(this.parse("d"), 10);
		var w2 = parseInt(this.parse("w"));
		total = total - w2 + (w - 1);
		var weeks = 0;
		if (w2 != 0) {
			weeks = (total/7 + 1).toString();
		}
		else {
			weeks = (total/7).toString();
		}
		if (weeks.length == 1) {
			weeks = "0" + weeks;
		}
		return weeks;
	};

	this._parse_z = function () {
		var m = parseInt(this.parse("m"), 10);
		var total = 0;
		for (var i=1; i<m; i++) {
			var date2 = new Tea.Date();
			date2.set("m", i);
			var t = parseInt(date2.parse("t"));
			total += t;
		}
		total += parseInt(this.parse("d"), 10) - 1;
		return total;
	};

	//minute
	this._parse_i = function () {
		var i = date.getMinutes().toString();
		if (i.length < 2) {
			i = "0" + i;
		}
		return i;
	};

	//second
	this._parse_s = function () {
		var s = date.getSeconds().toString();
		if (s.length < 2) {
			s = "0" + s;
		}
		return s;
	};

	//hour
	this._parse_H = function () {
		var H = this._parse_G();
		if (H.length < 2) {
			H = "0" + H;
		}
		return H;
	};

	this._parse_G = function () {
		return date.getHours().toString();
	};

	this._parse_h = function () {
		var h = this._parse_g();
		if (h.length < 2) {
			h = "0" + h;
		}
		return h;
	};

	this._parse_g = function () {
		var g = parseInt(this._parse_G(), 10);
		if (g > 12) {
			g = g - 12;
		}
		return g.toString();
	};

	//time
	this._parse_U = function () {
		return this.time().toString();
	};

	//am/pm
	this._parse_a = function () {
		var hour = this.parse("H");
		return (hour<12)?"am":"pm";
	};

	this._parse_A = function () {
		return this.parse("a").toUpperCase();
	};

	/**
	 * 取得当前时间对应的时间戳,代表了从 1970 年 1 月 1 日开始计算到 Date 对象中的时间之间的秒数
	 *
	 * @method time
	 * @return int
	 */
	this.time = function () {
		return Math.round(date.getTime()/1000);
	};

	/*
	 * 将该对象转换成字符串格式
	 *
	 * @method toString
	 * @return String 该对象的字符串表示形式
	 */
	this.toString = function () {
		return this.parse(format);
	};
};


Tea.Date.toTime = function (dateStr) {
	if (arguments.length == 1) {
		return Date.parse(dateStr);
	} else if (arguments.length == 3) {
		arguments[1] = parseInt(arguments[1], 10) - 1;
		return (new Date(arguments[0], arguments[1], arguments[2])).time();
	}
};

Number.prototype.dateFormat = function (format) {
	var date = new Tea.Date(format, this * 1000);
	return date.toString();
};

Date.prototype.format = function (format) {
	return new Tea.Date(format, this.getTime()).toString();
};