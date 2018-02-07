<template>
  <form method="GET" enctype="multipart/form-data" @submit.prevent="recoveryPass" >
      <div class="alert alert-info" role="alert" v-if="success">
        <strong>We send on your email a new password!</strong>
      </div>
      <div class="form-group row">
          <lable for="Email">Enter your email</lable>
          <input type="email" class="form-control" id="Email" placeholder="google@gmail.com" v-model="user.email" required>
      </div>
      <div class="alert alert-danger" role="alert" v-if="error">
        <strong>Not valid email!</strong>
      </div>
      <button type="submit" class="btn">Recovery password</button>

  </form>
</template>

<script>

import axsios from 'axios'

export default {
     data(){
        return {
            user:{
                email: ''
            },
            success: false,
            error: false
        }
     },
     methods: {
         recoveryPass: function() {

           this.error = false;
           this.success = false;

           if (!(/^[a-z0-9]+@[a-z]+[.][a-z]+$/.test(this.user.email))) {
             this.error = true;
             return;
           } else{
             this.success = true;
           }

           var axios = require('axios');

           var recoveryPass = 'http://localhost:8080/recoveryPass?email=' + this.user.email;
           let that = this;

           axios.get(recoveryPass, {headers: { 'Access-Control-Allow-Origin': '*'}})
             .then(function (response) {
               console.log(response.data);
               that.success = true;
             })
             .catch(function (error) {
               console.log(error.message);
             });
         }
    }
}

</script>

<style scoped>

    input{
      background-color: rgb(234, 241, 234);
    }
    .btn{
      font-size: 20px;
      font-weight: bold;
      border: none;
      background: rgb(184, 56, 56);
      color: white;
    }
    .btn:hover{
      background:white;
      border: 1px solid black;
      color: black;
    }

</style>
