import Vue from 'vue'
import App from './App.vue'

import AppSignUp from './Components/SignUp.vue'
import AppSignIn from './Components/SignIn.vue'


Vue.component('AppSignUp', AppSignUp)
Vue.component('AppSignIn', AppSignIn)

new Vue({
  el: '#app',
  render: h => h(App)
})
