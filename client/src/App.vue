<template>
  <div id="app">
    <img src="./assets/logo1.png">
    <div class="block">
      <div class="head">
        <button class="btn" :class="{active: isactive1}" @click="showForm(1)"><span>Sign in</span></button>
        <button class="btn" :class="{active: isactive2}" @click="showForm(2)"><span>Sign up</span></button>
      <div class="line-form">
        <app-sign-up v-if="state == 'sign-up'"></app-sign-up>
        <app-sign-in v-else-if="state == 'sign-in'"></app-sign-in>
      </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'app',
  data () {
    return {
      state: 'sign-up',
      isactive1: true,
      isactive2: false
    }
  },
  methods:{
    showForm(index){
      switch (index) {
        case 1:
          this.state = 'sign-in'
          this.isactive1 = true
          this.isactive2 = false
          break;
        case 2:
          this.state = 'sign-up'
          this.isactive1 = false
          this.isactive2 = true
          break;  
        default:
          break;
      }
    }
  }
}
</script>

<style>
body{
  background: rgb(234, 241, 234);
}
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 20px;
}
.block{
  box-shadow: 1px 1px 1px 1px gray;
  max-width: 500px;
  width: auto;
  height: auto;
  min-height: 200px;
  margin: 20px auto;
  padding: 10px;
  background: white;
}
.line-form{
  width: 85%;
  margin: 20px auto;
}
.btn.focus, .btn:focus {
    outline: none !important;
    box-shadow: none !important;;
}
span {
    font-size: 30px;
}
.active{
  border-bottom: 2px solid red;
}
</style>
