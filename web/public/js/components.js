Vue.component("traffic-map-box",{props:["v-stats","v-is-attack"],mounted:function(){this.render()},data:function(){let i=0;var e=this.vIsAttack,t=(this.vStats.forEach(function(e){var t=parseFloat(e.percent);t>i&&(i=t),e.formattedCountRequests=teaweb.formatCount(e.countRequests)+"次",e.formattedCountAttackRequests=teaweb.formatCount(e.countAttackRequests)+"次"}),i<100&&(i*=1.2),window.innerWidth<512);return{isAttack:e,stats:this.vStats,chart:null,minOpacity:.2,maxPercent:i,selectedCountryName:"",screenIsNarrow:t}},methods:{render:function(){this.chart=teaweb.initChart(document.getElementById("traffic-map-box"));let n=this;this.chart.setOption({backgroundColor:"white",grid:{top:0,bottom:0,left:0,right:0},roam:!1,tooltip:{trigger:"item"},series:[{type:"map",map:"world",zoom:1.3,selectedMode:!1,itemStyle:{areaColor:"#E9F0F9",borderColor:"#DDD"},label:{show:!1,fontSize:"10px",color:"#fff",backgroundColor:"#8B9BD3",padding:[2,2,2,2]},emphasis:{itemStyle:{areaColor:"#8B9BD3",opacity:1},label:{show:!0,fontSize:"10px",color:"#fff",backgroundColor:"#8B9BD3",padding:[2,2,2,2]}},tooltip:{formatter:function(e){let t=e.name,i=null;return n.stats.forEach(function(e){e.name==t&&(i=e)}),null!=i?t+"<br/>流量："+i.formattedBytes+"<br/>流量占比："+i.percent+"%<br/>请求数："+i.formattedCountRequests+"<br/>攻击数："+i.formattedCountAttackRequests:t}},data:this.stats.map(function(e){let t=parseFloat(e.percent)/n.maxPercent,i=3*(t=t<n.minOpacity?n.minOpacity:t);1<i&&(i=1);let s=n.vIsAttack?"#B03A5B":"#276AC6";return{name:e.name,value:e.bytes,percent:parseFloat(e.percent),itemStyle:{areaColor:s,opacity:t},emphasis:{itemStyle:{areaColor:s,opacity:i},label:{show:!0,formatter:function(e){return e.name}}},label:{show:!1,formatter:function(e){return e.name==n.selectedCountryName?e.name:""},fontSize:"10px",color:"#fff",backgroundColor:"#8B9BD3",padding:[2,2,2,2]}}}),nameMap:window.WorldCountriesMap}]}),this.chart.resize()},selectCountry:function(s){if(null!=this.chart){let e=this.chart.getOption(),i=this;e.series[0].data.forEach(function(e){let t=e.percent/i.maxPercent;if(t<i.minOpacity&&(t=i.minOpacity),e.name==s){if(e.isSelected)return e.itemStyle.opacity=t,e.isSelected=!1,e.label.show=!1,void(i.selectedCountryName="");e.isSelected=!0,i.selectedCountryName=s,(t=1<(t*=3)?1:t)<.5&&(t=.5),e.itemStyle.opacity=t,e.label.show=!0}else e.itemStyle.opacity=t,e.isSelected=!1,e.label.show=!1}),this.chart.setOption(e)}},select:function(e){this.selectCountry(e.countryName)}},template:`<div>
	<table style="width: 100%; border: 0; padding: 0; margin: 0">
		<tbody>
       	<tr>
           <td>
               <div class="traffic-map-box" id="traffic-map-box"></div>
           </td>
           <td style="width: 14em" v-if="!screenIsNarrow">
           		<traffic-map-box-table :v-stats="stats" :v-is-attack="isAttack" @select="select"></traffic-map-box-table>
           </td>
       </tr>
       </tbody>
       <tbody v-if="screenIsNarrow">
		   <tr>
				<td colspan="2">
					<traffic-map-box-table :v-stats="stats" :v-is-attack="isAttack" :v-screen-is-narrow="true" @select="select"></traffic-map-box-table>
				</td>
			</tr>
		</tbody>
   </table>
</div>`}),Vue.component("traffic-map-box-table",{props:["v-stats","v-is-attack","v-screen-is-narrow"],data:function(){return{stats:this.vStats,isAttack:this.vIsAttack}},methods:{select:function(e){this.$emit("select",{countryName:e})}},template:`<div style="overflow-y: auto" :style="{'max-height':vScreenIsNarrow ? 'auto' : '16em'}" class="narrow-scrollbar">
	   <table class="ui table selectable">
		  <thead>
			<tr>
				<th colspan="2">国家/地区排行&nbsp; <tip-icon content="只有开启了统计的服务才会有记录。"></tip-icon></th>
			</tr>
		  </thead>
		   <tbody v-if="stats.length == 0">
			   <tr>
				   <td colspan="2">暂无数据</td>
			   </tr>
		   </tbody>
		   <tbody>
			   <tr v-for="(stat, index) in stats.slice(0, 10)">
				   <td @click.prevent="select(stat.name)" style="cursor: pointer" colspan="2">
					   <div class="ui progress bar" :class="{red: vIsAttack, blue:!vIsAttack}" style="margin-bottom: 0.3em">
						   <div class="bar" style="min-width: 0; height: 4px;" :style="{width: stat.percent + '%'}"></div>
					   </div>
					  <div>{{stat.name}}</div> 
					   <div><span class="grey">{{stat.percent}}% </span>
					   <span class="small grey" v-if="isAttack">{{stat.formattedCountAttackRequests}}</span>
					   <span class="small grey" v-if="!isAttack">（{{stat.formattedBytes}}）</span></div>
				   </td>
			   </tr>
		   </tbody>
	   </table>
   </div>`}),Vue.component("ddos-protection-ports-config-box",{props:["v-ports"],data:function(){let e=this.vPorts;return{ports:e=null==e?[]:e,isAdding:!1,addingPort:{port:"",description:""}}},methods:{add:function(){this.isAdding=!0;let e=this;setTimeout(function(){e.$refs.addingPortInput.focus()})},confirm:function(){var e=this.addingPort.port;if(0==e.length)this.warn("请输入端口号");else if(/^\d+$/.test(e)){let i=parseInt(e,10);if(i<=0)this.warn("请输入正确的端口号");else if(65535<i)this.warn("请输入正确的端口号");else{let t=!1;this.ports.forEach(function(e){e.port==i&&(t=!0)}),t?this.warn("端口号已经存在"):(this.ports.push({port:i,description:this.addingPort.description}),this.notifyChange(),this.cancel())}}else this.warn("请输入正确的端口号")},cancel:function(){this.isAdding=!1,this.addingPort={port:"",description:""}},remove:function(e){this.ports.$remove(e),this.notifyChange()},warn:function(e){let t=this;teaweb.warn(e,function(){t.$refs.addingPortInput.focus()})},notifyChange:function(){this.$emit("change",this.ports)}},template:`<div>
	<div v-if="ports.length > 0">
		<div class="ui label basic tiny" v-for="(portConfig, index) in ports">
			{{portConfig.port}} <span class="grey small" v-if="portConfig.description.length > 0">（{{portConfig.description}}）</span> <a href="" @click.prevent="remove(index)" title="删除"><i class="icon remove"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding">
		<div class="ui fields inline">
			<div class="ui field">
				<div class="ui input left labeled">
					<span class="ui label">端口</span>
					<input type="text" v-model="addingPort.port" ref="addingPortInput" maxlength="5" size="5" placeholder="端口号" @keyup.enter="confirm" @keypress.enter.prevent="1"/>
				</div>
			</div>
			<div class="ui field">
				<div class="ui input left labeled">
					<span class="ui label">备注</span>
					<input type="text" v-model="addingPort.description" maxlength="12" size="12" placeholder="备注（可选）" @keyup.enter="confirm" @keypress.enter.prevent="1"/>
				</div>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>
				&nbsp;<a href="" @click.prevent="cancel()">取消</a>
			</div>
		</div>
	</div>
	<div v-if="!isAdding">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`}),Vue.component("node-clusters-labels",{props:["v-primary-cluster","v-secondary-clusters","size"],data:function(){var e=this.vPrimaryCluster;let t=this.vSecondaryClusters,i=(null==t&&(t=[]),this.size);return null==i&&(i="small"),{cluster:e,secondaryClusters:t,labelSize:i}},template:`<div>
	<a v-if="cluster != null" :href="'/clusters/cluster?clusterId=' + cluster.id" title="主集群" style="margin-bottom: 0.3em;">
		<span class="ui label basic grey" :class="labelSize" v-if="labelSize != 'tiny'">{{cluster.name}}</span>
		<grey-label v-if="labelSize == 'tiny'">{{cluster.name}}</grey-label>
	</a>
	<a v-for="c in secondaryClusters" :href="'/clusters/cluster?clusterId=' + c.id" :class="labelSize" title="从集群">
		<span class="ui label basic grey" :class="labelSize" v-if="labelSize != 'tiny'">{{c.name}}</span>
		<grey-label v-if="labelSize == 'tiny'">{{c.name}}</grey-label>
	</a>
</div>`}),Vue.component("cluster-selector",{props:["v-cluster-id"],mounted:function(){let t=this;Tea.action("/clusters/options").post().success(function(e){t.clusters=e.data.clusters})},data:function(){let e=this.vClusterId;return{clusters:[],clusterId:e=null==e?0:e}},template:`<div>
	<select class="ui dropdown" style="max-width: 10em" name="clusterId" v-model="clusterId">
		<option value="0">[选择集群]</option>
		<option v-for="cluster in clusters" :value="cluster.id">{{cluster.name}}</option>
	</select>
</div>`}),Vue.component("node-ddos-protection-config-box",{props:["v-ddos-protection-config","v-default-configs","v-is-node","v-cluster-is-on"],data:function(){let e=this.vDdosProtectionConfig;return null==(e=null==e?{tcp:{isPrior:!1,isOn:!1,maxConnections:0,maxConnectionsPerIP:0,newConnectionsRate:0,allowIPList:[],ports:[]}}:e).tcp&&(e.tcp={isPrior:!1,isOn:!1,maxConnections:0,maxConnectionsPerIP:0,newConnectionsRate:0,allowIPList:[],ports:[]}),{config:e,defaultConfigs:this.vDefaultConfigs,isNode:this.vIsNode,isAddingPort:!1}},methods:{changeTCPPorts:function(e){this.config.tcp.ports=e},changeTCPAllowIPList:function(e){this.config.tcp.allowIPList=e}},template:`<div>
 <input type="hidden" name="ddosProtectionJSON" :value="JSON.stringify(config)"/>

 <p class="comment">功能说明：此功能为<strong>试验性质</strong>，目前仅能防御简单的DDoS攻击，试验期间建议仅在被攻击时启用，仅支持已安装<code-label>nftables v0.9</code-label>以上的Linux系统。<pro-warning-label></pro-warning-label></p>

 <div class="ui message" v-if="vClusterIsOn">当前节点所在集群已设置DDoS防护。</div>

 <h4>TCP设置</h4>
 <table class="ui table definition selectable">
 	<prior-checkbox :v-config="config.tcp" v-if="isNode"></prior-checkbox>
 	<tbody v-show="config.tcp.isPrior || !isNode">
		<tr>
			<td class="title">启用</td>
			<td>
				<checkbox v-model="config.tcp.isOn"></checkbox>
			</td>
		</tr>
	</tbody>
	<tbody v-show="config.tcp.isOn && (config.tcp.isPrior || !isNode)">
		<tr>
			<td class="title">单节点TCP最大连接数</td>
			<td>
				<digit-input name="tcpMaxConnections" v-model="config.tcp.maxConnections" maxlength="6" size="6" style="width: 6em"></digit-input>
				<p class="comment">单个节点可以接受的TCP最大连接数。如果为0，则默认为{{defaultConfigs.tcpMaxConnections}}。</p>
			</td>
		</tr>
		<tr>
			<td>单IP TCP最大连接数</td>
			<td>
				<digit-input name="tcpMaxConnectionsPerIP" v-model="config.tcp.maxConnectionsPerIP" maxlength="6" size="6" style="width: 6em"></digit-input>
				<p class="comment">单个IP可以连接到节点的TCP最大连接数。如果为0，则默认为{{defaultConfigs.tcpMaxConnectionsPerIP}}；最小值为{{defaultConfigs.tcpMinConnectionsPerIP}}。</p>
			</td>
		</tr>
		<tr>
			<td>单IP TCP新连接速率</td>
			<td>
				<div class="ui input right labeled">
					<digit-input name="tcpNewConnectionsRate" v-model="config.tcp.newConnectionsRate" maxlength="6" size="6" style="width: 6em" :min="defaultConfigs.tcpNewConnectionsMinRate"></digit-input>
					<span class="ui label">个新连接/每分钟</span>
				</div>
				<p class="comment">单个IP可以创建TCP新连接的速率。如果为0，则默认为{{defaultConfigs.tcpNewConnectionsRate}}；最小值为{{defaultConfigs.tcpNewConnectionsMinRate}}。</p>
			</td>
		</tr>
		<tr>
			<td>TCP端口列表</td>
			<td>
				<ddos-protection-ports-config-box :v-ports="config.tcp.ports" @change="changeTCPPorts"></ddos-protection-ports-config-box>
				<p class="comment">默认为80和443两个端口。</p>
			</td>
		</tr>
		<tr>
			<td>IP白名单</td>
			<td>
				<ddos-protection-ip-list-config-box :v-ip-list="config.tcp.allowIPList" @change="changeTCPAllowIPList"></ddos-protection-ip-list-config-box>
				<p class="comment">在白名单中的IP不受当前设置的限制。</p>
			</td>
		</tr>
	</tbody>
</table>
<div class="margin"></div>
</div>`}),Vue.component("ddos-protection-ip-list-config-box",{props:["v-ip-list"],data:function(){let e=this.vIpList;return{list:e=null==e?[]:e,isAdding:!1,addingIP:{ip:"",description:""}}},methods:{add:function(){this.isAdding=!0;let e=this;setTimeout(function(){e.$refs.addingIPInput.focus()})},confirm:function(){let i=this.addingIP.ip;if(0==i.length)this.warn("请输入IP");else{let t=!1;if(this.list.forEach(function(e){e.ip==i&&(t=!0)}),t)this.warn("IP '"+i+"'已经存在");else{let e=this;Tea.Vue.$post("/ui/validateIPs").params({ips:[i]}).success(function(){e.list.push({ip:i,description:e.addingIP.description}),e.notifyChange(),e.cancel()}).fail(function(){e.warn("请输入正确的IP")})}}},cancel:function(){this.isAdding=!1,this.addingIP={ip:"",description:""}},remove:function(e){this.list.$remove(e),this.notifyChange()},warn:function(e){let t=this;teaweb.warn(e,function(){t.$refs.addingIPInput.focus()})},notifyChange:function(){this.$emit("change",this.list)}},template:`<div>
	<div v-if="list.length > 0">
		<div class="ui label basic tiny" v-for="(ipConfig, index) in list">
			{{ipConfig.ip}} <span class="grey small" v-if="ipConfig.description.length > 0">（{{ipConfig.description}}）</span> <a href="" @click.prevent="remove(index)" title="删除"><i class="icon remove"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding">
		<div class="ui fields inline">
			<div class="ui field">
				<div class="ui input left labeled">
					<span class="ui label">IP</span>
					<input type="text" v-model="addingIP.ip" ref="addingIPInput" maxlength="40" size="20" placeholder="IP" @keyup.enter="confirm" @keypress.enter.prevent="1"/>
				</div>
			</div>
			<div class="ui field">
				<div class="ui input left labeled">
					<span class="ui label">备注</span>
					<input type="text" v-model="addingIP.description" maxlength="10" size="10" placeholder="备注（可选）" @keyup.enter="confirm" @keypress.enter.prevent="1"/>
				</div>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>
				&nbsp;<a href="" @click.prevent="cancel()">取消</a>
			</div>
		</div>
	</div>
	<div v-if="!isAdding">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`}),Vue.component("node-cluster-combo-box",{props:["v-cluster-id"],data:function(){let t=this;return Tea.action("/clusters/options").post().success(function(e){t.clusters=e.data.clusters}),{clusters:[]}},methods:{change:function(e){null==e?this.$emit("change",0):this.$emit("change",e.value)}},template:`<div v-if="clusters.length > 0" style="min-width: 10.4em">
	<combo-box title="集群" placeholder="集群名称" :v-items="clusters" name="clusterId" :v-value="vClusterId" @change="change"></combo-box>
</div>`}),Vue.component("node-clusters-selector",{props:["v-primary-cluster","v-secondary-clusters"],data:function(){var e=this.vPrimaryCluster;let t=this.vSecondaryClusters;return null==t&&(t=[]),{primaryClusterId:null==e?0:e.id,secondaryClusterIds:t.map(function(e){return e.id}),primaryCluster:e,secondaryClusters:t}},methods:{addPrimary:function(){let t=this,e=[this.primaryClusterId].concat(this.secondaryClusterIds);teaweb.popup("/clusters/selectPopup?selectedClusterIds="+e.join(",")+"&mode=single",{height:"30em",width:"50em",callback:function(e){null!=e.data.cluster&&(t.primaryCluster=e.data.cluster,t.primaryClusterId=t.primaryCluster.id,t.notifyChange())}})},removePrimary:function(){this.primaryClusterId=0,this.primaryCluster=null,this.notifyChange()},addSecondary:function(){let t=this,e=[this.primaryClusterId].concat(this.secondaryClusterIds);teaweb.popup("/clusters/selectPopup?selectedClusterIds="+e.join(",")+"&mode=multiple",{height:"30em",width:"50em",callback:function(e){null!=e.data.cluster&&(t.secondaryClusterIds.push(e.data.cluster.id),t.secondaryClusters.push(e.data.cluster),t.notifyChange())}})},removeSecondary:function(e){this.secondaryClusterIds.$remove(e),this.secondaryClusters.$remove(e),this.notifyChange()},notifyChange:function(){this.$emit("change",{clusterId:this.primaryClusterId})}},template:`<div>
	<input type="hidden" name="primaryClusterId" :value="primaryClusterId"/>
	<input type="hidden" name="secondaryClusterIds" :value="JSON.stringify(secondaryClusterIds)"/>
	<table class="ui table">
		<tr>
			<td class="title">主集群</td>
			<td>
				<div v-if="primaryCluster != null">
					<div class="ui label basic small">{{primaryCluster.name}} &nbsp; <a href="" title="删除" @click.prevent="removePrimary"><i class="icon remove small"></i></a> </div>
				</div>
				<div style="margin-top: 0.6em" v-if="primaryClusterId == 0">
					<button class="ui button tiny" type="button" @click.prevent="addPrimary">+</button>
				</div>
				<p class="comment">多个集群配置有冲突时，优先使用主集群配置。</p>
			</td>
		</tr>
		<tr>
			<td>从集群</td>
			<td>
				<div v-if="secondaryClusters.length > 0">
					<div class="ui label basic small" v-for="(cluster, index) in secondaryClusters"><span class="grey">{{cluster.name}}</span> &nbsp; <a href="" title="删除" @click.prevent="removeSecondary(index)"><i class="icon remove small"></i></a> </div>
				</div>
				<div style="margin-top: 0.6em">
					<button class="ui button tiny" type="button" @click.prevent="addSecondary">+</button>
				</div>
			</td>
		</tr>
	</table>
</div>`}),Vue.component("message-media-selector",{props:["v-media-type"],mounted:function(){let i=this;Tea.action("/admins/recipients/mediaOptions").post().success(function(e){i.medias=e.data.medias,0<i.mediaType.length&&(null!=(e=i.medias.$find(function(e,t){return t.type==i.mediaType}))&&(i.description=e.description))})},data:function(){let e=this.vMediaType;return{medias:[],description:"",mediaType:e=null==e?"":e}},watch:{mediaType:function(i){var e=this.medias.$find(function(e,t){return t.type==i});this.description=null==e?"":e.description,this.$emit("change",e)}},template:`<div>
    <select class="ui dropdown auto-width" name="mediaType" v-model="mediaType">
        <option value="">[选择媒介类型]</option>
        <option v-for="media in medias" :value="media.type">{{media.name}}</option>
    </select>
    <p class="comment" v-html="description"></p>
</div>`}),Vue.component("message-receivers-box",{props:["v-node-cluster-id"],mounted:function(){let t=this;Tea.action("/clusters/cluster/settings/message/selectedReceivers").params({clusterId:this.clusterId}).post().success(function(e){t.receivers=e.data.receivers})},data:function(){let e=this.vNodeClusterId;return{clusterId:e=null==e?0:e,receivers:[]}},methods:{addReceiver:function(){let t=this,i=[],s=[];this.receivers.forEach(function(e){"recipient"==e.type?i.push(e.id.toString()):"group"==e.type&&s.push(e.id.toString())}),teaweb.popup("/clusters/cluster/settings/message/selectReceiverPopup?recipientIds="+i.join(",")+"&groupIds="+s.join(","),{callback:function(e){t.receivers.push(e.data)}})},removeReceiver:function(e){this.receivers.$remove(e)}},template:`<div>
        <input type="hidden" name="receiversJSON" :value="JSON.stringify(receivers)"/>           
        <div v-if="receivers.length > 0">
            <div v-for="(receiver, index) in receivers" class="ui label basic small">
               <span v-if="receiver.type == 'group'">分组：</span>{{receiver.name}} <span class="grey small" v-if="receiver.subName != null && receiver.subName.length > 0">({{receiver.subName}})</span> &nbsp; <a href="" title="删除" @click.prevent="removeReceiver(index)"><i class="icon remove"></i></a>
            </div>
             <div class="ui divider"></div>
        </div>
      <button type="button" class="ui button tiny" @click.prevent="addReceiver">+</button>
</div>`}),Vue.component("message-recipient-group-selector",{props:["v-groups"],data:function(){let e=this.vGroups,t=[];return 0<(e=null==e?[]:e).length&&(t=e.map(function(e){return e.id.toString()}).join(",")),{groups:e,groupIds:t}},methods:{addGroup:function(){let t=this;teaweb.popup("/admins/recipients/groups/selectPopup?groupIds="+this.groupIds,{callback:function(e){t.groups.push(e.data.group),t.update()}})},removeGroup:function(e){this.groups.$remove(e),this.update()},update:function(){let t=[];0<this.groups.length&&this.groups.forEach(function(e){t.push(e.id)}),this.groupIds=t.join(",")}},template:`<div>
    <input type="hidden" name="groupIds" :value="groupIds"/>
    <div v-if="groups.length > 0">
        <div>
            <div v-for="(group, index) in groups" class="ui label small basic">
                {{group.name}} &nbsp; <a href="" title="删除" @click.prevent="removeGroup(index)"><i class="icon remove"></i></a>
            </div>
        </div>
        <div class="ui divider"></div>
    </div>   
    <button class="ui button tiny" type="button" @click.prevent="addGroup()">+</button>
</div>`}),Vue.component("message-media-instance-selector",{props:["v-instance-id"],mounted:function(){let i=this;Tea.action("/admins/recipients/instances/options").post().success(function(e){i.instances=e.data.instances,0<i.instanceId&&(null!=(e=i.instances.$find(function(e,t){return t.id==i.instanceId}))&&(i.description=e.description,i.update(e.id)))})},data:function(){let e=this.vInstanceId;return{instances:[],description:"",instanceId:e=null==e?0:e}},watch:{instanceId:function(e){this.update(e)}},methods:{update:function(i){var e=this.instances.$find(function(e,t){return t.id==i});this.description=null==e?"":e.description,this.$emit("change",e)}},template:`<div>
    <select class="ui dropdown auto-width" name="instanceId" v-model="instanceId">
        <option value="0">[选择媒介]</option>
        <option v-for="instance in instances" :value="instance.id">{{instance.name}} ({{instance.media.name}})</option>
    </select>
    <p class="comment" v-html="description"></p>
</div>`}),Vue.component("message-row",{props:["v-message","v-can-close"],data:function(){var e=this.vMessage.params;let t=null;return null!=e&&0<e.length&&(t=JSON.parse(e)),{message:this.vMessage,params:t,isClosing:!1}},methods:{viewCert:function(e){teaweb.popup("/servers/certs/certPopup?certId="+e,{height:"28em",width:"48em"})},readMessage:function(e){let t=this;Tea.action("/messages/readPage").params({messageIds:[e]}).post().success(function(){null!=window.parent.Tea&&null!=window.parent.Tea.Vue&&window.parent.Tea.Vue.checkMessagesOnce(),t.vCanClose&&"undefined"!=typeof NotifyPopup?(t.isClosing=!0,setTimeout(function(){NotifyPopup({})},1e3)):teaweb.reload()})}},template:`<div>
<table class="ui table selectable" v-if="!isClosing">
	<tr :class="{error: message.level == 'error', positive: message.level == 'success', warning: message.level == 'warning'}">
		<td style="position: relative">
			<strong>{{message.datetime}}</strong>
			<span v-if="message.cluster != null && message.cluster.id != null">
				<span> | </span>
				<a :href="'/clusters/cluster?clusterId=' + message.cluster.id" target="_top" v-if="message.role == 'node'">集群：{{message.cluster.name}}</a>
				<a :href="'/ns/clusters/cluster?clusterId=' + message.cluster.id" target="_top" v-if="message.role == 'dns'">DNS集群：{{message.cluster.name}}</a>
			</span>
			<span v-if="message.node != null && message.node.id != null">
				<span> | </span>
				<a :href="'/clusters/cluster/node?clusterId=' + message.cluster.id + '&nodeId=' + message.node.id" target="_top" v-if="message.role == 'node'">节点：{{message.node.name}}</a>
				<a :href="'/ns/clusters/cluster/node?clusterId=' + message.cluster.id + '&nodeId=' + message.node.id" target="_top" v-if="message.role == 'dns'">DNS节点：{{message.node.name}}</a>
			</span>
			<a href=""  style="position: absolute; right: 1em" @click.prevent="readMessage(message.id)" title="标为已读"><i class="icon check"></i></a>
		</td>
	</tr>
	<tr :class="{error: message.level == 'error', positive: message.level == 'success', warning: message.level == 'warning'}">
		<td>
			{{message.body}}
			
			<!-- 健康检查 -->
			<div v-if="message.type == 'HealthCheckFailed'" style="margin-top: 0.8em">
				<a :href="'/clusters/cluster/node?clusterId=' + message.cluster.id + '&nodeId=' + param.node.id" v-for="param in params" class="ui label small basic" style="margin-bottom: 0.5em" target="_top">{{param.node.name}}: {{param.error}}</a>
			</div>
			
			<!-- 集群DNS设置 -->
			<div v-if="message.type == 'ClusterDNSSyncFailed'" style="margin-top: 0.8em">
				<a :href="'/dns/clusters/cluster?clusterId=' + message.cluster.id" target="_top">查看问题 &raquo;</a>
			</div>
			
			<!-- 证书即将过期 -->
			<div v-if="message.type == 'SSLCertExpiring'" style="margin-top: 0.8em">
				<a href="" @click.prevent="viewCert(params.certId)" target="_top">查看证书</a> &nbsp;|&nbsp; <a :href="'/servers/certs/acme'" v-if="params != null && params.acmeTaskId > 0" target="_top">查看任务&raquo;</a>
			</div>
			
			<!-- 证书续期成功 -->
			<div v-if="message.type == 'SSLCertACMETaskSuccess'" style="margin-top: 0.8em">
				<a href="" @click.prevent="viewCert(params.certId)" target="_top">查看证书</a> &nbsp;|&nbsp; <a :href="'/servers/certs/acme'" v-if="params != null && params.acmeTaskId > 0" target="_top">查看任务&raquo;</a>
			</div>
			
			<!-- 证书续期失败 -->
			<div v-if="message.type == 'SSLCertACMETaskFailed'" style="margin-top: 0.8em">
				<a href="" @click.prevent="viewCert(params.certId)" target="_top">查看证书</a> &nbsp;|&nbsp; <a :href="'/servers/certs/acme'" v-if="params != null && params.acmeTaskId > 0" target="_top">查看任务&raquo;</a>
			</div>
			
			<!-- 网站域名审核 -->
			<div v-if="message.type == 'serverNamesRequireAuditing'" style="margin-top: 0.8em">
				<a :href="'/servers/server/settings/serverNames?serverId=' + params.serverId" target="_top">去审核</a></a>
			</div>
		</td>
	</tr>
</table>
<div class="margin"></div>
</div>`}),Vue.component("ns-routes-selector",{props:["v-routes"],mounted:function(){let t=this;Tea.action("/ns/routes/options").post().success(function(e){t.routes=e.data.routes})},data:function(){let e=this.vRoutes;return{routeCode:"default",routes:[],isAdding:!1,routeType:"default",selectedRoutes:e=null==e?[]:e}},watch:{routeType:function(t){this.routeCode="";let i=this;this.routes.forEach(function(e){e.type==t&&0==i.routeCode.length&&(i.routeCode=e.code)})}},methods:{add:function(){this.isAdding=!0,this.routeType="default",this.routeCode="default"},cancel:function(){this.isAdding=!1},confirm:function(){if(0!=this.routeCode.length){let t=this;this.routes.forEach(function(e){e.code==t.routeCode&&t.selectedRoutes.push(e)}),this.cancel()}},remove:function(e){this.selectedRoutes.$remove(e)}},template:`<div>
	<div>
		<div class="ui label basic text small" v-for="(route, index) in selectedRoutes" style="margin-bottom: 0.3em">
			<input type="hidden" name="routeCodes" :value="route.code"/>
			{{route.name}} &nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding" style="margin-bottom: 1em">
		<div class="ui fields inline">
			<div class="ui field">
				<select class="ui dropdown" v-model="routeType">
					<option value="default">[默认线路]</option>
					<option value="user">自定义线路</option>
					<option value="isp">运营商</option>
					<option value="china">中国省市</option>
					<option value="world">全球国家地区</option>
				</select>
			</div>
			
			<div class="ui field">
				<select class="ui dropdown" v-model="routeCode" style="width: 10em">
					<option v-for="route in routes" :value="route.code" v-if="route.type == routeType">{{route.name}}</option>
				</select>
			</div>
			
			<div class="ui field">
				<button type="button" class="ui button tiny" @click.prevent="confirm">确定</button>
				&nbsp; <a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
	</div>
	<button class="ui button tiny" type="button" @click.prevent="add">+</button>
</div>`}),Vue.component("ns-recursion-config-box",{props:["v-recursion-config"],data:function(){let e=this.vRecursionConfig;return null==(e=null==e?{isOn:!1,hosts:[],allowDomains:[],denyDomains:[],useLocalHosts:!1}:e).hosts&&(e.hosts=[]),null==e.allowDomains&&(e.allowDomains=[]),null==e.denyDomains&&(e.denyDomains=[]),{config:e,hostIsAdding:!1,host:"",updatingHost:null}},methods:{changeHosts:function(e){this.config.hosts=e},changeAllowDomains:function(e){this.config.allowDomains=e},changeDenyDomains:function(e){this.config.denyDomains=e},removeHost:function(e){this.config.hosts.$remove(e)},addHost:function(){var t;this.updatingHost=null,this.host="",this.hostIsAdding=!this.hostIsAdding,this.hostIsAdding&&(t=this,setTimeout(function(){let e=t.$refs.hostRef;null!=e&&e.focus()},200))},updateHost:function(e){var t;this.updatingHost=e,this.host=e.host,this.hostIsAdding=!this.hostIsAdding,this.hostIsAdding&&(t=this,setTimeout(function(){let e=t.$refs.hostRef;null!=e&&e.focus()},200))},confirmHost:function(){0==this.host.length?teaweb.warn("请输入DNS地址"):(this.hostIsAdding=!1,null==this.updatingHost?this.config.hosts.push({host:this.host}):this.updatingHost.host=this.host)},cancelHost:function(){this.hostIsAdding=!1}},template:`<div>
	<input type="hidden" name="recursionJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<tbody>
			<tr>
				<td class="title">是否启用</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" name="isOn" value="1" v-model="config.isOn"/>
						<label></label>
					</div>
					<p class="comment">启用后，如果找不到某个域名的解析记录，则向上一级DNS查找。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="config.isOn">
			<tr>
				<td>从节点本机读取<br/>上级DNS主机</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" name="useLocalHosts" value="1" v-model="config.useLocalHosts"/>
						<label></label>
					</div>
					<p class="comment">选中后，节点会试图从<code-label>/etc/resolv.conf</code-label>文件中读取DNS配置。 </p>
				</td>
			</tr>
			<tr v-show="!config.useLocalHosts">
				<td>上级DNS主机地址 *</td>
				<td>
					<div v-if="config.hosts.length > 0">
						<div v-for="(host, index) in config.hosts" class="ui label tiny basic">
							{{host.host}} &nbsp;
							<a href="" title="修改" @click.prevent="updateHost(host)"><i class="icon pencil tiny"></i></a>
							<a href="" title="删除" @click.prevent="removeHost(index)"><i class="icon remove small"></i></a>
						</div>
						<div class="ui divider"></div>
					</div>
					<div v-if="hostIsAdding">
						<div class="ui fields inline">
							<div class="ui field">
								<input type="text" placeholder="DNS主机地址" v-model="host" ref="hostRef" @keyup.enter="confirmHost" @keypress.enter.prevent="1"/>
							</div>
							<div class="ui field">
								<button class="ui button tiny" type="button" @click.prevent="confirmHost">确认</button> &nbsp; <a href="" title="取消" @click.prevent="cancelHost"><i class="icon remove small"></i></a>
							</div>
						</div>
					</div>
					<div style="margin-top: 0.5em">
						<button type="button" class="ui button tiny" @click.prevent="addHost">+</button>
					</div>
				</td>
			</tr>
			<tr>
				<td>允许的域名</td>
				<td><values-box name="allowDomains" :values="config.allowDomains" @change="changeAllowDomains"></values-box>
					<p class="comment">支持星号通配符，比如<code-label>*.example.org</code-label>。</p>
				</td>
			</tr>
			<tr>
				<td>不允许的域名</td>
				<td>
					<values-box name="denyDomains" :values="config.denyDomains" @change="changeDenyDomains"></values-box>
					<p class="comment">支持星号通配符，比如<code-label>*.example.org</code-label>。优先级比允许的域名高。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`}),Vue.component("ns-access-log-ref-box",{props:["v-access-log-ref","v-is-parent"],data:function(){let e=this.vAccessLogRef;return void 0===(e=null==e?{isOn:!1,isPrior:!1,logMissingDomains:!1}:e).logMissingDomains&&(e.logMissingDomains=!1),{config:e}},template:`<div>
	<input type="hidden" name="accessLogJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="config" v-if="!vIsParent"></prior-checkbox>
		<tbody v-show="vIsParent || config.isPrior">
			<tr>
				<td class="title">是否启用</td>
				<td>
					<checkbox name="isOn" value="1" v-model="config.isOn"></checkbox>
				</td>
			</tr>
			<tr>
				<td>记录所有访问</td>
				<td>
					<checkbox name="logMissingDomains" value="1" v-model="config.logMissingDomains"></checkbox>
					<p class="comment">包括对没有在系统里创建的域名访问。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`}),Vue.component("ns-route-ranges-box",{props:["v-ranges"],data:function(){let e=this.vRanges;return{ranges:e=null==e?[]:e,isAdding:!1,ipRangeFrom:"",ipRangeTo:""}},methods:{add:function(){this.isAdding=!0;let e=this;setTimeout(function(){e.$refs.ipRangeFrom.focus()},100)},remove:function(e){this.ranges.$remove(e)},cancelIPRange:function(){this.isAdding=!1,this.ipRangeFrom="",this.ipRangeTo=""},confirmIPRange:function(){let e=this;this.ipRangeFrom=this.ipRangeFrom.trim(),this.validateIP(this.ipRangeFrom)?(this.ipRangeTo=this.ipRangeTo.trim(),this.validateIP(this.ipRangeTo)?(this.ranges.push({type:"ipRange",params:{ipFrom:this.ipRangeFrom,ipTo:this.ipRangeTo}}),this.cancelIPRange()):teaweb.warn("结束IP填写错误",function(){e.$refs.ipRangeTo.focus()})):teaweb.warn("开始IP填写错误",function(){e.$refs.ipRangeFrom.focus()})},validateIP:function(e){if(!e.match(/^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$/))return!1;let t=e.split("."),i=!0;return t.forEach(function(e){255<parseInt(e)&&(i=!1)}),i}},template:`<div>
	<input type="hidden" name="rangesJSON" :value="JSON.stringify(ranges)"/>
	<div v-if="ranges.length > 0">
		<div class="ui label tiny basic" v-for="(range, index) in ranges">
			<span v-if="range.type == 'ipRange'">IP范围：</span>
			{{range.params.ipFrom}} - {{range.params.ipTo}} &nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	
	<!-- IP 范围 -->
	<div style="margin-bottom: 1em" v-show="isAdding">
		<div class="ui fields inline">
			<div class="ui field">
				<input type="text" placeholder="开始IP" maxlength="15" size="15" v-model="ipRangeFrom" ref="ipRangeFrom"  @keyup.enter="confirmIPRange" @keypress.enter.prevent="1"/>
			</div>
			<div class="ui field">-</div>
			<div class="ui field">
				<input type="text" placeholder="结束IP" maxlength="15" size="15" v-model="ipRangeTo" ref="ipRangeTo" @keyup.enter="confirmIPRange" @keypress.enter.prevent="1"/>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirmIPRange">确定</button> &nbsp;
				<a href="" @click.prevent="cancelIPRange" title="取消"><i class="icon remove small"></i></a>
			</div>
		</div>
	</div>
	
	<button class="ui button tiny" type="button" @click.prevent="add">+</button>
</div>`}),Vue.component("ns-route-selector",{props:["v-route-code"],mounted:function(){let t=this;Tea.action("/ns/routes/options").post().success(function(e){t.routes=e.data.routes})},data:function(){let e=this.vRouteCode;return{routeCode:e=null==e?"":e,routes:[]}},template:`<div>
	<div v-if="routes.length > 0">
		<select class="ui dropdown" name="routeCode" v-model="routeCode">
			<option value="">[线路]</option>
			<option v-for="route in routes" :value="route.code">{{route.name}}</option>
		</select>
	</div>
</div>`}),Vue.component("ns-user-selector",{mounted:function(){let t=this;Tea.action("/ns/users/options").post().success(function(e){t.users=e.data.users})},props:["v-user-id"],data:function(){let e=this.vUserId;return{users:[],userId:e=null==e?0:e}},template:`<div>
	<select class="ui dropdown auto-width" name="userId" v-model="userId">
		<option value="0">[选择用户]</option>
		<option v-for="user in users" :value="user.id">{{user.fullname}} ({{user.username}})</option>
	</select>
</div>`}),Vue.component("ns-access-log-box",{props:["v-access-log","v-keyword"],data:function(){return{accessLog:this.vAccessLog}},methods:{showLog:function(){let e=this;var t=this.accessLog.requestId;this.$parent.$children.forEach(function(e){null!=e.deselect&&e.deselect()}),this.select(),teaweb.popup("/ns/clusters/accessLogs/viewPopup?requestId="+t,{width:"50em",height:"24em",onClose:function(){e.deselect()}})},select:function(){this.$refs.box.parentNode.style.cssText="background: rgba(0, 0, 0, 0.1)"},deselect:function(){this.$refs.box.parentNode.style.cssText=""}},template:`<div class="access-log-row" :style="{'color': (!accessLog.isRecursive && (accessLog.nsRecordId == null || accessLog.nsRecordId == 0) || (accessLog.isRecursive && accessLog.recordValue != null && accessLog.recordValue.length == 0)) ? '#dc143c' : ''}" ref="box">
	<span v-if="accessLog.region != null && accessLog.region.length > 0" class="grey">[{{accessLog.region}}]</span> <keyword :v-word="vKeyword">{{accessLog.remoteAddr}}</keyword> [{{accessLog.timeLocal}}] [{{accessLog.networking}}] <em>{{accessLog.questionType}} <keyword :v-word="vKeyword">{{accessLog.questionName}}</keyword></em> -&gt; <em>{{accessLog.recordType}} <keyword :v-word="vKeyword">{{accessLog.recordValue}}</keyword></em><!-- &nbsp; <a href="" @click.prevent="showLog" title="查看详情"><i class="icon expand"></i></a>-->
	<div v-if="(accessLog.nsRoutes != null && accessLog.nsRoutes.length > 0) || accessLog.isRecursive" style="margin-top: 0.3em">
		<span class="ui label tiny basic grey" v-for="route in accessLog.nsRoutes">线路: {{route.name}}</span>
		<span class="ui label tiny basic grey" v-if="accessLog.isRecursive">递归DNS</span>
	</div>
	<div v-if="accessLog.error != null && accessLog.error.length > 0" style="color:#dc143c">
		<i class="icon warning circle"></i>错误：[{{accessLog.error}}]
	</div>
</div>`}),Vue.component("ns-cluster-selector",{props:["v-cluster-id"],mounted:function(){let t=this;Tea.action("/ns/clusters/options").post().success(function(e){t.clusters=e.data.clusters})},data:function(){let e=this.vClusterId;return{clusters:[],clusterId:e=null==e?0:e}},template:`<div>
	<select class="ui dropdown auto-width" name="clusterId" v-model="clusterId">
		<option value="0">[选择集群]</option>
		<option v-for="cluster in clusters" :value="cluster.id">{{cluster.name}}</option>
	</select>
</div>`}),Vue.component("plan-user-selector",{mounted:function(){let t=this;Tea.action("/plans/users/options").post().success(function(e){t.users=e.data.users})},props:["v-user-id"],data:function(){let e=this.vUserId;return{users:[],userId:e=null==e?0:e}},watch:{userId:function(e){this.$emit("change",e)}},template:`<div>
	<select class="ui dropdown auto-width" name="userId" v-model="userId">
		<option value="0">[选择用户]</option>
		<option v-for="user in users" :value="user.id">{{user.fullname}} ({{user.username}})</option>
	</select>
</div>`}),Vue.component("plan-price-view",{props:["v-plan"],data:function(){return{plan:this.vPlan}},template:`<div>
	 <span v-if="plan.priceType == 'period'">
	 	按时间周期计费
	 	<div>
	 		<span class="grey small">
				<span v-if="plan.monthlyPrice > 0">月度：￥{{plan.monthlyPrice}}元<br/></span>
				<span v-if="plan.seasonallyPrice > 0">季度：￥{{plan.seasonallyPrice}}元<br/></span>
				<span v-if="plan.yearlyPrice > 0">年度：￥{{plan.yearlyPrice}}元</span>
			</span>
		</div>
	</span>
	<span v-if="plan.priceType == 'traffic'">
		按流量计费
		<div>
			<span class="grey small">基础价格：￥{{plan.trafficPrice.base}}元/GB</span>
		</div>
	</span>
	<div v-if="plan.priceType == 'bandwidth' && plan.bandwidthPrice != null && plan.bandwidthPrice.percentile > 0">
		按{{plan.bandwidthPrice.percentile}}th带宽计费 
		<div>
			<div v-for="range in plan.bandwidthPrice.ranges">
				<span class="small grey">{{range.minMB}} - <span v-if="range.maxMB > 0">{{range.maxMB}}MB</span><span v-else>&infin;</span>： {{range.pricePerMB}}元/MB</span>
			</div>
		</div>
	</div>
</div>`}),Vue.component("plan-bandwidth-ranges",{props:["v-ranges"],data:function(){let e=this.vRanges;return{ranges:e=null==e?[]:e,isAdding:!1,minMB:"",maxMB:"",pricePerMB:"",addingRange:{minMB:0,maxMB:0,pricePerMB:0,totalPrice:0}}},methods:{add:function(){this.isAdding=!this.isAdding;let e=this;setTimeout(function(){e.$refs.minMB.focus()})},cancelAdding:function(){this.isAdding=!1},confirm:function(){this.isAdding=!1,this.minMB="",this.maxMB="",this.pricePerMB="",this.ranges.push(this.addingRange),this.ranges.$sort(function(e,t){return e.minMB<t.minMB?-1:e.minMB==t.minMB?0:1}),this.change(),this.addingRange={minMB:0,maxMB:0,pricePerMB:0,totalPrice:0}},remove:function(e){this.ranges.$remove(e),this.change()},change:function(){this.$emit("change",this.ranges)}},watch:{minMB:function(e){let t=parseInt(e.toString());(isNaN(t)||t<0)&&(t=0),this.addingRange.minMB=t},maxMB:function(e){let t=parseInt(e.toString());(isNaN(t)||t<0)&&(t=0),this.addingRange.maxMB=t},pricePerMB:function(e){let t=parseFloat(e.toString());(isNaN(t)||t<0)&&(t=0),this.addingRange.pricePerMB=t}},template:`<div>
	<!-- 已有价格 -->
	<div v-if="ranges.length > 0">
		<div class="ui label basic small" v-for="(range, index) in ranges" style="margin-bottom: 0.5em">
			{{range.minMB}}MB - <span v-if="range.maxMB > 0">{{range.maxMB}}MB</span><span v-else>&infin;</span> &nbsp;  价格：{{range.pricePerMB}}元/MB
			&nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	
	<!-- 添加 -->
	<div v-if="isAdding">
		<table class="ui table">
			<tr>
				<td class="title">带宽下限</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" placeholder="最小带宽" style="width: 7em" maxlength="10" ref="minMB" @keyup.enter="confirm()" @keypress.enter.prevent="1" v-model="minMB"/>
						<span class="ui label">MB</span>
					</div>
				</td>
			</tr>
			<tr>
				<td class="title">带宽上限</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" placeholder="最大带宽" style="width: 7em" maxlength="10" @keyup.enter="confirm()" @keypress.enter.prevent="1" v-model="maxMB"/>
						<span class="ui label">MB</span>
					</div>
					<p class="comment">如果填0，表示上不封顶。</p>
				</td>
			</tr>
			<tr>
				<td class="title">单位价格</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" placeholder="单位价格" style="width: 7em" maxlength="10" @keyup.enter="confirm()" @keypress.enter.prevent="1" v-model="pricePerMB"/>
						<span class="ui label">元/MB</span>
					</div>
				</td>
			</tr>
		</table>
		<button class="ui button small" type="button" @click.prevent="confirm">确定</button> &nbsp;
		<a href="" title="取消" @click.prevent="cancelAdding"><i class="icon remove small"></i></a>
	</div>
	
	<!-- 按钮 -->
	<div v-if="!isAdding">
		<button class="ui button small" type="button" @click.prevent="add">+</button>
	</div>
</div>`}),Vue.component("plan-price-config-box",{props:["v-price-type","v-monthly-price","v-seasonally-price","v-yearly-price","v-traffic-price","v-bandwidth-price","v-disable-period"],data:function(){let e=this.vPriceType,t=(null==e&&(e="bandwidth"),0),i=this.vMonthlyPrice,s=(null==i||i<=0?i="":(i=i.toString(),t=parseFloat(i),isNaN(t)&&(t=0)),0),n=this.vSeasonallyPrice,o=(null==n||n<=0?n="":(n=n.toString(),s=parseFloat(n),isNaN(s)&&(s=0)),0),a=this.vYearlyPrice,l=(null==a||a<=0?a="":(a=a.toString(),o=parseFloat(a),isNaN(o)&&(o=0)),this.vTrafficPrice),r=0,c=(null!=l?r=l.base:l={base:0},""),d=(0<r&&(c=r.toString()),this.vBandwidthPrice);return null==d?d={percentile:95,ranges:[]}:null==d.ranges&&(d.ranges=[]),{priceType:e,monthlyPrice:i,seasonallyPrice:n,yearlyPrice:a,monthlyPriceNumber:t,seasonallyPriceNumber:s,yearlyPriceNumber:o,trafficPriceBase:c,trafficPrice:l,bandwidthPrice:d,bandwidthPercentile:d.percentile}},methods:{changeBandwidthPriceRanges:function(e){this.bandwidthPrice.ranges=e}},watch:{monthlyPrice:function(e){let t=parseFloat(e);isNaN(t)&&(t=0),this.monthlyPriceNumber=t},seasonallyPrice:function(e){let t=parseFloat(e);isNaN(t)&&(t=0),this.seasonallyPriceNumber=t},yearlyPrice:function(e){let t=parseFloat(e);isNaN(t)&&(t=0),this.yearlyPriceNumber=t},trafficPriceBase:function(e){let t=parseFloat(e);isNaN(t)&&(t=0),this.trafficPrice.base=t},bandwidthPercentile:function(e){let t=parseInt(e);isNaN(t)||t<=0?t=95:100<t&&(t=100),this.bandwidthPrice.percentile=t}},template:`<div>
	<input type="hidden" name="priceType" :value="priceType"/>
	<input type="hidden" name="monthlyPrice" :value="monthlyPriceNumber"/>
	<input type="hidden" name="seasonallyPrice" :value="seasonallyPriceNumber"/>
	<input type="hidden" name="yearlyPrice" :value="yearlyPriceNumber"/>
	<input type="hidden" name="trafficPriceJSON" :value="JSON.stringify(trafficPrice)"/>
	<input type="hidden" name="bandwidthPriceJSON" :value="JSON.stringify(bandwidthPrice)"/>
	
	<div>
		<radio :v-value="'bandwidth'" :value="priceType" v-model="priceType">&nbsp;按带宽</radio> &nbsp;
		<radio :v-value="'traffic'" :value="priceType" v-model="priceType">&nbsp;按流量</radio> &nbsp;
		<radio :v-value="'period'" :value="priceType" v-model="priceType" v-show="typeof(vDisablePeriod) != 'boolean' || !vDisablePeriod">&nbsp;按时间周期</radio>
	</div>
	
	<!-- 按时间周期 -->
	<div v-show="priceType == 'period'">
		<div class="ui divider"></div>
		<table class="ui table">
			<tr>
				<td class="title">月度价格</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 7em" maxlength="10" v-model="monthlyPrice"/>
						<span class="ui label">元</span>
					</div>
				</td>
			</tr>
			<tr>
				<td class="title">季度价格</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 7em" maxlength="10" v-model="seasonallyPrice"/>
						<span class="ui label">元</span>
					</div>
				</td>
			</tr>
			<tr>
				<td class="title">年度价格</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 7em" maxlength="10" v-model="yearlyPrice"/>
						<span class="ui label">元</span>
					</div>
				</td>
			</tr>
		</table>
	</div>
	
	<!-- 按流量 -->
	<div v-show="priceType =='traffic'">
		<div class="ui divider"></div>
		<table class="ui table">
			<tr>
				<td class="title">基础流量费用 *</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" v-model="trafficPriceBase" maxlength="10" style="width: 7em"/>
						<span class="ui label">元/GB</span>
					</div>
				</td>
			</tr>
		</table>
	</div>
	
	<!-- 按带宽 -->
	<div v-show="priceType == 'bandwidth'">
		<div class="ui divider"></div>
		<table class="ui table">
			<tr>
				<td class="title">带宽百分位 *</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 4em" maxlength="3" v-model="bandwidthPercentile"/>
						<span class="ui label">th</span>
					</div>
				</td>
			</tr>
			<tr>
				<td>带宽价格</td>
				<td>
					<plan-bandwidth-ranges :v-ranges="bandwidthPrice.ranges" @change="changeBandwidthPriceRanges"></plan-bandwidth-ranges>
				</td>
			</tr>
		</table>
	</div>
</div>`}),Vue.component("http-stat-config-box",{props:["v-stat-config","v-is-location","v-is-group"],data:function(){let e=this.vStatConfig;return{stat:e=null==e?{isPrior:!1,isOn:!1}:e}},template:`<div>
	<input type="hidden" name="statJSON" :value="JSON.stringify(stat)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="stat" v-if="vIsLocation || vIsGroup" ></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || stat.isPrior">
			<tr>
				<td class="title">是否开启统计</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="stat.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
	</table>
<div class="margin"></div>
</div>`}),Vue.component("http-request-conds-box",{props:["v-conds"],data:function(){let e=this.vConds;return{conds:e=null==e?{isOn:!0,connector:"or",groups:[]}:e,components:window.REQUEST_COND_COMPONENTS}},methods:{change:function(){this.$emit("change",this.conds)},addGroup:function(){window.UPDATING_COND_GROUP=null;let t=this;teaweb.popup("/servers/server/settings/conds/addGroupPopup",{height:"30em",callback:function(e){t.conds.groups.push(e.data.group),t.change()}})},updateGroup:function(t,e){window.UPDATING_COND_GROUP=e;let i=this;teaweb.popup("/servers/server/settings/conds/addGroupPopup",{height:"30em",callback:function(e){Vue.set(i.conds.groups,t,e.data.group),i.change()}})},removeGroup:function(e){let t=this;teaweb.confirm("确定要删除这一组条件吗？",function(){t.conds.groups.$remove(e),t.change()})},typeName:function(i){var e=this.components.$find(function(e,t){return t.type==i.type});return null!=e?e.name:i.param+" "+i.operator}},template:`<div>
		<input type="hidden" name="condsJSON" :value="JSON.stringify(conds)"/>
		<div v-if="conds.groups.length > 0">
			<table class="ui table">
				<tr v-for="(group, groupIndex) in conds.groups">
					<td class="title" :class="{'color-border':conds.connector == 'and'}" :style="{'border-bottom':(groupIndex < conds.groups.length-1) ? '1px solid rgba(34,36,38,.15)':''}">分组{{groupIndex+1}}</td>
					<td style="background: white; word-break: break-all" :style="{'border-bottom':(groupIndex < conds.groups.length-1) ? '1px solid rgba(34,36,38,.15)':''}">
						<var v-for="(cond, index) in group.conds" style="font-style: normal;display: inline-block; margin-bottom:0.5em">
							<span class="ui label tiny">
								<var v-if="cond.type.length == 0 || cond.type == 'params'" style="font-style: normal">{{cond.param}} <var>{{cond.operator}}</var></var>
								<var v-if="cond.type.length > 0 && cond.type != 'params'" style="font-style: normal">{{typeName(cond)}}: </var>
								{{cond.value}}
								<sup v-if="cond.isCaseInsensitive" title="不区分大小写"><i class="icon info small"></i></sup>
							</span>
							
							<var v-if="index < group.conds.length - 1"> {{group.connector}} &nbsp;</var>
						</var>
					</td>
					<td style="width: 5em; background: white" :style="{'border-bottom':(groupIndex < conds.groups.length-1) ? '1px solid rgba(34,36,38,.15)':''}">
						<a href="" title="修改分组" @click.prevent="updateGroup(groupIndex, group)"><i class="icon pencil small"></i></a> <a href="" title="删除分组" @click.prevent="removeGroup(groupIndex)"><i class="icon remove"></i></a>
					</td>
				</tr>
			</table>
			<div class="ui divider"></div>
		</div>
		
		<!-- 分组之间关系 -->
		<table class="ui table" v-if="conds.groups.length > 1">
			<tr>
				<td class="title">分组之间关系</td>
				<td>
					<select class="ui dropdown auto-width" v-model="conds.connector">
						<option value="and">和</option>
						<option value="or">或</option>
					</select>
					<p class="comment">
						<span v-if="conds.connector == 'or'">只要满足其中一个条件分组即可。</span>
						<span v-if="conds.connector == 'and'">需要满足所有条件分组。</span>
					</p>	
				</td>
			</tr>
		</table>
		
		<div>
			<button class="ui button tiny" type="button" @click.prevent="addGroup()">+添加分组</button>
		</div>
	</div>	
</div>`}),Vue.component("ssl-config-box",{props:["v-ssl-policy","v-protocol","v-server-id"],created:function(){let e=this;setTimeout(function(){e.sortableCipherSuites()},100)},data:function(){let e=this.vSslPolicy,t=(null==e?e={id:0,isOn:!0,certRefs:[],certs:[],clientCARefs:[],clientCACerts:[],clientAuthType:0,minVersion:"TLS 1.1",hsts:null,cipherSuitesIsOn:!1,cipherSuites:[],http2Enabled:!0,ocspIsOn:!1}:(null==e.certRefs&&(e.certRefs=[]),null==e.certs&&(e.certs=[]),null==e.clientCARefs&&(e.clientCARefs=[]),null==e.clientCACerts&&(e.clientCACerts=[]),null==e.cipherSuites&&(e.cipherSuites=[])),e.hsts);return null==t&&(t={isOn:!1,maxAge:31536e3,includeSubDomains:!1,preload:!1,domains:[]}),{policy:e,hsts:t,hstsOptionsVisible:!1,hstsDomainAdding:!1,addingHstsDomain:"",hstsDomainEditingIndex:-1,allVersions:window.SSL_ALL_VERSIONS,allCipherSuites:window.SSL_ALL_CIPHER_SUITES.$copy(),modernCipherSuites:window.SSL_MODERN_CIPHER_SUITES,intermediateCipherSuites:window.SSL_INTERMEDIATE_CIPHER_SUITES,allClientAuthTypes:window.SSL_ALL_CLIENT_AUTH_TYPES,cipherSuitesVisible:!1,moreOptionsVisible:!1}},watch:{hsts:{deep:!0,handler:function(){this.policy.hsts=this.hsts}}},methods:{removeCert:function(e){let t=this;teaweb.confirm("确定删除此证书吗？证书数据仍然保留，只是当前服务不再使用此证书。",function(){t.policy.certRefs.$remove(e),t.policy.certs.$remove(e)})},selectCert:function(){let t=this,i=[];null!=this.policy&&0<this.policy.certs.length&&this.policy.certs.forEach(function(e){i.push(e.id.toString())}),teaweb.popup("/servers/certs/selectPopup?selectedCertIds="+i,{width:"50em",height:"30em",callback:function(e){t.policy.certRefs.push(e.data.certRef),t.policy.certs.push(e.data.cert)}})},uploadCert:function(){let t=this;teaweb.popup("/servers/certs/uploadPopup",{height:"28em",callback:function(e){teaweb.success("上传成功",function(){t.policy.certRefs.push(e.data.certRef),t.policy.certs.push(e.data.cert)})}})},requestCert:function(){let t=[],e=(null!=this.policy&&0<this.policy.certs.length&&this.policy.certs.forEach(function(e){t.$pushAll(e.dnsNames)}),this);teaweb.popup("/servers/server/settings/https/requestCertPopup?serverId="+this.vServerId+"&excludeServerNames="+t.join(","),{callback:function(){e.policy.certRefs.push(resp.data.certRef),e.policy.certs.push(resp.data.cert)}})},changeOptionsVisible:function(){this.moreOptionsVisible=!this.moreOptionsVisible},formatTime:function(e){return new Date(1e3*e).format("Y-m-d")},formatCipherSuite:function(e){return e.replace(/(AES|3DES)/,'<var style="font-weight: bold">$1</var>')},addCipherSuite:function(e){this.policy.cipherSuites.$contains(e)||this.policy.cipherSuites.push(e),this.allCipherSuites.$removeValue(e)},removeCipherSuite:function(e){let i=this;teaweb.confirm("确定要删除此套件吗？",function(){i.policy.cipherSuites.$removeValue(e),i.allCipherSuites=window.SSL_ALL_CIPHER_SUITES.$findAll(function(e,t){return!i.policy.cipherSuites.$contains(t)})})},clearCipherSuites:function(){let e=this;teaweb.confirm("确定要清除所有已选套件吗？",function(){e.policy.cipherSuites=[],e.allCipherSuites=window.SSL_ALL_CIPHER_SUITES.$copy()})},addBatchCipherSuites:function(e){var i=this;teaweb.confirm("确定要批量添加套件？",function(){e.$each(function(e,t){i.policy.cipherSuites.$contains(t)||i.policy.cipherSuites.push(t)})})},sortableCipherSuites:function(){var e=document.querySelector(".cipher-suites-box");Sortable.create(e,{draggable:".label",handle:".icon.handle",onStart:function(){},onUpdate:function(e){}})},showAllCipherSuites:function(){this.cipherSuitesVisible=!this.cipherSuitesVisible},showMoreHSTS:function(){this.hstsOptionsVisible=!this.hstsOptionsVisible,this.hstsOptionsVisible&&this.changeHSTSMaxAge()},changeHSTSMaxAge:function(){var e=this.hsts.maxAge;isNaN(e)?this.hsts.days="-":(this.hsts.days=parseInt(e/86400),(isNaN(this.hsts.days)||this.hsts.days<0)&&(this.hsts.days="-"))},setHSTSMaxAge:function(e){this.hsts.maxAge=e,this.changeHSTSMaxAge()},addHstsDomain:function(){this.hstsDomainAdding=!0,this.hstsDomainEditingIndex=-1;let e=this;setTimeout(function(){e.$refs.addingHstsDomain.focus()},100)},editHstsDomain:function(e){this.hstsDomainEditingIndex=e,this.addingHstsDomain=this.hsts.domains[e],this.hstsDomainAdding=!0;let t=this;setTimeout(function(){t.$refs.addingHstsDomain.focus()},100)},confirmAddHstsDomain:function(){this.addingHstsDomain=this.addingHstsDomain.trim(),0!=this.addingHstsDomain.length&&(-1<this.hstsDomainEditingIndex?this.hsts.domains[this.hstsDomainEditingIndex]=this.addingHstsDomain:this.hsts.domains.push(this.addingHstsDomain),this.cancelHstsDomainAdding())},cancelHstsDomainAdding:function(){this.hstsDomainAdding=!1,this.addingHstsDomain="",this.hstsDomainEditingIndex=-1},removeHstsDomain:function(e){this.cancelHstsDomainAdding(),this.hsts.domains.$remove(e)},selectClientCACert:function(){let t=this;teaweb.popup("/servers/certs/selectPopup?isCA=1",{width:"50em",height:"30em",callback:function(e){t.policy.clientCARefs.push(e.data.certRef),t.policy.clientCACerts.push(e.data.cert)}})},uploadClientCACert:function(){let t=this;teaweb.popup("/servers/certs/uploadPopup?isCA=1",{height:"28em",callback:function(e){teaweb.success("上传成功",function(){t.policy.clientCARefs.push(e.data.certRef),t.policy.clientCACerts.push(e.data.cert)})}})},removeClientCACert:function(e){let t=this;teaweb.confirm("确定删除此证书吗？证书数据仍然保留，只是当前服务不再使用此证书。",function(){t.policy.clientCARefs.$remove(e),t.policy.clientCACerts.$remove(e)})}},template:`<div>
	<h4>SSL/TLS相关配置</h4>
	<input type="hidden" name="sslPolicyJSON" :value="JSON.stringify(policy)"/>
	<table class="ui table definition selectable">
		<tbody>
			<tr v-show="vProtocol == 'https'">
				<td class="title">启用HTTP/2</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="policy.http2Enabled"/>
						<label></label>
					</div>
				</td>
			</tr>
			<tr>
				<td class="title">选择证书</td>
				<td>
					<div v-if="policy.certs != null && policy.certs.length > 0">
						<div class="ui label small" v-for="(cert, index) in policy.certs">
							{{cert.name}} / {{cert.dnsNames}} / 有效至{{formatTime(cert.timeEndAt)}} &nbsp; <a href="" title="删除" @click.prevent="removeCert(index)"><i class="icon remove"></i></a>
						</div>
						<div class="ui divider"></div>
					</div>
					<div v-else>
						<span class="red">选择或上传证书后<span v-if="vProtocol == 'https'">HTTPS</span><span v-if="vProtocol == 'tls'">TLS</span>服务才能生效。</span>
						<div class="ui divider"></div>
					</div>
					<button class="ui button tiny" type="button" @click.prevent="selectCert()">选择已有证书</button> &nbsp;
					<button class="ui button tiny" type="button" @click.prevent="uploadCert()">上传新证书</button> &nbsp;
					<button class="ui button tiny" type="button" @click.prevent="requestCert()" v-if="vServerId != null && vServerId > 0">申请免费证书</button>
				</td>
			</tr>
			<tr>
				<td>TLS最低版本</td>
				<td>
					<select v-model="policy.minVersion" class="ui dropdown auto-width">
						<option v-for="version in allVersions" :value="version">{{version}}</option>
					</select>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeOptionsVisible"></more-options-tbody>
		<tbody v-show="moreOptionsVisible">
			<!-- 加密套件 -->
			<tr>
				<td>加密算法套件<em>（CipherSuites）</em></td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="policy.cipherSuitesIsOn" />
						<label>是否要自定义</label>
					</div>
					<div v-show="policy.cipherSuitesIsOn">
						<div class="ui divider"></div>
						<div class="cipher-suites-box">
							已添加套件({{policy.cipherSuites.length}})：
							<div v-for="cipherSuite in policy.cipherSuites" class="ui label tiny basic" style="margin-bottom: 0.5em">
								<input type="hidden" name="cipherSuites" :value="cipherSuite"/>
								<span v-html="formatCipherSuite(cipherSuite)"></span> &nbsp; <a href="" title="删除套件" @click.prevent="removeCipherSuite(cipherSuite)"><i class="icon remove"></i></a>
								<a href="" title="拖动改变顺序"><i class="icon bars handle"></i></a>
							</div>
						</div>
						<div>
							<div class="ui divider"></div>
							<span v-if="policy.cipherSuites.length > 0"><a href="" @click.prevent="clearCipherSuites()">[清除所有已选套件]</a> &nbsp; </span>
							<a href="" @click.prevent="addBatchCipherSuites(modernCipherSuites)">[添加推荐套件]</a> &nbsp;
							<a href="" @click.prevent="addBatchCipherSuites(intermediateCipherSuites)">[添加兼容套件]</a>
							<div class="ui divider"></div>
						</div>
		
						<div class="cipher-all-suites-box">
							<a href="" @click.prevent="showAllCipherSuites()"><span v-if="policy.cipherSuites.length == 0">所有</span>可选套件({{allCipherSuites.length}}) <i class="icon angle" :class="{down:!cipherSuitesVisible, up:cipherSuitesVisible}"></i></a>
							<a href="" v-if="cipherSuitesVisible" v-for="cipherSuite in allCipherSuites" class="ui label tiny" title="点击添加到自定义套件中" @click.prevent="addCipherSuite(cipherSuite)" v-html="formatCipherSuite(cipherSuite)" style="margin-bottom:0.5em"></a>
						</div>
						<p class="comment" v-if="cipherSuitesVisible">点击可选套件添加。</p>
					</div>
				</td>
			</tr>
			
			<!-- HSTS -->
			<tr v-show="vProtocol == 'https'">
				<td :class="{'color-border':hsts.isOn}">是否开启HSTS</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" name="hstsOn" v-model="hsts.isOn" value="1"/>
						<label></label>
					</div>
					<p class="comment">
						 开启后，会自动在响应Header中加入
						 <span class="ui label small">Strict-Transport-Security:
							 <var v-if="!hsts.isOn">...</var>
							 <var v-if="hsts.isOn"><span>max-age=</span>{{hsts.maxAge}}</var>
							 <var v-if="hsts.isOn && hsts.includeSubDomains">; includeSubDomains</var>
							 <var v-if="hsts.isOn && hsts.preload">; preload</var>
						 </span>
						  <span v-if="hsts.isOn">
							<a href="" @click.prevent="showMoreHSTS()">修改<i class="icon angle" :class="{down:!hstsOptionsVisible, up:hstsOptionsVisible}"></i> </a>
						 </span>
					</p>
				</td>
			</tr>
			<tr v-show="hsts.isOn && hstsOptionsVisible">
				<td class="color-border">HSTS有效时间<em>（max-age）</em></td>
				<td>
					<div class="ui fields inline">
						<div class="ui field">
							<input type="text" name="hstsMaxAge" v-model="hsts.maxAge" maxlength="10" size="10" @input="changeHSTSMaxAge()"/>
						</div>
						<div class="ui field">
							秒
						</div>
						<div class="ui field">{{hsts.days}}天</div>
					</div>
					<p class="comment">
						<a href="" @click.prevent="setHSTSMaxAge(31536000)" :class="{active:hsts.maxAge == 31536000}">[1年/365天]</a> &nbsp; &nbsp;
						<a href="" @click.prevent="setHSTSMaxAge(15768000)" :class="{active:hsts.maxAge == 15768000}">[6个月/182.5天]</a> &nbsp;  &nbsp;
						<a href="" @click.prevent="setHSTSMaxAge(2592000)"  :class="{active:hsts.maxAge == 2592000}">[1个月/30天]</a>
					</p>
				</td>
			</tr>
			<tr v-show="hsts.isOn && hstsOptionsVisible">
				<td class="color-border">HSTS包含子域名<em>（includeSubDomains）</em></td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" name="hstsIncludeSubDomains" value="1" v-model="hsts.includeSubDomains"/>
						<label></label>
					</div>
				</td>
			</tr>
			<tr v-show="hsts.isOn && hstsOptionsVisible">
				<td class="color-border">HSTS预加载<em>（preload）</em></td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" name="hstsPreload" value="1" v-model="hsts.preload"/>
						<label></label>
					</div>
				</td>
			</tr>
			<tr v-show="hsts.isOn && hstsOptionsVisible">
				<td class="color-border">HSTS生效的域名</td>
				<td colspan="2">
					<div class="names-box">
					<span class="ui label tiny basic" v-for="(domain, arrayIndex) in hsts.domains" :class="{blue:hstsDomainEditingIndex == arrayIndex}">{{domain}}
						<input type="hidden" name="hstsDomains" :value="domain"/> &nbsp;
						<a href="" @click.prevent="editHstsDomain(arrayIndex)" title="修改"><i class="icon pencil"></i></a>
						<a href="" @click.prevent="removeHstsDomain(arrayIndex)" title="删除"><i class="icon remove"></i></a>
					</span>
					</div>
					<div class="ui fields inline" v-if="hstsDomainAdding" style="margin-top:0.8em">
						<div class="ui field">
							<input type="text" name="addingHstsDomain" ref="addingHstsDomain" style="width:16em" maxlength="100" placeholder="域名，比如example.com" @keyup.enter="confirmAddHstsDomain()" @keypress.enter.prevent="1" v-model="addingHstsDomain" />
						</div>
						<div class="ui field">
							<button class="ui button tiny" type="button" @click="confirmAddHstsDomain()">确定</button>
							&nbsp; <a href="" @click.prevent="cancelHstsDomainAdding()">取消</a>
						</div>
					</div>
					<div class="ui field" style="margin-top: 1em">
						<button class="ui button tiny" type="button" @click="addHstsDomain()">+</button>
					</div>
					<p class="comment">如果没有设置域名的话，则默认支持所有的域名。</p>
				</td>
			</tr>
			
			<!-- OCSP -->
			<tr>
				<td>OCSP Stapling</td>
				<td><checkbox name="ocspIsOn" v-model="policy.ocspIsOn"></checkbox>
					<p class="comment">选中表示启用OCSP Stapling。</p>
				</td>
			</tr>
			
			<!-- 客户端认证 -->
			<tr>
				<td>客户端认证方式</td>
				<td>
					<select name="clientAuthType" v-model="policy.clientAuthType" class="ui dropdown auto-width">
						<option v-for="authType in allClientAuthTypes" :value="authType.type">{{authType.name}}</option>
					</select>
				</td>
			</tr>
			<tr>
				<td>客户端认证CA证书</td>
				<td>
					<div v-if="policy.clientCACerts != null && policy.clientCACerts.length > 0">
						<div class="ui label small" v-for="(cert, index) in policy.clientCACerts">
							{{cert.name}} / {{cert.dnsNames}} / 有效至{{formatTime(cert.timeEndAt)}} &nbsp; <a href="" title="删除" @click.prevent="removeClientCACert()"><i class="icon remove"></i></a>
						</div>
						<div class="ui divider"></div>
					</div>
					<button class="ui button tiny" type="button" @click.prevent="selectClientCACert()">选择已有证书</button> &nbsp;
					<button class="ui button tiny" type="button" @click.prevent="uploadClientCACert()">上传新证书</button>
					<p class="comment">用来校验客户端证书以增强安全性，通常不需要设置。</p>
				</td>
			</tr>
		</tbody>	
	</table>
	<div class="ui margin"></div>
</div>`}),Vue.component("http-firewall-actions-view",{props:["v-actions"],template:`<div>
		<div v-for="action in vActions" style="margin-bottom: 0.3em">
			<span :class="{red: action.category == 'block', orange: action.category == 'verify', green: action.category == 'allow'}">{{action.name}} ({{action.code.toUpperCase()}})</span>
		</div>             
</div>`}),Vue.component("http-request-scripts-config-box",{props:["vRequestScriptsConfig","v-is-location"],data:function(){let e=this.vRequestScriptsConfig;return{config:e=null==e?{}:e}},methods:{changeInitGroup:function(e){this.config.initGroup=e,this.$forceUpdate()},changeRequestGroup:function(e){this.config.requestGroup=e,this.$forceUpdate()}},template:`<div>
	<input type="hidden" name="requestScriptsJSON" :value="JSON.stringify(config)"/>
	<div class="margin"></div>
	<h4 style="margin-bottom: 0">请求初始化</h4>
	<p class="comment">在请求刚初始化时调用，此时自定义Header等尚未生效。</p>
	<div>
		<script-group-config-box :v-group="config.initGroup" @change="changeInitGroup" :v-is-location="vIsLocation"></script-group-config-box>
	</div>
	<h4 style="margin-bottom: 0">准备发送请求</h4>
	<p class="comment">在准备执行请求或者转发请求之前调用，此时自定义Header、源站等已准备好。</p>
	<div>
		<script-group-config-box :v-group="config.requestGroup" @change="changeRequestGroup" :v-is-location="vIsLocation"></script-group-config-box>
	</div>
	<div class="margin"></div>
</div>`}),Vue.component("http-firewall-rule-label",{props:["v-rule"],data:function(){return{rule:this.vRule}},methods:{showErr:function(e){teaweb.popupTip('规则校验错误，请修正：<span class="red">'+teaweb.encodeHTML(e)+"</span>")}},template:`<div>
	<div class="ui label tiny basic" style="line-height: 1.5">
		{{rule.name}}[{{rule.param}}] 

		<!-- cc2 -->
		<span v-if="rule.param == '\${cc2}'">
			{{rule.checkpointOptions.period}}秒/{{rule.checkpointOptions.threshold}}请求
		</span>

		<!-- refererBlock -->
		<span v-if="rule.param == '\${refererBlock}'">
			{{rule.checkpointOptions.allowDomains}}
		</span>

		<span v-else>
			<span v-if="rule.paramFilters != null && rule.paramFilters.length > 0" v-for="paramFilter in rule.paramFilters"> | {{paramFilter.code}}</span> 
		<var :class="{dash:rule.isCaseInsensitive}" :title="rule.isCaseInsensitive ? '大小写不敏感':''" v-if="!rule.isComposed">{{rule.operator}}</var> 
		{{rule.value}}
		</span>
		
		<!-- description -->
		<span v-if="rule.description != null && rule.description.length > 0" class="grey small">（{{rule.description}}）</span>
		
		<a href="" v-if="rule.err != null && rule.err.length > 0" @click.prevent="showErr(rule.err)" style="color: #db2828; opacity: 1; border-bottom: 1px #db2828 dashed; margin-left: 0.5em">规则错误</a>
	</div>
</div>`}),Vue.component("http-cache-refs-box",{props:["v-cache-refs"],data:function(){let e=this.vCacheRefs;return{refs:e=null==e?[]:e}},methods:{timeUnitName:function(e){switch(e){case"ms":return"毫秒";case"second":return"秒";case"minute":return"分钟";case"hour":return"小时";case"day":return"天";case"week":return"周 "}return e}},template:`<div>
	<input type="hidden" name="refsJSON" :value="JSON.stringify(refs)"/>
	
	<p class="comment" v-if="refs.length == 0">暂时还没有缓存条件。</p>
	<div v-show="refs.length > 0">
		<table class="ui table selectable celled">
			<thead>
				<tr>
					<th>缓存条件</th>
					<th class="two wide">分组关系</th>
					<th class="width10">缓存时间</th>
				</tr>
				<tr v-for="(cacheRef, index) in refs">
					<td :class="{'color-border': cacheRef.conds.connector == 'and', disabled: !cacheRef.isOn}" :style="{'border-left':cacheRef.isReverse ? '1px #db2828 solid' : ''}">
						<http-request-conds-view :v-conds="cacheRef.conds" :class="{disabled: !cacheRef.isOn}"></http-request-conds-view>
						<grey-label v-if="cacheRef.minSize != null && cacheRef.minSize.count > 0">
							{{cacheRef.minSize.count}}{{cacheRef.minSize.unit}}
							<span v-if="cacheRef.maxSize != null && cacheRef.maxSize.count > 0">- {{cacheRef.maxSize.count}}{{cacheRef.maxSize.unit}}</span>
						</grey-label>
						<grey-label v-else-if="cacheRef.maxSize != null && cacheRef.maxSize.count > 0">0 - {{cacheRef.maxSize.count}}{{cacheRef.maxSize.unit}}</grey-label>
						<grey-label v-if="cacheRef.methods != null && cacheRef.methods.length > 0">{{cacheRef.methods.join(", ")}}</grey-label>
						<grey-label v-if="cacheRef.expiresTime != null && cacheRef.expiresTime.isPrior && cacheRef.expiresTime.isOn">Expires</grey-label>
						<grey-label v-if="cacheRef.status != null && cacheRef.status.length > 0 && (cacheRef.status.length > 1 || cacheRef.status[0] != 200)">状态码：{{cacheRef.status.map(function(v) {return v.toString()}).join(", ")}}</grey-label>
						<grey-label v-if="cacheRef.allowPartialContent">区间缓存</grey-label>
					</td>
					<td :class="{disabled: !cacheRef.isOn}">
						<span v-if="cacheRef.conds.connector == 'and'">和</span>
						<span v-if="cacheRef.conds.connector == 'or'">或</span>
					</td>
					<td :class="{disabled: !cacheRef.isOn}">
						<span v-if="!cacheRef.isReverse">{{cacheRef.life.count}} {{timeUnitName(cacheRef.life.unit)}}</span>
						<span v-else class="red">不缓存</span>
					</td>
				</tr>
			</thead>
		</table>
	</div>
	<div class="margin"></div>
</div>`}),Vue.component("ssl-certs-box",{props:["v-certs","v-cert","v-protocol","v-view-size","v-single-mode","v-description"],data:function(){let e=this.vCerts,t=(null==e&&(e=[]),null!=this.vCert&&e.push(this.vCert),this.vDescription);return null!=t&&"string"==typeof t||(t=""),{certs:e,description:t}},methods:{certIds:function(){return this.certs.map(function(e){return e.id})},removeCert:function(e){let t=this;teaweb.confirm("确定删除此证书吗？证书数据仍然保留，只是当前服务不再使用此证书。",function(){t.certs.$remove(e)})},selectCert:function(){let t=this,e="50em",i="30em",s=this.vViewSize;"mini"==(s=null==s?"normal":s)&&(e="35em",i="20em"),teaweb.popup("/servers/certs/selectPopup?viewSize="+s,{width:e,height:i,callback:function(e){t.certs.push(e.data.cert)}})},uploadCert:function(){let t=this;teaweb.popup("/servers/certs/uploadPopup",{height:"28em",callback:function(e){teaweb.success("上传成功",function(){t.certs.push(e.data.cert)})}})},formatTime:function(e){return new Date(1e3*e).format("Y-m-d")},buttonsVisible:function(){return null==this.vSingleMode||!this.vSingleMode||null==this.certs||0==this.certs.length}},template:`<div>
	<input type="hidden" name="certIdsJSON" :value="JSON.stringify(certIds())"/>
	<div v-if="certs != null && certs.length > 0">
		<div class="ui label small basic" v-for="(cert, index) in certs">
			{{cert.name}} / {{cert.dnsNames}} / 有效至{{formatTime(cert.timeEndAt)}} &nbsp; <a href="" title="删除" @click.prevent="removeCert(index)"><i class="icon remove"></i></a>
		</div>
		<div class="ui divider" v-if="buttonsVisible()"></div>
	</div>
	<div v-else>
		<span class="red" v-if="description.length == 0">选择或上传证书后<span v-if="vProtocol == 'https'">HTTPS</span><span v-if="vProtocol == 'tls'">TLS</span>服务才能生效。</span>
		<span class="grey" v-if="description.length > 0">{{description}}</span>
		<div class="ui divider" v-if="buttonsVisible()"></div>
	</div>
	<div v-if="buttonsVisible()">
		<button class="ui button tiny" type="button" @click.prevent="selectCert()">选择已有证书</button> &nbsp;
		<button class="ui button tiny" type="button" @click.prevent="uploadCert()">上传新证书</button> &nbsp;
	</div>
</div>`}),Vue.component("http-host-redirect-box",{props:["v-redirects"],mounted:function(){let s=this;sortTable(function(e){let i=[];e.forEach(function(t){s.redirects.forEach(function(e){e.id==t&&i.push(e)})}),s.updateRedirects(i)})},data:function(){let e=this.vRedirects,t=(null==e&&(e=[]),0);return e.forEach(function(e){t++,e.id=t}),{redirects:e,statusOptions:[{code:301,text:"Moved Permanently"},{code:308,text:"Permanent Redirect"},{code:302,text:"Found"},{code:303,text:"See Other"},{code:307,text:"Temporary Redirect"}],id:t}},methods:{add:function(){let t=this;window.UPDATING_REDIRECT=null,teaweb.popup("/servers/server/settings/redirects/createPopup",{width:"50em",height:"30em",callback:function(e){t.id++,e.data.redirect.id=t.id,t.redirects.push(e.data.redirect),t.change()}})},update:function(t,i){let s=this;window.UPDATING_REDIRECT=i,teaweb.popup("/servers/server/settings/redirects/createPopup",{width:"50em",height:"30em",callback:function(e){e.data.redirect.id=i.id,Vue.set(s.redirects,t,e.data.redirect),s.change()}})},remove:function(e){let t=this;teaweb.confirm("确定要删除这条跳转规则吗？",function(){t.redirects.$remove(e),t.change()})},change:function(){let e=this;setTimeout(function(){e.$emit("change",e.redirects)},100)},updateRedirects:function(e){this.redirects=e,this.change()}},template:`<div>
	<input type="hidden" name="hostRedirectsJSON" :value="JSON.stringify(redirects)"/>
	
	<first-menu>
		<menu-item @click.prevent="add">[创建]</menu-item>
	</first-menu>
	<div class="margin"></div>

	<p class="comment" v-if="redirects.length == 0">暂时还没有URL跳转规则。</p>
	<div v-show="redirects.length > 0">
		<table class="ui table celled selectable" id="sortable-table">
			<thead>
				<tr>
					<th style="width: 1em"></th>
					<th>跳转前URL</th>
					<th style="width: 1em"></th>
					<th>跳转后URL</th>
					<th>匹配模式</th>
					<th>HTTP状态码</th>
					<th class="two wide">状态</th>
					<th class="two op">操作</th>
				</tr>
			</thead>
			<tbody v-for="(redirect, index) in redirects" :key="redirect.id" :v-id="redirect.id">
				<tr>
					<td style="text-align: center;"><i class="icon bars handle grey"></i> </td>
					<td>
						{{redirect.beforeURL}}
						<div style="margin-top: 0.5em" v-if="redirect.conds != null && redirect.conds.groups != null && redirect.conds.groups.length > 0">
							<span class="ui label text basic tiny">匹配条件</span>
						</div>
					</td>
					<td nowrap="">-&gt;</td>
					<td>{{redirect.afterURL}}</td>
					<td>
						<span v-if="redirect.matchPrefix">匹配前缀</span>
						<span v-if="redirect.matchRegexp">正则匹配</span>
						<span v-if="!redirect.matchPrefix && !redirect.matchRegexp">精准匹配</span>
					</td>
					<td>
						<span v-if="redirect.status > 0">{{redirect.status}}</span>
						<span v-else class="disabled">默认</span>
					</td>
					<td><label-on :v-is-on="redirect.isOn"></label-on></td>
					<td>
						<a href="" @click.prevent="update(index, redirect)">修改</a> &nbsp;
						<a href="" @click.prevent="remove(index)">删除</a>	
					</td>
				</tr>
			</tbody>
		</table>
		<p class="comment" v-if="redirects.length > 1">所有规则匹配顺序为从上到下，可以拖动左侧的<i class="icon bars"></i>排序。</p>
	</div>
	<div class="margin"></div>
</div>`}),Vue.component("http-cache-ref-box",{props:["v-cache-ref","v-is-reverse"],mounted:function(){this.$refs.variablesDescriber.update(this.ref.key)},data:function(){let e=this.vCacheRef;return null==(e=null==e?{isOn:!0,cachePolicyId:0,key:"${scheme}://${host}${requestPath}${isArgs}${args}",life:{count:2,unit:"hour"},status:[200],maxSize:{count:32,unit:"mb"},minSize:{count:0,unit:"kb"},skipCacheControlValues:["private","no-cache","no-store"],skipSetCookie:!0,enableRequestCachePragma:!1,conds:null,allowChunkedEncoding:!0,allowPartialContent:!1,isReverse:this.vIsReverse,methods:[],expiresTime:{isPrior:!1,isOn:!1,overwrite:!0,autoCalculate:!0,duration:{count:-1,unit:"hour"}}}:e).key&&(e.key=""),null==e.methods&&(e.methods=[]),null==e.life&&(e.life={count:2,unit:"hour"}),null==e.maxSize&&(e.maxSize={count:32,unit:"mb"}),null==e.minSize&&(e.minSize={count:0,unit:"kb"}),{ref:e,moreOptionsVisible:!1}},methods:{changeOptionsVisible:function(e){this.moreOptionsVisible=e},changeLife:function(e){this.ref.life=e},changeMaxSize:function(e){this.ref.maxSize=e},changeMinSize:function(e){this.ref.minSize=e},changeConds:function(e){this.ref.conds=e},changeStatusList:function(e){let t=[];e.forEach(function(e){e=parseInt(e);isNaN(e)||e<100||999<e||t.push(e)}),this.ref.status=t},changeMethods:function(e){this.ref.methods=e.map(function(e){return e.toUpperCase()})},changeKey:function(e){this.$refs.variablesDescriber.update(e)},changeExpiresTime:function(e){this.ref.expiresTime=e}},template:`<tbody>
	<tr>
		<td class="title">匹配条件分组 *</td>
		<td>
			<http-request-conds-box :v-conds="ref.conds" @change="changeConds"></http-request-conds-box>
			
			<input type="hidden" name="cacheRefJSON" :value="JSON.stringify(ref)"/>
		</td>
	</tr>
	<tr v-show="!vIsReverse">
		<td>缓存有效期 *</td>
		<td>
			<time-duration-box :v-value="ref.life" @change="changeLife"></time-duration-box>
		</td>
	</tr>
	<tr v-show="!vIsReverse">
		<td>缓存Key *</td>
		<td>
			<input type="text" v-model="ref.key" @input="changeKey(ref.key)"/>
			<p class="comment">用来区分不同缓存内容的唯一Key。<request-variables-describer ref="variablesDescriber"></request-variables-describer>。</p>
		</td>
	</tr>
	<tr v-show="!vIsReverse">
		<td colspan="2"><more-options-indicator @change="changeOptionsVisible"></more-options-indicator></td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>请求方法限制</td>
		<td>
			<values-box size="5" maxlength="10" :values="ref.methods" @change="changeMethods"></values-box>
			<p class="comment">允许请求的缓存方法，默认支持所有的请求方法。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>客户端过期时间<em>（Expires）</em></td>
		<td>
			<http-expires-time-config-box :v-expires-time="ref.expiresTime" @change="changeExpiresTime"></http-expires-time-config-box>		
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>可缓存的最大内容尺寸</td>
		<td>
			<size-capacity-box :v-value="ref.maxSize" @change="changeMaxSize"></size-capacity-box>
			<p class="comment">内容尺寸如果高于此值则不缓存。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>可缓存的最小内容尺寸</td>
		<td>
			<size-capacity-box :v-value="ref.minSize" @change="changeMinSize"></size-capacity-box>
			<p class="comment">内容尺寸如果低于此值则不缓存。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>支持分片内容</td>
		<td>
			<checkbox name="allowChunkedEncoding" value="1" v-model="ref.allowChunkedEncoding"></checkbox>
			<p class="comment">选中后，Gzip等压缩后的Chunked内容可以直接缓存，无需检查内容长度。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>支持缓存区间内容</td>
		<td>
			<checkbox name="allowPartialContent" value="1" v-model="ref.allowPartialContent"></checkbox>
			<p class="comment">选中后，支持缓存源站返回的某个区间的内容，该内容通过<code-label>206 Partial Content</code-label>状态码返回。此功能目前为<code-label>试验性质</code-label>。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>状态码列表</td>
		<td>
			<values-box name="statusList" size="3" maxlength="3" :values="ref.status" @change="changeStatusList"></values-box>
			<p class="comment">允许缓存的HTTP状态码列表。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>跳过的Cache-Control值</td>
		<td>
			<values-box name="skipResponseCacheControlValues" size="10" maxlength="100" :values="ref.skipCacheControlValues"></values-box>
			<p class="comment">当响应的Cache-Control为这些值时不缓存响应内容，而且不区分大小写。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>跳过Set-Cookie</td>
		<td>
			<div class="ui checkbox">
				<input type="checkbox" value="1" v-model="ref.skipSetCookie"/>
				<label></label>
			</div>
			<p class="comment">选中后，当响应的Header中有Set-Cookie时不缓存响应内容。</p>
		</td>
	</tr>
	<tr v-show="moreOptionsVisible && !vIsReverse">
		<td>支持请求no-cache刷新</td>
		<td>
			<div class="ui checkbox">
				<input type="checkbox" name="enableRequestCachePragma" value="1" v-model="ref.enableRequestCachePragma"/>
				<label></label>
			</div>
			<p class="comment">选中后，当请求的Header中含有Pragma: no-cache或Cache-Control: no-cache时，会跳过缓存直接读取源内容。</p>
		</td>
	</tr>	
</tbody>`}),Vue.component("http-request-limit-config-box",{props:["v-request-limit-config","v-is-group","v-is-location"],data:function(){let e=this.vRequestLimitConfig;return{config:e=null==e?{isPrior:!1,isOn:!1,maxConns:0,maxConnsPerIP:0,maxBodySize:{count:-1,unit:"kb"},outBandwidthPerConn:{count:-1,unit:"kb"}}:e,maxConns:e.maxConns,maxConnsPerIP:e.maxConnsPerIP}},watch:{maxConns:function(e){e=parseInt(e,10);isNaN(e)?this.config.maxConns=0:this.config.maxConns=e<0?0:e},maxConnsPerIP:function(e){e=parseInt(e,10);isNaN(e)?this.config.maxConnsPerIP=0:this.config.maxConnsPerIP=e<0?0:e}},methods:{isOn:function(){return(!this.vIsLocation&&!this.vIsGroup||this.config.isPrior)&&this.config.isOn}},template:`<div>
	<input type="hidden" name="requestLimitJSON" :value="JSON.stringify(config)"/>
	<table class="ui table selectable definition">
		<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
			<tr>
				<td class="title">是否启用</td>
				<td>
					<checkbox v-model="config.isOn"></checkbox>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td>最大并发连接数</td>
				<td>
					<input type="text" maxlength="6" v-model="maxConns"/>
					<p class="comment">当前服务最大并发连接数。为0表示不限制。</p>
				</td>
			</tr>
			<tr>
				<td>单IP最大并发连接数</td>
				<td>
					<input type="text" maxlength="6" v-model="maxConnsPerIP"/>
					<p class="comment">单IP最大连接数，统计单个IP总连接数时不区分服务。为0表示不限制。</p>
				</td>
			</tr>
			<tr>
				<td>单连接带宽限制</td>
				<td>
					<size-capacity-box :v-value="config.outBandwidthPerConn" :v-supported-units="['byte', 'kb', 'mb']"></size-capacity-box>
					<p class="comment">客户端单个请求每秒可以读取的下行流量。</p>
				</td>
			</tr>
			<tr>
				<td>单请求最大尺寸</td>
				<td>
					<size-capacity-box :v-value="config.maxBodySize" :v-supported-units="['byte', 'kb', 'mb', 'gb']"></size-capacity-box>
					<p class="comment">单个请求能发送的最大内容尺寸。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`}),Vue.component("http-header-replace-values",{props:["v-replace-values"],data:function(){let e=this.vReplaceValues;return{values:e=null==e?[]:e,isAdding:!1,addingValue:{pattern:"",replacement:"",isCaseInsensitive:!1,isRegexp:!1}}},methods:{add:function(){this.isAdding=!0;let e=this;setTimeout(function(){e.$refs.pattern.focus()})},remove:function(e){this.values.$remove(e)},confirm:function(){let e=this;0==this.addingValue.pattern.length?teaweb.warn("替换前内容不能为空",function(){e.$refs.pattern.focus()}):(this.values.push(this.addingValue),this.cancel())},cancel:function(){this.isAdding=!1,this.addingValue={pattern:"",replacement:"",isCaseInsensitive:!1,isRegexp:!1}}},template:`<div>
	<input type="hidden" name="replaceValuesJSON" :value="JSON.stringify(values)"/>
	<div>
		<div v-for="(value, index) in values" class="ui label small" style="margin-bottom: 0.5em">
			<var>{{value.pattern}}</var><sup v-if="value.isCaseInsensitive" title="不区分大小写"><i class="icon info tiny"></i></sup> =&gt; <var v-if="value.replacement.length > 0">{{value.replacement}}</var><var v-else><span class="small grey">[空]</span></var>
			<a href="" @click.prevent="remove(index)" title="删除"><i class="icon remove small"></i></a>
		</div>
	</div>
	<div v-if="isAdding">
		<table class="ui table">
			<tr>
				<td class="title">替换前内容 *</td>
				<td><input type="text" v-model="addingValue.pattern" placeholder="替换前内容" ref="pattern" @keyup.enter="confirm()" @keypress.enter.prevent="1"/></td>
			</tr>	
			<tr>
				<td>替换后内容</td>
				<td><input type="text" v-model="addingValue.replacement" placeholder="替换后内容" @keyup.enter="confirm()" @keypress.enter.prevent="1"/></td>
			</tr>
			<tr>
				<td>是否忽略大小写</td>
				<td>
					<checkbox v-model="addingValue.isCaseInsensitive"></checkbox>
				</td>
			</tr>
		</table>

		<div>
			<button type="button" class="ui button tiny" @click.prevent="confirm">确定</button> &nbsp;
			<a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
		</div>
	</div>
	<div v-if="!isAdding">
		<button type="button" class="ui button tiny" @click.prevent="add">+</button>
	</div>
</div>`}),Vue.component("http-request-conds-view",{props:["v-conds"],data:function(){let e=this.vConds,t=(null==e&&(e={isOn:!0,connector:"or",groups:[]}),this);return e.groups.forEach(function(e){e.conds.forEach(function(e){e.typeName=t.typeName(e)})}),{initConds:e}},computed:{conds:function(){return this.initConds}},methods:{typeName:function(i){var e=window.REQUEST_COND_COMPONENTS.$find(function(e,t){return t.type==i.type});return null!=e?e.name:i.param+" "+i.operator},updateConds:function(e){this.initConds=e},notifyChange:function(){let t=this;this.initConds.groups.forEach(function(e){e.conds.forEach(function(e){e.typeName=t.typeName(e)})}),this.$forceUpdate()}},template:`<div>
		<div v-if="conds.groups.length > 0">
			<div v-for="(group, groupIndex) in conds.groups">
				<var v-for="(cond, index) in group.conds" style="font-style: normal;display: inline-block; margin-bottom:0.5em">
					<span class="ui label small basic" style="line-height: 1.5">
						<var v-if="cond.type.length == 0 || cond.type == 'params'" style="font-style: normal">{{cond.param}} <var>{{cond.operator}}</var></var>
						<var v-if="cond.type.length > 0 && cond.type != 'params'" style="font-style: normal">{{cond.typeName}}: </var>
						{{cond.value}}
						<sup v-if="cond.isCaseInsensitive" title="不区分大小写"><i class="icon info small"></i></sup>
					</span>
					
					<var v-if="index < group.conds.length - 1"> {{group.connector}} &nbsp;</var>
				</var>
				<div class="ui divider" v-if="groupIndex != conds.groups.length - 1" style="margin-top:0.3em;margin-bottom:0.5em"></div>
				<div>
					<span class="ui label tiny olive" v-if="group.description != null && group.description.length > 0">{{group.description}}</span>
				</div>
			</div>	
		</div>
	</div>	
</div>`}),Vue.component("http-firewall-config-box",{props:["v-firewall-config","v-is-location","v-is-group","v-firewall-policy"],data:function(){let e=this.vFirewallConfig;return{firewall:e=null==e?{isPrior:!1,isOn:!1,firewallPolicyId:0}:e}},template:`<div>
	<input type="hidden" name="firewallJSON" :value="JSON.stringify(firewall)"/>
	<table class="ui table selectable definition">
		<prior-checkbox :v-config="firewall" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || firewall.isPrior">
			<tr v-show="!vIsGroup">
				<td>WAF策略</td>
				<td>
					<div v-if="vFirewallPolicy != null">{{vFirewallPolicy.name}} <span v-if="vFirewallPolicy.modeInfo != null">&nbsp; <span :class="{green: vFirewallPolicy.modeInfo.code == 'defend', blue: vFirewallPolicy.modeInfo.code == 'observe', grey: vFirewallPolicy.modeInfo.code == 'bypass'}">[{{vFirewallPolicy.modeInfo.name}}]</span>&nbsp;</span> <link-icon :href="'/servers/components/waf/policy?firewallPolicyId=' + vFirewallPolicy.id"></link-icon>
						<p class="comment">使用当前服务所在集群的设置。</p>
					</div>
					<span v-else class="red">当前集群没有设置WAF策略，当前配置无法生效。</span>
				</td>
			</tr>
			<tr>
				<td class="title">启用WAF</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="firewall.isOn"/>
						<label></label>
					</div>
					<p class="comment">启用WAF之后，各项WAF设置才会生效。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`}),Vue.component("metric-chart",{props:["v-chart","v-stats","v-item"],mounted:function(){this.load()},data:function(){let e=this.vStats;var t;0<(e=null==e?[]:e).length&&((t=e.$sum(function(e,t){return t.value}))<e[0].total&&"pie"==this.vChart.type&&e.push({keys:["其他"],value:e[0].total-t,total:e[0].total,time:e[0].time})),(e=0<this.vChart.maxItems?e.slice(0,this.vChart.maxItems):e.slice(0,10)).$rsort(function(e,t){return e.value-t.value});let i=100;return 0<this.vChart.widthDiv&&(i=100/this.vChart.widthDiv),{chart:this.vChart,stats:e,item:this.vItem,width:i+"%",chartId:"metric-chart-"+this.vChart.id,valueTypeName:null!=this.vItem&&null!=this.vItem.valueTypeName&&0<this.vItem.valueTypeName.length?this.vItem.valueTypeName:""}},methods:{load:function(){var e=document.getElementById(this.chartId);null==e||0==e.offsetWidth||0==e.offsetHeight?setTimeout(this.load,100):this.render(e)},render:function(e){let t=echarts.init(e);switch(window.addEventListener("resize",function(){t.resize()}),this.chart.type){case"pie":this.renderPie(t);break;case"bar":this.renderBar(t);break;case"timeBar":this.renderTimeBar(t);break;case"timeLine":this.renderTimeLine(t);break;case"table":this.renderTable(t)}},renderPie:function(e){var t=this.stats.map(function(e){return{name:e.keys[0],value:e.value}});let s=this;e.setOption({tooltip:{show:!0,trigger:"item",formatter:function(e){e=s.stats[e.dataIndex];let t=0,i=(0<e.total&&(t=Math.round(100*e.value/e.total*100)/100),e.value);switch(s.item.valueType){case"byte":i=teaweb.formatBytes(i);break;case"count":i=teaweb.formatNumber(i)}return e.keys[0]+"<br/>"+s.valueTypeName+": "+i+"<br/>占比："+t+"%"}},series:[{name:name,type:"pie",data:t,areaStyle:{},color:["#9DD3E8","#B2DB9E","#F39494","#FBD88A","#879BD7"]}]})},renderTimeBar:function(e){this.stats.$sort(function(e,t){return e.time<t.time?-1:1});let t=this.stats.map(function(e){return e.value}),i={unit:"",divider:1};switch(this.item.valueType){case"count":i=teaweb.countAxis(t,function(e){return e});break;case"byte":i=teaweb.bytesAxis(t,function(e){return e})}let s=this;e.setOption({xAxis:{data:this.stats.map(function(e){return s.formatTime(e.time)})},yAxis:{axisLabel:{formatter:function(e){return e+i.unit}}},tooltip:{show:!0,trigger:"item",formatter:function(e){e=s.stats[e.dataIndex];let t=e.value;return"byte"===s.item.valueType&&(t=teaweb.formatBytes(t)),s.formatTime(e.time)+": "+t}},grid:{left:50,top:10,right:20,bottom:25},series:[{name:name,type:"bar",data:t.map(function(e){return e/i.divider}),itemStyle:{color:teaweb.DefaultChartColor},areaStyle:{},barWidth:"20em"}]})},renderTimeLine:function(e){this.stats.$sort(function(e,t){return e.time<t.time?-1:1});let t=this.stats.map(function(e){return e.value}),i={unit:"",divider:1};switch(this.item.valueType){case"count":i=teaweb.countAxis(t,function(e){return e});break;case"byte":i=teaweb.bytesAxis(t,function(e){return e})}let s=this;e.setOption({xAxis:{data:this.stats.map(function(e){return s.formatTime(e.time)})},yAxis:{axisLabel:{formatter:function(e){return e+i.unit}}},tooltip:{show:!0,trigger:"item",formatter:function(e){e=s.stats[e.dataIndex];let t=e.value;return"byte"===s.item.valueType&&(t=teaweb.formatBytes(t)),s.formatTime(e.time)+": "+t}},grid:{left:50,top:10,right:20,bottom:25},series:[{name:name,type:"line",data:t.map(function(e){return e/i.divider}),itemStyle:{color:teaweb.DefaultChartColor},areaStyle:{}}]})},renderBar:function(e){let t=this.stats.map(function(e){return e.value}),i={unit:"",divider:1};switch(this.item.valueType){case"count":i=teaweb.countAxis(t,function(e){return e});break;case"byte":i=teaweb.bytesAxis(t,function(e){return e})}let s=24,n=0;var o=teaweb.xRotation(e,this.stats.map(function(e){return e.keys[0]}));null!=o&&(s=o[0],n=o[1]);let a=this;if(e.setOption({xAxis:{data:this.stats.map(function(e){return e.keys[0]}),axisLabel:{interval:0,rotate:n}},tooltip:{show:!0,trigger:"item",formatter:function(e){e=a.stats[e.dataIndex];let t=0,i=(0<e.total&&(t=Math.round(100*e.value/e.total*100)/100),e.value);switch(a.item.valueType){case"byte":i=teaweb.formatBytes(i);break;case"count":i=teaweb.formatNumber(i)}return e.keys[0]+"<br/>"+a.valueTypeName+"："+i+"<br/>占比："+t+"%"}},yAxis:{axisLabel:{formatter:function(e){return e+i.unit}}},grid:{left:40,top:10,right:20,bottom:s},series:[{name:name,type:"bar",data:t.map(function(e){return e/i.divider}),itemStyle:{color:teaweb.DefaultChartColor},areaStyle:{},barWidth:"20em"}]}),null!=this.item.keys&&this.item.keys.$contains("${remoteAddr}")){let i=this;e.on("click",function(e){var t=i.item.keys.$indexesOf("${remoteAddr}")[0],e=i.stats[e.dataIndex].keys[t];teaweb.popup("/servers/ipbox?ip="+e,{width:"50em",height:"30em"})})}},renderTable:function(e){let s=`<table class="ui table celled">
	<thead>
		<tr>
			<th>对象</th>
			<th>数值</th>
			<th>占比</th>
		</tr>
	</thead>`,n=this;this.stats.forEach(function(e){let t=e.value,i=("byte"===n.item.valueType&&(t=teaweb.formatBytes(t)),s+="<tr><td>"+e.keys[0]+"</td><td>"+t+"</td>",0);0<e.total&&(i=Math.round(100*e.value/e.total*100)/100),s=s+('<td><div class="ui progress blue"><div class="bar" style="min-width: 0; height: 4px; width: '+i+'%"></div></div>'+i)+"%</td></tr>"}),s+="</table>",document.getElementById(this.chartId).innerHTML=s},formatTime:function(e){if(null==e)return"";switch(this.item.periodUnit){case"month":case"week":return e.substring(0,4)+"-"+e.substring(4,6);case"day":return e.substring(0,4)+"-"+e.substring(4,6)+"-"+e.substring(6,8);case"hour":return e.substring(0,4)+"-"+e.substring(4,6)+"-"+e.substring(6,8)+" "+e.substring(8,10);case"minute":return e.substring(0,4)+"-"+e.substring(4,6)+"-"+e.substring(6,8)+" "+e.substring(8,10)+":"+e.substring(10,12)}return e}},template:`<div style="float: left" :style="{'width': width}">
	<h4>{{chart.name}} <span>（{{valueTypeName}}）</span></h4>
	<div class="ui divider"></div>
	<div style="height: 14em; padding-bottom: 1em; " :id="chartId" :class="{'scroll-box': chart.type == 'table'}"></div>
</div>`}),Vue.component("metric-board",{template:"<div><slot></slot></div>"}),Vue.component("http-cache-config-box",{props:["v-cache-config","v-is-location","v-is-group","v-cache-policy","v-web-id"],data:function(){let e=this.vCacheConfig;return null==(e=null==e?{isPrior:!1,isOn:!1,addStatusHeader:!0,addAgeHeader:!1,enableCacheControlMaxAge:!1,cacheRefs:[],purgeIsOn:!1,purgeKey:"",disablePolicyRefs:!1}:e).cacheRefs&&(e.cacheRefs=[]),{cacheConfig:e,moreOptionsVisible:!1,enablePolicyRefs:!e.disablePolicyRefs}},watch:{enablePolicyRefs:function(e){this.cacheConfig.disablePolicyRefs=!e}},methods:{isOn:function(){return(!this.vIsLocation&&!this.vIsGroup||this.cacheConfig.isPrior)&&this.cacheConfig.isOn},isPlus:function(){return Tea.Vue.teaIsPlus},generatePurgeKey:function(){let e=Math.random().toString()+Math.random().toString(),t=e.replace(/0\./g,"").replace(/\./g,""),i="";for(let e=0;e<t.length;e++)i+=String.fromCharCode(parseInt(t.substring(e,e+1))+(Math.random()<.5?"a":"A").charCodeAt(0));this.cacheConfig.purgeKey=i},showMoreOptions:function(){this.moreOptionsVisible=!this.moreOptionsVisible},changeStale:function(e){this.cacheConfig.stale=e}},template:`<div>
	<input type="hidden" name="cacheJSON" :value="JSON.stringify(cacheConfig)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="cacheConfig" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || cacheConfig.isPrior">
			<tr v-show="!vIsGroup">
				<td>缓存策略</td>
				<td>
					<div v-if="vCachePolicy != null">{{vCachePolicy.name}} <link-icon :href="'/servers/components/cache/policy?cachePolicyId=' + vCachePolicy.id"></link-icon>
						<p class="comment">使用当前服务所在集群的设置。</p>
					</div>
					<span v-else class="red">当前集群没有设置缓存策略，当前配置无法生效。</span>
				</td>
			</tr>
			<tr>
				<td class="title">启用缓存</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="cacheConfig.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td colspan="2">
					<a href="" @click.prevent="showMoreOptions"><span v-if="moreOptionsVisible">收起选项</span><span v-else>更多选项</span><i class="icon angle" :class="{up: moreOptionsVisible, down:!moreOptionsVisible}"></i></a>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn() && moreOptionsVisible">
			<tr>
				<td>使用默认缓存条件</td>
				<td>	
					<checkbox v-model="enablePolicyRefs"></checkbox>
					<p class="comment">选中后使用系统中已经定义的默认缓存条件。</p>
				</td>
			</tr>
			<tr>
				<td>添加X-Cache Header</td>
				<td>
					<checkbox v-model="cacheConfig.addStatusHeader"></checkbox>
					<p class="comment">选中后自动在响应Header中增加<code-label>X-Cache: BYPASS|MISS|HIT|PURGE</code-label>。</p>
				</td>
			</tr>
			<tr>
				<td>添加Age Header</td>
				<td>
					<checkbox v-model="cacheConfig.addAgeHeader"></checkbox>
					<p class="comment">选中后自动在响应Header中增加<code-label>Age: [存活时间秒数]</code-label>。</p>
				</td>
			</tr>
			<tr>
				<td>支持源站控制有效时间</td>
				<td>
					<checkbox v-model="cacheConfig.enableCacheControlMaxAge"></checkbox>
					<p class="comment">选中后表示支持源站在Header中设置的<code-label>Cache-Control: max-age=[有效时间秒数]</code-label>。</p>
				</td>
			</tr>
			<tr>
				<td class="color-border">允许PURGE</td>
				<td>
					<checkbox v-model="cacheConfig.purgeIsOn"></checkbox>
					<p class="comment">允许使用PURGE方法清除某个URL缓存。</p>
				</td>
			</tr>
			<tr v-show="cacheConfig.purgeIsOn">
				<td class="color-border">PURGE Key *</td>
				<td>
					<input type="text" maxlength="200" v-model="cacheConfig.purgeKey"/>
					<p class="comment"><a href="" @click.prevent="generatePurgeKey">[随机生成]</a>。需要在PURGE方法调用时加入<code-label>Edge-Purge-Key: {{cacheConfig.purgeKey}}</code-label> Header。只能包含字符、数字、下划线。</p>
				</td>
			</tr>
		</tbody>
	</table>
	
	<div v-if="isOn() && moreOptionsVisible && isPlus()">
		<h4>过时缓存策略</h4>
		<http-cache-stale-config :v-cache-stale-config="cacheConfig.stale" @change="changeStale"></http-cache-stale-config>
	</div>
	
	<div v-show="isOn()" style="margin-top: 1em">
		<h4>缓存条件</h4>
		<http-cache-refs-config-box :v-cache-config="cacheConfig" :v-cache-refs="cacheConfig.cacheRefs" :v-web-id="vWebId"></http-cache-refs-config-box>
	</div>
	<div class="margin"></div>
</div>`});let defaultGeneralHeaders=["Cache-Control","Connection","Date","Pragma","Trailer","Transfer-Encoding","Upgrade","Via","Warning"];function sortTable(n){let e=document.createElement("script");e.setAttribute("src","/js/sortable.min.js"),e.addEventListener("load",function(){let s=document.querySelector("#sortable-table");null!=s&&Sortable.create(s,{draggable:"tbody",handle:".icon.handle",onStart:function(){},onUpdate:function(e){let t=s.querySelectorAll("tbody"),i=[];t.forEach(function(e){i.push(parseInt(e.getAttribute("v-id")))}),n(i)}})}),document.head.appendChild(e)}function sortLoad(e){let t=document.createElement("script");t.setAttribute("src","/js/sortable.min.js"),t.addEventListener("load",function(){"function"==typeof e&&e()}),document.head.appendChild(t)}function emitClick(e,arguments){let t=["click"];for(let e=0;e<arguments.length;e++)t.push(arguments[e]);e.$emit.apply(e,t)}Vue.component("http-cond-general-header-length",{props:["v-checkpoint"],data:function(){let e=null,t=null;var i;null!=window.parent.UPDATING_RULE&&(null!=(i=window.parent.UPDATING_RULE.checkpointOptions).headers&&Array.$isArray(i.headers)&&(e=i.headers),null!=i.length&&(t=i.length)),null==e&&(e=defaultGeneralHeaders),null==t&&(t=128);let s=this;return setTimeout(function(){s.change()},100),{headers:e,length:t}},watch:{length:function(e){let t=parseInt(e);(t=isNaN(t)?0:t)<0&&(t=0),this.length=t,this.change()}},methods:{change:function(){this.vCheckpoint.options=[{code:"headers",value:this.headers},{code:"length",value:this.length}]}},template:`<div>
	<table class="ui table">
		<tr>
			<td class="title">通用Header列表</td>
			<td>
				<values-box :values="headers" :placeholder="'Header'" @change="change"></values-box>
				<p class="comment">需要检查的Header列表。</p>
			</td>
		</tr>
		<tr>
			<td>Header值超出长度</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" name="" style="width: 5em" v-model="length" maxlength="6"/>
					<span class="ui label">字节</span>
				</div>
				<p class="comment">超出此长度认为匹配成功，0表示不限制。</p>
			</td>
		</tr>
	</table>
</div>`}),Vue.component("http-firewall-checkpoint-cc",{props:["v-checkpoint"],data:function(){let e=[],t=60,i=1e3,s={},n=(null==(s=null!=window.parent.UPDATING_RULE?window.parent.UPDATING_RULE.checkpointOptions:s)&&(s={}),0==(e=null!=s.keys?s.keys:e).length&&(e=["${remoteAddr}","${requestPath}"]),null!=s.period&&(t=s.period),null!=s.threshold&&(i=s.threshold),this);return setTimeout(function(){n.change()},100),{keys:e,period:t,threshold:i,options:{},value:i}},watch:{period:function(){this.change()},threshold:function(){this.change()}},methods:{changeKeys:function(e){this.keys=e,this.change()},change:function(){let e=parseInt(this.period.toString()),t=((isNaN(e)||e<=0)&&(e=60),parseInt(this.threshold.toString()));(isNaN(t)||t<=0)&&(t=1e3),this.value=t,this.vCheckpoint.options=[{code:"keys",value:this.keys},{code:"period",value:e},{code:"threshold",value:t}]}},template:`<div>
	<input type="hidden" name="operator" value="gt"/>
	<input type="hidden" name="value" :value="value"/>
	<table class="ui table">
		<tr>
			<td class="title">统计对象组合 *</td>
			<td>
				<metric-keys-config-box :v-keys="keys" @change="changeKeys"></metric-keys-config-box>
			</td>
		</tr>
		<tr>
			<td>统计周期 *</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" v-model="period" style="width: 6em" maxlength="8"/>
					<span class="ui label">秒</span>
				</div>
			</td>
		</tr>
		<tr>
			<td>阈值 *</td>
			<td>
				<input type="text" v-model="threshold" style="width: 6em" maxlength="8"/>
			</td>
		</tr>
	</table>
</div>`}),Vue.component("http-firewall-checkpoint-referer-block",{props:["v-checkpoint"],data:function(){let e=!0,t=!0,i=[],s={},n=("boolean"==typeof(s=null==(s=null!=window.parent.UPDATING_RULE?window.parent.UPDATING_RULE.checkpointOptions:s)?{}:s).allowEmpty&&(e=s.allowEmpty),"boolean"==typeof s.allowSameDomain&&(t=s.allowSameDomain),null!=s.allowDomains&&"object"==typeof s.allowDomains&&(i=s.allowDomains),this);return setTimeout(function(){n.change()},100),{allowEmpty:e,allowSameDomain:t,allowDomains:i,options:{},value:0}},watch:{allowEmpty:function(){this.change()},allowSameDomain:function(){this.change()}},methods:{changeAllowDomains:function(e){this.allowDomains=e,this.change()},change:function(){this.vCheckpoint.options=[{code:"allowEmpty",value:this.allowEmpty},{code:"allowSameDomain",value:this.allowSameDomain},{code:"allowDomains",value:this.allowDomains}]}},template:`<div>
	<input type="hidden" name="operator" value="eq"/>
	<input type="hidden" name="value" :value="value"/>
	<table class="ui table">
		<tr>
			<td class="title">来源域名允许为空</td>
			<td>
				<checkbox v-model="allowEmpty"></checkbox>
				<p class="comment">允许不带来源的访问。</p>
			</td>
		</tr>
		<tr>
			<td>来源域名允许一致</td>
			<td>
				<checkbox v-model="allowSameDomain"></checkbox>
				<p class="comment">允许来源域名和当前访问的域名一致，相当于在站内访问。</p>
			</td>
		</tr>
		<tr>
			<td>允许的来源域名</td>
			<td>
				<values-box :values="allowDomains" @change="changeAllowDomains"></values-box>
				<p class="comment">允许的来源域名列表，比如<code-label>example.com</code-label>、<code-label>*.example.com</code-label>。单个星号<code-label>*</code-label>表示允许所有域名。</p>
			</td>
		</tr>
	</table>
</div>`}),Vue.component("http-access-log-partitions-box",{props:["v-partition","v-day"],mounted:function(){let t=this;Tea.action("/servers/logs/partitionData").params({day:this.vDay}).success(function(e){t.partitions=[],e.data.partitions.reverse().forEach(function(e){t.partitions.push({code:e,isDisabled:!1})}),0<t.partitions.length&&(null==t.vPartition||t.vPartition<0)&&(t.selectedPartition=t.partitions[0].code)}).post()},data:function(){return{partitions:[],selectedPartition:this.vPartition}},methods:{url:function(e){let t=window.location.toString();return 0<(t=(t=(t=(t=t.replace(/\?partition=-?\d+/,"?")).replace(/\?requestId=-?\d+/,"?")).replace(/&partition=-?\d+/,"")).replace(/&requestId=-?\d+/,"")).indexOf("?")?t+="&partition="+e:t+="?partition="+e,t},disable:function(t){this.partitions.forEach(function(e){e.code==t&&(e.isDisabled=!0)})}},template:`<div v-if="partitions.length > 1">
	<div class="ui divider" style="margin-bottom: 0"></div>
	<div class="ui menu text small blue" style="margin-bottom: 0; margin-top: 0">
		<a v-for="(p, index) in partitions" :href="url(p.code)" class="item" :class="{active: selectedPartition == p.code, disabled: p.isDisabled}">分表{{p.code+1}} &nbsp; &nbsp; <span class="disabled" v-if="index != partitions.length - 1">|</span></a>
	</div>
	<div class="ui divider" style="margin-top: 0"></div>
</div>`}),Vue.component("http-cache-refs-config-box",{props:["v-cache-refs","v-cache-config","v-cache-policy-id","v-web-id"],mounted:function(){let s=this;sortTable(function(e){let i=[];e.forEach(function(t){s.refs.forEach(function(e){e.id==t&&i.push(e)})}),s.updateRefs(i),s.change()})},data:function(){let e=this.vCacheRefs,t=(null==e&&(e=[]),0);return e.forEach(function(e){t++,e.id=t}),{refs:e,id:t}},methods:{addRef:function(e){window.UPDATING_CACHE_REF=null;let t=window.innerWidth,i=(1024<t&&(t=1024),window.innerHeight),n=(500<i&&(i=500),this);teaweb.popup("/servers/server/settings/cache/createPopup?isReverse="+(e?1:0),{width:t+"px",height:i+"px",callback:function(e){let s=e.data.cacheRef;if(null!=s.conds){if(n.id++,s.id=n.id,s.isReverse){let t=[],i=!1;n.refs.forEach(function(e){e.isReverse||i||(t.push(s),i=!0),t.push(e)}),i||t.push(s),n.updateRefs(t)}else n.refs.push(s);n.change()}}})},updateRef:function(t,e){window.UPDATING_CACHE_REF=e;let i=window.innerWidth,s=(1024<i&&(i=1024),window.innerHeight),n=(500<s&&(s=500),this);teaweb.popup("/servers/server/settings/cache/createPopup",{width:i+"px",height:s+"px",callback:function(e){e.data.cacheRef.id=n.refs[t].id,Vue.set(n.refs,t,e.data.cacheRef),n.change(),n.$refs.cacheRef[t].updateConds(e.data.cacheRef.conds),n.$refs.cacheRef[t].notifyChange()}})},disableRef:function(e){e.isOn=!1,this.change()},enableRef:function(e){e.isOn=!0,this.change()},removeRef:function(e){let t=this;teaweb.confirm("确定要删除此缓存设置吗？",function(){t.refs.$remove(e),t.change()})},updateRefs:function(e){this.refs=e,null!=this.vCacheConfig&&(this.vCacheConfig.cacheRefs=e)},timeUnitName:function(e){switch(e){case"ms":return"毫秒";case"second":return"秒";case"minute":return"分钟";case"hour":return"小时";case"day":return"天";case"week":return"周 "}return e},change:function(){this.$forceUpdate(),null!=this.vCachePolicyId&&0<this.vCachePolicyId?Tea.action("/servers/components/cache/updateRefs").params({cachePolicyId:this.vCachePolicyId,refsJSON:JSON.stringify(this.refs)}).post():null!=this.vWebId&&0<this.vWebId&&Tea.action("/servers/server/settings/cache/updateRefs").params({webId:this.vWebId,refsJSON:JSON.stringify(this.refs)}).success(function(e){e.data.isUpdated&&teaweb.successToast("保存成功")}).post()}},template:`<div>
	<input type="hidden" name="refsJSON" :value="JSON.stringify(refs)"/>
	
	<div>
		<p class="comment" v-if="refs.length == 0">暂时还没有缓存条件。</p>
		<table class="ui table selectable celled" v-show="refs.length > 0" id="sortable-table">
			<thead>
				<tr>
					<th style="width:1em"></th>
					<th>缓存条件</th>
					<th class="two wide">分组关系</th>
					<th class="width10">缓存时间</th>
					<th class="three op">操作</th>
				</tr>
			</thead>	
			<tbody v-for="(cacheRef, index) in refs" :key="cacheRef.id" :v-id="cacheRef.id">
				<tr>
					<td style="text-align: center;"><i class="icon bars handle grey"></i> </td>
					<td :class="{'color-border': cacheRef.conds.connector == 'and', disabled: !cacheRef.isOn}" :style="{'border-left':cacheRef.isReverse ? '1px #db2828 solid' : ''}">
						<http-request-conds-view :v-conds="cacheRef.conds" ref="cacheRef" :class="{disabled: !cacheRef.isOn}"></http-request-conds-view>
						<grey-label v-if="cacheRef.minSize != null && cacheRef.minSize.count > 0">
							{{cacheRef.minSize.count}}{{cacheRef.minSize.unit}}
							<span v-if="cacheRef.maxSize != null && cacheRef.maxSize.count > 0">- {{cacheRef.maxSize.count}}{{cacheRef.maxSize.unit}}</span>
						</grey-label>
						<grey-label v-else-if="cacheRef.maxSize != null && cacheRef.maxSize.count > 0">0 - {{cacheRef.maxSize.count}}{{cacheRef.maxSize.unit}}</grey-label>
						<grey-label v-if="cacheRef.methods != null && cacheRef.methods.length > 0">{{cacheRef.methods.join(", ")}}</grey-label>
						<grey-label v-if="cacheRef.expiresTime != null && cacheRef.expiresTime.isPrior && cacheRef.expiresTime.isOn">Expires</grey-label>
						<grey-label v-if="cacheRef.status != null && cacheRef.status.length > 0 && (cacheRef.status.length > 1 || cacheRef.status[0] != 200)">状态码：{{cacheRef.status.map(function(v) {return v.toString()}).join(", ")}}</grey-label>
						<grey-label v-if="cacheRef.allowPartialContent">区间缓存</grey-label>
					</td>
					<td :class="{disabled: !cacheRef.isOn}">
						<span v-if="cacheRef.conds.connector == 'and'">和</span>
						<span v-if="cacheRef.conds.connector == 'or'">或</span>
					</td>
					<td :class="{disabled: !cacheRef.isOn}">
						<span v-if="!cacheRef.isReverse">{{cacheRef.life.count}} {{timeUnitName(cacheRef.life.unit)}}</span>
						<span v-else class="red">不缓存</span>
					</td>
					<td>
						<a href="" @click.prevent="updateRef(index, cacheRef)">修改</a> &nbsp;
						<a href="" v-if="cacheRef.isOn" @click.prevent="disableRef(cacheRef)">暂停</a><a href="" v-if="!cacheRef.isOn" @click.prevent="enableRef(cacheRef)"><span class="red">恢复</span></a> &nbsp;
						<a href="" @click.prevent="removeRef(index)">删除</a>
					</td>
				</tr>
			</tbody>
		</table>
		<p class="comment" v-if="refs.length > 1">所有条件匹配顺序为从上到下，可以拖动左侧的<i class="icon bars"></i>排序。服务设置的优先级比全局缓存策略设置的优先级要高。</p>
		
		<button class="ui button tiny" @click.prevent="addRef(false)" type="button">+添加缓存设置</button> &nbsp; &nbsp; <a href="" @click.prevent="addRef(true)">+添加不缓存设置</a>
	</div>
	<div class="margin"></div>
</div>`}),Vue.component("origin-list-box",{props:["v-primary-origins","v-backup-origins","v-server-type","v-params"],data:function(){return{primaryOrigins:this.vPrimaryOrigins,backupOrigins:this.vBackupOrigins}},methods:{createPrimaryOrigin:function(){teaweb.popup("/servers/server/settings/origins/addPopup?originType=primary&"+this.vParams,{height:"27em",callback:function(e){teaweb.success("保存成功",function(){window.location.reload()})}})},createBackupOrigin:function(){teaweb.popup("/servers/server/settings/origins/addPopup?originType=backup&"+this.vParams,{height:"27em",callback:function(e){teaweb.success("保存成功",function(){window.location.reload()})}})},updateOrigin:function(e,t){teaweb.popup("/servers/server/settings/origins/updatePopup?originType="+t+"&"+this.vParams+"&originId="+e,{height:"27em",callback:function(e){teaweb.success("保存成功",function(){window.location.reload()})}})},deleteOrigin:function(e,t){let i=this;teaweb.confirm("确定要删除此源站吗？",function(){Tea.action("/servers/server/settings/origins/delete?"+i.vParams+"&originId="+e+"&originType="+t).post().success(function(){teaweb.success("删除成功",function(){window.location.reload()})})})}},template:`<div>
	<h3>主要源站 <a href="" @click.prevent="createPrimaryOrigin()">[添加主要源站]</a> </h3>
	<p class="comment" v-if="primaryOrigins.length == 0">暂时还没有主要源站。</p>
	<origin-list-table v-if="primaryOrigins.length > 0" :v-origins="vPrimaryOrigins" :v-origin-type="'primary'" @deleteOrigin="deleteOrigin" @updateOrigin="updateOrigin"></origin-list-table>

	<h3>备用源站 <a href="" @click.prevent="createBackupOrigin()">[添加备用源站]</a></h3>
	<p class="comment" v-if="backupOrigins.length == 0" :v-origins="primaryOrigins">暂时还没有备用源站。</p>
	<origin-list-table v-if="backupOrigins.length > 0" :v-origins="backupOrigins" :v-origin-type="'backup'" @deleteOrigin="deleteOrigin" @updateOrigin="updateOrigin"></origin-list-table>
</div>`}),Vue.component("origin-list-table",{props:["v-origins","v-origin-type"],data:function(){return{}},methods:{deleteOrigin:function(e){this.$emit("deleteOrigin",e,this.vOriginType)},updateOrigin:function(e){this.$emit("updateOrigin",e,this.vOriginType)}},template:`
<table class="ui table selectable">
	<thead>
		<tr>
			<th>源站地址</th>
			<th>权重</th>
			<th class="width10">状态</th>
			<th class="two op">操作</th>
		</tr>	
	</thead>
	<tr v-for="origin in vOrigins">
		<td :class="{disabled:!origin.isOn}"><a href="" @click.prevent="updateOrigin(origin.id)">{{origin.addr}} &nbsp;<i class="icon expand small"></i></a>
			<div style="margin-top: 0.3em" v-if="origin.name.length > 0 || origin.hasCert || (origin.host != null && origin.host.length > 0) || (origin.domains != null && origin.domains.length > 0)">
				<tiny-basic-label v-if="origin.name.length > 0">{{origin.name}}</tiny-basic-label>
				<tiny-basic-label v-if="origin.hasCert">证书</tiny-basic-label>
				<tiny-basic-label v-if="origin.host != null && origin.host.length > 0">主机名: {{origin.host}}</tiny-basic-label>
				<span v-if="origin.domains != null && origin.domains.length > 0"><tiny-basic-label v-for="domain in origin.domains">匹配: {{domain}}</tiny-basic-label></span>
			</div>
		</td>
		<td :class="{disabled:!origin.isOn}">{{origin.weight}}</td>
		<td>
			<label-on :v-is-on="origin.isOn"></label-on>
		</td>
		<td>
			<a href="" @click.prevent="updateOrigin(origin.id)">修改</a> &nbsp;
			<a href="" @click.prevent="deleteOrigin(origin.id)">删除</a>
		</td>
	</tr>
</table>`}),Vue.component("http-firewall-policy-selector",{props:["v-http-firewall-policy"],mounted:function(){let t=this;Tea.action("/servers/components/waf/count").post().success(function(e){t.count=e.data.count})},data:function(){return{count:0,firewallPolicy:this.vHttpFirewallPolicy}},methods:{remove:function(){this.firewallPolicy=null},select:function(){let t=this;teaweb.popup("/servers/components/waf/selectPopup",{callback:function(e){t.firewallPolicy=e.data.firewallPolicy}})},create:function(){let t=this;teaweb.popup("/servers/components/waf/createPopup",{height:"26em",callback:function(e){t.firewallPolicy=e.data.firewallPolicy}})}},template:`<div>
	<div v-if="firewallPolicy != null" class="ui label basic">
		<input type="hidden" name="httpFirewallPolicyId" :value="firewallPolicy.id"/>
		{{firewallPolicy.name}} &nbsp; <a :href="'/servers/components/waf/policy?firewallPolicyId=' + firewallPolicy.id" target="_blank" title="修改"><i class="icon pencil small"></i></a>&nbsp; <a href="" @click.prevent="remove()" title="删除"><i class="icon remove small"></i></a>
	</div>
	<div v-if="firewallPolicy == null">
		<span v-if="count > 0"><a href="" @click.prevent="select">[选择已有策略]</a> &nbsp; &nbsp; </span><a href="" @click.prevent="create">[创建新策略]</a>
	</div>
</div>`}),Vue.component("http-websocket-box",{props:["v-websocket-ref","v-websocket-config","v-is-location","v-is-group"],data:function(){let e=this.vWebsocketRef,t=(null==e&&(e={isPrior:!1,isOn:!1,websocketId:0}),this.vWebsocketConfig);return null==t?t={id:0,isOn:!1,handshakeTimeout:{count:30,unit:"second"},allowAllOrigins:!0,allowedOrigins:[],requestSameOrigin:!0,requestOrigin:""}:(null==t.handshakeTimeout&&(t.handshakeTimeout={count:30,unit:"second"}),null==t.allowedOrigins&&(t.allowedOrigins=[])),{websocketRef:e,websocketConfig:t,handshakeTimeoutCountString:t.handshakeTimeout.count.toString(),advancedVisible:!1}},watch:{handshakeTimeoutCountString:function(e){e=parseInt(e);!isNaN(e)&&0<=e?this.websocketConfig.handshakeTimeout.count=e:this.websocketConfig.handshakeTimeout.count=0}},methods:{isOn:function(){return(!this.vIsLocation&&!this.vIsGroup||this.websocketRef.isPrior)&&this.websocketRef.isOn},changeAdvancedVisible:function(e){this.advancedVisible=e},createOrigin:function(){let t=this;teaweb.popup("/servers/server/settings/websocket/createOrigin",{height:"12.5em",callback:function(e){t.websocketConfig.allowedOrigins.push(e.data.origin)}})},removeOrigin:function(e){this.websocketConfig.allowedOrigins.$remove(e)}},template:`<div>
	<input type="hidden" name="websocketRefJSON" :value="JSON.stringify(websocketRef)"/>
	<input type="hidden" name="websocketJSON" :value="JSON.stringify(websocketConfig)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="websocketRef" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="((!vIsLocation && !vIsGroup) || websocketRef.isPrior)">
			<tr>
				<td class="title">是否启用配置</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="websocketRef.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td class="color-border">允许所有来源域<em>（Origin）</em></td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="websocketConfig.allowAllOrigins"/>
						<label></label>
					</div>
					<p class="comment">选中表示允许所有的来源域。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn() && !websocketConfig.allowAllOrigins">
			<tr>
				<td class="color-border">允许的来源域列表<em>（Origin）</em></td>
				<td>
					<div v-if="websocketConfig.allowedOrigins.length > 0">
						<div class="ui label tiny" v-for="(origin, index) in websocketConfig.allowedOrigins">
							{{origin}} <a href="" title="删除" @click.prevent="removeOrigin(index)"><i class="icon remove"></i></a>
						</div>
						<div class="ui divider"></div>
					</div>
					<button class="ui button tiny" type="button" @click.prevent="createOrigin()">+</button>
					<p class="comment">只允许在列表中的来源域名访问Websocket服务。</p>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-show="isOn()"></more-options-tbody>
		<tbody v-show="isOn() && advancedVisible">
			<tr>
				<td class="color-border">是否传递请求来源域</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="websocketConfig.requestSameOrigin"/>
						<label></label>
					</div>
					<p class="comment">选中表示把接收到的请求中的<span class="ui label tiny">Origin</span>字段传递到源站。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn() && advancedVisible && !websocketConfig.requestSameOrigin">
			<tr>
				<td class="color-border">指定传递的来源域</td>
				<td>
					<input type="text" v-model="websocketConfig.requestOrigin" maxlength="200"/>
					<p class="comment">指定向源站传递的<span class="ui label tiny">Origin</span>字段值。</p>
				</td>
			</tr>
		</tbody>
		<!-- TODO 这个选项暂时保留 -->
		<tbody v-show="isOn() && false">
			<tr>
				<td>握手超时时间<em>（Handshake）</em></td>
				<td>
					<div class="ui fields inline">
						<div class="ui field">
							<input type="text" maxlength="10" v-model="handshakeTimeoutCountString" style="width:6em"/>
						</div>
						<div class="ui field">
							秒
						</div>
					</div>
					<p class="comment">0表示使用默认的时间设置。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`}),Vue.component("http-rewrite-rule-list",{props:["v-web-id","v-rewrite-rules"],mounted:function(){setTimeout(this.sort,1e3)},data:function(){let e=this.vRewriteRules;return{rewriteRules:e=null==e?[]:e}},methods:{updateRewriteRule:function(e){teaweb.popup("/servers/server/settings/rewrite/updatePopup?webId="+this.vWebId+"&rewriteRuleId="+e,{height:"26em",callback:function(){window.location.reload()}})},deleteRewriteRule:function(e){let t=this;teaweb.confirm("确定要删除此重写规则吗？",function(){Tea.action("/servers/server/settings/rewrite/delete").params({webId:t.vWebId,rewriteRuleId:e}).post().refresh()})},sort:function(){if(0!=this.rewriteRules.length){let t=this;sortTable(function(e){Tea.action("/servers/server/settings/rewrite/sort").post().params({webId:t.vWebId,rewriteRuleIds:e}).success(function(){teaweb.success("保存成功")})})}}},template:`<div>
	<div class="margin"></div>
	<p class="comment" v-if="rewriteRules.length == 0">暂时还没有重写规则。</p>
	<table class="ui table selectable" v-if="rewriteRules.length > 0" id="sortable-table">
		<thead>
			<tr>
				<th style="width:1em"></th>
				<th>匹配规则</th>
				<th>转发目标</th>
				<th>转发方式</th>
				<th class="two wide">状态</th>
				<th class="two op">操作</th>
			</tr>
		</thead>
		<tbody v-for="rule in rewriteRules" :v-id="rule.id">
			<tr>
				<td><i class="icon bars grey handle"></i></td>
				<td>{{rule.pattern}}
				<br/>
					<http-rewrite-labels-label class="ui label tiny" v-if="rule.isBreak">BREAK</http-rewrite-labels-label>
					<http-rewrite-labels-label class="ui label tiny" v-if="rule.mode == 'redirect' && rule.redirectStatus != 307">{{rule.redirectStatus}}</http-rewrite-labels-label>
					<http-rewrite-labels-label class="ui label tiny" v-if="rule.proxyHost.length > 0">Host: {{rule.proxyHost}}</http-rewrite-labels-label>
				</td>
				<td>{{rule.replace}}</td>
				<td>
					<span v-if="rule.mode == 'proxy'">隐式</span>
					<span v-if="rule.mode == 'redirect'">显示</span>
				</td>
				<td>
					<label-on :v-is-on="rule.isOn"></label-on>
				</td>
				<td>
					<a href="" @click.prevent="updateRewriteRule(rule.id)">修改</a> &nbsp;
					<a href="" @click.prevent="deleteRewriteRule(rule.id)">删除</a>
				</td>
			</tr>
		</tbody>
	</table>
	<p class="comment" v-if="rewriteRules.length > 0">拖动左侧的<i class="icon bars grey"></i>图标可以对重写规则进行排序。</p>

</div>`}),Vue.component("http-rewrite-labels-label",{props:["v-class"],template:'<span class="ui label tiny" :class="vClass" style="font-size:0.7em;padding:4px;margin-top:0.3em;margin-bottom:0.3em"><slot></slot></span>'}),Vue.component("server-name-box",{props:["v-server-names"],data:function(){let e=this.vServerNames;return{serverNames:e=null==e?[]:e,isSearching:!1,keyword:""}},methods:{addServerName:function(){window.UPDATING_SERVER_NAME=null;let t=this;teaweb.popup("/servers/addServerNamePopup",{callback:function(e){e=e.data.serverName;t.serverNames.push(e)}})},removeServerName:function(e){this.serverNames.$remove(e)},updateServerName:function(t,e){window.UPDATING_SERVER_NAME=e;let i=this;teaweb.popup("/servers/addServerNamePopup",{callback:function(e){e=e.data.serverName;Vue.set(i.serverNames,t,e)}})},showSearchBox:function(){if(this.isSearching=!this.isSearching,this.isSearching){let e=this;setTimeout(function(){e.$refs.keywordRef.focus()},200)}else this.keyword=""}},watch:{keyword:function(i){this.serverNames.forEach(function(e){if(0==i.length)e.isShowing=!0;else if(null==e.subNames||0==e.subNames.length)teaweb.match(e.name,i)||(e.isShowing=!1);else{let t=!1;e.subNames.forEach(function(e){teaweb.match(e,i)&&(t=!0)}),e.isShowing=t}})}},template:`<div>
	<input type="hidden" name="serverNames" :value="JSON.stringify(serverNames)"/>
	<div v-if="serverNames.length > 0">
		<div v-for="(serverName, index) in serverNames" class="ui label small basic" :class="{hidden: serverName.isShowing === false}">
			<em v-if="serverName.type != 'full'">{{serverName.type}}</em>  
			<span v-if="serverName.subNames == null || serverName.subNames.length == 0" :class="{disabled: serverName.isShowing === false}">{{serverName.name}}</span>
			<span v-else :class="{disabled: serverName.isShowing === false}">{{serverName.subNames[0]}}等{{serverName.subNames.length}}个域名</span>
			<a href="" title="修改" @click.prevent="updateServerName(index, serverName)"><i class="icon pencil small"></i></a> <a href="" title="删除" @click.prevent="removeServerName(index)"><i class="icon remove"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div class="ui fields inline">
	    <div class="ui field"><a href="" @click.prevent="addServerName()">[添加域名绑定]</a></div>
	    <div class="ui field" v-if="serverNames.length > 0"><span class="grey">|</span> </div>
	    <div class="ui field" v-if="serverNames.length > 0">
	        <a href="" @click.prevent="showSearchBox()" v-if="!isSearching"><i class="icon search small"></i></a>
	        <a href="" @click.prevent="showSearchBox()" v-if="isSearching"><i class="icon close small"></i></a>
        </div>
        <div class="ui field" v-if="isSearching">
            <input type="text" placeholder="搜索域名" ref="keywordRef" class="ui input tiny" v-model="keyword"/>
        </div>
    </div>
</div>`}),Vue.component("http-cache-stale-config",{props:["v-cache-stale-config"],data:function(){let e=this.vCacheStaleConfig;return{config:e=null==e?{isPrior:!1,isOn:!1,status:[],supportStaleIfErrorHeader:!0,life:{count:1,unit:"day"}}:e}},watch:{config:{deep:!0,handler:function(){this.$emit("change",this.config)}}},methods:{},template:`<table class="ui table definition selectable">
	<tbody>
		<tr>
			<td class="title">启用过时缓存</td>
			<td>
				<checkbox v-model="config.isOn"></checkbox>
				<p class="comment"><plus-label></plus-label>选中后，在更新缓存失败后会尝试读取过时的缓存。</p>
			</td>
		</tr>
		<tr v-show="config.isOn">
			<td>有效期</td>
			<td>
				<time-duration-box :v-value="config.life"></time-duration-box>
				<p class="comment">缓存在过期之后，仍然保留的时间。</p>
			</td>
		</tr>
		<tr v-show="config.isOn">
			<td>状态码</td>
			<td><http-status-box :v-status-list="config.status"></http-status-box>
				<p class="comment">在这些状态码出现时使用过时缓存，默认支持<code-label>50x</code-label>状态码。</p>
			</td>
		</tr>
		<tr v-show="config.isOn">
			<td>支持stale-if-error</td>
			<td>
				<checkbox v-model="config.supportStaleIfErrorHeader"></checkbox>
				<p class="comment">选中后，支持在Cache-Control中通过<code-label>stale-if-error</code-label>指定过时缓存有效期。</p>
			</td>
		</tr>
	</tbody>
</table>`}),Vue.component("firewall-syn-flood-config-viewer",{props:["v-syn-flood-config"],data:function(){let e=this.vSynFloodConfig;return{config:e=null==e?{isOn:!1,minAttempts:10,timeoutSeconds:600,ignoreLocal:!0}:e}},template:`<div>
	<span v-if="config.isOn">
		已启用 / <span>空连接次数：{{config.minAttempts}}次/分钟</span> / 封禁时间：{{config.timeoutSeconds}}秒 <span v-if="config.ignoreLocal">/ 忽略局域网访问</span>
	</span>
	<span v-else>未启用</span>
</div>`}),Vue.component("domains-box",{props:["v-domains"],data:function(){let e=this.vDomains;return{domains:e=null==e?[]:e,isAdding:!1,addingDomain:""}},methods:{add:function(){this.isAdding=!0;let e=this;setTimeout(function(){e.$refs.addingDomain.focus()},100)},confirm:function(){let t=this;if(this.addingDomain=this.addingDomain.replace(/\s/g,""),0==this.addingDomain.length)teaweb.warn("请输入要添加的域名",function(){t.$refs.addingDomain.focus()});else{if("~"==this.addingDomain[0]){var e=this.addingDomain.substring(1);try{new RegExp(e)}catch(e){return void teaweb.warn("正则表达式错误："+e.message,function(){t.$refs.addingDomain.focus()})}}this.domains.push(this.addingDomain),this.cancel()}},remove:function(e){this.domains.$remove(e)},cancel:function(){this.isAdding=!1,this.addingDomain=""}},template:`<div>
	<input type="hidden" name="domainsJSON" :value="JSON.stringify(domains)"/>
	<div v-if="domains.length > 0">
		<span class="ui label small basic" v-for="(domain, index) in domains">
			<span v-if="domain.length > 0 && domain[0] == '~'" class="grey" style="font-style: normal">[正则]</span>
			<span v-if="domain.length > 0 && domain[0] == '.'" class="grey" style="font-style: normal">[后缀]</span>
			<span v-if="domain.length > 0 && domain[0] == '*'" class="grey" style="font-style: normal">[泛域名]</span>
			{{domain}}
			&nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</span>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding">
		<div class="ui fields">
			<div class="ui field">
				<input type="text" v-model="addingDomain" @keyup.enter="confirm()" @keypress.enter.prevent="1" ref="addingDomain" placeholder="*.xxx.com" size="30"/>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>
				&nbsp; <a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
		<p class="comment">支持普通域名（<code-label>example.com</code-label>）、泛域名（<code-label>*.example.com</code-label>）、域名后缀（以点号开头，如<code-label>.example.com</code-label>）和正则表达式（以波浪号开头，如<code-label>~.*.example.com</code-label>）。</p>
		<div class="ui divider"></div>
	</div>
	<div style="margin-top: 0.5em" v-if="!isAdding">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`}),Vue.component("http-redirect-to-https-box",{props:["v-redirect-to-https-config","v-is-location"],data:function(){let e=this.vRedirectToHttpsConfig;return null==e?e={isPrior:!1,isOn:!1,host:"",port:0,status:0,onlyDomains:[],exceptDomains:[]}:(null==e.onlyDomains&&(e.onlyDomains=[]),null==e.exceptDomains&&(e.exceptDomains=[])),{redirectToHttpsConfig:e,portString:0<e.port?e.port.toString():"",moreOptionsVisible:!1,statusOptions:[{code:301,text:"Moved Permanently"},{code:308,text:"Permanent Redirect"},{code:302,text:"Found"},{code:303,text:"See Other"},{code:307,text:"Temporary Redirect"}]}},watch:{"redirectToHttpsConfig.status":function(){this.redirectToHttpsConfig.status=parseInt(this.redirectToHttpsConfig.status)},portString:function(e){e=parseInt(e);isNaN(e)?this.redirectToHttpsConfig.port=0:this.redirectToHttpsConfig.port=e}},methods:{changeMoreOptions:function(e){this.moreOptionsVisible=e},changeOnlyDomains:function(e){this.redirectToHttpsConfig.onlyDomains=e,this.$forceUpdate()},changeExceptDomains:function(e){this.redirectToHttpsConfig.exceptDomains=e,this.$forceUpdate()}},template:`<div>
	<input type="hidden" name="redirectToHTTPSJSON" :value="JSON.stringify(redirectToHttpsConfig)"/>
	
	<!-- Location -->
	<table class="ui table selectable definition" v-if="vIsLocation">
		<prior-checkbox :v-config="redirectToHttpsConfig"></prior-checkbox>
		<tbody v-show="redirectToHttpsConfig.isPrior">
			<tr>
				<td class="title">自动跳转到HTTPS</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="redirectToHttpsConfig.isOn"/>
						<label></label>
					</div>
					<p class="comment">开启后，所有HTTP的请求都会自动跳转到对应的HTTPS URL上，<more-options-angle @change="changeMoreOptions"></more-options-angle></p>
					
					<!--  TODO 如果已经设置了特殊设置，需要在界面上显示 -->
					<table class="ui table" v-show="moreOptionsVisible">
						<tr>
							<td class="title">状态码</td>
							<td>
								<select class="ui dropdown auto-width" v-model="redirectToHttpsConfig.status">
									<option value="0">[使用默认]</option>
									<option v-for="option in statusOptions" :value="option.code">{{option.code}} {{option.text}}</option>
								</select>
							</td>
						</tr>
						<tr>
							<td>域名或IP地址</td>
							<td>
								<input type="text" name="host" v-model="redirectToHttpsConfig.host"/>
								<p class="comment">默认和用户正在访问的域名或IP地址一致。</p>
							</td>
						</tr>
						<tr>
							<td>端口</td>
							<td>
								<input type="text" name="port" v-model="portString" maxlength="5" style="width:6em"/>
								<p class="comment">默认端口为443。</p>
							</td>
						</tr>
					</table>
				</td>
			</tr>	
		</tbody>
	</table>
	
	<!-- 非Location -->
	<div v-if="!vIsLocation">
		<div class="ui checkbox">
			<input type="checkbox" v-model="redirectToHttpsConfig.isOn"/>
			<label></label>
		</div>
		<p class="comment">开启后，所有HTTP的请求都会自动跳转到对应的HTTPS URL上，<more-options-angle @change="changeMoreOptions"></more-options-angle></p>
		
		<!--  TODO 如果已经设置了特殊设置，需要在界面上显示 -->
		<table class="ui table" v-show="moreOptionsVisible">
			<tr>
				<td class="title">状态码</td>
				<td>
					<select class="ui dropdown auto-width" v-model="redirectToHttpsConfig.status">
						<option value="0">[使用默认]</option>
						<option v-for="option in statusOptions" :value="option.code">{{option.code}} {{option.text}}</option>
					</select>
				</td>
			</tr>
			<tr>
				<td>跳转后域名或IP地址</td>
				<td>
					<input type="text" name="host" v-model="redirectToHttpsConfig.host"/>
					<p class="comment">默认和用户正在访问的域名或IP地址一致，不填写就表示使用当前的域名。</p>
				</td>
			</tr>
			<tr>
				<td>端口</td>
				<td>
					<input type="text" name="port" v-model="portString" maxlength="5" style="width:6em"/>
					<p class="comment">默认端口为443。</p>
				</td>
			</tr>
			<tr>
				<td>允许的域名</td>
				<td>
					<values-box :values="redirectToHttpsConfig.onlyDomains" @change="changeOnlyDomains"></values-box>
					<p class="comment">如果填写了允许的域名，那么只有这些域名可以自动跳转。</p>
				</td>
			</tr>
			<tr>
				<td>排除的域名</td>
				<td>
					<values-box :values="redirectToHttpsConfig.exceptDomains" @change="changeExceptDomains"></values-box>
					<p class="comment">如果填写了排除的域名，那么这些域名将不跳转。</p>
				</td>
			</tr>
		</table>
	</div>
	<div class="margin"></div>
</div>`}),Vue.component("http-firewall-actions-box",{props:["v-actions","v-firewall-policy","v-action-configs"],mounted:function(){let o=this;Tea.action("/servers/iplists/levelOptions").success(function(e){o.ipListLevels=e.data.levels}).post(),this.loadJS(function(){let n=document.getElementById("actions-box");Sortable.create(n,{draggable:".label",handle:".icon.handle",onStart:function(){o.cancel()},onUpdate:function(e){let t=n.getElementsByClassName("label"),i=[];for(let e=0;e<t.length;e++){var s=parseInt(t[e].getAttribute("data-index"));i.push(o.configs[s])}o.configs=i}})})},data:function(){null==this.vFirewallPolicy.inbound&&(this.vFirewallPolicy.inbound={}),null==this.vFirewallPolicy.inbound.groups&&(this.vFirewallPolicy.inbound.groups=[]);let t=0,e=[];null!=this.vActionConfigs&&(e=this.vActionConfigs).forEach(function(e){e.id=t++});var i=`<!DOCTYPE html>
<html>
<title>403 Forbidden</title>
<body>
<h1>403 Forbidden</h1>
<address>Request ID: \${requestId}.</address>
</body>
</html>`;return{id:t,actions:this.vActions,configs:e,isAdding:!1,editingIndex:-1,action:null,actionCode:"",actionOptions:{},ipListLevels:[],blockTimeout:"",blockScope:"global",captchaLife:"",captchaMaxFails:"",captchaFailBlockTimeout:"",get302Life:"",post307Life:"",recordIPType:"black",recordIPLevel:"critical",recordIPTimeout:"",recordIPListId:0,recordIPListName:"",tagTags:[],pageStatus:403,pageBody:i,defaultPageBody:i,goGroupName:"",goGroupId:0,goGroup:null,goSetId:0,goSetName:""}},watch:{actionCode:function(i){this.action=this.actions.$find(function(e,t){return t.code==i}),this.actionOptions={}},blockTimeout:function(e){e=parseInt(e),isNaN(e)?this.actionOptions.timeout=0:this.actionOptions.timeout=e},blockScope:function(e){this.actionOptions.scope=e},captchaLife:function(e){e=parseInt(e),isNaN(e)?this.actionOptions.life=0:this.actionOptions.life=e},captchaMaxFails:function(e){e=parseInt(e),isNaN(e)?this.actionOptions.maxFails=0:this.actionOptions.maxFails=e},captchaFailBlockTimeout:function(e){e=parseInt(e),isNaN(e)?this.actionOptions.failBlockTimeout=0:this.actionOptions.failBlockTimeout=e},get302Life:function(e){e=parseInt(e),isNaN(e)?this.actionOptions.life=0:this.actionOptions.life=e},post307Life:function(e){e=parseInt(e),isNaN(e)?this.actionOptions.life=0:this.actionOptions.life=e},recordIPType:function(e){this.recordIPListId=0},recordIPTimeout:function(e){e=parseInt(e),isNaN(e)?this.actionOptions.timeout=0:this.actionOptions.timeout=e},goGroupId:function(i){var e=this.vFirewallPolicy.inbound.groups.$find(function(e,t){return t.id==i});this.goGroup=e,this.goGroupName=null==e?"":e.name,this.goSetId=0,this.goSetName=""},goSetId:function(i){var e;null!=this.goGroup&&(null==(e=this.goGroup.sets.$find(function(e,t){return t.id==i}))?(this.goSetId=0,this.goSetName=""):this.goSetName=e.name)}},methods:{add:function(){this.action=null,this.actionCode="block",this.isAdding=!0,this.actionOptions={},this.blockTimeout="",this.blockScope="global",this.captchaLife="",this.captchaMaxFails="",this.captchaFailBlockTimeout="",this.get302Life="",this.post307Life="",this.recordIPLevel="critical",this.recordIPType="black",this.recordIPTimeout="",this.recordIPListId=0,this.recordIPListName="",this.tagTags=[],this.pageStatus=403,this.pageBody=this.defaultPageBody,this.goGroupName="",this.goGroupId=0,this.goGroup=null,this.goSetId=0,this.goSetName="";let i=this;this.action=this.vActions.$find(function(e,t){return t.code==i.actionCode}),this.scroll()},remove:function(e){this.isAdding=!1,this.editingIndex=-1,this.configs.$remove(e)},update:function(e,i){if(this.isAdding&&this.editingIndex==e)this.cancel();else{switch(this.add(),this.isAdding=!0,this.editingIndex=e,this.actionCode=i.code,i.code){case"block":this.blockTimeout="",(null!=i.options.timeout||0<i.options.timeout)&&(this.blockTimeout=i.options.timeout.toString()),null!=i.options.scope&&0<i.options.scope.length?this.blockScope=i.options.scope:this.blockScope="global";break;case"allow":case"log":break;case"captcha":this.captchaLife="",(null!=i.options.life||0<i.options.life)&&(this.captchaLife=i.options.life.toString()),this.captchaMaxFails="",(null!=i.options.maxFails||0<i.options.maxFails)&&(this.captchaMaxFails=i.options.maxFails.toString()),this.captchaFailBlockTimeout="",(null!=i.options.failBlockTimeout||0<i.options.failBlockTimeout)&&(this.captchaFailBlockTimeout=i.options.failBlockTimeout.toString());break;case"notify":break;case"get_302":this.get302Life="",(null!=i.options.life||0<i.options.life)&&(this.get302Life=i.options.life.toString());break;case"post_307":this.post307Life="",(null!=i.options.life||0<i.options.life)&&(this.post307Life=i.options.life.toString());break;case"record_ip":if(null!=i.options){this.recordIPLevel=i.options.level,this.recordIPType=i.options.type,0<i.options.timeout&&(this.recordIPTimeout=i.options.timeout.toString());let e=this;setTimeout(function(){e.recordIPListId=i.options.ipListId,e.recordIPListName=i.options.ipListName})}break;case"tag":this.tagTags=[],null!=i.options.tags&&(this.tagTags=i.options.tags);break;case"page":this.pageStatus=403,this.pageBody=this.defaultPageBody,null!=i.options.status&&(this.pageStatus=i.options.status),null!=i.options.body&&(this.pageBody=i.options.body);break;case"go_group":null!=i.options&&(this.goGroupName=i.options.groupName,this.goGroupId=i.options.groupId,this.goGroup=this.vFirewallPolicy.inbound.groups.$find(function(e,t){return t.id==i.options.groupId}));break;case"go_set":if(null!=i.options){this.goGroupName=i.options.groupName,this.goGroupId=i.options.groupId,this.goGroup=this.vFirewallPolicy.inbound.groups.$find(function(e,t){return t.id==i.options.groupId});let t=this;setTimeout(function(){var e;t.goSetId=i.options.setId,null!=t.goGroup&&null!=(e=t.goGroup.sets.$find(function(e,t){return t.id==i.options.setId}))&&(t.goSetName=e.name)})}}this.scroll()}},cancel:function(){this.isAdding=!1,this.editingIndex=-1},confirm:function(){if(null!=this.action){if(null==this.actionOptions&&(this.actionOptions={}),"record_ip"==this.actionCode){let e=parseInt(this.recordIPTimeout);if(isNaN(e)&&(e=0),this.recordIPListId<=0)return;this.actionOptions={type:this.recordIPType,level:this.recordIPLevel,timeout:e,ipListId:this.recordIPListId,ipListName:this.recordIPListName}}else if("tag"==this.actionCode){if(null==this.tagTags||0==this.tagTags.length)return;this.actionOptions={tags:this.tagTags}}else if("page"==this.actionCode){let e=this.pageStatus.toString();e=e.match(/^\d{3}$/)?parseInt(e):403,this.actionOptions={status:e,body:this.pageBody}}else if("go_group"==this.actionCode){let e=this.goGroupId;if("string"==typeof e&&(e=parseInt(e),isNaN(e)&&(e=0)),e<=0)return;this.actionOptions={groupId:e.toString(),groupName:this.goGroupName}}else if("go_set"==this.actionCode){let e=this.goGroupId,t=("string"==typeof e&&(e=parseInt(e),isNaN(e)&&(e=0)),this.goSetId);if("string"==typeof t&&(t=parseInt(t),isNaN(t)&&(t=0)),t<=0)return;this.actionOptions={groupId:e.toString(),groupName:this.goGroupName,setId:t.toString(),setName:this.goSetName}}let e={};for(var t in this.actionOptions)this.actionOptions.hasOwnProperty(t)&&(e[t]=this.actionOptions[t]);-1<this.editingIndex?this.configs[this.editingIndex]={id:this.configs[this.editingIndex].id,code:this.actionCode,name:this.action.name,options:e}:this.configs.push({id:this.id++,code:this.actionCode,name:this.action.name,options:e}),this.cancel()}},removeRecordIPList:function(){this.recordIPListId=0},selectRecordIPList:function(){let t=this;teaweb.popup("/servers/iplists/selectPopup?type="+this.recordIPType,{width:"50em",height:"30em",callback:function(e){t.recordIPListId=e.data.list.id,t.recordIPListName=e.data.list.name}})},changeTags:function(e){this.tagTags=e},loadJS:function(t){if("undefined"!=typeof Sortable)t();else{let e=document.createElement("script");e.setAttribute("src","/js/sortable.min.js"),e.addEventListener("load",function(){t()}),document.head.appendChild(e)}},scroll:function(){setTimeout(function(){let e=document.getElementsByClassName("main");0<e.length&&e[0].scrollTo(0,1e3)},10)}},template:`<div>
	<input type="hidden" name="actionsJSON" :value="JSON.stringify(configs)"/>
	<div v-show="configs.length > 0" style="margin-bottom: 0.5em" id="actions-box"> 
		<div v-for="(config, index) in configs" :data-index="index" :key="config.id" class="ui label small basic" :class="{blue: index == editingIndex}" style="margin-bottom: 0.4em">
			{{config.name}} <span class="small">({{config.code.toUpperCase()}})</span> 
			
			<!-- block -->
			<span v-if="config.code == 'block' && config.options.timeout > 0">：有效期{{config.options.timeout}}秒</span>
			
			<!-- captcha -->
			<span v-if="config.code == 'captcha' && config.options.life > 0">：有效期{{config.options.life}}秒
				<span v-if="config.options.maxFails > 0"> / 最多失败{{config.options.maxFails}}次</span>
			</span>
			
			<!-- get 302 -->
			<span v-if="config.code == 'get_302' && config.options.life > 0">：有效期{{config.options.life}}秒</span>
			
			<!-- post 307 -->
			<span v-if="config.code == 'post_307' && config.options.life > 0">：有效期{{config.options.life}}秒</span>
			
			<!-- record_ip -->
			<span v-if="config.code == 'record_ip'">：{{config.options.ipListName}}</span>
			
			<!-- tag -->
			<span v-if="config.code == 'tag'">：{{config.options.tags.join(", ")}}</span>
			
			<!-- page -->
			<span v-if="config.code == 'page'">：[{{config.options.status}}]</span>
			
			<!-- go_group -->
			<span v-if="config.code == 'go_group'">：{{config.options.groupName}}</span>
			
			<!-- go_set -->
			<span v-if="config.code == 'go_set'">：{{config.options.groupName}} / {{config.options.setName}}</span>
			
			<!-- 范围 -->
			<span v-if="config.options.scope != null && config.options.scope.length > 0" class="small grey">
				&nbsp; 
				<span v-if="config.options.scope == 'global'">[所有服务]</span>
				<span v-if="config.options.scope == 'service'">[当前服务]</span>
			</span>
			
			<!-- 操作按钮 -->
			 &nbsp; <a href="" title="修改" @click.prevent="update(index, config)"><i class="icon pencil small"></i></a> &nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a> &nbsp; <a href="" title="拖动改变顺序"><i class="icon bars handle"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div style="margin-bottom: 0.5em" v-if="isAdding">
		<table class="ui table" :class="{blue: editingIndex > -1}">
			<tr>
				<td class="title">动作类型 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="actionCode">
						<option v-for="action in actions" :value="action.code">{{action.name}} ({{action.code.toUpperCase()}})</option>
					</select>
					<p class="comment" v-if="action != null && action.description.length > 0">{{action.description}}</p>
				</td>
			</tr>
			
			<!-- block -->
			<tr v-if="actionCode == 'block'">
				<td>封锁时间</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 5em" maxlength="9" v-model="blockTimeout" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">秒</span>
					</div>
				</td>
			</tr>
			<tr v-if="actionCode == 'block'">
				<td>封锁范围</td>
				<td>
					<select class="ui dropdown auto-width" v-model="blockScope">
						<option value="service">当前服务</option>
						<option value="global">所有服务</option>
					</select>
					<p class="comment" v-if="blockScope == 'service'">只封锁用户对当前网站服务的访问，其他服务不受影响。</p>
					<p class="comment" v-if="blockScope =='global'">封锁用户对所有网站服务的访问。</p>
				</td>
			</tr>
			
			<!-- captcha -->
			<tr v-if="actionCode == 'captcha'">
				<td>有效时间</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 5em" maxlength="9" v-model="captchaLife" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">秒</span>
					</div>
					<p class="comment">验证通过后在这个时间内不再验证；如果为空或者为0表示默认。</p>
				</td>
			</tr>
			<tr v-if="actionCode == 'captcha'">
				<td>最多失败次数</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 5em" maxlength="9" v-model="captchaMaxFails" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">次</span>
					</div>
					<p class="comment">允许用户失败尝试的最多次数，超过这个次数将被自动加入黑名单；如果为空或者为0表示默认。</p>
				</td>
			</tr>
			<tr v-if="actionCode == 'captcha'">
				<td>失败拦截时间</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 5em" maxlength="9" v-model="captchaFailBlockTimeout" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">秒</span>
					</div>
					<p class="comment">在达到最多失败次数（大于0）时，自动拦截的时间；如果为空或者为0表示默认。</p>
				</td>
			</tr>
			
			<!-- get_302 -->
			<tr v-if="actionCode == 'get_302'">
				<td>有效时间</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 5em" maxlength="9" v-model="get302Life" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">秒</span>
					</div>
					<p class="comment">验证通过后在这个时间内不再验证。</p>
				</td>
			</tr>
			
			<!-- post_307 -->
			<tr v-if="actionCode == 'post_307'">
				<td>有效时间</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 5em" maxlength="9" v-model="post307Life" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">秒</span>
					</div>
					<p class="comment">验证通过后在这个时间内不再验证。</p>
				</td>
			</tr>
			
			<!-- record_ip -->
			<tr v-if="actionCode == 'record_ip'">
				<td>IP名单类型 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="recordIPType">
					<option value="black">黑名单</option>
					<option value="white">白名单</option>
					</select>
				</td>
			</tr>
			<tr v-if="actionCode == 'record_ip'">
				<td>选择IP名单 *</td>
				<td>
					<div v-if="recordIPListId > 0" class="ui label basic small">{{recordIPListName}} <a href="" @click.prevent="removeRecordIPList"><i class="icon remove small"></i></a></div>
					<button type="button" class="ui button tiny" @click.prevent="selectRecordIPList">+</button>
					<p class="comment">如不选择，则自动添加到当前策略的IP名单中。</p>
				</td>
			</tr>
			<tr v-if="actionCode == 'record_ip'">
				<td>级别</td>
				<td>
					<select class="ui dropdown auto-width" v-model="recordIPLevel">
						<option v-for="level in ipListLevels" :value="level.code">{{level.name}}</option>
					</select>
				</td>
			</tr>
			<tr v-if="actionCode == 'record_ip'">
				<td>超时时间</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" style="width: 6em" maxlength="9" v-model="recordIPTimeout" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<span class="ui label">秒</span>
					</div>
					<p class="comment">0表示不超时。</p>
				</td>
			</tr>
			
			<!-- tag -->
			<tr v-if="actionCode == 'tag'">
				<td>标签 *</td>
				<td>
					<values-box @change="changeTags" :values="tagTags"></values-box>
				</td>
			</tr>
			
			<!-- page -->
			<tr v-if="actionCode == 'page'">
				<td>状态码 *</td>
				<td><input type="text" style="width: 4em" maxlength="3" v-model="pageStatus"/></td>
			</tr>
			<tr v-if="actionCode == 'page'">
				<td>网页内容</td>
				<td>
					<textarea v-model="pageBody"></textarea>
				</td>
			</tr>
			
			<!-- 规则分组 -->
			<tr v-if="actionCode == 'go_group'">
				<td>下一个分组 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="goGroupId">
						<option value="0">[选择分组]</option>
						<option v-for="group in vFirewallPolicy.inbound.groups" :value="group.id">{{group.name}}</option>
					</select>
				</td>
			</tr>
			
			<!-- 规则集 -->
			<tr v-if="actionCode == 'go_set'">
				<td>下一个分组 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="goGroupId">
						<option value="0">[选择分组]</option>
						<option v-for="group in vFirewallPolicy.inbound.groups" :value="group.id">{{group.name}}</option>
					</select>
				</td>
			</tr>
			<tr v-if="actionCode == 'go_set' && goGroup != null">
				<td>下一个规则集 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="goSetId">
						<option value="0">[选择规则集]</option>
						<option v-for="set in goGroup.sets" :value="set.id">{{set.name}}</option>
					</select>
				</td>
			</tr>
		</table>
		<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button> &nbsp;
		<a href="" @click.prevent="cancel" title="取消"><i class="icon remove small"></i></a>
	</div>
	<div v-if="!isAdding">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
	<p class="comment">系统总是会先执行记录日志、标签等不会修改请求的动作，再执行阻止、验证码等可能改变请求的动作。</p>
</div>`}),Vue.component("http-auth-config-box",{props:["v-auth-config","v-is-location"],data:function(){let e=this.vAuthConfig;return null==(e=null==e?{isPrior:!1,isOn:!1}:e).policyRefs&&(e.policyRefs=[]),{authConfig:e}},methods:{isOn:function(){return(!this.vIsLocation||this.authConfig.isPrior)&&this.authConfig.isOn},add:function(){let t=this;teaweb.popup("/servers/server/settings/access/createPopup",{callback:function(e){t.authConfig.policyRefs.push(e.data.policyRef)},height:"28em"})},update:function(e,t){teaweb.popup("/servers/server/settings/access/updatePopup?policyId="+t,{callback:function(e){teaweb.success("保存成功",function(){teaweb.reload()})},height:"28em"})},remove:function(e){this.authConfig.policyRefs.$remove(e)},methodName:function(e){switch(e){case"basicAuth":return"BasicAuth";case"subRequest":return"子请求"}return""}},template:`<div>
<input type="hidden" name="authJSON" :value="JSON.stringify(authConfig)"/> 
<table class="ui table selectable definition">
	<prior-checkbox :v-config="authConfig" v-if="vIsLocation"></prior-checkbox>
	<tbody v-show="!vIsLocation || authConfig.isPrior">
		<tr>
			<td class="title">启用认证</td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" v-model="authConfig.isOn"/>
					<label></label>
				</div>
			</td>
		</tr>
	</tbody>
</table>
<div class="margin"></div>
<!-- 认证方式 -->
<div v-show="isOn()">
	<h4>认证方式</h4>
	<table class="ui table selectable celled" v-show="authConfig.policyRefs.length > 0">
		<thead>
			<tr>
				<th class="three wide">名称</th>
				<th class="three wide">认证方法</th>
				<th>参数</th>
				<th class="two wide">状态</th>
				<th class="two op">操作</th>
			</tr>
		</thead>
		<tbody v-for="(ref, index) in authConfig.policyRefs" :key="ref.authPolicyId">
			<tr>
				<td>{{ref.authPolicy.name}}</td>
				<td>
					{{methodName(ref.authPolicy.type)}}
				</td>
				<td>
					<span v-if="ref.authPolicy.type == 'basicAuth'">{{ref.authPolicy.params.users.length}}个用户</span>
					<span v-if="ref.authPolicy.type == 'subRequest'">
						<span v-if="ref.authPolicy.params.method.length > 0" class="grey">[{{ref.authPolicy.params.method}}]</span>
						{{ref.authPolicy.params.url}}
					</span>
				</td>
				<td>
					<label-on :v-is-on="ref.authPolicy.isOn"></label-on>
				</td>
				<td>
					<a href="" @click.prevent="update(index, ref.authPolicyId)">修改</a> &nbsp;
					<a href="" @click.prevent="remove(index)">删除</a>
				</td>
			</tr>
		</tbody>
	</table>
	<button class="ui button small" type="button" @click.prevent="add">+添加认证方式</button>
</div>
<div class="margin"></div>
</div>`}),Vue.component("user-selector",{mounted:function(){let t=this;Tea.action("/servers/users/options").post().success(function(e){t.users=e.data.users})},props:["v-user-id"],data:function(){let e=this.vUserId;return{users:[],userId:e=null==e?0:e}},watch:{userId:function(e){this.$emit("change",e)}},template:`<div>
	<select class="ui dropdown auto-width" name="userId" v-model="userId">
		<option value="0">[选择用户]</option>
		<option v-for="user in users" :value="user.id">{{user.fullname}} ({{user.username}})</option>
	</select>
</div>`}),Vue.component("uam-config-box",{props:["v-uam-config","v-is-location","v-is-group"],data:function(){let e=this.vUamConfig;return{config:e=null==e?{isPrior:!1,isOn:!1}:e}},template:`<div>
<input type="hidden" name="uamJSON" :value="JSON.stringify(config)"/>
<table class="ui table definition selectable">
	<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
	<tbody v-show="((!vIsLocation && !vIsGroup) || config.isPrior)">
		<tr>
			<td class="title">启用5秒盾</td>
			<td>
				<checkbox v-model="config.isOn"></checkbox>
				<p class="comment"><plus-label></plus-label>启用后，访问网站时，自动检查浏览器环境，阻止非正常访问。</p>
			</td>
		</tr>
	</tbody>
</table>
<div class="margin"></div>
</div>`}),Vue.component("http-header-policy-box",{props:["v-request-header-policy","v-request-header-ref","v-response-header-policy","v-response-header-ref","v-params","v-is-location","v-is-group","v-has-group-request-config","v-has-group-response-config","v-group-setting-url"],data:function(){let e="response";"#request"==window.location.hash&&(e="request");let t=this.vRequestHeaderRef,i=(null==t&&(t={isPrior:!1,isOn:!0,headerPolicyId:0}),this.vResponseHeaderRef),s=(null==i&&(i={isPrior:!1,isOn:!0,headerPolicyId:0}),[]),n=[];var o=this.vRequestHeaderPolicy;null!=o&&(null!=o.setHeaders&&(s=o.setHeaders),null!=o.deleteHeaders&&(n=o.deleteHeaders));let a=[],l=[];o=this.vResponseHeaderPolicy;return null!=o&&(null!=o.setHeaders&&(a=o.setHeaders),null!=o.deleteHeaders&&(l=o.deleteHeaders)),{type:e,typeName:"request"==e?"请求":"响应",requestHeaderRef:t,responseHeaderRef:i,requestSettingHeaders:s,requestDeletingHeaders:n,responseSettingHeaders:a,responseDeletingHeaders:l}},methods:{selectType:function(e){this.type=e,window.location.hash="#"+e,window.location.reload()},addSettingHeader:function(e){teaweb.popup("/servers/server/settings/headers/createSetPopup?"+this.vParams+"&headerPolicyId="+e+"&type="+this.type,{callback:function(){teaweb.successRefresh("保存成功")}})},addDeletingHeader:function(e,t){teaweb.popup("/servers/server/settings/headers/createDeletePopup?"+this.vParams+"&headerPolicyId="+e+"&type="+t,{callback:function(){teaweb.successRefresh("保存成功")}})},updateSettingPopup:function(e,t){teaweb.popup("/servers/server/settings/headers/updateSetPopup?"+this.vParams+"&headerPolicyId="+e+"&headerId="+t+"&type="+this.type,{callback:function(){teaweb.successRefresh("保存成功")}})},deleteDeletingHeader:function(e,t){teaweb.confirm("确定要删除'"+t+"'吗？",function(){Tea.action("/servers/server/settings/headers/deleteDeletingHeader").params({headerPolicyId:e,headerName:t}).post().refresh()})},deleteHeader:function(e,t,i){teaweb.confirm("确定要删除此Header吗？",function(){this.$post("/servers/server/settings/headers/delete").params({headerPolicyId:e,type:t,headerId:i}).refresh()})}},template:`<div>
	<div class="ui menu tabular small">
		<a class="item" :class="{active:type == 'response'}" @click.prevent="selectType('response')">响应Header<span v-if="responseSettingHeaders.length > 0">({{responseSettingHeaders.length}})</span></a>
		<a class="item" :class="{active:type == 'request'}" @click.prevent="selectType('request')">请求Header<span v-if="requestSettingHeaders.length > 0">({{requestSettingHeaders.length}})</span></a>
	</div>
	
	<div class="margin"></div>
	
	<input type="hidden" name="type" :value="type"/>
	
	<!-- 请求 -->
	<div v-if="(vIsLocation || vIsGroup) && type == 'request'">
		<input type="hidden" name="requestHeaderJSON" :value="JSON.stringify(requestHeaderRef)"/>
		<table class="ui table definition selectable">
			<prior-checkbox :v-config="requestHeaderRef"></prior-checkbox>
		</table>
		<submit-btn></submit-btn>
	</div>
	
	<div v-if="((!vIsLocation && !vIsGroup) || requestHeaderRef.isPrior) && type == 'request'">
		<div v-if="vHasGroupRequestConfig">
        	<div class="margin"></div>
        	<warning-message>由于已经在当前<a :href="vGroupSettingUrl + '#request'">服务分组</a>中进行了对应的配置，在这里的配置将不会生效。</warning-message>
    	</div>
    	<div :class="{'opacity-mask': vHasGroupRequestConfig}">
		<h3>设置请求Header <a href="" @click.prevent="addSettingHeader(vRequestHeaderPolicy.id)">[添加新Header]</a></h3>
			<p class="comment" v-if="requestSettingHeaders.length == 0">暂时还没有Header。</p>
			<table class="ui table selectable celled" v-if="requestSettingHeaders.length > 0">
				<thead>
					<tr>
						<th>名称</th>
						<th>值</th>
						<th class="two op">操作</th>
					</tr>
				</thead>
				<tr v-for="header in requestSettingHeaders">
					<td class="five wide">
						{{header.name}}
						<div>
							<span v-if="header.status != null && header.status.codes != null && !header.status.always"><grey-label v-for="code in header.status.codes" :key="code">{{code}}</grey-label></span>
							<span v-if="header.methods != null && header.methods.length > 0"><grey-label v-for="method in header.methods" :key="method">{{method}}</grey-label></span>
							<span v-if="header.domains != null && header.domains.length > 0"><grey-label v-for="domain in header.domains" :key="domain">{{domain}}</grey-label></span>
							<grey-label v-if="header.shouldAppend">附加</grey-label>
							<grey-label v-if="header.disableRedirect">跳转禁用</grey-label>
							<grey-label v-if="header.shouldReplace && header.replaceValues != null && header.replaceValues.length > 0">替换</grey-label>
						</div>
					</td>
					<td>{{header.value}}</td>
					<td><a href="" @click.prevent="updateSettingPopup(vRequestHeaderPolicy.id, header.id)">修改</a> &nbsp; <a href="" @click.prevent="deleteHeader(vRequestHeaderPolicy.id, 'setHeader', header.id)">删除</a> </td>
				</tr>
			</table>
			
			<h3>删除请求Header</h3>
			<p class="comment">这里可以设置需要从请求中删除的Header。</p>
			
			<table class="ui table definition selectable">
				<td class="title">需要删除的Header</td>
				<td>
					<div v-if="requestDeletingHeaders.length > 0">
						<div class="ui label small basic" v-for="headerName in requestDeletingHeaders">{{headerName}} <a href=""><i class="icon remove" title="删除" @click.prevent="deleteDeletingHeader(vRequestHeaderPolicy.id, headerName)"></i></a> </div>
						<div class="ui divider" ></div>
					</div>
					<button class="ui button small" type="button" @click.prevent="addDeletingHeader(vRequestHeaderPolicy.id, 'request')">+</button>
				</td>
			</table>
		</div>			
	</div>
	
	<!-- 响应 -->
	<div v-if="(vIsLocation || vIsGroup) && type == 'response'">
		<input type="hidden" name="responseHeaderJSON" :value="JSON.stringify(responseHeaderRef)"/>
		<table class="ui table definition selectable">
			<prior-checkbox :v-config="responseHeaderRef"></prior-checkbox>
		</table>
		<submit-btn></submit-btn>
	</div>
	
	<div v-if="((!vIsLocation && !vIsGroup) || responseHeaderRef.isPrior) && type == 'response'">
		<div v-if="vHasGroupResponseConfig">
        	<div class="margin"></div>
        	<warning-message>由于已经在当前<a :href="vGroupSettingUrl + '#response'">服务分组</a>中进行了对应的配置，在这里的配置将不会生效。</warning-message>
    	</div>
    	<div :class="{'opacity-mask': vHasGroupResponseConfig}">
			<h3>设置响应Header <a href="" @click.prevent="addSettingHeader(vResponseHeaderPolicy.id)">[添加新Header]</a></h3>
			<p class="comment" style="margin-top: 0; padding-top: 0">将会覆盖已有的同名Header。</p>
			<p class="comment" v-if="responseSettingHeaders.length == 0">暂时还没有Header。</p>
			<table class="ui table selectable celled" v-if="responseSettingHeaders.length > 0">
				<thead>
					<tr>
						<th>名称</th>
						<th>值</th>
						<th class="two op">操作</th>
					</tr>
				</thead>
				<tr v-for="header in responseSettingHeaders">
					<td class="five wide">
						{{header.name}}
						<div>
							<span v-if="header.status != null && header.status.codes != null && !header.status.always"><grey-label v-for="code in header.status.codes" :key="code">{{code}}</grey-label></span>
							<span v-if="header.methods != null && header.methods.length > 0"><grey-label v-for="method in header.methods" :key="method">{{method}}</grey-label></span>
							<span v-if="header.domains != null && header.domains.length > 0"><grey-label v-for="domain in header.domains" :key="domain">{{domain}}</grey-label></span>
							<grey-label v-if="header.shouldAppend">附加</grey-label>
							<grey-label v-if="header.disableRedirect">跳转禁用</grey-label>
							<grey-label v-if="header.shouldReplace && header.replaceValues != null && header.replaceValues.length > 0">替换</grey-label>
						</div>
					</td>
					<td>{{header.value}}</td>
					<td><a href="" @click.prevent="updateSettingPopup(vResponseHeaderPolicy.id, header.id)">修改</a> &nbsp; <a href="" @click.prevent="deleteHeader(vResponseHeaderPolicy.id, 'setHeader', header.id)">删除</a> </td>
				</tr>
			</table>
			
			<h3>删除响应Header</h3>
			<p class="comment">这里可以设置需要从响应中删除的Header。</p>
			
			<table class="ui table definition selectable">
				<td class="title">需要删除的Header</td>
				<td>
					<div v-if="responseDeletingHeaders.length > 0">
						<div class="ui label small basic" v-for="headerName in responseDeletingHeaders">{{headerName}} <a href=""><i class="icon remove" title="删除" @click.prevent="deleteDeletingHeader(vResponseHeaderPolicy.id, headerName)"></i></a> </div>
						<div class="ui divider" ></div>
					</div>
					<button class="ui button small" type="button" @click.prevent="addDeletingHeader(vResponseHeaderPolicy.id, 'response')">+</button>
				</td>
			</table>
		</div>			
	</div>
	<div class="margin"></div>
</div>`}),Vue.component("http-common-config-box",{props:["v-common-config"],data:function(){let e=this.vCommonConfig;return{config:e=null==e?{mergeSlashes:!1}:e}},template:`<div>
	<table class="ui table definition selectable">
		<tr>
			<td class="title">合并重复的路径分隔符</td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" name="mergeSlashes" value="1" v-model="config.mergeSlashes"/>
					<label></label>
				</div>
				<p class="comment">合并URL中重复的路径分隔符为一个，比如<code-label>//hello/world</code-label>中的<code-label>//</code-label>。</p>
			</td>
		</tr>
	</table>
	<div class="margin"></div>
</div>`}),Vue.component("http-cache-policy-selector",{props:["v-cache-policy"],mounted:function(){let t=this;Tea.action("/servers/components/cache/count").post().success(function(e){t.count=e.data.count})},data:function(){return{count:0,cachePolicy:this.vCachePolicy}},methods:{remove:function(){this.cachePolicy=null},select:function(){let t=this;teaweb.popup("/servers/components/cache/selectPopup",{callback:function(e){t.cachePolicy=e.data.cachePolicy}})},create:function(){let t=this;teaweb.popup("/servers/components/cache/createPopup",{height:"26em",callback:function(e){t.cachePolicy=e.data.cachePolicy}})}},template:`<div>
	<div v-if="cachePolicy != null" class="ui label basic">
		<input type="hidden" name="cachePolicyId" :value="cachePolicy.id"/>
		{{cachePolicy.name}} &nbsp; <a :href="'/servers/components/cache/policy?cachePolicyId=' + cachePolicy.id" target="_blank" title="修改"><i class="icon pencil small"></i></a>&nbsp; <a href="" @click.prevent="remove()" title="删除"><i class="icon remove small"></i></a>
	</div>
	<div v-if="cachePolicy == null">
		<span v-if="count > 0"><a href="" @click.prevent="select">[选择已有策略]</a> &nbsp; &nbsp; </span><a href="" @click.prevent="create">[创建新策略]</a>
	</div>
</div>`}),Vue.component("http-pages-and-shutdown-box",{props:["v-pages","v-shutdown-config","v-is-location"],data:function(){let e=[],t=(null!=this.vPages&&(e=this.vPages),{isPrior:!1,isOn:!1,bodyType:"url",url:"",body:"",status:0}),i=(null!=this.vShutdownConfig&&(null==this.vShutdownConfig.body&&(this.vShutdownConfig.body=""),null==this.vShutdownConfig.bodyType&&(this.vShutdownConfig.bodyType="url"),t=this.vShutdownConfig),"");return 0<t.status&&(i=t.status.toString()),{pages:e,shutdownConfig:t,shutdownStatus:i}},watch:{shutdownStatus:function(e){e=parseInt(e);!isNaN(e)&&0<e&&e<1e3?this.shutdownConfig.status=e:this.shutdownConfig.status=0}},methods:{addPage:function(){let t=this;teaweb.popup("/servers/server/settings/pages/createPopup",{height:"26em",callback:function(e){t.pages.push(e.data.page)}})},updatePage:function(t,e){let i=this;teaweb.popup("/servers/server/settings/pages/updatePopup?pageId="+e,{height:"26em",callback:function(e){Vue.set(i.pages,t,e.data.page)}})},removePage:function(e){let t=this;teaweb.confirm("确定要移除此页面吗？",function(){t.pages.$remove(e)})},addShutdownHTMLTemplate:function(){this.shutdownConfig.body=`<!DOCTYPE html>
<html>
<head>
	<title>升级中</title>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
</head>
<body>

<h1>网站升级中</h1>
<p>为了给您提供更好的服务，我们正在升级网站，请稍后重新访问。</p>

<address>Request ID: \${requestId}.</address>

</body>
</html>`}},template:`<div>
<input type="hidden" name="pagesJSON" :value="JSON.stringify(pages)"/>
<input type="hidden" name="shutdownJSON" :value="JSON.stringify(shutdownConfig)"/>
<table class="ui table selectable definition">
	<tr>
		<td class="title">自定义页面</td>
		<td>
			<div v-if="pages.length > 0">
				<div class="ui label small basic" v-for="(page,index) in pages">
					{{page.status}} -&gt; <span v-if="page.bodyType == 'url'">{{page.url}}</span><span v-if="page.bodyType == 'html'">[HTML内容]</span> <a href="" title="修改" @click.prevent="updatePage(index, page.id)"><i class="icon pencil small"></i></a> <a href="" title="删除" @click.prevent="removePage(index)"><i class="icon remove"></i></a>
				</div>
				<div class="ui divider"></div>
			</div>
			<div>
				<button class="ui button small" type="button" @click.prevent="addPage()">+</button>
			</div>
			<p class="comment">根据响应状态码返回一些自定义页面，比如404，500等错误页面。</p>
		</td>
	</tr>	
	<tr>
		<td>临时关闭页面</td>
		<td>
			<div>
				<table class="ui table selectable definition">
					<prior-checkbox :v-config="shutdownConfig" v-if="vIsLocation"></prior-checkbox>
					<tbody v-show="!vIsLocation || shutdownConfig.isPrior">
						<tr>
							<td class="title">是否开启</td>
							<td>
								<div class="ui checkbox">
									<input type="checkbox" value="1" v-model="shutdownConfig.isOn" />
									<label></label>
								</div>
							</td>
						</tr>
					</tbody>
					<tbody v-show="(!vIsLocation || shutdownConfig.isPrior) && shutdownConfig.isOn">
						<tr>
							<td>内容类型 *</td>
							<td>
								<select class="ui dropdown auto-width" v-model="shutdownConfig.bodyType">
									<option value="url">读取URL</option>
									<option value="html">HTML</option>
								</select>
							</td>
						</tr>
						<tr v-show="shutdownConfig.bodyType == 'url'">
							<td class="title">页面URL *</td>
							<td>
								<input type="text" v-model="shutdownConfig.url" placeholder="页面文件路径或一个完整URL"/>
								<p class="comment">页面文件是相对于节点安装目录的页面文件比如pages/40x.html，或者一个完整的URL。</p>
							</td>
						</tr>
						<tr v-show="shutdownConfig.bodyType == 'html'">
							<td>HTML *</td>
							<td>
								<textarea name="body" ref="shutdownHTMLBody" v-model="shutdownConfig.body"></textarea>
								<p class="comment"><a href="" @click.prevent="addShutdownHTMLTemplate">[使用模板]</a>。填写页面的HTML内容，支持请求变量。</p>
							</td>
						</tr>
						<tr>
							<td>状态码</td>
							<td><input type="text" size="3" maxlength="3" name="shutdownStatus" style="width:5.2em" placeholder="状态码" v-model="shutdownStatus"/></td>
						</tr>
					</tbody>
				</table>
				<p class="comment">开启临时关闭页面时，所有请求都会直接显示此页面。可用于临时升级网站或者禁止用户访问某个网页。</p>
			</div>
		</td>
	</tr>
</table>
<div class="ui margin"></div>
</div>`}),Vue.component("http-compression-config-box",{props:["v-compression-config","v-is-location","v-is-group"],mounted:function(){let e=this;sortLoad(function(){e.initSortableTypes()})},data:function(){let t=this.vCompressionConfig,e=(null==(t=null==t?{isPrior:!1,isOn:!1,useDefaultTypes:!0,types:["brotli","gzip","deflate"],level:5,decompressData:!1,gzipRef:null,deflateRef:null,brotliRef:null,minLength:{count:0,unit:"kb"},maxLength:{count:0,unit:"kb"},mimeTypes:["text/*","application/*","font/*"],extensions:[".js",".json",".html",".htm",".xml",".css",".woff2",".txt"],conds:null}:t).types&&(t.types=[]),null==t.mimeTypes&&(t.mimeTypes=[]),null==t.extensions&&(t.extensions=[]),[{name:"Gzip",code:"gzip",isOn:!0},{name:"Deflate",code:"deflate",isOn:!0},{name:"Brotli",code:"brotli",isOn:!0}]),i=[];return t.types.forEach(function(t){e.forEach(function(e){t==e.code&&(e.isOn=!0,i.push(e))})}),e.forEach(function(e){t.types.$contains(e.code)||(e.isOn=!1,i.push(e))}),{config:t,moreOptionsVisible:!1,allTypes:i}},watch:{"config.level":function(e){let t=parseInt(e);isNaN(t)||t<1?t=1:10<t&&(t=10),this.config.level=t}},methods:{isOn:function(){return(!this.vIsLocation&&!this.vIsGroup||this.config.isPrior)&&this.config.isOn},changeExtensions:function(i){i.forEach(function(e,t){0<e.length&&"."!=e[0]&&(i[t]="."+e)}),this.config.extensions=i},changeMimeTypes:function(e){this.config.mimeTypes=e},changeAdvancedVisible:function(){this.moreOptionsVisible=!this.moreOptionsVisible},changeConds:function(e){this.config.conds=e},changeType:function(){this.config.types=[];let t=this;this.allTypes.forEach(function(e){e.isOn&&t.config.types.push(e.code)})},initSortableTypes:function(){let s=document.querySelector("#compression-types-box"),n=this;Sortable.create(s,{draggable:".checkbox",handle:".icon.handle",onStart:function(){},onUpdate:function(e){let t=s.querySelectorAll(".checkbox"),i=[];t.forEach(function(e){e=e.getAttribute("data-code");i.push(e)}),n.config.types=i}})}},template:`<div>
	<input type="hidden" name="compressionJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
			<tr>
				<td class="title">是否启用</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="config.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td>压缩级别</td>
				<td>
					<select class="ui dropdown auto-width" v-model="config.level">
						<option v-for="i in 10" :value="i">{{i}}</option>	
					</select>
					<p class="comment">级别越高，压缩比例越大。</p>
				</td>
			</tr>
			<tr>
				<td>支持的扩展名</td>
				<td>
					<values-box :values="config.extensions" @change="changeExtensions" placeholder="比如 .html"></values-box>
					<p class="comment">含有这些扩展名的URL将会被压缩，不区分大小写。</p>
				</td>
			</tr>
			<tr>
				<td>支持的MimeType</td>
				<td>
					<values-box :values="config.mimeTypes" @change="changeMimeTypes" placeholder="比如 text/*"></values-box>
					<p class="comment">响应的Content-Type里包含这些MimeType的内容将会被压缩。</p>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-if="isOn()"></more-options-tbody>
		<tbody v-show="isOn() && moreOptionsVisible">
			<tr>
				<td>压缩算法</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="config.useDefaultTypes" id="compression-use-default"/>
						<label v-if="config.useDefaultTypes" for="compression-use-default">使用默认顺序<span class="grey small">（brotli、gzip、deflate）</span></label>
						<label v-if="!config.useDefaultTypes" for="compression-use-default">使用默认顺序</label>
					</div>
					<div v-show="!config.useDefaultTypes">
						<div class="ui divider"></div>
						<div id="compression-types-box">
							<div class="ui checkbox" v-for="t in allTypes" style="margin-right: 2em" :data-code="t.code">
								<input type="checkbox" v-model="t.isOn" :id="'compression-type-' + t.code" @change="changeType"/>
								<label :for="'compression-type-' + t.code">{{t.name}} &nbsp; <i class="icon list small grey handle"></i></label>
							</div>
						</div>
					</div>
					
					<p class="comment" v-show="!config.useDefaultTypes">选择支持的压缩算法和优先顺序，拖动<i class="icon list small grey"></i>图表排序。</p>
				</td>
			</tr>
			<tr>
				<td>支持已压缩内容</td>
				<td>
					<checkbox v-model="config.decompressData"></checkbox>
					<p class="comment">支持对已压缩内容尝试重新使用新的算法压缩；不选中表示保留当前的压缩格式。</p>
				</td>
			</tr>
			<tr>
				<td>内容最小长度</td>
				<td>
					<size-capacity-box :v-name="'minLength'" :v-value="config.minLength" :v-unit="'kb'"></size-capacity-box>
					<p class="comment">0表示不限制，内容长度从文件尺寸或Content-Length中获取。</p>
				</td>
			</tr>
			<tr>
				<td>内容最大长度</td>
				<td>
					<size-capacity-box :v-name="'maxLength'" :v-value="config.maxLength" :v-unit="'mb'"></size-capacity-box>
					<p class="comment">0表示不限制，内容长度从文件尺寸或Content-Length中获取。</p>
				</td>
			</tr>
			<tr>
				<td>匹配条件</td>
				<td>
					<http-request-conds-box :v-conds="config.conds" @change="changeConds"></http-request-conds-box>
	</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`}),Vue.component("firewall-event-level-options",{props:["v-value"],mounted:function(){let t=this;Tea.action("/ui/eventLevelOptions").post().success(function(e){t.levels=e.data.eventLevels,t.change()})},data:function(){let e=this.vValue;return{levels:[],description:"",level:e=null!=e&&0!=e.length?e:""}},methods:{change:function(){this.$emit("change");let i=this;var e=this.levels.$find(function(e,t){return t.code==i.level});this.description=null!=e?e.description:""}},template:`<div>
    <select class="ui dropdown auto-width" name="eventLevel" v-model="level" @change="change">
        <option v-for="level in levels" :value="level.code">{{level.name}}</option>
    </select>
    <p class="comment">{{description}}</p>
</div>`}),Vue.component("prior-checkbox",{props:["v-config","description"],data:function(){let e=this.description;return null==e&&(e="打开后可以覆盖父级或子级配置"),{isPrior:this.vConfig.isPrior,realDescription:e}},watch:{isPrior:function(e){this.vConfig.isPrior=e}},template:`<tbody>
	<tr :class="{active:isPrior}">
		<td class="title">打开独立配置</td>
		<td>
			<div class="ui toggle checkbox">
				<input type="checkbox" v-model="isPrior"/>
				<label class="red"></label>
			</div>
			<p class="comment"><strong v-if="isPrior">[已打开]</strong> {{realDescription}}。</p>
		</td>
	</tr>
</tbody>`}),Vue.component("http-charsets-box",{props:["v-usual-charsets","v-all-charsets","v-charset-config","v-is-location","v-is-group"],data:function(){let e=this.vCharsetConfig;return{charsetConfig:e=null==e?{isPrior:!1,isOn:!1,charset:"",isUpper:!1}:e,advancedVisible:!1}},methods:{changeAdvancedVisible:function(e){this.advancedVisible=e}},template:`<div>
	<input type="hidden" name="charsetJSON" :value="JSON.stringify(charsetConfig)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="charsetConfig" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || charsetConfig.isPrior">
			<tr>
				<td class="title">是否启用</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="charsetConfig.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="((!vIsLocation && !vIsGroup) || charsetConfig.isPrior) && charsetConfig.isOn">	
			<tr>
				<td class="title">选择字符编码</td>
				<td><select class="ui dropdown" style="width:20em" name="charset" v-model="charsetConfig.charset">
						<option value="">[未选择]</option>
						<optgroup label="常用字符编码"></optgroup>
						<option v-for="charset in vUsualCharsets" :value="charset.charset">{{charset.charset}}（{{charset.name}}）</option>
						<optgroup label="全部字符编码"></optgroup>
						<option v-for="charset in vAllCharsets" :value="charset.charset">{{charset.charset}}（{{charset.name}}）</option>
					</select>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-if="((!vIsLocation && !vIsGroup) || charsetConfig.isPrior) && charsetConfig.isOn"></more-options-tbody>
		<tbody v-show="((!vIsLocation && !vIsGroup) || charsetConfig.isPrior) && charsetConfig.isOn && advancedVisible">
			<tr>
				<td>字符编码是否大写</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="charsetConfig.isUpper"/>
						<label></label>
					</div>
					<p class="comment">选中后将指定的字符编码转换为大写，比如默认为<span class="ui label tiny">utf-8</span>，选中后将改为<span class="ui label tiny">UTF-8</span>。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`}),Vue.component("http-expires-time-config-box",{props:["v-expires-time"],data:function(){let e=this.vExpiresTime;return{expiresTime:e=null==e?{isPrior:!1,isOn:!1,overwrite:!0,autoCalculate:!0,duration:{count:-1,unit:"hour"}}:e}},watch:{"expiresTime.isPrior":function(){this.notifyChange()},"expiresTime.isOn":function(){this.notifyChange()},"expiresTime.overwrite":function(){this.notifyChange()},"expiresTime.autoCalculate":function(){this.notifyChange()}},methods:{notifyChange:function(){this.$emit("change",this.expiresTime)}},template:`<div>
	<table class="ui table">
		<prior-checkbox :v-config="expiresTime"></prior-checkbox>
		<tbody v-show="expiresTime.isPrior">
			<tr>
				<td class="title">是否启用</td>
				<td><checkbox v-model="expiresTime.isOn"></checkbox>
					<p class="comment">启用后，将会在响应的Header中添加<code-label>Expires</code-label>字段，浏览器据此会将内容缓存在客户端；同时，在管理后台执行清理缓存时，也将无法清理客户端已有的缓存。</p>
				</td>
			</tr>
			<tr v-show="expiresTime.isPrior && expiresTime.isOn">
				<td>覆盖源站设置</td>
				<td>
					<checkbox v-model="expiresTime.overwrite"></checkbox>
					<p class="comment">选中后，会覆盖源站Header中已有的<code-label>Expires</code-label>字段。</p>
				</td>
			</tr>
			<tr v-show="expiresTime.isPrior && expiresTime.isOn">
				<td>自动计算时间</td>
				<td><checkbox v-model="expiresTime.autoCalculate"></checkbox>
					<p class="comment">根据已设置的缓存有效期进行计算。</p>
				</td>
			</tr>
			<tr v-show="expiresTime.isPrior && expiresTime.isOn && !expiresTime.autoCalculate">
				<td>强制缓存时间</td>
				<td>
					<time-duration-box :v-value="expiresTime.duration" @change="notifyChange"></time-duration-box>
					<p class="comment">从客户端访问的时间开始要缓存的时长。</p>
				</td>
			</tr>
		</tbody>
	</table>
</div>`}),Vue.component("http-access-log-box",{props:["v-access-log","v-keyword","v-show-server-link"],data:function(){let e=this.vAccessLog;return null!=e.header&&null!=e.header.Upgrade&&null!=e.header.Upgrade.values&&e.header.Upgrade.values.$contains("websocket")&&("http"==e.scheme?e.scheme="ws":"https"==e.scheme&&(e.scheme="wss")),{accessLog:e}},methods:{formatCost:function(e){if(null==e)return"0";let t=(1e3*e).toString(),i=t.split(".");return i.length<2?t:i[0]+"."+i[1].substring(0,3)},showLog:function(){let e=this;var t=this.accessLog.requestId;this.$parent.$children.forEach(function(e){null!=e.deselect&&e.deselect()}),this.select(),teaweb.popup("/servers/server/log/viewPopup?requestId="+t,{width:"50em",height:"28em",onClose:function(){e.deselect()}})},select:function(){this.$refs.box.parentNode.style.cssText="background: rgba(0, 0, 0, 0.1)"},deselect:function(){this.$refs.box.parentNode.style.cssText=""}},template:`<div style="word-break: break-all" :style="{'color': (accessLog.status >= 400) ? '#dc143c' : ''}" ref="box">
	<div>
		<a v-if="accessLog.node != null && accessLog.node.nodeCluster != null" :href="'/clusters/cluster/node?nodeId=' + accessLog.node.id + '&clusterId=' + accessLog.node.nodeCluster.id" title="点击查看节点详情" target="_top"><span class="grey">[{{accessLog.node.name}}<span v-if="!accessLog.node.name.endsWith('节点')">节点</span>]</span></a>
		<a :href="'/servers/server/log?serverId=' + accessLog.serverId" title="点击到网站服务" v-if="vShowServerLink"><span class="grey">[服务]</span></a>
		<span v-if="accessLog.region != null && accessLog.region.length > 0" class="grey"><ip-box :v-ip="accessLog.remoteAddr">[{{accessLog.region}}]</ip-box></span> <ip-box><keyword :v-word="vKeyword">{{accessLog.remoteAddr}}</keyword></ip-box> [{{accessLog.timeLocal}}] <em>&quot;<keyword :v-word="vKeyword">{{accessLog.requestMethod}}</keyword> {{accessLog.scheme}}://<keyword :v-word="vKeyword">{{accessLog.host}}</keyword><keyword :v-word="vKeyword">{{accessLog.requestURI}}</keyword> <a :href="accessLog.scheme + '://' + accessLog.host + accessLog.requestURI" target="_blank" title="新窗口打开" class="disabled"><i class="external icon tiny"></i> </a> {{accessLog.proto}}&quot; </em> <keyword :v-word="vKeyword">{{accessLog.status}}</keyword> <code-label v-if="accessLog.attrs != null && (accessLog.attrs['cache.status'] == 'HIT' || accessLog.attrs['cache.status'] == 'STALE')">cache {{accessLog.attrs['cache.status'].toLowerCase()}}</code-label> <code-label v-if="accessLog.firewallActions != null && accessLog.firewallActions.length > 0">waf {{accessLog.firewallActions}}</code-label> <span v-if="accessLog.tags != null && accessLog.tags.length > 0">- <code-label v-for="tag in accessLog.tags" :key="tag">{{tag}}</code-label></span>
		
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
</div>`}),Vue.component("http-firewall-block-options-viewer",{props:["v-block-options"],data:function(){return{options:this.vBlockOptions}},template:`<div>
	<span v-if="options == null">默认设置</span>
	<div v-else>
		状态码：{{options.statusCode}} / 提示内容：<span v-if="options.body != null && options.body.length > 0">[{{options.body.length}}字符]</span><span v-else class="disabled">[无]</span>  / 超时时间：{{options.timeout}}秒
	</div>
</div>	
`}),Vue.component("http-access-log-config-box",{props:["v-access-log-config","v-fields","v-default-field-codes","v-is-location","v-is-group"],data:function(){let t=this,i=(setTimeout(function(){t.changeFields()},100),{isPrior:!1,isOn:!1,fields:[1,2,6,7],status1:!0,status2:!0,status3:!0,status4:!0,status5:!0,firewallOnly:!1,enableClientClosed:!1});return null!=this.vAccessLogConfig&&(i=this.vAccessLogConfig),this.vFields.forEach(function(e){null==t.vAccessLogConfig?e.isChecked=t.vDefaultFieldCodes.$contains(e.code):e.isChecked=i.fields.$contains(e.code)}),{accessLog:i,hasRequestBodyField:this.vFields.$contains(8)}},methods:{changeFields:function(){this.accessLog.fields=this.vFields.filter(function(e){return e.isChecked}).map(function(e){return e.code}),this.hasRequestBodyField=this.accessLog.fields.$contains(8)}},template:`<div>
	<input type="hidden" name="accessLogJSON" :value="JSON.stringify(accessLog)"/>
	<table class="ui table definition selectable" :class="{'opacity-mask': this.accessLog.firewallOnly}">
		<prior-checkbox :v-config="accessLog" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || accessLog.isPrior">
			<tr>
				<td class="title">开启访问日志</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="accessLog.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody  v-show="((!vIsLocation && !vIsGroup) || accessLog.isPrior) && accessLog.isOn">
			<tr>
				<td>基础信息</td>
				<td><p class="comment" style="padding-top: 0">默认记录客户端IP、请求URL等基础信息。</p></td>
			</tr>
			<tr>
				<td>高级信息</td>
				<td>
					<div class="ui checkbox" v-for="(field, index) in vFields" style="width:10em;margin-bottom:0.8em">
						<input type="checkbox" v-model="field.isChecked" @change="changeFields" :id="'access-log-field-' + index"/>
						<label :for="'access-log-field-' + index">{{field.name}}</label>
					</div>
					<p class="comment">在基础信息之外要存储的信息。
						<span class="red" v-if="hasRequestBodyField">记录"请求Body"将会显著消耗更多的系统资源，建议仅在调试时启用，最大记录尺寸为2MB。</span>
					</p>
				</td>
			</tr>
			<tr>
				<td>要存储的访问日志状态码</td>
				<td>
					<div class="ui checkbox" style="width:3.5em">
						<input type="checkbox" v-model="accessLog.status1"/>
						<label>1xx</label>
					</div>
					<div class="ui checkbox" style="width:3.5em">
						<input type="checkbox" v-model="accessLog.status2"/>
						<label>2xx</label>
					</div>
					<div class="ui checkbox" style="width:3.5em">
						<input type="checkbox" v-model="accessLog.status3"/>
						<label>3xx</label>
					</div>
					<div class="ui checkbox" style="width:3.5em">
						<input type="checkbox" v-model="accessLog.status4"/>
						<label>4xx</label>
					</div>
					<div class="ui checkbox" style="width:3.5em">
						<input type="checkbox" v-model="accessLog.status5"/>
						<label>5xx</label>
					</div>
				</td>
			</tr>
			<tr>
				<td>记录客户端中断日志</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="accessLog.enableClientClosed"/>
						<label></label>
					</div>
					<p class="comment">以<code-label>499</code-label>的状态码记录客户端主动中断日志。</p>
				</td>
			</tr>
		</tbody>
	</table>
	
	<div v-show="((!vIsLocation && !vIsGroup) || accessLog.isPrior) && accessLog.isOn">
        <h4>WAF相关</h4>
        <table class="ui table definition selectable">
            <tr>
                <td class="title">只记录WAF相关日志</td>
                <td>
                    <checkbox v-model="accessLog.firewallOnly"></checkbox>
                    <p class="comment">选中后只记录WAF相关的日志。通过此选项可有效减少访问日志数量，降低网络带宽和存储压力。</p>
                </td>
            </tr>
        </table>
    </div>
	<div class="margin"></div>
</div>`}),Vue.component("traffic-limit-view",{props:["v-traffic-limit"],data:function(){return{config:this.vTrafficLimit}},template:`<div>
	<div v-if="config.isOn">
		<span v-if="config.dailySize != null && config.dailySize.count > 0">日流量限制：{{config.dailySize.count}}{{config.dailySize.unit.toUpperCase()}}<br/></span>
		<span v-if="config.monthlySize != null && config.monthlySize.count > 0">月流量限制：{{config.monthlySize.count}}{{config.monthlySize.unit.toUpperCase()}}<br/></span>
	</div>
	<span v-else class="disabled">没有限制。</span>
</div>`}),Vue.component("http-auth-basic-auth-user-box",{props:["v-users"],data:function(){let e=this.vUsers;return{users:e=null==e?[]:e,isAdding:!1,updatingIndex:-1,username:"",password:""}},methods:{add:function(){this.isAdding=!0,this.username="",this.password="";let e=this;setTimeout(function(){e.$refs.username.focus()},100)},cancel:function(){this.isAdding=!1,this.updatingIndex=-1},confirm:function(){let e=this;0==this.username.length?teaweb.warn("请输入用户名",function(){e.$refs.username.focus()}):0==this.password.length?teaweb.warn("请输入密码",function(){e.$refs.password.focus()}):(this.updatingIndex<0?this.users.push({username:this.username,password:this.password}):(this.users[this.updatingIndex].username=this.username,this.users[this.updatingIndex].password=this.password),this.cancel())},update:function(e,t){this.updatingIndex=e,this.isAdding=!0,this.username=t.username,this.password=t.password;let i=this;setTimeout(function(){i.$refs.username.focus()},100)},remove:function(e){this.users.$remove(e)}},template:`<div>
	<input type="hidden" name="httpAuthBasicAuthUsersJSON" :value="JSON.stringify(users)"/>
	<div v-if="users.length > 0">
		<div class="ui label small basic" v-for="(user, index) in users">
			{{user.username}} <a href="" title="修改" @click.prevent="update(index, user)"><i class="icon pencil tiny"></i></a>
			<a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div v-show="isAdding">
		<div class="ui fields inline">
			<div class="ui field">
				<input type="text" placeholder="用户名" v-model="username" size="15" ref="username"/>
			</div>
			<div class="ui field">
				<input type="password" placeholder="密码" v-model="password" size="15" ref="password"/>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>&nbsp;
				<a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
	</div>
	<div v-if="!isAdding" style="margin-top: 1em">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`}),Vue.component("http-location-labels",{props:["v-location-config","v-server-id"],data:function(){return{location:this.vLocationConfig}},methods:{configIsOn:function(e){return null!=e&&e.isPrior&&e.isOn},refIsOn:function(e,t){return this.configIsOn(e)&&null!=t&&t.isOn},len:function(e){return null==e?0:e.length},url:function(e){return"/servers/server/settings/locations"+e+"?serverId="+this.vServerId+"&locationId="+this.location.id}},template:`	<div class="labels-box">
	<!-- 基本信息 -->
	<http-location-labels-label v-if="location.name != null && location.name.length > 0" :class="'olive'" :href="url('/location')">{{location.name}}</http-location-labels-label>
	
	<!-- domains -->
	<div v-if="location.domains != null && location.domains.length > 0">
		<grey-label v-for="domain in location.domains">{{domain}}</grey-label>
	</div>
	
	<!-- break -->
	<http-location-labels-label v-if="location.isBreak" :href="url('/location')">BREAK</http-location-labels-label>
	
	<!-- redirectToHTTPS -->
	<http-location-labels-label v-if="location.web != null && configIsOn(location.web.redirectToHTTPS)" :href="url('/http')">自动跳转HTTPS</http-location-labels-label>
	
	<!-- Web -->
	<http-location-labels-label v-if="location.web != null && configIsOn(location.web.root)" :href="url('/web')">文档根目录</http-location-labels-label>
	
	<!-- 反向代理 -->
	<http-location-labels-label v-if="refIsOn(location.reverseProxyRef, location.reverseProxy)" :v-href="url('/reverseProxy')">反向代理</http-location-labels-label>
	
	<!-- WAF -->
	<!-- TODO -->
	
	<!-- Cache -->
	<http-location-labels-label v-if="location.web != null && configIsOn(location.web.cache)" :v-href="url('/cache')">CACHE</http-location-labels-label>
	
	<!-- Charset -->
	<http-location-labels-label v-if="location.web != null && configIsOn(location.web.charset) && location.web.charset.charset.length > 0" :href="url('/charset')">{{location.web.charset.charset}}</http-location-labels-label>
	
	<!-- 访问日志 -->
	<!-- TODO -->
	
	<!-- 统计 -->
	<!-- TODO -->
	
	<!-- Gzip -->
	<http-location-labels-label v-if="location.web != null && refIsOn(location.web.gzipRef, location.web.gzip) && location.web.gzip.level > 0" :href="url('/gzip')">Gzip:{{location.web.gzip.level}}</http-location-labels-label>
	
	<!-- HTTP Header -->
	<http-location-labels-label v-if="location.web != null && refIsOn(location.web.requestHeaderPolicyRef, location.web.requestHeaderPolicy) && (len(location.web.requestHeaderPolicy.addHeaders) > 0 || len(location.web.requestHeaderPolicy.setHeaders) > 0 || len(location.web.requestHeaderPolicy.replaceHeaders) > 0 || len(location.web.requestHeaderPolicy.deleteHeaders) > 0)" :href="url('/headers')">请求Header</http-location-labels-label>
	<http-location-labels-label v-if="location.web != null && refIsOn(location.web.responseHeaderPolicyRef, location.web.responseHeaderPolicy) && (len(location.web.responseHeaderPolicy.addHeaders) > 0 || len(location.web.responseHeaderPolicy.setHeaders) > 0 || len(location.web.responseHeaderPolicy.replaceHeaders) > 0 || len(location.web.responseHeaderPolicy.deleteHeaders) > 0)" :href="url('/headers')">响应Header</http-location-labels-label>
	
	<!-- Websocket -->
	<http-location-labels-label v-if="location.web != null && refIsOn(location.web.websocketRef, location.web.websocket)" :href="url('/websocket')">Websocket</http-location-labels-label>
	
	<!-- 请求脚本 -->
	<http-location-labels-label v-if="location.web != null && location.web.requestScripts != null && ((location.web.requestScripts.initGroup != null && location.web.requestScripts.initGroup.isPrior) || (location.web.requestScripts.requestGroup != null && location.web.requestScripts.requestGroup.isPrior))" :href="url('/requestScripts')">请求脚本</http-location-labels-label>
	
	<!-- 自定义页面 -->
	<div v-if="location.web != null && location.web.pages != null && location.web.pages.length > 0">
		<div v-for="page in location.web.pages" :key="page.id"><http-location-labels-label :href="url('/pages')">PAGE [状态码{{page.status[0]}}] -&gt; {{page.url}}</http-location-labels-label></div>
	</div>
	<div v-if="location.web != null && configIsOn(location.web.shutdown)">
		<http-location-labels-label :v-class="'red'" :href="url('/pages')">临时关闭</http-location-labels-label>
	</div>
	
	<!-- 重写规则 -->
	<div v-if="location.web != null && location.web.rewriteRules != null && location.web.rewriteRules.length > 0">
		<div v-for="rewriteRule in location.web.rewriteRules">
			<http-location-labels-label :href="url('/rewrite')">REWRITE {{rewriteRule.pattern}} -&gt; {{rewriteRule.replace}}</http-location-labels-label>
		</div>
	</div>
</div>`}),Vue.component("http-location-labels-label",{props:["v-class","v-href"],template:'<a :href="vHref" class="ui label tiny basic" :class="vClass" style="font-size:0.7em;padding:4px;margin-top:0.3em;margin-bottom:0.3em"><slot></slot></a>'}),Vue.component("http-gzip-box",{props:["v-gzip-config","v-gzip-ref","v-is-location"],data:function(){let e=this.vGzipConfig;return{gzip:e=null==e?{isOn:!0,level:0,minLength:null,maxLength:null,conds:null}:e,advancedVisible:!1}},methods:{isOn:function(){return(!this.vIsLocation||this.vGzipRef.isPrior)&&this.vGzipRef.isOn},changeAdvancedVisible:function(e){this.advancedVisible=e}},template:`<div>
<input type="hidden" name="gzipRefJSON" :value="JSON.stringify(vGzipRef)"/> 
<table class="ui table selectable definition">
	<prior-checkbox :v-config="vGzipRef" v-if="vIsLocation"></prior-checkbox>
	<tbody v-show="!vIsLocation || vGzipRef.isPrior">
		<tr>
			<td class="title">启用Gzip压缩</td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" v-model="vGzipRef.isOn"/>
					<label></label>
				</div>
			</td>
		</tr>
	</tbody>
	<tbody v-show="isOn()">
		<tr>
			<td class="title">压缩级别</td>
			<td>
				<select class="dropdown auto-width" name="level" v-model="gzip.level">
					<option value="0">不压缩</option>
					<option v-for="i in 9" :value="i">{{i}}</option>
				</select>
				<p class="comment">级别越高，压缩比例越大。</p>
			</td>
		</tr>
	</tbody>
	<more-options-tbody @change="changeAdvancedVisible" v-if="isOn()"></more-options-tbody>
	<tbody v-show="isOn() && advancedVisible">
		<tr>
			<td>Gzip内容最小长度</td>
			<td>
				<size-capacity-box :v-name="'minLength'" :v-value="gzip.minLength" :v-unit="'kb'"></size-capacity-box>
				<p class="comment">0表示不限制，内容长度从文件尺寸或Content-Length中获取。</p>
			</td>
		</tr>
		<tr>
			<td>Gzip内容最大长度</td>
			<td>
				<size-capacity-box :v-name="'maxLength'" :v-value="gzip.maxLength" :v-unit="'mb'"></size-capacity-box>
				<p class="comment">0表示不限制，内容长度从文件尺寸或Content-Length中获取。</p>
			</td>
		</tr>
		<tr>
			<td>匹配条件</td>
			<td>
				<http-request-conds-box :v-conds="gzip.conds"></http-request-conds-box>
</td>
		</tr>
	</tbody>
</table>
</div>`}),Vue.component("script-config-box",{props:["id","v-script-config","comment"],data:function(){let e=this.vScriptConfig;return 0==(e=null==e?{isPrior:!1,isOn:!1,code:""}:e).code.length&&(e.code="\n\n\n\n"),{config:e}},watch:{"config.isOn":function(){this.change()}},methods:{change:function(){this.$emit("change",this.config)},changeCode:function(e){this.config.code=e,this.change()}},template:`<div>
	<table class="ui table definition selectable">
		<tbody>
			<tr>
				<td class="title">是否启用</td>
				<td><checkbox v-model="config.isOn"></checkbox></td>
			</tr>
		</tbody>
		<tbody>
			<tr :style="{opacity: !config.isOn ? 0.5 : 1}">
				<td>脚本代码</td>	
				<td><source-code-box :id="id" type="text/javascript" :read-only="false" @change="changeCode">{{config.code}}</source-code-box>
					<p class="comment">{{comment}}</p>
				</td>
			</tr>
		</tbody>
	</table>
</div>`}),Vue.component("ssl-certs-view",{props:["v-certs"],data:function(){let e=this.vCerts;return{certs:e=null==e?[]:e}},methods:{formatTime:function(e){return new Date(1e3*e).format("Y-m-d")},viewCert:function(e){teaweb.popup("/servers/certs/certPopup?certId="+e,{height:"28em",width:"48em"})}},template:`<div>
	<div v-if="certs != null && certs.length > 0">
		<div class="ui label small basic" v-for="(cert, index) in certs">
			{{cert.name}} / {{cert.dnsNames}} / 有效至{{formatTime(cert.timeEndAt)}} &nbsp;<a href="" title="查看" @click.prevent="viewCert(cert.id)"><i class="icon expand blue"></i></a>
		</div>
	</div>
</div>`}),Vue.component("http-firewall-captcha-options-viewer",{props:["v-captcha-options"],mounted:function(){this.updateSummary()},data:function(){let e=this.vCaptchaOptions;return{options:e=null==e?{life:0,maxFails:0,failBlockTimeout:0,failBlockScopeAll:!1,uiIsOn:!1,uiTitle:"",uiPrompt:"",uiButtonTitle:"",uiShowRequestId:!1,uiCss:"",uiFooter:"",uiBody:"",cookieId:"",lang:""}:e,summary:""}},methods:{updateSummary:function(){let e=[];0<this.options.life&&e.push("有效时间"+this.options.life+"秒"),0<this.options.maxFails&&e.push("最多失败"+this.options.maxFails+"次"),0<this.options.failBlockTimeout&&e.push("失败拦截"+this.options.failBlockTimeout+"秒"),this.options.failBlockScopeAll&&e.push("全局封禁"),this.options.uiIsOn&&e.push("定制UI"),0==e.length?this.summary="默认配置":this.summary=e.join(" / ")}},template:`<div>{{summary}}</div>
`}),Vue.component("reverse-proxy-box",{props:["v-reverse-proxy-ref","v-reverse-proxy-config","v-is-location","v-is-group","v-family"],data:function(){let e=this.vReverseProxyRef,t=(null==e&&(e={isPrior:!1,isOn:!1,reverseProxyId:0}),this.vReverseProxyConfig),i=(null==(t=null==t?{requestPath:"",stripPrefix:"",requestURI:"",requestHost:"",requestHostType:0,addHeaders:[],connTimeout:{count:0,unit:"second"},readTimeout:{count:0,unit:"second"},idleTimeout:{count:0,unit:"second"},maxConns:0,maxIdleConns:0,followRedirects:!1}:t).addHeaders&&(t.addHeaders=[]),null==t.connTimeout&&(t.connTimeout={count:0,unit:"second"}),null==t.readTimeout&&(t.readTimeout={count:0,unit:"second"}),null==t.idleTimeout&&(t.idleTimeout={count:0,unit:"second"}),null==t.proxyProtocol&&Vue.set(t,"proxyProtocol",{isOn:!1,version:1}),[{name:"X-Real-IP",isChecked:!1},{name:"X-Forwarded-For",isChecked:!1},{name:"X-Forwarded-By",isChecked:!1},{name:"X-Forwarded-Host",isChecked:!1},{name:"X-Forwarded-Proto",isChecked:!1}]);return i.forEach(function(e){e.isChecked=t.addHeaders.$contains(e.name)}),{reverseProxyRef:e,reverseProxyConfig:t,advancedVisible:!1,family:this.vFamily,forwardHeaders:i}},watch:{"reverseProxyConfig.requestHostType":function(e){let t=parseInt(e);isNaN(t)&&(t=0),this.reverseProxyConfig.requestHostType=t},"reverseProxyConfig.connTimeout.count":function(e){let t=parseInt(e);(isNaN(t)||t<0)&&(t=0),this.reverseProxyConfig.connTimeout.count=t},"reverseProxyConfig.readTimeout.count":function(e){let t=parseInt(e);(isNaN(t)||t<0)&&(t=0),this.reverseProxyConfig.readTimeout.count=t},"reverseProxyConfig.idleTimeout.count":function(e){let t=parseInt(e);(isNaN(t)||t<0)&&(t=0),this.reverseProxyConfig.idleTimeout.count=t},"reverseProxyConfig.maxConns":function(e){let t=parseInt(e);(isNaN(t)||t<0)&&(t=0),this.reverseProxyConfig.maxConns=t},"reverseProxyConfig.maxIdleConns":function(e){let t=parseInt(e);(isNaN(t)||t<0)&&(t=0),this.reverseProxyConfig.maxIdleConns=t},"reverseProxyConfig.proxyProtocol.version":function(e){let t=parseInt(e);isNaN(t)&&(t=1),this.reverseProxyConfig.proxyProtocol.version=t}},methods:{isOn:function(){return(!this.vIsLocation&&!this.vIsGroup||this.reverseProxyRef.isPrior)&&this.reverseProxyRef.isOn},changeAdvancedVisible:function(e){this.advancedVisible=e},changeAddHeader:function(){this.reverseProxyConfig.addHeaders=this.forwardHeaders.filter(function(e){return e.isChecked}).map(function(e){return e.name})}},template:`<div>
	<input type="hidden" name="reverseProxyRefJSON" :value="JSON.stringify(reverseProxyRef)"/>
	<input type="hidden" name="reverseProxyJSON" :value="JSON.stringify(reverseProxyConfig)"/>
	<table class="ui table selectable definition">
		<prior-checkbox :v-config="reverseProxyRef" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || reverseProxyRef.isPrior">
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
				<td>回源跟随</td>
				<td>
					<checkbox v-model="reverseProxyConfig.followRedirects"></checkbox>
					<p class="comment">选中后，自动读取源站跳转后的网页内容。</p>
				</td>
			</tr>
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
			<tr v-if="family == null || family == 'http'">
				<td>是否自动刷新缓存区<em>（AutoFlush）</em></td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="reverseProxyConfig.autoFlush"/>
						<label></label>
					</div>
					<p class="comment">开启后将自动刷新缓冲区数据到客户端，在类似于SSE（server-sent events）等场景下很有用。</p>
				</td>
			</tr>
			<tr v-if="family == null || family == 'http'">
                <td class="color-border">源站连接失败超时时间</td>
                <td>
                    <div class="ui fields inline">
                        <div class="ui field">
                            <input type="text" name="connTimeout" value="10" size="6" v-model="reverseProxyConfig.connTimeout.count"/>
                        </div>
                        <div class="ui field">
                            秒
                        </div>
                    </div>
                    <p class="comment">连接源站失败的最大超时时间，0表示不限制。</p>
                </td>
            </tr>
            <tr v-if="family == null || family == 'http'">
                <td class="color-border">源站读取超时时间</td>
                <td>
                    <div class="ui fields inline">
                        <div class="ui field">
                            <input type="text" name="readTimeout" value="0" size="6" v-model="reverseProxyConfig.readTimeout.count"/>
                        </div>
                        <div class="ui field">
                            秒
                        </div>
                    </div>
                    <p class="comment">读取内容时的最大超时时间，0表示不限制。</p>
                </td>
            </tr>
            <tr v-if="family == null || family == 'http'">
                <td class="color-border">源站最大并发连接数</td>
                <td>
                    <div class="ui fields inline">
                        <div class="ui field">
                            <input type="text" name="maxConns" value="0" size="6" maxlength="10" v-model="reverseProxyConfig.maxConns"/>
                        </div>
                    </div>
                    <p class="comment">源站可以接受到的最大并发连接数，0表示使用系统默认。</p>
                </td>
            </tr>
            <tr v-if="family == null || family == 'http'">
                <td class="color-border">源站最大空闲连接数</td>
                <td>
                    <div class="ui fields inline">
                        <div class="ui field">
                            <input type="text" name="maxIdleConns" value="0" size="6" maxlength="10" v-model="reverseProxyConfig.maxIdleConns"/>
                        </div>
                    </div>
                    <p class="comment">当没有请求时，源站保持等待的最大空闲连接数量，0表示使用系统默认。</p>
                </td>
            </tr>
            <tr v-if="family == null || family == 'http'">
                <td class="color-border">源站最大空闲超时时间</td>
                <td>
                    <div class="ui fields inline">
                        <div class="ui field">
                            <input type="text" name="idleTimeout" value="0" size="6" v-model="reverseProxyConfig.idleTimeout.count"/>
                        </div>
                        <div class="ui field">
                            秒
                        </div>
                    </div>
                    <p class="comment">源站保持等待的空闲超时时间，0表示使用默认时间。</p>
                </td>
            </tr>
            <tr v-show="family != 'unix'">
            	<td>PROXY Protocol</td>
            	<td>
            		<checkbox name="proxyProtocolIsOn" v-model="reverseProxyConfig.proxyProtocol.isOn"></checkbox>
            		<p class="comment">选中后表示启用PROXY Protocol，每次连接源站时都会在头部写入客户端地址信息。</p>
				</td>
			</tr>
			<tr v-show="family != 'unix' && reverseProxyConfig.proxyProtocol.isOn">
				<td>PROXY Protocol版本</td>
				<td>
					<select class="ui dropdown auto-width" name="proxyProtocolVersion" v-model="reverseProxyConfig.proxyProtocol.version">
						<option value="1">1</option>
						<option value="2">2</option>
					</select>
					<p class="comment" v-if="reverseProxyConfig.proxyProtocol.version == 1">发送类似于<code-label>PROXY TCP4 192.168.1.1 192.168.1.10 32567 443</code-label>的头部信息。</p>
					<p class="comment" v-if="reverseProxyConfig.proxyProtocol.version == 2">发送二进制格式的头部信息。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`}),Vue.component("http-firewall-param-filters-box",{props:["v-filters"],data:function(){let e=this.vFilters;return{filters:e=null==e?[]:e,isAdding:!1,options:[{name:"MD5",code:"md5"},{name:"URLEncode",code:"urlEncode"},{name:"URLDecode",code:"urlDecode"},{name:"BASE64Encode",code:"base64Encode"},{name:"BASE64Decode",code:"base64Decode"},{name:"UNICODE编码",code:"unicodeEncode"},{name:"UNICODE解码",code:"unicodeDecode"},{name:"HTML实体编码",code:"htmlEscape"},{name:"HTML实体解码",code:"htmlUnescape"},{name:"计算长度",code:"length"},{name:"十六进制->十进制",code:"hex2dec"},{name:"十进制->十六进制",code:"dec2hex"},{name:"SHA1",code:"sha1"},{name:"SHA256",code:"sha256"}],addingCode:""}},methods:{add:function(){this.isAdding=!0,this.addingCode=""},confirm:function(){if(0!=this.addingCode.length){let i=this;this.filters.push(this.options.$find(function(e,t){return t.code==i.addingCode})),this.isAdding=!1}},cancel:function(){this.isAdding=!1},remove:function(e){this.filters.$remove(e)}},template:`<div>
		<input type="hidden" name="paramFiltersJSON" :value="JSON.stringify(filters)" />
		<div v-if="filters.length > 0">
			<div v-for="(filter, index) in filters" class="ui label small basic">
				{{filter.name}} <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove"></i></a>
			</div>
			<div class="ui divider"></div>
		</div>
		<div v-if="isAdding">
			<div class="ui fields inline">
				<div class="ui field">
					<select class="ui dropdown auto-width" v-model="addingCode">
						<option value="">[请选择]</option>
						<option v-for="option in options" :value="option.code">{{option.name}}</option>
					</select>
				</div>
				<div class="ui field">
					<button class="ui button tiny" type="button" @click.prevent="confirm()">确定</button>
					&nbsp; <a href="" @click.prevent="cancel()" title="取消"><i class="icon remove"></i></a>
				</div>
			</div>
		</div>
		<div v-if="!isAdding">
			<button class="ui button tiny" type="button" @click.prevent="add">+</button>
		</div>
		<p class="comment">可以对参数值进行特定的编解码处理。</p>
</div>`}),Vue.component("http-remote-addr-config-box",{props:["v-remote-addr-config","v-is-location","v-is-group"],data:function(){let e=this.vRemoteAddrConfig,t="";return(e=null==e?{isPrior:!1,isOn:!1,value:"${rawRemoteAddr}",isCustomized:!1}:e).isCustomized||"${remoteAddr}"!=e.value&&"${rawRemoteAddr}"!=e.value||(t=e.value),{config:e,options:[{name:"直接获取",description:'用户直接访问边缘节点，即 "用户 --\x3e 边缘节点" 模式，这时候可以直接从连接中读取到真实的IP地址。',value:"${rawRemoteAddr}"},{name:"从上级代理中获取",description:'用户和边缘节点之间有别的代理服务转发，即 "用户 --\x3e [第三方代理服务] --\x3e 边缘节点"，这时候只能从上级代理中获取传递的IP地址。',value:"${remoteAddr}"},{name:"[自定义]",description:"通过自定义变量来获取客户端真实的IP地址。",value:""}],optionValue:t}},methods:{isOn:function(){return(!this.vIsLocation&&!this.vIsGroup||this.config.isPrior)&&this.config.isOn},changeOptionValue:function(){0<this.optionValue.length?(this.config.value=this.optionValue,this.config.isCustomized=!1):this.config.isCustomized=!0}},template:`<div>
	<input type="hidden" name="remoteAddrJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
			<tr>
				<td class="title">是否启用</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="config.isOn"/>
						<label></label>
					</div>
					<p class="comment">选中后表示使用自定义的请求变量获取客户端IP。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td>获取IP方式 *</td>
				<td>
					<select class="ui dropdown auto-width" v-model="optionValue" @change="changeOptionValue">
						<option v-for="option in options" :value="option.value">{{option.name}}</option>
					</select>
					<p class="comment" v-for="option in options" v-if="option.value == optionValue && option.description.length > 0">{{option.description}}</p>
				</td>
			</tr>
			<tr v-show="optionValue.length == 0">
				<td>读取IP变量值 *</td>
				<td>
					<input type="hidden" v-model="config.value" maxlength="100"/>
					<div v-if="optionValue == ''" style="margin-top: 1em">
						<input type="text" v-model="config.value" maxlength="100"/>
						<p class="comment">通过此变量获取用户的IP地址。具体可用的请求变量列表可参考官方网站文档。</p>
					</div>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>		
</div>`}),Vue.component("http-access-log-search-box",{props:["v-ip","v-domain","v-keyword","v-cluster-id","v-node-id"],data:function(){let e=this.vIp,t=(null==e&&(e=""),this.vDomain),i=(null==t&&(t=""),this.vKeyword);return null==i&&(i=""),{ip:e,domain:t,keyword:i,clusterId:this.vClusterId}},methods:{cleanIP:function(){this.ip="",this.submit()},cleanDomain:function(){this.domain="",this.submit()},cleanKeyword:function(){this.keyword="",this.submit()},submit:function(){let e=this.$el.parentNode;for(;;){if(null==e)break;if("FORM"==e.tagName)break;e=e.parentNode}null!=e&&setTimeout(function(){e.submit()},500)},changeCluster:function(e){this.clusterId=e}},template:`<div style="z-index: 10">
	<div class="margin"></div>
	<div class="ui fields inline">
		<div class="ui field">
			<div class="ui input left right labeled small">
				<span class="ui label basic" style="font-weight: normal">IP</span>
				<input type="text" name="ip" placeholder="x.x.x.x" size="15" v-model="ip"/>
				<a class="ui label basic" :class="{disabled: ip.length == 0}" @click.prevent="cleanIP"><i class="icon remove small"></i></a>
			</div>
		</div>
		<div class="ui field">
			<div class="ui input left right labeled small" >
				<span class="ui label basic" style="font-weight: normal">域名</span>
				<input type="text" name="domain" placeholder="xxx.com" size="15" v-model="domain"/>
				<a class="ui label basic" :class="{disabled: domain.length == 0}" @click.prevent="cleanDomain"><i class="icon remove small"></i></a>
			</div>
		</div>
		<div class="ui field">
			<div class="ui input left right labeled small">
				<span class="ui label basic" style="font-weight: normal">关键词</span>
				<input type="text" name="keyword" v-model="keyword" placeholder="路径、UserAgent等..." size="30"/>
				<a class="ui label basic" :class="{disabled: keyword.length == 0}" @click.prevent="cleanKeyword"><i class="icon remove small"></i></a>
			</div>
		</div>
		<div class="ui field"><tip-icon content="一些特殊的关键词：<br/>单个状态码：status:200<br/>状态码范围：status:500-504<br/>查询IP：ip:192.168.1.100<br/>查询URL：https://goedge.cn/docs"></tip-icon></div>
	</div>
	<div class="ui fields inline" style="margin-top: 0.5em">
		<div class="ui field">
			<node-cluster-combo-box :v-cluster-id="clusterId" @change="changeCluster"></node-cluster-combo-box>
		</div>
		<div class="ui field" v-if="clusterId > 0">
			<node-combo-box :v-cluster-id="clusterId" :v-node-id="vNodeId"></node-combo-box>
		</div>
		<slot></slot>
		<div class="ui field">
			<button class="ui button small" type="submit">搜索日志</button>
		</div>
	</div>
</div>`}),Vue.component("metric-key-label",{props:["v-key"],data:function(){return{keyDefs:window.METRIC_HTTP_KEYS}},methods:{keyName:function(i){let s=this,n="";var e=this.keyDefs.$find(function(e,t){return t.code==i||(i.startsWith("${arg.")&&t.code.startsWith("${arg.")?(n=s.getSubKey("arg.",i),!0):i.startsWith("${header.")&&t.code.startsWith("${header.")?(n=s.getSubKey("header.",i),!0):!(!i.startsWith("${cookie.")||!t.code.startsWith("${cookie."))&&(n=s.getSubKey("cookie.",i),!0))});return null!=e?0<n.length?e.name+": "+n:e.name:i},getSubKey:function(e,t){var i=t.indexOf(e="${"+e);return 0<=i?(t=t.substring(i+e.length)).substring(0,t.length-1):""}},template:`<div class="ui label basic small">
	{{keyName(this.vKey)}}
</div>`}),Vue.component("metric-keys-config-box",{props:["v-keys"],data:function(){let e=this.vKeys;return{keys:e=null==e?[]:e,isAdding:!1,key:"",subKey:"",keyDescription:"",keyDefs:window.METRIC_HTTP_KEYS}},watch:{keys:function(){this.$emit("change",this.keys)}},methods:{cancel:function(){this.key="",this.subKey="",this.keyDescription="",this.isAdding=!1},confirm:function(){if(0!=this.key.length){if(0<this.key.indexOf(".NAME")){if(0==this.subKey.length)return void teaweb.warn("请输入参数值");this.key=this.key.replace(".NAME","."+this.subKey)}this.keys.push(this.key),this.cancel()}},add:function(){this.isAdding=!0;let e=this;setTimeout(function(){null!=e.$refs.key&&e.$refs.key.focus()},100)},remove:function(e){this.keys.$remove(e)},changeKey:function(){if(0!=this.key.length){let i=this;var e=this.keyDefs.$find(function(e,t){return t.code==i.key});null!=e&&(this.keyDescription=e.description)}},keyName:function(i){let s=this,n="";var e=this.keyDefs.$find(function(e,t){return t.code==i||(i.startsWith("${arg.")&&t.code.startsWith("${arg.")?(n=s.getSubKey("arg.",i),!0):i.startsWith("${header.")&&t.code.startsWith("${header.")?(n=s.getSubKey("header.",i),!0):!(!i.startsWith("${cookie.")||!t.code.startsWith("${cookie."))&&(n=s.getSubKey("cookie.",i),!0))});return null!=e?0<n.length?e.name+": "+n:e.name:i},getSubKey:function(e,t){var i=t.indexOf(e="${"+e);return 0<=i?(t=t.substring(i+e.length)).substring(0,t.length-1):""}},template:`<div>
	<input type="hidden" name="keysJSON" :value="JSON.stringify(keys)"/>
	<div>
		<div v-for="(key, index) in keys" class="ui label small basic">
			{{keyName(key)}} &nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</div>
	</div>
	<div v-if="isAdding" style="margin-top: 1em">
		<div class="ui fields inline">
			<div class="ui field">
				<select class="ui dropdown" v-model="key" @change="changeKey">
					<option value="">[选择对象]</option>
					<option v-for="def in keyDefs" :value="def.code">{{def.name}}</option>
				</select>
			</div>
			<div class="ui field" v-if="key == '\${arg.NAME}'">
				<input type="text" v-model="subKey" placeholder="参数名" size="15"/>
			</div>
			<div class="ui field" v-if="key == '\${header.NAME}'">
				<input type="text" v-model="subKey" placeholder="Header名" size="15">
			</div>
			<div class="ui field" v-if="key == '\${cookie.NAME}'">
				<input type="text" v-model="subKey" placeholder="Cookie名" size="15">
			</div>
			<div class="ui field">
				<button type="button" class="ui button tiny" @click.prevent="confirm">确定</button>
				<a href="" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
		<p class="comment" v-if="keyDescription.length > 0">{{keyDescription}}</p>
	</div>
	<div style="margin-top: 1em" v-if="!isAdding">
		<button type="button" class="ui button tiny" @click.prevent="add">+</button>
	</div>
</div>`}),Vue.component("http-web-root-box",{props:["v-root-config","v-is-location","v-is-group"],data:function(){let e=this.vRootConfig;return null==(e=null==e?{isPrior:!1,isOn:!0,dir:"",indexes:[],stripPrefix:"",decodePath:!1,isBreak:!1}:e).indexes&&(e.indexes=[]),{rootConfig:e,advancedVisible:!1}},methods:{changeAdvancedVisible:function(e){this.advancedVisible=e},addIndex:function(){let t=this;teaweb.popup("/servers/server/settings/web/createIndex",{height:"10em",callback:function(e){t.rootConfig.indexes.push(e.data.index)}})},removeIndex:function(e){this.rootConfig.indexes.$remove(e)},isOn:function(){return(!this.vIsLocation&&!this.vIsGroup||this.rootConfig.isPrior)&&this.rootConfig.isOn}},template:`<div>
	<input type="hidden" name="rootJSON" :value="JSON.stringify(rootConfig)"/>
	<table class="ui table selectable definition">
		<prior-checkbox :v-config="rootConfig" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || rootConfig.isPrior">
			<tr>
				<td class="title">是否开启静态资源分发</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="rootConfig.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td class="title">静态资源根目录</td>
				<td>
					<input type="text" name="root" v-model="rootConfig.dir" ref="focus" placeholder="类似于 /home/www"/>
					<p class="comment">可以访问此根目录下的静态资源。</p>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-if="isOn()"></more-options-tbody>

		<tbody v-show="isOn() && advancedVisible">
			<tr>
				<td>首页文件</td>
				<td>
					<!-- TODO 支持排序 -->
					<div v-if="rootConfig.indexes.length > 0">
						<div v-for="(index, i) in rootConfig.indexes" class="ui label tiny">
							{{index}} <a href="" title="删除" @click.prevent="removeIndex(i)"><i class="icon remove"></i></a>
						</div>
						<div class="ui divider"></div>
					</div>
					<button class="ui button tiny" type="button" @click.prevent="addIndex()">+</button>
					<p class="comment">在URL中只有目录没有文件名时默认查找的首页文件。</p>
				</td>
			</tr>
			<tr>
				<td>去除URL前缀</td>
				<td>
					<input type="text" v-model="rootConfig.stripPrefix" placeholder="/PREFIX"/>
					<p class="comment">可以把请求的路径部分前缀去除后再查找文件，比如把 <span class="ui label tiny">/web/app/index.html</span> 去除前缀 <span class="ui label tiny">/web</span> 后就变成 <span class="ui label tiny">/app/index.html</span>。 </p>
				</td>
			</tr>
			<tr>
				<td>路径解码</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="rootConfig.decodePath"/>
						<label></label>	
					</div>
					<p class="comment">是否对请求路径进行URL解码，比如把 <span class="ui label tiny">/Web+App+Browser.html</span> 解码成 <span class="ui label tiny">/Web App Browser.html</span> 再查找文件。</p>
				</td>
			</tr>
			<tr>
				<td>是否终止请求</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="rootConfig.isBreak"/>
						<label></label>	
					</div>
					<p class="comment">在找不到要访问的文件的情况下是否终止请求并返回404，如果选择终止请求，则不再尝试反向代理等设置。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`}),Vue.component("http-webp-config-box",{props:["v-webp-config","v-is-location","v-is-group","v-require-cache"],data:function(){let e=this.vWebpConfig;return null==(e=null==e?{isPrior:!1,isOn:!1,quality:50,minLength:{count:0,unit:"kb"},maxLength:{count:0,unit:"kb"},mimeTypes:["image/png","image/jpeg","image/bmp","image/x-ico","image/gif"],extensions:[".png",".jpeg",".jpg",".bmp",".ico"],conds:null}:e).mimeTypes&&(e.mimeTypes=[]),null==e.extensions&&(e.extensions=[]),{config:e,moreOptionsVisible:!1,quality:e.quality}},watch:{quality:function(e){let t=parseInt(e);isNaN(t)?t=90:t<1?t=1:100<t&&(t=100),this.config.quality=t}},methods:{isOn:function(){return(!this.vIsLocation&&!this.vIsGroup||this.config.isPrior)&&this.config.isOn},changeExtensions:function(i){i.forEach(function(e,t){0<e.length&&"."!=e[0]&&(i[t]="."+e)}),this.config.extensions=i},changeMimeTypes:function(e){this.config.mimeTypes=e},changeAdvancedVisible:function(){this.moreOptionsVisible=!this.moreOptionsVisible},changeConds:function(e){this.config.conds=e}},template:`<div>
	<input type="hidden" name="webpJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="config" v-if="vIsLocation || vIsGroup"></prior-checkbox>
		<tbody v-show="(!vIsLocation && !vIsGroup) || config.isPrior">
			<tr>
				<td class="title">启用</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" value="1" v-model="config.isOn"/>
						<label></label>
					</div>
					<p class="comment">选中后表示开启自动WebP压缩<span v-if="vRequireCache">；只有满足缓存条件的图片内容才会被转换</span>。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="isOn()">
			<tr>
				<td>图片质量</td>
				<td>
					<div class="ui input right labeled">
						<input type="text" v-model="quality" style="width: 5em" maxlength="4"/>
						<span class="ui label">%</span>
					</div>
					<p class="comment">取值在0到100之间，数值越大生成的图像越清晰，同时文件尺寸也会越大。</p>
				</td>
			</tr>
			<tr>
				<td>支持的扩展名</td>
				<td>
					<values-box :values="config.extensions" @change="changeExtensions" placeholder="比如 .html"></values-box>
					<p class="comment">含有这些扩展名的URL将会被转成WebP，不区分大小写。</p>
				</td>
			</tr>
			<tr>
				<td>支持的MimeType</td>
				<td>
					<values-box :values="config.mimeTypes" @change="changeMimeTypes" placeholder="比如 text/*"></values-box>
					<p class="comment">响应的Content-Type里包含这些MimeType的内容将会被转成WebP。</p>
				</td>
			</tr>
		</tbody>
		<more-options-tbody @change="changeAdvancedVisible" v-if="isOn()"></more-options-tbody>
		<tbody v-show="isOn() && moreOptionsVisible">
				<tr>
					<td>内容最小长度</td>
				<td>
					<size-capacity-box :v-name="'minLength'" :v-value="config.minLength" :v-unit="'kb'"></size-capacity-box>
					<p class="comment">0表示不限制，内容长度从文件尺寸或Content-Length中获取。</p>
				</td>
			</tr>
			<tr>
				<td>内容最大长度</td>
				<td>
					<size-capacity-box :v-name="'maxLength'" :v-value="config.maxLength" :v-unit="'mb'"></size-capacity-box>
					<p class="comment">0表示不限制，内容长度从文件尺寸或Content-Length中获取。</p>
				</td>
			</tr>
			<tr>
				<td>匹配条件</td>
				<td>
					<http-request-conds-box :v-conds="config.conds" @change="changeConds"></http-request-conds-box>
	</td>
			</tr>
		</tbody>
	</table>			
	<div class="ui margin"></div>
</div>`}),Vue.component("origin-scheduling-view-box",{props:["v-scheduling","v-params"],data:function(){let e=this.vScheduling;return{scheduling:e=null==e?{}:e}},methods:{update:function(){teaweb.popup("/servers/server/settings/reverseProxy/updateSchedulingPopup?"+this.vParams,{height:"21em",callback:function(){window.location.reload()}})}},template:`<div>
	<div class="margin"></div>
	<table class="ui table selectable definition">
		<tr>
			<td class="title">当前正在使用的算法</td>
			<td>
				{{scheduling.name}} &nbsp; <a href="" @click.prevent="update()"><span>[修改]</span></a>
				<p class="comment">{{scheduling.description}}</p>
			</td>
		</tr>
	</table>
</div>`}),Vue.component("http-firewall-block-options",{props:["v-block-options"],data:function(){return{blockOptions:this.vBlockOptions,statusCode:this.vBlockOptions.statusCode,timeout:this.vBlockOptions.timeout,isEditing:!1}},watch:{statusCode:function(e){e=parseInt(e);isNaN(e)?this.blockOptions.statusCode=403:this.blockOptions.statusCode=e},timeout:function(e){e=parseInt(e);isNaN(e)?this.blockOptions.timeout=0:this.blockOptions.timeout=e}},methods:{edit:function(){this.isEditing=!this.isEditing}},template:`<div>
	<input type="hidden" name="blockOptionsJSON" :value="JSON.stringify(blockOptions)"/>
	<a href="" @click.prevent="edit">状态码：{{statusCode}} / 提示内容：<span v-if="blockOptions.body != null && blockOptions.body.length > 0">[{{blockOptions.body.length}}字符]</span><span v-else class="disabled">[无]</span>  / 超时时间：{{timeout}}秒 <i class="icon angle" :class="{up: isEditing, down: !isEditing}"></i></a>
	<table class="ui table" v-show="isEditing">
		<tr>
			<td class="title">状态码</td>
			<td>
				<input type="text" v-model="statusCode" style="width:4.5em" maxlength="3"/>
			</td>
		</tr>
		<tr>
			<td>提示内容</td>
			<td>
				<textarea rows="3" v-model="blockOptions.body"></textarea>
			</td>
		</tr>
		<tr>
			<td>超时时间</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" v-model="timeout" style="width: 5em" maxlength="6"/>
					<span class="ui label">秒</span>
				</div>
				<p class="comment">触发阻止动作时，封锁客户端IP的时间。</p>
			</td>
		</tr>
	</table>
</div>	
`}),Vue.component("http-firewall-rules-box",{props:["v-rules","v-type"],data:function(){let e=this.vRules;return{rules:e=null==e?[]:e}},methods:{addRule:function(){window.UPDATING_RULE=null;let t=this;teaweb.popup("/servers/components/waf/createRulePopup?type="+this.vType,{callback:function(e){t.rules.push(e.data.rule)}})},updateRule:function(t,e){window.UPDATING_RULE=e;let i=this;teaweb.popup("/servers/components/waf/createRulePopup?type="+this.vType,{callback:function(e){Vue.set(i.rules,t,e.data.rule)}})},removeRule:function(e){let t=this;teaweb.confirm("确定要删除此规则吗？",function(){t.rules.$remove(e)})}},template:`<div>
		<input type="hidden" name="rulesJSON" :value="JSON.stringify(rules)"/>
		<div v-if="rules.length > 0">
			<div v-for="(rule, index) in rules" class="ui label small basic" style="margin-bottom: 0.5em">
				{{rule.name}}[{{rule.param}}] 
				
				<!-- cc2 -->
				<span v-if="rule.param == '\${cc2}'">
					{{rule.checkpointOptions.period}}秒/{{rule.checkpointOptions.threshold}}请求
				</span>	
				
				<!-- refererBlock -->
				<span v-if="rule.param == '\${refererBlock}'">
					{{rule.checkpointOptions.allowDomains}}
				</span>
				
				<span v-else>
					<span v-if="rule.paramFilters != null && rule.paramFilters.length > 0" v-for="paramFilter in rule.paramFilters"> | {{paramFilter.code}}</span> <var>{{rule.operator}}</var> {{rule.value}}
				</span>
				
				<!-- description -->
				<span v-if="rule.description != null && rule.description.length > 0" class="grey small">（{{rule.description}}）</span>
				
				<a href="" title="修改" @click.prevent="updateRule(index, rule)"><i class="icon pencil small"></i></a>
				<a href="" title="删除" @click.prevent="removeRule(index)"><i class="icon remove"></i></a>
			</div>
			<div class="ui divider"></div>
		</div>
		<button class="ui button tiny" type="button" @click.prevent="addRule()">+</button>
</div>`}),Vue.component("http-fastcgi-box",{props:["v-fastcgi-ref","v-fastcgi-configs","v-is-location"],data:function(){let e=this.vFastcgiRef,t=(null==e&&(e={isPrior:!1,isOn:!1,fastcgiIds:[]}),this.vFastcgiConfigs);return null==t?t=[]:e.fastcgiIds=t.map(function(e){return e.id}),{fastcgiRef:e,fastcgiConfigs:t,advancedVisible:!1}},methods:{isOn:function(){return(!this.vIsLocation||this.fastcgiRef.isPrior)&&this.fastcgiRef.isOn},createFastcgi:function(){let t=this;teaweb.popup("/servers/server/settings/fastcgi/createPopup",{height:"26em",callback:function(e){teaweb.success("添加成功",function(){t.fastcgiConfigs.push(e.data.fastcgi),t.fastcgiRef.fastcgiIds.push(e.data.fastcgi.id)})}})},updateFastcgi:function(e,t){let i=this;teaweb.popup("/servers/server/settings/fastcgi/updatePopup?fastcgiId="+e,{callback:function(e){teaweb.success("修改成功",function(){Vue.set(i.fastcgiConfigs,t,e.data.fastcgi)})}})},removeFastcgi:function(e){this.fastcgiRef.fastcgiIds.$remove(e),this.fastcgiConfigs.$remove(e)}},template:`<div>
	<input type="hidden" name="fastcgiRefJSON" :value="JSON.stringify(fastcgiRef)"/>
	<table class="ui table definition selectable">
		<prior-checkbox :v-config="fastcgiRef" v-if="vIsLocation"></prior-checkbox>
		<tbody v-show="(!this.vIsLocation || this.fastcgiRef.isPrior)">
			<tr>
				<td class="title">是否启用配置</td>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" v-model="fastcgiRef.isOn"/>
						<label></label>
					</div>
				</td>
			</tr>
		</tbody>
		<tbody v-if="isOn()">
			<tr>
				<td>Fastcgi服务</td>
				<td>
					<div v-show="fastcgiConfigs.length > 0" style="margin-bottom: 0.5em">
						<div class="ui label basic small" :class="{disabled: !fastcgi.isOn}" v-for="(fastcgi, index) in fastcgiConfigs">
							{{fastcgi.address}} &nbsp; <a href="" title="修改" @click.prevent="updateFastcgi(fastcgi.id, index)"><i class="ui icon pencil small"></i></a> &nbsp; <a href="" title="删除" @click.prevent="removeFastcgi(index)"><i class="ui icon remove"></i></a>
						</div>
						<div class="ui divided"></div>
					</div>
					<button type="button" class="ui button tiny" @click.prevent="createFastcgi()">+</button>
				</td>
			</tr>
		</tbody>
	</table>	
	<div class="margin"></div>
</div>`}),Vue.component("http-methods-box",{props:["v-methods"],data:function(){let e=this.vMethods;return{methods:e=null==e?[]:e,isAdding:!1,addingMethod:""}},methods:{add:function(){this.isAdding=!0;let e=this;setTimeout(function(){e.$refs.addingMethod.focus()},100)},confirm:function(){let e=this;this.addingMethod=this.addingMethod.replace(/\s/g,"").toUpperCase(),0==this.addingMethod.length?teaweb.warn("请输入要添加的请求方法",function(){e.$refs.addingMethod.focus()}):this.methods.$contains(this.addingMethod)?teaweb.warn("此请求方法已经存在，无需重复添加",function(){e.$refs.addingMethod.focus()}):(this.methods.push(this.addingMethod),this.cancel())},remove:function(e){this.methods.$remove(e)},cancel:function(){this.isAdding=!1,this.addingMethod=""}},template:`<div>
	<input type="hidden" name="methodsJSON" :value="JSON.stringify(methods)"/>
	<div v-if="methods.length > 0">
		<span class="ui label small basic" v-for="(method, index) in methods">
			{{method}}
			&nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</span>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding">
		<div class="ui fields">
			<div class="ui field">
				<input type="text" v-model="addingMethod" @keyup.enter="confirm()" @keypress.enter.prevent="1" ref="addingMethod" placeholder="如GET" size="10"/>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>
				&nbsp; <a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
		<p class="comment">格式为大写，比如<code-label>GET</code-label>、<code-label>POST</code-label>等。</p>
		<div class="ui divider"></div>
	</div>
	<div style="margin-top: 0.5em" v-if="!isAdding">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`}),Vue.component("http-cond-url-extension",{props:["v-cond"],data:function(){let e={isRequest:!0,param:"${requestPathExtension}",operator:"in",value:"[]"},t=(null!=this.vCond&&this.vCond.param==e.param&&(e.value=this.vCond.value),[]);try{t=JSON.parse(e.value)}catch(e){}return{cond:e,extensions:t,isAdding:!1,addingExt:""}},watch:{extensions:function(){this.cond.value=JSON.stringify(this.extensions)}},methods:{addExt:function(){if(this.isAdding=!this.isAdding,this.isAdding){let e=this;setTimeout(function(){e.$refs.addingExt.focus()},100)}},cancelAdding:function(){this.isAdding=!1,this.addingExt=""},confirmAdding:function(){0!=this.addingExt.length&&("."!=this.addingExt[0]&&(this.addingExt="."+this.addingExt),this.addingExt=this.addingExt.replace(/\s+/g,"").toLowerCase(),this.extensions.push(this.addingExt),this.cancelAdding())},removeExt:function(e){this.extensions.$remove(e)}},template:`<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<div v-if="extensions.length > 0">
		<div class="ui label small" v-for="(ext, index) in extensions">{{ext}} <a href="" title="删除" @click.prevent="removeExt(index)"><i class="icon remove"></i></a></div>
		<div class="ui divider"></div>
	</div>
	<div class="ui fields inline" v-if="isAdding">
		<div class="ui field">
			<input type="text" size="6" maxlength="100" v-model="addingExt" ref="addingExt" placeholder=".xxx" @keyup.enter="confirmAdding" @keypress.enter.prevent="1" />
		</div>
		<div class="ui field">
			<button class="ui button tiny" type="button" @click.prevent="confirmAdding">确认</button>
			<a href="" title="取消" @click.prevent="cancelAdding"><i class="icon remove"></i></a>
		</div> 
	</div>
	<div style="margin-top: 1em">
		<button class="ui button tiny" type="button" @click.prevent="addExt()">+添加扩展名</button>
	</div>
	<p class="comment">扩展名需要包含点（.）符号，例如<span class="ui label tiny">.jpg</span>、<span class="ui label tiny">.png</span>之类。</p>
</div>`}),Vue.component("http-cond-url-not-extension",{props:["v-cond"],data:function(){let e={isRequest:!0,param:"${requestPathExtension}",operator:"not in",value:"[]"},t=(null!=this.vCond&&this.vCond.param==e.param&&(e.value=this.vCond.value),[]);try{t=JSON.parse(e.value)}catch(e){}return{cond:e,extensions:t,isAdding:!1,addingExt:""}},watch:{extensions:function(){this.cond.value=JSON.stringify(this.extensions)}},methods:{addExt:function(){if(this.isAdding=!this.isAdding,this.isAdding){let e=this;setTimeout(function(){e.$refs.addingExt.focus()},100)}},cancelAdding:function(){this.isAdding=!1,this.addingExt=""},confirmAdding:function(){0!=this.addingExt.length&&("."!=this.addingExt[0]&&(this.addingExt="."+this.addingExt),this.addingExt=this.addingExt.replace(/\s+/g,"").toLowerCase(),this.extensions.push(this.addingExt),this.cancelAdding())},removeExt:function(e){this.extensions.$remove(e)}},template:`<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<div v-if="extensions.length > 0">
		<div class="ui label small" v-for="(ext, index) in extensions">{{ext}} <a href="" title="删除" @click.prevent="removeExt(index)"><i class="icon remove"></i></a></div>
		<div class="ui divider"></div>
	</div>
	<div class="ui fields inline" v-if="isAdding">
		<div class="ui field">
			<input type="text" size="6" maxlength="100" v-model="addingExt" ref="addingExt" placeholder=".xxx" @keyup.enter="confirmAdding" @keypress.enter.prevent="1" />
		</div>
		<div class="ui field">
			<button class="ui button tiny" type="button" @click.prevent="confirmAdding">确认</button>
			<a href="" title="取消" @click.prevent="cancelAdding"><i class="icon remove"></i></a>
		</div> 
	</div>
	<div style="margin-top: 1em">
		<button class="ui button tiny" type="button" @click.prevent="addExt()">+添加扩展名</button>
	</div>
	<p class="comment">扩展名需要包含点（.）符号，例如<span class="ui label tiny">.jpg</span>、<span class="ui label tiny">.png</span>之类。</p>
</div>`}),Vue.component("http-cond-url-prefix",{props:["v-cond"],data:function(){let e={isRequest:!0,param:"${requestPath}",operator:"prefix",value:"",isCaseInsensitive:!1};return null!=this.vCond&&"string"==typeof this.vCond.value&&(e.value=this.vCond.value),{cond:e}},methods:{changeCaseInsensitive:function(e){this.cond.isCaseInsensitive=e}},template:`<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value"/>
	<p class="comment">URL前缀，有此前缀的URL都将会被匹配，通常以<code-label>/</code-label>开头，比如<code-label>/static</code-label>。</p>
</div>`}),Vue.component("http-cond-url-not-prefix",{props:["v-cond"],data:function(){let e={isRequest:!0,param:"${requestPath}",operator:"prefix",value:"",isReverse:!0,isCaseInsensitive:!1};return null!=this.vCond&&"string"==typeof this.vCond.value&&(e.value=this.vCond.value),{cond:e}},methods:{changeCaseInsensitive:function(e){this.cond.isCaseInsensitive=e}},template:`<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value"/>
	<p class="comment">要排除的URL前缀，有此前缀的URL都将会被匹配，通常以<code-label>/</code-label>开头，比如<code-label>/static</code-label>。</p>
</div>`}),Vue.component("http-cond-url-eq",{props:["v-cond"],data:function(){let e={isRequest:!0,param:"${requestPath}",operator:"eq",value:"",isCaseInsensitive:!1};return null!=this.vCond&&"string"==typeof this.vCond.value&&(e.value=this.vCond.value),{cond:e}},methods:{changeCaseInsensitive:function(e){this.cond.isCaseInsensitive=e}},template:`<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value"/>
	<p class="comment">完整的URL路径，通常以<code-label>/</code-label>开头，比如<code-label>/static/ui.js</code-label>。</p>
</div>`}),Vue.component("http-cond-url-not-eq",{props:["v-cond"],data:function(){let e={isRequest:!0,param:"${requestPath}",operator:"eq",value:"",isReverse:!0,isCaseInsensitive:!1};return null!=this.vCond&&"string"==typeof this.vCond.value&&(e.value=this.vCond.value),{cond:e}},methods:{changeCaseInsensitive:function(e){this.cond.isCaseInsensitive=e}},template:`<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value"/>
	<p class="comment">要排除的完整的URL路径，通常以<code-label>/</code-label>开头，比如<code-label>/static/ui.js</code-label>。</p>
</div>`}),Vue.component("http-cond-url-regexp",{props:["v-cond"],data:function(){let e={isRequest:!0,param:"${requestPath}",operator:"regexp",value:"",isCaseInsensitive:!1};return null!=this.vCond&&"string"==typeof this.vCond.value&&(e.value=this.vCond.value),{cond:e}},methods:{changeCaseInsensitive:function(e){this.cond.isCaseInsensitive=e}},template:`<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value"/>
	<p class="comment">匹配URL的正则表达式，比如<code-label>^/static/(.*).js$</code-label>。</p>
</div>`}),Vue.component("http-cond-url-not-regexp",{props:["v-cond"],data:function(){let e={isRequest:!0,param:"${requestPath}",operator:"not regexp",value:"",isCaseInsensitive:!1};return null!=this.vCond&&"string"==typeof this.vCond.value&&(e.value=this.vCond.value),{cond:e}},methods:{changeCaseInsensitive:function(e){this.cond.isCaseInsensitive=e}},template:`<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value"/>
	<p class="comment"><strong>不要</strong>匹配URL的正则表达式，意即只要匹配成功则排除此条件，比如<code-label>^/static/(.*).js$</code-label>。</p>
</div>`}),Vue.component("http-cond-user-agent-regexp",{props:["v-cond"],data:function(){let e={isRequest:!0,param:"${userAgent}",operator:"regexp",value:"",isCaseInsensitive:!1};return null!=this.vCond&&"string"==typeof this.vCond.value&&(e.value=this.vCond.value),{cond:e}},methods:{changeCaseInsensitive:function(e){this.cond.isCaseInsensitive=e}},template:`<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value"/>
	<p class="comment">匹配User-Agent的正则表达式，比如<code-label>Android|iPhone</code-label>。</p>
</div>`}),Vue.component("http-cond-user-agent-not-regexp",{props:["v-cond"],data:function(){let e={isRequest:!0,param:"${userAgent}",operator:"not regexp",value:"",isCaseInsensitive:!1};return null!=this.vCond&&"string"==typeof this.vCond.value&&(e.value=this.vCond.value),{cond:e}},methods:{changeCaseInsensitive:function(e){this.cond.isCaseInsensitive=e}},template:`<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<input type="text" v-model="cond.value"/>
	<p class="comment">匹配User-Agent的正则表达式，比如<code-label>Android|iPhone</code-label>，如果匹配，则排除此条件。</p>
</div>`}),Vue.component("http-cond-mime-type",{props:["v-cond"],data:function(){let e={isRequest:!1,param:"${response.contentType}",operator:"mime type",value:"[]"};return null!=this.vCond&&this.vCond.param==e.param&&(e.value=this.vCond.value),{cond:e,mimeTypes:JSON.parse(e.value),isAdding:!1,addingMimeType:""}},watch:{mimeTypes:function(){this.cond.value=JSON.stringify(this.mimeTypes)}},methods:{addMimeType:function(){if(this.isAdding=!this.isAdding,this.isAdding){let e=this;setTimeout(function(){e.$refs.addingMimeType.focus()},100)}},cancelAdding:function(){this.isAdding=!1,this.addingMimeType=""},confirmAdding:function(){0!=this.addingMimeType.length&&(this.addingMimeType=this.addingMimeType.replace(/\s+/g,""),this.mimeTypes.push(this.addingMimeType),this.cancelAdding())},removeMimeType:function(e){this.mimeTypes.$remove(e)}},template:`<div>
	<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
	<div v-if="mimeTypes.length > 0">
		<div class="ui label small" v-for="(mimeType, index) in mimeTypes">{{mimeType}} <a href="" title="删除" @click.prevent="removeMimeType(index)"><i class="icon remove"></i></a></div>
		<div class="ui divider"></div>
	</div>
	<div class="ui fields inline" v-if="isAdding">
		<div class="ui field">
			<input type="text" size="16" maxlength="100" v-model="addingMimeType" ref="addingMimeType" placeholder="类似于image/png" @keyup.enter="confirmAdding" @keypress.enter.prevent="1" />
		</div>
		<div class="ui field">
			<button class="ui button tiny" type="button" @click.prevent="confirmAdding">确认</button>
			<a href="" title="取消" @click.prevent="cancelAdding"><i class="icon remove"></i></a>
		</div> 
	</div>
	<div style="margin-top: 1em">
		<button class="ui button tiny" type="button" @click.prevent="addMimeType()">+添加MimeType</button>
	</div>
	<p class="comment">服务器返回的内容的MimeType，比如<span class="ui label tiny">text/html</span>、<span class="ui label tiny">image/*</span>等。</p>
</div>`}),Vue.component("http-cond-params",{props:["v-cond"],mounted:function(){let i=this.vCond;if(null!=i)if(this.operator=i.operator,["regexp","not regexp","eq","not","prefix","suffix","contains","not contains","eq ip","gt ip","gte ip","lt ip","lte ip","ip range"].$contains(i.operator))this.stringValue=i.value;else if(["eq int","eq float","gt","gte","lt","lte","mod 10","ip mod 10","mod 100","ip mod 100"].$contains(i.operator))this.numberValue=i.value;else{var e;if(["mod","ip mod"].$contains(i.operator))return e=i.value.split(","),this.modDivValue=e[0],void(1<e.length&&(this.modRemValue=e[1]));let t=this;if(["in","not in","file ext","mime type"].$contains(i.operator))try{let e=JSON.parse(i.value);null!=e&&e instanceof Array&&e.forEach(function(e){t.stringValues.push(e)})}catch(e){}else["version range"].$contains(i.operator)&&(e=i.value.split(","),this.versionRangeMinValue=e[0],1<e.length&&(this.versionRangeMaxValue=e[1]))}},data:function(){let e={isRequest:!0,param:"",operator:window.REQUEST_COND_OPERATORS[0].op,value:"",isCaseInsensitive:!1};return{cond:e=null!=this.vCond?this.vCond:e,operators:window.REQUEST_COND_OPERATORS,operator:window.REQUEST_COND_OPERATORS[0].op,operatorDescription:window.REQUEST_COND_OPERATORS[0].description,variables:window.REQUEST_VARIABLES,variable:"",stringValue:"",numberValue:"",modDivValue:"",modRemValue:"",stringValues:[],versionRangeMinValue:"",versionRangeMaxValue:""}},methods:{changeVariable:function(){let e=this.cond.param;null==e&&(e=""),this.cond.param=e+this.variable},changeOperator:function(){let t=this,i=(this.operators.forEach(function(e){e.op==t.operator&&(t.operatorDescription=e.description)}),this.cond.operator=this.operator,document.getElementById("variables-value-box"));null!=i&&setTimeout(function(){let e=i.getElementsByTagName("INPUT");0<e.length&&e[0].focus()},100)},changeStringValues:function(e){this.stringValues=e,this.cond.value=JSON.stringify(e)}},watch:{stringValue:function(e){this.cond.value=e},numberValue:function(e){this.cond.value=e},modDivValue:function(t){if(0!=t.length){let e=parseInt(t);isNaN(e)&&(e=1),this.modDivValue=e,this.cond.value=e+","+this.modRemValue}},modRemValue:function(t){if(0!=t.length){let e=parseInt(t);isNaN(e)&&(e=0),this.modRemValue=e,this.cond.value=this.modDivValue+","+e}},versionRangeMinValue:function(e){this.cond.value=this.versionRangeMinValue+","+this.versionRangeMaxValue},versionRangeMaxValue:function(e){this.cond.value=this.versionRangeMinValue+","+this.versionRangeMaxValue}},template:`<tbody>
	<tr>
		<td>参数值</td>
		<td>
			<input type="hidden" name="condJSON" :value="JSON.stringify(cond)"/>
			<div>
				<div class="ui fields inline">
					<div class="ui field">
						<input type="text" placeholder="\${xxx}" v-model="cond.param"/>
					</div>
					<div class="ui field">
						<select class="ui dropdown" style="width: 7em; color: grey" v-model="variable" @change="changeVariable">
							<option value="">[常用参数]</option>
							<option v-for="v in variables" :value="v.code">{{v.code}} - {{v.name}}</option>
						</select>
					</div>
				</div>			
			</div>
			<p class="comment">其中可以使用变量，类似于<code-label>\${requestPath}</code-label>，也可以是多个变量的组合。</p>
		</td>
	</tr>
	<tr>
		<td>操作符</td>
		<td>
			<div>
				<select class="ui dropdown auto-width" v-model="operator" @change="changeOperator">
					<option v-for="operator in operators" :value="operator.op">{{operator.name}}</option>
				</select>
				<p class="comment">{{operatorDescription}}</p>
			</div>
		</td>
	</tr>
	<tr v-show="!['file exist', 'file not exist'].$contains(cond.operator)">
		<td>对比值</td>
		<td id="variables-value-box">
			<!-- 正则表达式 -->
			<div v-if="['regexp', 'not regexp'].$contains(cond.operator)">
				<input type="text" v-model="stringValue"/>
				<p class="comment">要匹配的正则表达式，比如<code-label>^/static/(.+).js</code-label>。</p>
			</div>
			
			<!-- 数字相关 -->
			<div v-if="['eq int', 'eq float', 'gt', 'gte', 'lt', 'lte'].$contains(cond.operator)">
				<input type="text" maxlength="11" size="11" style="width: 5em" v-model="numberValue"/>
				<p class="comment">要对比的数字。</p>
			</div>
			
			<!-- 取模 -->
			<div v-if="['mod 10'].$contains(cond.operator)">
				<input type="text" maxlength="11" size="11" style="width: 5em" v-model="numberValue"/>
				<p class="comment">参数值除以10的余数，在0-9之间。</p>
			</div>
			<div v-if="['mod 100'].$contains(cond.operator)">
				<input type="text" maxlength="11" size="11" style="width: 5em" v-model="numberValue"/>
				<p class="comment">参数值除以100的余数，在0-99之间。</p>
			</div>
			<div v-if="['mod', 'ip mod'].$contains(cond.operator)">
				<div class="ui fields inline">
					<div class="ui field">除：</div>
					<div class="ui field">
						<input type="text" maxlength="11" size="11" style="width: 5em" v-model="modDivValue" placeholder="除数"/>
					</div>
					<div class="ui field">余：</div>
					<div class="ui field">
						<input type="text" maxlength="11" size="11" style="width: 5em" v-model="modRemValue" placeholder="余数"/>
					</div>
				</div>
			</div>
			
			<!-- 字符串相关 -->
			<div v-if="['eq', 'not', 'prefix', 'suffix', 'contains', 'not contains'].$contains(cond.operator)">
				<input type="text" v-model="stringValue"/>
				<p class="comment" v-if="cond.operator == 'eq'">和参数值一致的字符串。</p>
				<p class="comment" v-if="cond.operator == 'not'">和参数值不一致的字符串。</p>
				<p class="comment" v-if="cond.operator == 'prefix'">参数值的前缀。</p>
				<p class="comment" v-if="cond.operator == 'suffix'">参数值的后缀为此字符串。</p>
				<p class="comment" v-if="cond.operator == 'contains'">参数值包含此字符串。</p>
				<p class="comment" v-if="cond.operator == 'not contains'">参数值不包含此字符串。</p>
			</div>
			<div v-if="['in', 'not in', 'file ext', 'mime type'].$contains(cond.operator)">
				<values-box @change="changeStringValues" :values="stringValues" size="15"></values-box>
				<p class="comment" v-if="cond.operator == 'in'">添加参数值列表。</p>
				<p class="comment" v-if="cond.operator == 'not in'">添加参数值列表。</p>
				<p class="comment" v-if="cond.operator == 'file ext'">添加扩展名列表，比如<code-label>png</code-label>、<code-label>html</code-label>，不包括点。</p>
				<p class="comment" v-if="cond.operator == 'mime type'">添加MimeType列表，类似于<code-label>text/html</code-label>、<code-label>image/*</code-label>。</p>
			</div>
			<div v-if="['version range'].$contains(cond.operator)">
				<div class="ui fields inline">
					<div class="ui field"><input type="text" v-model="versionRangeMinValue" maxlength="200" placeholder="最小版本" style="width: 10em"/></div>
					<div class="ui field">-</div>
					<div class="ui field"><input type="text" v-model="versionRangeMaxValue" maxlength="200" placeholder="最大版本" style="width: 10em"/></div>
				</div>
			</div>
			
			<!-- IP相关 -->
			<div v-if="['eq ip', 'gt ip', 'gte ip', 'lt ip', 'lte ip', 'ip range'].$contains(cond.operator)">
				<input type="text" style="width: 10em" v-model="stringValue" placeholder="x.x.x.x"/>
				<p class="comment">要对比的IP。</p>
			</div>
			<div v-if="['ip mod 10'].$contains(cond.operator)">
				<input type="text" maxlength="11" size="11" style="width: 5em" v-model="numberValue"/>
				<p class="comment">参数中IP转换成整数后除以10的余数，在0-9之间。</p>
			</div>
			<div v-if="['ip mod 100'].$contains(cond.operator)">
				<input type="text" maxlength="11" size="11" style="width: 5em" v-model="numberValue"/>
				<p class="comment">参数中IP转换成整数后除以100的余数，在0-99之间。</p>
			</div>
		</td>
	</tr>
	<tr v-if="['regexp', 'not regexp', 'eq', 'not', 'prefix', 'suffix', 'contains', 'not contains', 'in', 'not in'].$contains(cond.operator)">
		<td>不区分大小写</td>
		<td>
		   <div class="ui checkbox">
				<input type="checkbox" v-model="cond.isCaseInsensitive"/>
				<label></label>
			</div>
			<p class="comment">选中后表示对比时忽略参数值的大小写。</p>
		</td>
	</tr>
</tbody>`}),Vue.component("http-status-box",{props:["v-status-list"],data:function(){let e=this.vStatusList;return{statusList:e=null==e?[]:e,isAdding:!1,addingStatus:""}},methods:{add:function(){this.isAdding=!0;let e=this;setTimeout(function(){e.$refs.addingStatus.focus()},100)},confirm:function(){let e=this;this.addingStatus=this.addingStatus.replace(/\s/g,"").toUpperCase(),0==this.addingStatus.length?teaweb.warn("请输入要添加的状态码",function(){e.$refs.addingStatus.focus()}):this.statusList.$contains(this.addingStatus)?teaweb.warn("此状态码已经存在，无需重复添加",function(){e.$refs.addingStatus.focus()}):this.addingStatus.match(/^\d{3}$/)?(this.statusList.push(parseInt(this.addingStatus,10)),this.cancel()):teaweb.warn("请输入正确的状态码",function(){e.$refs.addingStatus.focus()})},remove:function(e){this.statusList.$remove(e)},cancel:function(){this.isAdding=!1,this.addingStatus=""}},template:`<div>
	<input type="hidden" name="statusListJSON" :value="JSON.stringify(statusList)"/>
	<div v-if="statusList.length > 0">
		<span class="ui label small basic" v-for="(status, index) in statusList">
			{{status}}
			&nbsp; <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove small"></i></a>
		</span>
		<div class="ui divider"></div>
	</div>
	<div v-if="isAdding">
		<div class="ui fields">
			<div class="ui field">
				<input type="text" v-model="addingStatus" @keyup.enter="confirm()" @keypress.enter.prevent="1" ref="addingStatus" placeholder="如200" size="3" maxlength="3" style="width: 5em"/>
			</div>
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>
				&nbsp; <a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
			</div>
		</div>
		<p class="comment">格式为三位数字，比如<code-label>200</code-label>、<code-label>404</code-label>等。</p>
		<div class="ui divider"></div>
	</div>
	<div style="margin-top: 0.5em" v-if="!isAdding">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`}),Vue.component("server-group-selector",{props:["v-groups"],data:function(){let e=this.vGroups;return{groups:e=null==e?[]:e}},methods:{selectGroup:function(){let t=this;var e=this.groups.map(function(e){return e.id.toString()}).join(",");teaweb.popup("/servers/groups/selectPopup?selectedGroupIds="+e,{callback:function(e){t.groups.push(e.data.group)}})},addGroup:function(){let t=this;teaweb.popup("/servers/groups/createPopup",{callback:function(e){t.groups.push(e.data.group)}})},removeGroup:function(e){this.groups.$remove(e)},groupIds:function(){return this.groups.map(function(e){return e.id})}},template:`<div>
	<div v-if="groups.length > 0">
		<div class="ui label small basic" v-if="groups.length > 0" v-for="(group, index) in groups">
			<input type="hidden" name="groupIds" :value="group.id"/>
			{{group.name}} &nbsp;<a href="" title="删除" @click.prevent="removeGroup(index)"><i class="icon remove"></i></a>
		</div>
		<div class="ui divider"></div>
	</div>
	<div>
		<a href="" @click.prevent="selectGroup()">[选择分组]</a> &nbsp; <a href="" @click.prevent="addGroup()">[添加分组]</a>
	</div>
</div>`}),Vue.component("script-group-config-box",{props:["v-group","v-is-location"],data:function(){let e=this.vGroup,t=(null==(e=null==e?{isPrior:!1,isOn:!0,scripts:[]}:e).scripts&&(e.scripts=[]),null);return 0<e.scripts.length&&(t=e.scripts[e.scripts.length-1]),{group:e,script:t}},methods:{changeScript:function(e){this.group.scripts=[e],this.change()},change:function(){this.$emit("change",this.group)}},template:`<div>
		<table class="ui table definition selectable">
			<prior-checkbox :v-config="group" v-if="vIsLocation"></prior-checkbox>
		</table>
		<div :style="{opacity: (!vIsLocation || group.isPrior) ? 1 : 0.5}">
			<script-config-box :v-script-config="script" comment="在接收到客户端请求之后立即调用。预置req、resp变量。" @change="changeScript" :v-is-location="vIsLocation"></script-config-box>
		</div>
</div>`}),Vue.component("metric-period-config-box",{props:["v-period","v-period-unit"],data:function(){let e=this.vPeriod,t=this.vPeriodUnit;return null!=e&&0!=e.toString().length||(e=1),null!=t&&0!=t.length||(t="day"),{periodConfig:{period:e,unit:t}}},watch:{"periodConfig.period":function(e){e=parseInt(e),(isNaN(e)||e<=0)&&(e=1),this.periodConfig.period=e}},template:`<div>
	<input type="hidden" name="periodJSON" :value="JSON.stringify(periodConfig)"/>
	<div class="ui fields inline">
		<div class="ui field">
			<input type="text" v-model="periodConfig.period" maxlength="4" size="4"/>
		</div>
		<div class="ui field">
			<select class="ui dropdown" v-model="periodConfig.unit">
				<option value="minute">分钟</option>
				<option value="hour">小时</option>
				<option value="day">天</option>
				<option value="week">周</option>
				<option value="month">月</option>
			</select>
		</div>
	</div>
	<p class="comment">在此周期内同一对象累积为同一数据。</p>
</div>`}),Vue.component("traffic-limit-config-box",{props:["v-traffic-limit"],data:function(){let e=this.vTrafficLimit;return null==(e=null==e?{isOn:!1,dailySize:{count:-1,unit:"gb"},monthlySize:{count:-1,unit:"gb"},totalSize:{count:-1,unit:"gb"},noticePageBody:""}:e).dailySize&&(e.dailySize={count:-1,unit:"gb"}),null==e.monthlySize&&(e.monthlySize={count:-1,unit:"gb"}),null==e.totalSize&&(e.totalSize={count:-1,unit:"gb"}),{config:e}},methods:{showBodyTemplate:function(){this.config.noticePageBody=`<!DOCTYPE html>
<html>
<head>
<title>Traffic Limit Exceeded Warning</title>
<body>

<h1>Traffic Limit Exceeded Warning</h1>
<p>The site traffic has exceeded the limit. Please contact with the site administrator.</p>
<address>Request ID: \${requestId}.</address>

</body>
</html>`}},template:`<div>
	<input type="hidden" name="trafficLimitJSON" :value="JSON.stringify(config)"/>
	<table class="ui table selectable definition">
		<tbody>
			<tr>
				<td class="title">是否启用</td>
				<td>
					<checkbox v-model="config.isOn"></checkbox>
					<p class="comment">注意：由于流量统计是每5分钟统计一次，所以超出流量限制后，对用户的提醒也会有所延迟。</p>
				</td>
			</tr>
		</tbody>
		<tbody v-show="config.isOn">
			<tr>
				<td>日流量限制</td>
				<td>
					<size-capacity-box :v-value="config.dailySize"></size-capacity-box>
				</td>
			</tr>
			<tr>
				<td>月流量限制</td>
				<td>
					<size-capacity-box :v-value="config.monthlySize"></size-capacity-box>
				</td>
			</tr>
			<!--<tr>
				<td>总体限制</td>
				<td>
					<size-capacity-box :v-value="config.totalSize"></size-capacity-box>
					<p class="comment"></p>
				</td>
			</tr>-->
			<tr>
				<td>网页提示内容</td>
				<td>
					<textarea v-model="config.noticePageBody"></textarea>
					<p class="comment"><a href="" @click.prevent="showBodyTemplate">[使用模板]</a>。当达到流量限制时网页显示的HTML内容，不填写则显示默认的提示内容。</p>
				</td>
			</tr>
		</tbody>
	</table>
	<div class="margin"></div>
</div>`}),Vue.component("http-firewall-captcha-options",{props:["v-captcha-options"],mounted:function(){this.updateSummary()},data:function(){let e=this.vCaptchaOptions;return(e=null==e?{countLetters:0,life:0,maxFails:0,failBlockTimeout:0,failBlockScopeAll:!1,uiIsOn:!1,uiTitle:"",uiPrompt:"",uiButtonTitle:"",uiShowRequestId:!0,uiCss:"",uiFooter:"",uiBody:"",cookieId:"",lang:""}:e).countLetters<=0&&(e.countLetters=6),{options:e,isEditing:!1,summary:""}},watch:{"options.countLetters":function(e){let t=parseInt(e,10);isNaN(t)||t<0?t=0:10<t&&(t=10),this.options.countLetters=t},"options.life":function(e){let t=parseInt(e,10);isNaN(t)&&(t=0),this.options.life=t,this.updateSummary()},"options.maxFails":function(e){let t=parseInt(e,10);isNaN(t)&&(t=0),this.options.maxFails=t,this.updateSummary()},"options.failBlockTimeout":function(e){let t=parseInt(e,10);isNaN(t)&&(t=0),this.options.failBlockTimeout=t,this.updateSummary()},"options.failBlockScopeAll":function(e){this.updateSummary()},"options.uiIsOn":function(e){this.updateSummary()}},methods:{edit:function(){this.isEditing=!this.isEditing},updateSummary:function(){let e=[];0<this.options.life&&e.push("有效时间"+this.options.life+"秒"),0<this.options.maxFails&&e.push("最多失败"+this.options.maxFails+"次"),0<this.options.failBlockTimeout&&e.push("失败拦截"+this.options.failBlockTimeout+"秒"),this.options.failBlockScopeAll&&e.push("全局封禁"),this.options.uiIsOn&&e.push("定制UI"),0==e.length?this.summary="默认配置":this.summary=e.join(" / ")},confirm:function(){this.isEditing=!1}},template:`<div>
	<input type="hidden" name="captchaOptionsJSON" :value="JSON.stringify(options)"/>
	<a href="" @click.prevent="edit">{{summary}} <i class="icon angle" :class="{up: isEditing, down: !isEditing}"></i></a>
	<div v-show="isEditing" style="margin-top: 0.5em">
		<table class="ui table definition selectable">
			<tbody>
				<tr>
					<td class="title">有效时间</td>
					<td>
						<div class="ui input right labeled">
							<input type="text" style="width: 5em" maxlength="9" v-model="options.life" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
							<span class="ui label">秒</span>
						</div>
						<p class="comment">验证通过后在这个时间内不再验证，默认600秒。</p>
					</td>
				</tr>
				<tr>
					<td>最多失败次数</td>
					<td>
						<div class="ui input right labeled">
							<input type="text" style="width: 5em" maxlength="9" v-model="options.maxFails" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
							<span class="ui label">次</span>
						</div>
						<p class="comment">允许用户失败尝试的最多次数，超过这个次数将被自动加入黑名单。如果为空或者为0，表示不限制。</p>
					</td>
				</tr>
				<tr>
					<td>失败拦截时间</td>
					<td>
						<div class="ui input right labeled">
							<input type="text" style="width: 5em" maxlength="9" v-model="options.failBlockTimeout" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
							<span class="ui label">秒</span>
						</div>
						<p class="comment">在达到最多失败次数（大于0）时，自动拦截的时间；如果为0表示不自动拦截。</p>
					</td>
				</tr>
				<tr>
					<td>失败全局封禁</td>
					<td>
						<checkbox v-model="options.failBlockScopeAll"></checkbox>
						<p class="comment">是否在失败时全局封禁，默认为只封禁对单个网站服务的访问。</p>
					</td>
				</tr>
				<tr>
					<td>验证码中数字个数</td>
					<td>
						<select class="ui dropdown auto-width" v-model="options.countLetters">
							<option v-for="i in 10" :value="i">{{i}}</option>
						</select>
					</td>
				</tr>
				<tr>
					<td class="color-border">定制UI</td>
					<td><checkbox v-model="options.uiIsOn"></checkbox></td>
				</tr>
			</tbody>
			<tbody v-show="options.uiIsOn">
				<tr>
					<td class="color-border">页面标题</td>
					<td>
						<input type="text" v-model="options.uiTitle" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
					</td>
				</tr>
				<tr>
					<td class="color-border">按钮标题</td>
					<td>
						<input type="text" v-model="options.uiButtonTitle" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
						<p class="comment">类似于<code-label>提交验证</code-label>。</p>
					</td>
				</tr>
				<tr>
					<td class="color-border">显示请求ID</td>
					<td>
						<checkbox v-model="options.uiShowRequestId"></checkbox>
						<p class="comment">在界面上显示请求ID，方便用户报告问题。</p>
					</td>
				</tr>
				<tr>
					<td class="color-border">CSS样式</td>
					<td>
						<textarea spellcheck="false" v-model="options.uiCss" rows="2"></textarea>
					</td>
				</tr>
				<tr>
					<td class="color-border">页头提示</td>
					<td>
						<textarea spellcheck="false" v-model="options.uiPrompt" rows="2"></textarea>
						<p class="comment">类似于<code-label>请输入上面的验证码</code-label>，支持HTML。</p>
					</td>
				</tr>
				<tr>
					<td class="color-border">页尾提示</td>
					<td>
						<textarea spellcheck="false" v-model="options.uiFooter" rows="2"></textarea>
						<p class="comment">支持HTML。</p>
					</td>
				</tr>
				<tr>
					<td class="color-border">页面模板</td>
					<td>
						<textarea spellcheck="false" rows="2" v-model="options.uiBody"></textarea>
						<p class="comment"><span v-if="options.uiBody.length > 0 && options.uiBody.indexOf('\${body}') < 0 " class="red">模板中必须包含\${body}表示验证码表单！</span>整个页面的模板，支持HTML，其中必须使用<code-label>\${body}</code-label>变量代表验证码表单，否则将无法正常显示验证码。</p>
					</td>
				</tr>
			</tbody>
		</table>
	</div>
</div>
`}),Vue.component("firewall-syn-flood-config-box",{props:["v-syn-flood-config"],data:function(){let e=this.vSynFloodConfig;return{config:e=null==e?{isOn:!1,minAttempts:10,timeoutSeconds:600,ignoreLocal:!0}:e,isEditing:!1,minAttempts:e.minAttempts,timeoutSeconds:e.timeoutSeconds}},methods:{edit:function(){this.isEditing=!this.isEditing}},watch:{minAttempts:function(e){let t=parseInt(e);(t=isNaN(t)?10:t)<5&&(t=5),this.config.minAttempts=t},timeoutSeconds:function(e){let t=parseInt(e);(t=isNaN(t)?10:t)<60&&(t=60),this.config.timeoutSeconds=t}},template:`<div>
	<input type="hidden" name="synFloodJSON" :value="JSON.stringify(config)"/>
	<a href="" @click.prevent="edit">
		<span v-if="config.isOn">
			已启用 / <span>空连接次数：{{config.minAttempts}}次/分钟</span> / 封禁时间：{{config.timeoutSeconds}}秒 <span v-if="config.ignoreLocal">/ 忽略局域网访问</span>
		</span>
		<span v-else>未启用</span>
		<i class="icon angle" :class="{up: isEditing, down: !isEditing}"></i>
	</a>
	
	<table class="ui table selectable" v-show="isEditing">
		<tr>
			<td class="title">是否启用</td>
			<td>
				<checkbox v-model="config.isOn"></checkbox>
				<p class="comment">启用后，WAF将会尝试自动检测并阻止SYN Flood攻击。此功能需要节点已安装并启用Firewalld。</p>
			</td>
		</tr>
		<tr>
			<td>空连接次数</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" v-model="minAttempts" style="width: 5em" maxlength="4"/>
					<span class="ui label">次/分钟</span>
				</div>
				<p class="comment">超过此数字的"空连接"将被视为SYN Flood攻击，为了防止误判，此数值默认不小于5。</p>
			</td>
		</tr>
		<tr>
			<td>封禁时间</td>
			<td>
				<div class="ui input right labeled">
					<input type="text" v-model="timeoutSeconds" style="width: 5em" maxlength="8"/>
					<span class="ui label">秒</span>
				</div>
			</td>
		</tr>
		<tr>
			<td>忽略局域网访问</td>
			<td>
				<checkbox v-model="config.ignoreLocal"></checkbox>
			</td>
		</tr>
	</table>
</div>`}),Vue.component("admin-selector",{props:["v-admin-id"],mounted:function(){let t=this;Tea.action("/admins/options").post().success(function(e){t.admins=e.data.admins})},data:function(){let e=this.vAdminId;return{admins:[],adminId:e=null==e?0:e}},template:`<div>
    <select class="ui dropdown auto-width" name="adminId" v-model="adminId">
        <option value="0">[选择系统用户]</option>
        <option v-for="admin in admins" :value="admin.id">{{admin.name}}（{{admin.username}}）</option>
    </select>
</div>`}),Vue.component("ip-list-bind-box",{props:["v-http-firewall-policy-id","v-type"],mounted:function(){this.refresh()},data:function(){return{policyId:this.vHttpFirewallPolicyId,type:this.vType,lists:[]}},methods:{bind:function(){let e=this;teaweb.popup("/servers/iplists/bindHTTPFirewallPopup?httpFirewallPolicyId="+this.policyId+"&type="+this.type,{width:"50em",height:"34em",callback:function(){},onClose:function(){e.refresh()}})},remove:function(t,e){let i=this;teaweb.confirm("确定要删除这个绑定的IP名单吗？",function(){Tea.action("/servers/iplists/unbindHTTPFirewall").params({httpFirewallPolicyId:i.policyId,listId:e}).post().success(function(e){i.lists.$remove(t)})})},refresh:function(){let t=this;Tea.action("/servers/iplists/httpFirewall").params({httpFirewallPolicyId:this.policyId,type:this.vType}).post().success(function(e){t.lists=e.data.lists})}},template:`<div>
	<a href="" @click.prevent="bind()" style="color: rgba(0,0,0,.6)">绑定+</a> &nbsp; <span v-if="lists.length > 0"><span class="disabled small">|&nbsp;</span> 已绑定：</span>
	<div class="ui label basic small" v-for="(list, index) in lists">
		<a :href="'/servers/iplists/list?listId=' + list.id" title="点击查看详情" style="opacity: 1">{{list.name}}</a>
		<a href="" title="删除" @click.prevent="remove(index, list.id)"><i class="icon remove small"></i></a>
	</div>
</div>`}),Vue.component("ip-list-table",{props:["v-items","v-keyword","v-show-search-button"],data:function(){return{items:this.vItems,keyword:null!=this.vKeyword?this.vKeyword:"",selectedAll:!1,hasSelectedItems:!1}},methods:{updateItem:function(e){this.$emit("update-item",e)},deleteItem:function(e){this.$emit("delete-item",e)},viewLogs:function(e){teaweb.popup("/servers/iplists/accessLogsPopup?itemId="+e,{width:"50em",height:"30em"})},changeSelectedAll:function(){let e=this.$refs.itemCheckBox;if(null!=e){let t=this;e.forEach(function(e){e.checked=t.selectedAll}),this.hasSelectedItems=this.selectedAll}},changeSelected:function(e){let t=this,i=(t.hasSelectedItems=!1,t.$refs.itemCheckBox);null!=i&&i.forEach(function(e){e.checked&&(t.hasSelectedItems=!0)})},deleteAll:function(){let e=this.$refs.itemCheckBox;if(null!=e){let t=[];e.forEach(function(e){e.checked&&t.push(e.value)}),0!=t.length&&Tea.action("/servers/iplists/deleteItems").post().params({itemIds:t}).success(function(){teaweb.successToast("批量删除成功",1200,teaweb.reload)})}},formatSeconds:function(e){return e<60?e+"秒":e<3600?Math.ceil(e/60)+"分钟":e<86400?Math.ceil(e/3600)+"小时":Math.ceil(e/86400)+"天"}},template:`<div>
 <div v-show="hasSelectedItems">
 	<a href="" @click.prevent="deleteAll">[批量删除]</a>
</div>
 <table class="ui table selectable celled" v-if="items.length > 0">
        <thead>
            <tr>
            	<th style="width: 1em">
            		<div class="ui checkbox">
						<input type="checkbox" v-model="selectedAll" @change="changeSelectedAll"/>
						<label></label>
					</div>
				</th>
                <th style="width:18em">IP</th>
                <th style="width: 6em">类型</th>
                <th style="width: 6em">级别</th>
                <th style="width: 12em">过期时间</th>
                <th>备注</th>
                <th class="three op">操作</th>
            </tr>
        </thead>
		<tbody v-for="item in items">
			<tr>
				<td>
					<div class="ui checkbox">
						<input type="checkbox" :value="item.id" @change="changeSelected" ref="itemCheckBox"/>
						<label></label>
					</div>
				</td>
				<td>
					<span v-if="item.type != 'all'" :class="{green: item.list != null && item.list.type == 'white'}">
					<keyword :v-word="keyword">{{item.ipFrom}}</keyword> <span> <span class="small red" v-if="item.isRead != null && !item.isRead">&nbsp;New&nbsp;</span>&nbsp;<a :href="'/servers/iplists?ip=' + item.ipFrom" v-if="vShowSearchButton" title="搜索此IP"><span><i class="icon search small" style="color: #ccc"></i></span></a></span>
					<span v-if="item.ipTo.length > 0"> - <keyword :v-word="keyword">{{item.ipTo}}</keyword></span></span>
					<span v-else class="disabled">*</span>
					<div v-if="item.region != null && item.region.length > 0">
						<span class="grey small">{{item.region}}</span>
						<span v-if="item.isp != null && item.isp.length > 0 && item.isp != '内网IP'" class="grey small"><span class="disabled">|</span> {{item.isp}}</span>
					</div>
					<div v-if="item.createdTime != null">
						<span class="small grey">添加于 {{item.createdTime}}
							<span v-if="item.list != null && item.list.id > 0">
								@ 
								
								<a :href="'/servers/iplists/list?listId=' + item.list.id" v-if="item.policy.id == 0"><span>[<span v-if="item.list.type == 'black'">黑</span><span v-if="item.list.type == 'white'">白</span>名单：{{item.list.name}}]</span></a>
								<span v-else>[<span v-if="item.list.type == 'black'">黑</span><span v-if="item.list.type == 'white'">白</span>名单：{{item.list.name}}</span>
								
								<span v-if="item.policy.id > 0">
									<span v-if="item.policy.server != null">
										<a :href="'/servers/server/settings/waf/ipadmin/allowList?serverId=' + item.policy.server.id + '&firewallPolicyId=' + item.policy.id" v-if="item.list.type == 'white'">[服务：{{item.policy.server.name}}]</a>
										<a :href="'/servers/server/settings/waf/ipadmin/denyList?serverId=' + item.policy.server.id + '&firewallPolicyId=' + item.policy.id" v-if="item.list.type == 'black'">[服务：{{item.policy.server.name}}]</a>
									</span>
									<span v-else>
										<a :href="'/servers/components/waf/ipadmin/lists?firewallPolicyId=' + item.policy.id +  '&type=' + item.list.type">[策略：{{item.policy.name}}]</a>
									</span>
								</span>
							</span>
						</span>
					</div>
				</td>
				<td>
					<span v-if="item.type.length == 0">IPv4</span>
					<span v-else-if="item.type == 'ipv4'">IPv4</span>
					<span v-else-if="item.type == 'ipv6'">IPv6</span>
					<span v-else-if="item.type == 'all'"><strong>所有IP</strong></span>
				</td>
				<td>
					<span v-if="item.eventLevelName != null && item.eventLevelName.length > 0">{{item.eventLevelName}}</span>
					<span v-else class="disabled">-</span>
				</td>
				<td>
					<div v-if="item.expiredTime.length > 0">
						{{item.expiredTime}}
						<div v-if="item.isExpired" style="margin-top: 0.5em">
							<span class="ui label tiny basic red">已过期</span>
						</div>
						<div  v-if="item.lifeSeconds != null && item.lifeSeconds > 0">
							<span class="small grey">{{formatSeconds(item.lifeSeconds)}}</span>
						</div>
					</div>
					<span v-else class="disabled">不过期</span>
				</td>
				<td>
					<span v-if="item.reason.length > 0">{{item.reason}}</span>
					<span v-else class="disabled">-</span>
					
					<div v-if="item.sourceNode != null && item.sourceNode.id > 0" style="margin-top: 0.4em">
						<a :href="'/clusters/cluster/node?clusterId=' + item.sourceNode.clusterId + '&nodeId=' + item.sourceNode.id"><span class="small"><i class="icon cloud"></i>{{item.sourceNode.name}}</span></a>
					</div>
					<div style="margin-top: 0.4em" v-if="item.sourceServer != null && item.sourceServer.id > 0">
						<a :href="'/servers/server?serverId=' + item.sourceServer.id" style="border: 0"><span class="small "><i class="icon clone outline"></i>{{item.sourceServer.name}}</span></a>
					</div>
					<div v-if="item.sourcePolicy != null && item.sourcePolicy.id > 0" style="margin-top: 0.4em">
						<a :href="'/servers/components/waf/group?firewallPolicyId=' +  item.sourcePolicy.id + '&type=inbound&groupId=' + item.sourceGroup.id + '#set' + item.sourceSet.id" v-if="item.sourcePolicy.serverId == 0"><span class="small "><i class="icon shield"></i>{{item.sourcePolicy.name}} &raquo; {{item.sourceGroup.name}} &raquo; {{item.sourceSet.name}}</span></a>
						<a :href="'/servers/server/settings/waf/group?serverId=' + item.sourcePolicy.serverId + '&firewallPolicyId=' + item.sourcePolicy.id + '&type=inbound&groupId=' + item.sourceGroup.id + '#set' + item.sourceSet.id" v-if="item.sourcePolicy.serverId > 0"><span class="small "><i class="icon shield"></i> {{item.sourcePolicy.name}} &raquo; {{item.sourceGroup.name}} &raquo; {{item.sourceSet.name}}</span></a>
					</div>
				</td>
				<td>
					<a href="" @click.prevent="viewLogs(item.id)">日志</a> &nbsp;
					<a href="" @click.prevent="updateItem(item.id)">修改</a> &nbsp;
					<a href="" @click.prevent="deleteItem(item.id)">删除</a>
				</td>
			</tr>
        </tbody>
    </table>
</div>`}),Vue.component("ip-item-text",{props:["v-item"],template:`<span>
    <span v-if="vItem.type == 'all'">*</span>
    <span v-if="vItem.type == 'ipv4' || vItem.type.length == 0">
        {{vItem.ipFrom}}
        <span v-if="vItem.ipTo.length > 0">- {{vItem.ipTo}}</span>
    </span>
    <span v-if="vItem.type == 'ipv6'">{{vItem.ipFrom}}</span>
    <span v-if="vItem.eventLevelName != null && vItem.eventLevelName.length > 0">&nbsp; 级别：{{vItem.eventLevelName}}</span>
</span>`}),Vue.component("ip-box",{props:["v-ip"],methods:{popup:function(){let e=this.vIp;var t;null!=e&&0!=e.length||(t=this.$refs.container,null==(e=t.innerText)&&(e=t.textContent)),teaweb.popup("/servers/ipbox?ip="+e,{width:"50em",height:"30em"})}},template:'<span @click.prevent="popup()" ref="container"><slot></slot></span>'}),Vue.component("api-node-selector",{props:[],data:function(){return{}},template:`<div>
	暂未实现
</div>`}),Vue.component("api-node-addresses-box",{props:["v-addrs","v-name"],data:function(){let e=this.vAddrs;return{addrs:e=null==e?[]:e}},methods:{addAddr:function(){let t=this;teaweb.popup("/api/node/createAddrPopup",{height:"16em",callback:function(e){t.addrs.push(e.data.addr)}})},updateAddr:function(t,e){let i=this;window.UPDATING_ADDR=e,teaweb.popup("/api/node/updateAddrPopup?addressId=",{callback:function(e){Vue.set(i.addrs,t,e.data.addr)}})},removeAddr:function(e){this.addrs.$remove(e)}},template:`<div>
	<input type="hidden" :name="vName" :value="JSON.stringify(addrs)"/>
	<div v-if="addrs.length > 0">
		<div>
			<div v-for="(addr, index) in addrs" class="ui label small">
				{{addr.protocol}}://{{addr.host.quoteIP()}}:{{addr.portRange}}</span>
				<a href="" title="修改" @click.prevent="updateAddr(index, addr)"><i class="icon pencil small"></i></a>
				<a href="" title="删除" @click.prevent="removeAddr(index)"><i class="icon remove"></i></a>
			</div>
		</div>
		<div class="ui divider"></div>
	</div>
	<div>
		<button class="ui button small" type="button" @click.prevent="addAddr()">+</button>
	</div>
</div>`}),Vue.component("page-box",{data:function(){return{page:""}},created:function(){let e=this;setTimeout(function(){e.page=Tea.Vue.page})},template:`<div>
	<div class="page" v-html="page"></div>
</div>`}),Vue.component("network-addresses-box",{props:["v-server-type","v-addresses","v-protocol","v-name","v-from","v-support-range"],data:function(){let e=this.vAddresses,t=(null==e&&(e=[]),this.vProtocol),i=(null==t&&(t=""),this.vName),s=(null==i&&(i="addresses"),this.vFrom);return null==s&&(s=""),{addresses:e,protocol:t,name:i,from:s}},watch:{vServerType:function(){this.addresses=[]},vAddresses:function(){null!=this.vAddresses&&(this.addresses=this.vAddresses)}},methods:{addAddr:function(){let t=this;window.UPDATING_ADDR=null,teaweb.popup("/servers/addPortPopup?serverType="+this.vServerType+"&protocol="+this.protocol+"&from="+this.from+"&supportRange="+(this.supportRange()?1:0),{height:"18em",callback:function(e){var i=e.data.address;null!=t.addresses.$find(function(e,t){return i.host==t.host&&i.portRange==t.portRange&&i.protocol==t.protocol})?teaweb.warn("要添加的网络地址已经存在"):(t.addresses.push(i),["https","https4","https6"].$contains(i.protocol)?this.tlsProtocolName="HTTPS":["tls","tls4","tls6"].$contains(i.protocol)&&(this.tlsProtocolName="TLS"),t.$emit("change",t.addresses))}})},removeAddr:function(e){this.addresses.$remove(e),this.$emit("change",this.addresses)},updateAddr:function(t,e){let i=this;window.UPDATING_ADDR=e,teaweb.popup("/servers/addPortPopup?serverType="+this.vServerType+"&protocol="+this.protocol+"&from="+this.from+"&supportRange="+(this.supportRange()?1:0),{height:"18em",callback:function(e){e=e.data.address;Vue.set(i.addresses,t,e),["https","https4","https6"].$contains(e.protocol)?this.tlsProtocolName="HTTPS":["tls","tls4","tls6"].$contains(e.protocol)&&(this.tlsProtocolName="TLS"),i.$emit("change",i.addresses)}}),this.$emit("change",this.addresses)},supportRange:function(){return this.vSupportRange||"tcpProxy"==this.vServerType||"udpProxy"==this.vServerType}},template:`<div>
	<input type="hidden" :name="name" :value="JSON.stringify(addresses)"/>
	<div v-if="addresses.length > 0">
		<div class="ui label small basic" v-for="(addr, index) in addresses">
			{{addr.protocol}}://<span v-if="addr.host.length > 0">{{addr.host.quoteIP()}}</span><span v-if="addr.host.length == 0">*</span>:<span v-if="addr.portRange.indexOf('-')<0">{{addr.portRange}}</span><span v-else style="font-style: italic">{{addr.portRange}}</span>
			<a href="" @click.prevent="updateAddr(index, addr)" title="修改"><i class="icon pencil small"></i></a>
			<a href="" @click.prevent="removeAddr(index)" title="删除"><i class="icon remove"></i></a> </div>
		<div class="ui divider"></div>
	</div>
	<a href="" @click.prevent="addAddr()">[添加端口绑定]</a>
</div>`}),Vue.component("submit-btn",{template:'<button class="ui button primary" type="submit"><slot>保存</slot></button>'}),Vue.component("more-items-angle",{props:["v-data-url","v-url"],data:function(){return{visible:!1}},methods:{show:function(){this.visible=!this.visible,this.visible?this.showBox():this.hideBox()},showBox:function(){let a=this;this.visible=!0,Tea.action(this.vDataUrl).params({url:this.vUrl}).post().success(function(e){let t=e.data.groups;var e=a.$el.offsetLeft+120,i=a.$el.offsetTop+70;function s(e){"I"==e.target.tagName||a.isInBox(n,e.target)||(document.removeEventListener("click",s),a.hideBox())}let n=document.createElement("div"),o=(n.setAttribute("id","more-items-box"),n.style.cssText="z-index: 100; position: absolute; left: "+e+"px; top: "+i+"px; max-height: 30em; overflow: auto; border-bottom: 1px solid rgba(34,36,38,.15)",document.body.append(n),'<ul class="ui labeled menu vertical borderless" style="padding: 0">');t.forEach(function(e){o+='<div class="item header">'+teaweb.encodeHTML(e.name)+"</div>",e.items.forEach(function(e){o+='<a href="'+e.url+'" class="item '+(e.isActive?"active":"")+'" style="font-size: 0.9em;">'+teaweb.encodeHTML(e.name)+'<i class="icon right angle"></i></a>'})}),o+="</ul>",n.innerHTML=o;document.addEventListener("click",s)})},hideBox:function(){let e=document.getElementById("more-items-box");null!=e&&e.parentNode.removeChild(e),this.visible=!1},isInBox:function(e,t){for(;;){if(null==t)break;if(t.parentNode==e)return!0;t=t.parentNode}return!1}},template:'<a href="" class="item" @click.prevent="show"><i class="icon angle" :class="{down: !visible, up: visible}"></i></a>'}),Vue.component("menu-item",{props:["href","active","code"],data:function(){let e=this.active;var t;void 0===e&&(t="",null!=(t=void 0!==window.TEA.ACTION.data.firstMenuItem?window.TEA.ACTION.data.firstMenuItem:t)&&0<t.length&&null!=this.code&&0<this.code.length&&(e=0<t.indexOf(",")?t.split(",").$contains(this.code):t==this.code));let i=null==this.href?"":this.href;return"string"==typeof i&&0<i.length&&i.startsWith(".")&&(t=i.indexOf("?"),i=0<=t?Tea.url(i.substring(0,t))+i.substring(t):Tea.url(i)),{vHref:i,vActive:e}},methods:{click:function(e){this.$emit("click",e)}},template:'\t\t<a :href="vHref" class="item" :class="{active:vActive}" @click="click"><slot></slot></a> \t\t'}),Vue.component("link-icon",{props:["href","title","target"],data:function(){return{vTitle:null==this.title?"打开链接":this.title}},template:'<span><slot></slot>&nbsp;<a :href="href" :title="vTitle" class="link grey" :target="target"><i class="icon linkify small"></i></a></span>'}),Vue.component("link-red",{props:["href","title"],data:function(){let e=this.href;return{vHref:e=null==e?"":e}},methods:{clickPrevent:function(){emitClick(this,arguments)}},template:'<a :href="vHref" :title="title" style="border-bottom: 1px #db2828 dashed" @click.prevent="clickPrevent"><span class="red"><slot></slot></span></a>'}),Vue.component("link-popup",{props:["title"],methods:{clickPrevent:function(){emitClick(this,arguments)}},template:'<a href="" :title="title" @click.prevent="clickPrevent"><slot></slot></a>'}),Vue.component("popup-icon",{props:["title","href","height"],methods:{clickPrevent:function(){null!=this.href&&0<this.href.length&&teaweb.popup(this.href,{height:this.height})}},template:'<span><slot></slot>&nbsp;<a href="" :title="title" @click.prevent="clickPrevent"><i class="icon expand small"></i></a></span>'}),Vue.component("tip-icon",{props:["content"],methods:{showTip:function(){teaweb.popupTip(this.content)}},template:'<a href="" title="查看帮助" @click.prevent="showTip"><i class="icon question circle grey"></i></a>'}),Vue.component("countries-selector",{props:["v-countries"],data:function(){let e=this.vCountries;var t=(e=null==e?[]:e).$map(function(e,t){return t.id});return{countries:e,countryIds:t}},methods:{add:function(){let e=this.countryIds.map(function(e){return e.toString()}),t=this;teaweb.popup("/ui/selectCountriesPopup?countryIds="+e.join(","),{width:"48em",height:"23em",callback:function(e){t.countries=e.data.countries,t.change()}})},remove:function(e){this.countries.$remove(e),this.change()},change:function(){this.countryIds=this.countries.$map(function(e,t){return t.id})}},template:`<div>
	<input type="hidden" name="countryIdsJSON" :value="JSON.stringify(countryIds)"/>
	<div v-if="countries.length > 0" style="margin-bottom: 0.5em">
		<div v-for="(country, index) in countries" class="ui label tiny basic">{{country.name}} <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove"></i></a></div>
		<div class="ui divider"></div>
	</div>
	<div>
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`}),Vue.component("raquo-item",{template:'<span class="item disabled" style="padding: 0">&raquo;</span>'}),Vue.component("more-options-tbody",{data:function(){return{isVisible:!1}},methods:{show:function(){this.isVisible=!this.isVisible,this.$emit("change",this.isVisible)}},template:`<tbody>
	<tr>
		<td colspan="2"><a href="" @click.prevent="show()"><span v-if="!isVisible">更多选项</span><span v-if="isVisible">收起选项</span><i class="icon angle" :class="{down:!isVisible, up:isVisible}"></i></a></td>
	</tr>
</tbody>`}),Vue.component("download-link",{props:["v-element","v-file","v-value"],created:function(){let e=this;setTimeout(function(){e.url=e.composeURL()},1e3)},data:function(){let e=this.vFile;return{file:e=null!=e&&0!=e.length?e:"unknown-file",url:this.composeURL()}},methods:{composeURL:function(){let e="";if(null!=this.vValue)e=this.vValue;else{var t=document.getElementById(this.vElement);if(null==t)return void teaweb.warn("找不到要下载的内容");null==(e=t.innerText)&&(e=t.textContent)}return Tea.url("/ui/download",{file:this.file,text:e})}},template:'<a :href="url" target="_blank" style="font-weight: normal"><slot></slot></a>'}),Vue.component("values-box",{props:["values","size","maxlength","name","placeholder"],data:function(){let e=this.values;return{vValues:e=null==e?[]:e,isUpdating:!1,isAdding:!1,index:0,value:"",isEditing:!1}},methods:{create:function(){this.isAdding=!0;var e=this;setTimeout(function(){e.$refs.value.focus()},200)},update:function(e){this.cancel(),this.isUpdating=!0,this.index=e,this.value=this.vValues[e];var t=this;setTimeout(function(){t.$refs.value.focus()},200)},confirm:function(){0!=this.value.length&&(this.isUpdating?Vue.set(this.vValues,this.index,this.value):this.vValues.push(this.value),this.cancel(),this.$emit("change",this.vValues))},remove:function(e){this.vValues.$remove(e),this.$emit("change",this.vValues)},cancel:function(){this.isUpdating=!1,this.isAdding=!1,this.value=""},updateAll:function(e){this.vValeus=e},addValue:function(e){this.vValues.push(e)},startEditing:function(){this.isEditing=!this.isEditing}},template:`<div>
	<div v-show="!isEditing && vValues.length > 0">
		<div class="ui label tiny basic" v-for="(value, index) in vValues" style="margin-top:0.4em;margin-bottom:0.4em">{{value}}</div>
		<a href="" @click.prevent="startEditing" style="font-size: 0.8em; margin-left: 0.2em">[修改]</a>
	</div>
	<div v-show="isEditing || vValues.length == 0">
		<div style="margin-bottom: 1em" v-if="vValues.length > 0">
			<div class="ui label tiny basic" v-for="(value, index) in vValues" style="margin-top:0.4em;margin-bottom:0.4em">{{value}}
				<input type="hidden" :name="name" :value="value"/>
				&nbsp; <a href="" @click.prevent="update(index)" title="修改"><i class="icon pencil small" ></i></a> 
				<a href="" @click.prevent="remove(index)" title="删除"><i class="icon remove"></i></a> 
			</div> 
			<div class="ui divider"></div>
		</div> 
		<!-- 添加|修改 -->
		<div v-if="isAdding || isUpdating">
			<div class="ui fields inline">
				<div class="ui field">
					<input type="text" :size="size" :maxlength="maxlength" :placeholder="placeholder" v-model="value" ref="value" @keyup.enter="confirm()" @keypress.enter.prevent="1"/>
				</div> 
				<div class="ui field">
					<button class="ui button small" type="button" @click.prevent="confirm()">确定</button> 
				</div>
				<div class="ui field">
					<a href="" @click.prevent="cancel()" title="取消"><i class="icon remove small"></i></a> 
				</div> 
			</div> 
		</div> 
		<div v-if="!isAdding && !isUpdating">
			<button class="ui button tiny" type="button" @click.prevent="create()">+</button> 
		</div>
	</div>	
</div>`}),Vue.component("datetime-input",{props:["v-name","v-timestamp"],mounted:function(){let t=this;teaweb.datepicker(this.$refs.dayInput,function(e){t.day=e,t.hour="23",t.minute="59",t.second="59",t.change()})},data:function(){let t=this.vTimestamp,i=(null!=t?(t=parseInt(t),isNaN(t)&&(t=0)):t=0,""),s="",n="",o="";if(0<t){let e=new Date;e.setTime(1e3*t);var a=e.getFullYear().toString(),l=this.leadingZero((e.getMonth()+1).toString(),2);i=a+"-"+l+"-"+this.leadingZero(e.getDate().toString(),2),s=this.leadingZero(e.getHours().toString(),2),n=this.leadingZero(e.getMinutes().toString(),2),o=this.leadingZero(e.getSeconds().toString(),2)}return{timestamp:t,day:i,hour:s,minute:n,second:o,hasDayError:!1,hasHourError:!1,hasMinuteError:!1,hasSecondError:!1}},methods:{change:function(){let e=new Date;var t,i;/^\d{4}-\d{1,2}-\d{1,2}$/.test(this.day)?(i=this.day.split("-"),t=parseInt(i[0]),e.setFullYear(t),(t=parseInt(i[1]))<1||12<t?this.hasDayError=!0:(e.setMonth(t-1),(t=parseInt(i[2]))<1||32<t?this.hasDayError=!0:(e.setDate(t),this.hasDayError=!1,/^\d+$/.test(this.hour)?(i=parseInt(this.hour),isNaN(i)||i<0||24<=i?this.hasHourError=!0:(this.hasHourError=!1,e.setHours(i),/^\d+$/.test(this.minute)?(t=parseInt(this.minute),isNaN(t)||t<0||60<=t?this.hasMinuteError=!0:(this.hasMinuteError=!1,e.setMinutes(t),/^\d+$/.test(this.second)?(i=parseInt(this.second),isNaN(i)||i<0||60<=i?this.hasSecondError=!0:(this.hasSecondError=!1,e.setSeconds(i),this.timestamp=Math.floor(e.getTime()/1e3))):this.hasSecondError=!0)):this.hasMinuteError=!0)):this.hasHourError=!0))):this.hasDayError=!0},leadingZero:function(t,i){if(i<=(t=t.toString()).length)return t;for(let e=0;e<i-t.length;e++)t="0"+t;return t},resultTimestamp:function(){return this.timestamp},nextDays:function(e){let t=new Date;t.setTime(t.getTime()+86400*e*1e3),this.day=t.getFullYear()+"-"+this.leadingZero(t.getMonth()+1,2)+"-"+this.leadingZero(t.getDate(),2),this.hour=this.leadingZero(t.getHours(),2),this.minute=this.leadingZero(t.getMinutes(),2),this.second=this.leadingZero(t.getSeconds(),2),this.change()}},template:`<div>
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
	<p class="comment">常用时间：<a href="" @click.prevent="nextDays(1)"> &nbsp;1天&nbsp; </a> <span class="disabled">|</span> <a href="" @click.prevent="nextDays(3)"> &nbsp;3天&nbsp; </a> <span class="disabled">|</span> <a href="" @click.prevent="nextDays(7)"> &nbsp;一周&nbsp; </a> <span class="disabled">|</span> <a href="" @click.prevent="nextDays(30)"> &nbsp;30天&nbsp; </a> </p>
</div>`}),Vue.component("label-on",{props:["v-is-on"],template:'<div><span v-if="vIsOn" class="ui label tiny green basic">已启用</span><span v-if="!vIsOn" class="ui label tiny red basic">已停用</span></div>'}),Vue.component("code-label",{methods:{click:function(e){this.$emit("click",e)}},template:'<span class="ui label basic tiny" style="padding: 3px;margin-left:2px;margin-right:2px" @click.prevent="click"><slot></slot></span>'}),Vue.component("code-label-plain",{template:'<span class="ui label basic tiny" style="padding: 3px;margin-left:2px;margin-right:2px"><slot></slot></span>'}),Vue.component("tiny-label",{template:'<span class="ui label tiny" style="margin-bottom: 0.5em"><slot></slot></span>'}),Vue.component("tiny-basic-label",{template:'<span class="ui label tiny basic" style="margin-bottom: 0.5em"><slot></slot></span>'}),Vue.component("micro-basic-label",{template:'<span class="ui label tiny basic" style="margin-bottom: 0.5em; font-size: 0.7em; padding: 4px"><slot></slot></span>'}),Vue.component("grey-label",{props:["color"],data:function(){let e="grey";return{labelColor:e=null!=this.color&&0<this.color.length?"red":e}},template:'<span class="ui label basic tiny" :class="labelColor" style="margin-top: 0.4em; font-size: 0.7em; border: 1px solid #ddd!important; font-weight: normal;"><slot></slot></span>'}),Vue.component("optional-label",{template:'<em><span class="grey">（可选）</span></em>'}),Vue.component("plus-label",{template:'<span style="color: #B18701;">Plus专属功能。</span>'}),Vue.component("pro-warning-label",{template:'<span><i class="icon warning circle"></i>注意：通常不需要修改；如要修改，请在专家指导下进行。</span>'}),Vue.component("first-menu",{props:[],template:' \t\t<div class="first-menu"> \t\t\t<div class="ui menu text blue small">\t\t\t\t<slot></slot>\t\t\t</div> \t\t\t<div class="ui divider"></div> \t\t</div>'}),Vue.component("more-options-indicator",{data:function(){return{visible:!1}},methods:{changeVisible:function(){this.visible=!this.visible,null!=Tea.Vue&&(Tea.Vue.moreOptionsVisible=this.visible),this.$emit("change",this.visible)}},template:'<a href="" style="font-weight: normal" @click.prevent="changeVisible()"><slot><span v-if="!visible">更多选项</span><span v-if="visible">收起选项</span></slot> <i class="icon angle" :class="{down:!visible, up:visible}"></i> </a>'}),Vue.component("page-size-selector",{data:function(){let t=window.location.search,i=10;if(0<t.length){let e=(t=t.substr(1)).split("&");e.forEach(function(t){t=t.split("=");if(2==t.length&&"pageSize"==t[0]){let e=t[1];e.match(/^\d+$/)&&(i=parseInt(e,10),(isNaN(i)||i<1)&&(i=10))}})}return{pageSize:i}},watch:{pageSize:function(){window.ChangePageSize(this.pageSize)}},template:`<select class="ui dropdown" style="height:34px;padding-top:0;padding-bottom:0;margin-left:1em;color:#666" v-model="pageSize">
	<option value="10">[每页]</option><option value="10" selected="selected">10条</option><option value="20">20条</option><option value="30">30条</option><option value="40">40条</option><option value="50">50条</option><option value="60">60条</option><option value="70">70条</option><option value="80">80条</option><option value="90">90条</option><option value="100">100条</option>
</select>`}),Vue.component("second-menu",{template:' \t\t<div class="second-menu"> \t\t\t<div class="ui menu text blue small">\t\t\t\t<slot></slot>\t\t\t</div> \t\t\t<div class="ui divider"></div> \t\t</div>'}),Vue.component("loading-message",{template:`<div class="ui message loading">
        <div class="ui active inline loader small"></div>  &nbsp; <slot></slot>
    </div>`}),Vue.component("more-options-angle",{data:function(){return{isVisible:!1}},methods:{show:function(){this.isVisible=!this.isVisible,this.$emit("change",this.isVisible)}},template:'<a href="" @click.prevent="show()"><span v-if="!isVisible">更多选项</span><span v-if="isVisible">收起选项</span><i class="icon angle" :class="{down:!isVisible, up:isVisible}"></i></a>'}),Vue.component("inner-menu-item",{props:["href","active","code"],data:function(){var e,t=this.active;return void 0===t&&(e="",t=(e=void 0!==window.TEA.ACTION.data.firstMenuItem?window.TEA.ACTION.data.firstMenuItem:e)==this.code),{vHref:null==this.href?"":this.href,vActive:t}},template:'\t\t<a :href="vHref" class="item right" style="color:#4183c4" :class="{active:vActive}">[<slot></slot>]</a> \t\t'}),Vue.component("health-check-config-box",{props:["v-health-check-config"],data:function(){let t=this.vHealthCheckConfig,i="http",s="",n="/",o="";if(null==t){t={isOn:!1,url:"",interval:{count:60,unit:"second"},statusCodes:[200],timeout:{count:10,unit:"second"},countTries:3,tryDelay:{count:100,unit:"ms"},autoDown:!0,countUp:1,countDown:3,userAgent:"",onlyBasicRequest:!0,accessLogIsOn:!0};let e=this;setTimeout(function(){e.changeURL()},500)}else{try{let e=new URL(t.url);i=e.protocol.substring(0,e.protocol.length-1);var a=(o="%24%7Bhost%7D"==(o=e.host)?"${host}":o).indexOf(":");0<a&&(o=o.substring(0,a)),s=e.port,n=e.pathname,0<e.search.length&&(n+=e.search)}catch(e){}null==t.statusCodes&&(t.statusCodes=[200]),null==t.interval&&(t.interval={count:60,unit:"second"}),null==t.timeout&&(t.timeout={count:10,unit:"second"}),null==t.tryDelay&&(t.tryDelay={count:100,unit:"ms"}),(null==t.countUp||t.countUp<1)&&(t.countUp=1),(null==t.countDown||t.countDown<1)&&(t.countDown=3)}return{healthCheck:t,advancedVisible:!1,urlProtocol:i,urlHost:o,urlPort:s,urlRequestURI:n,urlIsEditing:0==t.url.length}},watch:{urlRequestURI:function(){0<this.urlRequestURI.length&&"/"!=this.urlRequestURI[0]&&(this.urlRequestURI="/"+this.urlRequestURI),this.changeURL()},urlPort:function(e){let t=parseInt(e);isNaN(t)?this.urlPort="":this.urlPort=t.toString(),this.changeURL()},urlProtocol:function(){this.changeURL()},urlHost:function(){this.changeURL()},"healthCheck.countTries":function(e){e=parseInt(e);isNaN(e)?this.healthCheck.countTries=0:this.healthCheck.countTries=e},"healthCheck.countUp":function(e){e=parseInt(e);isNaN(e)?this.healthCheck.countUp=0:this.healthCheck.countUp=e},"healthCheck.countDown":function(e){e=parseInt(e);isNaN(e)?this.healthCheck.countDown=0:this.healthCheck.countDown=e}},methods:{showAdvanced:function(){this.advancedVisible=!this.advancedVisible},changeURL:function(){let e=this.urlHost;0==e.length&&(e="${host}"),this.healthCheck.url=this.urlProtocol+"://"+e+(0<this.urlPort.length?":"+this.urlPort:"")+this.urlRequestURI},changeStatus:function(e){this.healthCheck.statusCodes=e.$map(function(e,t){t=parseInt(t);return isNaN(t)?0:t})},editURL:function(){this.urlIsEditing=!this.urlIsEditing}},template:`<div>
<input type="hidden" name="healthCheckJSON" :value="JSON.stringify(healthCheck)"/>
<table class="ui table definition selectable">
	<tbody>
		<tr>
			<td class="title">启用</td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" value="1" v-model="healthCheck.isOn"/>
					<label></label>
				</div>
			</td>
		</tr>
	</tbody>
	<tbody v-show="healthCheck.isOn">
		<tr>
			<td>检测URL *</td>
			<td>
				<div v-if="healthCheck.url.length > 0" style="margin-bottom: 1em"><code-label>{{healthCheck.url}}</code-label> &nbsp; <a href="" @click.prevent="editURL"><span class="small">修改 <i class="icon angle" :class="{down: !urlIsEditing, up: urlIsEditing}"></i></span></a> </div>
				<div v-show="urlIsEditing">
					<table class="ui table">
						 <tr>
							<td class="title">协议</td> 
							<td>
								<select class="ui dropdown auto-width" v-model="urlProtocol">
								<option value="http">http://</option>
								<option value="https">https://</option>
								</select>
							</td>
						</tr>
						<tr>
							<td>域名</td>
							<td>
								<input type="text" v-model="urlHost"/>
								<p class="comment">已经绑定到此集群的一个域名；如果为空则使用节点IP作为域名。</p>
							</td>
						</tr>
						<tr>
							<td>端口</td>
							<td>
								<input type="text" maxlength="5" style="width:5.4em" placeholder="端口" v-model="urlPort"/>
								<p class="comment">域名或者IP的端口，可选项，默认为80/443。</p>
							</td>
						</tr>
						<tr>
							<td>RequestURI</td>
							<td><input type="text" v-model="urlRequestURI" placeholder="/" style="width:20em"/>
								<p class="comment">请求的路径，可以带参数，可选项。</p>
							</td>
						</tr>
					</table>
					<div class="ui divider"></div>
					<p class="comment" v-if="healthCheck.url.length > 0">拼接后的检测URL：<code-label>{{healthCheck.url}}</code-label>，其中\${host}指的是域名。</p>
				</div>
			</td>
		</tr>
		<tr>
			<td>检测时间间隔</td>
			<td>
				<time-duration-box :v-value="healthCheck.interval"></time-duration-box>
				<p class="comment">两次检查之间的间隔。</p>
			</td>
		</tr>
		<tr>
			<td>是否自动下线</td>
			<td>
				<div class="ui checkbox">
					<input type="checkbox" value="1" v-model="healthCheck.autoDown"/>
					<label></label>
				</div>
				<p class="comment">选中后系统会根据健康检查的结果自动标记节点的上线/下线状态，并可能自动同步DNS设置。</p>
			</td>
		</tr>
		<tr v-show="healthCheck.autoDown">
			<td>连续上线次数</td>
			<td>
				<input type="text" v-model="healthCheck.countUp" style="width:5em" maxlength="6"/>
				<p class="comment">连续{{healthCheck.countUp}}次检查成功后自动恢复上线。</p>
			</td>
		</tr>
		<tr v-show="healthCheck.autoDown">
			<td>连续下线次数</td>
			<td>
				<input type="text" v-model="healthCheck.countDown" style="width:5em" maxlength="6"/>
				<p class="comment">连续{{healthCheck.countDown}}次检查失败后自动下线。</p>
			</td>
		</tr>
	</tbody>
	<tbody v-show="healthCheck.isOn">
		<tr>
			<td colspan="2"><more-options-angle @change="showAdvanced"></more-options-angle></td>
		</tr>
	</tbody>
	<tbody v-show="advancedVisible && healthCheck.isOn">
		<tr>
			<td>允许的状态码</td>
			<td>
				<values-box :values="healthCheck.statusCodes" maxlength="3" @change="changeStatus"></values-box>
			</td>
		</tr>
		<tr>
			<td>超时时间</td>
			<td>
				<time-duration-box :v-value="healthCheck.timeout"></time-duration-box>
			</td>	
		</tr>
		<tr>
			<td>连续尝试次数</td>
			<td>
				<input type="text" v-model="healthCheck.countTries" style="width: 5em" maxlength="2"/>
			</td>
		</tr>
		<tr>
			<td>每次尝试间隔</td>
			<td>
				<time-duration-box :v-value="healthCheck.tryDelay"></time-duration-box>
			</td>
		</tr>
		<tr>
			<td>终端信息<em>（User-Agent）</em></td>
			<td>
				<input type="text" v-model="healthCheck.userAgent" maxlength="200"/>
				<p class="comment">发送到服务器的User-Agent值，不填写表示使用默认值。</p>
			</td>
		</tr>
		<tr>
			<td>只基础请求</td>
			<td>
				<checkbox v-model="healthCheck.onlyBasicRequest"></checkbox>
				<p class="comment">只做基础的请求，不处理反向代理（不检查源站）、WAF等。</p>
			</td>
		</tr>
		<tr>
			<td>记录访问日志</td>
			<td>
				<checkbox v-model="healthCheck.accessLogIsOn"></checkbox>
				<p class="comment">是否记录健康检查的访问日志。</p>
			</td>
		</tr>
	</tbody>
</table>
<div class="margin"></div>
</div>`}),Vue.component("request-variables-describer",{data:function(){return{vars:[]}},methods:{update:function(e){this.vars=[];let i=this;e.replace(/\${.+?}/g,function(e){var t=i.findVar(e);if(null==t)return e;i.vars.push(t)})},findVar:function(t){let i=null;return window.REQUEST_VARIABLES.forEach(function(e){e.code==t&&(i=e)}),i}},template:`<span>
	<span v-for="(v, index) in vars"><code-label :title="v.description">{{v.code}}</code-label> - {{v.name}}<span v-if="index < vars.length-1">；</span></span>
</span>`}),Vue.component("combo-box",{props:["name","title","placeholder","size","v-items","v-value"],data:function(){let e=this.vItems,i=((e=null!=e&&e instanceof Array?e:[]).forEach(function(e){null==e.value&&(e.value=e.id)}),null);if(null!=this.vValue){let t=this;e.forEach(function(e){e.value==t.vValue&&(i=e)})}return{allItems:e,items:e.$copy(),selectedItem:i,keyword:"",visible:!1,hideTimer:null,hoverIndex:0}},methods:{reset:function(){this.selectedItem=null,this.change(),this.hoverIndex=0;let e=this;setTimeout(function(){e.$refs.searchBox&&e.$refs.searchBox.focus()})},changeKeyword:function(){this.hoverIndex=0;let t=this.keyword;0==t.length?this.items=this.allItems.$copy():this.items=this.allItems.$copy().filter(function(e){return teaweb.match(e.name,t)})},selectItem:function(e){this.selectedItem=e,this.change(),this.hoverIndex=0,this.keyword="",this.changeKeyword()},confirm:function(){this.items.length>this.hoverIndex&&this.selectItem(this.items[this.hoverIndex])},show:function(){this.visible=!0},hide:function(){let e=this;this.hideTimer=setTimeout(function(){e.visible=!1},500)},downItem:function(){this.hoverIndex++,this.hoverIndex>this.items.length-1&&(this.hoverIndex=0),this.focusItem()},upItem:function(){this.hoverIndex--,this.hoverIndex<0&&(this.hoverIndex=0),this.focusItem()},focusItem:function(){if(this.hoverIndex<this.items.length){this.$refs.itemRef[this.hoverIndex].focus();let e=this;setTimeout(function(){e.$refs.searchBox.focus(),null!=e.hideTimer&&(clearTimeout(e.hideTimer),e.hideTimer=null)})}},change:function(){this.$emit("change",this.selectedItem);let e=this;setTimeout(function(){null!=e.$refs.selectedLabel&&e.$refs.selectedLabel.focus()})},submitForm:function(e){if("A"==e.target.tagName){let e=this.$refs.selectedLabel.parentNode;for(;;){if(null==(e=e.parentNode)||"BODY"==e.tagName)return;if("FORM"==e.tagName){e.submit();break}}}}},template:`<div style="display: inline; z-index: 10; background: white">
	<!-- 搜索框 -->
	<div v-if="selectedItem == null">
		<input type="text" v-model="keyword" :placeholder="placeholder" :size="size" style="width: 11em" @input="changeKeyword" @focus="show" @blur="hide" @keyup.enter="confirm()" @keypress.enter.prevent="1" ref="searchBox" @keyup.down="downItem" @keyup.up="upItem"/>
	</div>
	
	<!-- 当前选中 -->
	<div v-if="selectedItem != null">
		<input type="hidden" :name="name" :value="selectedItem.value"/>
		<a href="" class="ui label basic" style="line-height: 1.4; font-weight: normal; font-size: 1em" ref="selectedLabel" @click.prevent="submitForm"><span>{{title}}：{{selectedItem.name}}</span>
			<span title="清除" @click.prevent="reset"><i class="icon remove small"></i></span>
		</a>
	</div>
	
	<!-- 菜单 -->
	<div v-if="selectedItem == null && items.length > 0 && visible">
		<div class="ui menu vertical small narrow-scrollbar" style="width: 11em; max-height: 17em; overflow-y: auto; position: absolute; border: rgba(129, 177, 210, 0.81) 1px solid; border-top: 0; z-index: 100">
			<a href="" v-for="(item, index) in items" ref="itemRef" class="item" :class="{active: index == hoverIndex, blue: index == hoverIndex}" @click.prevent="selectItem(item)" style="line-height: 1.4">{{item.name}}</a>
		</div>
	</div>
</div>`}),Vue.component("time-duration-box",{props:["v-name","v-value","v-count","v-unit"],mounted:function(){this.change()},data:function(){let e=this.vValue;return"number"!=typeof(e=null==e?{count:this.vCount,unit:this.vUnit}:e).count&&(e.count=-1),{duration:e,countString:0<=e.count?e.count.toString():""}},watch:{countString:function(e){var e=e.trim();0==e.length?this.duration.count=-1:(e=parseInt(e),isNaN(e)||(this.duration.count=e),this.change())}},methods:{change:function(){this.$emit("change",this.duration)}},template:`<div class="ui fields inline" style="padding-bottom: 0; margin-bottom: 0">
	<input type="hidden" :name="vName" :value="JSON.stringify(duration)"/>
	<div class="ui field">
		<input type="text" v-model="countString" maxlength="11" size="11" @keypress.enter.prevent="1"/>
	</div>
	<div class="ui field">
		<select class="ui dropdown" v-model="duration.unit" @change="change">
			<option value="ms">毫秒</option>
			<option value="second">秒</option>
			<option value="minute">分钟</option>
			<option value="hour">小时</option>
			<option value="day">天</option>
			<option value="week">周</option>
		</select>
	</div>
</div>`}),Vue.component("not-found-box",{props:["message"],template:`<div style="text-align: center; margin-top: 5em;">
	<div style="font-size: 2em; margin-bottom: 1em"><i class="icon exclamation triangle large grey"></i></div>
	<p class="comment">{{message}}<slot></slot></p>
</div>`}),Vue.component("warning-message",{template:'<div class="ui icon message warning"><i class="icon warning circle"></i><div class="content"><slot></slot></div></div>'});let checkboxId=0,radioId=(Vue.component("checkbox",{props:["name","value","v-value","id","checked"],data:function(){checkboxId++;let e=this.id,t=(null==e&&(e="checkbox"+checkboxId),this.vValue),i=(null==t&&(t="1"),this.value);return null==i&&"checked"==this.checked&&(i=t),{elementId:e,elementValue:t,newValue:i}},methods:{change:function(){this.$emit("input",this.newValue)},check:function(){this.newValue=this.elementValue},uncheck:function(){this.newValue=""},isChecked:function(){return this.newValue==this.elementValue}},watch:{value:function(e){"boolean"==typeof e&&(this.newValue=e)}},template:`<div class="ui checkbox">
	<input type="checkbox" :name="name" :value="elementValue" :id="elementId" @change="change" v-model="newValue"/>
	<label :for="elementId" style="font-size: 0.85em!important;"><slot></slot></label>
</div>`}),Vue.component("network-addresses-view",{props:["v-addresses"],template:`<div>
	<div class="ui label tiny basic" v-if="vAddresses != null" v-for="addr in vAddresses">
		{{addr.protocol}}://<span v-if="addr.host.length > 0">{{addr.host.quoteIP()}}</span><span v-else>*</span>:{{addr.portRange}}
	</div>
</div>`}),Vue.component("size-capacity-view",{props:["v-default-text","v-value"],template:`<div>
	<span v-if="vValue != null && vValue.count > 0">{{vValue.count}}{{vValue.unit.toUpperCase()}}</span>
	<span v-else>{{vDefaultText}}</span>
</div>`}),Vue.component("tip-message-box",{props:["code"],mounted:function(){let t=this;Tea.action("/ui/showTip").params({code:this.code}).success(function(e){t.visible=e.data.visible}).post()},data:function(){return{visible:!1}},methods:{close:function(){this.visible=!1,Tea.action("/ui/hideTip").params({code:this.code}).post()}},template:`<div class="ui icon message" v-if="visible">
	<i class="icon info circle"></i>
	<i class="close icon" title="取消" @click.prevent="close" style="margin-top: 1em"></i>
	<div class="content">
		<slot></slot>
	</div>
</div>`}),Vue.component("digit-input",{props:["value","maxlength","size","min","max","required","placeholder"],mounted:function(){let e=this;setTimeout(function(){e.check()})},data:function(){let e=this.maxlength,t=(null==e&&(e=20),this.size);return null==t&&(t=6),{realValue:this.value,realMaxLength:e,realSize:t,isValid:!0}},watch:{realValue:function(e){this.notifyChange()}},methods:{notifyChange:function(){let e=parseInt(this.realValue.toString(),10);isNaN(e)&&(e=0),this.check(),this.$emit("input",e)},check:function(){var e;null!=this.realValue&&(e=this.realValue.toString(),/^\d+$/.test(e)?(e=parseInt(e,10),isNaN(e)?this.isValid=!1:this.required?this.isValid=(null==this.min||this.min<=e)&&(null==this.max||this.max>=e):this.isValid=0==e||(null==this.min||this.min<=e)&&(null==this.max||this.max>=e)):this.isValid=!1)}},template:'<input type="text" v-model="realValue" :maxlength="realMaxLength" :size="realSize" :class="{error: !this.isValid}" :placeholder="placeholder"/>'}),Vue.component("keyword",{props:["v-word"],data:function(){let e=this.vWord;e=null==e?"":(e=(e=(e=(e=(e=(e=(e=(e=(e=e.replace(/\)/g,"\\)")).replace(/\(/g,"\\(")).replace(/\+/g,"\\+")).replace(/\^/g,"\\^")).replace(/\$/g,"\\$")).replace(/\?/,"\\?")).replace(/\*/,"\\*")).replace(/\[/,"\\[")).replace(/{/,"\\{")).replace(/\./,"\\.");let t=this.$slots.default[0].text;if(0<e.length){let i=this,s=[],n=0;t=t.replaceAll(new RegExp("("+e+")","ig"),function(e){n++;var e='<span style="border: 1px #ccc dashed; color: #ef4d58">'+i.encodeHTML(e)+"</span>",t="$TMP__KEY__"+n.toString()+"$";return s.push([t,e]),t}),t=this.encodeHTML(t),s.forEach(function(e){t=t.replace(e[0],e[1])})}else t=this.encodeHTML(t);return{word:e,text:t}},methods:{encodeHTML:function(e){return e=(e=(e=(e=e.replace(/&/g,"&amp;")).replace(/</g,"&lt;")).replace(/>/g,"&gt;")).replace(/"/g,"&quot;")}},template:'<span><span style="display: none"><slot></slot></span><span v-html="text"></span></span>'}),Vue.component("node-log-row",{props:["v-log","v-keyword"],data:function(){return{log:this.vLog,keyword:this.vKeyword}},template:`<div>
	<pre class="log-box" style="margin: 0; padding: 0"><span :class="{red:log.level == 'error', orange:log.level == 'warning', green: log.level == 'success'}"><span v-if="!log.isToday">[{{log.createdTime}}]</span><strong v-if="log.isToday">[{{log.createdTime}}]</strong><keyword :v-word="keyword">[{{log.tag}}]{{log.description}}</keyword></span> &nbsp; <span v-if="log.count > 1" class="ui label tiny" :class="{red:log.level == 'error', orange:log.level == 'warning'}">共{{log.count}}条</span> <span v-if="log.server != null && log.server.id > 0"><a :href="'/servers/server?serverId=' + log.server.id" class="ui label tiny basic">{{log.server.name}}</a></span></pre>
</div>`}),Vue.component("provinces-selector",{props:["v-provinces"],data:function(){let e=this.vProvinces;var t=(e=null==e?[]:e).$map(function(e,t){return t.id});return{provinces:e,provinceIds:t}},methods:{add:function(){let e=this.provinceIds.map(function(e){return e.toString()}),t=this;teaweb.popup("/ui/selectProvincesPopup?provinceIds="+e.join(","),{width:"48em",height:"23em",callback:function(e){t.provinces=e.data.provinces,t.change()}})},remove:function(e){this.provinces.$remove(e),this.change()},change:function(){this.provinceIds=this.provinces.$map(function(e,t){return t.id})}},template:`<div>
	<input type="hidden" name="provinceIdsJSON" :value="JSON.stringify(provinceIds)"/>
	<div v-if="provinces.length > 0" style="margin-bottom: 0.5em">
		<div v-for="(province, index) in provinces" class="ui label tiny basic">{{province.name}} <a href="" title="删除" @click.prevent="remove(index)"><i class="icon remove"></i></a></div>
		<div class="ui divider"></div>
	</div>
	<div>
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`}),Vue.component("csrf-token",{created:function(){this.refreshToken()},mounted:function(){let e=this;this.$refs.token.form.addEventListener("submit",function(){e.refreshToken()}),setInterval(function(){e.refreshToken()},6e5)},data:function(){return{token:""}},methods:{refreshToken:function(){let t=this;Tea.action("/csrf/token").get().success(function(e){t.token=e.data.token})}},template:'<input type="hidden" name="csrfToken" :value="token" ref="token"/>'}),Vue.component("labeled-input",{props:["name","size","maxlength","label","value"],template:'<div class="ui input right labeled"> \t<input type="text" :name="name" :size="size" :maxlength="maxlength" :value="value"/>\t<span class="ui label">{{label}}</span></div>'}),0),sourceCodeBoxIndex=(Vue.component("radio",{props:["name","value","v-value","id"],data:function(){radioId++;let e=this.id;return{elementId:e=null==e?"radio"+radioId:e}},methods:{change:function(){this.$emit("input",this.vValue)}},template:`<div class="ui checkbox radio">
	<input type="radio" :name="name" :value="vValue" :id="elementId" @change="change" :checked="(vValue == value)"/>
	<label :for="elementId"><slot></slot></label>
</div>`}),Vue.component("copy-to-clipboard",{props:["v-target"],created:function(){if("undefined"==typeof ClipboardJS){let e=document.createElement("script");e.setAttribute("src","/js/clipboard.min.js"),document.head.appendChild(e)}},methods:{copy:function(){new ClipboardJS("[data-clipboard-target]"),teaweb.successToast("已复制到剪切板")}},template:`<a href="" title="拷贝到剪切板" :data-clipboard-target="'#' + vTarget" @click.prevent="copy"><i class="ui icon copy small"></i></em></a>`}),Vue.component("node-role-name",{props:["v-role"],data:function(){let e="";switch(this.vRole){case"node":e="边缘节点";break;case"monitor":e="监控节点";break;case"api":e="API节点";break;case"user":e="用户平台";break;case"admin":e="管理平台";break;case"database":e="数据库节点";break;case"dns":e="DNS节点";break;case"report":e="区域监控终端"}return{roleName:e}},template:"<span>{{roleName}}</span>"}),0);Vue.component("source-code-box",{props:["name","type","id","read-only","width","height","focus"],mounted:function(){let e=this.readOnly;"boolean"!=typeof e&&(e=!0);var t=document.getElementById("source-code-box-"+this.index),i=document.getElementById(this.valueBoxId);let s="";null!=i.textContent?s=i.textContent:null!=i.innerText&&(s=i.innerText),this.createEditor(t,s,e)},data:function(){var e=sourceCodeBoxIndex++;let t="source-code-box-value-"+sourceCodeBoxIndex;return{index:e,valueBoxId:t=null!=this.id?this.id:t}},methods:{createEditor:function(e,t,i){let s=CodeMirror.fromTextArea(e,{theme:"idea",lineNumbers:!0,value:"",readOnly:i,showCursorWhenSelecting:!0,height:"auto",viewportMargin:1/0,lineWrapping:!0,highlightFormatting:!1,indentUnit:4,indentWithTabs:!0}),n=this,o=(s.on("change",function(){n.change(s.getValue())}),s.setValue(t),this.focus&&s.focus(),this.width),a=this.height;null!=o&&null!=a?(o=parseInt(o),a=parseInt(a),isNaN(o)||isNaN(a)||(o<=0&&(o=e.parentNode.offsetWidth),s.setSize(o,a))):null!=a&&(a=parseInt(a),isNaN(a)||s.setSize("100%",a));i=CodeMirror.findModeByMIME(this.type);null!=i&&(s.setOption("mode",i.mode),CodeMirror.modeURL="/codemirror/mode/%N/%N.js",CodeMirror.autoLoadMode(s,i.mode))},change:function(e){this.$emit("change",e)}},template:`<div class="source-code-box">
	<div style="display: none" :id="valueBoxId"><slot></slot></div>
	<textarea :id="'source-code-box-' + index" :name="name"></textarea>
</div>`}),Vue.component("size-capacity-box",{props:["v-name","v-value","v-count","v-unit","size","maxlength","v-supported-units"],data:function(){let e=this.vValue,t=("number"!=typeof(e=null==e?{count:this.vCount,unit:this.vUnit}:e).count&&(e.count=-1),this.size),i=(null==t&&(t=6),this.maxlength),s=(null==i&&(i=10),this.vSupportedUnits);return null==s&&(s=[]),{capacity:e,countString:0<=e.count?e.count.toString():"",vSize:t,vMaxlength:i,supportedUnits:s}},watch:{countString:function(e){e=e.trim();if(0==e.length)return this.capacity.count=-1,void this.change();e=parseInt(e);isNaN(e)||(this.capacity.count=e),this.change()}},methods:{change:function(){this.$emit("change",this.capacity)}},template:`<div class="ui fields inline">
	<input type="hidden" :name="vName" :value="JSON.stringify(capacity)"/>
	<div class="ui field">
		<input type="text" v-model="countString" :maxlength="vMaxlength" :size="vSize"/>
	</div>
	<div class="ui field">
		<select class="ui dropdown" v-model="capacity.unit" @change="change">
			<option value="byte" v-if="supportedUnits.length == 0 || supportedUnits.$contains('byte')">字节</option>
			<option value="kb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('kb')">KB</option>
			<option value="mb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('mb')">MB</option>
			<option value="gb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('gb')">GB</option>
			<option value="tb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('tb')">TB</option>
			<option value="pb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('pb')">PB</option>
			<option value="eb" v-if="supportedUnits.length == 0 || supportedUnits.$contains('eb')">EB</option>
		</select>
	</div>
</div>`}),Vue.component("inner-menu",{template:`
		<div class="second-menu" style="width:80%;position: absolute;top:-8px;right:1em"> 
			<div class="ui menu text blue small">
				<slot></slot>
			</div> 
		</div>`}),Vue.component("datepicker",{props:["v-name","v-value","v-bottom-left"],mounted:function(){let t=this;teaweb.datepicker(this.$refs.dayInput,function(e){t.day=e,t.change()},!!this.vBottomLeft)},data:function(){let e=this.vName,t=(null==e&&(e="day"),this.vValue);return null==t&&(t=""),{name:e,day:t}},methods:{change:function(){this.$emit("change",this.day)}},template:`<div style="display: inline-block">
	<input type="text" :name="name" v-model="day" placeholder="YYYY-MM-DD" style="width:8.6em" maxlength="10" @input="change" ref="dayInput" autocomplete="off"/>
</div>`}),Vue.component("sort-arrow",{props:["name"],data:function(){let e=window.location.toString(),n="",o=[];if(null!=window.location.search&&0<window.location.search.length){let e=window.location.search.substring(1),t=e.split("&"),s=this;t.forEach(function(e){var t,i=e.indexOf("=");0<i?(t=e.substring(0,i),i=e.substring(i+1),t==s.name?n=i:"page"!=t&&"asc"!=i&&"desc"!=i&&o.push(e)):o.push(e)})}"asc"!=n&&"desc"==n?o.push(this.name+"=asc"):o.push(this.name+"=desc");var t=e.indexOf("?");return e=0<t?e.substring(0,t)+"?"+o.join("&"):e+"?"+o.join("&"),{order:n,url:e}},template:`<a :href="url" title="排序">&nbsp; <i class="ui icon long arrow small" :class="{down: order == 'asc', up: order == 'desc', 'down grey': order == '' || order == null}"></i></a>`}),Vue.component("user-link",{props:["v-user","v-keyword"],data:function(){let e=this.vUser;return{user:e=null==e?{id:0,username:"",fullname:""}:e}},template:`<div style="display: inline-block">
	<span v-if="user.id > 0"><keyword :v-word="vKeyword">{{user.fullname}}</keyword><span class="small grey">（<keyword :v-word="vKeyword">{{user.username}}</keyword>）</span></span>
	<span v-else class="disabled">[已删除]</span>
</div>`}),Vue.component("report-node-groups-selector",{props:["v-group-ids"],mounted:function(){let t=this;Tea.action("/clusters/monitors/groups/options").post().success(function(e){t.groups=e.data.groups.map(function(e){return e.isChecked=t.groupIds.$contains(e.id),e}),t.isLoaded=!0})},data:function(){var e=this.vGroupIds;return{groups:[],groupIds:e=null==e?[]:e,isLoaded:!1,allGroups:0==e.length}},methods:{check:function(e){e.isChecked=!e.isChecked,this.groupIds=[];let t=this;this.groups.forEach(function(e){e.isChecked&&t.groupIds.push(e.id)}),this.change()},change:function(){let t=this,s=[];this.groupIds.forEach(function(i){var e=t.groups.$find(function(e,t){return t.id==i});null!=e&&s.push({id:e.id,name:e.name})}),this.$emit("change",s)}},watch:{allGroups:function(e){e&&(this.groupIds=[],this.groups.forEach(function(e){e.isChecked=!1})),this.change()}},template:`<div>
	<input type="hidden" name="reportNodeGroupIdsJSON" :value="JSON.stringify(groupIds)"/>
	<span class="disabled" v-if="isLoaded && groups.length == 0">还没有分组。</span>
	<div v-if="groups.length > 0">
		<div>
			<div class="ui checkbox">
				<input type="checkbox" v-model="allGroups" id="all-group"/>
				<label for="all-group">全部分组</label>
			</div>
			<div class="ui divider" v-if="!allGroups"></div>
		</div>
		<div v-show="!allGroups">
			<div v-for="group in groups" :key="group.id" style="float: left; width: 7.6em; margin-bottom: 0.5em">
				<div class="ui checkbox">
					<input type="checkbox" v-model="group.isChecked" value="1" :id="'report-node-group-' + group.id" @click.prevent="check(group)"/>
					<label :for="'report-node-group-' + group.id">{{group.name}}</label>
				</div>
			</div>
		</div>
	</div>
</div>`}),Vue.component("finance-user-selector",{mounted:function(){let t=this;Tea.action("/finance/users/options").post().success(function(e){t.users=e.data.users})},props:["v-user-id"],data:function(){let e=this.vUserId;return{users:[],userId:e=null==e?0:e}},watch:{userId:function(e){this.$emit("change",e)}},template:`<div>
	<select class="ui dropdown auto-width" name="userId" v-model="userId">
		<option value="0">[选择用户]</option>
		<option v-for="user in users" :value="user.id">{{user.fullname}} ({{user.username}})</option>
	</select>
</div>`}),Vue.component("node-login-suggest-ports",{data:function(){return{ports:[],availablePorts:[],autoSelected:!1,isLoading:!1}},methods:{reload:function(e){let t=this;this.autoSelected=!1,this.isLoading=!0,Tea.action("/clusters/cluster/suggestLoginPorts").params({host:e}).success(function(e){null!=e.data.availablePorts&&(t.availablePorts=e.data.availablePorts,0<t.availablePorts.length&&(t.autoSelectPort(t.availablePorts[0]),t.autoSelected=!0)),null!=e.data.ports&&(t.ports=e.data.ports,0<t.ports.length&&!t.autoSelected&&(t.autoSelectPort(t.ports[0]),t.autoSelected=!0))}).done(function(){t.isLoading=!1}).post()},selectPort:function(e){this.$emit("select",e)},autoSelectPort:function(e){this.$emit("auto-select",e)}},template:`<span>
	<span v-if="isLoading">正在检查端口...</span>
	<span v-if="availablePorts.length > 0">
		可能端口：<a href="" v-for="port in availablePorts" @click.prevent="selectPort(port)" class="ui label tiny basic blue" style="border: 1px #2185d0 dashed; font-weight: normal">{{port}}</a>
		&nbsp; &nbsp;
	</span>
	<span v-if="ports.length > 0">
		常用端口：<a href="" v-for="port in ports" @click.prevent="selectPort(port)" class="ui label tiny basic blue" style="border: 1px #2185d0 dashed;  font-weight: normal">{{port}}</a>
	</span>
	<span v-if="ports.length == 0">常用端口有22等。</span>
	<span v-if="ports.length > 0" class="grey small">（可以点击要使用的端口）</span>
</span>`}),Vue.component("node-group-selector",{props:["v-cluster-id","v-group"],data:function(){return{selectedGroup:this.vGroup}},methods:{selectGroup:function(){let t=this;teaweb.popup("/clusters/cluster/groups/selectPopup?clusterId="+this.vClusterId,{callback:function(e){t.selectedGroup=e.data.group}})},addGroup:function(){let t=this;teaweb.popup("/clusters/cluster/groups/createPopup?clusterId="+this.vClusterId,{callback:function(e){t.selectedGroup=e.data.group}})},removeGroup:function(){this.selectedGroup=null}},template:`<div>
	<div class="ui label small basic" v-if="selectedGroup != null">
		<input type="hidden" name="groupId" :value="selectedGroup.id"/>
		{{selectedGroup.name}} &nbsp;<a href="" title="删除" @click.prevent="removeGroup()"><i class="icon remove"></i></a>
	</div>
	<div v-if="selectedGroup == null">
		<a href="" @click.prevent="selectGroup()">[选择分组]</a> &nbsp; <a href="" @click.prevent="addGroup()">[添加分组]</a>
	</div>
</div>`}),Vue.component("node-ip-addresses-box",{props:["v-ip-addresses","role"],data:function(){return{ipAddresses:null==this.vIpAddresses?[]:this.vIpAddresses,supportThresholds:"ns"!=this.role}},methods:{addIPAddress:function(){window.UPDATING_NODE_IP_ADDRESS=null;let t=this;teaweb.popup("/nodes/ipAddresses/createPopup?supportThresholds="+(this.supportThresholds?1:0),{callback:function(e){t.ipAddresses.push(e.data.ipAddress)},height:"24em",width:"44em"})},updateIPAddress:function(t,e){window.UPDATING_NODE_IP_ADDRESS=e;let i=this;teaweb.popup("/nodes/ipAddresses/updatePopup?supportThresholds="+(this.supportThresholds?1:0),{callback:function(e){Vue.set(i.ipAddresses,t,e.data.ipAddress)},height:"24em",width:"44em"})},removeIPAddress:function(e){this.ipAddresses.$remove(e)},isIPv6:function(e){return-1<e.indexOf(":")}},template:`<div>
	<input type="hidden" name="ipAddressesJSON" :value="JSON.stringify(ipAddresses)"/>
	<div v-if="ipAddresses.length > 0">
		<div>
			<div v-for="(address, index) in ipAddresses" class="ui label tiny basic">
				<span v-if="isIPv6(address.ip)" class="grey">[IPv6]</span> {{address.ip}}
				<span class="small grey" v-if="address.name.length > 0">（{{address.name}}<span v-if="!address.canAccess">，不可访问</span>）</span>
				<span class="small grey" v-if="address.name.length == 0 && !address.canAccess">（不可访问）</span>
				<span class="small red" v-if="!address.isOn" title="未启用">[off]</span>
				<span class="small red" v-if="!address.isUp" title="已下线">[down]</span>
				<span class="small" v-if="address.thresholds != null && address.thresholds.length > 0">[{{address.thresholds.length}}个阈值]</span>
				&nbsp;
				<a href="" title="修改" @click.prevent="updateIPAddress(index, address)"><i class="icon pencil small"></i></a>
				<a href="" title="删除" @click.prevent="removeIPAddress(index)"><i class="icon remove"></i></a>
			</div>
		</div>
		<div class="ui divider"></div>
	</div>
	<div>
		<button class="ui button small" type="button" @click.prevent="addIPAddress()">+</button>
	</div>
</div>`}),Vue.component("node-ip-address-thresholds-view",{props:["v-thresholds"],data:function(){let e=this.vThresholds;return null==e?e=[]:e.forEach(function(e){null==e.items&&(e.items=[]),null==e.actions&&(e.actions=[])}),{thresholds:e,allItems:window.IP_ADDR_THRESHOLD_ITEMS,allOperators:[{name:"小于等于",code:"lte"},{name:"大于",code:"gt"},{name:"不等于",code:"neq"},{name:"小于",code:"lt"},{name:"大于等于",code:"gte"}],allActions:window.IP_ADDR_THRESHOLD_ACTIONS}},methods:{itemName:function(t){let i="";return this.allItems.forEach(function(e){e.code==t&&(i=e.name)}),i},itemUnitName:function(t){let i="";return this.allItems.forEach(function(e){e.code==t&&(i=e.unit)}),i},itemDurationUnitName:function(e){switch(e){case"minute":return"分钟";case"second":return"秒";case"hour":return"小时";case"day":return"天"}return e},itemOperatorName:function(t){let i="";return this.allOperators.forEach(function(e){e.code==t&&(i=e.name)}),i},actionName:function(t){let i="";return this.allActions.forEach(function(e){e.code==t&&(i=e.name)}),i}},template:`<div>
	<!-- 已有条件 -->
	<div v-if="thresholds.length > 0">
		<div class="ui label basic small" v-for="(threshold, index) in thresholds" style="margin-bottom: 0.8em">
			<span v-for="(item, itemIndex) in threshold.items">
				<span>
					<span v-if="item.item != 'nodeHealthCheck'">
						[{{item.duration}}{{itemDurationUnitName(item.durationUnit)}}]
					</span>	 
					{{itemName(item.item)}}
					
					<span v-if="item.item == 'nodeHealthCheck'">
						<!-- 健康检查 -->
						<span v-if="item.value == 1">成功</span>
						<span v-if="item.value == 0">失败</span>
					</span>
					<span v-else>
						<!-- 连通性 -->
						<span v-if="item.item == 'connectivity' && item.options != null && item.options.groups != null && item.options.groups.length > 0">[<span v-for="(group, groupIndex) in item.options.groups">{{group.name}} <span v-if="groupIndex != item.options.groups.length - 1">&nbsp; </span></span>]</span>
						
						 <span class="grey">[{{itemOperatorName(item.operator)}}]</span> {{item.value}}{{itemUnitName(item.item)}} &nbsp;
					 </span>
				 </span>
				 <span v-if="itemIndex != threshold.items.length - 1" style="font-style: italic">AND &nbsp;</span></span>
				-&gt;
				<span v-for="(action, actionIndex) in threshold.actions">{{actionName(action.action)}}
				<span v-if="action.action == 'switch'">到{{action.options.ips.join(", ")}}</span>
				<span v-if="action.action == 'webHook'" class="small grey">({{action.options.url}})</span>
				 &nbsp;					 
				 <span v-if="actionIndex != threshold.actions.length - 1" style="font-style: italic">AND &nbsp;</span>
			 </span>
		</div>
	</div>
</div>`}),Vue.component("node-ip-address-thresholds-box",{props:["v-thresholds"],data:function(){let e=this.vThresholds;return null==e?e=[]:e.forEach(function(e){null==e.items&&(e.items=[]),null==e.actions&&(e.actions=[])}),{editingIndex:-1,thresholds:e,addingThreshold:{items:[],actions:[]},isAdding:!1,isAddingItem:!1,isAddingAction:!1,itemCode:"nodeAvgRequests",itemReportGroups:[],itemOperator:"lte",itemValue:"",itemDuration:"5",allItems:window.IP_ADDR_THRESHOLD_ITEMS,allOperators:[{name:"小于等于",code:"lte"},{name:"大于",code:"gt"},{name:"不等于",code:"neq"},{name:"小于",code:"lt"},{name:"大于等于",code:"gte"}],allActions:window.IP_ADDR_THRESHOLD_ACTIONS,actionCode:"up",actionBackupIPs:"",actionWebHookURL:""}},methods:{add:function(){this.isAdding=!this.isAdding},cancel:function(){this.isAdding=!1,this.editingIndex=-1,this.addingThreshold={items:[],actions:[]}},confirm:function(){0==this.addingThreshold.items.length?teaweb.warn("需要至少添加一个阈值"):0==this.addingThreshold.actions.length?teaweb.warn("需要至少添加一个动作"):(0<=this.editingIndex?(this.thresholds[this.editingIndex].items=this.addingThreshold.items,this.thresholds[this.editingIndex].actions=this.addingThreshold.actions):this.thresholds.push({items:this.addingThreshold.items,actions:this.addingThreshold.actions}),this.cancel())},remove:function(e){this.cancel(),this.thresholds.$remove(e)},update:function(e){this.editingIndex=e,this.addingThreshold={items:this.thresholds[e].items.$copy(),actions:this.thresholds[e].actions.$copy()},this.isAdding=!0},addItem:function(){this.isAddingItem=!this.isAddingItem;let e=this;setTimeout(function(){e.$refs.itemValue.focus()},100)},cancelItem:function(){this.isAddingItem=!1,this.itemCode="nodeAvgRequests",this.itmeOperator="lte",this.itemValue="",this.itemDuration="5",this.itemReportGroups=[]},confirmItem:function(){if(["nodeHealthCheck"].$contains(this.itemCode)){if(0==this.itemValue.toString().length)return void teaweb.warn("请选择检查结果");let e=parseInt(this.itemValue);return isNaN(e)||e<0?e=0:1<e&&(e=1),this.addingThreshold.items.push({item:this.itemCode,operator:this.itemOperator,value:e,duration:0,durationUnit:"minute",options:{}}),void this.cancelItem()}if(0==this.itemDuration.length){let e=this;void teaweb.warn("请输入统计周期",function(){e.$refs.itemDuration.focus()})}else{var t=parseInt(this.itemDuration);if(isNaN(t)||t<=0)teaweb.warn("请输入正确的统计周期",function(){that.$refs.itemDuration.focus()});else if(0==this.itemValue.length){let e=this;void teaweb.warn("请输入对比值",function(){e.$refs.itemValue.focus()})}else{var i=parseFloat(this.itemValue);if(isNaN(i))teaweb.warn("请输入正确的对比值",function(){that.$refs.itemValue.focus()});else{let e={};if("connectivity"===this.itemCode){if(100<i){let e=this;return void teaweb.warn("连通性对比值不能超过100",function(){e.$refs.itemValue.focus()})}e.groups=this.itemReportGroups}this.addingThreshold.items.push({item:this.itemCode,operator:this.itemOperator,value:i,duration:t,durationUnit:"minute",options:e}),this.cancelItem()}}}},removeItem:function(e){this.cancelItem(),this.addingThreshold.items.$remove(e)},changeReportGroups:function(e){this.itemReportGroups=e},itemName:function(t){let i="";return this.allItems.forEach(function(e){e.code==t&&(i=e.name)}),i},itemUnitName:function(t){let i="";return this.allItems.forEach(function(e){e.code==t&&(i=e.unit)}),i},itemDurationUnitName:function(e){switch(e){case"minute":return"分钟";case"second":return"秒";case"hour":return"小时";case"day":return"天"}return e},itemOperatorName:function(t){let i="";return this.allOperators.forEach(function(e){e.code==t&&(i=e.name)}),i},addAction:function(){this.isAddingAction=!this.isAddingAction},cancelAction:function(){this.isAddingAction=!1,this.actionCode="up",this.actionBackupIPs="",this.actionWebHookURL=""},confirmAction:function(){this.doConfirmAction(!1)},doConfirmAction:function(e,t){let i=!1,s=this;if(this.addingThreshold.actions.forEach(function(e){e.action==s.actionCode&&(i=!0)}),i)teaweb.warn("此动作已经添加过了，无需重复添加");else{switch(null==t&&(t={}),this.actionCode){case"switch":if(e)break;return void Tea.action("/ui/validateIPs").params({ips:this.actionBackupIPs}).success(function(e){0==e.data.ips.length?teaweb.warn("请输入备用IP",function(){s.$refs.actionBackupIPs.focus()}):(t.ips=e.data.ips,s.doConfirmAction(!0,t))}).fail(function(e){teaweb.warn("输入的IP '"+e.data.failIP+"' 格式不正确，请改正后提交",function(){s.$refs.actionBackupIPs.focus()})}).post();case"webHook":if(0==this.actionWebHookURL.length)return void teaweb.warn("请输入WebHook URL",function(){s.$refs.webHookURL.focus()});if(!this.actionWebHookURL.match(/^(http|https):\/\//i))return void teaweb.warn("URL开头必须是http://或者https://",function(){s.$refs.webHookURL.focus()});t.url=this.actionWebHookURL}this.addingThreshold.actions.push({action:this.actionCode,options:t}),this.cancelAction()}},removeAction:function(e){this.cancelAction(),this.addingThreshold.actions.$remove(e)},actionName:function(t){let i="";return this.allActions.forEach(function(e){e.code==t&&(i=e.name)}),i}},template:`<div>
	<input type="hidden" name="thresholdsJSON" :value="JSON.stringify(thresholds)"/>
		
	<!-- 已有条件 -->
	<div v-if="thresholds.length > 0">
		<div class="ui label basic small" v-for="(threshold, index) in thresholds">
			<span v-for="(item, itemIndex) in threshold.items">
				<span v-if="item.item != 'nodeHealthCheck'">
					[{{item.duration}}{{itemDurationUnitName(item.durationUnit)}}]
				</span> 
				{{itemName(item.item)}}
				
				<span v-if="item.item == 'nodeHealthCheck'">
					<!-- 健康检查 -->
					<span v-if="item.value == 1">成功</span>
					<span v-if="item.value == 0">失败</span>
				</span>
				<span v-else>
					<!-- 连通性 -->
					<span v-if="item.item == 'connectivity' && item.options != null && item.options.groups != null && item.options.groups.length > 0">[<span v-for="(group, groupIndex) in item.options.groups">{{group.name}} <span v-if="groupIndex != item.options.groups.length - 1">&nbsp; </span></span>]</span>
				
					<span  class="grey">[{{itemOperatorName(item.operator)}}]</span> &nbsp;{{item.value}}{{itemUnitName(item.item)}} 
			 	</span>
			 	&nbsp;<span v-if="itemIndex != threshold.items.length - 1" style="font-style: italic">AND &nbsp;</span>
			</span>
			-&gt;
			<span v-for="(action, actionIndex) in threshold.actions">{{actionName(action.action)}}
			<span v-if="action.action == 'switch'">到{{action.options.ips.join(", ")}}</span>
			<span v-if="action.action == 'webHook'" class="small grey">({{action.options.url}})</span>
			 &nbsp;<span v-if="actionIndex != threshold.actions.length - 1" style="font-style: italic">AND &nbsp;</span></span>
			&nbsp;
			<a href="" title="修改" @click.prevent="update(index)"><i class="icon pencil small"></i></a> 
			<a href="" title="删除" @click.prevent="remove(index)"><i class="icon small remove"></i></a>
		</div>
	</div>
	
	<!-- 新阈值 -->
	<div v-if="isAdding" style="margin-top: 0.5em">
		<table class="ui table celled">
			<thead>
				<tr>
					<td style="width: 50%; background: #f9fafb; border-bottom: 1px solid rgba(34,36,38,.1)">阈值</td>
					<th>动作</th>
				</tr>
			</thead>
			<tr>
				<td style="background: white">
					<!-- 已经添加的项目 -->
					<div>
						<div v-for="(item, index) in addingThreshold.items" class="ui label basic small" style="margin-bottom: 0.5em;">
							<span v-if="item.item != 'nodeHealthCheck'">
								[{{item.duration}}{{itemDurationUnitName(item.durationUnit)}}]
							</span> 
							{{itemName(item.item)}}
							
							<span v-if="item.item == 'nodeHealthCheck'">
								<!-- 健康检查 -->
								<span v-if="item.value == 1">成功</span>
								<span v-if="item.value == 0">失败</span>
							</span>
							<span v-else>
								<!-- 连通性 -->
								<span v-if="item.item == 'connectivity' && item.options != null && item.options.groups != null && item.options.groups.length > 0">[<span v-for="(group, groupIndex) in item.options.groups">{{group.name}} <span v-if="groupIndex != item.options.groups.length - 1">&nbsp; </span></span>]</span>
								 <span class="grey">[{{itemOperatorName(item.operator)}}]</span> {{item.value}}{{itemUnitName(item.item)}}
							 </span> 
							 &nbsp;
							<a href="" title="删除" @click.prevent="removeItem(index)"><i class="icon remove small"></i></a>
						</div>
					</div>
					
					<!-- 正在添加的项目 -->
					<div v-if="isAddingItem" style="margin-top: 0.8em">
						<table class="ui table">
							<tr>
								<td style="width: 6em">统计项目</td>
								<td>
									<select class="ui dropdown auto-width" v-model="itemCode">
									<option v-for="item in allItems" :value="item.code">{{item.name}}</option>
									</select>
									<p class="comment" style="font-weight: normal" v-for="item in allItems" v-if="item.code == itemCode">{{item.description}}</p>
								</td>
							</tr>
							<tr v-show="itemCode != 'nodeHealthCheck'">
								<td>统计周期</td>
								<td>
									<div class="ui input right labeled">
										<input type="text" v-model="itemDuration" style="width: 4em" maxlength="4" ref="itemDuration" @keyup.enter="confirmItem()" @keypress.enter.prevent="1"/>
										<span class="ui label">分钟</span>
									</div>
								</td>
							</tr>
							<tr v-show="itemCode != 'nodeHealthCheck'">
								<td>操作符</td>
								<td>
									<select class="ui dropdown auto-width" v-model="itemOperator">
										<option v-for="operator in allOperators" :value="operator.code">{{operator.name}}</option>
									</select>
								</td>
							</tr>
							<tr v-show="itemCode != 'nodeHealthCheck'">
								<td>对比值</td>
								<td>
									<div class="ui input right labeled">
										<input type="text" maxlength="20" style="width: 5em" v-model="itemValue" ref="itemValue" @keyup.enter="confirmItem()" @keypress.enter.prevent="1"/>
										<span class="ui label" v-for="item in allItems" v-if="item.code == itemCode">{{item.unit}}</span>
									</div>
								</td>
							</tr>
							<tr v-show="itemCode == 'nodeHealthCheck'">
								<td>检查结果</td>
								<td>
									<select class="ui dropdown auto-width" v-model="itemValue">
										<option value="1">成功</option>
										<option value="0">失败</option>
									</select>
									<p class="comment" style="font-weight: normal">只有状态发生改变的时候才会触发。</p>
								</td>
							</tr>
							
							<!-- 连通性 -->
							<tr v-if="itemCode == 'connectivity'">
								<td>终端分组</td>
								<td style="font-weight: normal">
									<div style="zoom: 0.8"><report-node-groups-selector @change="changeReportGroups"></report-node-groups-selector></div>
								</td>
							</tr>
						</table>
						<div style="margin-top: 0.8em">
							<button class="ui button tiny" type="button" @click.prevent="confirmItem">确定</button>							 &nbsp;
							<a href="" title="取消" @click.prevent="cancelItem"><i class="icon remove small"></i></a>
						</div>
					</div>
					<div style="margin-top: 0.8em" v-if="!isAddingItem">
						<button class="ui button tiny" type="button" @click.prevent="addItem">+</button>
					</div>
				</td>
				<td style="background: white">
					<!-- 已经添加的动作 -->
					<div>
						<div v-for="(action, index) in addingThreshold.actions" class="ui label basic small" style="margin-bottom: 0.5em">
							{{actionName(action.action)}} &nbsp;
							<span v-if="action.action == 'switch'">到{{action.options.ips.join(", ")}}</span>
							<span v-if="action.action == 'webHook'" class="small grey">({{action.options.url}})</span>
							<a href="" title="删除" @click.prevent="removeAction(index)"><i class="icon remove small"></i></a>
						</div>
					</div>
					
					<!-- 正在添加的动作 -->
					<div v-if="isAddingAction" style="margin-top: 0.8em">
						<table class="ui table">
							<tr>
								<td style="width: 6em">动作类型</td>
								<td>
									<select class="ui dropdown auto-width" v-model="actionCode">
										<option v-for="action in allActions" :value="action.code">{{action.name}}</option>
									</select>
									<p class="comment" v-for="action in allActions" v-if="action.code == actionCode">{{action.description}}</p>
								</td>
							</tr>
							
							<!-- 切换 -->
							<tr v-if="actionCode == 'switch'">
								<td>备用IP *</td>
								<td>
									<textarea rows="2" v-model="actionBackupIPs" ref="actionBackupIPs"></textarea>
									<p class="comment">每行一个备用IP。</p>
								</td>
							</tr>
							
							<!-- WebHook -->
							<tr v-if="actionCode == 'webHook'">
								<td>URL *</td>
								<td>
									<input type="text" maxlength="1000" placeholder="https://..." v-model="actionWebHookURL" ref="webHookURL" @keyup.enter="confirmAction()" @keypress.enter.prevent="1"/>
									<p class="comment">完整的URL，比如<code-label>https://example.com/webhook/api</code-label>，系统会在触发阈值的时候通过GET调用此URL。</p>
								</td>
							</tr>
						</table>
						<div style="margin-top: 0.8em">
							<button class="ui button tiny" type="button" @click.prevent="confirmAction">确定</button>	 &nbsp;
							<a href="" title="取消" @click.prevent="cancelAction"><i class="icon remove small"></i></a>
						</div>
					</div>
					
					<div style="margin-top: 0.8em" v-if="!isAddingAction">
						<button class="ui button tiny" type="button" @click.prevent="addAction">+</button>
					</div>	
				</td>
			</tr>
		</table>
		
		<!-- 添加阈值 -->
		<div>
			<button class="ui button tiny" :class="{disabled: (isAddingItem || isAddingAction)}" type="button" @click.prevent="confirm">确定</button> &nbsp;
			<a href="" title="取消" @click.prevent="cancel"><i class="icon remove small"></i></a>
		</div>
	</div>
	
	<div v-if="!isAdding" style="margin-top: 0.5em">
		<button class="ui button tiny" type="button" @click.prevent="add">+</button>
	</div>
</div>`}),Vue.component("node-region-selector",{props:["v-region"],data:function(){return{selectedRegion:this.vRegion}},methods:{selectRegion:function(){let t=this;teaweb.popup("/clusters/regions/selectPopup?clusterId="+this.vClusterId,{callback:function(e){t.selectedRegion=e.data.region}})},addRegion:function(){let t=this;teaweb.popup("/clusters/regions/createPopup?clusterId="+this.vClusterId,{callback:function(e){t.selectedRegion=e.data.region}})},removeRegion:function(){this.selectedRegion=null}},template:`<div>
	<div class="ui label small basic" v-if="selectedRegion != null">
		<input type="hidden" name="regionId" :value="selectedRegion.id"/>
		{{selectedRegion.name}} &nbsp;<a href="" title="删除" @click.prevent="removeRegion()"><i class="icon remove"></i></a>
	</div>
	<div v-if="selectedRegion == null">
		<a href="" @click.prevent="selectRegion()">[选择区域]</a> &nbsp; <a href="" @click.prevent="addRegion()">[添加区域]</a>
	</div>
</div>`}),Vue.component("node-combo-box",{props:["v-cluster-id","v-node-id"],data:function(){let t=this;return Tea.action("/clusters/nodeOptions").params({clusterId:this.vClusterId}).post().success(function(e){t.nodes=e.data.nodes}),{nodes:[]}},template:`<div v-if="nodes.length > 0">
	<combo-box title="节点" placeholder="节点名称" :v-items="nodes" name="nodeId" :v-value="vNodeId"></combo-box>
</div>`}),Vue.component("node-level-selector",{props:["v-node-level"],data:function(){let e=this.vNodeLevel;return{levels:[{name:"边缘节点",code:1,description:"普通的边缘节点。"},{name:"L2节点",code:2,description:"特殊的边缘节点，同时负责同组上一级节点的回源。"}],levelCode:e=null==e||e<1?1:e}},template:`<div>
	<select class="ui dropdown auto-width" name="level" v-model="levelCode">
	<option v-for="level in levels" :value="level.code">{{level.name}}</option>
</select>
<p class="comment" v-if="typeof(levels[levelCode - 1]) != null"><plus-label
></plus-label>{{levels[levelCode - 1].description}}</p>
</div>`}),Vue.component("dns-route-selector",{props:["v-all-routes","v-routes"],data:function(){let e=this.vRoutes;return(e=null==e?[]:e).$sort(function(e,t){return e.domainId==t.domainId?e.code<t.code:e.domainId<t.domainId?1:-1}),{routes:e,routeCodes:e.$map(function(e,t){return t.code+"@"+t.domainId}),isAdding:!1,routeCode:"",keyword:"",searchingRoutes:this.vAllRoutes.$copy()}},methods:{add:function(){this.isAdding=!0,this.keyword="",this.routeCode="";let e=this;setTimeout(function(){e.$refs.keywordRef.focus()},200)},cancel:function(){this.isAdding=!1},confirm:function(){if(0!=this.routeCode.length)if(this.routeCodes.$contains(this.routeCode))teaweb.warn("已经添加过此线路，不能重复添加");else{let i=this;var e=this.vAllRoutes.$find(function(e,t){return t.code+"@"+t.domainId==i.routeCode});null!=e&&(this.routeCodes.push(this.routeCode),this.routes.push(e),this.routes.$sort(function(e,t){return e.domainId==t.domainId?e.code<t.code:e.domainId<t.domainId?1:-1}),this.routeCode="",this.isAdding=!1)}},remove:function(i){this.routeCodes.$removeValue(i.code+"@"+i.domainId),this.routes.$removeIf(function(e,t){return t.code+"@"+t.domainId==i.code+"@"+i.domainId})}},watch:{keyword:function(t){if(0==t.length)return this.searchingRoutes=this.vAllRoutes.$copy(),void(this.routeCode="");this.searchingRoutes=this.vAllRoutes.filter(function(e){return teaweb.match(e.name,t)||teaweb.match(e.domainName,t)}),0<this.searchingRoutes.length?this.routeCode=this.searchingRoutes[0].code+"@"+this.searchingRoutes[0].domainId:this.routeCode=""}},template:`<div>
	<input type="hidden" name="dnsRoutesJSON" :value="JSON.stringify(routeCodes)"/>
	<div v-if="routes.length > 0">
		<tiny-basic-label v-for="route in routes" :key="route.code + '@' + route.domainId">
			{{route.name}} <span class="grey small">（{{route.domainName}}）</span><a href="" @click.prevent="remove(route)"><i class="icon remove"></i></a>
		</tiny-basic-label>
		<div class="ui divider"></div>
	</div>
	<button type="button" class="ui button small" @click.prevent="add" v-if="!isAdding">+</button>
	<div v-if="isAdding">
		<div class="ui fields inline">
			<div class="ui field">
				<select class="ui dropdown" style="width: 18em" v-model="routeCode">
					<option value="" v-if="keyword.length == 0">[请选择]</option>
					<option v-for="route in searchingRoutes" :value="route.code + '@' + route.domainId">{{route.name}}（{{route.domainName}}）</option>
				</select>
			</div>
			<div class="ui field">
				<input type="text" placeholder="搜索..." size="10" v-model="keyword" ref="keywordRef" @keyup.enter="confirm" @keypress.enter.prevent="1"/>
			</div>
			
			<div class="ui field">
				<button class="ui button tiny" type="button" @click.prevent="confirm">确定</button>
			</div>
			<div class="ui field">
				<a href="" @click.prevent="cancel()"><i class="icon remove small"></i></a>
			</div>
		</div>
	</div>
</div>`}),Vue.component("dns-domain-selector",{props:["v-domain-id","v-domain-name"],data:function(){let e=this.vDomainId,t=(null==e&&(e=0),this.vDomainName);return null==t&&(t=""),{domainId:e,domainName:t}},methods:{select:function(){let t=this;teaweb.popup("/dns/domains/selectPopup",{callback:function(e){t.domainId=e.data.domainId,t.domainName=e.data.domainName,t.change()}})},remove:function(){this.domainId=0,this.domainName="",this.change()},update:function(){let t=this;teaweb.popup("/dns/domains/selectPopup?domainId="+this.domainId,{callback:function(e){t.domainId=e.data.domainId,t.domainName=e.data.domainName,t.change()}})},change:function(){this.$emit("change",{id:this.domainId,name:this.domainName})}},template:`<div>
	<input type="hidden" name="dnsDomainId" :value="domainId"/>
	<div v-if="domainName.length > 0">
		<span class="ui label small basic">
			{{domainName}}
			<a href="" @click.prevent="update"><i class="icon pencil small"></i></a>
			<a href="" @click.prevent="remove()"><i class="icon remove"></i></a>
		</span>
	</div>
	<div v-if="domainName.length == 0">
		<a href="" @click.prevent="select()">[选择域名]</a>
	</div>
</div>`}),Vue.component("dns-resolver-config-box",{props:["v-dns-resolver-config"],data:function(){let e=this.vDnsResolverConfig;return{config:e=null==e?{type:"default"}:e,types:[{name:"默认",code:"default"},{name:"CGO",code:"cgo"},{name:"Go原生",code:"goNative"}]}},template:`<div>
	<input type="hidden" name="dnsResolverJSON" :value="JSON.stringify(config)"/>
	<table class="ui table definition selectable">
		<tr>
			<td class="title">使用的DNS解析库</td>
			<td>
				<select class="ui dropdown auto-width" v-model="config.type">
					<option v-for="t in types" :value="t.code">{{t.name}}</option>
				</select>
				<p class="comment">边缘节点使用的DNS解析库。修改此项配置后，需要重启节点进程才会生效。<pro-warning-label></pro-warning-label></p>
			</td>
		</tr>
	</table>
	<div class="margin"></div>
</div>`}),Vue.component("grant-selector",{props:["v-grant","v-node-cluster-id","v-ns-cluster-id"],data:function(){return{grantId:null==this.vGrant?0:this.vGrant.id,grant:this.vGrant,nodeClusterId:null!=this.vNodeClusterId?this.vNodeClusterId:0,nsClusterId:null!=this.vNsClusterId?this.vNsClusterId:0}},methods:{select:function(){let t=this;teaweb.popup("/clusters/grants/selectPopup?nodeClusterId="+this.nodeClusterId+"&nsClusterId="+this.nsClusterId,{callback:e=>{t.grantId=e.data.grant.id,0<t.grantId&&(t.grant=e.data.grant),t.notifyUpdate()},height:"26em"})},create:function(){let t=this;teaweb.popup("/clusters/grants/createPopup",{height:"26em",callback:e=>{t.grantId=e.data.grant.id,0<t.grantId&&(t.grant=e.data.grant),t.notifyUpdate()}})},update:function(){if(null==this.grant)window.location.reload();else{let t=this;teaweb.popup("/clusters/grants/updatePopup?grantId="+this.grant.id,{height:"26em",callback:e=>{t.grant=e.data.grant,t.notifyUpdate()}})}},remove:function(){this.grant=null,this.grantId=0,this.notifyUpdate()},notifyUpdate:function(){this.$emit("change",this.grant)}},template:`<div>
	<input type="hidden" name="grantId" :value="grantId"/>
	<div class="ui label small basic" v-if="grant != null">{{grant.name}}<span class="small grey">（{{grant.methodName}}）</span><span class="small grey" v-if="grant.username != null && grant.username.length > 0">（{{grant.username}}）</span> <a href="" title="修改" @click.prevent="update()"><i class="icon pencil small"></i></a> <a href="" title="删除" @click.prevent="remove()"><i class="icon remove"></i></a> </div>
	<div v-if="grant == null">
		<a href="" @click.prevent="select()">[选择已有认证]</a> &nbsp; &nbsp; <a href="" @click.prevent="create()">[添加新认证]</a>
	</div>
</div>`}),window.REQUEST_COND_COMPONENTS=[{type:"url-extension",name:"URL扩展名",description:"根据URL中的文件路径扩展名进行过滤",component:"http-cond-url-extension",paramsTitle:"扩展名列表",isRequest:!0,caseInsensitive:!1},{type:"url-prefix",name:"URL前缀",description:"根据URL中的文件路径前缀进行过滤",component:"http-cond-url-prefix",paramsTitle:"URL前缀",isRequest:!0,caseInsensitive:!0},{type:"url-eq",name:"URL精准匹配",description:"检查URL中的文件路径是否一致",component:"http-cond-url-eq",paramsTitle:"URL完整路径",isRequest:!0,caseInsensitive:!0},{type:"url-regexp",name:"URL正则匹配",description:"使用正则表达式检查URL中的文件路径是否一致",component:"http-cond-url-regexp",paramsTitle:"正则表达式",isRequest:!0,caseInsensitive:!0},{type:"user-agent-regexp",name:"User-Agent正则匹配",description:"使用正则表达式检查User-Agent中是否含有某些浏览器和系统标识",component:"http-cond-user-agent-regexp",paramsTitle:"正则表达式",isRequest:!0,caseInsensitive:!0},{type:"params",name:"参数匹配",description:"根据参数值进行匹配",component:"http-cond-params",paramsTitle:"参数配置",isRequest:!0,caseInsensitive:!1},{type:"url-not-extension",name:"排除：URL扩展名",description:"根据URL中的文件路径扩展名进行过滤",component:"http-cond-url-not-extension",paramsTitle:"扩展名列表",isRequest:!0,caseInsensitive:!1},{type:"url-not-prefix",name:"排除：URL前缀",description:"根据URL中的文件路径前缀进行过滤",component:"http-cond-url-not-prefix",paramsTitle:"URL前缀",isRequest:!0,caseInsensitive:!0},{type:"url-not-eq",name:"排除：URL精准匹配",description:"检查URL中的文件路径是否一致",component:"http-cond-url-not-eq",paramsTitle:"URL完整路径",isRequest:!0,caseInsensitive:!0},{type:"url-not-regexp",name:"排除：URL正则匹配",description:"使用正则表达式检查URL中的文件路径是否一致，如果一致，则不匹配",component:"http-cond-url-not-regexp",paramsTitle:"正则表达式",isRequest:!0,caseInsensitive:!0},{type:"user-agent-not-regexp",name:"排除：User-Agent正则匹配",description:"使用正则表达式检查User-Agent中是否含有某些浏览器和系统标识，如果含有，则不匹配",component:"http-cond-user-agent-not-regexp",paramsTitle:"正则表达式",isRequest:!0,caseInsensitive:!0},{type:"mime-type",name:"内容MimeType",description:"根据服务器返回的内容的MimeType进行过滤。注意：当用于缓存条件时，此条件需要结合别的请求条件使用。",component:"http-cond-mime-type",paramsTitle:"MimeType列表",isRequest:!1,caseInsensitive:!1}],window.REQUEST_COND_OPERATORS=[{description:"判断是否正则表达式匹配",name:"正则表达式匹配",op:"regexp"},{description:"判断是否正则表达式不匹配",name:"正则表达式不匹配",op:"not regexp"},{description:"使用字符串对比参数值是否相等于某个值",name:"字符串等于",op:"eq"},{description:"参数值包含某个前缀",name:"字符串前缀",op:"prefix"},{description:"参数值包含某个后缀",name:"字符串后缀",op:"suffix"},{description:"参数值包含另外一个字符串",name:"字符串包含",op:"contains"},{description:"参数值不包含另外一个字符串",name:"字符串不包含",op:"not contains"},{description:"使用字符串对比参数值是否不相等于某个值",name:"字符串不等于",op:"not"},{description:"判断参数值在某个列表中",name:"在列表中",op:"in"},{description:"判断参数值不在某个列表中",name:"不在列表中",op:"not in"},{description:"判断小写的扩展名（不带点）在某个列表中",name:"扩展名",op:"file ext"},{description:"判断MimeType在某个列表中，支持类似于image/*的语法",name:"MimeType",op:"mime type"},{description:"判断版本号在某个范围内，格式为version1,version2",name:"版本号范围",op:"version range"},{description:"将参数转换为整数数字后进行对比",name:"整数等于",op:"eq int"},{description:"将参数转换为可以有小数的浮点数字进行对比",name:"浮点数等于",op:"eq float"},{description:"将参数转换为数字进行对比",name:"数字大于",op:"gt"},{description:"将参数转换为数字进行对比",name:"数字大于等于",op:"gte"},{description:"将参数转换为数字进行对比",name:"数字小于",op:"lt"},{description:"将参数转换为数字进行对比",name:"数字小于等于",op:"lte"},{description:"对整数参数值取模，除数为10，对比值为余数",name:"整数取模10",op:"mod 10"},{description:"对整数参数值取模，除数为100，对比值为余数",name:"整数取模100",op:"mod 100"},{description:"对整数参数值取模，对比值格式为：除数,余数，比如10,1",name:"整数取模",op:"mod"},{description:"将参数转换为IP进行对比",name:"IP等于",op:"eq ip"},{description:"将参数转换为IP进行对比",name:"IP大于",op:"gt ip"},{description:"将参数转换为IP进行对比",name:"IP大于等于",op:"gte ip"},{description:"将参数转换为IP进行对比",name:"IP小于",op:"lt ip"},{description:"将参数转换为IP进行对比",name:"IP小于等于",op:"lte ip"},{description:"IP在某个范围之内，范围格式可以是英文逗号分隔的ip1,ip2，或者CIDR格式的ip/bits",name:"IP范围",op:"ip range"},{description:"对IP参数值取模，除数为10，对比值为余数",name:"IP取模10",op:"ip mod 10"},{description:"对IP参数值取模，除数为100，对比值为余数",name:"IP取模100",op:"ip mod 100"},{description:"对IP参数值取模，对比值格式为：除数,余数，比如10,1",name:"IP取模",op:"ip mod"},{description:"判断参数值解析后的文件是否存在",name:"文件存在",op:"file exist"},{description:"判断参数值解析后的文件是否不存在",name:"文件不存在",op:"file not exist"}],window.REQUEST_VARIABLES=[{code:"${edgeVersion}",description:"",name:"边缘节点版本"},{code:"${remoteAddr}",description:"会依次根据X-Forwarded-For、X-Real-IP、RemoteAddr获取，适合前端有别的反向代理服务时使用，存在伪造的风险",name:"客户端地址（IP）"},{code:"${rawRemoteAddr}",description:"返回直接连接服务的客户端原始IP地址",name:"客户端地址（IP）"},{code:"${remotePort}",description:"",name:"客户端端口"},{code:"${remoteUser}",description:"",name:"客户端用户名"},{code:"${requestURI}",description:"比如/hello?name=lily",name:"请求URI"},{code:"${requestPath}",description:"比如/hello",name:"请求路径（不包括参数）"},{code:"${requestURL}",description:"比如https://example.com/hello?name=lily",name:"完整的请求URL"},{code:"${requestLength}",description:"",name:"请求内容长度"},{code:"${requestMethod}",description:"比如GET、POST",name:"请求方法"},{code:"${requestFilename}",description:"",name:"请求文件路径"},{code:"${scheme}",description:"",name:"请求协议，http或https"},{code:"${proto}","description:":"类似于HTTP/1.0",name:"包含版本的HTTP请求协议"},{code:"${timeISO8601}",description:"比如2018-07-16T23:52:24.839+08:00",name:"ISO 8601格式的时间"},{code:"${timeLocal}",description:"比如17/Jul/2018:09:52:24 +0800",name:"本地时间"},{code:"${msec}",description:"比如1531756823.054",name:"带有毫秒的时间"},{code:"${timestamp}",description:"",name:"unix时间戳，单位为秒"},{code:"${host}",description:"",name:"主机名"},{code:"${serverName}",description:"",name:"接收请求的服务器名"},{code:"${serverPort}",description:"",name:"接收请求的服务器端口"},{code:"${referer}",description:"",name:"请求来源URL"},{code:"${referer.host}",description:"",name:"请求来源URL域名"},{code:"${userAgent}",description:"",name:"客户端信息"},{code:"${contentType}",description:"",name:"请求头部的Content-Type"},{code:"${cookies}",description:"",name:"所有cookie组合字符串"},{code:"${cookie.NAME}",description:"",name:"单个cookie值"},{code:"${isArgs}",description:"如果URL有参数，则值为`?`；否则，则值为空",name:"问号（?）标记"},{code:"${args}",description:"",name:"所有参数组合字符串"},{code:"${arg.NAME}",description:"",name:"单个参数值"},{code:"${headers}",description:"",name:"所有Header信息组合字符串"},{code:"${header.NAME}",description:"",name:"单个Header值"},{code:"${geo.country.name}",description:"",name:"国家/地区名称"},{code:"${geo.country.id}",description:"",name:"国家/地区ID"},{code:"${geo.province.name}",description:"目前只包含中国省份",name:"省份名称"},{code:"${geo.province.id}",description:"目前只包含中国省份",name:"省份ID"},{code:"${geo.city.name}",description:"目前只包含中国城市",name:"城市名称"},{code:"${geo.city.id}",description:"目前只包含中国城市",name:"城市名称"},{code:"${isp.name}",description:"",name:"ISP服务商名称"},{code:"${isp.id}",description:"",name:"ISP服务商ID"},{code:"${browser.os.name}",description:"客户端所在操作系统名称",name:"操作系统名称"},{code:"${browser.os.version}",description:"客户端所在操作系统版本",name:"操作系统版本"},{code:"${browser.name}",description:"客户端浏览器名称",name:"浏览器名称"},{code:"${browser.version}",description:"客户端浏览器版本",name:"浏览器版本"},{code:"${browser.isMobile}",description:"如果客户端是手机，则值为1，否则为0",name:"手机标识"}],window.METRIC_HTTP_KEYS=[{name:"客户端地址（IP）",code:"${remoteAddr}",description:"会依次根据X-Forwarded-For、X-Real-IP、RemoteAddr获取，适用于前端可能有别的反向代理的情形，存在被伪造的可能",icon:""},{name:"直接客户端地址（IP）",code:"${rawRemoteAddr}",description:"返回直接连接服务的客户端原始IP地址",icon:""},{name:"客户端用户名",code:"${remoteUser}",description:"通过基本认证填入的用户名",icon:""},{name:"请求URI",code:"${requestURI}",description:"包含参数，比如/hello?name=lily",icon:""},{name:"请求路径",code:"${requestPath}",description:"不包含参数，比如/hello",icon:""},{name:"完整URL",code:"${requestURL}",description:"比如https://example.com/hello?name=lily",icon:""},{name:"请求方法",code:"${requestMethod}",description:"比如GET、POST等",icon:""},{name:"请求协议Scheme",code:"${scheme}",description:"http或https",icon:""},{name:"文件扩展名",code:"${requestPathExtension}",description:"请求路径中的文件扩展名，包括点符号，比如.html、.png",icon:""},{name:"主机名",code:"${host}",description:"通常是请求的域名",icon:""},{name:"请求协议Proto",code:"${proto}",description:"包含版本的HTTP请求协议，类似于HTTP/1.0",icon:""},{name:"HTTP协议",code:"${proto}",description:"包含版本的HTTP请求协议，类似于HTTP/1.0",icon:""},{name:"URL参数值",code:"${arg.NAME}",description:"单个URL参数值",icon:""},{name:"请求来源URL",code:"${referer}",description:"请求来源Referer URL",icon:""},{name:"请求来源URL域名",code:"${referer.host}",description:"请求来源Referer URL域名",icon:""},{name:"Header值",code:"${header.NAME}",description:"单个Header值，比如${header.User-Agent}",icon:""},{name:"Cookie值",code:"${cookie.NAME}",description:"单个cookie值，比如${cookie.sid}",icon:""},{name:"状态码",code:"${status}",description:"",icon:""},{name:"响应的Content-Type值",code:"${response.contentType}",description:"",icon:""}],window.IP_ADDR_THRESHOLD_ITEMS=[{code:"nodeAvgRequests",description:"当前节点在单位时间内接收到的平均请求数。",name:"节点平均请求数",unit:"个"},{code:"nodeAvgTrafficOut",description:"当前节点在单位时间内发送的下行流量。",name:"节点平均下行流量",unit:"M"},{code:"nodeAvgTrafficIn",description:"当前节点在单位时间内接收的上行流量。",name:"节点平均上行流量",unit:"M"},{code:"nodeHealthCheck",description:"当前节点健康检查结果。",name:"节点健康检查结果",unit:""},{code:"connectivity",description:"通过区域监控得到的当前IP地址的连通性数值，取值在0和100之间。",name:"IP连通性",unit:"%"},{code:"groupAvgRequests",description:"当前节点所在分组在单位时间内接收到的平均请求数。",name:"分组平均请求数",unit:"个"},{code:"groupAvgTrafficOut",description:"当前节点所在分组在单位时间内发送的下行流量。",name:"分组平均下行流量",unit:"M"},{code:"groupAvgTrafficIn",description:"当前节点所在分组在单位时间内接收的上行流量。",name:"分组平均上行流量",unit:"M"},{code:"clusterAvgRequests",description:"当前节点所在集群在单位时间内接收到的平均请求数。",name:"集群平均请求数",unit:"个"},{code:"clusterAvgTrafficOut",description:"当前节点所在集群在单位时间内发送的下行流量。",name:"集群平均下行流量",unit:"M"},{code:"clusterAvgTrafficIn",description:"当前节点所在集群在单位时间内接收的上行流量。",name:"集群平均上行流量",unit:"M"}],window.IP_ADDR_THRESHOLD_ACTIONS=[{code:"up",description:"上线当前IP。",name:"上线"},{code:"down",description:"下线当前IP。",name:"下线"},{code:"notify",description:"发送已达到阈值通知。",name:"通知"},{code:"switch",description:"在DNS中记录中将IP切换到指定的备用IP。",name:"切换"},{code:"webHook",description:"调用外部的WebHook。",name:"WebHook"}];
