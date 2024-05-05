Vue.component("http-access-log-box", {
	props: ["v-access-log", "v-keyword", "v-show-server-link"],
	data: function () {
		let accessLog = this.vAccessLog
		if (accessLog.header != null && accessLog.header.Upgrade != null && accessLog.header.Upgrade.values != null && accessLog.header.Upgrade.values.$contains("websocket")) {
			if (accessLog.scheme == "http") {
				accessLog.scheme = "ws"
			} else if (accessLog.scheme == "https") {
				accessLog.scheme = "wss"
			}
		}

		// 对TAG去重
		if (accessLog.tags != null && accessLog.tags.length > 0) {
			let tagMap = {}
			accessLog.tags = accessLog.tags.$filter(function (k, tag) {
				let b = (typeof (tagMap[tag]) == "undefined")
				tagMap[tag] = true
				return b
			})
		}

		// 域名
		accessLog.unicodeHost = ""
		if (accessLog.host != null && accessLog.host.startsWith("xn--")) {
			// port
			let portIndex = accessLog.host.indexOf(":")
			if (portIndex > 0) {
				accessLog.unicodeHost = punycode.ToUnicode(accessLog.host.substring(0, portIndex))
			} else {
				accessLog.unicodeHost = punycode.ToUnicode(accessLog.host)
			}
		}

		return {
			accessLog: accessLog
		}
	},
	methods: {
		formatCost: function (seconds) {
			if (seconds == null) {
				return "0"
			}
			let s = (seconds * 1000).toString();
			let pieces = s.split(".");
			if (pieces.length < 2) {
				return s;
			}

			return pieces[0] + "." + pieces[1].substring(0, 3);
		},
		showLog: function () {
			let that = this
			let requestId = this.accessLog.requestId
			this.$parent.$children.forEach(function (v) {
				if (v.deselect != null) {
					v.deselect()
				}
			})
			this.select()
			teaweb.popup("/servers/server/log/viewPopup?requestId=" + requestId, {
				width: "50em",
				height: "28em",
				onClose: function () {
					that.deselect()
				}
			})
		},
		select: function () {
			this.$refs.box.parentNode.style.cssText = "background: rgba(0, 0, 0, 0.1)"
		},
		deselect: function () {
			this.$refs.box.parentNode.style.cssText = ""
		},
		mismatch: function () {
			teaweb.warn("当前访问没有匹配到任何网站")
		}
	},
	template: `<div style="word-break: break-all" :style="{'color': (accessLog.status >= 400) ? '#dc143c' : ''}" ref="box">
	<div>
		<a v-if="accessLog.node != null && accessLog.node.nodeCluster != null" :href="'/clusters/cluster/node?nodeId=' + accessLog.node.id + '&clusterId=' + accessLog.node.nodeCluster.id" title="点击查看节点详情" target="_top"><span class="grey">[{{accessLog.node.name}}<span v-if="!accessLog.node.name.endsWith('节点')">节点</span>]</span></a>
		
		<!-- 网站 -->
		<a :href="'/servers/server/log?serverId=' + accessLog.serverId" title="点击到网站" v-if="vShowServerLink && accessLog.serverId > 0"><span class="grey">[网站]</span></a>
		<span v-if="vShowServerLink && (accessLog.serverId == null || accessLog.serverId == 0)" @click.prevent="mismatch()"><span class="disabled">[网站]</span></span>
		
		<span v-if="accessLog.region != null && accessLog.region.length > 0" class="grey"><ip-box :v-ip="accessLog.remoteAddr">[{{accessLog.region}}]</ip-box></span> 
		<ip-box><keyword :v-word="vKeyword">{{accessLog.remoteAddr}}</keyword></ip-box> [{{accessLog.timeLocal}}] <em>&quot;<keyword :v-word="vKeyword">{{accessLog.requestMethod}}</keyword> {{accessLog.scheme}}://<keyword :v-word="vKeyword">{{accessLog.host}}</keyword><keyword :v-word="vKeyword">{{accessLog.requestURI}}</keyword> <a :href="accessLog.scheme + '://' + accessLog.host + accessLog.requestURI" target="_blank" title="新窗口打开" class="disabled"><i class="external icon tiny"></i> </a> {{accessLog.proto}}&quot; </em> <keyword :v-word="vKeyword">{{accessLog.status}}</keyword> 
		
		<code-label v-if="accessLog.unicodeHost != null && accessLog.unicodeHost.length > 0">{{accessLog.unicodeHost}}</code-label>
		
		<!-- attrs -->
		<code-label v-if="accessLog.attrs != null && (accessLog.attrs['cache.status'] == 'HIT' || accessLog.attrs['cache.status'] == 'STALE')">cache {{accessLog.attrs['cache.status'].toLowerCase()}}</code-label> 
		<!-- waf -->
		<code-label v-if="accessLog.firewallActions != null && accessLog.firewallActions.length > 0">waf {{accessLog.firewallActions}}</code-label> 
		
		<!-- tags -->
		<span v-if="accessLog.tags != null && accessLog.tags.length > 0">- <code-label v-for="tag in accessLog.tags" :key="tag">{{tag}}</code-label>
		</span>
		<span  v-if="accessLog.wafInfo != null">
			<a :href="(accessLog.wafInfo.policy.serverId == 0) ? '/servers/components/waf/group?firewallPolicyId=' +  accessLog.firewallPolicyId + '&type=inbound&groupId=' + accessLog.firewallRuleGroupId+ '#set' + accessLog.firewallRuleSetId : '/servers/server/settings/waf/group?serverId=' + accessLog.serverId + '&firewallPolicyId=' + accessLog.firewallPolicyId + '&type=inbound&groupId=' + accessLog.firewallRuleGroupId + '#set' + accessLog.firewallRuleSetId" target="_blank">
				<code-label-plain>
					<span>
						WAF -
						<span v-if="accessLog.wafInfo.group != null">{{accessLog.wafInfo.group.name}} -</span>
						<span v-if="accessLog.wafInfo.set != null">{{accessLog.wafInfo.set.name}}</span>
					</span>
				</code-label-plain>
			</a>
		</span>
			
		<span v-if="accessLog.requestTime != null"> - 耗时:{{formatCost(accessLog.requestTime)}} ms </span><span v-if="accessLog.humanTime != null && accessLog.humanTime.length > 0" class="grey small">&nbsp; ({{accessLog.humanTime}})</span>
		&nbsp; <a href="" @click.prevent="showLog" title="查看详情"><i class="icon expand"></i></a>
	</div>
</div>`
})

