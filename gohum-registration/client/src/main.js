import Vue from 'vue'
import App from './App.vue'

import AppSignUp from './Components/SignUp.vue'
import AppSignIn from './Components/SignIn.vue'
import AppRecoveryPass from './Components/RecoveryPass.vue'

Vue.component('AppSignUp', AppSignUp)
Vue.component('AppSignIn', AppSignIn)
Vue.component('AppRecoveryPass', AppRecoveryPass)

new Vue({
  el: '#app',
  render: h => h(App)
})
