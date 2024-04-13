Vue.component("ip-item-text", {
    props: ["v-item"],
    template: `<span>
    <span v-if="vItem.type == 'all'">*</span>
    <span v-else>
    	<span v-if="vItem.value != null && vItem.value.length > 0">{{vItem.value}}</span>
    	<span v-else>
			{{vItem.ipFrom}}
			<span v-if="vItem.ipTo != null &&vItem.ipTo.length > 0">- {{vItem.ipTo}}</span>
		</span>
	</span>
    <span v-if="vItem.eventLevelName != null && vItem.eventLevelName.length > 0">&nbsp; 级别：{{vItem.eventLevelName}}</span>
</span>`
})