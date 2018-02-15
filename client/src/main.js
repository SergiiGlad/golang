import Vue from 'vue'
import App from './App.vue'
import router  from './routers'
import BootstrapVue from 'bootstrap-vue'
import AsyncComputed from 'vue-async-computed'

Vue.use(BootstrapVue);
Vue.use(AsyncComputed);


new Vue({
  router,
  el: '#app',
  render: h => h(App)
})
