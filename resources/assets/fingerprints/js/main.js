import '../css/fingerprints.scss'
import '../../shared/js/polyfills'
import api from "./api"
import Vue from "vue"
import FingerprintTable from "../components/FingerprintTable"
import VueTippy from "vue-tippy";

Vue.use(VueTippy);

api.init('/admin/api')

const Application = new Vue({
  render: h => h(FingerprintTable)
}).$mount('#fingerprint-app');
