Vue.component("reverse-proxy-box", {
    props: ["v-reverse-proxy-ref", "v-reverse-proxy-config", "v-is-location", "v-family"],
    data: function () {
        let reverseProxyRef = this.vReverseProxyRef
        if (reverseProxyRef == null) {
            reverseProxyRef = {
                isPrior: false,
                isOn: false,
                reverseProxyId: 0
            }
        }

        let reverseProxyConfig = this.vReverseProxyConfig
        if (reverseProxyConfig == null) {
            reverseProxyConfig = {
                requestPath: "",
                stripPrefix: "",
                requestURI: "",
                requestHost: "",
                requestHostType: 0,
                addHeaders: []
            }
        }

        let forwardHeaders = [
            {
                name: "X-Real-IP",
                isChecked: false
            },
            {
                name: "X-Forwarded-For",
                isChecked: false
            },
            {
                name: "X-Forwarded-By",
                isChecked: false
            },
            {
                name: "X-Forwarded-Host",
                isChecked: false
            },
            {
                name: "X-Forwarded-Proto",
                isChecked: false
            }
        ]
        forwardHeaders.forEach(function (v) {
            v.isChecked = reverseProxyConfig.addHeaders.$contains(v.name)
        })

        return {
            reverseProxyRef: reverseProxyRef,
            reverseProxyConfig: reverseProxyConfig,
            advancedVisible: false,
            family: this.vFamily,
            forwardHeaders: forwardHeaders
        }
    },
    watch: {
        "reverseProxyConfig.requestHostType": function (v) {
            let requestHostType = parseInt(v)
            if (isNaN(requestHostType)) {
                requestHostType = 0
            }
            this.reverseProxyConfig.requestHostType = requestHostType
        }
    },
    methods: {
        isOn: function () {
            return (!this.vIsLocation || this.reverseProxyRef.isPrior) && this.reverseProxyRef.isOn
        },
        changeAdvancedVisible: function (v) {
            this.advancedVisible = v
        },
        changeAddHeader: function () {
            this.reverseProxyConfig.addHeaders = this.forwardHeaders.filter(function (v) {
                return v.isChecked
            }).map(function (v) {
                return v.name
            })
        }
    },
    template: `<div>
	<input type="hidden" name="reverseProxyRefJSON" :value="JSON.stringify(reverseProxyRef)"/>
	<input type="hidden" name="reverseProxyJSON" :value="JSON.stringify(reverseProxyConfig)"/>
	<table class="ui table selectable definition">
		<prior-checkbox :v-config="reverseProxyRef" v-if="vIsLocation"></prior-checkbox>
		<tbody v-show="!vIsLocation || reverseProxyRef.isPrior">
			<tr>
				<td class="title">是否启用反向代理</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="reverseProxyRef.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
			<tr v-show="family == null || family == 'http'">
				<td>回源主机名<em>（Host）</em></td>
				<td>	
					<radio :v-value="0" v-model="reverseProxyConfig.requestHostType">跟随代理服务</radio> &nbsp;
					<radio :v-value="1" v-model="reverseProxyConfig.requestHostType">跟随源站</radio> &nbsp;
					<radio :v-value="2" v-model="reverseProxyConfig.requestHostType">自定义</radio>
					<div v-show="reverseProxyConfig.requestHostType == 2" style="margin-top: 0.8em">
						<input type="text" placeholder="比如example.com" v-model="reverseProxyConfig.requestHost"/>
					</div>
					<p class="comment">请求源站时的Host，用于修改源站接收到的域名
					<span v-if="reverseProxyConfig.requestHostType == 0">，"跟随代理服务"是指源站接收到的域名和当前代理服务保持一致</span>
					<span v-if="reverseProxyConfig.requestHostType == 1">，"跟随源站"是指源站接收到的域名仍然是填写的源站地址中的信息，不随代理服务域名改变而改变</span>					
					<span v-if="reverseProxyConfig.requestHostType == 2">，自定义Host内容中支持请求变量</span>。</p>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-if="isOn()"></more-options-tbody>
		<tbody v-show="isOn() && advancedVisible">
		    <tr v-show="family == null || family == 'http'">
		        <td>自动添加的Header</td>
		        <td>
		            <div>
		                <div style="width: 14em; float: left; margin-bottom: 1em" v-for="header in forwardHeaders" :key="header.name">
		                    <checkbox v-model="header.isChecked" @input="changeAddHeader">{{header.name}}</checkbox>
                        </div>
                        <div style="clear: both"></div>
                    </div>
                    <p class="comment">选中后，会自动向源站请求添加这些Header。</p>
                </td> 
            </tr>
			<tr v-show="family == null || family == 'http'">
				<td>请求URI<em>（RequestURI）</em></td>
				<td>
					<input type="text" placeholder="\${requestURI}" v-model="reverseProxyConfig.requestURI"/>
					<p class="comment">\${requestURI}为完整的请求URI，可以使用类似于"\${requestURI}?arg1=value1&arg2=value2"的形式添加你的参数。</p>
				</td>
			</tr>
			<tr v-show="family == null || family == 'http'">
				<td>去除URL前缀<em>（StripPrefix）</em></td>
				<td>
					<input type="text" v-model="reverseProxyConfig.stripPrefix" placeholder="/PREFIX"/>
					<p class="comment">可以把请求的路径部分前缀去除后再查找文件，比如把 <span class="ui label tiny">/web/app/index.html</span> 去除前缀 <span class="ui label tiny">/web</span> 后就变成 <span class="ui label tiny">/app/index.html</span>。 </p>
				</td>
			</tr>
			<tr>
				<td>是否自动刷新缓存区<em>（AutoFlush）</em></td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="reverseProxyConfig.autoFlush"/>
						<label></label>
					</div>
					<p class="comment">开启后将自动刷新缓冲区数据到客户端，在类似于SSE（server-sent events）等场景下很有用。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`
})