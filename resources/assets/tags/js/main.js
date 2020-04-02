import '../css/tags.scss'
import '../../shared/js/polyfills'
import Vue from "vue"
import TagsTable from "../components/TagsTable"

const Application = new Vue({
  render: h => h(TagsTable)
}).$mount('#tags-table-app');