// Javascript Punycode converter derived from example in RFC3492.
// This implementation is created by some@domain.name and released into public domain
// 代码来自：https://stackoverflow.com/questions/183485/converting-punycode-with-dash-character-to-unicode
var punycode = new function Punycode() {
	// This object converts to and from puny-code used in IDN
	//
	// punycode.ToASCII ( domain )
	//
	// Returns a puny coded representation of "domain".
	// It only converts the part of the domain name that
	// has non ASCII characters. I.e. it dosent matter if
	// you call it with a domain that already is in ASCII.
	//
	// punycode.ToUnicode (domain)
	//
	// Converts a puny-coded domain name to unicode.
	// It only converts the puny-coded parts of the domain name.
	// I.e. it dosent matter if you call it on a string
	// that already has been converted to unicode.
	//
	//
	this.utf16 = {
		// The utf16-class is necessary to convert from javascripts internal character representation to unicode and back.
		decode: function (input) {
			var output = [], i = 0, len = input.length, value, extra;
			while (i < len) {
				value = input.charCodeAt(i++);
				if ((value & 0xF800) === 0xD800) {
					extra = input.charCodeAt(i++);
					if (((value & 0xFC00) !== 0xD800) || ((extra & 0xFC00) !== 0xDC00)) {
						throw new RangeError("UTF-16(decode): Illegal UTF-16 sequence");
					}
					value = ((value & 0x3FF) << 10) + (extra & 0x3FF) + 0x10000;
				}
				output.push(value);
			}
			return output;
		},
		encode: function (input) {
			var output = [], i = 0, len = input.length, value;
			while (i < len) {
				value = input[i++];
				if ((value & 0xF800) === 0xD800) {
					throw new RangeError("UTF-16(encode): Illegal UTF-16 value");
				}
				if (value > 0xFFFF) {
					value -= 0x10000;
					output.push(String.fromCharCode(((value >>> 10) & 0x3FF) | 0xD800));
					value = 0xDC00 | (value & 0x3FF);
				}
				output.push(String.fromCharCode(value));
			}
			return output.join("");
		}
	}

	//Default parameters
	var initial_n = 0x80;
	var initial_bias = 72;
	var delimiter = "\x2D";
	var base = 36;
	var damp = 700;
	var tmin = 1;
	var tmax = 26;
	var skew = 38;
	var maxint = 0x7FFFFFFF;

	// decode_digit(cp) returns the numeric value of a basic code
	// point (for use in representing integers) in the range 0 to
	// base-1, or base if cp is does not represent a value.

	function decode_digit(cp) {
		return cp - 48 < 10 ? cp - 22 : cp - 65 < 26 ? cp - 65 : cp - 97 < 26 ? cp - 97 : base;
	}

	// encode_digit(d,flag) returns the basic code point whose value
	// (when used for representing integers) is d, which needs to be in
	// the range 0 to base-1. The lowercase form is used unless flag is
	// nonzero, in which case the uppercase form is used. The behavior
	// is undefined if flag is nonzero and digit d has no uppercase form.

	function encode_digit(d, flag) {
		return d + 22 + 75 * (d < 26) - ((flag != 0) << 5);
		//  0..25 map to ASCII a..z or A..Z
		// 26..35 map to ASCII 0..9
	}

	//** Bias adaptation function **
	function adapt(delta, numpoints, firsttime) {
		var k;
		delta = firsttime ? Math.floor(delta / damp) : (delta >> 1);
		delta += Math.floor(delta / numpoints);

		for (k = 0; delta > (((base - tmin) * tmax) >> 1); k += base) {
			delta = Math.floor(delta / (base - tmin));
		}
		return Math.floor(k + (base - tmin + 1) * delta / (delta + skew));
	}

	// encode_basic(bcp,flag) forces a basic code point to lowercase if flag is zero,
	// uppercase if flag is nonzero, and returns the resulting code point.
	// The code point is unchanged if it is caseless.
	// The behavior is undefined if bcp is not a basic code point.

	function encode_basic(bcp, flag) {
		bcp -= (bcp - 97 < 26) << 5;
		return bcp + ((!flag && (bcp - 65 < 26)) << 5);
	}

	// Main decode
	this.decode = function (input, preserveCase) {
		// Dont use utf16
		var output = [];
		var case_flags = [];
		var input_length = input.length;

		var n, out, i, bias, basic, j, ic, oldi, w, k, digit, t, len;

		// Initialize the state:

		n = initial_n;
		i = 0;
		bias = initial_bias;

		// Handle the basic code points: Let basic be the number of input code
		// points before the last delimiter, or 0 if there is none, then
		// copy the first basic code points to the output.

		basic = input.lastIndexOf(delimiter);
		if (basic < 0) basic = 0;

		for (j = 0; j < basic; ++j) {
			if (preserveCase) case_flags[output.length] = (input.charCodeAt(j) - 65 < 26);
			if (input.charCodeAt(j) >= 0x80) {
				throw new RangeError("Illegal input >= 0x80");
			}
			output.push(input.charCodeAt(j));
		}

		// Main decoding loop: Start just after the last delimiter if any
		// basic code points were copied; start at the beginning otherwise.

		for (ic = basic > 0 ? basic + 1 : 0; ic < input_length;) {

			// ic is the index of the next character to be consumed,

			// Decode a generalized variable-length integer into delta,
			// which gets added to i. The overflow checking is easier
			// if we increase i as we go, then subtract off its starting
			// value at the end to obtain delta.
			for (oldi = i, w = 1, k = base; ; k += base) {
				if (ic >= input_length) {
					throw RangeError("punycode_bad_input(1)");
				}
				digit = decode_digit(input.charCodeAt(ic++));

				if (digit >= base) {
					throw RangeError("punycode_bad_input(2)");
				}
				if (digit > Math.floor((maxint - i) / w)) {
					throw RangeError("punycode_overflow(1)");
				}
				i += digit * w;
				t = k <= bias ? tmin : k >= bias + tmax ? tmax : k - bias;
				if (digit < t) {
					break;
				}
				if (w > Math.floor(maxint / (base - t))) {
					throw RangeError("punycode_overflow(2)");
				}
				w *= (base - t);
			}

			out = output.length + 1;
			bias = adapt(i - oldi, out, oldi === 0);

			// i was supposed to wrap around from out to 0,
			// incrementing n each time, so we'll fix that now:
			if (Math.floor(i / out) > maxint - n) {
				throw RangeError("punycode_overflow(3)");
			}
			n += Math.floor(i / out);
			i %= out;

			// Insert n at position i of the output:
			// Case of last character determines uppercase flag:
			if (preserveCase) {
				case_flags.splice(i, 0, input.charCodeAt(ic - 1) - 65 < 26);
			}

			output.splice(i, 0, n);
			i++;
		}
		if (preserveCase) {
			for (i = 0, len = output.length; i < len; i++) {
				if (case_flags[i]) {
					output[i] = (String.fromCharCode(output[i]).toUpperCase()).charCodeAt(0);
				}
			}
		}
		return this.utf16.encode(output);
	};

	//** Main encode function **

	this.encode = function (input, preserveCase) {
		//** Bias adaptation function **

		var n, delta, h, b, bias, j, m, q, k, t, ijv, case_flags;

		if (preserveCase) {
			// Preserve case, step1 of 2: Get a list of the unaltered string
			case_flags = this.utf16.decode(input);
		}
		// Converts the input in UTF-16 to Unicode
		input = this.utf16.decode(input.toLowerCase());

		var input_length = input.length; // Cache the length

		if (preserveCase) {
			// Preserve case, step2 of 2: Modify the list to true/false
			for (j = 0; j < input_length; j++) {
				case_flags[j] = input[j] != case_flags[j];
			}
		}

		var output = [];


		// Initialize the state:
		n = initial_n;
		delta = 0;
		bias = initial_bias;

		// Handle the basic code points:
		for (j = 0; j < input_length; ++j) {
			if (input[j] < 0x80) {
				output.push(
					String.fromCharCode(
						case_flags ? encode_basic(input[j], case_flags[j]) : input[j]
					)
				);
			}
		}

		h = b = output.length;

		// h is the number of code points that have been handled, b is the
		// number of basic code points

		if (b > 0) output.push(delimiter);

		// Main encoding loop:
		//
		while (h < input_length) {
			// All non-basic code points < n have been
			// handled already. Find the next larger one:

			for (m = maxint, j = 0; j < input_length; ++j) {
				ijv = input[j];
				if (ijv >= n && ijv < m) m = ijv;
			}

			// Increase delta enough to advance the decoder's
			// <n,i> state to <m,0>, but guard against overflow:

			if (m - n > Math.floor((maxint - delta) / (h + 1))) {
				throw RangeError("punycode_overflow (1)");
			}
			delta += (m - n) * (h + 1);
			n = m;

			for (j = 0; j < input_length; ++j) {
				ijv = input[j];

				if (ijv < n) {
					if (++delta > maxint) return Error("punycode_overflow(2)");
				}

				if (ijv == n) {
					// Represent delta as a generalized variable-length integer:
					for (q = delta, k = base; ; k += base) {
						t = k <= bias ? tmin : k >= bias + tmax ? tmax : k - bias;
						if (q < t) break;
						output.push(String.fromCharCode(encode_digit(t + (q - t) % (base - t), 0)));
						q = Math.floor((q - t) / (base - t));
					}
					output.push(String.fromCharCode(encode_digit(q, preserveCase && case_flags[j] ? 1 : 0)));
					bias = adapt(delta, h + 1, h == b);
					delta = 0;
					++h;
				}
			}

			++delta, ++n;
		}
		return output.join("");
	}

	this.ToASCII = function (domain) {
		var domain_array = domain.split(".");
		var out = [];
		for (var i = 0; i < domain_array.length; ++i) {
			var s = domain_array[i];
			out.push(
				s.match(/[^A-Za-z0-9-]/) ?
					"xn--" + punycode.encode(s) :
					s
			);
		}
		return out.join(".");
	}
	this.ToUnicode = function (domain) {
		var domain_array = domain.split(".");
		var out = [];
		for (var i = 0; i < domain_array.length; ++i) {
			var s = domain_array[i];
			out.push(
				s.match(/^xn--/) ?
					punycode.decode(s.slice(4)) :
					s
			);
		}
		return out.join(".");
	}
}();