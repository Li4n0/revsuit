<template>
  <div
      slot="filterDropdown"
      style="padding: 8px"
  >
    <a-input
        :placeholder="`Search ${column.dataIndex}`"
        :value="selectedKeys[0]"
        style="width: 188px; margin-bottom: 8px; display: block;"
        @change="e => setSelectedKeys(e.target.value ? [e.target.value] : [])"
        @pressEnter="() => {
            filters[column.dataIndex] = selectedKeys[0];
            fetch()
          }"
    />
    <a-button
        type="primary"
        icon="search"
        size="small"
        style="width: 90px; margin-right: 8px"
        @click="() => {
            filters[column.dataIndex] = selectedKeys[0];
            fetch()
          }"
    >
      Search
    </a-button>
    <a-button size="small" style="width: 90px" @click="() =>{
        clearFilters();
        delete filters[column.dataIndex];
        fetch()
      }">
      Reset
    </a-button>
  </div>
</template>

<script>
export default {
name: "FilterDropdown",
  props:["setSelectedKeys","selectedKeys","clearFilters","column","filters","fetch"],

}
</script>

<style scoped>

</style>