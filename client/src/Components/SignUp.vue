<template>
    <form  method="POST" enctype="multipart/form-data" @submit.prevent="registerUser()">
      <div class="alert alert-success" role="alert" v-if="success">
        <strong>Thanks for registration! Check your Email!</strong>
      </div>
        <div class="form-group row">
            <label for="Name" >First name</label>
            <input type="text"  class="form-control" id="Name" placeholder="Enter your first name"  v-model="user.first_name">
        </div>
        <div class="form-group row">
            <label for="Surname">Last name</label>
            <input type="text" class="form-control" id="Surname" placeholder="Enter your last name"  v-model="user.last_name">
        </div>
      <div class="form-group row">
        <label for="Phone">Phone number</label>
        <input  type="tel" class="form-control" id="Phone" placeholder="+38050XXXXXXX"  v-model="user.phone">
      </div>
        <div class="form-group row">
            <label for="Email">Email</label>
            <input type="email" class="form-control" id="Email" placeholder="Enter your email" v-model="user.email">
        </div>
        <div class="form-group row">
            <label for="Pass">Password</label>
            <input type="password" class="form-control" id="Pass" placeholder="******"  v-model="user.password">
        </div>
        <div class="form-group row">
            <label for="ConfirmPass">Confirm password</label>
            <input type="password" class="form-control" id="ConfirmPass" placeholder="******"  v-model="confirmPass">
        </div>
        <div class="alert alert-danger" role="alert" v-if="error">
            <strong>{{message}}</strong>
        </div>
        <button type="submit" class="btn btn-primary">Register</button>
    </form>
</template>

<script>

  import axsios from 'axios'

  export default {
      name: 'sign-up',
      data(){
          return {
              user:{
                  first_name: '',
                  last_name: '',
                  phone: '',
                  email: '',
                  password: ''
              },
            confirmPass: '',
            message: '',
            error: false,
            success: false
          }
      },
      methods: {
          registerUser(){

            this.error = false;
            this.success = false;

            if (!(/^[a-zA-Z]{1,50}$/.test(this.user.first_name)) && !(/^[а-яА-Я]{1,50}$/.test(this.user.first_name))) {
              this.error = true;
              this.message = "Invalid name!";
              return;
            }
            if (!(/^[a-zA-Z]{1,50}$/.test(this.user.last_name)) && !(/^[а-яА-Я]{1,50}$/.test(this.user.last_name))) {
              this.error = true;
              this.message = "Invalid last name!";
              return;
            }
            if (!(/^[+][0-9]{12}$/.test(this.user.phone))) {
              this.error = true;
              this.message = "Invalid phone number!";
              return;
            }
            if (!(/^[a-z0-9]+@[a-z]+[.][a-z]+$/.test(this.user.email))) {
              this.error = true;
              this.message = "Invalid email!";
              return;
            }
            if (this.user.password !== this.confirmPass || this.user.password.length < 6) {
              this.error = true;
              this.message = "Password are different or less than 6 characters!";
              return;
            }
            else{
              this.user.first_name = this.user.first_name.charAt(0).toUpperCase() + this.user.first_name.slice(1);
              this.user.last_name = this.user.last_name.charAt(0).toUpperCase() + this.user.last_name.slice(1);

              var axios = require('axios');
              let that = this;
              axios.post('http://localhost:8080/api/registration', this.user, {
                headers: {
                            "accept": "application/json",
                            "Content-Type": "application/json" }
              })
                .then(function (response) {
                  console.log(response.data);
                  that.success = true;
                })
                .catch(function (error){
                  console.log(error.message);
                });
            }
          },
      }
  }

</script>

<style scoped>
  input{
    background-color: rgb(234, 241, 234);
  }
  button.btn{
    font-size: 20px;
    font-weight: bold;
    border: none;
    background: rgb(184, 56, 56);
  }
  button:hover{
      background: white;
      border: 1px solid black;
      color: black;
  }

</style>
