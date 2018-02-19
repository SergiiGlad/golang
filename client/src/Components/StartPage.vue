<template>
  <div class="start">
    <img src="../assets/logo1.png">
    <div class="block">
      <div class="head">
        <button class="btn-top" :class="{active: isactive1}" @click="showForm(1)"><span>Sign in</span></button>
        <button class="btn-top" :class="{active: isactive2}" @click="showForm(2)"><span>Sign up</span></button>
      <div class="line-form">
        <app-sign-up v-if="state == 'sign-up'"></app-sign-up>
        <app-sign-in v-else-if="state == 'sign-in'"
                       @recovery = "recoveryPass"
        ></app-sign-in>
        <app-recovery-pass v-else-if="state == 'recovery-pass'"></app-recovery-pass>
      </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'start',
  data () {
    return {
      state: 'sign-in',
      isactive1: true,
      isactive2: false
    }
  },
  methods:{
    showForm(index){
      switch (index) {
        case 1:
          this.state = 'sign-in';
          this.isactive1 = true;
          this.isactive2 = false;
          break;
        case 2:
          this.state = 'sign-up';
          this.isactive1 = false;
          this.isactive2 = true;
          break;
        default:
          break;
      }
    },
    recoveryPass(){
      this.state = "recovery-pass"
    }
  }
}
</script>


<style scoped>
    .start {
      text-align: center;
    }
    .block{
      box-shadow: 7px 7px 10px -3px gray;
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
    .btn-top{
        outline: none !important;
        box-shadow: none !important;
        background: none !important;
        border: none;
    }
    span {
        font-size: 30px;
    }
    .active{
      border-bottom: 2px solid red;
      font-weight: bolder;
    }
</style>
