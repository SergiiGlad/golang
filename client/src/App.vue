<template>
  <div id="app">
    <top-menu v-show="currentPath != '/'"></top-menu>
    <div class="container">
    <router-view class="view"></router-view>
    </div>

  </div>
</template>

<script>

import  router  from './routers'

export default {
  name: 'app',
  data () {
    return {
      currentPath: router.currentRoute.path
    }
  },
  watch: {
    '$route' (to, from) {
      console.log(to.path)
      this.currentRoute = to.path
    }

  },
}
</script>

<style scoped>
    #app {
      font-family: 'Avenir', Helvetica, Arial, sans-serif;
      -webkit-font-smoothing: antialiased;
      -moz-osx-font-smoothing: grayscale;
      color: #2c3e50;
      margin-top: 20px;
    }
    .container{
      margin-top: 90px;
      max-width: 700px;
      background-color: aliceblue;
}
</style>
