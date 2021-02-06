Vue.component("ip-item-text", {
    props: ["v-item"],
    template: `<span>
    <span v-if="vItem.type == 'all'">*</span>
    <span v-if="vItem.type == 'ipv4' || vItem.type.length == 0">
        {{vItem.ipFrom}}
        <span v-if="vItem.ipTo.length > 0">- {{vItem.ipTo}}</span>
    </span>
    <span v-if="vItem.type == 'ipv6'">{{vItem.ipFrom}}</span>
    <span v-if="vItem.eventLevelName != null && vItem.eventLevelName.length > 0">&nbsp; 级别：{{vItem.eventLevelName}}</span>
</span>`
})