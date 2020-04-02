import '../css/locations.scss'
import '../../shared/js/polyfills'
import api from "./api"
import Vue from "vue"
import LocationsTable from "../components/LocationsTable"

api.init('/admin/api');

const Application = new Vue({
  render: h => h(LocationsTable)
}).$mount('#locations-table-app');
