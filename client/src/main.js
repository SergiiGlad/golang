import Vue from 'vue'
import App from './App.vue'
import router  from './routers'
import BootstrapVue from 'bootstrap-vue'
import AsyncComputed from 'vue-async-computed'

// import Icon from 'vue-awesome/icons'


Vue.use(BootstrapVue);
Vue.use(AsyncComputed);
//Vue.component('icon', Icon)


new Vue({
  router,
  el: '#app',
  render: h => h(App)
})
