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
      //console.log(to.path)
      this.currentRoute = to.path
    }

 },
}
</script>

<style>
  body {
    background-color: rgb(234, 241, 234);
  }
</style>
<style scoped>
    #app {
      font-family: 'Avenir', Helvetica, Arial, sans-serif;
      -webkit-font-smoothing: antialiased;
      -moz-osx-font-smoothing: grayscale;
      color: #2c3e50;
      margin-top: 20px;
      /* background-color: rgb(234, 241, 234); */

    }
.container {
    margin-top: 80px;
    max-width: 1000px;
    width: 100%;
    /* box-shadow: 1px 1px 1px 1px black; */
    height: auto;
}
</style>
