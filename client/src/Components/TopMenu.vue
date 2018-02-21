<template>
    <div>
        <header>
            <div class="container">
            <img src="../assets/logo1.png">
            <button type="button">Logout</button>
            </div>
        </header>
        <div class="menu">
            <div class="container-2">
            <router-link :to="{ name: 'Profile', params: {id: user_id}}">Profile</router-link>
            <router-link :to="'/friends'">Friends</router-link>
            <router-link :to="'/messages'">Messages</router-link>
            <router-link :to="'/photos'">Photos</router-link>
            </div> 
        </div>
    </div>
</template>

<script>
export default {
  name: 'topMenu',
  data () {
    return {
    }
  },
  computed:{
     user_id(){
         console.log(localStorage.getItem("id"));
          return localStorage.getItem("id");
      }
  }
 
}
</script>

<style scoped>
header {
    background-color: rgb(245, 235, 235);
    box-shadow: 0px 0px 15px -1px gray ;
    height: 50px;
    position: fixed;
    top:0;
    width: 100%;
    z-index: 21;    
}
img {
    height: 40px;
    width: auto;
    float: left;

}
button {
    float: right;
    border-radius: 4px;
    background-color: white;
    color: black;
    margin-top: 10px;
    outline: none;
    cursor: pointer;

     background-color: #49afcd;
        text-shadow: 0 -1px 0 rgba(0, 0, 0, 0.25);
        background-color: #49afcd;
        background-image: -moz-linear-gradient(top, #5bc0de, #2f96b4);
        background-image: -webkit-gradient(linear, 0 0, 0 100%, from(#5bc0de), to(#2f96b4));
        background-image: -webkit-linear-gradient(top, #5bc0de, #2f96b4);
        background-image: -o-linear-gradient(top, #5bc0de, #2f96b4);
        background-image: linear-gradient(to bottom, #5bc0de, #2f96b4);
        background-repeat: repeat-x;
        border-color: #2f96b4 #2f96b4 #1f6377;
        border-color: rgba(0, 0, 0, 0.1) rgba(0, 0, 0, 0.1) rgba(0, 0, 0, 0.25);
        filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#ff5bc0de', endColorstr='#ff2f96b4', GradientType=0);
        filter: progid:DXImageTransform.Microsoft.gradient(enabled=false);
        padding: 4px 12px;
        margin-bottom: 0;
        font-size: 14px;
        line-height: 20px;
        box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.2), 0 1px 2px rgba(0, 0, 0, 0.05);
        border-radius: 4px;
}
.menu {
    position: fixed;
    top: 50px;
    height: 40px;
    width: 100%;
    background-color: #fafafa;
    margin-bottom: 20px;
    display: flex;
    justify-content: center;
    border-bottom: 1px gray;
    z-index: 20;
    }
.container{
    max-width: 1000px;
    width: auto;
}
.container-2 {
    margin-right: auto;
    margin-left: auto;
    padding-right: 15px;
    padding-left: 15px;
    width: 100%;
    max-width: 400px;
    display: flex;
    justify-content: space-around;
    margin: auto;
    }
a{
    text-decoration: none;
    color:black;
    cursor: pointer;
    text-shadow: black;
}

a.router-link-exact-active.router-link-active{
    font-weight: bold;
}
</style>
