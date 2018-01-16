<template>
    <form class="" @submit.prevent="registerUser">
        <div class="form-group row">
            <lable for="Name" >First name</lable>
            <input type="text" class="form-control" id="Name" placeholder="Enter your first name" required v-model="user.firstName">
        </div>
        <div class="form-group row">
            <lable for="Surname">Last name</lable>
            <input type="text" class="form-control" id="Surname" placeholder="Enter your last name" required v-model="user.lastName">
        </div>
        <div class="form-group row">
            <lable for="Email">Email</lable>
            <input type="email" class="form-control" id="Email" placeholder="Enter your email" required v-model="user.email">
        </div>
        <div class="form-group row">
            <lable for="Pass">Password</lable>
            <input type="password" class="form-control" id="Pass" placeholder="******" required v-model="user.password">
        </div>
        <div class="form-group row">
            <lable for="ConfirmPass">Confirm password</lable>
            <input type="password" class="form-control" id="ConfirmPass" placeholder="******" required v-model="user.confirmPassword">
        </div>
        <div class="alert alert-danger" role="alert" v-if="error">
            <strong>Passwords are differents!</strong>
        </div>
        <button type="submit" class="btn btn-primary">Register</button>
    </form>
</template>

<script>
export default {
    name: 'sign up',
    data(){
        return {
            user:{
                firstName: '',
                lastName: '',
                email: '',
                password: '',
                confirmpassword: ''
            },
         error: false
        }
    },
    methods: {
        registerUse(){
            if (this.user.password != this.user.confirmpassword || this.user.password.length < 6){
                this.error = true;
            }
        }
    }
}
</script>

<style scoped>
input{
  background-color: rgb(245, 244, 244);
}
button.btn{
  font-size: 20px;
  border: none;
  background: rgb(184, 56, 56);
}
button:hover{
    background:white;
    border: 1px solid black;
    color: black;
}
</style>
