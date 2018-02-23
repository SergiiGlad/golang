"use strict";

import Vue from 'vue'
import VueRouter from 'vue-router'

import AppStartPage from '../Components/StartPage.vue'
import AppSignUp from '../Components/SignUp.vue'
import AppSignIn from '../Components/SignIn.vue'
import AppRecoveryPass from '../Components/RecoveryPass.vue'
import AppProfile from '../Components/Profile.vue'
import TopMenu from '../Components/TopMenu.vue'

import Friends from '../Components/Friends.vue'
import Messages from '../Components/Messages.vue'
import Photos from '../Components/Photos.vue'




Vue.component('AppSignUp', AppSignUp)
Vue.component('AppSignIn', AppSignIn)
Vue.component('AppRecoveryPass', AppRecoveryPass)
Vue.component('StartPage', AppStartPage)
Vue.component('Profile', AppProfile)
Vue.component('TopMenu', TopMenu)

Vue.component('Friends', Friends)
Vue.component('Messages', Messages)
Vue.component('Photos', Photos)




Vue.use(VueRouter);

export default new VueRouter({
  mode: 'history',
  base: __dirname,
  routes: [
    {
      path: '/profile/:id',
      name: 'Profile',
      component: AppProfile
    },
    {
        path: '/',
        name: 'StartPage',
        component: AppStartPage
      },
    // {
    //     path: '/login',
    //     name: 'Login',
    //     component: AppSignIn
    // },
    // {
    //     path: '/register',
    //     name: 'Register',
    //     component: AppSignUp
    // },
    // {
    //     path: '/recoveryPass',
    //     name: 'RecoveryPass',
    //     component: AppRecoveryPass
    // },
    // inside pages
    {
      path: '/friends',
      name: 'Friends',
      component: Friends
  },
  {
    path: '/messages',
    name: 'Messages',
    component: Messages
},
{
  path: '/photos',
  name: 'Photos',
  component: Photos
},
  ]
})
