<template>
    <form class="" @submit.prevent="signInUser">
        <div class="form-group row">
            <lable for="Email">Enter email</lable>
            <input type="email" class="form-control" id="Email" placeholder="google@gmail.com" required v-model="user.email">
        </div>
        <div class="form-group row">
            <lable for="Pass">Enter password</lable>
            <input type="password" class="form-control" id="Pass" placeholder="password" required v-model="user.password">
        </div>
        <div class="checkbox">
            <div class="form-check">
                <input type="checkbox" class="form-check-input" id="exampleCheck">
                <label class="form-check-label" for="exampleCheck">Remember me</label>
            </div>
            <div>
                <a href="#" class="forgotPass">Forgot your password?</a>
            </div>
        </div>
        <button type="submit" class="btn btn-primary">Sign in</button>
    </form>
</template>

<script>
export default {
    name: 'sign in',
    data(){
        return {
            user:{
                email: '',
                password: ''
            },
        }
    },
    // methods: {
    //     registerUse(){
    //         if (this.user.password != this.user.confirmpassword || this.user.password.length < 6){
    //             this.error = true;
    //         }
    //     }
    // }

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
.checkbox{
    display: flex;
    justify-content: space-between;
}
label.form-check-label{
    margin-left: -15px;
}
#exampleCheck{
    margin-top: 5px;
}

</style>
