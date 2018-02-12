<template>
    <form  method="post" enctype="multipart/form-data" @submit.prevent="signInUser">
        <div class="form-group row">
            <label for="Email">Enter email</label>
            <input type="email" class="form-control" id="Email" placeholder="google@gmail.com" v-model="user.email">
        </div>
        <div class="form-group row">
            <label for="Pass">Enter password</label>
            <input type="password" class="form-control" id="Pass" placeholder="password" v-model="user.password">
        </div>
      <div class="alert alert-danger" role="alert" v-if="error">
        <strong>{{message}}</strong>
      </div>
        <div class="checkbox">
            <div class="form-check">
                <input type="checkbox" class="form-check-input" id="exampleCheck" checked>
                <label class="form-check-label" for="exampleCheck">Remember me</label>
            </div>
            <div>
                <a href="#" class="forgotPass"  @click="$emit('recovery')">Forgot your password?</a>
            </div>
        </div>
        <button type="submit" class="btn">Sign in</button>
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
          error: false,
          message: ''
        }
    },
    method: {
      signInUser(){

        this.error = false;

        if (!(/^[a-z0-9]+@[a-z]+[.][a-z]+$/.test(this.user.email))) {
          this.error = true;
          this.message = "Invalid email!";
          return;
        }
        if (this.user.password.length < 6) {
          this.error = true;
          this.message = "Password must be more than 6 characters!";
          return;
        }

      }
    }
}
</script>

<style scoped>

    input{
      background-color: rgb(234, 241, 234);
    }
    .btn {
      font-size: 20px;
      font-weight: bold;
      border: none;
      background: rgb(184, 56, 56);
      color: white;
    }
    .btn:hover, .btn:active{
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
    input#exampleCheck.form-check-input{
        margin-top: 5px;
        margin-left: -13px;
    }

</style>
