import '../css/corpmappings.scss'
import '../../shared/js/polyfills'
import api from "./api"
import Vue from "vue"
import CorpMappingInterface from "../components/CorpMappingInterface";
import vSelect from 'vue-select'

Vue.component('v-select', vSelect)
api.init('/admin/api')

const Application = new Vue({
  render: h => h(CorpMappingInterface)
}).$mount('#corpmappings-app');
